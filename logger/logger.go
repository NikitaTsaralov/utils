package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Instance *zap.SugaredLogger
var Level zap.AtomicLevel

const defaultLevel = zap.InfoLevel

func init() {
	var level zapcore.Level
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = defaultLevel
	}

	Level = zap.NewAtomicLevelAt(level)

	config := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			zapcore.RFC3339NanoTimeEncoder(time.UTC(), encoder)
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), zapcore.AddSync(os.Stderr), Level)
	Instance = zap.New(core).
		Sugar().
		Named("fd")

	saramaConfig := config
	// omit level key for sarama logger
	saramaConfig.LevelKey = zapcore.OmitKey
	saramaCore := zapcore.NewCore(zapcore.NewJSONEncoder(saramaConfig), zapcore.AddSync(os.Stderr), Level)
	saramaLogger := zap.New(saramaCore).Named("sarama")
	sarama.Logger, _ = zap.NewStdLogAt(saramaLogger, zapcore.InfoLevel)
	sarama.DebugLogger, _ = zap.NewStdLogAt(saramaLogger, zapcore.DebugLevel)

	Instance.Infof("Logger initialized with level: %s", level)
}

func Logger() *zap.SugaredLogger {
	return Instance
}

func Debug(ctx context.Context, args ...any) {
	Instance.With(ctx).Debug(args...)
}

func Info(ctx context.Context, args ...any) {
	Instance.With(ctx).Info(args...)
}

func Warn(ctx context.Context, args ...any) {
	Instance.With(ctx).Warn(args...)
}

func Error(ctx context.Context, args ...any) {
	Instance.With(ctx).Error(args...)
}

func Panic(ctx context.Context, args ...any) {
	Instance.With(ctx).Panic(args...)
}

func Fatal(ctx context.Context, args ...any) {
	Instance.With(ctx).Fatal(args...)
}

func Debugf(ctx context.Context, template string, args ...any) {
	Instance.With(ctx).Debugf(template, args...)
}

func Infof(ctx context.Context, template string, args ...any) {
	Instance.With(ctx).Infof(template, args...)
}

func Warnf(ctx context.Context, template string, args ...any) {
	Instance.With(ctx).With(ctx).Warnf(template, args...)
}

func Errorf(ctx context.Context, template string, args ...any) {
	Instance.With(ctx).With(ctx).Errorf(template, args...)
}

func Panicf(ctx context.Context, template string, args ...any) {
	Instance.With(ctx).Panicf(template, args...)
}

func Fatalf(ctx context.Context, template string, args ...any) {
	Instance.With(ctx).Fatalf(template, args...)
}
