package logger

import (
	"context"

	"github.com/xarest/gobs"
	"github.com/xarest/gobs-template/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	*zap.SugaredLogger
}

var _ gobs.IServiceSetup = (*ZapLogger)(nil)
var _ gobs.IServiceStart = (*ZapLogger)(nil)
var _ gobs.IServiceStop = (*ZapLogger)(nil)

func (l *ZapLogger) Setup(c context.Context, _ gobs.Dependencies) error {
	var config zap.Config
	env := utils.GetAppMode(c)
	if env == "development" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Level.SetLevel(zapcore.DebugLevel)
	} else {
		config = zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.WarnLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}
	l.SugaredLogger = logger.Sugar()
	return nil
}

func (l *ZapLogger) Start(c context.Context) error {
	l.Infof("Start logger at mode %s ", utils.GetAppMode(c))
	return nil
}

func (l *ZapLogger) Stop(c context.Context) error {
	l.Info("End of logger. Flush all logs")
	l.Sync()
	return nil
}
