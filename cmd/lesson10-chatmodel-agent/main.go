package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type WeatherQuery struct {
	City string `json:"city" jsonschema:"required,description=city name"`
}

type WeatherResult struct {
	Summary     string `json:"summary"`
	Temperature int    `json:"temperature"`
	Unit        string `json:"unit"`
}

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

	weatherTool, err := toolutils.InferTool("get_weather", "return mock weather data by city", getWeather)
	if err != nil {
		log.Fatalf("create tool failed: %v", err)
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "weather_agent",
		Description: "answer weather questions with a weather tool",
		Instruction: "你是一个天气助教。必要时调用 get_weather 工具。拿到工具结果后，用中文给出简短答案。",
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{weatherTool},
			},
		},
	})
	if err != nil {
		log.Fatalf("create chat model agent failed: %v", err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})

	iter := runner.Query(ctx, "Shanghai 的天气怎么样？")

	var finalMessage *schema.Message
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("agent run failed: %v", event.Err)
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}

		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Fatalf("read agent event message failed: %v", err)
		}
		if message == nil {
			continue
		}

		fmt.Printf("event role=%s\n", event.Output.MessageOutput.Role)
		fmt.Printf("content=%s\n\n", message.Content)

		if event.Output.MessageOutput.Role == schema.Assistant {
			finalMessage = message
		}
	}

	if finalMessage != nil {
		fmt.Println("final agent answer:")
		fmt.Println(finalMessage.Content)
	}
}

func getWeather(ctx context.Context, input *WeatherQuery) (*WeatherResult, error) {
	_ = ctx

	return &WeatherResult{
		Summary:     fmt.Sprintf("%s is sunny today", input.City),
		Temperature: 26,
		Unit:        "celsius",
	}, nil
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required env: %s", key)
	}

	return value
}
