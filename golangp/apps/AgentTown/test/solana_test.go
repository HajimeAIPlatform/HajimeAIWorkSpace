package test

import (
	"context"
	"sync"
	"testing"
	"time"

	"hajime/golangp/apps/AgentTown/agent"
	"hajime/golangp/apps/AgentTown/runtime"
	"hajime/golangp/apps/AgentTown/task"
	"hajime/golangp/apps/solana"
)

func TestGenerateWallet(t *testing.T) {
	wallet, err := solana.GenerateWallet()
	if err != nil {
		t.Fatalf("Failed to generate wallet: %v", err)
	}

	if wallet.PublicKey == "" || wallet.PrivateKey == "" {
		t.Errorf("Generated wallet has empty public or private key")
	}
}

func TestGenerateMultipleWallets(t *testing.T) {
	count := 5
	wallets, err := solana.GenerateMultipleWallets(count)
	if err != nil {
		t.Fatalf("Failed to generate multiple wallets: %v", err)
	}

	if len(wallets) != count {
		t.Errorf("Expected %d wallets, but got %d", count, len(wallets))
	}

	for _, wallet := range wallets {
		if wallet.PublicKey == "" || wallet.PrivateKey == "" {
			t.Errorf("Generated wallet has empty public or private key")
		}
	}
}

func TestNewConnection(t *testing.T) {
	conn := solana.NewConnection("")
	if conn == nil {
		t.Errorf("Failed to create new connection")
	} else if conn.Client == nil {
		t.Errorf("Connection client is nil")
	}
}

func TestGetWalletInfo(t *testing.T) {
	conn := solana.NewConnection("")

	// 使用一个已知的公钥进行测试，例如 Solana 的系统程序公钥
	publicKey := "11111111111111111111111111111111"
	err := conn.GetWalletInfo(publicKey)
	if err != nil {
		t.Errorf("Failed to get wallet info: %v", err)
	}
}

func TestSolana(t *testing.T) {
	// Use a WaitGroup to wait for tasks to complete.
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create agents and start them.
	agentA := agent.NewAgent(nil)
	runtime.AddAgent(agentA)
	runtime.StartAgents(ctx)

	// Create tasks with different execution times.
	now := time.Now()
	task1 := &task.Task{
		ID:          "task1",
		AssigneeIDs: []string{agentA.ID},
		Description: "Generate 5 wallets for agent1",
		Execute: func(params ...any) {
			solana.GenerateMultipleWalletsWrapper(params...)
			wg.Done() // Call wg.Done() after the task is completed
		},
		Parameters:  []any{5},
		ExecuteTime: now.Add(1 * time.Second),
		CreatedAt:   now,
	}

	// Create a TaskQueue and add tasks.
	tq := runtime.NewTaskQueue()
	if err := tq.AddTask(task1); err != nil {
		t.Fatalf("Error adding task %s: %v", task1.ID, err)
	}
	wg.Add(len(task1.AssigneeIDs))

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
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out waiting for tasks to complete")
	}

	// Stop the scheduler.
	s.Stop()
}
