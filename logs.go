package logs

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
)

type Logger interface {
	Info(msg string, fields ...zapcore.Field)

	Debug(msg string, fields ...zapcore.Field)

	Warn(msg string, fields ...zapcore.Field)

	Error(msg string, fields ...zapcore.Field)

	Fatal(msg string, fields ...zapcore.Field)

	With(fields ...zapcore.Field) Logger

	Ctx(ctx context.Context) Logger
}

type logger struct {
	jsonStacktrace bool // If true, output stacktrace use embedded json style.
	stackLevel     zapcore.Level
	zapLogger      *zap.Logger
}

func (l *logger) Info(msg string, fields ...zapcore.Field) {
	fields = l.withStack(zapcore.InfoLevel, fields...)

	l.zapLogger.Info(msg, fields...)
}

func (l *logger) Debug(msg string, fields ...zapcore.Field) {
	fields = l.withStack(zapcore.DebugLevel, fields...)

	l.zapLogger.Debug(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zapcore.Field) {
	fields = l.withStack(zapcore.WarnLevel, fields...)

	l.zapLogger.Warn(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zapcore.Field) {
	fields = l.withStack(zapcore.ErrorLevel, fields...)

	l.zapLogger.Error(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...zapcore.Field) {
	fields = l.withStack(zapcore.FatalLevel, fields...)

	l.zapLogger.Fatal(msg, fields...)
}

func (l *logger) With(fields ...zapcore.Field) Logger {
	ll := l.clone()
	ll.zapLogger = ll.zapLogger.With(fields...)

	return ll
}

func (l *logger) Ctx(ctx context.Context) Logger {
	if ctx != nil {
		var fields []zap.Field
		if span := opentracing.SpanFromContext(ctx); span != nil {
			if sc, ok := span.Context().(jaeger.SpanContext); ok {
				fields = append(fields, zap.String("traceId", sc.TraceID().String()))
			}
		}

		var (
			md metadata.MD
			ok bool
		)

		md, ok = metadata.FromIncomingContext(ctx)
		if !ok {
			md, ok = metadata.FromOutgoingContext(ctx)
		}

		for _, key := range fieldList {
			if value := ctx.Value(key); value != nil {
				fields = append(fields, zap.Any(key, value))
			}

			if ok && md != nil {
				if s := md.Get(key); len(s) > 0 {
					fields = append(fields, zap.Any(key, s[0]))
				}
			}
		}

		return l.With(fields...)
	}

	return l
}

func (l *logger) clone() *logger {
	ll := *l
	return &ll
}

func (l *logger) withStack(level zapcore.Level, fields ...zapcore.Field) []zapcore.Field {
	if l.stackLevel <= level {
		fields = append(fields, stackSkip(2, l.jsonStacktrace))
	}

	return fields
}
