package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type WeatherQuery struct {
	City string `json:"city" jsonschema:"required,description=city name"`
}

type WeatherResult struct {
	Summary string `json:"summary"`
}

func main() {
	ctx := context.Background()

	weatherTool, err := toolutils.InferTool("get_weather", "return mock weather data by city", getWeather)
	if err != nil {
		log.Fatalf("create tool failed: %v", err)
	}

	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{weatherTool},
	})
	if err != nil {
		log.Fatalf("create tools node failed: %v", err)
	}

	assistantMessage := schema.AssistantMessage("", []schema.ToolCall{
		{
			ID:   "call_weather_1",
			Type: "function",
			Function: schema.FunctionCall{
				Name:      "get_weather",
				Arguments: `{"city":"Shanghai"}`,
			},
		},
	})

	toolMessages, err := toolsNode.Invoke(ctx, assistantMessage)
	if err != nil {
		log.Fatalf("invoke tools node failed: %v", err)
	}

	fmt.Println("tool messages:")
	for i, message := range toolMessages {
		fmt.Printf("%d. role=%s tool_name=%s tool_call_id=%s\n", i+1, message.Role, message.ToolName, message.ToolCallID)
		fmt.Printf("   content=%s\n\n", message.Content)
	}
}

func getWeather(ctx context.Context, input *WeatherQuery) (*WeatherResult, error) {
	_ = ctx

	return &WeatherResult{
		Summary: fmt.Sprintf("%s is sunny today", input.City),
	}, nil
}
