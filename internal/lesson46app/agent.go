package lesson46app

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type prefixAgent struct {
	name   string
	prefix string
}

func newAgent(cfg Config) adk.Agent {
	return &prefixAgent{
		name:   cfg.AgentName,
		prefix: cfg.Prefix,
	}
}

func (a *prefixAgent) Name(context.Context) string {
	return a.name
}

func (a *prefixAgent) Description(context.Context) string {
	return "agent created by lesson46 bootstrap package"
}

func (a *prefixAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	_ = ctx

	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	last := ""
	if len(input.Messages) > 0 && input.Messages[len(input.Messages)-1] != nil {
		last = input.Messages[len(input.Messages)-1].Content
	}

	event := adk.EventFromMessage(
		schema.AssistantMessage(fmt.Sprintf("%s agent handled: %s", a.prefix, last), nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name
	gen.Send(event)
	gen.Close()
	return iter
}
