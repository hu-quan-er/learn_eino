package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type SearchInput struct {
	Query string `json:"query" jsonschema:"required,description=search query"`
}

func main() {
	ctx := context.Background()

	searchTool, err := toolutils.InferStreamTool("search_docs", "stream mock search result chunks", streamSearchDocs)
	if err != nil {
		log.Fatalf("create streamable tool failed: %v", err)
	}

	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{searchTool},
	})
	if err != nil {
		log.Fatalf("create tools node failed: %v", err)
	}

	assistantMessage := schema.AssistantMessage("", []schema.ToolCall{
		{
			ID:   "call_search_1",
			Type: "function",
			Function: schema.FunctionCall{
				Name:      "search_docs",
				Arguments: `{"query":"Eino StreamableTool"}`,
			},
		},
	})

	stream, err := toolsNode.Stream(ctx, assistantMessage)
	if err != nil {
		log.Fatalf("stream tools node failed: %v", err)
	}
	defer stream.Close()

	var chunks [][]*schema.Message

	fmt.Println("stream chunks:")
	for {
		messages, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("recv tool chunk failed: %v", err)
		}

		chunks = append(chunks, messages)
		for _, message := range messages {
			if message == nil {
				continue
			}
			fmt.Printf("- tool_name=%s content=%s\n", message.ToolName, message.Content)
		}
	}

	merged, err := schema.ConcatMessageArray(chunks)
	if err != nil {
		log.Fatalf("concat tool chunks failed: %v", err)
	}

	fmt.Println("\nmerged tool message:")
	for _, message := range merged {
		if message == nil {
			continue
		}
		fmt.Printf("role=%s tool_name=%s content=%s\n", message.Role, message.ToolName, message.Content)
	}
}

func streamSearchDocs(ctx context.Context, input *SearchInput) (*schema.StreamReader[string], error) {
	_ = ctx

	reader, writer := schema.Pipe[string](0)

	go func() {
		defer writer.Close()

		chunks := []string{
			fmt.Sprintf("收到查询：%s", input.Query),
			"检索到 3 篇和 StreamableTool 相关的资料",
			"总结：StreamableTool 适合把长结果按 chunk 持续返回给上游",
		}

		for _, chunk := range chunks {
			writer.Send(chunk, nil)
			time.Sleep(150 * time.Millisecond)
		}
	}()

	return reader, nil
}
