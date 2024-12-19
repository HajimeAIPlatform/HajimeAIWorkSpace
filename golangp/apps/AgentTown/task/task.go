package task

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Task represents a task assigned to an agent
type Task struct {
	ID          string
	AssigneeIDs []string
	Description string
	Execute     func(...any) // Function to be executed
	Parameters  []any        // Parameters to pass to the function
	ExecuteTime time.Time
	CreatedAt   time.Time
	Done        bool
	mu          sync.RWMutex
}

// Even though task is parsed by pointer, task should not be shared between goroutines, same task should be cloned and passed to each agent
func NewTask(description string) *Task {
	return &Task{
		ID:          uuid.New().String(),
		AssigneeIDs: []string{},
		Description: description,
		Execute:     nil,
		Parameters:  nil,
		CreatedAt:   time.Now(),
		Done:        false,
	}
}

func setAssigneeIDs(t *Task, ids []string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.AssigneeIDs = ids
}

func addAssigneeID(t *Task, id string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.AssigneeIDs = append(t.AssigneeIDs, id)
}

func setExecute(t *Task, f func(...any)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Execute = f
}

func setParameters(t *Task, params []any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Parameters = params
}

func setExecuteTime(t *Task, et time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ExecuteTime = et
}

func getAssigneeIDs(t *Task) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.AssigneeIDs
}

func getExecute(t *Task) func(...any) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Execute
}

func getParameters(t *Task) []any {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Parameters
}

// CloneToNewTask clones a task to a new task
func CloneToNewTask(t *Task) *Task {
	return &Task{
		ID:          uuid.New().String(),
		AssigneeIDs: []string{},
		Description: t.Description,
		Execute:     t.Execute,
		Parameters:  t.Parameters,
		CreatedAt:   time.Now(),
		Done:        false,
	}
}
