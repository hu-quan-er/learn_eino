package main

import (
	"context"
	"fmt"
	"log"

	toolcomponent "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type inMemoryStore struct {
	data map[string][]byte
}

type wrapperTool struct {
	compiledGraph    compose.Runnable[string, string]
	resumeTargetLogs []bool
}

func main() {
	ctx := context.Background()
	checkpointID := "lesson30-tool-composite-interrupt"
	store := newInMemoryStore()

	innerGraph := buildInnerGraph(ctx)
	toolNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []toolcomponent.BaseTool{&wrapperTool{compiledGraph: innerGraph}},
	})
	if err != nil {
		log.Fatalf("create tools node failed: %v", err)
	}

	graph := compose.NewGraph[*schema.Message, string]()
	_ = graph.AddToolsNode("tools", toolNode)
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
		compose.WithGraphName("lesson30_root"),
	)
	if err != nil {
		log.Fatalf("compile outer graph failed: %v", err)
	}

	input := &schema.Message{
		Role: schema.Assistant,
		ToolCalls: []schema.ToolCall{
			{
				ID: "call_1",
				Function: schema.FunctionCall{
					Name:      "wrapper_tool",
					Arguments: "lesson30 input",
				},
			},
		},
	}

	fmt.Println("first run:")
	_, err = runner.Invoke(ctx, input, compose.WithCheckPointID(checkpointID))
	if err == nil {
		log.Fatal("expected composite interrupt, got nil")
	}

	info, ok := compose.ExtractInterruptInfo(err)
	if !ok {
		log.Fatalf("expected interrupt info, got: %v", err)
	}

	root := firstRootCause(info)
	fmt.Printf("root cause: id=%s address=%s info=%v\n", root.ID, root.Address.String(), root.Info)
	fmt.Printf("parent infos: %v\n", collectParentInfos(root))

	resumeCtx := compose.Resume(ctx, root.ID)

	fmt.Println("\nresume from checkpoint:")
	output, err := runner.Invoke(resumeCtx, input, compose.WithCheckPointID(checkpointID))
	if err != nil {
		log.Fatalf("resume outer graph failed: %v", err)
	}

	fmt.Println(output)
}

func buildInnerGraph(ctx context.Context) compose.Runnable[string, string] {
	graph := compose.NewGraph[string, string]()

	_ = graph.AddLambdaNode("pause", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = input

		wasInterrupted, _, _ := compose.GetInterruptState[any](ctx)
		if !wasInterrupted {
			return "", compose.Interrupt(ctx, "inner graph needs resume")
		}

		isResumeTarget, _, _ := compose.GetResumeContext[any](ctx)
		return fmt.Sprintf("inner resumed target=%v", isResumeTarget), nil
	}))

	_ = graph.AddEdge(compose.START, "pause")
	_ = graph.AddEdge("pause", compose.END)

	runner, err := graph.Compile(ctx, compose.WithGraphName("lesson30_inner"))
	if err != nil {
		log.Fatalf("compile inner graph failed: %v", err)
	}

	return runner
}

func (t *wrapperTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	_ = ctx
	return &schema.ToolInfo{
		Name: "wrapper_tool",
		Desc: "wrap a nested graph and forward its interrupt as a composite interrupt",
	}, nil
}

func (t *wrapperTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...toolcomponent.Option) (string, error) {
	_ = opts

	isResumeTarget, _, _ := toolcomponent.GetResumeContext[any](ctx)
	t.resumeTargetLogs = append(t.resumeTargetLogs, isResumeTarget)

	result, err := t.compiledGraph.Invoke(ctx, argumentsInJSON)
	if err != nil {
		if _, ok := compose.ExtractInterruptInfo(err); ok {
			return "", toolcomponent.CompositeInterrupt(ctx, "wrapper tool interrupt", nil, err)
		}
		return "", err
	}

	return result, nil
}

func firstRootCause(info *compose.InterruptInfo) *compose.InterruptCtx {
	for _, interruptCtx := range info.InterruptContexts {
		if interruptCtx.IsRootCause {
			return interruptCtx
		}
	}
	if len(info.InterruptContexts) == 0 {
		return nil
	}
	return info.InterruptContexts[0]
}

func collectParentInfos(ctx *compose.InterruptCtx) []string {
	infos := make([]string, 0)
	for parent := ctx.Parent; parent != nil; parent = parent.Parent {
		if info, ok := parent.Info.(string); ok && info != "" {
			infos = append(infos, info)
		}
	}
	return infos
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
