package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
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

	toolInfo, err := weatherTool.Info(ctx)
	if err != nil {
		log.Fatalf("get tool info failed: %v", err)
	}

	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{weatherTool},
	})
	if err != nil {
		log.Fatalf("create tools node failed: %v", err)
	}

	history := []*schema.Message{
		{
			Role: schema.System,
			Content: "你是一个天气助教。" +
				"当用户问天气时，先调用 get_weather 工具。" +
				"拿到工具结果后，再用中文给出简短总结，不要重复调用工具。",
		},
		{
			Role:    schema.User,
			Content: "请告诉我 Shanghai 的天气。",
		},
	}

	firstResp, err := chatModel.Generate(
		ctx,
		history,
		model.WithTools([]*schema.ToolInfo{toolInfo}),
		model.WithToolChoice(schema.ToolChoiceForced),
	)
	if err != nil {
		log.Fatalf("first generate failed: %v", err)
	}

	fmt.Println("first model response:")
	fmt.Printf("tool calls: %d\n", len(firstResp.ToolCalls))
	for i, toolCall := range firstResp.ToolCalls {
		fmt.Printf("%d. tool=%s args=%s\n", i+1, toolCall.Function.Name, toolCall.Function.Arguments)
	}
	fmt.Println()

	toolMessages, err := toolsNode.Invoke(ctx, firstResp)
	if err != nil {
		log.Fatalf("execute tool calls failed: %v", err)
	}

	history = append(history, firstResp)
	history = append(history, toolMessages...)

	finalResp, err := chatModel.Generate(ctx, history)
	if err != nil {
		log.Fatalf("final generate failed: %v", err)
	}

	fmt.Println("final answer:")
	fmt.Println(finalResp.Content)
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
