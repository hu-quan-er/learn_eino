package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/compose"
)

type NormalizeInput struct {
	Question string
}

type NoteInput struct {
	Student string
}

func main() {
	ctx := context.Background()

	workflow := compose.NewWorkflow[map[string]any, map[string]any]()

	workflow.
		AddLambdaNode("normalize_question", compose.InvokableLambda(func(ctx context.Context, input *NormalizeInput) (map[string]any, error) {
			_ = ctx
			return map[string]any{
				"normalized_question": strings.TrimSpace(input.Question),
			}, nil
		})).
		AddInput(compose.START, compose.MapFields("question", "Question"))

	workflow.
		AddLambdaNode("build_note", compose.InvokableLambda(func(ctx context.Context, input *NoteInput) (map[string]any, error) {
			_ = ctx
			return map[string]any{
				"note": fmt.Sprintf("%s is learning Eino workflow.", input.Student),
			}, nil
		})).
		AddInput(compose.START, compose.MapFields("student", "Student"))

	workflow.End().AddInput("normalize_question", compose.MapFields("normalized_question", "normalized_question"))
	workflow.End().AddInput("build_note", compose.MapFields("note", "note"))

	runner, err := workflow.Compile(ctx)
	if err != nil {
		log.Fatalf("compile workflow failed: %v", err)
	}

	output, err := runner.Invoke(ctx, map[string]any{
		"question": "  什么是 Eino Workflow？  ",
		"student":  "Alice",
	})
	if err != nil {
		log.Fatalf("invoke workflow failed: %v", err)
	}

	fmt.Printf("workflow output: %#v\n", output)
}
