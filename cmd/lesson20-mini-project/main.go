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

	workflow := buildProjectWorkflow()

	runner, err := workflow.Compile(ctx)
	if err != nil {
		log.Fatalf("compile mini project workflow failed: %v", err)
	}

	output, err := runner.Invoke(ctx, map[string]any{
		"topic":    "  第 20 课：完整项目骨架  ",
		"audience": "初学者",
	})
	if err != nil {
		log.Fatalf("invoke mini project workflow failed: %v", err)
	}

	fmt.Printf("mini project output: %#v\n", output)
}

func buildProjectWorkflow() *compose.Workflow[map[string]any, map[string]any] {
	workflow := compose.NewWorkflow[map[string]any, map[string]any]()

	workflow.
		AddLambdaNode("prepare_request", compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			_ = ctx

			topic, _ := input["topic"].(string)
			audience, _ := input["audience"].(string)

			return map[string]any{
				"topic":    strings.TrimSpace(topic),
				"audience": audience,
			}, nil
		})).
		AddInput(compose.START)

	workflow.
		AddGraphNode("draft_pipeline", buildDraftGraph(), compose.WithGraphCompileOptions(compose.WithGraphName("draft_pipeline"))).
		AddInput("prepare_request", compose.FromField("topic"))

	workflow.
		AddLambdaNode("review", compose.InvokableLambda(func(ctx context.Context, input string) (map[string]any, error) {
			_ = ctx

			return map[string]any{
				"draft":       input,
				"review_note": "审阅通过，可以发布",
			}, nil
		})).
		AddInput("draft_pipeline")

	workflow.
		AddLambdaNode("summary", compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			_ = ctx

			draft, _ := input["draft"].(string)
			reviewNote, _ := input["review_note"].(string)

			return map[string]any{
				"summary": fmt.Sprintf("项目骨架完成：%s | %s", draft, reviewNote),
			}, nil
		})).
		AddInput("review")

	workflow.End().AddInput("review", compose.MapFields("draft", "draft"))
	workflow.End().AddInput("summary", compose.MapFields("summary", "summary"))

	return workflow
}

func buildDraftGraph() *compose.Graph[string, string] {
	graph := compose.NewGraph[string, string]()

	_ = graph.AddLambdaNode("plan", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return "提纲：" + input, nil
	}))

	_ = graph.AddLambdaNode("write", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return input + " -> 初稿完成", nil
	}))

	_ = graph.AddEdge(compose.START, "plan")
	_ = graph.AddEdge("plan", "write")
	_ = graph.AddEdge("write", compose.END)

	return graph
}
