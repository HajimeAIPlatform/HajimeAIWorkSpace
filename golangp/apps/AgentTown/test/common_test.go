package test

import (
	"context"
	"fmt"
	"hajime/golangp/apps/AgentTown/agent"
	"hajime/golangp/apps/AgentTown/runtime"
	"hajime/golangp/apps/AgentTown/task"
	"sync"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestSmokeTest(t *testing.T) {
	assert.Equal(t, 10, 9+1)
}

func TestSmokeTest2(t *testing.T) {
	assert.Equal(t, 12, 9+3)
}

func TestTaskGeneration(t *testing.T) {
	tsk := task.NewTask("Test Task")
	fmt.Printf("Task ID: %s\n\n", tsk.ID)
	assert.Equal(t, "Test Task", tsk.Description)
}

func TestScheduler(t *testing.T) {
	// Use a WaitGroup to wait for tasks to complete.
	var wg sync.WaitGroup

	// Create a mock task function that signals the WaitGroup when done.
	mockTaskFunc := func(args ...any) {
		fmt.Println("Mock task executed with arguments:", args)
		if len(args) > 0 {
			fmt.Println("Mock Task ID:", args[0])
		}
		wg.Done()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create agents and start them.
	agentA := agent.NewAgent(nil)
	agentB := agent.NewAgent(nil)
	runtime.AddAgent(agentA)
	runtime.AddAgent(agentB)
	runtime.StartAgents(ctx)

	// Create tasks with different execution times.
	now := time.Now()
	task1 := &task.Task{
		ID:          "task1",
		AssigneeIDs: []string{agentA.ID},
		Description: "Task for agent1",
		Execute:     mockTaskFunc,
		Parameters:  []any{"task1"},
		ExecuteTime: now.Add(1 * time.Second),
		CreatedAt:   now,
	}
	task2 := &task.Task{
		ID:          "task2",
		AssigneeIDs: []string{agentB.ID},
		Description: "Task for agent2",
		Execute:     mockTaskFunc,
		Parameters:  []any{"task2"},
		ExecuteTime: now.Add(2 * time.Second),
		CreatedAt:   now,
	}
	task3 := &task.Task{
		ID:          "task3",
		AssigneeIDs: []string{agentA.ID, agentB.ID},
		Description: "Task for both agent1 and agent2",
		Execute:     mockTaskFunc,
		Parameters:  []any{"task3"},
		ExecuteTime: now.Add(4 * time.Second),
		CreatedAt:   now,
	}

	// task4 := &task.Task{
	// 	ID:          "task4",
	// 	AssigneeIDs: []string{agentA.ID, agentB.ID},
	// 	Description: "Task for both agent1 and agent2",
	// 	Execute:     task.TestBinanceConnectivy,
	// 	Parameters:  []any{"task4"},
	// 	ExecuteTime: now.Add(3 * time.Second),
	// 	CreatedAt:   now,
	// }

	// Create a TaskQueue and add tasks.
	tq := runtime.NewTaskQueue()
	for _, task := range []*task.Task{task1, task2, task3} {
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
	case <-time.After(12 * time.Second):
		t.Fatal("Test timed out waiting for tasks to complete")
	}

	// Stop the scheduler.
	s.Stop()
}
