package background

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/api/handler/common"
	"github.com/xarest/gobs-template/api/validator"
	"github.com/xarest/gobs-template/schema"
	"github.com/xarest/gobs-template/worker"
)

type Check struct {
	worker worker.IClient
}

func (b *Check) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			worker.NewIClient(),
		},
	}, nil
}

func (b *Check) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&b.worker)
}

// Route implements common.IHandler.
func (b *Check) Route(r *echo.Group) {
	g := r.Group("/check")
	g.GET("", b.Overall)
	g.GET("/:id", b.GetTaskDetail)
}

func (b *Check) Overall(c echo.Context) error {
	var page schema.Page
	validator.BindAndValidate(c, &page)
	page.LoadDefault()

	taskRuns, err := b.worker.GetTasks(schema.TaskStatusDone, page)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, taskRuns)
}

func (b *Check) GetTaskDetail(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing task ID")
	}
	taskID, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID").WithInternal(err)
	}
	task, err := b.worker.GetTask(taskID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, task)
}

var _ gobs.IServiceInit = (*Check)(nil)
var _ gobs.IServiceSetup = (*Check)(nil)
var _ common.IHandler = (*Check)(nil)
