package main

import (
	"context"
	"fmt"
	"log"

	"eino-tutorial/internal/lesson48testing"
)

func main() {
	ctx := context.Background()

	outputs, err := lesson48testing.RunOnce(ctx, "lesson48_demo_agent", "review", "怎么给 agent 写测试")
	if err != nil {
		log.Fatalf("run lesson48 demo failed: %v", err)
	}

	fmt.Println("lesson48 sample run:")
	for _, output := range outputs {
		fmt.Println(output)
	}

	fmt.Println()
	fmt.Println("test command:")
	fmt.Println("go test ./internal/lesson48testing")
}
