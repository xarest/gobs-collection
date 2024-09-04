package public

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/api/dto"
	hCommon "github.com/xarest/gobs-template/api/handler/common"
	"github.com/xarest/gobs-template/api/validator"
	"github.com/xarest/gobs-template/lib/logger"
	"github.com/xarest/gobs-template/service/auth"
)

type Auth struct {
	log     logger.ILogger
	service *auth.Auth
}

func (a *Auth) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
			&auth.Auth{},
		},
	}, nil
}

func (a *Auth) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&a.log, &a.service)
}

func (a *Auth) Route(r *echo.Group) {
	g := r.Group("/auth")
	g.POST("", a.SignUp)
	g.PUT("", a.SignIn)
}

func (a *Auth) SignUp(c echo.Context) error {
	var creds dto.Credentials
	if err := validator.BindAndValidate(c, &creds); nil != err {
		a.log.Error(err)
		return err
	}

	resp, err := a.service.Register(c.Request().Context(), creds)
	if nil != err {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *Auth) SignIn(c echo.Context) error {
	var creds dto.Credentials
	if err := validator.BindAndValidate(c, &creds); nil != err {
		a.log.Error(err)
		return err
	}

	resp, err := a.service.SignIn(c.Request().Context(), creds)
	if nil != err {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

var _ gobs.IServiceInit = (*Auth)(nil)
var _ gobs.IServiceSetup = (*Auth)(nil)
var _ hCommon.IHandler = (*Auth)(nil)
