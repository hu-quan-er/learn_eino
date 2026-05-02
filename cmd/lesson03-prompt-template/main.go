package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	chatTemplate := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}。回答要短、准、清晰。"),
		schema.MessagesPlaceholder("history", true),
		schema.UserMessage("请基于上面的上下文回答：{question}"),
	)

	messages, err := chatTemplate.Format(ctx, map[string]any{
		"role": "Eino 入门老师",
		"history": []*schema.Message{
			schema.UserMessage("我已经学完了 Generate 和 Stream。"),
			{
				Role:    schema.Assistant,
				Content: "很好，下一步应该学会怎么把输入组织成模板。",
			},
		},
		"question": "Prompt Template 在 Eino 里负责什么？",
	})
	if err != nil {
		log.Fatalf("format prompt failed: %v", err)
	}

	fmt.Println("formatted messages:")
	for i, message := range messages {
		fmt.Printf("%d. role=%s\n", i+1, message.Role)
		fmt.Printf("   content=%s\n\n", message.Content)
	}
}
