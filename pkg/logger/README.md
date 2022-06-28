# Logger module

This module, pkg or folder contain customs types and methods for log managements using as a base [`zap`](go.uber.org/zap)

Useful method 

```go

// Info wrap for log.info
func (al *AppLogger) Info(message string, fields ...F) {
	al.Log.Info(message, parseFieldToZap(fields...)...)
}

// Debug wrap for log.Debug
func (al *AppLogger) Debug(message string, fields ...F) {
	al.Log.Debug(message, parseFieldToZap(fields...)...)
}

// Error wrap for log.Error
func (al *AppLogger) Error(message string, fields ...F) {
	al.Log.Error(message, parseFieldToZap(fields...)...)
}

// Fatal wrap for log.Error
func (al *AppLogger) Fatal(message string, fields ...F) {
	al.Log.Error(message, parseFieldToZap(fields...)...)
}

// Warn wrap for log.Warn
func (al *AppLogger) Warn(message string, fields ...F) {
	al.Log.Warn(message, parseFieldToZap(fields...)...)
}
```

