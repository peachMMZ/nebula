package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func Init(mode string) {
	// 使用绝对路径，确保日志文件位置可预期
	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}
	logDir := filepath.Join(filepath.Dir(execPath), "logs")

	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("failed to create log directory: " + err.Error())
	}

	logFilePath := filepath.Join(logDir, "app.log")

	// 创建 lumberjack writer（支持日志轮转）
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
		Compress:   true,
	})

	var core zapcore.Core

	if mode == gin.DebugMode {
		// Debug 模式：同时输出到控制台和文件
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleWriter := zapcore.Lock(os.Stdout)

		// 合并文件和控制台输出
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleWriter, zap.DebugLevel),
			zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				fileWriter,
				zap.DebugLevel,
			),
		)
	} else {
		// 生产模式：仅输出到文件
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "time"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			fileWriter,
			zap.InfoLevel,
		)
	}

	Log = zap.New(core, zap.AddCaller())
	fmt.Printf("[Logger] Log file location: %s\n", logFilePath)
}
