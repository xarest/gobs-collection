package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xarest/gobs"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	eMiddleware "github.com/labstack/echo/v4/middleware"
	admin_api "github.com/xarest/gobs-collection/api/handler/admin"
	hCommon "github.com/xarest/gobs-collection/api/handler/common"
	public_api "github.com/xarest/gobs-collection/api/handler/public"
	superadmin_api "github.com/xarest/gobs-collection/api/handler/super-admin"
	user_api "github.com/xarest/gobs-collection/api/handler/user"
	"github.com/xarest/gobs-collection/api/middleware"
	"github.com/xarest/gobs-collection/api/validator"

	"github.com/xarest/gobs-collection/lib/config"
	httpclient "github.com/xarest/gobs-collection/lib/http-client"
	"github.com/xarest/gobs-collection/lib/logger"
	gCommon "github.com/xarest/gobs/common"
)

type APIConfig struct {
	AllowOrigins string `env:"ALLOW_ORIGINS" mapstructure:"ALLOW_ORIGINS" envDefault:"*"`
	AllowHeaders string `env:"ALLOW_HEADERS" mapstructure:"ALLOW_HEADERS" envDefault:"*"`

	IdleTimeout      int `env:"IDLE_TIMEOUT" mapstructure:"IDLE_TIMEOUT" envDefault:"10"`
	MaxConcurrent    int `env:"MAX_CONCURRENT" mapstructure:"MAX_CONCURRENT" envDefault:"1000"`
	MaxReadFrameSize int `env:"MAX_READ_FRAME_SIZE" mapstructure:"MAX_READ_FRAME_SIZE" envDefault:"1048576"`
	Port             int `env:"PORT" mapstructure:"PORT" envDefault:"8080"`
}

type API struct {
	config     APIConfig
	log        logger.ILogger
	httpServer *http.Server
	httpClient *httpclient.HTTPCient
}

func (a *API) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{
			logger.NewILogger(),
			config.NewIConfig(),
			&httpclient.HTTPCient{},

			&middleware.HTTPErrorHandling{},
			&middleware.MWLogger{},
			&validator.Validator{},

			&public_api.PublicHandler{},
			&superadmin_api.SuperAdminHandler{},
			&user_api.UserHandler{},
			&admin_api.AdminHandler{},
		},
		AsyncMode: map[gCommon.ServiceStatus]bool{
			gCommon.StatusStart: true,
		},
		OnInterrupt: func(errno int) {
			// handle interrupt signal before shutting down API server
			if a.log != nil {
				a.log.Warnf("API server got interrupt signal (%d)", errno)
			}
		},
	}, nil
}

func (a *API) Setup(ctx context.Context, deps ...gobs.IService) error {
	// gobs parse all dependencies
	var (
		cfgService    config.IConfiguration
		mErrorHandler *middleware.HTTPErrorHandling
		validator     *validator.Validator
		mLogger       *middleware.MWLogger
		handlers      []hCommon.IHandler
	)
	if err := gobs.Dependencies(deps).Assign(
		&a.log,
		&cfgService,
		&a.httpClient,

		&mErrorHandler,
		&mLogger,
		&validator,
	); err != nil {
		return err
	}
	for _, d := range deps {
		if h, ok := d.(hCommon.IHandler); ok {
			handlers = append(handlers, h)
		}
	}

	// parse API server config
	if err := cfgService.Parse(&a.config); err != nil {
		return err
	}

	// setup echo engine
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Pre(eMiddleware.CORSWithConfig(eMiddleware.CORSConfig{
		AllowOrigins: strings.Split(a.config.AllowOrigins, ","),
		AllowHeaders: strings.Split(a.config.AllowHeaders, ","),
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
	}))

	e.Use(eMiddleware.Recover())
	e.Use(eMiddleware.RequestID())
	e.Use(eMiddleware.Gzip())
	e.Use(eMiddleware.Decompress())
	e.Use(eMiddleware.CSRF())

	e.HTTPErrorHandler = mErrorHandler.CatchErr
	e.Validator = validator
	e.Use(mLogger.Handler())

	// setup routes
	g := e.Group("/api")
	for _, h := range handlers {
		h.Route(g)
	}

	// setup http server from echo engine
	h2s := &http2.Server{
		MaxConcurrentStreams: uint32(a.config.MaxConcurrent),
		MaxReadFrameSize:     uint32(a.config.MaxReadFrameSize),
		IdleTimeout:          time.Duration(a.config.IdleTimeout) * time.Second,
	}
	a.httpServer = &http.Server{
		Addr:    ":" + strconv.Itoa(a.config.Port),
		Handler: h2c.NewHandler(e, h2s),
	}
	return nil
}

func (a *API) StartServer(ctx context.Context, onReady func(err error)) error {
	a.log.Infof("API server is running on port %d", a.config.Port)
	go func(ctx context.Context) {
		defer onReady(nil)
		for {
			resp, err := a.httpClient.Get(ctx, fmt.Sprintf("http://localhost:%d/api/ping", a.config.Port), nil)
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			a.log.Infof("API server is ready: %s", string(resp))
			break
		}
	}(ctx)
	return a.httpServer.ListenAndServe()
}

func (a *API) Stop(ctx context.Context) error {
	a.log.Info("API server is shutting down")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return a.httpServer.Shutdown(ctx)
}

var _ gobs.IServiceInit = (*API)(nil)
var _ gobs.IServiceSetup = (*API)(nil)
var _ gobs.IServiceStartServer = (*API)(nil)
var _ gobs.IServiceStop = (*API)(nil)
