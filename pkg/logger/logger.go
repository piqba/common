package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

// NewLogger return a custom AppLogger using zap lib and receive by params the log lvl
// by default is info lvl
// # E.g (not sensitive case)
//  -1 debug
//   0 info
//   1 warn
//   2 error
func NewLogger(lvl string) *AppLogger {
	config := zap.NewProductionConfig()
	config.Level.SetLevel(parseLvl(lvl))
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.StacktraceKey = "tracerKey"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig

	log, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err.Error())
	}
	return &AppLogger{Log: log}
}

func parseLvl(lvl string) zapcore.Level {
	switch strings.ToLower(lvl) {
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

// parseFieldToZap transform ...KVLog (map[string]interface{}) to []zap.Field
func parseFieldToZap(fields ...KVLog) []zap.Field {
	var fieldsZap []zap.Field
	for _, field := range fields {
		for k, v := range field {
			switch v.(type) {
			case int:
				newField := zap.Int(k, v.(int))
				fieldsZap = append(fieldsZap, newField)
			case float64:
				newField := zap.Float64(k, v.(float64))
				fieldsZap = append(fieldsZap, newField)
			case float32:
				newField := zap.Float32(k, v.(float32))
				fieldsZap = append(fieldsZap, newField)
			case time.Duration:
				newField := zap.Duration(k, v.(time.Duration))
				fieldsZap = append(fieldsZap, newField)
			default:
				newField := zap.String(k, v.(string))
				fieldsZap = append(fieldsZap, newField)
			}
		}
	}
	return fieldsZap
}

// Info wrap for log.info
func (al *AppLogger) Info(message string, fields ...KVLog) {
	al.Log.Info(message, parseFieldToZap(fields...)...)
}

// Debug wrap for log.Debug
func (al *AppLogger) Debug(message string, fields ...KVLog) {
	al.Log.Debug(message, parseFieldToZap(fields...)...)
}

// Error wrap for log.Error
func (al *AppLogger) Error(message string, fields ...KVLog) {
	al.Log.Error(message, parseFieldToZap(fields...)...)
}

// Fatal wrap for log.Error
func (al *AppLogger) Fatal(message string, fields ...KVLog) {
	al.Log.Error(message, parseFieldToZap(fields...)...)
}

// Warn wrap for log.Warn
func (al *AppLogger) Warn(message string, fields ...KVLog) {
	al.Log.Warn(message, parseFieldToZap(fields...)...)
}

// Sync wrap for log.Sync
func (al *AppLogger) Sync() error {
	return al.Log.Sync()
}
