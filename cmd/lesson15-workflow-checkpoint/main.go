package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/compose"
)

type inMemoryStore struct {
	data map[string][]byte
}

func main() {
	ctx := context.Background()
	checkpointID := "lesson15-workflow-checkpoint"
	store := newInMemoryStore()

	workflow := compose.NewWorkflow[map[string]any, map[string]any]()

	workflow.
		AddLambdaNode("build_draft", compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			_ = ctx

			topic, _ := input["topic"].(string)
			return map[string]any{
				"draft": "发布草稿：" + strings.TrimSpace(topic),
			}, nil
		})).
		AddInput(compose.START)

	workflow.
		AddLambdaNode("send_notification", compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			_ = ctx

			draft, _ := input["draft"].(string)
			return map[string]any{
				"status": "已发送通知 -> " + draft,
			}, nil
		})).
		AddInput("build_draft")

	workflow.End().AddInput("send_notification", compose.MapFields("status", "status"))

	runner, err := workflow.Compile(ctx,
		compose.WithCheckPointStore(store),
		compose.WithInterruptBeforeNodes([]string{"send_notification"}),
	)
	if err != nil {
		log.Fatalf("compile workflow failed: %v", err)
	}

	fmt.Println("first run:")
	_, err = runner.Invoke(ctx, map[string]any{
		"topic": "今晚 8 点发布 Eino 第 15 课",
	}, compose.WithCheckPointID(checkpointID))
	if err == nil {
		log.Fatal("expected workflow interrupt, got nil")
	}

	info, ok := compose.ExtractInterruptInfo(err)
	if !ok {
		log.Fatalf("expected interrupt info, got: %v", err)
	}

	fmt.Printf("interrupt before nodes: %v\n", info.BeforeNodes)
	for _, interruptCtx := range info.InterruptContexts {
		fmt.Printf("- id=%s address=%s root=%v\n", interruptCtx.ID, interruptCtx.Address.String(), interruptCtx.IsRootCause)
	}

	fmt.Println("\nresume from checkpoint:")
	output, err := runner.Invoke(ctx, nil, compose.WithCheckPointID(checkpointID))
	if err != nil {
		log.Fatalf("resume workflow failed: %v", err)
	}

	fmt.Printf("workflow output: %#v\n", output)
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
