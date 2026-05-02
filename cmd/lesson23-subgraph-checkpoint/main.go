package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type ReviewState struct {
	Reviewer string
}

type inMemoryStore struct {
	data map[string][]byte
}

func init() {
	schema.Register[ReviewState]()
}

func main() {
	ctx := context.Background()
	checkpointID := "lesson23-subgraph-checkpoint"
	store := newInMemoryStore()

	subGraph := compose.NewGraph[string, string](
		compose.WithGenLocalState(func(ctx context.Context) *ReviewState {
			_ = ctx
			return &ReviewState{}
		}),
	)

	err := subGraph.AddLambdaNode("draft", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return input + " -> 初稿", nil
	}))
	if err != nil {
		log.Fatalf("add subgraph draft node failed: %v", err)
	}

	err = subGraph.AddLambdaNode("review", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return input + " -> 已审", nil
	}), compose.WithStatePreHandler(func(ctx context.Context, input string, state *ReviewState) (string, error) {
		_ = ctx
		return input + " [" + strings.TrimSpace(state.Reviewer) + "]", nil
	}))
	if err != nil {
		log.Fatalf("add subgraph review node failed: %v", err)
	}

	err = subGraph.AddEdge(compose.START, "draft")
	if err != nil {
		log.Fatalf("add subgraph start->draft failed: %v", err)
	}
	err = subGraph.AddEdge("draft", "review")
	if err != nil {
		log.Fatalf("add subgraph draft->review failed: %v", err)
	}
	err = subGraph.AddEdge("review", compose.END)
	if err != nil {
		log.Fatalf("add subgraph review->end failed: %v", err)
	}

	graph := compose.NewGraph[string, string]()

	err = graph.AddLambdaNode("prepare", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return strings.TrimSpace(input), nil
	}))
	if err != nil {
		log.Fatalf("add outer prepare node failed: %v", err)
	}

	err = graph.AddGraphNode("content_pipeline", subGraph, compose.WithGraphCompileOptions(
		compose.WithGraphName("lesson23_content_pipeline"),
		compose.WithInterruptAfterNodes([]string{"draft"}),
	))
	if err != nil {
		log.Fatalf("add outer subgraph node failed: %v", err)
	}

	err = graph.AddLambdaNode("publish", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return input + " -> 发布", nil
	}))
	if err != nil {
		log.Fatalf("add outer publish node failed: %v", err)
	}

	err = graph.AddEdge(compose.START, "prepare")
	if err != nil {
		log.Fatalf("add outer start->prepare failed: %v", err)
	}
	err = graph.AddEdge("prepare", "content_pipeline")
	if err != nil {
		log.Fatalf("add outer prepare->content_pipeline failed: %v", err)
	}
	err = graph.AddEdge("content_pipeline", "publish")
	if err != nil {
		log.Fatalf("add outer content_pipeline->publish failed: %v", err)
	}
	err = graph.AddEdge("publish", compose.END)
	if err != nil {
		log.Fatalf("add outer publish->end failed: %v", err)
	}

	runner, err := graph.Compile(ctx,
		compose.WithCheckPointStore(store),
		compose.WithGraphName("lesson23_root"),
	)
	if err != nil {
		log.Fatalf("compile graph failed: %v", err)
	}

	fmt.Println("first run:")
	_, err = runner.Invoke(ctx, "  第 23 课：嵌套 Checkpoint  ", compose.WithCheckPointID(checkpointID))
	if err == nil {
		log.Fatal("expected subgraph interrupt, got nil")
	}

	info, ok := compose.ExtractInterruptInfo(err)
	if !ok {
		log.Fatalf("expected interrupt info, got: %v", err)
	}

	fmt.Printf("subgraphs: %v\n", sortedSubgraphNames(info.SubGraphs))
	for _, interruptCtx := range info.InterruptContexts {
		fmt.Printf("- id=%s address=%s root=%v\n", interruptCtx.ID, interruptCtx.Address.String(), interruptCtx.IsRootCause)
	}

	resumeCtx := compose.ResumeWithData(ctx, rootInterruptID(info), &ReviewState{Reviewer: "Alice"})

	fmt.Println("\nresume from checkpoint:")
	output, err := runner.Invoke(resumeCtx, "第 23 课：嵌套 Checkpoint", compose.WithCheckPointID(checkpointID))
	if err != nil {
		log.Fatalf("resume graph failed: %v", err)
	}

	fmt.Println("subgraph output:")
	fmt.Println(output)
}

func rootInterruptID(info *compose.InterruptInfo) string {
	for _, interruptCtx := range info.InterruptContexts {
		if interruptCtx.IsRootCause {
			return interruptCtx.ID
		}
	}
	if len(info.InterruptContexts) == 0 {
		return ""
	}
	return info.InterruptContexts[0].ID
}

func sortedSubgraphNames(subgraphs map[string]*compose.InterruptInfo) []string {
	names := make([]string, 0, len(subgraphs))
	for name := range subgraphs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		data: make(map[string][]byte),
	}
}

func (s *inMemoryStore) Get(_ context.Context, checkPointID string) ([]byte, bool, error) {
	value, ok := s.data[checkPointID]
	return value, ok, nil
}

func (s *inMemoryStore) Set(_ context.Context, checkPointID string, checkPoint []byte) error {
	s.data[checkPointID] = checkPoint
	return nil
}
