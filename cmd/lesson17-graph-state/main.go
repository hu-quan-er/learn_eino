package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/compose"
)

type LessonState struct {
	Normalized string
	StepCount  int
}

func main() {
	ctx := context.Background()

	graph := compose.NewGraph[string, string](
		compose.WithGenLocalState(func(ctx context.Context) *LessonState {
			_ = ctx
			return &LessonState{}
		}),
	)

	err := graph.AddLambdaNode("normalize", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return strings.ToLower(strings.TrimSpace(input)), nil
	}), compose.WithStatePostHandler(func(ctx context.Context, output string, state *LessonState) (string, error) {
		_ = ctx
		state.Normalized = output
		state.StepCount++
		return output, nil
	}))
	if err != nil {
		log.Fatalf("add normalize node failed: %v", err)
	}

	err = graph.AddLambdaNode("finalize", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return input + " -> ready", nil
	}), compose.WithStatePreHandler(func(ctx context.Context, input string, state *LessonState) (string, error) {
		_ = ctx
		state.StepCount++
		return fmt.Sprintf("[normalized=%s steps=%d] %s", state.Normalized, state.StepCount, input), nil
	}))
	if err != nil {
		log.Fatalf("add finalize node failed: %v", err)
	}

	err = graph.AddEdge(compose.START, "normalize")
	if err != nil {
		log.Fatalf("add edge start->normalize failed: %v", err)
	}

	err = graph.AddEdge("normalize", "finalize")
	if err != nil {
		log.Fatalf("add edge normalize->finalize failed: %v", err)
	}

	err = graph.AddEdge("finalize", compose.END)
	if err != nil {
		log.Fatalf("add edge finalize->end failed: %v", err)
	}

	runner, err := graph.Compile(ctx, compose.WithGraphName("lesson17_graph_state"))
	if err != nil {
		log.Fatalf("compile graph failed: %v", err)
	}

	output, err := runner.Invoke(ctx, "  Eino Graph State  ")
	if err != nil {
		log.Fatalf("invoke graph failed: %v", err)
	}

	fmt.Println("graph state output:")
	fmt.Println(output)
}
