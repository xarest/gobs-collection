package public

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/api/handler/common"
)

type PublicHandler struct {
	handlers []common.IHandler
}

var _ gobs.IServiceInit = (*PublicHandler)(nil)

func (a *PublicHandler) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			// Add handlers for public here
			&HealthCheck{},
			&Auth{},
		},
	}, nil
}

var _ gobs.IServiceSetup = (*PublicHandler)(nil)

func (a *PublicHandler) Setup(ctx context.Context, deps gobs.Dependencies) error {
	for _, d := range deps {
		a.handlers = append(a.handlers, d.(common.IHandler))
	}
	return nil
}

var _ common.IHandler = (*PublicHandler)(nil)

func (a *PublicHandler) Route(r *echo.Group) {
	for _, h := range a.handlers {
		h.Route(r)
	}
}
