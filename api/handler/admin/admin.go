package admin

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/api/handler/common"
)

type AdminHandler struct {
	handlers []common.IHandler
}

var _ gobs.IServiceInit = (*AdminHandler)(nil)

func (a *AdminHandler) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			// Add handlers for admin here
		},
	}, nil
}

var _ gobs.IServiceSetup = (*AdminHandler)(nil)

func (a *AdminHandler) Setup(ctx context.Context, deps ...gobs.IService) error {
	for _, d := range deps {
		a.handlers = append(a.handlers, d.(common.IHandler))
	}
	return nil
}

var _ common.IHandler = (*AdminHandler)(nil)

func (a *AdminHandler) Route(r *echo.Group) {
	for _, h := range a.handlers {
		h.Route(r)
	}
}
