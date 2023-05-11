package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Interface = &Logger{}

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	logger *zap.Logger
}

func (l *Logger) Debug(message interface{}, args ...interface{}) {

}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Sugar().Infof("%s: %v", message, args)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Sugar().Warnf("%s: %v", message, args)
}

func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.logger.Sugar().Errorf("%s: %v", message, args)
}

func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.logger.Sugar().Fatalf("%s: %v", message, args)
}

func New(level string) *Logger {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}

	cfg := zap.Config{
		Level:       lvl,
		Development: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	op := zap.AddCallerSkip(1)
	l, err := cfg.Build(op)
	if err != nil {
		panic(err)
	}

	return &Logger{
		logger: l,
	}
}
