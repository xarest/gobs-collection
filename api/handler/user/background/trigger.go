package background

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/api/handler/common"
	"github.com/xarest/gobs-template/worker"
)

type Trigger struct {
	wc worker.IClient
}

func (b *Trigger) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			worker.NewIClient(),
		},
	}, nil
}

func (b *Trigger) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&b.wc)
}

// Route implements common.IHandler.
func (t *Trigger) Route(r *echo.Group) {
	g := r.Group("/trigger")
	g.GET("", t.TriggerTask)
}

func (t *Trigger) TriggerTask(c echo.Context) error {
	params := map[string]any{
		"delay": 5000,
	}
	if err := t.wc.AddTask("worker1", params, uuid.New()); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Trigger task ID")
}

var _ common.IHandler = (*Trigger)(nil)
var _ gobs.IServiceInit = (*Trigger)(nil)
var _ gobs.IServiceSetup = (*Trigger)(nil)
