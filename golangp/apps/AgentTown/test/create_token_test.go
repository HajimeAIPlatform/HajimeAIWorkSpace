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

func TestCreateToken(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create agents and start them.
	agentA := agent.NewAgent(nil)
	runtime.AddAgent(agentA)
	runtime.StartAgents(ctx)

	// Set the private key for the agent.
	privateKey := "WEZT6Wdau5GDz2HCygJxZheWzZodkGUX5Yz3bgqWJCuJL7ccVk4cfP1oFyDxpv2Ak8hacyvTyPspQ3f66oxNfHd"
	tokenName := "HAJIME_P"
	tokenSymbol := "HAJIME_P"
	description := "HAJIME_P.\n\nTELEGRAM:  \nTWITTER:  \n WEBSITE: "
	uri := "https://devixyz.github.io/telegram/hajime.json"
	tokenSupply := int64(1_500_000)
	tokenDecimals := int64(6)

	// Create a task for creating a token.
	now := time.Now()
	task := &task.Task{
		ID:          "create token",
		AssigneeIDs: []string{agentA.ID},
		Description: "Task for creating a new token",
		Execute: func(args ...any) {
			params := args[0].([]any)
			privateKey := params[0].(string)
			tokenName := params[1].(string)
			tokenSymbol := params[2].(string)
			description := params[3].(string)
			uri := params[4].(string)
			tokenSupply := params[5].(int64)
			tokenDecimals := params[6].(int64)
			wg := params[7].(*sync.WaitGroup)

			defer wg.Done()

			res, err := raydium.CallCreateToken(privateKey, tokenName, tokenSymbol, description, uri, tokenSupply, tokenDecimals)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			fmt.Printf("Token created successfully. TxId: %s\n", res.TxId)
		},
		Parameters:  []any{privateKey, tokenName, tokenSymbol, description, uri, tokenSupply, tokenDecimals, &wg},
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
