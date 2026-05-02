package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

type plannerAgent struct {
	name string
}

type executorAgent struct {
	name string
}

func main() {
	ctx := context.Background()

	workflow, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        "lesson34_session_workflow",
		Description: "share session values across custom agents",
		SubAgents: []adk.Agent{
			&plannerAgent{name: "planner_agent"},
			&executorAgent{name: "executor_agent"},
		},
	})
	if err != nil {
		log.Fatalf("create sequential agent failed: %v", err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: workflow})
	iter := runner.Query(ctx, "准备一节 ADK 教程", adk.WithSessionValues(map[string]any{
		"tenant": "tutorial-team",
	}))

	fmt.Println("agent events:")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run agent failed: %v", event.Err)
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

func (a *plannerAgent) Name(context.Context) string {
	return a.name
}

func (a *plannerAgent) Description(context.Context) string {
	return "write plan into ADK session values"
}

func (a *plannerAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	plan := fmt.Sprintf("%s -> 收集资料 -> 写示例 -> 校验输出", lastUserContent(input.Messages))
	adk.AddSessionValue(ctx, "plan", plan)
	adk.AddSessionValue(ctx, "owner", a.name)

	event := adk.EventFromMessage(
		schema.AssistantMessage("planner 写入 session: "+plan, nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name

	gen.Send(event)
	gen.Close()
	return iter
}

func (a *executorAgent) Name(context.Context) string {
	return a.name
}

func (a *executorAgent) Description(context.Context) string {
	return "read session values prepared by previous agents"
}

func (a *executorAgent) Run(ctx context.Context, _ *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	plan, _ := adk.GetSessionValue(ctx, "plan")
	owner, _ := adk.GetSessionValue(ctx, "owner")
	tenant, _ := adk.GetSessionValue(ctx, "tenant")
	keys := sessionKeys(adk.GetSessionValues(ctx))

	content := fmt.Sprintf(
		"executor 读取 session: tenant=%v owner=%v plan=%v keys=%s",
		tenant,
		owner,
		plan,
		strings.Join(keys, ","),
	)

	event := adk.EventFromMessage(
		schema.AssistantMessage(content, nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name

	gen.Send(event)
	gen.Close()
	return iter
}

func lastUserContent(messages []adk.Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i] != nil && messages[i].Role == schema.User {
			return messages[i].Content
		}
	}

	return ""
}

func sessionKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
