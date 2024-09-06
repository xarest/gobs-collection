package user

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	hCommon "github.com/xarest/gobs-collection/api/handler/common"
	"github.com/xarest/gobs-collection/api/middleware"
	"github.com/xarest/gobs-collection/lib/logger"
	"github.com/xarest/gobs-collection/service/auth"
)

type Auth struct {
	log     logger.ILogger
	mwAuth  *middleware.Authentication
	service *auth.Auth
}

// Init implements gobs.IServiceInit.
func (a *Auth) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
			&middleware.Authentication{},
			&auth.Auth{},
		},
	}, nil
}

// Setup implements gobs.IServiceSetup.
func (a *Auth) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&a.log, &a.mwAuth, &a.service)
}

// Route implements common.IHandler.
func (a *Auth) Route(r *echo.Group) {
	g := r.Group("/auth")
	g.Use(a.mwAuth.CheckJWTToken())
	g.PATCH("", a.RefreshToken)
	g.DELETE("", a.SignOut)
}

func (a *Auth) RefreshToken(c echo.Context) error {
	tokenStr, err := hCommon.ExtractToken(c.Request().Header["Authorization"][0])
	if nil != err {
		return err
	}

	resp, err := a.service.RefreshToken(c.Request().Context(), tokenStr)
	if nil != err {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *Auth) SignOut(c echo.Context) error {
	return c.JSON(http.StatusOK, "sign out")
}

var _ gobs.IServiceInit = (*Auth)(nil)
var _ gobs.IServiceSetup = (*Auth)(nil)
var _ hCommon.IHandler = (*Auth)(nil)
