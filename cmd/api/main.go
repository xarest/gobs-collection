package main

import (
	"context"
	"fmt"
	"os"

	"syscall"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/api"
	"github.com/xarest/gobs-template/lib/logger"
	"github.com/xarest/gobs-template/worker/local"
)

type keyType string

const ENV_KEY keyType = "mode"

func main() {
	fmt.Println("Starting API server")
	ctx := context.Background()

	app_mode := os.Getenv("APP_MODE")
	if app_mode == "" {
		app_mode = "production"
	}

	appCtx := context.WithValue(ctx, ENV_KEY, app_mode)

	log := logger.Logrus{}
	if err := log.Setup(appCtx, nil); err != nil {
		panic(err)
	}

	bs := gobs.NewBootstrap(gobs.Config{
		NumOfConcurrencies: gobs.DEFAULT_MAX_CONCURRENT,
		Logger:             log.Infof,
		// EnableLogDetail:    true,
	})

	bs.AddOrPanic(&api.API{})
	bs.AddOrPanic(&local.Scheduler{})

	bs.StartBootstrap(appCtx, syscall.SIGINT, syscall.SIGTERM)
}
