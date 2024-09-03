package common

import "github.com/labstack/echo/v4"

type IHandler interface {
	Route(r *echo.Group)
}
