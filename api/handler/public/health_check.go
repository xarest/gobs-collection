package public

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs-collection/api/handler/common"
)

type HealthCheck struct{}

var _ common.IHandler = (*HealthCheck)(nil)

func (a *HealthCheck) Route(r *echo.Group) {
	g := r.Group("/ping")
	g.GET("", a.Ping)
}

func (a *HealthCheck) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}
