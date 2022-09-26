package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	L *logger
	fieldList = []string{"userId", "traceId"}
)

type Factory struct {
	logger *logger
}

func GlobalFactory() *logger {
	return L
}

func NewFactory(zapLogger *zap.Logger) *logger {
	forFactory := newFactoryFromZap(zapLogger)

	L = forFactory.logger
	return forFactory.logger
}

func (b Factory) GetZapLogger() *zap.Logger {
	return b.logger.zapLogger
}

func (b Factory) With(fields ...zapcore.Field) Factory {
	return Factory{
		logger: b.logger.With(fields...).(*logger),
	}
}

func newFactoryFromZap(zl *zap.Logger) Factory {
	return Factory{
		logger: &logger{
			zapLogger:  zl,
			stackLevel: zapcore.ErrorLevel,
		},
	}
}
