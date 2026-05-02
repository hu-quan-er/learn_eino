package main

import (
	"context"
	"fmt"
	"log"

	"eino-tutorial/internal/lesson46app"
)

func main() {
	ctx := context.Background()

	app := lesson46app.New(ctx, lesson46app.DefaultConfig())
	outputs, err := app.Run(ctx, "为什么工程化示例要先做 bootstrap")
	if err != nil {
		log.Fatalf("run lesson46 app failed: %v", err)
	}

	fmt.Println("lesson46 outputs:")
	for _, output := range outputs {
		fmt.Println(output)
	}
}
