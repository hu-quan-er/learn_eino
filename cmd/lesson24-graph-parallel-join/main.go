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

	graph := compose.NewGraph[string, string]()

	err := graph.AddLambdaNode("outline", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		time.Sleep(180 * time.Millisecond)
		return "大纲：" + strings.TrimSpace(input), nil
	}), compose.WithOutputKey("outline"))
	if err != nil {
		log.Fatalf("add outline node failed: %v", err)
	}

	err = graph.AddLambdaNode("keywords", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		time.Sleep(80 * time.Millisecond)
		return "关键词：graph / fan-in / all_predecessor", nil
	}), compose.WithOutputKey("keywords"))
	if err != nil {
		log.Fatalf("add keywords node failed: %v", err)
	}

	err = graph.AddLambdaNode("merge", compose.InvokableLambda(func(ctx context.Context, input map[string]any) (string, error) {
		_ = ctx

		outline, _ := input["outline"].(string)
		keywords, _ := input["keywords"].(string)

		return fmt.Sprintf("合并结果：%s | %s", outline, keywords), nil
	}))
	if err != nil {
		log.Fatalf("add merge node failed: %v", err)
	}

	err = graph.AddEdge(compose.START, "outline")
	if err != nil {
		log.Fatalf("add start->outline failed: %v", err)
	}
	err = graph.AddEdge(compose.START, "keywords")
	if err != nil {
		log.Fatalf("add start->keywords failed: %v", err)
	}
	err = graph.AddEdge("outline", "merge")
	if err != nil {
		log.Fatalf("add outline->merge failed: %v", err)
	}
	err = graph.AddEdge("keywords", "merge")
	if err != nil {
		log.Fatalf("add keywords->merge failed: %v", err)
	}
	err = graph.AddEdge("merge", compose.END)
	if err != nil {
		log.Fatalf("add merge->end failed: %v", err)
	}

	runner, err := graph.Compile(ctx,
		compose.WithNodeTriggerMode(compose.AllPredecessor),
		compose.WithGraphName("lesson24_graph_parallel_join"),
	)
	if err != nil {
		log.Fatalf("compile graph failed: %v", err)
	}

	startedAt := time.Now()
	output, err := runner.Invoke(ctx, "第 24 课：Graph 并行汇聚")
	if err != nil {
		log.Fatalf("invoke graph failed: %v", err)
	}

	fmt.Println("parallel join output:")
	fmt.Println(output)
	fmt.Printf("elapsed: %s\n", time.Since(startedAt).Round(10*time.Millisecond))
}
