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

	rt := runtime.GetInstance()

	// Create and add agents
	rt.AddAgent(config.Config{Name: "Agent_A"})
	rt.AddAgent(config.Config{Name: "Agent_B"})
	rt.AddAgent(config.Config{Name: "Agent_C"})

	// Start agents
	go rt.StartAgents(ctx)

	// Assign tasks to agents
	go func() {
		time.Sleep(2 * time.Second)
		rt.AssignTask("Agent_A", task.Task{ID: "1", Content: "Fetch Data"})
		rt.AssignTask("Agent_B", task.Task{ID: "2", Content: "Process Data"})
		rt.AssignTask("Agent_C", task.Task{ID: "3", Content: "Export Results"})
	}()

	// Log activities
	go func() {
		for {
			time.Sleep(5 * time.Second)
			rt.LogActivity()
		}
	}()

	// Start tm monitoring
	tm := telemetry.NewTelemetry()
	tm.Monitor()

	// For test purpose
	time.Sleep(15 * time.Second)
	fmt.Printf("System shutting down...\n")
	cancel()
	time.Sleep(10 * time.Second)
}
