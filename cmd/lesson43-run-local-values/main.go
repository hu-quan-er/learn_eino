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

type runLocalReaderAgent struct {
	name string
}

func main() {
	ctx := context.Background()

	chatAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "lesson43_chat_agent",
		Description: "show run local values in handlers",
		Instruction: "你是一个 run-local 演示 agent。",
		Model:       &runLocalModel{},
		Handlers: []adk.ChatModelAgentMiddleware{
			&runLocalHandler{BaseChatModelAgentMiddleware: &adk.BaseChatModelAgentMiddleware{}},
		},
	})
	if err != nil {
		log.Fatalf("create chat model agent failed: %v", err)
	}

	reader := &runLocalReaderAgent{name: "lesson43_reader"}
	workflow, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        "lesson43_workflow",
		Description: "read run-local summary written by handler",
		SubAgents:   []adk.Agent{chatAgent, reader},
	})
	if err != nil {
		log.Fatalf("create lesson43 workflow failed: %v", err)
	}

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: workflow}).Query(ctx, "lesson43 怎么理解")

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
	_ = state
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

	adk.AddSessionValue(ctx, "run_local_summary", fmt.Sprintf("trace_found=%v trace=%v deleted=%v", found, value, !stillFound))
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

func (a *runLocalReaderAgent) Name(context.Context) string {
	return a.name
}

func (a *runLocalReaderAgent) Description(context.Context) string {
	return "read run-local summary from session"
}

func (a *runLocalReaderAgent) Run(ctx context.Context, _ *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	value, _ := adk.GetSessionValue(ctx, "run_local_summary")
	event := adk.EventFromMessage(
		schema.AssistantMessage(fmt.Sprintf("reader saw run_local_summary=%v", value), nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name
	gen.Send(event)
	gen.Close()
	return iter
}
