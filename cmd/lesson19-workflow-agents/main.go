package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type TimedAgent struct {
	name        string
	description string
	prefix      string
	delay       time.Duration
}

func main() {
	ctx := context.Background()

	researchAgent := &TimedAgent{
		name:        "research_agent",
		description: "collect background information",
		prefix:      "research",
		delay:       200 * time.Millisecond,
	}
	writeAgent := &TimedAgent{
		name:        "write_agent",
		description: "write the final answer",
		prefix:      "write",
		delay:       80 * time.Millisecond,
	}

	sequentialAgent, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        "sequential_demo",
		Description: "run sub agents one by one",
		SubAgents:   []adk.Agent{researchAgent, writeAgent},
	})
	if err != nil {
		log.Fatalf("create sequential agent failed: %v", err)
	}

	fmt.Println("sequential events:")
	consumeEvents(adk.NewRunner(ctx, adk.RunnerConfig{Agent: sequentialAgent}).Query(ctx, "Eino 多 agent 怎么编排？"))

	searchAgent := &TimedAgent{
		name:        "search_agent",
		description: "search in parallel",
		prefix:      "search",
		delay:       220 * time.Millisecond,
	}
	reviewAgent := &TimedAgent{
		name:        "review_agent",
		description: "review in parallel",
		prefix:      "review",
		delay:       100 * time.Millisecond,
	}

	parallelAgent, err := adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
		Name:        "parallel_demo",
		Description: "run sub agents in parallel",
		SubAgents:   []adk.Agent{searchAgent, reviewAgent},
	})
	if err != nil {
		log.Fatalf("create parallel agent failed: %v", err)
	}

	fmt.Println("\nparallel events:")
	consumeEvents(adk.NewRunner(ctx, adk.RunnerConfig{Agent: parallelAgent}).Query(ctx, "Eino 多 agent 怎么编排？"))
}

func (a *TimedAgent) Name(ctx context.Context) string {
	_ = ctx
	return a.name
}

func (a *TimedAgent) Description(ctx context.Context) string {
	_ = ctx
	return a.description
}

func (a *TimedAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	_ = ctx
	_ = options

	lastMessage := ""
	if len(input.Messages) > 0 && input.Messages[len(input.Messages)-1] != nil {
		lastMessage = input.Messages[len(input.Messages)-1].Content
	}

	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	go func() {
		defer gen.Close()
		time.Sleep(a.delay)
		gen.Send(&adk.AgentEvent{
			AgentName: a.name,
			Output: &adk.AgentOutput{
				MessageOutput: &adk.MessageVariant{
					IsStreaming: false,
					Message:     schema.AssistantMessage(fmt.Sprintf("%s -> %s", a.prefix, lastMessage), nil),
					Role:        schema.Assistant,
				},
			},
		})
	}()

	return iter
}

func consumeEvents(iter *adk.AsyncIterator[*adk.AgentEvent]) {
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("agent event failed: %v", event.Err)
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}

		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Fatalf("read agent event message failed: %v", err)
		}
		if message == nil {
			continue
		}

		fmt.Printf("- %s: %s\n", event.AgentName, message.Content)
	}
}
