package task

import (
	"time"

	"github.com/google/uuid"
)

// Task represents a task assigned to an agent
type Task struct {
	ID          string
	Description string
	Execute     func(...any) // Function to be executed
	Parameters  []any        // Parameters to pass to the function
	Timestamp   time.Time
	Done        bool
}

// Even though task is parsed by pointer, task should not be shared between goroutines, same task should be cloned and passed to each agent
func NewTask(description string) *Task {
	return &Task{
		ID:          uuid.New().String(),
		Description: description,
		Execute:     nil,
		Parameters:  nil,
		Timestamp:   time.Now(),
		Done:        false,
	}
}

// CloneToNewTask clones a task to a new task
func CloneToNewTask(t *Task) *Task {
	return &Task{
		ID:          uuid.New().String(),
		Description: t.Description,
		Execute:     t.Execute,
		Parameters:  t.Parameters,
		Timestamp:   time.Now(),
		Done:        false,
	}
}
