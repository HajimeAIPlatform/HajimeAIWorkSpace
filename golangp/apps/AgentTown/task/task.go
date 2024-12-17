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
}

func NewTask(content string) *Task {
	return &Task{
		ID:          uuid.New().String(),
		Description: content,
		Execute:     nil,
		Parameters:  nil,
		Timestamp:   time.Now(),
	}
}
