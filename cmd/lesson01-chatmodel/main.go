package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  mustEnv("OPENAI_API_KEY"),
		Model:   mustEnv("OPENAI_MODEL"),
		BaseURL: os.Getenv("OPENAI_BASE_URL"),
	})
	if err != nil {
		log.Fatalf("create chat model failed: %v", err)
	}

	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: "你是一个讲解 Eino 的入门老师，回答要短、准、清晰。",
		},
		{
			Role:    schema.User,
			Content: "请用一句话解释什么是 Eino。",
		},
	}

	resp, err := chatModel.Generate(ctx, messages, model.WithTemperature(0.2))
	if err != nil {
		log.Fatalf("generate failed: %v", err)
	}

	fmt.Println("assistant:")
	fmt.Println(resp.Content)
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required env: %s", key)
	}

	return value
}
