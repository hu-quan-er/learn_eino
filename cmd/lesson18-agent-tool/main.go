package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type FAQAgent struct{}

func main() {
	ctx := context.Background()

	agentTool := adk.NewAgentTool(ctx, &FAQAgent{})

	toolInfo, err := agentTool.Info(ctx)
	if err != nil {
		log.Fatalf("read agent tool info failed: %v", err)
	}

	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{agentTool},
	})
	if err != nil {
		log.Fatalf("create tools node failed: %v", err)
	}

	assistantMessage := schema.AssistantMessage("", []schema.ToolCall{
		{
			ID:   "call_agent_1",
			Type: "function",
			Function: schema.FunctionCall{
				Name:      toolInfo.Name,
				Arguments: `{"request":"Eino 的 AgentTool 有什么用？"}`,
			},
		},
	})

	toolMessages, err := toolsNode.Invoke(ctx, assistantMessage)
	if err != nil {
		log.Fatalf("invoke agent tool failed: %v", err)
	}

	fmt.Printf("agent tool name: %s\n", toolInfo.Name)
	fmt.Printf("agent tool desc: %s\n\n", toolInfo.Desc)

	for _, message := range toolMessages {
		fmt.Printf("role=%s tool_name=%s content=%s\n", message.Role, message.ToolName, message.Content)
	}
}

func (a *FAQAgent) Name(ctx context.Context) string {
	_ = ctx
	return "faq_agent"
}

func (a *FAQAgent) Description(ctx context.Context) string {
	_ = ctx
	return "answer simple FAQ as a reusable tool"
}

func (a *FAQAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	_ = ctx
	_ = options

	request := ""
	if len(input.Messages) > 0 && input.Messages[len(input.Messages)-1] != nil {
		request = input.Messages[len(input.Messages)-1].Content
	}

	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	gen.Send(&adk.AgentEvent{
		AgentName: a.Name(context.Background()),
		Output: &adk.AgentOutput{
			MessageOutput: &adk.MessageVariant{
				IsStreaming: false,
				Message:     schema.AssistantMessage("AgentTool 内部 agent 的回答："+request, nil),
				Role:        schema.Assistant,
			},
		},
	})
	gen.Close()

	return iter
}
