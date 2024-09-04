package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/lib/logger"
	"golang.org/x/time/rate"
)

type RateLimit struct {
	log logger.ILogger
}

// Init implements gobs.IService.
func (r *RateLimit) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
		},
	}, nil
}

func (r *RateLimit) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&r.log)
}

func (r *RateLimit) Handler(limit float64, burst int, duration time.Duration) echo.MiddlewareFunc {
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(limit), Burst: burst, ExpiresIn: duration},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	return middleware.RateLimiterWithConfig(config)
}

var _ gobs.IServiceInit = (*RateLimit)(nil)
var _ gobs.IServiceSetup = (*RateLimit)(nil)
