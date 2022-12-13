package logger

import (
	"go.uber.org/zap"
	"log"
)

var zapLog *zap.Logger

func Init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	zapLog = logger
}

func Info(message string, fields ...zap.Field) {
	zapLog.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	zapLog.Debug(message, fields...)
}

func Any(fields ...interface{}) {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any("", field))
	}
	zapLog.Debug("", zapFields...)
}

func Warn(message string, fields ...zap.Field) {
	zapLog.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	zapLog.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	zapLog.Fatal(message, fields...)
}
