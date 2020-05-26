package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var (
	logger *zap.Logger
)

func New() {
	logConfig := zap.NewProductionEncoderConfig()
	logConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(logConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer("errors")),
		zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zap.ErrorLevel
		}),
	)
	warnCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(logConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer("warn")),
		zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level >= zap.WarnLevel && level < zap.ErrorLevel
		}),
	)
	infoCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(logConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer("info")),
		zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level > zap.DebugLevel && level < zap.WarnLevel
		}),
	)
	debugCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(logConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writeSyncer("defug")),
		zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level <= zap.DebugLevel
		}),
	)

	logger = zap.New(zapcore.NewTee(errorCore, warnCore, infoCore, debugCore))
	logger.Sync()
}

func writeSyncer(level string) zapcore.WriteSyncer {
	now := time.Now()
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("data/logs/lchat-%s-%04d%02d%02d.log", level, now.Year(), now.Month(), now.Day()),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress: true,
	})

}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
	defer logger.Sync()
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
	defer logger.Sync()
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
	defer logger.Sync()
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
	defer logger.Sync()
}