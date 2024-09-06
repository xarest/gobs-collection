package user

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/api/handler/common"
	"github.com/xarest/gobs-collection/api/handler/user/background"
	"github.com/xarest/gobs-collection/api/middleware"
)

type UserHandler struct {
	authen   *middleware.Authentication
	author   *middleware.Authorization
	handlers []common.IHandler
}

func (a *UserHandler) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			&middleware.Authentication{},
			&middleware.Authorization{},
			// Add handlers for User here
			&Auth{},
			&background.BackgroundHandler{},
		},
	}, nil
}

func (a *UserHandler) Setup(ctx context.Context, deps ...gobs.IService) error {
	if err := gobs.Dependencies(deps).Assign(&a.authen, &a.author); err != nil {
		return err
	}
	for _, d := range deps {
		if h, ok := d.(common.IHandler); ok {
			a.handlers = append(a.handlers, h)
		}
	}
	return nil
}

func (a *UserHandler) Route(r *echo.Group) {
	uGroup := r.Group("/user")
	uGroup.Use(a.authen.CheckJWTToken())
	uGroup.Use(a.author.VerifyPermission("user"))
	for _, h := range a.handlers {
		h.Route(uGroup)
	}
}

var _ gobs.IServiceInit = (*UserHandler)(nil)
var _ gobs.IServiceSetup = (*UserHandler)(nil)
var _ common.IHandler = (*UserHandler)(nil)
