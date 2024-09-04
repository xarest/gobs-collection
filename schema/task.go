package schema

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TaskStatus int

const (
	TaskStatusInit TaskStatus = iota
	TaskStatusRunning
	TaskStatusPending
	TaskStatusDone
	TaskStatusFailed
)

type Task struct {
	ID        uuid.UUID       `json:"id"`
	WorkerID  string          `json:"worker_id"`
	Status    TaskStatus      `json:"status"`
	Params    json.RawMessage `json:"params"`
	Result    json.RawMessage `json:"result"`
	Error     string          `json:"error"`
	CreatedAt time.Time       `json:"created_at"`
	CreatedBy uuid.UUID       `json:"created_by"`
}
