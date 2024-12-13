package runtime

import (
	"context"
	"fmt"
	"hajime/golangp/apps/AgentTown/agent"
	"hajime/golangp/apps/AgentTown/config"
	"hajime/golangp/apps/AgentTown/task"
	"sync"
)

type Runtime struct {
	agents map[string]*agent.Agent
	mu     sync.Mutex
}

var instance *Runtime
var once sync.Once

// GetInstance returns the singleton instance of Runtime
func GetInstance() *Runtime {
	once.Do(func() {
		instance = &Runtime{
			agents: make(map[string]*agent.Agent),
		}
	})
	return instance
}

// AddAgent adds a new agent to the runtime
func (r *Runtime) AddAgent(cfg config.Config) {
	r.mu.Lock()
	defer r.mu.Unlock()
	newAgent := agent.NewAgent(cfg.Name)
	r.agents[cfg.Name] = newAgent
}

// StartAgents starts all agents
func (r *Runtime) StartAgents(ctx context.Context) {
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
func (r *Runtime) AssignTask(agentName string, t task.Task) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ag, exists := r.agents[agentName]; exists {
		ag.AssignTask(t.Content)
	}
}

// LogActivity logs activities of all agents
func (r *Runtime) LogActivity() {
	for _, ag := range r.agents {
		if ag.IsActive {
			fmt.Printf("Agent %s is active\n", ag.Name)
		} else {
			fmt.Printf("Agent %s is inactive\n", ag.Name)
		}
	}
}
