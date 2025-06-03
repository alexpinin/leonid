package logger

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
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
