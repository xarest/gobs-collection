package gocronwork

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/lib/logger"
	"github.com/xarest/gobs-collection/schema"
	"github.com/xarest/gobs-collection/worker/pool"
)

type Scheduler struct {
	log       logger.ILogger
	scheduler gocron.Scheduler
	workers   []pool.IWorker
}

func (g *Scheduler) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
			&pool.Worker1{},
			&pool.Worker2{},
		},
	}, nil
}

func (g *Scheduler) Setup(ctx context.Context, deps ...gobs.IService) error {
	if err := gobs.Dependencies(deps).Assign(&g.log); err != nil {
		return err
	}
	for _, dep := range deps {
		if w, ok := dep.(pool.IWorker); ok {
			g.workers = append(g.workers, w)
		}
	}
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}
	g.scheduler = s
	return nil
}

func (g *Scheduler) Start(ctx context.Context) error {
	g.scheduler.Start()
	return nil
}

func (g *Scheduler) Stop(ctx context.Context) error {
	return g.scheduler.Shutdown()
}

func (s *Scheduler) AddTask(task *schema.Task) error {
	s.log.Debug("Adding task to scheduler")
	for _, w := range s.workers {
		if w.ID() == task.WorkerID {
			j, err := s.scheduler.NewJob(
				gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Now().Add(5*time.Second))),
				gocron.NewTask(w.Execute, context.TODO(), task.Params),
				gocron.WithEventListeners(
					gocron.AfterJobRuns(
						func(jobID uuid.UUID, jobName string) {
							task.Status = schema.TaskStatusDone
						},
					),
					gocron.AfterJobRunsWithError(
						func(jobID uuid.UUID, jobName string, err error) {
							task.Status = schema.TaskStatusFailed
							task.Error = err.Error()
						},
					),
					gocron.BeforeJobRuns(
						func(jobID uuid.UUID, jobName string) {
							task.Status = schema.TaskStatusRunning
						},
					),
				),
			)
			if err != nil {
				return err
			}
			task.ID = j.ID()
			return nil
		}
	}
	return fmt.Errorf("worker not found for ID %s", task.WorkerID)
}

var _ gobs.IServiceInit = (*Scheduler)(nil)
var _ gobs.IServiceSetup = (*Scheduler)(nil)
var _ gobs.IServiceStart = (*Scheduler)(nil)
