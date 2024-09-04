package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/lib/logger"
)

type MWLogger struct {
	log logger.ILogger
}

// Init implements gobs.IService.
func (l *MWLogger) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
		},
	}, nil
}

func (l *MWLogger) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&l.log)
}

func (l *MWLogger) Handler() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, id=${id} status=${status} datetime=${time_rfc3339} latency=${latency_human} ${error}\n",
		})
}

func (l *MWLogger) Log(c echo.Context, v middleware.RequestLoggerValues) error {
	l.log.Debugf("%s %s %d %s",
		v.Method,
		v.URI,
		v.Status,
		v.Latency.String(),
	)
	return nil
}

var _ gobs.IServiceInit = (*MWLogger)(nil)
var _ gobs.IServiceSetup = (*MWLogger)(nil)
