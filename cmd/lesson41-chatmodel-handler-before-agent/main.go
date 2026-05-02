package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type PolicyQuery struct {
	Topic string `json:"topic" jsonschema:"required,description=topic name"`
}

type handlerToolModel struct {
	step int
}

type policyInjectionHandler struct {
	*adk.BaseChatModelAgentMiddleware
	tool tool.BaseTool
}

func main() {
	ctx := context.Background()

	policyTool, err := toolutils.InferTool("lookup_policy", "return policy by topic", lookupPolicy)
	if err != nil {
		log.Fatalf("create policy tool failed: %v", err)
	}

	handler := &policyInjectionHandler{
		BaseChatModelAgentMiddleware: &adk.BaseChatModelAgentMiddleware{},
		tool:                         policyTool,
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "lesson41_agent",
		Description: "inject instruction and tool in BeforeAgent",
		Instruction: "你是一个讲解 ADK 的助手。",
		Model:       &handlerToolModel{},
		Handlers:    []adk.ChatModelAgentMiddleware{handler},
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{},
			},
		},
	})
	if err != nil {
		log.Fatalf("create handler agent failed: %v", err)
	}

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent}).Query(ctx, "什么时候该用 handler")

	fmt.Println("agent events:")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run lesson41 failed: %v", event.Err)
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}

		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Fatalf("read message failed: %v", err)
		}
		if strings.TrimSpace(message.Content) == "" && len(message.ToolCalls) == 0 {
			continue
		}
		fmt.Printf("agent=%s role=%s content=%s\n", event.AgentName, event.Output.MessageOutput.Role, message.Content)
	}
}

func (h *policyInjectionHandler) BeforeAgent(ctx context.Context, runCtx *adk.ChatModelAgentContext) (context.Context, *adk.ChatModelAgentContext, error) {
	runCtx.Instruction += "\n\n运行时规则：必须调用 lookup_policy 工具后再回答。"
	runCtx.Tools = append(runCtx.Tools, h.tool)
	return ctx, runCtx, nil
}

func lookupPolicy(ctx context.Context, input *PolicyQuery) (string, error) {
	_ = ctx
	return "policy[" + input.Topic + "] = 先讲接口，再讲状态，再讲事件", nil
}

func (m *handlerToolModel) Generate(_ context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	if m.step == 0 {
		m.step++
		options := model.GetCommonOptions(&model.Options{}, opts...)
		toolName := "lookup_policy"
		if len(options.Tools) > 0 {
			toolName = options.Tools[0].Name
		}
		return schema.AssistantMessage("", []schema.ToolCall{
			{
				ID: "policy_call_1",
				Function: schema.FunctionCall{
					Name:      toolName,
					Arguments: `{"topic":"handler"}`,
				},
			},
		}), nil
	}

	instruction := ""
	last := ""
	for _, message := range input {
		if message == nil {
			continue
		}
		if message.Role == schema.System {
			instruction = message.Content
		}
		last = message.Content
	}

	return schema.AssistantMessage(
		fmt.Sprintf("instruction=%s | tool_result=%s", instruction, last),
		nil,
	), nil
}

func (m *handlerToolModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}

func (m *handlerToolModel) WithTools(_ []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	return m, nil
}
