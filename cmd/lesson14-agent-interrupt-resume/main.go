package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type ApprovalInput struct {
	Action string `json:"action" jsonschema:"required,description=action waiting for approval"`
}

type inMemoryStore struct {
	data map[string][]byte
}

type scriptedToolCallModel struct {
	step int
}

func main() {
	ctx := context.Background()
	checkpointID := "lesson14-agent-interrupt"

	approvalTool, err := toolutils.InferTool("request_approval", "pause until resume data arrives", requestApproval)
	if err != nil {
		log.Fatalf("create approval tool failed: %v", err)
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "approval_agent",
		Description: "agent interrupt and resume demo",
		Instruction: "你负责发布任务。调用 request_approval 拿到审批结果后，再给出最终结论。",
		Model:       &scriptedToolCallModel{},
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{approvalTool},
			},
		},
	})
	if err != nil {
		log.Fatalf("create chat model agent failed: %v", err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		CheckPointStore: newInMemoryStore(),
	})

	fmt.Println("first run:")
	iter := runner.Query(ctx, "请发布第 14 课", adk.WithCheckPointID(checkpointID))
	rootCauseID := consumeAgentEvents(iter)
	if rootCauseID == "" {
		log.Fatal("did not receive interrupt root cause id")
	}

	fmt.Println("\nresume with decision=approved:")
	iter, err = runner.ResumeWithParams(ctx, checkpointID, &adk.ResumeParams{
		Targets: map[string]any{
			rootCauseID: "approved",
		},
	})
	if err != nil {
		log.Fatalf("resume agent failed: %v", err)
	}

	consumeAgentEvents(iter)
}

func requestApproval(ctx context.Context, input *ApprovalInput) (string, error) {
	wasInterrupted, hasState, state := tool.GetInterruptState[string](ctx)
	if !wasInterrupted {
		return "", tool.StatefulInterrupt(ctx, fmt.Sprintf("需要审批：%s", input.Action), "pending:"+input.Action)
	}

	isResumeTarget, hasData, decision := tool.GetResumeContext[string](ctx)
	if !isResumeTarget {
		if hasState {
			return "", tool.StatefulInterrupt(ctx, fmt.Sprintf("继续等待审批：%s", state), state)
		}
		return "", tool.StatefulInterrupt(ctx, fmt.Sprintf("继续等待审批：%s", input.Action), "pending:"+input.Action)
	}

	if !hasData {
		return "", fmt.Errorf("missing resume data for approval tool")
	}

	return fmt.Sprintf("%s -> %s", state, decision), nil
}

func consumeAgentEvents(iter *adk.AsyncIterator[*adk.AgentEvent]) string {
	rootCauseID := ""

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("agent event failed: %v", event.Err)
		}

		if event.Output != nil && event.Output.MessageOutput != nil {
			message, err := event.Output.MessageOutput.GetMessage()
			if err != nil {
				log.Fatalf("read event message failed: %v", err)
			}
			if message != nil {
				fmt.Printf("event role=%s content=%s\n", event.Output.MessageOutput.Role, message.Content)
			}
		}

		if event.Action != nil && event.Action.Interrupted != nil {
			fmt.Println("interrupt contexts:")
			for _, interruptCtx := range event.Action.Interrupted.InterruptContexts {
				fmt.Printf("- id=%s root=%v info=%v\n", interruptCtx.ID, interruptCtx.IsRootCause, interruptCtx.Info)
				if interruptCtx.IsRootCause {
					rootCauseID = interruptCtx.ID
				}
			}
		}
	}

	return rootCauseID
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

func (m *scriptedToolCallModel) Generate(_ context.Context, input []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	switch m.step {
	case 0:
		m.step++
		return schema.AssistantMessage("", []schema.ToolCall{
			{
				ID:   "approval_call_1",
				Type: "function",
				Function: schema.FunctionCall{
					Name:      "request_approval",
					Arguments: `{"action":"发布第 14 课"}`,
				},
			},
		}), nil
	case 1:
		m.step++
		lastMessage := input[len(input)-1]
		return schema.AssistantMessage("最终结果："+lastMessage.Content, nil), nil
	default:
		return schema.AssistantMessage("流程结束", nil), nil
	}
}

func (m *scriptedToolCallModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}

	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}
