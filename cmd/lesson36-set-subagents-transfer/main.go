package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type closureAgent struct {
	name        string
	description string
	runFn       func(context.Context, *adk.AgentInput, ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent]
}

func main() {
	ctx := context.Background()

	router := &closureAgent{
		name:        "router_agent",
		description: "analyze request then transfer to writer",
		runFn: func(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
			iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

			note := adk.EventFromMessage(
				schema.AssistantMessage("router 已完成初步分析，准备转交 writer_agent。", nil),
				nil,
				schema.Assistant,
				"",
			)
			note.AgentName = "router_agent"
			gen.Send(note)

			gen.Send(&adk.AgentEvent{
				AgentName: "router_agent",
				Action:    adk.NewTransferToAgentAction("writer_agent"),
			})
			gen.Close()
			return iter
		},
	}

	writer := &closureAgent{
		name:        "writer_agent",
		description: "receive rewritten history and produce final answer",
		runFn: func(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
			iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

			message := adk.EventFromMessage(
				schema.AssistantMessage("writer 收到输入："+summarizeMessages(input.Messages), nil),
				nil,
				schema.Assistant,
				"",
			)
			message.AgentName = "writer_agent"
			gen.Send(message)
			gen.Close()
			return iter
		},
	}

	root, err := adk.SetSubAgents(ctx, router, []adk.Agent{writer})
	if err != nil {
		log.Fatalf("set sub agents failed: %v", err)
	}

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: root}).Query(ctx, "请写一段关于 ADK 分层的说明")

	fmt.Println("agent events:")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run transfer demo failed: %v", event.Err)
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

func (a *closureAgent) Name(context.Context) string {
	return a.name
}

func (a *closureAgent) Description(context.Context) string {
	return a.description
}

func (a *closureAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	return a.runFn(ctx, input, options...)
}

func summarizeMessages(messages []adk.Message) string {
	parts := make([]string, 0, len(messages))
	for i, message := range messages {
		if message == nil {
			continue
		}
		parts = append(parts, fmt.Sprintf("[%d]%s=%s", i, message.Role, message.Content))
	}
	return strings.Join(parts, " | ")
}
