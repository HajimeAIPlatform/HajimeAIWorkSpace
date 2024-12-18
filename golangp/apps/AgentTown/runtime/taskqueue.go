package runtime

import (
	"container/heap"
	"errors"
	"sync"
	"time"

	"hajime/golangp/apps/AgentTown/task"
)

// TaskQueue is a thread-safe task queue that assigns tasks to agents.
type TaskQueue struct {
	tasks []*task.Task // The list of tasks.
	mu    sync.Mutex   // Mutex to protect concurrent access to the task queue.
	cond  *sync.Cond   // Condition variable for managing worker goroutine.
	timer *time.Timer  // Timer to trigger the execution of the nearest task.
}

// NewTaskQueue creates a new task queue.
func NewTaskQueue() *TaskQueue {
	tq := &TaskQueue{
		tasks: make([]*task.Task, 0),
	}
	tq.cond = sync.NewCond(&tq.mu)
	return tq
}

// AddTask adds a task to the queue.
func (tq *TaskQueue) AddTask(task *task.Task) error {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	// Check if the task's execution time has already passed.
	if task.ExecuteTime.Before(time.Now()) {
		return errors.New("task's execution time has already passed")
	}

	heap.Push(&taskHeap{tq: tq}, task)
	tq.cond.Signal() // Wake up the potentially waiting worker goroutine.

	// If the newly added task is the nearest to be executed, reset the timer.
	if tq.timer != nil && task.ExecuteTime.Before(tq.tasks[0].ExecuteTime) {
		tq.timer.Stop()
		tq.timer = nil
	}

	return nil
}

// GetDueTasks retrieves all tasks that are due for execution.
func (tq *TaskQueue) GetDueTasks() []*task.Task {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	var dueTasks []*task.Task
	now := time.Now()

	for len(tq.tasks) > 0 && tq.tasks[0].ExecuteTime.Before(now) {
		dueTasks = append(dueTasks, heap.Pop(&taskHeap{tq: tq}).(*task.Task))
	}

	return dueTasks
}

// RemoveTask removes a specific task from the queue.
func (tq *TaskQueue) RemoveTask(taskID string) error {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	for i, t := range tq.tasks {
		if t.ID == taskID {
			// Remove the task from the slice without using heap.Remove to avoid unnecessary reordering.
			tq.tasks = append(tq.tasks[:i], tq.tasks[i+1:]...)
			heap.Init(&taskHeap{tq: tq}) // Reinitialize the heap after removing an element.
			return nil
		}
	}

	return errors.New("task not found")
}

// taskHeap is a helper type that implements heap.Interface to sort tasks by ExecuteTime.
type taskHeap struct {
	tq *TaskQueue
}

func (h taskHeap) Len() int { return len(h.tq.tasks) }
func (h taskHeap) Less(i, j int) bool {
	return h.tq.tasks[i].ExecuteTime.Before(h.tq.tasks[j].ExecuteTime)
}
func (h taskHeap) Swap(i, j int) {
	h.tq.tasks[i], h.tq.tasks[j] = h.tq.tasks[j], h.tq.tasks[i]
}

func (h *taskHeap) Push(x interface{}) {
	h.tq.tasks = append(h.tq.tasks, x.(*task.Task))
}

func (h *taskHeap) Pop() interface{} {
	old := h.tq.tasks
	n := len(old)
	x := old[n-1]
	h.tq.tasks = old[0 : n-1]
	return x
}
