package main

import (
	"context"
	"fmt"
	"hajime/golangp/apps/AgentTown/config"
	"hajime/golangp/apps/AgentTown/runtime"
	"hajime/golangp/apps/AgentTown/task"
	"hajime/golangp/apps/AgentTown/telemetry"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Create and add agents
	runtime.AddAgent(config.Config{Name: "Agent_A"})
	runtime.AddAgent(config.Config{Name: "Agent_B"})
	runtime.AddAgent(config.Config{Name: "Agent_C"})

	// Start agents
	go runtime.StartAgents(ctx)

	// Assign tasks to agents
	go func() {
		time.Sleep(2 * time.Second)
		runtime.AssignTask("Agent_A", task.Task{ID: "1", Content: "Fetch Data"})
		runtime.AssignTask("Agent_B", task.Task{ID: "2", Content: "Process Data"})
		runtime.AssignTask("Agent_C", task.Task{ID: "3", Content: "Export Results"})
	}()

	// Log activities
	go func() {
		for {
			time.Sleep(5 * time.Second)
			runtime.LogActivity()
		}
	}()

	// Start telemetry monitoring
	telemetry.Monitor(5 * time.Second)

	// For test purpose
	time.Sleep(15 * time.Second)
	fmt.Printf("System shutting down...\n")
	cancel()
	time.Sleep(10 * time.Second)
}
