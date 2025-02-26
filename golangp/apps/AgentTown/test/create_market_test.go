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

func TestCreateMarket(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create agents and start them.
	agentA := agent.NewAgent(nil)
	runtime.AddAgent(agentA)
	runtime.StartAgents(ctx)

	// Set the parameters for creating a market.
	privateKey := "WEZT6Wdau5GDz2HCygJxZheWzZodkGUX5Yz3bgqWJCuJL7ccVk4cfP1oFyDxpv2Ak8hacyvTyPspQ3f66oxNfHd"
	mintAAddress := "3uAv9qSsUdz2RFkVx99Fe81dqpDUbxQRzN1kNwdykTuf"
	mintADecimals := 6
	mintBAddress := "So11111111111111111111111111111111111111112"
	mintBDecimals := 9

	// Create a task for creating a market.
	now := time.Now()
	task := &task.Task{
		ID:          "create market",
		AssigneeIDs: []string{agentA.ID},
		Description: "Task for creating a new market",
		Execute: func(args ...any) {
			params := args[0].([]any)
			privateKey := params[0].(string)
			mintAAddress := params[1].(string)
			mintADecimals := params[2].(int)
			mintBAddress := params[3].(string)
			mintBDecimals := params[4].(int)
			wg := params[5].(*sync.WaitGroup)

			defer wg.Done()

			res, err := raydium.CallCreateMarket(privateKey, mintAAddress, mintADecimals, mintBAddress, mintBDecimals)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			fmt.Printf("Market created successfully. MarketId: %s, TxIds: %s\n", res.MarketId, res.TxIds)
		},
		Parameters:  []any{privateKey, mintAAddress, mintADecimals, mintBAddress, mintBDecimals, &wg},
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
