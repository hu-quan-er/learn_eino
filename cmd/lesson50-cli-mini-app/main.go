package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"eino-tutorial/internal/lesson50app"
)

func main() {
	ctx := context.Background()
	query := "请生成第 50 课的小项目说明"
	if len(os.Args) > 1 {
		query = strings.Join(os.Args[1:], " ")
	}

	app, err := lesson50app.New(ctx, lesson50app.DefaultConfig())
	if err != nil {
		log.Fatalf("create lesson50 app failed: %v", err)
	}

	result, err := app.Run(ctx, query)
	if err != nil {
		log.Fatalf("run lesson50 app failed: %v", err)
	}

	fmt.Println("trace:")
	for _, trace := range result.Trace {
		fmt.Println(trace)
	}

	fmt.Println()
	fmt.Println("final:")
	fmt.Println(result.Final)
}
