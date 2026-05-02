package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	toolcomponent "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type capturingAgent struct {
	name          string
	description   string
	capturedInput []adk.Message
}

type historyParentAgent struct {
	name      string
	childTool toolcomponent.InvokableTool
}

func main() {
	ctx := context.Background()

	runDefaultSchemaDemo(ctx)
	fmt.Println()
	runCustomSchemaDemo(ctx)
	fmt.Println()
	runFullHistoryDemo(ctx)
}

func runDefaultSchemaDemo(ctx context.Context) {
	agent := &capturingAgent{
		name:        "default_schema_agent",
		description: "capture default agent tool input",
	}
	agentTool := adk.NewAgentTool(ctx, agent).(toolcomponent.InvokableTool)

	result, err := agentTool.InvokableRun(ctx, `{"request":"请解释默认 request 字段"}`)
	if err != nil {
		log.Fatalf("run default schema tool failed: %v", err)
	}

	fmt.Println("default schema:")
	fmt.Println("tool result:", result)
	printCapturedInput(agent.capturedInput)
}

func runCustomSchemaDemo(ctx context.Context) {
	agent := &capturingAgent{
		name:        "custom_schema_agent",
		description: "capture custom schema agent tool input",
	}
	customSchema := schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
		"topic": {
			Desc:     "topic to explain",
			Required: true,
			Type:     schema.String,
		},
		"level": {
			Desc:     "difficulty level",
			Required: true,
			Type:     schema.String,
		},
	})
	agentTool := adk.NewAgentTool(ctx, agent, adk.WithAgentInputSchema(customSchema)).(toolcomponent.InvokableTool)

	result, err := agentTool.InvokableRun(ctx, `{"topic":"agent tool","level":"advanced"}`)
	if err != nil {
		log.Fatalf("run custom schema tool failed: %v", err)
	}

	fmt.Println("custom schema:")
	fmt.Println("tool result:", result)
	printCapturedInput(agent.capturedInput)
}

func runFullHistoryDemo(ctx context.Context) {
	childAgent := &capturingAgent{
		name:        "history_child_agent",
		description: "capture full chat history from parent context",
	}
	childTool := adk.NewAgentTool(ctx, childAgent, adk.WithFullChatHistoryAsInput()).(toolcomponent.InvokableTool)
	parentAgent := &historyParentAgent{
		name:      "history_parent_agent",
		childTool: childTool,
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: parentAgent})
	iter := runner.Query(ctx, "什么时候该用 WithFullChatHistoryAsInput")

	fmt.Println("full chat history:")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("run history parent agent failed: %v", event.Err)
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}

		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Fatalf("read history parent message failed: %v", err)
		}

		fmt.Println("tool result:", message.Content)
	}

	printCapturedInput(childAgent.capturedInput)
}

func (a *capturingAgent) Name(context.Context) string {
	return a.name
}

func (a *capturingAgent) Description(context.Context) string {
	return a.description
}

func (a *capturingAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	a.capturedInput = append([]adk.Message(nil), input.Messages...)
	summary := fmt.Sprintf("child received %d messages, last=%q", len(input.Messages), lastMessageContent(input.Messages))

	event := adk.EventFromMessage(
		schema.AssistantMessage(summary, nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name

	gen.Send(event)
	gen.Close()
	return iter
}

func (a *historyParentAgent) Name(context.Context) string {
	return a.name
}

func (a *historyParentAgent) Description(context.Context) string {
	return "invoke child agent tool with full chat history"
}

func (a *historyParentAgent) Run(ctx context.Context, input *adk.AgentInput, _ ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	history := []adk.Message{
		schema.UserMessage("父 agent 收到问题：" + lastMessageContent(input.Messages)),
		schema.AssistantMessage("父 agent 先解释一半，然后准备把任务转给子 agent。", nil),
		schema.AssistantMessage("tool call placeholder", nil),
	}

	graph := compose.NewGraph[string, string](
		compose.WithGenLocalState(func(context.Context) *adk.State {
			return &adk.State{Messages: history}
		}),
	)

	err := graph.AddLambdaNode("invoke_tool", compose.InvokableLambda(func(ctx context.Context, _ string) (string, error) {
		return a.childTool.InvokableRun(ctx, `{"request":"this request will be ignored"}`)
	}))
	if err != nil {
		gen.Send(&adk.AgentEvent{Err: err})
		gen.Close()
		return iter
	}

	if err = graph.AddEdge(compose.START, "invoke_tool"); err != nil {
		gen.Send(&adk.AgentEvent{Err: err})
		gen.Close()
		return iter
	}
	if err = graph.AddEdge("invoke_tool", compose.END); err != nil {
		gen.Send(&adk.AgentEvent{Err: err})
		gen.Close()
		return iter
	}

	graphRunner, err := graph.Compile(ctx, compose.WithGraphName("lesson35_history_graph"))
	if err != nil {
		gen.Send(&adk.AgentEvent{Err: err})
		gen.Close()
		return iter
	}

	output, err := graphRunner.Invoke(ctx, "")
	if err != nil {
		gen.Send(&adk.AgentEvent{Err: err})
		gen.Close()
		return iter
	}

	event := adk.EventFromMessage(
		schema.AssistantMessage(output, nil),
		nil,
		schema.Assistant,
		"",
	)
	event.AgentName = a.name

	gen.Send(event)
	gen.Close()
	return iter
}

func printCapturedInput(messages []adk.Message) {
	for i, message := range messages {
		if message == nil {
			fmt.Printf("input[%d]: <nil>\n", i)
			continue
		}
		fmt.Printf("input[%d]: role=%s content=%s\n", i, message.Role, message.Content)
	}
}

func lastMessageContent(messages []adk.Message) string {
	if len(messages) == 0 || messages[len(messages)-1] == nil {
		return ""
	}

	return messages[len(messages)-1].Content
}
