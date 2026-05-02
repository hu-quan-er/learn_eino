package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudwego/eino/compose"
)

type PublishRequest struct {
	Topic    string
	Audience string
}

type AudiencePayload struct {
	Value string
}

type MergeInput struct {
	Topic    string
	Audience string
	Original *PublishRequest
}

type PublishResult struct {
	Result string
}

func main() {
	ctx := context.Background()

	workflow := compose.NewWorkflow[*PublishRequest, *PublishResult]()

	workflow.
		AddLambdaNode("normalize_topic", compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
			_ = ctx
			return strings.TrimSpace(input), nil
		})).
		AddInput(compose.START, compose.FromField("Topic"))

	workflow.
		AddLambdaNode("audience_payload", compose.InvokableLambda(func(ctx context.Context, input *AudiencePayload) (*AudiencePayload, error) {
			_ = ctx
			return &AudiencePayload{Value: strings.ToUpper(strings.TrimSpace(input.Value))}, nil
		})).
		AddInput(compose.START, compose.MapFields("Audience", "Value"))

	workflow.
		AddLambdaNode("merge", compose.InvokableLambda(func(ctx context.Context, input *MergeInput) (string, error) {
			_ = ctx
			return fmt.Sprintf("发布：%s [%s] (original=%s/%s)", input.Topic, input.Audience, input.Original.Topic, input.Original.Audience), nil
		})).
		AddInput("normalize_topic", compose.ToField("Topic")).
		AddInput("audience_payload", compose.MapFields("Value", "Audience")).
		AddInput(compose.START, compose.ToField("Original"))

	workflow.End().AddInput("merge", compose.ToField("Result"))

	runner, err := workflow.Compile(ctx)
	if err != nil {
		log.Fatalf("compile workflow failed: %v", err)
	}

	output, err := runner.Invoke(ctx, &PublishRequest{
		Topic:    "  第 28 课：Workflow 字段映射  ",
		Audience: "初学者",
	})
	if err != nil {
		log.Fatalf("invoke workflow failed: %v", err)
	}

	fmt.Printf("workflow output: %#v\n", output)
}
