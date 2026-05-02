package lesson48testing

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type ReviewAgent struct {
	name   string
	prefix string
}

func NewReviewAgent(name, prefix string) *ReviewAgent {
	return &ReviewAgent{name: name, prefix: prefix}
}

func (a *ReviewAgent) Name(context.Context) string {
	return a.name
}

func (a *ReviewAgent) Description(context.Context) string {
	return "review agent used in lesson48 tests"
}

func (a *ReviewAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	_ = ctx

	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	last := ""
	if len(input.Messages) > 0 && input.Messages[len(input.Messages)-1] != nil {
		last = input.Messages[len(input.Messages)-1].Content
	}

	event := adk.EventFromMessage(
		schema.AssistantMessage(fmt.Sprintf("%s:%s", a.prefix, last), nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name
	gen.Send(event)
	gen.Close()
	return iter
}

func RunOnce(ctx context.Context, name, prefix, query string) ([]string, error) {
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: NewReviewAgent(name, prefix),
	})
	iter := runner.Query(ctx, query)

	var outputs []string
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			return nil, event.Err
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}

		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, message.Content)
	}

	return outputs, nil
}
