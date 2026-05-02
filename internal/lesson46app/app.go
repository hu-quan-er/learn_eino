package lesson46app

import (
	"context"

	"github.com/cloudwego/eino/adk"
)

type App struct {
	runner *adk.Runner
}

func New(ctx context.Context, cfg Config) *App {
	return &App{
		runner: adk.NewRunner(ctx, adk.RunnerConfig{
			Agent: newAgent(cfg),
		}),
	}
}

func (a *App) Run(ctx context.Context, query string) ([]string, error) {
	iter := a.runner.Query(ctx, query)

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
