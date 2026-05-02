package main

import (
	"context"
	"fmt"
	"log"

	"eino-tutorial/internal/lesson49service"
)

func main() {
	ctx := context.Background()
	service := lesson49service.New(ctx)

	pending, err := service.StartPublish(ctx, "发布第 49 课")
	if err != nil {
		log.Fatalf("start publish failed: %v", err)
	}

	fmt.Println("pending approval:")
	fmt.Printf("checkpoint=%s\n", pending.CheckPointID)
	fmt.Printf("interrupt_id=%s\n", pending.InterruptID)
	fmt.Printf("preview=%s\n", pending.Preview)
	fmt.Printf("reason=%s\n", pending.Reason)

	final, err := service.ResumePublish(ctx, "approved by service")
	if err != nil {
		log.Fatalf("resume publish failed: %v", err)
	}

	fmt.Println()
	fmt.Println("final result:")
	fmt.Println(final)
}
