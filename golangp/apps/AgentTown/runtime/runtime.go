package runtime

import (
	"context"
	"fmt"
	"hajime/golangp/apps/AgentTown/agent"
	"hajime/golangp/apps/AgentTown/config"
	"hajime/golangp/apps/AgentTown/task"
	"sync"
	"time"
)

type Runtime struct {
	agents    map[string]*agent.Agent
	mu        sync.RWMutex
	CreatedAt time.Time
}

var rt *Runtime
var once sync.Once

// AddAgent adds a new agent to the runtime
func AddAgentByConfig(cfg *config.Config) {
	r := GetInstance()
	r.mu.Lock()
	defer r.mu.Unlock()
	newAgent := agent.NewAgent(cfg)
	r.agents[newAgent.ID] = newAgent
}

func AddAgent(ag *agent.Agent) {
	r := GetInstance()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[ag.ID] = ag
}

// GetAgents returns all agents
func GetAgents() map[string]*agent.Agent {
	r := GetInstance()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.agents
}

// GetInstance returns the singleton instance of Runtime
func GetInstance() *Runtime {
	once.Do(func() {
		rt = &Runtime{
			agents:    make(map[string]*agent.Agent),
			CreatedAt: time.Now(),
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
			ag.MarkAsDone()
		}
	}()
	wg.Wait()
}

// AssignTask assigns a task to a specific agent
func AssignTaskByAgentName(agentName string, tsk *task.Task) {
	ag := GetAgentByName(agentName)
	if ag != nil {
		ag.AssignTask(tsk)
	}

}

// AssignTaskByID assigns a task to a specific agent
func AssignTaskByID(agentID string, tsk *task.Task) {
	ag := GetAgentByID(agentID)
	if ag != nil {
		ag.AssignTask(tsk)
	}
}

func GetAgentByID(agentID string) *agent.Agent {
	r := GetInstance()
	r.mu.RLock()
	defer r.mu.RUnlock()
	if ag, exists := r.agents[agentID]; exists {
		return ag
	}
	return nil
}

func ActivateAgentByID(agentID string) {
	r := GetInstance()
	r.mu.RLock()
	defer r.mu.RUnlock()
	ag := GetAgentByID(agentID)
	if ag != nil {
		ag.Activate()
	}
}

func DeactivateAgentByID(agentID string) {
	r := GetInstance()
	r.mu.RLock()
	defer r.mu.RUnlock()
	ag := GetAgentByID(agentID)
	if ag != nil {
		ag.Deactivate()
	}
}

func GetAgentByName(agentName string) *agent.Agent {
	r := GetInstance()
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ag := range r.agents {
		if ag.Name == agentName {
			return ag
		}
	}
	return nil
}

// LogActivity logs activities of all agents
func LogActivity() {
	r := GetInstance()
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ag := range r.agents {
		if ag.IsActiveAgent() {
			fmt.Printf("Agent %s is active\n", ag.Name)
		} else {
			fmt.Printf("Agent %s is inactive\n", ag.Name)
		}
	}
}
