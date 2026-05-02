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

	graph := compose.NewGraph[string, map[string]any]()

	err := graph.AddLambdaNode("summary", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return "概要：" + strings.TrimSpace(input), nil
	}), compose.WithOutputKey("summary"))
	if err != nil {
		log.Fatalf("add summary node failed: %v", err)
	}

	err = graph.AddLambdaNode("examples", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return "示例：给 " + strings.TrimSpace(input) + " 增加 demo", nil
	}), compose.WithOutputKey("examples"))
	if err != nil {
		log.Fatalf("add examples node failed: %v", err)
	}

	err = graph.AddLambdaNode("faq", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return "答疑：围绕 " + strings.TrimSpace(input) + " 整理常见问题", nil
	}), compose.WithOutputKey("faq"))
	if err != nil {
		log.Fatalf("add faq node failed: %v", err)
	}

	err = graph.AddBranch(compose.START, compose.NewGraphMultiBranch(func(ctx context.Context, input string) (map[string]bool, error) {
		_ = ctx

		selected := map[string]bool{
			"summary": true,
		}

		if strings.Contains(input, "实战") {
			selected["examples"] = true
		}
		if strings.Contains(input, "答疑") {
			selected["faq"] = true
		}

		return selected, nil
	}, map[string]bool{
		"summary":  true,
		"examples": true,
		"faq":      true,
	}))
	if err != nil {
		log.Fatalf("add graph multibranch failed: %v", err)
	}

	err = graph.AddEdge("summary", compose.END)
	if err != nil {
		log.Fatalf("add summary->end failed: %v", err)
	}
	err = graph.AddEdge("examples", compose.END)
	if err != nil {
		log.Fatalf("add examples->end failed: %v", err)
	}
	err = graph.AddEdge("faq", compose.END)
	if err != nil {
		log.Fatalf("add faq->end failed: %v", err)
	}

	runner, err := graph.Compile(ctx, compose.WithGraphName("lesson22_graph_multibranch"))
	if err != nil {
		log.Fatalf("compile graph failed: %v", err)
	}

	fmt.Println("case 1:")
	output, err := runner.Invoke(ctx, "Eino 实战")
	if err != nil {
		log.Fatalf("invoke case 1 failed: %v", err)
	}
	fmt.Printf("%#v\n", output)

	fmt.Println("\ncase 2:")
	output, err = runner.Invoke(ctx, "Eino 答疑")
	if err != nil {
		log.Fatalf("invoke case 2 failed: %v", err)
	}
	fmt.Printf("%#v\n", output)

	fmt.Println("\ncase 3:")
	output, err = runner.Invoke(ctx, "Eino 入门")
	if err != nil {
		log.Fatalf("invoke case 3 failed: %v", err)
	}
	fmt.Printf("%#v\n", output)
}
