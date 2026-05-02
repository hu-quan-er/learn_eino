package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type baseWrapModel struct{}

type wrapHandler struct {
	*adk.BaseChatModelAgentMiddleware
}

type modelPrefixWrapper struct {
	inner model.BaseChatModel
}

func main() {
	ctx := context.Background()

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "lesson42_agent",
		Description: "emit custom event and wrap the model output",
		Instruction: "你是一个会发自定义事件的 agent。",
		Model:       &baseWrapModel{},
		Handlers: []adk.ChatModelAgentMiddleware{
			&wrapHandler{BaseChatModelAgentMiddleware: &adk.BaseChatModelAgentMiddleware{}},
		},
	})
	if err != nil {
		log.Fatalf("create chat model agent failed: %v", err)
	}

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent}).Query(ctx, "lesson42 的重点是什么")

	fmt.Println("agent events:")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run lesson42 failed: %v", event.Err)
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

func (h *wrapHandler) BeforeModelRewriteState(ctx context.Context, state *adk.ChatModelAgentState, mc *adk.ModelContext) (context.Context, *adk.ChatModelAgentState, error) {
	_ = mc

	event := adk.EventFromMessage(
		schema.AssistantMessage(fmt.Sprintf("custom event before model, messages=%d", len(state.Messages)), nil),
		nil,
		schema.Assistant,
		"",
	)
	if err := adk.SendEvent(ctx, event); err != nil {
		return ctx, state, err
	}

	return ctx, state, nil
}

func (h *wrapHandler) WrapModel(ctx context.Context, m model.BaseChatModel, mc *adk.ModelContext) (model.BaseChatModel, error) {
	_ = ctx
	_ = mc
	return &modelPrefixWrapper{inner: m}, nil
}

func (m *baseWrapModel) Generate(_ context.Context, input []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	last := ""
	if len(input) > 0 && input[len(input)-1] != nil {
		last = input[len(input)-1].Content
	}
	return schema.AssistantMessage("base model answer for "+last, nil), nil
}

func (m *baseWrapModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}

func (m *modelPrefixWrapper) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	message, err := m.inner.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	message.Content = "wrapped -> " + message.Content
	return message, nil
}

func (m *modelPrefixWrapper) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}
