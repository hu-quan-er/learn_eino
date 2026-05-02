package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type ReviewRoleAgent struct {
	name        string
	description string
	prefix      string
	delay       time.Duration
}

func main() {
	ctx := context.Background()

	workflow := buildDeliveryWorkflow()

	runner, err := workflow.Compile(ctx)
	if err != nil {
		log.Fatalf("compile workflow failed: %v", err)
	}

	workflowOutput, err := runner.Invoke(ctx, map[string]any{
		"topic":    "  Eino 第 25 课：可扩展脚手架  ",
		"audience": "进阶同学",
	})
	if err != nil {
		log.Fatalf("invoke workflow failed: %v", err)
	}

	draft, _ := workflowOutput["draft"].(string)

	reviewTeam, err := buildReviewTeam(ctx)
	if err != nil {
		log.Fatalf("build review team failed: %v", err)
	}

	reviewEvents, reviewSummary, err := runReviewTeam(ctx, reviewTeam, draft)
	if err != nil {
		log.Fatalf("run review team failed: %v", err)
	}

	fmt.Println("review events:")
	for _, event := range reviewEvents {
		fmt.Println("-", event)
	}

	fmt.Println("\nextensible scaffold output:")
	fmt.Printf("%#v\n", map[string]any{
		"topic":          workflowOutput["topic"],
		"audience":       workflowOutput["audience"],
		"draft":          draft,
		"review_summary": reviewSummary,
		"final_package":  fmt.Sprintf("%s | %s", draft, reviewSummary),
	})
}

func buildDeliveryWorkflow() *compose.Workflow[map[string]any, map[string]any] {
	workflow := compose.NewWorkflow[map[string]any, map[string]any]()

	workflow.
		AddLambdaNode("prepare_request", compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
			_ = ctx

			topic, _ := input["topic"].(string)
			audience, _ := input["audience"].(string)

			return map[string]any{
				"topic":    strings.TrimSpace(topic),
				"audience": audience,
			}, nil
		})).
		AddInput(compose.START)

	workflow.
		AddGraphNode("draft_pipeline", buildDraftGraph(), compose.WithGraphCompileOptions(
			compose.WithNodeTriggerMode(compose.AllPredecessor),
			compose.WithGraphName("lesson25_draft_pipeline"),
		)).
		AddInput("prepare_request", compose.FromField("topic"))

	workflow.
		AddLambdaNode("collect_draft", compose.InvokableLambda(func(ctx context.Context, input string) (map[string]any, error) {
			_ = ctx
			return map[string]any{
				"draft": input,
			}, nil
		})).
		AddInput("draft_pipeline")

	workflow.End().
		AddInput("prepare_request",
			compose.MapFields("topic", "topic"),
			compose.MapFields("audience", "audience"),
		)

	workflow.End().AddInput("collect_draft", compose.MapFields("draft", "draft"))

	return workflow
}

func buildDraftGraph() *compose.Graph[string, string] {
	graph := compose.NewGraph[string, string]()

	_ = graph.AddLambdaNode("outline", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return "提纲：" + input, nil
	}), compose.WithOutputKey("outline"))

	_ = graph.AddLambdaNode("examples", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
		_ = ctx
		return "案例：" + input + " 的两个最小 demo", nil
	}), compose.WithOutputKey("examples"))

	_ = graph.AddLambdaNode("merge", compose.InvokableLambda(func(ctx context.Context, input map[string]any) (string, error) {
		_ = ctx

		outline, _ := input["outline"].(string)
		examples, _ := input["examples"].(string)

		return fmt.Sprintf("初稿：%s | %s", outline, examples), nil
	}))

	_ = graph.AddEdge(compose.START, "outline")
	_ = graph.AddEdge(compose.START, "examples")
	_ = graph.AddEdge("outline", "merge")
	_ = graph.AddEdge("examples", "merge")
	_ = graph.AddEdge("merge", compose.END)

	return graph
}

func buildReviewTeam(ctx context.Context) (adk.ResumableAgent, error) {
	factChecker := &ReviewRoleAgent{
		name:        "fact_checker",
		description: "check whether the draft is coherent",
		prefix:      "核对完成",
		delay:       120 * time.Millisecond,
	}
	editor := &ReviewRoleAgent{
		name:        "editor",
		description: "turn the previous result into a publishable summary",
		prefix:      "编辑结论",
		delay:       80 * time.Millisecond,
	}

	return adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        "review_team",
		Description: "review the draft in sequence",
		SubAgents:   []adk.Agent{factChecker, editor},
	})
}

func runReviewTeam(ctx context.Context, agent adk.Agent, draft string) ([]string, string, error) {
	iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent}).Query(ctx, draft)

	events := make([]string, 0)
	lastMessage := ""

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			return nil, "", event.Err
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}

		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			return nil, "", err
		}
		if message == nil {
			continue
		}

		line := fmt.Sprintf("%s: %s", event.AgentName, message.Content)
		events = append(events, line)
		lastMessage = message.Content
	}

	return events, lastMessage, nil
}

func (a *ReviewRoleAgent) Name(ctx context.Context) string {
	_ = ctx
	return a.name
}

func (a *ReviewRoleAgent) Description(ctx context.Context) string {
	_ = ctx
	return a.description
}

func (a *ReviewRoleAgent) Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	_ = ctx
	_ = options

	lastMessage := ""
	if input != nil && len(input.Messages) > 0 && input.Messages[len(input.Messages)-1] != nil {
		lastMessage = input.Messages[len(input.Messages)-1].Content
	}

	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	go func() {
		defer gen.Close()
		time.Sleep(a.delay)
		gen.Send(&adk.AgentEvent{
			AgentName: a.name,
			Output: &adk.AgentOutput{
				MessageOutput: &adk.MessageVariant{
					IsStreaming: false,
					Message:     schema.AssistantMessage(fmt.Sprintf("%s -> %s", a.prefix, lastMessage), nil),
					Role:        schema.Assistant,
				},
			},
		})
	}()

	return iter
}
