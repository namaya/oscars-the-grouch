package log

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger = New()
var logCtxKey = "logger"

func New() *zap.Logger {
	zcfg := zap.NewProductionConfig()
	zcfg.OutputPaths = []string{"stdout"}
	zcfg.ErrorOutputPaths = []string{"stderr"}
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	ll := os.Getenv("LOG_LEVEL")
	if ll != "" {
		al := zap.NewAtomicLevel()
		err := al.UnmarshalText([]byte(ll))
		if err != nil {
			panic(err)
		}
		zcfg.Level = al
	}

	logger, err := zcfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func Default() *zap.Logger {
	return logger
}

func NewContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, logCtxKey, logger)
}

func With(ctx context.Context, zapLogger *zap.Logger) context.Context {
	clone := *zapLogger
	return NewContext(ctx, &clone)
}

func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	return NewContext(ctx, FromContext(ctx).With(fields...))
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}

	ctxLogger := ctx.Value(logCtxKey)
	if ctxLogger == nil {
		return logger
	}

	return ctxLogger.(*zap.Logger)
}

func Get(ctx context.Context) *zap.SugaredLogger {
	return FromContext(ctx).Sugar()
}
