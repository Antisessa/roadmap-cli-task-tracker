package task

import (
	"encoding/json"
	"time"
)

type Status int

const (
	Todo = iota
	InProgress
	Done
)

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (t *Task) UpdateDescription(newDescription string) {
	t.Description = newDescription
	t.UpdatedAt = time.Now()
}

func (t *Task) UpdateStatus(newStatus Status) {
	t.Status = newStatus
	t.UpdatedAt = time.Now()
}

func (t *Task) ToJson() ([]byte, error) {
	return json.Marshal(t)
}
