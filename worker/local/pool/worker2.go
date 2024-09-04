package pool

import (
	"context"
	"encoding/json"
	"time"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/lib/logger"
)

type Worker2 struct {
	log logger.ILogger
}

type ParamsModel2 struct {
	Message time.Duration `json:"message"`
}

func (w *Worker2) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
		},
	}, nil
}

func (w *Worker2) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&w.log)
}

func (w *Worker2) Execute(ctx context.Context, jsonParam []byte) (any, error) {
	var params ParamsModel2
	if err := json.Unmarshal(jsonParam, &params); err != nil {
		return nil, err
	}
	w.log.Info("Worker2.Execute start")
	time.Sleep(5 * time.Second)
	w.log.Info(params.Message)
	w.log.Info("Worker2.Execute end")
	return nil, nil
}

func (w *Worker2) ID() string {
	return "worker2"
}

var _ IWorker = (*Worker2)(nil)
var _ gobs.IServiceInit = (*Worker2)(nil)
var _ gobs.IServiceSetup = (*Worker2)(nil)
