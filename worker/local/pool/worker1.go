package pool

import (
	"context"
	"encoding/json"
	"time"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/lib/logger"
)

type Worker1 struct {
	log logger.ILogger
}

type ParamsModel struct {
	Delay time.Duration `json:"delay"`
}

func (w *Worker1) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
		},
	}, nil
}

func (w *Worker1) Setup(ctx context.Context, deps gobs.Dependencies) error {
	return deps.Assign(&w.log)
}

func (w *Worker1) Execute(ctx context.Context, jsParams []byte) (any, error) {
	var params ParamsModel
	if err := json.Unmarshal(jsParams, &params); err != nil {
		return nil, err
	}
	w.log.Info("Worker1.Execute start")
	time.Sleep(params.Delay * time.Millisecond)
	w.log.Info("Worker1.Execute end")
	return nil, nil
}

func (w *Worker1) ID() string {
	return "worker1"
}

var _ IWorker = (*Worker1)(nil)
var _ gobs.IServiceInit = (*Worker1)(nil)
var _ gobs.IServiceSetup = (*Worker1)(nil)
