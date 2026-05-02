package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/eino/compose"
)

func main() {
	ctx := context.Background()

	parallel := compose.NewParallel()
	parallel.
		AddLambda("outline", compose.InvokableLambda(func(ctx context.Context, input map[string]any) (string, error) {
			_ = ctx
			time.Sleep(180 * time.Millisecond)
			topic, _ := input["topic"].(string)
			return "提纲：" + topic, nil
		})).
		AddLambda("keywords", compose.InvokableLambda(func(ctx context.Context, input map[string]any) (string, error) {
			_ = ctx
			time.Sleep(80 * time.Millisecond)
			return "关键词：chain / parallel / output key", nil
		}))

	chain := compose.NewChain[string, string]()
	chain.
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (map[string]any, error) {
			_ = ctx
			return map[string]any{
				"topic": strings.TrimSpace(input),
			}, nil
		})).
		AppendParallel(parallel).
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input map[string]any) (string, error) {
			_ = ctx

			outline, _ := input["outline"].(string)
			keywords, _ := input["keywords"].(string)

			return fmt.Sprintf("parallel output: %s | %s", outline, keywords), nil
		}))

	runner, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("compile chain failed: %v", err)
	}

	startedAt := time.Now()
	output, err := runner.Invoke(ctx, "第 27 课：Chain Parallel")
	if err != nil {
		log.Fatalf("invoke chain failed: %v", err)
	}

	fmt.Println(output)
	fmt.Printf("elapsed: %s\n", time.Since(startedAt).Round(10*time.Millisecond))
}
