package task

import "github.com/google/uuid"

// Task represents a task assigned to an agent
type Task struct {
	ID          string
	Description string
	Execute     func(...any) // Function to be executed
	Parameters  []any        // Parameters to pass to the function
}

func NewTask(content string) *Task {
	return &Task{
		ID:          uuid.New().String(),
		Description: content,
		Execute:     nil,
		Parameters:  nil,
	}
}
