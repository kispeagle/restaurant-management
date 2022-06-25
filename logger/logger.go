package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func GetCustomProductionLogger() {
	conf := zap.NewProductionConfig()

	// Below list is default config of Production logger
	// TimeKey:        "ts",
	// LevelKey:       "level",
	// NameKey:        "logger",
	// CallerKey:      "caller",
	// FunctionKey:    zapcore.OmitKey,
	// MessageKey:     "msg",
	// StacktraceKey:  "stacktrace",
	// LineEnding:     zapcore.DefaultLineEnding,
	// EncodeLevel:    zapcore.LowercaseLevelEncoder,
	// EncodeTime:     zapcore.EpochTimeEncoder,
	// EncodeDuration: zapcore.SecondsDurationEncoder,
	// EncodeCaller:   zapcore.ShortCallerEncoder,

	// modify time structure in logger
	conf.EncoderConfig.TimeKey = "timestamp"
	conf.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	// change from logging to console to logging to file
	conf.OutputPaths = []string{"log.txt"}

	// conf.EncoderConfig.EncodeLevel =

	logger, _ := conf.Build()

	Logger = logger.Sugar()
}
