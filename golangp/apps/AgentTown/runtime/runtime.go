package runtime

import (
	"context"
	"fmt"
	"hajime/golangp/apps/AgentTown/agent"
	"hajime/golangp/apps/AgentTown/config"
	"hajime/golangp/apps/AgentTown/task"
	"hajime/golangp/apps/AgentTown/telemetry"
	"sync"
)

type Runtime struct {
	agents map[string]*agent.Agent
	mu     sync.Mutex
}

var rt *Runtime
var once sync.Once

// AddAgent adds a new agent to the runtime
func AddAgent(cfg config.Config) {
	r := GetInstance()
	r.mu.Lock()
	defer r.mu.Unlock()
	newAgent := agent.NewAgent(cfg.Name)
	r.agents[cfg.Name] = newAgent
}

// GetInstance returns the singleton instance of Runtime
func GetInstance() *Runtime {
	once.Do(func() {
		rt = &Runtime{
			agents: make(map[string]*agent.Agent),
		}
	})
	return rt
}

// StartAgents starts all agents
func StartAgents(ctx context.Context) {
	r := GetInstance()
	fmt.Printf("Starting all agents from runtime ...\n")
	var wg sync.WaitGroup
	for _, ag := range r.agents {
		wg.Add(1)
		go ag.Start(&wg, ctx)
	}
	go func() {
		<-ctx.Done()
		fmt.Printf("Shutting down all agents from runtime ...\n")
		for _, ag := range r.agents {
			close(ag.Done)
		}
	}()
	wg.Wait()
}

// AssignTask assigns a task to a specific agent
func AssignTask(agentName string, t task.Task) {
	r := GetInstance()
	r.mu.Lock()
	defer r.mu.Unlock()
	if ag, exists := r.agents[agentName]; exists {
		ag.AssignTask(t.Description)
		telemetry.RecordMetricInc("tasks_assigned", 1)
	}
}

// LogActivity logs activities of all agents
func LogActivity() {
	r := GetInstance()
	for _, ag := range r.agents {
		if ag.IsActive {
			fmt.Printf("Agent %s is active\n", ag.Name)
		} else {
			fmt.Printf("Agent %s is inactive\n", ag.Name)
		}
	}
}
