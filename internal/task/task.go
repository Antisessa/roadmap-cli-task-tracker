package task

import (
	"encoding/json"
	"fmt"
	"time"
)

type Status string

const (
	None       = `none`
	Todo       = `todo`
	InProgress = `in-progress`
	Done       = `done`
)

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (t *Task) UpdateDescription(newDescription string) error {
	if t.Description == newDescription {
		return fmt.Errorf("error! description already set")
	}
	t.Description = newDescription
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) UpdateStatus(newStatus Status) error {
	if t.Status == newStatus {
		return fmt.Errorf("error! status already set")
	}
	t.Status = newStatus
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) ToJson() ([]byte, error) {
	return json.MarshalIndent(t, "", "")
}

func (t *Task) String() string {
	return fmt.Sprintf(
		"ID: %-3d | Description: %-10s | Status: %-12s | CreatedAt: %s | UpdatedAt: %s",
		t.Id,
		t.Description,
		t.Status,
		t.CreatedAt.Format("02.01.2006 15:04:05"),
		t.UpdatedAt.Format("02.01.2006 15:04:05"),
	)
}
