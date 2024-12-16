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
}

// Agent represents a single autonomous entity
type Agent struct {
	ID        string
	Name      string
	MessageCh chan Message    // Channel to receive messages
	ProcessCh chan *task.Task // Channel for tasks to process
	Done      chan struct{}   // Signaling channel indicating work completion
	Receivers map[string]*Agent
	IsActive  bool
	Config    config.Config
	mu        sync.Mutex
}

// NewAgent creates a new agent
func NewAgent(name string) *Agent {
	return &Agent{
		ID:        uuid.New().String(),
		Name:      name,
		MessageCh: make(chan Message, 10),    // Buffered channel for communication
		ProcessCh: make(chan *task.Task, 20), // Buffered channel for tasks
		Done:      make(chan struct{}),
		Receivers: make(map[string]*Agent),
		IsActive:  false,
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
	agent.IsActive = true
	telemetry.RecordMetricInc("agents_active", 1)
	go func() {
		for {
			select {
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
				fmt.Printf("[%s] Finished task %s : %s\n", agent.Name, tsk.ID, tsk.Description)
				continue

			case <-agent.Done:
				fmt.Printf("[%s] Done signal received\n", agent.Name)
				fmt.Printf("[%s] Quiting the agent goroutine\n", agent.Name)
				agent.IsActive = false
				telemetry.RecordMetricInc("agents_active", -1)
				return

			// Graceful shutdown on context cancellation
			case <-ctx.Done():
				fmt.Printf("[%s] Context cancelled\n", agent.Name)
				fmt.Printf("[%s] Quiting the agent goroutine\n", agent.Name)
				agent.IsActive = false
				telemetry.RecordMetricInc("agents_active", -1)
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
					}
					receiver.MessageCh <- msg
				}
			}
		}
	}()
}

// AssignTask sends a task to the agent for processing
func (agent *Agent) AssignTask(taskDescription string) {
	agent.ProcessCh <- task.NewTask(taskDescription)
}
