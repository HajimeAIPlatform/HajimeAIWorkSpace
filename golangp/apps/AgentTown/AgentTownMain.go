package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Message structure to model communication between agents
type Message struct {
	Sender    string
	Recipient string
	Content   string
}

// Agent represents a single autonomous entity
type Agent struct {
	Name       string
	MessageCh  chan Message // Channel to receive messages
	ProcessCh  chan string  // Channel for tasks to process
	Done       chan bool    // Signaling channel indicating work completion
	Receivers  map[string]*Agent
}

// NewAgent creates a new agent
func NewAgent(name string) *Agent {
	return &Agent{
		Name:      name,
		MessageCh: make(chan Message, 10), // Buffered channel for communication
		ProcessCh: make(chan string, 5),  // Buffered channel for tasks
		Done:      make(chan bool),
		Receivers: make(map[string]*Agent),
	}
}

// RegisterReceiver binds communication with other agents
func (agent *Agent) RegisterReceiver(other *Agent) {
	agent.Receivers[other.Name] = other
}

// Start runs the agent's loop as a goroutine
func (agent *Agent) Start(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	go func() {
		for {
			select {
			// Process incoming messages
			case msg := <-agent.MessageCh:
				fmt.Printf("[%s] Received message from %s: %s\n", agent.Name, msg.Sender, msg.Content)

			// Process assigned tasks
			case task := <-agent.ProcessCh:
				fmt.Printf("[%s] Processing task: %s\n", agent.Name, task)
				time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second) // Simulate task processing
				fmt.Printf("[%s] Finished task: %s\n", agent.Name, task)

			// Graceful shutdown on context cancellation
			case <-ctx.Done():
				fmt.Printf("[%s] Shutting down...\n", agent.Name)
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
func (agent *Agent) AssignTask(task string) {
	agent.ProcessCh <- task
}

func main() {
	// Random seed for generating delays
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Create agents
	agentA := NewAgent("Agent_A")
	agentB := NewAgent("Agent_B")
	agentC := NewAgent("Agent_C")

	// Register receivers (peer-to-peer communication)
	agentA.RegisterReceiver(agentB)
	agentB.RegisterReceiver(agentC)
	agentC.RegisterReceiver(agentA)

	// Start agents
	wg.Add(3)
	go agentA.Start(&wg, ctx)
	go agentB.Start(&wg, ctx)
	go agentC.Start(&wg, ctx)

	// Assign tasks to agents
	go func() {
		time.Sleep(2 * time.Second)
		agentA.AssignTask("Fetch Data")
		agentB.AssignTask("Process Data")
		agentC.AssignTask("Export Results")
	}()

	// Allow system to run for 10 seconds before shutting down
	time.Sleep(10 * time.Second)
	fmt.Println("Shutting down agents...")
	cancel()
	wg.Wait()
	fmt.Println("All agents are shut down.")
}