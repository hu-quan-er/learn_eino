package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/schema"
)

type callbackAgent struct {
	name string
}

type callbackRecorder struct {
	mu         sync.Mutex
	runName    string
	startInput string
	events     []string
	done       chan struct{}
}

func main() {
	ctx := context.Background()
	recorder := &callbackRecorder{done: make(chan struct{})}

	handler := callbacks.NewHandlerBuilder().
		OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
			if info.Component != adk.ComponentOfAgent {
				return ctx
			}
			recorder.mu.Lock()
			recorder.runName = info.Name
			if converted := adk.ConvAgentCallbackInput(input); converted != nil && converted.Input != nil && len(converted.Input.Messages) > 0 {
				recorder.startInput = converted.Input.Messages[0].Content
			}
			recorder.mu.Unlock()
			return ctx
		}).
		OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
			if info.Component != adk.ComponentOfAgent {
				return ctx
			}
			converted := adk.ConvAgentCallbackOutput(output)
			if converted == nil || converted.Events == nil {
				close(recorder.done)
				return ctx
			}

			go func() {
				defer close(recorder.done)
				for {
					event, ok := converted.Events.Next()
					if !ok {
						break
					}
					if event.Output == nil || event.Output.MessageOutput == nil {
						continue
					}
					message, err := event.Output.MessageOutput.GetMessage()
					if err != nil {
						continue
					}
					recorder.mu.Lock()
					recorder.events = append(recorder.events, fmt.Sprintf("%s:%s", event.AgentName, message.Content))
					recorder.mu.Unlock()
				}
			}()
			return ctx
		}).
		Build()

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: &callbackAgent{name: "callback_agent"}}).
		Query(ctx, "lesson44 callback demo", adk.WithCallbacks(handler))

	fmt.Println("runner events:")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run lesson44 failed: %v", event.Err)
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

	<-recorder.done

	recorder.mu.Lock()
	defer recorder.mu.Unlock()

	fmt.Println()
	fmt.Printf("callback run name: %s\n", recorder.runName)
	fmt.Printf("callback start input: %s\n", recorder.startInput)
	fmt.Printf("callback copied events: %v\n", recorder.events)
}

func (a *callbackAgent) Name(context.Context) string {
	return a.name
}

func (a *callbackAgent) Description(context.Context) string {
	return "emit two messages so callback can observe them"
}

func (a *callbackAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	last := ""
	if len(input.Messages) > 0 && input.Messages[0] != nil {
		last = input.Messages[0].Content
	}

	first := adk.EventFromMessage(schema.AssistantMessage("callback step 1: "+last, nil), nil, schema.Assistant, "")
	first.AgentName = a.name
	second := adk.EventFromMessage(schema.AssistantMessage("callback step 2", nil), nil, schema.Assistant, "")
	second.AgentName = a.name

	gen.Send(first)
	gen.Send(second)
	gen.Close()
	return iter
}
