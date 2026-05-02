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

	graph := compose.NewGraph[string, string]()

	err := graph.AddLambdaNode("normalize", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return strings.TrimSpace(input), nil
	}))
	if err != nil {
		log.Fatalf("add normalize node failed: %v", err)
	}

	err = graph.AddLambdaNode("reply", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return "Graph demo -> " + input, nil
	}))
	if err != nil {
		log.Fatalf("add reply node failed: %v", err)
	}

	err = graph.AddEdge(compose.START, "normalize")
	if err != nil {
		log.Fatalf("add edge start->normalize failed: %v", err)
	}

	err = graph.AddEdge("normalize", "reply")
	if err != nil {
		log.Fatalf("add edge normalize->reply failed: %v", err)
	}

	err = graph.AddEdge("reply", compose.END)
	if err != nil {
		log.Fatalf("add edge reply->end failed: %v", err)
	}

	runner, err := graph.Compile(ctx, compose.WithGraphName("lesson16_graph"))
	if err != nil {
		log.Fatalf("compile graph failed: %v", err)
	}

	output, err := runner.Invoke(ctx, "  什么情况下要直接使用 Graph？  ")
	if err != nil {
		log.Fatalf("invoke graph failed: %v", err)
	}

	fmt.Println("graph output:")
	fmt.Println(output)
}
