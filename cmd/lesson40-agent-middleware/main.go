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

type middlewareEchoModel struct{}

func main() {
	ctx := context.Background()

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "lesson40_agent",
		Description: "show AgentMiddleware hooks",
		Instruction: "你是一个会回显输入的 agent。",
		Model:       &middlewareEchoModel{},
		Middlewares: []adk.AgentMiddleware{
			{
				AdditionalInstruction: "额外规则：回答里必须显式出现 middleware。",
				BeforeChatModel: func(ctx context.Context, state *adk.ChatModelAgentState) error {
					_ = ctx
					state.Messages = append(state.Messages, schema.UserMessage("before:middleware"))
					return nil
				},
				AfterChatModel: func(ctx context.Context, state *adk.ChatModelAgentState) error {
					_ = ctx
					if len(state.Messages) == 0 {
						return nil
					}
					last := state.Messages[len(state.Messages)-1]
					last.Content += " | after:middleware"
					return nil
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("create chat model agent failed: %v", err)
	}

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent}).Query(ctx, "lesson40 原始问题")

	fmt.Println("agent events:")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run lesson40 failed: %v", event.Err)
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}

		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Fatalf("read message failed: %v", err)
		}
		fmt.Printf("agent=%s content=%s\n", event.AgentName, message.Content)
	}
}

func (m *middlewareEchoModel) Generate(_ context.Context, input []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	system := ""
	userInputs := make([]string, 0)
	for _, message := range input {
		if message == nil {
			continue
		}
		if message.Role == schema.System {
			system = message.Content
		}
		if message.Role == schema.User {
			userInputs = append(userInputs, message.Content)
		}
	}

	return schema.AssistantMessage(
		fmt.Sprintf("system=%s users=%s", system, strings.Join(userInputs, " / ")),
		nil,
	), nil
}

func (m *middlewareEchoModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}
