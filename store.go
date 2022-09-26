package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	outputType string

	store struct {
		cfg   *StoreConfig
		write zapcore.WriteSyncer
		level zapcore.LevelEnabler
	}

	StoreConfig struct {
		Level      string             `json:"level" yaml:"level"`
		Output     outputType         `json:"output" yaml:"output"`
		StackLevel string             `json:"stackLevel" yaml:"stackLevel"`
		Lumberjack *lumberjack.Logger `json:"lumberjack" yaml:"lumberjack"`
	}
)

const (
	OutputConsole outputType = "console"
	OutputFile    outputType = "file"
)

func NewStore(cfg *StoreConfig) *store {
	return &store{
		cfg:   cfg,
		write: zapcore.AddSync(cfg.Lumberjack),
		level: unmarshalTextLevel(cfg.Level),
	}
}

func (s *store) JsonEncoder() *zap.Logger {
	var encoder zapcore.Encoder

	switch s.cfg.Output {
	case OutputFile:
		encoder = zapcore.NewJSONEncoder(newProductionEncoderConfig())
	default:
		encoder = zapcore.NewConsoleEncoder(newConsoleZapLogger())
	}

	return zap.New(zapcore.NewCore(encoder, s.write, s.level), zap.AddCaller(), zap.AddCallerSkip(2))
}

func unmarshalTextLevel(text string) zapcore.Level {
	switch text {
	case "debug", "DEBUG":
		return zapcore.DebugLevel

	case "info", "INFO", "":
		return zapcore.InfoLevel

	case "warn", "WARN":
		return zapcore.WarnLevel

	case "error", "ERROR":
		return zapcore.ErrorLevel

	case "dpanic", "DPANIC":
		return zapcore.DPanicLevel

	case "panic", "PANIC":
		return zapcore.PanicLevel

	case "fatal", "FATAL":
		return zapcore.FatalLevel

	default:
		return zapcore.InfoLevel
	}
}

func newProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  stacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newConsoleZapLogger() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  stacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

}
