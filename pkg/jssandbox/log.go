package jssandbox

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// defaultLogger 是默认的全局日志实例
	defaultLogger *logrus.Logger
)

func init() {
	// 初始化默认日志实例
	defaultLogger = logrus.New()
	defaultLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
	defaultLogger.SetLevel(logrus.InfoLevel)
	defaultLogger.SetOutput(os.Stdout)
}

// GetLogger 获取默认的日志实例
func GetLogger() *logrus.Logger {
	return defaultLogger
}

// SetLoggerLevel 设置日志级别
func SetLoggerLevel(level logrus.Level) {
	defaultLogger.SetLevel(level)
}

// SetLoggerFormatter 设置日志格式化器
func SetLoggerFormatter(formatter logrus.Formatter) {
	defaultLogger.SetFormatter(formatter)
}

// SetLoggerOutput 设置日志输出
func SetLoggerOutput(output io.Writer) {
	defaultLogger.SetOutput(output)
}

// NewLogger 创建一个新的日志实例
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
	return logger
}

// NewLoggerWithConfig 使用指定配置创建新的日志实例
func NewLoggerWithConfig(level logrus.Level, formatter logrus.Formatter, output io.Writer) *logrus.Logger {
	logger := logrus.New()
	if formatter != nil {
		logger.SetFormatter(formatter)
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}
	logger.SetLevel(level)
	if output != nil {
		logger.SetOutput(output)
	} else {
		logger.SetOutput(os.Stdout)
	}
	return logger
}
