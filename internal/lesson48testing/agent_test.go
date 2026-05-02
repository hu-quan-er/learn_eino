package lesson48testing

import (
	"context"
	"testing"
)

func TestRunOnce(t *testing.T) {
	outputs, err := RunOnce(context.Background(), "test_agent", "review", "lesson48 query")
	if err != nil {
		t.Fatalf("RunOnce returned error: %v", err)
	}
	if len(outputs) != 1 {
		t.Fatalf("expected 1 output, got %d", len(outputs))
	}
	if outputs[0] != "review:lesson48 query" {
		t.Fatalf("unexpected output: %s", outputs[0])
	}
}

func TestReviewAgentMetadata(t *testing.T) {
	agent := NewReviewAgent("meta_agent", "review")
	if agent.Name(context.Background()) != "meta_agent" {
		t.Fatalf("unexpected name: %s", agent.Name(context.Background()))
	}
	if agent.Description(context.Background()) == "" {
		t.Fatal("expected non-empty description")
	}
}
