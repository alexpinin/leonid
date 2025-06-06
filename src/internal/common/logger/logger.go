package logger

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"leonid.log",
	}
	var err error
	logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}

func Info(msg string) {
	logger.WithOptions(zap.AddStacktrace(zap.FatalLevel)).Info(msg)
}

func Error(msg string) {
	logger.WithOptions(zap.AddStacktrace(zap.FatalLevel)).Error(msg)
}

func Panic(msg string) {
	logger.WithOptions(zap.AddStacktrace(zap.FatalLevel)).Panic(msg)
}
