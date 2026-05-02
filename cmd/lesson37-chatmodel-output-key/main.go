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

type outputKeyModel struct{}

type sessionReaderAgent struct {
	name string
}

func main() {
	ctx := context.Background()

	chatAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "draft_agent",
		Description: "write a draft and store it into session",
		Instruction: "你负责先生成一个 draft。",
		Model:       &outputKeyModel{},
		OutputKey:   "draft_output",
	})
	if err != nil {
		log.Fatalf("create chat model agent failed: %v", err)
	}

	reader := &sessionReaderAgent{name: "session_reader"}
	workflow, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        "lesson37_workflow",
		Description: "read OutputKey from session in next agent",
		SubAgents:   []adk.Agent{chatAgent, reader},
	})
	if err != nil {
		log.Fatalf("create workflow failed: %v", err)
	}

	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: workflow}).Query(ctx, "请为 lesson37 写一句介绍")

	fmt.Println("agent events:")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run workflow failed: %v", event.Err)
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

func (m *outputKeyModel) Generate(_ context.Context, input []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	system := ""
	var userParts []string
	for _, message := range input {
		if message == nil {
			continue
		}
		switch message.Role {
		case schema.System:
			system = message.Content
		case schema.User:
			userParts = append(userParts, message.Content)
		}
	}

	return schema.AssistantMessage(
		fmt.Sprintf("draft generated: system=%s users=%s", system, strings.Join(userParts, " / ")),
		nil,
	), nil
}

func (m *outputKeyModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	message, err := m.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return schema.StreamReaderFromArray([]*schema.Message{message}), nil
}

func (a *sessionReaderAgent) Name(context.Context) string {
	return a.name
}

func (a *sessionReaderAgent) Description(context.Context) string {
	return "read draft output from ADK session"
}

func (a *sessionReaderAgent) Run(ctx context.Context, _ *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	value, _ := adk.GetSessionValue(ctx, "draft_output")
	event := adk.EventFromMessage(
		schema.AssistantMessage(fmt.Sprintf("reader got session draft_output=%v", value), nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name
	gen.Send(event)
	gen.Close()
	return iter
}
