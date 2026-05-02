package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type streamingAgent struct {
	name string
}

func main() {
	ctx := context.Background()
	agent := &streamingAgent{name: "lesson33_streaming_agent"}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})

	iter := runner.Query(ctx, "演示自定义流式 Agent")

	fmt.Println("agent events:")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run agent failed: %v", event.Err)
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
			fmt.Printf("agent=%s role=%s content=%s\n", event.AgentName, output.Role, message.Content)
			continue
		}

		fmt.Printf("agent=%s role=%s stream=", event.AgentName, output.Role)

		var chunks []*schema.Message
		func() {
			defer output.MessageStream.Close()

			for {
				chunk, err := output.MessageStream.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					log.Fatalf("recv stream chunk failed: %v", err)
				}

				fmt.Print(chunk.Content)
				chunks = append(chunks, chunk)
			}
		}()

		fmt.Println()

		finalMessage, err := schema.ConcatMessages(chunks)
		if err != nil {
			log.Fatalf("concat stream chunks failed: %v", err)
		}

		fmt.Printf("final=%s\n", finalMessage.Content)
	}
}

func (a *streamingAgent) Name(context.Context) string {
	return a.name
}

func (a *streamingAgent) Description(context.Context) string {
	return "a custom agent that emits MessageStream directly"
}

func (a *streamingAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	if !input.EnableStreaming {
		event := adk.EventFromMessage(
			schema.AssistantMessage("runner 没有开启 streaming", nil),
			nil,
			schema.Assistant,
			"",
		)
		event.AgentName = a.name
		gen.Send(event)
		gen.Close()
		return iter
	}

	stream := schema.StreamReaderFromArray([]*schema.Message{
		schema.AssistantMessage("自定义 ", nil),
		schema.AssistantMessage("streaming ", nil),
		schema.AssistantMessage("agent", nil),
	})
	stream.SetAutomaticClose()

	event := adk.EventFromMessage(nil, stream, schema.Assistant, "")
	event.AgentName = a.name

	gen.Send(event)
	gen.Close()
	return iter
}
