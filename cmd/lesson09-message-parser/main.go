package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type LessonSummary struct {
	Topic      string   `json:"topic"`
	KeyPoints  []string `json:"key_points"`
	Difficulty int      `json:"difficulty"`
}

func main() {
	ctx := context.Background()

	parser := schema.NewMessageJSONParser[LessonSummary](&schema.MessageJSONParseConfig{
		ParseFrom: schema.MessageParseFromContent,
	})

	chain := compose.NewChain[*schema.Message, LessonSummary]()
	chain.AppendLambda(compose.MessageParser(parser))

	runner, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("compile parser chain failed: %v", err)
	}

	message := &schema.Message{
		Role:    schema.Assistant,
		Content: `{"topic":"ToolsNode","key_points":["executes tool calls","returns tool messages"],"difficulty":2}`,
	}

	summary, err := runner.Invoke(ctx, message)
	if err != nil {
		log.Fatalf("parse message failed: %v", err)
	}

	fmt.Printf("parsed struct: %+v\n", summary)
}
