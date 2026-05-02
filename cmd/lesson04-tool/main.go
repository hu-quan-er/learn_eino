package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	toolutils "github.com/cloudwego/eino/components/tool/utils"
)

type WeatherQuery struct {
	City string `json:"city" jsonschema:"required,description=city name"`
	Unit string `json:"unit" jsonschema:"description=temperature unit,enum=celsius,enum=fahrenheit"`
}

type WeatherResult struct {
	Summary     string `json:"summary"`
	Temperature int    `json:"temperature"`
	Unit        string `json:"unit"`
}

func main() {
	ctx := context.Background()

	weatherTool, err := toolutils.InferTool("get_weather", "return mock weather data by city", getWeather)
	if err != nil {
		log.Fatalf("create tool failed: %v", err)
	}

	info, err := weatherTool.Info(ctx)
	if err != nil {
		log.Fatalf("get tool info failed: %v", err)
	}

	schemaJSON, err := info.ToJSONSchema()
	if err != nil {
		log.Fatalf("convert tool schema failed: %v", err)
	}

	prettySchema, err := json.MarshalIndent(schemaJSON, "", "  ")
	if err != nil {
		log.Fatalf("marshal schema failed: %v", err)
	}

	fmt.Printf("tool name: %s\n", info.Name)
	fmt.Printf("tool desc: %s\n\n", info.Desc)
	fmt.Println("inferred input schema:")
	fmt.Println(string(prettySchema))
	fmt.Println()

	result, err := weatherTool.InvokableRun(ctx, `{"city":"Shanghai","unit":"celsius"}`)
	if err != nil {
		log.Fatalf("run tool failed: %v", err)
	}

	fmt.Println("tool result:")
	fmt.Println(result)
}

func getWeather(ctx context.Context, input *WeatherQuery) (*WeatherResult, error) {
	_ = ctx

	temperature := 26
	if input.Unit == "fahrenheit" {
		temperature = 79
	}

	return &WeatherResult{
		Summary:     fmt.Sprintf("%s is sunny today", input.City),
		Temperature: temperature,
		Unit:        input.Unit,
	}, nil
}
