package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type outerToolModel struct {
	step int
}

type innerAgent struct {
	name string
}

func main() {
	ctx := context.Background()

	runCase(ctx, "EmitInternalEvents=false", false)
	fmt.Println()
	runCase(ctx, "EmitInternalEvents=true", true)
}

func runCase(ctx context.Context, title string, emit bool) {
	inner := &innerAgent{name: "inner_agent"}
	innerTool := adk.NewAgentTool(ctx, inner)

	outer, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "outer_agent",
		Description: "call inner agent tool once",
		Instruction: "你会把问题交给 inner agent tool。",
		Model:       &outerToolModel{},
		ToolsConfig: adk.ToolsConfig{
			EmitInternalEvents: emit,
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{innerTool},
			},
		},
	})
	if err != nil {
		log.Fatalf("create outer agent failed: %v", err)
	}

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: outer}).Query(ctx, "lesson45 inner event 如何透出")

	fmt.Println(title + ":")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run %s failed: %v", title, event.Err)
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
		if strings.TrimSpace(message.Content) == "" && len(message.ToolCalls) > 0 {
			fmt.Printf("agent=%s role=%s tool_calls=%s\n", event.AgentName, event.Output.MessageOutput.Role, joinToolCallNames(message.ToolCalls))
			continue
		}
		fmt.Printf("agent=%s role=%s content=%s\n", event.AgentName, event.Output.MessageOutput.Role, message.Content)
	}
}

func (m *outerToolModel) Generate(_ context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	if m.step == 0 {
		m.step++
		options := model.GetCommonOptions(&model.Options{}, opts...)
		toolName := "inner_agent"
		if len(options.Tools) > 0 {
			toolName = options.Tools[0].Name
		}
		return schema.AssistantMessage("", []schema.ToolCall{
			{
				ID: "inner_call_1",
				Function: schema.FunctionCall{
					Name:      toolName,
					Arguments: `{"request":"请给我一个 inner answer"}`,
				},
			},
		}), nil
	}

	last := ""
	if len(input) > 0 && input[len(input)-1] != nil {
		last = input[len(input)-1].Content
	}
	return schema.AssistantMessage("outer final -> "+last, nil), nil
}

func (m *outerToolModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}

func (m *outerToolModel) WithTools(_ []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	return m, nil
}

func (a *innerAgent) Name(context.Context) string {
	return a.name
}

func (a *innerAgent) Description(context.Context) string {
	return "return one assistant message so outer agent can decide whether to emit it"
}

func (a *innerAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	last := ""
	if len(input.Messages) > 0 && input.Messages[len(input.Messages)-1] != nil {
		last = input.Messages[len(input.Messages)-1].Content
	}

	event := adk.EventFromMessage(
		schema.AssistantMessage("inner answer for "+last, nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name

	gen.Send(event)
	gen.Close()
	return iter
}

func joinToolCallNames(toolCalls []schema.ToolCall) string {
	names := make([]string, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		names = append(names, toolCall.Function.Name)
	}
	return strings.Join(names, ",")
}
