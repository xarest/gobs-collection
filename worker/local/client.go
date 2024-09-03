package local

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/schema"
)

type WorkerClient struct {
	scheduler *Scheduler
	taskList  []*schema.Task
}

func (w *WorkerClient) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			&Scheduler{},
		},
	}, nil
}

func (w *WorkerClient) Setup(ctx context.Context, deps gobs.Dependencies) error {
	return deps.Assign(&w.scheduler)
}

func (w *WorkerClient) Stop(ctx context.Context) error {
	// TODO: Wait for multiple dependencies for each stages
	return nil
}

// AddTask implements worker.IClient.
func (w *WorkerClient) AddTask(workerID string, params any) error {
	jsParams, err := json.Marshal(params)
	if err != nil {
		return err
	}
	task := schema.Task{
		ID:       uuid.New(),
		WorkerID: workerID,
		Params:   jsParams,
		Status:   schema.TaskStatusPending,
	}
	w.scheduler.AddTask(&task)
	w.taskList = append(w.taskList, &task)
	return nil
}

// GetTask implements worker.IClient.
func (w *WorkerClient) GetTask(id uuid.UUID) (schema.Task, error) {
	for _, t := range w.taskList {
		if t.ID == id {
			return *t, nil
		}
	}
	return schema.Task{}, fmt.Errorf("task not found")
}

// GetTasks implements worker.IClient.
func (w *WorkerClient) GetTasks(status schema.TaskStatus, page schema.Page) ([]schema.Task, error) {
	var results []schema.Task
	for i, t := range w.taskList {
		if i < page.Offset {
			continue
		}
		if t.Status == status {
			results = append(results, *t)
		}
		if len(results) >= page.Limit {
			break
		}
	}
	return results, nil
}

var _ gobs.IServiceInit = (*WorkerClient)(nil)
var _ gobs.IServiceSetup = (*WorkerClient)(nil)
var _ gobs.IServiceStop = (*WorkerClient)(nil)
