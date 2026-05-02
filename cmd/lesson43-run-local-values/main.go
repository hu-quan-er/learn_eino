package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type runLocalModel struct{}

type runLocalHandler struct {
	*adk.BaseChatModelAgentMiddleware
}

func main() {
	ctx := context.Background()

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "lesson43_agent",
		Description: "show run local values in ChatModelAgent handlers",
		Instruction: "你是一个 run-local 演示 agent。",
		Model:       &runLocalModel{},
		Handlers: []adk.ChatModelAgentMiddleware{
			&runLocalHandler{BaseChatModelAgentMiddleware: &adk.BaseChatModelAgentMiddleware{}},
		},
	})
	if err != nil {
		log.Fatalf("create chat model agent failed: %v", err)
	}

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent}).Query(ctx, "lesson43 怎么理解")

	fmt.Println("agent events:")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run lesson43 failed: %v", event.Err)
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

func (h *runLocalHandler) BeforeModelRewriteState(ctx context.Context, state *adk.ChatModelAgentState, mc *adk.ModelContext) (context.Context, *adk.ChatModelAgentState, error) {
	_ = state
	_ = mc
	return ctx, state, adk.SetRunLocalValue(ctx, "trace_id", "trace-lesson43")
}

func (h *runLocalHandler) AfterModelRewriteState(ctx context.Context, state *adk.ChatModelAgentState, mc *adk.ModelContext) (context.Context, *adk.ChatModelAgentState, error) {
	_ = mc

	value, found, err := adk.GetRunLocalValue(ctx, "trace_id")
	if err != nil {
		return ctx, state, err
	}
	if err := adk.DeleteRunLocalValue(ctx, "trace_id"); err != nil {
		return ctx, state, err
	}
	_, stillFound, err := adk.GetRunLocalValue(ctx, "trace_id")
	if err != nil {
		return ctx, state, err
	}

	if len(state.Messages) > 0 {
		last := state.Messages[len(state.Messages)-1]
		last.Content = fmt.Sprintf("%s | trace_found=%v trace=%v deleted=%v", last.Content, found, value, !stillFound)
	}

	return ctx, state, nil
}

func (m *runLocalModel) Generate(_ context.Context, input []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	last := ""
	if len(input) > 0 && input[len(input)-1] != nil {
		last = input[len(input)-1].Content
	}
	return schema.AssistantMessage("base answer for "+last, nil), nil
}

func (m *runLocalModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}
