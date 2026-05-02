package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/compose"
)

func main() {
	ctx := context.Background()

	chain := compose.NewChain[string, string]()
	chain.
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			_ = ctx
			return strings.TrimSpace(input), nil
		})).
		AppendGraph(buildDraftGraph()).
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			_ = ctx
			return input + " -> 已发布", nil
		}))

	runner, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("compile chain failed: %v", err)
	}

	output, err := runner.Invoke(ctx, "  第 26 课：Chain AppendGraph  ")
	if err != nil {
		log.Fatalf("invoke chain failed: %v", err)
	}

	fmt.Println("chain output:")
	fmt.Println(output)
}

func buildDraftGraph() *compose.Graph[string, string] {
	graph := compose.NewGraph[string, string]()

	_ = graph.AddLambdaNode("outline", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return "提纲：" + input, nil
	}))

	_ = graph.AddLambdaNode("write", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return input + " -> 初稿完成", nil
	}))

	_ = graph.AddEdge(compose.START, "outline")
	_ = graph.AddEdge("outline", "write")
	_ = graph.AddEdge("write", compose.END)

	return graph
}
