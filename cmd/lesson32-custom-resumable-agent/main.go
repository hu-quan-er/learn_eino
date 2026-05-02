package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type ApprovalState struct {
	Request string
}

type resumableAgent struct {
	name string
}

type inMemoryStore struct {
	data map[string][]byte
}

func init() {
	schema.Register[ApprovalState]()
}

func main() {
	ctx := context.Background()
	checkPointID := "lesson32-custom-resumable-agent"
	agent := &resumableAgent{name: "lesson32_resumable_agent"}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		CheckPointStore: newInMemoryStore(),
	})

	fmt.Println("first run:")
	iter := runner.Query(ctx, "发布第 32 课", adk.WithCheckPointID(checkPointID))
	rootCauseID := consumeEvents(iter)
	if rootCauseID == "" {
		log.Fatal("expected interrupt root cause id, got empty")
	}

	fmt.Println("\nresume from checkpoint:")
	iter, err := runner.ResumeWithParams(ctx, checkPointID, &adk.ResumeParams{
		Targets: map[string]any{
			rootCauseID: "approved by reviewer",
		},
	})
	if err != nil {
		log.Fatalf("resume agent failed: %v", err)
	}

	consumeEvents(iter)
}

func (a *resumableAgent) Name(context.Context) string {
	return a.name
}

func (a *resumableAgent) Description(context.Context) string {
	return "interrupt once and continue after external approval"
}

func (a *resumableAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	request := lastUserContent(input.Messages)

	messageEvent := adk.EventFromMessage(
		schema.AssistantMessage("收到请求："+request, nil),
		nil,
		schema.Assistant,
		"",
	)
	messageEvent.AgentName = a.name
	gen.Send(messageEvent)

	interruptEvent := adk.StatefulInterrupt(ctx, "需要人工审批", ApprovalState{Request: request})
	interruptEvent.AgentName = a.name
	gen.Send(interruptEvent)

	gen.Close()
	return iter
}

func (a *resumableAgent) Resume(ctx context.Context, info *adk.ResumeInfo, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	if info == nil || !info.WasInterrupted {
		gen.Send(&adk.AgentEvent{Err: fmt.Errorf("resume called without interrupt state")})
		gen.Close()
		return iter
	}

	state, ok := approvalStateFromAny(info.InterruptState)
	if !ok {
		gen.Send(&adk.AgentEvent{Err: fmt.Errorf("unexpected interrupt state type: %T", info.InterruptState)})
		gen.Close()
		return iter
	}

	if !info.IsResumeTarget {
		interruptEvent := adk.StatefulInterrupt(ctx, "resume target not reached", state)
		interruptEvent.AgentName = a.name
		gen.Send(interruptEvent)
		gen.Close()
		return iter
	}

	decision, _ := info.ResumeData.(string)
	if decision == "" {
		decision = "approved(default)"
	}

	messageEvent := adk.EventFromMessage(
		schema.AssistantMessage(fmt.Sprintf("审批通过(%s)：%s", decision, state.Request), nil),
		nil,
		schema.Assistant,
		"",
	)
	messageEvent.AgentName = a.name
	gen.Send(messageEvent)

	gen.Close()
	return iter
}

func approvalStateFromAny(v any) (ApprovalState, bool) {
	switch state := v.(type) {
	case ApprovalState:
		return state, true
	case *ApprovalState:
		if state != nil {
			return *state, true
		}
	}

	return ApprovalState{}, false
}

func consumeEvents(iter *adk.AsyncIterator[*adk.AgentEvent]) string {
	var rootCauseID string

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("agent execution failed: %v", event.Err)
		}

		if event.Output != nil && event.Output.MessageOutput != nil {
			message, err := event.Output.MessageOutput.GetMessage()
			if err != nil {
				log.Fatalf("read message failed: %v", err)
			}
			fmt.Printf("message: %s\n", message.Content)
		}

		if event.Action == nil || event.Action.Interrupted == nil {
			continue
		}

		for _, interruptCtx := range event.Action.Interrupted.InterruptContexts {
			fmt.Printf(
				"interrupt: id=%s address=%s info=%v root=%v\n",
				interruptCtx.ID,
				interruptCtx.Address.String(),
				interruptCtx.Info,
				interruptCtx.IsRootCause,
			)
			if interruptCtx.IsRootCause {
				rootCauseID = interruptCtx.ID
			}
		}
	}

	return rootCauseID
}

func lastUserContent(messages []adk.Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i] != nil && messages[i].Role == schema.User {
			return messages[i].Content
		}
	}

	if len(messages) == 0 || messages[len(messages)-1] == nil {
		return ""
	}

	return messages[len(messages)-1].Content
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
