package gocronwork

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/schema"
)

type Cient struct {
	scheduler *Scheduler
	taskList  []*schema.Task
}

func (w *Cient) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			&Scheduler{},
		},
	}, nil
}

func (w *Cient) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&w.scheduler)
}

func (w *Cient) AddTask(workerID string, params any, createdBy uuid.UUID) error {
	jsParams, err := json.Marshal(params)
	if err != nil {
		return err
	}
	task := schema.Task{
		ID:        uuid.New(),
		WorkerID:  workerID,
		Params:    jsParams,
		Status:    schema.TaskStatusPending,
		CreatedAt: time.Now(),
	}
	w.taskList = append(w.taskList, &task)
	fmt.Println("Adding task to scheduler", w.taskList)
	return w.scheduler.AddTask(&task)
}

func (w *Cient) GetTask(id uuid.UUID) (schema.Task, error) {
	for _, t := range w.taskList {
		if t.ID == id {
			return *t, nil
		}
	}
	return schema.Task{}, fmt.Errorf("task not found")
}

func (w *Cient) GetTasks(status schema.TaskStatus, page schema.Page) ([]schema.Task, error) {
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

var _ gobs.IServiceInit = (*Cient)(nil)
var _ gobs.IServiceSetup = (*Cient)(nil)
