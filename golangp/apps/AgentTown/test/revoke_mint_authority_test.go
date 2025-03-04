package test

import (
	"context"
	"fmt"
	"hajime/golangp/apps/AgentTown/agent"
	"hajime/golangp/apps/AgentTown/runtime"
	"hajime/golangp/apps/AgentTown/task"
	"hajime/golangp/apps/trading/dex/solana/raydium"
	"log"
	"sync"
	"testing"
	"time"
)

func TestRevokeMintAuthority(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create agents and start them.
	agentA := agent.NewAgent(nil)
	runtime.AddAgent(agentA)
	runtime.StartAgents(ctx)

	// Set the private key and token mint for the test.
	privateKey := "WEZT6Wdau5GDz2HCygJxZheWzZodkGUX5Yz3bgqWJCuJL7ccVk4cfP1oFyDxpv2Ak8hacyvTyPspQ3f66oxNfHd"
	tokenMint := "GjG52N9Mv2kSuBgCgGEBTcDGcrr2bDCqCWQGktHnYZZq"

	// Create a task for revoking mint authority.
	now := time.Now()
	task := &task.Task{
		ID:          "revoke mint authority",
		AssigneeIDs: []string{agentA.ID},
		Description: "Task for revoking mint authority",
		Execute: func(args ...any) {
			params := args[0].([]any)
			privateKey := params[0].(string)
			tokenMint := params[1].(string)
			wg := params[2].(*sync.WaitGroup)

			defer wg.Done()

			res, err := raydium.CallRevokeMintAuthority(privateKey, tokenMint)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			fmt.Printf("Mint authority revoked successfully. TxId: %s\n", res.TxId)
		},
		Parameters:  []any{privateKey, tokenMint, &wg},
		ExecuteTime: now.Add(2 * time.Second),
		CreatedAt:   now,
	}

	// Create a TaskQueue and add the task.
	tq := runtime.NewTaskQueue()
	if err := tq.AddTask(task); err != nil {
		t.Fatalf("Error adding task %s: %v", task.ID, err)
	}
	wg.Add(len(task.AssigneeIDs))

	// Create and start the scheduler.
	s := runtime.NewScheduler(tq)
	go s.Start()

	// Wait for a sufficient time for the task to be executed.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Task completed.
	case <-time.After(360 * time.Second):
		t.Fatal("Test timed out waiting for task to complete")
	}

	// Stop the scheduler.
	s.Stop()
}
