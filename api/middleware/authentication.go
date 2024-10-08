package middleware

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/lib/cache"
	"github.com/xarest/gobs-collection/lib/config"
	"github.com/xarest/gobs-collection/lib/logger"
)

type Authentication struct {
	log    logger.ILogger
	cache  cache.ICache
	config *JWTSecret
}

func (a *Authentication) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: []gobs.IService{
			logger.NewILogger(),
			config.NewIConfig(),
			cache.NewICache(),
		},
	}, nil
}

func (a *Authentication) Setup(ctx context.Context, deps ...gobs.IService) error {
	var (
		config config.IConfiguration
		cfg    JWTSecret
	)
	if err := gobs.Dependencies(deps).Assign(&a.log, &config, &a.cache); err != nil {
		return err
	}
	if err := config.Parse(&cfg); err != nil {
		return err
	}
	a.config = &cfg
	return nil
}

type JWTSecret struct {
	Secret string `env:"JWT_SECRET" mapstructure:"JWT_SECRET" envDefault:"mysecretjwt"`
}

func (a *Authentication) CheckJWTToken() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
		SigningKey:     []byte(a.config.Secret),
		SuccessHandler: func(c echo.Context) {},
	}
	return echojwt.WithConfig(config)
}

func (a *Authentication) CheckRedisToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

var _ gobs.IServiceInit = (*Authentication)(nil)
