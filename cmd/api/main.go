package main

import (
	"context"
	"fmt"
	"os"

	"syscall"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/api"
	"github.com/xarest/gobs-template/lib/logger"
	gCommon "github.com/xarest/gobs/common"
	gUtils "github.com/xarest/gobs/utils"
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

	var log logger.ILogger = logger.NewILogger()
	if err := log.(gobs.IServiceSetup).Setup(appCtx, nil); err != nil {
		fmt.Println("Error setting up logger")
		return
	}

	bs := gobs.NewBootstrap(gobs.Config{
		NumOfConcurrencies: gobs.DEFAULT_MAX_CONCURRENT,
		Logger:             log.Infof,
		// EnableLogDetail:    true,
	})

	bs.Add(log, gCommon.StatusSetup, gUtils.DefaultServiceName(logger.NewILogger()))
	bs.AddOrPanic(&api.API{})

	bs.StartBootstrap(appCtx, syscall.SIGINT, syscall.SIGTERM)
}
