package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type customAgent struct {
	name string
}

func main() {
	ctx := context.Background()
	agent := &customAgent{name: "lesson31_custom_agent"}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent})

	iter := runner.Query(ctx, "请解释自定义 Agent 的最小实现")

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

		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Fatalf("read message failed: %v", err)
		}

		fmt.Printf("agent=%s role=%s content=%s\n", event.AgentName, event.Output.MessageOutput.Role, message.Content)
	}
}

func (a *customAgent) Name(context.Context) string {
	return a.name
}

func (a *customAgent) Description(context.Context) string {
	return "a minimal custom agent without any model dependency"
}

func (a *customAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	content := fmt.Sprintf(
		"messages=%d streaming=%v last_user=%q",
		len(input.Messages),
		input.EnableStreaming,
		lastUserContent(input.Messages),
	)

	event := adk.EventFromMessage(
		schema.AssistantMessage(content, nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name

	gen.Send(event)
	gen.Close()
	return iter
}

func lastUserContent(messages []adk.Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i] != nil && messages[i].Role == schema.User {
			return messages[i].Content
		}
	}

	if len(messages) == 0 || messages[len(messages)-1] == nil {
		return ""
	}

	return messages[len(messages)-1].Content
}
