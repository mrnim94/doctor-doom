package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type DoomLogger struct{}

var zapLogger *zap.Logger

func DoomLoggerInit(logDir string) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logDir + "/doom-" + time.Now().GoString() + ".log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	zapLogger := zap.New(core)
	zapLogger.Info("Logger initialized")
}

func (dl DoomLogger) Info(event string, filePath string, lastModifiedUnix int64, size int64) {
	zapLogger.Info(event, zap.String("filePath", filePath), zap.Int64("lastModifiedUnix", lastModifiedUnix), zap.Int64("size", size))
}

func (dl DoomLogger) Error(event string, filePath string, lastModifiedUnix int64, size int64) {
	zapLogger.Error(event, zap.String("filePath", filePath), zap.Int64("lastModifiedUnix", lastModifiedUnix), zap.Int64("size", size))
}
