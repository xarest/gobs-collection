package logger

type ILogger interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	DPanicf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Info(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	DPanic(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
}

func NewILogger() ILogger {
	return &ZapLogger{}
}
