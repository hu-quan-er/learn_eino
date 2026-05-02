package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type parentTransferModel struct{}

type childEchoModel struct{}

func main() {
	ctx := context.Background()

	runCase(ctx, "default transfer history")
	fmt.Println()
	runCase(ctx, "skip transfer messages", adk.WithSkipTransferMessages())
}

func runCase(ctx context.Context, title string, opts ...adk.AgentRunOption) {
	parentAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ParentAgent",
		Description: "transfer to child agent",
		Instruction: "你负责把任务转交给 ChildAgent。",
		Model:       &parentTransferModel{},
	})
	if err != nil {
		log.Fatalf("create parent agent failed: %v", err)
	}

	childAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ChildAgent",
		Description: "echo all messages it receives",
		Instruction: "你负责打印自己收到的全部输入。",
		Model:       &childEchoModel{},
	})
	if err != nil {
		log.Fatalf("create child agent failed: %v", err)
	}

	root, err := adk.SetSubAgents(ctx, parentAgent, []adk.Agent{childAgent})
	if err != nil {
		log.Fatalf("set sub agents failed: %v", err)
	}

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: root}).Query(ctx, "lesson39 要强调什么", opts...)

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
		fmt.Printf("agent=%s role=%s content=%s\n", event.AgentName, event.Output.MessageOutput.Role, message.Content)
	}
}

func (m *parentTransferModel) Generate(_ context.Context, _ []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	return schema.AssistantMessage("我会把这个请求交给 ChildAgent。", []schema.ToolCall{
		{
			ID: "transfer_call_1",
			Function: schema.FunctionCall{
				Name:      adk.TransferToAgentToolName,
				Arguments: `{"agent_name":"ChildAgent"}`,
			},
		},
	}), nil
}

func (m *parentTransferModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}

func (m *parentTransferModel) WithTools(_ []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	return m, nil
}

func (m *childEchoModel) Generate(_ context.Context, input []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	parts := make([]string, 0, len(input))
	for _, message := range input {
		if message == nil {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", message.Role, message.Content))
	}
	return schema.AssistantMessage(strings.Join(parts, " | "), nil), nil
}

func (m *childEchoModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}

func (m *childEchoModel) WithTools(_ []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	return m, nil
}
