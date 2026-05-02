package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type streamingLessonModel struct {
	chunks []string
}

func main() {
	ctx := context.Background()

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "stream_agent",
		Description: "streaming agent demo",
		Instruction: "你是一个简洁的 Eino 助手。",
		Model: &streamingLessonModel{
			chunks: []string{
				"Eino ",
				"Agent ",
				"也可以流式输出。",
			},
		},
	})
	if err != nil {
		log.Fatalf("create chat model agent failed: %v", err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})

	iter := runner.Query(ctx, "请用一句话说明 Eino Agent 的流式输出")

	fmt.Println("agent events:")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("agent run failed: %v", event.Err)
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}

		output := event.Output.MessageOutput
		if !output.IsStreaming {
			message, err := output.GetMessage()
			if err != nil {
				log.Fatalf("read message failed: %v", err)
			}
			fmt.Printf("role=%s content=%s\n", output.Role, message.Content)
			continue
		}

		fmt.Printf("role=%s stream=", output.Role)

		var chunks []*schema.Message
		func() {
			defer output.MessageStream.Close()

			for {
				chunk, err := output.MessageStream.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					log.Fatalf("recv message chunk failed: %v", err)
				}

				fmt.Print(chunk.Content)
				chunks = append(chunks, chunk)
			}
		}()

		fmt.Println()

		finalMessage, err := schema.ConcatMessages(chunks)
		if err != nil {
			log.Fatalf("concat message chunks failed: %v", err)
		}

		fmt.Printf("final=%s\n", finalMessage.Content)
	}
}

func (m *streamingLessonModel) Generate(_ context.Context, _ []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	return schema.AssistantMessage(strings.Join(m.chunks, ""), nil), nil
}

func (m *streamingLessonModel) Stream(_ context.Context, _ []*schema.Message, _ ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	reader, writer := schema.Pipe[*schema.Message](0)

	go func() {
		defer writer.Close()
		for _, chunk := range m.chunks {
			writer.Send(schema.AssistantMessage(chunk, nil), nil)
			time.Sleep(120 * time.Millisecond)
		}
	}()

	return reader, nil
}
