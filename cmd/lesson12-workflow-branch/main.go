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

	workflow := compose.NewWorkflow[string, map[string]any]()

	workflow.
		AddLambdaNode("answer_eino", compose.InvokableLambda(func(ctx context.Context, input string) (map[string]any, error) {
			_ = ctx
			return map[string]any{
				"answer": "这是一个 Eino 问题，优先进入框架知识分支。",
			}, nil
		})).
		AddInputWithOptions(compose.START, nil, compose.WithNoDirectDependency())

	workflow.
		AddLambdaNode("answer_general", compose.InvokableLambda(func(ctx context.Context, input string) (map[string]any, error) {
			_ = ctx
			return map[string]any{
				"answer": "这不是 Eino 专属问题，进入通用知识分支。",
			}, nil
		})).
		AddInputWithOptions(compose.START, nil, compose.WithNoDirectDependency())

	workflow.AddBranch(compose.START, compose.NewGraphBranch(func(ctx context.Context, input string) (string, error) {
		_ = ctx

		if strings.Contains(strings.ToLower(input), "eino") {
			return "answer_eino", nil
		}

		return "answer_general", nil
	}, map[string]bool{
		"answer_eino":    true,
		"answer_general": true,
	}))

	workflow.End().AddInput("answer_eino", compose.MapFields("answer", "answer"))
	workflow.End().AddInput("answer_general", compose.MapFields("answer", "general_answer"))

	runner, err := workflow.Compile(ctx)
	if err != nil {
		log.Fatalf("compile workflow failed: %v", err)
	}

	output1, err := runner.Invoke(ctx, "Eino 的 ToolsNode 是做什么的？")
	if err != nil {
		log.Fatalf("invoke workflow route 1 failed: %v", err)
	}

	output2, err := runner.Invoke(ctx, "向量数据库有什么用？")
	if err != nil {
		log.Fatalf("invoke workflow route 2 failed: %v", err)
	}

	fmt.Printf("route 1 output: %#v\n", output1)
	fmt.Printf("route 2 output: %#v\n", output2)
}
