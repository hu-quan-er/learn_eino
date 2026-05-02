package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type LoopDemoAgent struct {
	name        string
	description string
	prefix      string
	breakAfter  int
	current     int
}

func main() {
	ctx := context.Background()

	plainLoop, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
		Name:        "plain_loop",
		Description: "run the same agent for a fixed number of iterations",
		SubAgents: []adk.Agent{
			&LoopDemoAgent{
				name:        "loop_worker",
				description: "emit one message every iteration",
				prefix:      "plain iteration",
			},
		},
		MaxIterations: 3,
	})
	if err != nil {
		log.Fatalf("create plain loop agent failed: %v", err)
	}

	fmt.Println("plain loop:")
	consumeEvents(adk.NewRunner(ctx, adk.RunnerConfig{Agent: plainLoop}).Query(ctx, "第 21 课：LoopAgent"))

	breakLoop, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
		Name:        "break_loop",
		Description: "stop the loop from inside a sub agent",
		SubAgents: []adk.Agent{
			&LoopDemoAgent{
				name:        "break_worker",
				description: "request break after the second iteration",
				prefix:      "break iteration",
				breakAfter:  2,
			},
		},
		MaxIterations: 5,
	})
	if err != nil {
		log.Fatalf("create break loop agent failed: %v", err)
	}

	fmt.Println("\nbreak loop:")
	consumeEvents(adk.NewRunner(ctx, adk.RunnerConfig{Agent: breakLoop}).Query(ctx, "第 21 课：BreakLoop"))
}

func (a *LoopDemoAgent) Name(ctx context.Context) string {
	_ = ctx
	return a.name
}

func (a *LoopDemoAgent) Description(ctx context.Context) string {
	_ = ctx
	return a.description
}

func (a *LoopDemoAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	_ = ctx
	_ = input
	_ = options

	a.current++

	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	go func() {
		defer gen.Close()

		event := &adk.AgentEvent{
			AgentName: a.name,
			Output: &adk.AgentOutput{
				MessageOutput: &adk.MessageVariant{
					IsStreaming: false,
					Message:     schema.AssistantMessage(fmt.Sprintf("%s #%d", a.prefix, a.current), nil),
					Role:        schema.Assistant,
				},
			},
		}

		if a.breakAfter > 0 && a.current >= a.breakAfter {
			event.Action = adk.NewBreakLoopAction(a.name)
		}

		gen.Send(event)
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
			log.Fatalf("read message failed: %v", err)
		}
		if message == nil {
			continue
		}

		suffix := ""
		if event.Action != nil && event.Action.BreakLoop != nil {
			suffix = fmt.Sprintf(" [break done=%v current_iteration=%d]", event.Action.BreakLoop.Done, event.Action.BreakLoop.CurrentIterations)
		}

		fmt.Printf("- %s: %s%s\n", event.AgentName, message.Content, suffix)
	}
}
