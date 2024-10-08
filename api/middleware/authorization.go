package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/lib/cache"
	"github.com/xarest/gobs-collection/lib/db"
	"github.com/xarest/gobs-collection/lib/logger"
	"github.com/xarest/gobs-collection/schema"
)

type Authorization struct {
	log   logger.ILogger
	db    *db.DB
	cache cache.ICache
}

func (a *Authorization) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
			&db.DB{},
			cache.NewICache(),
		},
	}, nil
}

func (a *Authorization) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&a.log, &a.db, &a.cache)
}

func (a *Authorization) VerifyPermission(permissions ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userCtx, err := a.getuserCtx(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized").WithInternal(err)
			}

			if len(permissions) > 0 {
				mPermissions := make(map[string]*schema.Permission)
				for _, r := range userCtx.Roles {
					for _, p := range r.Permissions {
						mPermissions[p.Name] = p
					}
				}
				for _, p := range permissions {
					if _, ok := mPermissions[p]; !ok {
						return echo.NewHTTPError(http.StatusForbidden, "forbidden")
					}
				}
			}
			return next(c)
		}
	}
}

func (a *Authorization) getuserCtx(c echo.Context) (schema.User, error) {
	userCtx, ok := c.Get("user_context").(schema.User)
	if ok {
		return userCtx, nil
	}

	ctx := c.Request().Context()
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return userCtx, fmt.Errorf("token not found")
	}

	uCtx, err := a.cache.Get(ctx, token.Raw, &userCtx)
	if err != nil {
		return userCtx, err
	}
	puCtx, ok := uCtx.(*schema.User)
	if !ok {
		return userCtx, fmt.Errorf("invalid user context")
	}
	userCtx = *puCtx
	c.Set("user_context", userCtx)

	return userCtx, nil
}

var _ gobs.IServiceInit = (*Authorization)(nil)
var _ gobs.IServiceSetup = (*Authorization)(nil)
