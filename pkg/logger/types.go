package logger

import "go.uber.org/zap"

type KVLog map[string]interface{}

type Logger interface {
	Info(message string, fields ...KVLog)
	Debug(message string, fields ...KVLog)
	Error(message string, fields ...KVLog)
	Fatal(message string, fields ...KVLog)
	Warn(message string, fields ...KVLog)
	Sync() error
}

type AppLogger struct {
	Log *zap.Logger
}
