package background

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/api/handler/common"
)

type BackgroundHandler struct {
	handlers []common.IHandler
}

func (a *BackgroundHandler) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			// Add handlers for Background here
			&Check{},
			&Trigger{},
		},
	}, nil
}

func (a *BackgroundHandler) Setup(ctx context.Context, deps ...gobs.IService) error {
	for _, d := range deps {
		a.handlers = append(a.handlers, d.(common.IHandler))
	}
	return nil
}

func (a *BackgroundHandler) Route(r *echo.Group) {
	g := r.Group("/background")
	for _, h := range a.handlers {
		h.Route(g)
	}
}

var _ gobs.IServiceInit = (*BackgroundHandler)(nil)
var _ gobs.IServiceSetup = (*BackgroundHandler)(nil)
var _ common.IHandler = (*BackgroundHandler)(nil)
