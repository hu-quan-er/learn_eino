package lesson49service

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type ApprovalState struct {
	Request string
}

type PendingApproval struct {
	CheckPointID string
	InterruptID  string
	Preview      string
	Reason       string
}

type Service struct {
	runner        *adk.Runner
	checkPointID  string
	lastInterrupt string
}

type approvalAgent struct {
	name string
}

type inMemoryStore struct {
	data map[string][]byte
}

func init() {
	schema.Register[ApprovalState]()
}

func New(ctx context.Context) *Service {
	store := &inMemoryStore{data: make(map[string][]byte)}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           &approvalAgent{name: "lesson49_service_agent"},
		CheckPointStore: store,
	})

	return &Service{
		runner:       runner,
		checkPointID: "lesson49-service-checkpoint",
	}
}

func (s *Service) StartPublish(ctx context.Context, request string) (*PendingApproval, error) {
	iter := s.runner.Query(ctx, request, adk.WithCheckPointID(s.checkPointID))

	result := &PendingApproval{CheckPointID: s.checkPointID}
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			return nil, event.Err
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			message, err := event.Output.MessageOutput.GetMessage()
			if err != nil {
				return nil, err
			}
			if result.Preview == "" {
				result.Preview = message.Content
			}
		}
		if event.Action != nil && event.Action.Interrupted != nil {
			for _, interruptCtx := range event.Action.Interrupted.InterruptContexts {
				if interruptCtx.IsRootCause {
					result.InterruptID = interruptCtx.ID
					result.Reason = fmt.Sprint(interruptCtx.Info)
				}
			}
		}
	}

	s.lastInterrupt = result.InterruptID
	return result, nil
}

func (s *Service) ResumePublish(ctx context.Context, decision string) (string, error) {
	iter, err := s.runner.ResumeWithParams(ctx, s.checkPointID, &adk.ResumeParams{
		Targets: map[string]any{
			s.lastInterrupt: decision,
		},
	})
	if err != nil {
		return "", err
	}

	final := ""
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			return "", event.Err
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			return "", err
		}
		final = message.Content
	}

	return final, nil
}

func (a *approvalAgent) Name(context.Context) string {
	return a.name
}

func (a *approvalAgent) Description(context.Context) string {
	return "service-oriented resumable agent"
}

func (a *approvalAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	request := ""
	if len(input.Messages) > 0 && input.Messages[len(input.Messages)-1] != nil {
		request = input.Messages[len(input.Messages)-1].Content
	}

	message := adk.EventFromMessage(schema.AssistantMessage("service preview: "+request, nil), nil, schema.Assistant, "")
	message.AgentName = a.name
	gen.Send(message)

	interrupt := adk.StatefulInterrupt(ctx, "service needs approval", ApprovalState{Request: request})
	interrupt.AgentName = a.name
	gen.Send(interrupt)

	gen.Close()
	return iter
}

func (a *approvalAgent) Resume(ctx context.Context, info *adk.ResumeInfo, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	state, ok := info.InterruptState.(ApprovalState)
	if !ok {
		if ptr, okPtr := info.InterruptState.(*ApprovalState); okPtr && ptr != nil {
			state = *ptr
			ok = true
		}
	}
	if !ok {
		gen.Send(&adk.AgentEvent{Err: fmt.Errorf("unexpected interrupt state: %T", info.InterruptState)})
		gen.Close()
		return iter
	}

	decision, _ := info.ResumeData.(string)
	event := adk.EventFromMessage(
		schema.AssistantMessage(fmt.Sprintf("service approved(%s): %s", decision, state.Request), nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name
	gen.Send(event)
	gen.Close()
	return iter
}

func (s *inMemoryStore) Get(_ context.Context, checkPointID string) ([]byte, bool, error) {
	value, ok := s.data[checkPointID]
	return value, ok, nil
}

func (s *inMemoryStore) Set(_ context.Context, checkPointID string, checkPoint []byte) error {
	s.data[checkPointID] = checkPoint
	return nil
}
