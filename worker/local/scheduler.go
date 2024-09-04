package local

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/lib/config"
	"github.com/xarest/gobs-template/lib/logger"
	"github.com/xarest/gobs-template/schema"
	"github.com/xarest/gobs-template/worker/local/pool"
)

type SchedulerConfig struct {
	ExecuteInterval int `env:"SCHEDULER_EXECUTE_INTERVAL" envDefault:"10000"`
}

type SchedulerStatus int

const (
	SchedulerStatuInit SchedulerStatus = iota
	SchedulerStatusRunning
	SchedulerStatusWaiting
)

type Scheduler struct {
	log             logger.ILogger
	ExecuteInterval int
	workers         []pool.IWorker
	tasks           []*schema.Task
	mu              *sync.Mutex
	status          SchedulerStatus
	ch              chan any
}

func (s *Scheduler) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
			config.NewIConfig(),
			&pool.Worker1{},
			&pool.Worker2{},
		},
	}, nil
}

func (s *Scheduler) Setup(ctx context.Context, deps ...gobs.IService) error {
	var sConfig config.IConfiguration
	if err := gobs.Dependencies(deps).Assign(&s.log, &sConfig); err != nil {
		return err
	}
	var cfg SchedulerConfig
	if err := sConfig.Parse(&cfg); err != nil {
		return err
	}
	s.ExecuteInterval = cfg.ExecuteInterval

	for _, d := range deps {
		if w, ok := d.(pool.IWorker); ok {
			s.workers = append(s.workers, w)
		}
	}
	s.mu = &sync.Mutex{}
	s.ch = make(chan any, 5)
	s.tasks = []*schema.Task{}
	return nil
}

func (s *Scheduler) StartServer(ctx context.Context, onReady func(err error)) error {
	s.mu.Lock()
	if s.status == SchedulerStatuInit {
		s.status = SchedulerStatusWaiting
	}
	s.mu.Unlock()

	onReady(nil)
	for {
		timer, cancel := context.WithTimeout(ctx, time.Duration(s.ExecuteInterval)*time.Millisecond)
		select {
		case <-timer.Done():
			s.log.Debug("Scheduler run periodically")
		case _, ok := <-s.ch:
			if !ok {
				cancel()
				return nil
			}
			s.log.Debug("Scheduler run for task incomming")
		}
		cancel()
		if len(s.tasks) == 0 {
			continue
		}
		s.mu.Lock()
		s.status = SchedulerStatusRunning
		tasks := s.tasks
		s.tasks = nil
		s.mu.Unlock()

		for _, task := range tasks {
			for _, w := range s.workers {
				if w.ID() != task.WorkerID {
					continue
				}
				res, err := w.Execute(ctx, task.Params)
				if err != nil {
					task.Error = err.Error()
					task.Status = schema.TaskStatusFailed
				} else {
					task.Status = schema.TaskStatusDone
				}
				if res != nil {
					task.Result, err = json.Marshal(res)
					if err != nil {
						s.log.Error("Failed to marshal task result: %v", err)
					}
				}
				// Save task results to DB
				out, _ := json.Marshal(task)
				s.log.Info("Finish task", task.ID, string(out))
			}
		}

		s.mu.Lock()
		s.status = SchedulerStatusWaiting
		s.mu.Unlock()
	}
}

func (s *Scheduler) AddTask(task *schema.Task) {
	s.mu.Lock()
	s.tasks = append(s.tasks, task)
	if s.status == SchedulerStatusWaiting {
		s.ch <- nil
	}
	s.mu.Unlock()
}

func (s *Scheduler) Stop(ctx context.Context) error {
	close(s.ch)
	return nil
}

var _ gobs.IServiceInit = (*Scheduler)(nil)
var _ gobs.IServiceSetup = (*Scheduler)(nil)
var _ gobs.IServiceStartServer = (*Scheduler)(nil)
