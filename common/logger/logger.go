package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type DoomLogger struct {
	logger *zap.Logger
}

func (d *DoomLogger) New(logDir string) {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logDir + "/doom-destroy.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)

	d.logger = zap.New(core)
}

func (dl DoomLogger) InfoVictim(event string, filePath string, lastModifiedUnix int64, size int64) {
	dl.logger.Info(event, zap.String("filePath", filePath), zap.Int64("lastModifiedUnix", lastModifiedUnix), zap.Int64("size", size))
}

func (dl DoomLogger) ErrorVictim(event string, filePath string, lastModifiedUnix int64, size int64) {
	dl.logger.Error(event, zap.String("filePath", filePath), zap.Int64("lastModifiedUnix", lastModifiedUnix), zap.Int64("size", size))
}

func (dl DoomLogger) Info(event string, message string) {
	dl.logger.Info(event, zap.String("message", message))
}
