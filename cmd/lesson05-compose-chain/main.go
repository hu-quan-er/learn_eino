package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
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

	chatTemplate := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}，回答时给出 3 条要点。"),
		schema.UserMessage("{question}"),
	)

	chain := compose.NewChain[map[string]any, string]()
	chain.
		AppendChatTemplate(chatTemplate).
		AppendChatModel(chatModel).
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, message *schema.Message) (string, error) {
			return message.Content, nil
		}))

	runner, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("compile chain failed: %v", err)
	}

	answer, err := runner.Invoke(ctx, map[string]any{
		"role":     "Eino 助教",
		"question": "请简要解释为什么 Eino 要把 Prompt 和 Model 拆成不同组件。",
	})
	if err != nil {
		log.Fatalf("invoke chain failed: %v", err)
	}

	fmt.Println("chain output:")
	fmt.Println(answer)
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required env: %s", key)
	}

	return value
}
