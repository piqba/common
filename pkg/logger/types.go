package logger

import "go.uber.org/zap"

// F minimal key value type for a field representation in our logger.Logger interface
type F map[string]interface{}

type Logger interface {
	Info(message string, fields ...F)
	Debug(message string, fields ...F)
	Error(message string, fields ...F)
	Fatal(message string, fields ...F)
	Warn(message string, fields ...F)
	Sync() error
}

type AppLogger struct {
	Log *zap.Logger
}
