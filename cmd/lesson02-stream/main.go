package main

import (
	"context"
	"errors"
	"fmt"
	"io"
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
			Content: "请分 3 条简短说明 Generate 和 Stream 的区别。",
		},
	}

	stream, err := chatModel.Stream(ctx, messages, model.WithTemperature(0.2))
	if err != nil {
		log.Fatalf("stream failed: %v", err)
	}
	defer stream.Close()

	fmt.Println("assistant (streaming):")

	chunks := make([]*schema.Message, 0, 16)
	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("stream recv failed: %v", err)
		}

		chunks = append(chunks, chunk)
		fmt.Print(chunk.Content)
	}

	fmt.Println()
	fmt.Println()

	fullMessage, err := schema.ConcatMessages(chunks)
	if err != nil {
		log.Fatalf("concat streamed chunks failed: %v", err)
	}

	fmt.Println("assistant (merged):")
	fmt.Println(fullMessage.Content)
	fmt.Printf("\nreceived %d chunks\n", len(chunks))
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required env: %s", key)
	}

	return value
}
