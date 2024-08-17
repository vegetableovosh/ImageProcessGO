package types

import "github.com/google/uuid"

type TaskStatus string

const (
	StatusInProgress TaskStatus = "in_progress"
	StatusReady      TaskStatus = "ready"
)

type Task struct {
	ID     string     `json:"id"`
	Status TaskStatus `json:"status"`
	Result string     `json:"result"`
}

type TaskStorage interface {
	PostTask(task *Task)
	GetTask(id string) (*Task, bool)
}

func NewTask() *Task {
	return &Task{
		ID:     uuid.New().String(),
		Status: StatusInProgress,
	}
}
