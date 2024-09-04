package logger

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/utils"
)

type Logrus struct {
	*log.Logger
}

// Setup implements gobs.IServiceSetup.
func (l *Logrus) Setup(ctx context.Context, _ ...gobs.IService) error {
	l.Logger = log.New()
	appMode := utils.GetAppMode(ctx)
	if appMode == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

func (l *Logrus) Start(c context.Context) error {
	l.Infof("Start logger at mode %s ", utils.GetAppMode(c))
	return nil
}

func (l *Logrus) Stop(c context.Context) error {
	l.Info("End of logger. Flush all logs")
	return nil
}

func (l *Logrus) DPanic(args ...interface{}) {
	l.Logger.Panic(args...)
}

func (l *Logrus) DPanicf(format string, args ...interface{}) {
	l.Logger.Panicf(format, args...)
}

var _ ILogger = (*Logrus)(nil)
var _ gobs.IServiceSetup = (*Logrus)(nil)
var _ gobs.IServiceStart = (*Logrus)(nil)
var _ gobs.IServiceStop = (*Logrus)(nil)
