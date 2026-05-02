package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type skipTransferAgent struct {
	name        string
	description string
	runFn       func(context.Context, *adk.AgentInput, ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent]
}

func main() {
	ctx := context.Background()

	runCase(ctx, "default history with transfer messages")
	fmt.Println()
	runCase(ctx, "skip transfer messages", adk.WithSkipTransferMessages())
}

func runCase(ctx context.Context, title string, opts ...adk.AgentRunOption) {
	router := &skipTransferAgent{
		name:        "skip_router",
		description: "emit note and transfer",
		runFn: func(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
			iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

			event := adk.EventFromMessage(
				schema.AssistantMessage("router note: 这段消息默认会进入下游 history。", nil),
				nil,
				schema.Assistant,
				"",
			)
			event.AgentName = "skip_router"
			gen.Send(event)
			gen.Send(&adk.AgentEvent{
				AgentName: "skip_router",
				Action:    adk.NewTransferToAgentAction("skip_child"),
			})
			gen.Close()
			return iter
		},
	}

	child := &skipTransferAgent{
		name:        "skip_child",
		description: "inspect messages it received",
		runFn: func(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
			iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
			event := adk.EventFromMessage(
				schema.AssistantMessage(fmt.Sprintf("child inputs(%d)=%s", len(input.Messages), summarizeTransferInputs(input.Messages)), nil),
				nil,
				schema.Assistant,
				"",
			)
			event.AgentName = "skip_child"
			gen.Send(event)
			gen.Close()
			return iter
		},
	}

	root, err := adk.SetSubAgents(ctx, router, []adk.Agent{child})
	if err != nil {
		log.Fatalf("set sub agents failed: %v", err)
	}

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: root}).Query(ctx, "lesson39 要强调什么", opts...)

	fmt.Println(title + ":")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run %s failed: %v", title, event.Err)
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}

		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Fatalf("read message failed: %v", err)
		}
		fmt.Printf("agent=%s content=%s\n", event.AgentName, message.Content)
	}
}

func (a *skipTransferAgent) Name(context.Context) string {
	return a.name
}

func (a *skipTransferAgent) Description(context.Context) string {
	return a.description
}

func (a *skipTransferAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	return a.runFn(ctx, input, options...)
}

func summarizeTransferInputs(messages []adk.Message) string {
	parts := make([]string, 0, len(messages))
	for i, message := range messages {
		if message == nil {
			continue
		}
		parts = append(parts, fmt.Sprintf("[%d]%s=%s", i, message.Role, message.Content))
	}
	return strings.Join(parts, " | ")
}
