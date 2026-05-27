package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Sync() error
}

type loggerZap struct {
	logger *zap.Logger
}

type Field struct {
	Key   string
	Value interface{}
}

func New(cfg *Config) (Logger, error) {
	var zapConfig zapcore.EncoderConfig
	switch cfg.Format {
	case JSONFormat:
		zapConfig = zap.NewProductionEncoderConfig()
	default:
		zapConfig = zap.NewDevelopmentEncoderConfig()
	}
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	zapConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zapConfig.TimeKey = "time"
	zapConfig.LevelKey = "level"
	zapConfig.MessageKey = "message"
	zapConfig.CallerKey = "row"

	var writer zapcore.WriteSyncer
	if cfg.OutputPath == "" {
		writer = zapcore.AddSync(os.Stdout)
	} else {
		file, err := os.OpenFile(cfg.OutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %s, err: %+v", cfg.OutputPath, err)
		}
		writer = zapcore.AddSync(file)
	}
	level := parseLevel(cfg.Level)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zapConfig), writer, level)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &loggerZap{logger: logger}, nil
}

func (l *loggerZap) Info(msg string, fields ...Field) {
	l.logger.Info(msg, toZapFields(fields)...)
}

func (l *loggerZap) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, toZapFields(fields)...)
}

func (l *loggerZap) Error(msg string, fields ...Field) {
	l.logger.Error(msg, toZapFields(fields)...)
}

func (l *loggerZap) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, toZapFields(fields)...)
}

func (l *loggerZap) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, toZapFields(fields)...)
}

func (l *loggerZap) Sync() error {
	return l.logger.Sync()
}

func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	case "warn":
		return zapcore.WarnLevel
	default:
		return zapcore.InfoLevel
	}
}
