package worker

import (
	"github.com/google/uuid"
	"github.com/xarest/gobs-template/schema"
	"github.com/xarest/gobs-template/worker/local"
)

type IClient interface {
	GetTasks(status schema.TaskStatus, page schema.Page) ([]schema.Task, error)
	GetTask(id uuid.UUID) (schema.Task, error)
	AddTask(wokerID string, params any, createdBy uuid.UUID) error
}

func NewIClient() IClient {
	return &local.WorkerClient{}
}
