package lesson50app

import (
	"context"

	"github.com/cloudwego/eino/adk"
)

type App struct {
	runner        *adk.Runner
	sessionValues map[string]any
}

func New(ctx context.Context, cfg Config) (*App, error) {
	workflow, err := newWorkflow(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		runner: adk.NewRunner(ctx, adk.RunnerConfig{Agent: workflow}),
		sessionValues: map[string]any{
			"tenant": cfg.Tenant,
		},
	}, nil
}

func (a *App) Run(ctx context.Context, query string) (*Result, error) {
	iter := a.runner.Query(ctx, query, adk.WithSessionValues(a.sessionValues))

	result := &Result{}
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
		result.Trace = append(result.Trace, event.AgentName+": "+message.Content)
		result.Final = message.Content
	}

	return result, nil
}
