package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type ApprovalState struct {
	Arguments string
}

type inMemoryStore struct {
	data map[string][]byte
}

type approvalTool struct{}

func init() {
	schema.Register[ApprovalState]()
}

func main() {
	ctx := context.Background()
	checkpointID := "lesson29-tool-interrupt-resume"
	store := newInMemoryStore()

	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{&approvalTool{}},
	})
	if err != nil {
		log.Fatalf("create tools node failed: %v", err)
	}

	graph := compose.NewGraph[*schema.Message, string]()
	_ = graph.AddToolsNode("tools", toolsNode)
	_ = graph.AddLambdaNode("collect", compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (string, error) {
		_ = ctx
		if len(input) == 0 {
			return "", nil
		}
		return input[0].Content, nil
	}))
	_ = graph.AddEdge(compose.START, "tools")
	_ = graph.AddEdge("tools", "collect")
	_ = graph.AddEdge("collect", compose.END)

	runner, err := graph.Compile(ctx,
		compose.WithCheckPointStore(store),
		compose.WithGraphName("lesson29_root"),
	)
	if err != nil {
		log.Fatalf("compile graph failed: %v", err)
	}

	input := &schema.Message{
		Role: schema.Assistant,
		ToolCalls: []schema.ToolCall{
			{
				ID: "call_1",
				Function: schema.FunctionCall{
					Name:      "approval_tool",
					Arguments: "发布第 29 课",
				},
			},
		},
	}

	fmt.Println("first run:")
	_, err = runner.Invoke(ctx, input, compose.WithCheckPointID(checkpointID))
	if err == nil {
		log.Fatal("expected tool interrupt, got nil")
	}

	info, ok := compose.ExtractInterruptInfo(err)
	if !ok {
		log.Fatalf("expected interrupt info, got: %v", err)
	}

	for _, interruptCtx := range info.InterruptContexts {
		fmt.Printf("- id=%s address=%s info=%v root=%v\n", interruptCtx.ID, interruptCtx.Address.String(), interruptCtx.Info, interruptCtx.IsRootCause)
	}

	resumeCtx := compose.ResumeWithData(ctx, rootInterruptID(info), "approved by reviewer")

	fmt.Println("\nresume from checkpoint:")
	output, err := runner.Invoke(resumeCtx, input, compose.WithCheckPointID(checkpointID))
	if err != nil {
		log.Fatalf("resume graph failed: %v", err)
	}

	fmt.Println(output)
}

func (t *approvalTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	_ = ctx
	return &schema.ToolInfo{
		Name: "approval_tool",
		Desc: "interrupt once and wait for an approval decision",
	}, nil
}

func (t *approvalTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	_ = opts

	callID := compose.GetToolCallID(ctx)
	wasInterrupted, hasState, state := tool.GetInterruptState[ApprovalState](ctx)
	isResumeTarget, hasData, decision := tool.GetResumeContext[string](ctx)

	if !wasInterrupted {
		return "", tool.StatefulInterrupt(ctx, "need approval for "+callID, ApprovalState{Arguments: argumentsInJSON})
	}

	if !hasState {
		return "", fmt.Errorf("approval tool expected saved state for %s", callID)
	}

	if !isResumeTarget {
		return "", tool.StatefulInterrupt(ctx, "waiting for resume "+callID, state)
	}

	if hasData && decision != "" {
		return fmt.Sprintf("approved(%s): %s", decision, state.Arguments), nil
	}

	return "approved(default): " + state.Arguments, nil
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
