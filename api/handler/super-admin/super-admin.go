package superadmin

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/api/handler/common"
)

type SuperAdminHandler struct {
	handlers []common.IHandler
}

var _ gobs.IServiceInit = (*SuperAdminHandler)(nil)

func (a *SuperAdminHandler) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			// Add handlers for super admin here
		},
	}, nil
}

var _ gobs.IServiceSetup = (*SuperAdminHandler)(nil)

func (a *SuperAdminHandler) Setup(ctx context.Context, deps ...gobs.IService) error {
	for _, d := range deps {
		a.handlers = append(a.handlers, d.(common.IHandler))
	}
	return nil
}

var _ common.IHandler = (*SuperAdminHandler)(nil)

func (a *SuperAdminHandler) Route(r *echo.Group) {
	for _, h := range a.handlers {
		h.Route(r)
	}
}
