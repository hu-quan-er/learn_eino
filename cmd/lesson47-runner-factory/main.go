package main

import (
	"context"
	"fmt"
	"log"

	"eino-tutorial/internal/lesson47factory"
)

func main() {
	ctx := context.Background()

	factory := lesson47factory.NewRunnerFactory(lesson47factory.Dependencies{
		SessionValues: map[string]any{
			"tenant": "tutorial-team",
		},
	})

	draftService := factory.Build(ctx, lesson47factory.BuildConfig{
		AgentName: "draft_factory_agent",
		Prefix:    "draft",
	})
	reviewService := factory.Build(ctx, lesson47factory.BuildConfig{
		AgentName: "review_factory_agent",
		Prefix:    "review",
	})

	draftOutputs, err := draftService.Query(ctx, "第 47 课如何讲 factory")
	if err != nil {
		log.Fatalf("run draft service failed: %v", err)
	}
	reviewOutputs, err := reviewService.Query(ctx, "第 47 课如何讲 factory")
	if err != nil {
		log.Fatalf("run review service failed: %v", err)
	}

	fmt.Println("draft service:")
	for _, output := range draftOutputs {
		fmt.Println(output)
	}

	fmt.Println()
	fmt.Println("review service:")
	for _, output := range reviewOutputs {
		fmt.Println(output)
	}
}
