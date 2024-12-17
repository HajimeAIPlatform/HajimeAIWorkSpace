package agent

import (
	"context"
	"fmt"
	"hajime/golangp/apps/AgentTown/config"
	"hajime/golangp/apps/AgentTown/task"
	"hajime/golangp/apps/AgentTown/telemetry"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Message structure to model communication between agents
type Message struct {
	Sender    string
	Recipient string
	Content   string
	Timestamp time.Time
}

// Agent represents a single autonomous entity
// Name is the readable identifier for the agent, may not be unique
// ID is the unique identifier for the agent
// AI logic can be implemented in the Execute function and message handlers
type Agent struct {
	ID        string
	Name      string
	MessageCh chan Message    // Channel to receive messages
	ProcessCh chan *task.Task // Channel for tasks to process
	Done      chan struct{}   // Signaling channel indicating work completion
	Receivers map[string]*Agent
	IsActive  bool
	Config    *config.Config
	CreatedAt time.Time
	mu        sync.RWMutex
}

// NewAgent creates a new agent
func NewAgent(config *config.Config) *Agent {
	return &Agent{
		ID:        uuid.New().String(),
		Name:      config.Name + uuid.New().String(),
		MessageCh: make(chan Message, 10),    // Buffered channel for communication
		ProcessCh: make(chan *task.Task, 20), // Buffered channel for tasks
		Done:      make(chan struct{}),
		Receivers: make(map[string]*Agent),
		Config:    config,
		IsActive:  false,
		CreatedAt: time.Now(),
	}
}

// RegisterReceiver binds communication with other agents
func (agent *Agent) RegisterReceiver(other *Agent) {
	agent.mu.Lock()
	defer agent.mu.Unlock()
	agent.Receivers[other.Name] = other
}

// Start runs the agent's loop as a goroutine
func (agent *Agent) Start(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	agent.Activate()
	go func() {
		for {
			select {

			// Done should be processed before tasks
			case <-agent.Done:
				fmt.Printf("[%s] Done signal received, ignoring remaining tasks\n", agent.Name)
				fmt.Printf("[%s] Quiting the agent goroutine\n", agent.Name)
				agent.Deactivate()
				return

			// Process incoming messages
			case msg := <-agent.MessageCh:
				fmt.Printf("[%s] Received message from %s: %s\n", agent.Name, msg.Sender, msg.Content)
				continue

			// Process assigned tasks
			case tsk := <-agent.ProcessCh:
				fmt.Printf("[%s] Processing task %s : %s\n", agent.Name, tsk.ID, tsk.Description)
				if tsk.Execute != nil {
					tsk.Execute(tsk.Parameters, agent.Config.PrivateData)
				}
				time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second) // Simulate task processing
				telemetry.RecordMetricInc("tasks_completed", 1)
				fmt.Printf("[%s] Finished task %s : %s\n", agent.Name, tsk.ID, tsk.Description)
				continue

			// Graceful shutdown on context cancellation
			case <-ctx.Done():
				fmt.Printf("[%s] Context cancelled\n", agent.Name)
				fmt.Printf("[%s] Quiting the agent goroutine\n", agent.Name)
				agent.Deactivate()
				return
			}
		}
	}()

	// Simulate sending messages to other agents
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Randomly send a message every 2-3 seconds
				time.Sleep(time.Duration(rand.Intn(3)+2) * time.Second)
				for _, receiver := range agent.Receivers {
					msg := Message{
						Sender:    agent.Name,
						Recipient: receiver.Name,
						Content:   fmt.Sprintf("Hello from %s!", agent.Name),
						Timestamp: time.Now(),
					}
					receiver.MessageCh <- msg
				}
			}
		}
	}()
}

// AssignTask sends a task to the agent for processing
func (agent *Agent) AssignTask(tsk *task.Task) {
	if agent.IsActiveAgent() {
		agent.ProcessCh <- tsk
		fmt.Printf("Task assigned to agent %s : %s\n", agent.Name, tsk.Description)
		return
	}
	fmt.Printf("Agent %s is not active, task cannot be assigned: %s \n", agent.Name, tsk.Description)
}

func (agent *Agent) GetConfig() *config.Config {
	agent.mu.RLock()
	defer agent.mu.RUnlock()
	return agent.Config
}

func (agent *Agent) GetID() string {
	agent.mu.RLock()
	defer agent.mu.RUnlock()
	return agent.ID
}

func (agent *Agent) GetName() string {
	agent.mu.RLock()
	defer agent.mu.RUnlock()
	return agent.Name
}

// MarkAsDone signals the agent to stop processing tasks and exit
func (agent *Agent) MarkAsDone() {
	agent.mu.Lock()
	defer agent.mu.Unlock()
	close(agent.Done)
}

func (agent *Agent) IsActiveAgent() bool {
	agent.mu.RLock()
	defer agent.mu.RUnlock()
	return agent.IsActive
}

func (agent *Agent) Deactivate() {
	agent.mu.Lock()
	defer agent.mu.Unlock()
	agent.IsActive = false
	telemetry.RecordMetricInc("agents_active", -1)

}

func (agent *Agent) Activate() {
	agent.mu.Lock()
	defer agent.mu.Unlock()
	agent.IsActive = true
	telemetry.RecordMetricInc("agents_active", 1)
}
