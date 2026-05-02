package lesson50app

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type plannerAgent struct {
	name string
}

type writerAgent struct {
	name string
}

func newWorkflow(ctx context.Context, cfg Config) (adk.ResumableAgent, error) {
	return adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        cfg.AppName + "_workflow",
		Description: "planner + writer mini app",
		SubAgents: []adk.Agent{
			&plannerAgent{name: "planner_agent"},
			&writerAgent{name: "writer_agent"},
		},
	})
}

func (a *plannerAgent) Name(context.Context) string {
	return a.name
}

func (a *plannerAgent) Description(context.Context) string {
	return "plan the response and store it into session"
}

func (a *plannerAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	query := ""
	if len(input.Messages) > 0 && input.Messages[len(input.Messages)-1] != nil {
		query = input.Messages[len(input.Messages)-1].Content
	}
	plan := fmt.Sprintf("plan(%s): 1.analyze 2.write 3.review", query)
	adk.AddSessionValue(ctx, "plan", plan)

	event := adk.EventFromMessage(schema.AssistantMessage("planner stored "+plan, nil), nil, schema.Assistant, "")
	event.AgentName = a.name
	gen.Send(event)
	gen.Close()
	return iter
}

func (a *writerAgent) Name(context.Context) string {
	return a.name
}

func (a *writerAgent) Description(context.Context) string {
	return "write final answer from session plan"
}

func (a *writerAgent) Run(ctx context.Context, _ *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	plan, _ := adk.GetSessionValue(ctx, "plan")
	tenant, _ := adk.GetSessionValue(ctx, "tenant")
	content := fmt.Sprintf("writer final: tenant=%v using %v", tenant, plan)

	event := adk.EventFromMessage(schema.AssistantMessage(content, nil), nil, schema.Assistant, "")
	event.AgentName = a.name
	gen.Send(event)
	gen.Close()
	return iter
}
