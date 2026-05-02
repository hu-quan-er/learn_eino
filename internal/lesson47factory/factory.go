package lesson47factory

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type Dependencies struct {
	SessionValues map[string]any
}

type BuildConfig struct {
	AgentName string
	Prefix    string
}

type RunnerFactory struct {
	deps Dependencies
}

type Service struct {
	runner *adk.Runner
	opts   []adk.AgentRunOption
}

func NewRunnerFactory(deps Dependencies) *RunnerFactory {
	return &RunnerFactory{deps: deps}
}

func (f *RunnerFactory) Build(ctx context.Context, cfg BuildConfig) *Service {
	agent := &factoryAgent{
		name:   cfg.AgentName,
		prefix: cfg.Prefix,
	}

	return &Service{
		runner: adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent}),
		opts: []adk.AgentRunOption{
			adk.WithSessionValues(f.deps.SessionValues),
		},
	}
}

func (s *Service) Query(ctx context.Context, query string) ([]string, error) {
	iter := s.runner.Query(ctx, query, s.opts...)

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

type factoryAgent struct {
	name   string
	prefix string
}

func (a *factoryAgent) Name(context.Context) string {
	return a.name
}

func (a *factoryAgent) Description(context.Context) string {
	return "agent built by lesson47 runner factory"
}

func (a *factoryAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	tenant, _ := adk.GetSessionValue(ctx, "tenant")
	last := ""
	if len(input.Messages) > 0 && input.Messages[len(input.Messages)-1] != nil {
		last = input.Messages[len(input.Messages)-1].Content
	}

	event := adk.EventFromMessage(
		schema.AssistantMessage(fmt.Sprintf("%s tenant=%v query=%s", a.prefix, tenant, last), nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name
	gen.Send(event)
	gen.Close()
	return iter
}
