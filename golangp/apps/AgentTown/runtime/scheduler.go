package runtime

import (
	"time"
)

type Scheduler struct {
	taskQueue *TaskQueue
	stop      chan struct{}
}

type Agent struct {
	ID string
}

func NewScheduler(tq *TaskQueue) *Scheduler {
	// Create a new TaskQueue if tq is nil.
	if tq == nil {
		tq = NewTaskQueue()
	}

	return &Scheduler{
		taskQueue: tq,
		stop:      make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(1 * time.Second) // Check for due tasks every second.
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// due task is already being removed from the queue
			dueTasks := s.taskQueue.GetDueTasks()
			for _, t := range dueTasks {
				for _, agentID := range t.AssigneeIDs {
					AssignTaskByAgentID(agentID, t) // Assign the task to the agent.
				}
			}
		case <-s.stop:
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stop)
}
