package test

import (
	"context"
	"fmt"
	"hajime/golangp/apps/AgentTown/agent"
	"hajime/golangp/apps/AgentTown/runtime"
	"hajime/golangp/apps/AgentTown/task"
	"hajime/golangp/apps/trading/dex/solana/raydium"
	"log"
	"math"
	"sync"
	"testing"
	"time"
)

func TestCreatePool(t *testing.T) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create agents and start them.
	agentA := agent.NewAgent(nil)
	runtime.AddAgent(agentA)
	runtime.StartAgents(ctx)

	// Set the parameters for creating a pool.
	privateKey := "WEZT6Wdau5GDz2HCygJxZheWzZodkGUX5Yz3bgqWJCuJL7ccVk4cfP1oFyDxpv2Ak8hacyvTyPspQ3f66oxNfHd"
	mintAAddress := "3uAv9qSsUdz2RFkVx99Fe81dqpDUbxQRzN1kNwdykTuf"
	mintADecimals := 6
	mintAInitialAmount := int64(1_000 * math.Pow(10, 6))
	mintBAddress := "So11111111111111111111111111111111111111112"
	mintBDecimals := 9
	mintBInitialAmount := int64(1 * math.Pow(10, 9))
	marketId := "DRSbrtzZPoAwwJ36kEfRfQka5jaBmfcm46NmBt5ASnSu"

	// Create a task for creating a pool.
	now := time.Now()
	task := &task.Task{
		ID:          "create pool",
		AssigneeIDs: []string{agentA.ID},
		Description: "Task for creating a new pool",
		Execute: func(args ...any) {
			params := args[0].([]any)
			privateKey := params[0].(string)
			mintAAddress := params[1].(string)
			mintADecimals := params[2].(int)
			mintAInitialAmount := params[3].(int64)
			mintBAddress := params[4].(string)
			mintBDecimals := params[5].(int)
			mintBInitialAmount := params[6].(int64)
			marketId := params[7].(string)
			wg := params[8].(*sync.WaitGroup)

			defer wg.Done()

			res, err := raydium.CallCreatePool(privateKey, mintAAddress, mintADecimals, mintAInitialAmount, mintBAddress, mintBDecimals, mintBInitialAmount, marketId)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			fmt.Printf("Pool created successfully. PoolId: %s, TxId: %s\n", res.PoolId, res.TxId)
		},
		Parameters:  []any{privateKey, mintAAddress, mintADecimals, mintAInitialAmount, mintBAddress, mintBDecimals, mintBInitialAmount, marketId, &wg},
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
