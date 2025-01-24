package test

import (
	"context"
	"fmt"
	"hajime/golangp/apps/AgentTown/agent"
	"hajime/golangp/apps/AgentTown/runtime"
	"hajime/golangp/apps/AgentTown/task"
	"hajime/golangp/apps/trading/dex/solana/raydium"
	"log"
	"os"

	"sync"
	"testing"
	"time"
)

func SetEnvironmentVariables(grpcServer string) {
	err := os.Setenv("GRPC_SERVER", grpcServer)
	if err != nil {
		log.Fatalf("Error setting GRPC_SERVER: %v", err)
	}
}

func init() {
	SetEnvironmentVariables("0.0.0.0:50051")
}

// Please set these two environment variables under execution. default values are shown below.
// export LOG_FILE_PATH=/tmp/logs
// GRPC_SERVER=0.0.0.0:50051
func ExecuteSolanaSwap(tokenIn string, tokenOut string, privateKey string, amountIn int64, microLamports int64, wg *sync.WaitGroup) {
	fmt.Printf("Executing solana swap from %s to %s\n", tokenIn, tokenOut)
	defer wg.Done()

	txId, err := raydium.CallSwap(tokenIn, tokenOut, privateKey, amountIn, microLamports)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Swap successful. TxId: %s\n", txId)
}

func TestSolanaSwap(t *testing.T) {
	// Use a WaitGroup to wait for tasks to complete.
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create agents and start them.
	agentA := agent.NewAgent(nil)
	runtime.AddAgent(agentA)
	runtime.StartAgents(ctx)

	// Set the private key for the agent.
	privateKey := "YOUR_PRIVATE_KEY"

	// Create tasks with different execution times.
	now := time.Now()
	task1 := &task.Task{
		ID:          "solana swap1",
		AssigneeIDs: []string{agentA.ID},
		Description: "Task for solana swap SOL to USDT",
		Execute: func(args ...any) {
			params := args[0].([]any)
			tokenIn := params[0].(string)
			tokenOut := params[1].(string)
			privateKey := params[2].(string)
			amountIn := params[3].(int64)
			microLamports := params[4].(int64)
			wg := params[5].(*sync.WaitGroup)

			ExecuteSolanaSwap(tokenIn, tokenOut, privateKey, amountIn, microLamports, wg)
		},
		Parameters:  []any{"SOL", "USDT", privateKey, int64(100000), int64(300000), &wg},
		ExecuteTime: now.Add(2 * time.Second),
		CreatedAt:   now,
	}

	task2 := &task.Task{
		ID:          "solana swap2",
		AssigneeIDs: []string{agentA.ID},
		Description: "Task for solana swap SOL to USDC",
		Execute: func(args ...any) {
			params := args[0].([]any)
			tokenIn := params[0].(string)
			tokenOut := params[1].(string)
			privateKey := params[2].(string)
			amountIn := params[3].(int64)
			microLamports := params[4].(int64)
			wg := params[5].(*sync.WaitGroup)
			ExecuteSolanaSwap(tokenIn, tokenOut, privateKey, amountIn, microLamports, wg)
		},
		Parameters:  []any{"SOL", "USDC", privateKey, int64(100000), int64(300000), &wg},
		ExecuteTime: now.Add(1 * time.Second),
		CreatedAt:   now,
	}

	// Create a TaskQueue and add tasks.
	tq := runtime.NewTaskQueue()
	for _, task := range []*task.Task{task1, task2} {
		if err := tq.AddTask(task); err != nil {
			t.Fatalf("Error adding task %s: %v", task.ID, err)
		}
		wg.Add(len(task.AssigneeIDs))
	}

	// Create and start the scheduler.
	s := runtime.NewScheduler(tq)
	go s.Start()

	// Wait for a sufficient time for the tasks to be executed.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All tasks completed.
	case <-time.After(360 * time.Second):
		t.Fatal("Test timed out waiting for tasks to complete")
	}

	// Stop the scheduler.
	s.Stop()
}
