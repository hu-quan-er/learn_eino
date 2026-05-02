package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type historyAgent struct {
	name        string
	description string
	runFn       func(context.Context, *adk.AgentInput, ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent]
}

func main() {
	ctx := context.Background()

	defaultChild := &historyAgent{
		name:        "history_child_default",
		description: "show default rewritten history",
		runFn:       explainHistoryRunFn("history_child_default"),
	}
	rewrittenChild := adk.AgentWithOptions(ctx, &historyAgent{
		name:        "history_child_compact",
		description: "show compact rewritten history",
		runFn:       explainHistoryRunFn("history_child_compact"),
	}, adk.WithHistoryRewriter(compactHistory))

	runScenario(ctx, "default history", defaultChild)
	fmt.Println()
	runScenario(ctx, "custom history rewriter", rewrittenChild)
}

func runScenario(ctx context.Context, title string, child adk.Agent) {
	router := &historyAgent{
		name:        "history_router",
		description: "emit one note then transfer",
		runFn: func(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
			iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

			event := adk.EventFromMessage(
				schema.AssistantMessage("router 的原始分析：先讲 ADK，再讲 middleware。", nil),
				nil,
				schema.Assistant,
				"",
			)
			event.AgentName = "history_router"
			gen.Send(event)
			gen.Send(&adk.AgentEvent{
				AgentName: "history_router",
				Action:    adk.NewTransferToAgentAction(child.Name(ctx)),
			})
			gen.Close()
			return iter
		},
	}

	root, err := adk.SetSubAgents(ctx, router, []adk.Agent{child})
	if err != nil {
		log.Fatalf("set sub agents failed for %s: %v", title, err)
	}

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: root}).Query(ctx, "请安排第 38 课的讲解顺序")

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

func explainHistoryRunFn(agentName string) func(context.Context, *adk.AgentInput, ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	return func(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
		iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
		event := adk.EventFromMessage(
			schema.AssistantMessage("history="+joinContents(input.Messages), nil),
			nil,
			schema.Assistant,
			"",
		)
		event.AgentName = agentName
		gen.Send(event)
		gen.Close()
		return iter
	}
}

func compactHistory(_ context.Context, entries []*adk.HistoryEntry) ([]adk.Message, error) {
	question := ""
	contexts := make([]string, 0)

	for _, entry := range entries {
		if entry == nil || entry.Message == nil {
			continue
		}
		if entry.IsUserInput && question == "" {
			question = entry.Message.Content
			continue
		}
		contexts = append(contexts, fmt.Sprintf("%s:%s", entry.AgentName, entry.Message.Content))
	}

	return []adk.Message{
		schema.UserMessage("question=" + question),
		schema.UserMessage("compressed_context=" + strings.Join(contexts, " || ")),
	}, nil
}

func (a *historyAgent) Name(context.Context) string {
	return a.name
}

func (a *historyAgent) Description(context.Context) string {
	return a.description
}

func (a *historyAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	return a.runFn(ctx, input, options...)
}

func joinContents(messages []adk.Message) string {
	parts := make([]string, 0, len(messages))
	for _, message := range messages {
		if message == nil {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", message.Role, message.Content))
	}
	return strings.Join(parts, " | ")
}
