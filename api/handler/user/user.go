package user

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/api/handler/common"
	"github.com/xarest/gobs-template/api/handler/user/background"
)

type UserHandler struct {
	handlers []common.IHandler
}

func (a *UserHandler) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			// Add handlers for User here
			&Auth{},
			&background.BackgroundHandler{},
		},
	}, nil
}

func (a *UserHandler) Setup(ctx context.Context, deps ...gobs.IService) error {
	for _, d := range deps {
		a.handlers = append(a.handlers, d.(common.IHandler))
	}
	return nil
}

func (a *UserHandler) Route(r *echo.Group) {
	for _, h := range a.handlers {
		h.Route(r)
	}
}

var _ gobs.IServiceInit = (*UserHandler)(nil)
var _ gobs.IServiceSetup = (*UserHandler)(nil)
var _ common.IHandler = (*UserHandler)(nil)
