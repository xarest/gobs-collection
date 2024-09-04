package middleware

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/lib/logger"
)

type HTTPErrorHandling struct {
	log logger.ILogger
}

func (r *HTTPErrorHandling) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
		},
	}, nil
}

func (r *HTTPErrorHandling) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&r.log)
}

func (r *HTTPErrorHandling) CatchErr(err error, c echo.Context) {
	if errors.Is(err, net.ErrClosed) {
		r.log.Debug("Connection is closed by client")
		return
	}

	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		httpErr = echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	if httpErr.Code == http.StatusInternalServerError {
		r.log.Errorf("Error: %v", err)
	}
	if err := c.JSON(httpErr.Code, httpErr.Message); err != nil {
		r.log.Errorf("Error: %v", err)
	}
}

var _ gobs.IServiceInit = (*HTTPErrorHandling)(nil)
var _ gobs.IServiceSetup = (*HTTPErrorHandling)(nil)
