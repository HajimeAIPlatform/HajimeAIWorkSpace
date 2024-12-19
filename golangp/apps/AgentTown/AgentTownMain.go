package main

import (
	"context"
	"fmt"
	"hajime/golangp/apps/AgentTown/agent"
	"hajime/golangp/apps/AgentTown/config"
	"hajime/golangp/apps/AgentTown/runtime"
	"hajime/golangp/apps/AgentTown/task"
	"hajime/golangp/apps/AgentTown/telemetry"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.NewConfig("Config_Empty")

	// Create all agents using the same configuration
	agentA := agent.NewAgent(cfg)
	agentB := agent.NewAgent(cfg)
	agentC := agent.NewAgent(cfg)

	// Create and add agents
	runtime.AddAgent(agentA)
	runtime.AddAgent(agentB)
	runtime.AddAgent(agentC)

	// Start agents
	go runtime.StartAgents(ctx)

	// Assign tasks to agents
	go func() {
		time.Sleep(2 * time.Second)
		runtime.AssignTaskByAgentID(agentA.ID, task.NewTask("Collect Data"))
		runtime.AssignTaskByAgentName(agentB.Name, task.NewTask("Process Data"))
		runtime.AssignTaskByAgentID(agentC.ID, task.NewTask("Report Data"))
		runtime.DeactivateAgentByID(agentB.ID)
		runtime.AssignTaskByAgentID(agentB.ID, task.NewTask("Test Deactivation"))
		runtime.ActivateAgentByID(agentB.ID)
		runtime.AssignTaskByAgentID(agentB.ID, task.NewTask("Test Activation"))
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
