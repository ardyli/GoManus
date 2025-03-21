package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	// 日志级别
	LevelDebug = 0
	LevelInfo  = 1
	LevelWarn  = 2
	LevelError = 3
	LevelFatal = 4

	// 当前日志级别
	currentLevel = LevelInfo

	// 日志输出
	logger = log.New(os.Stdout, "", 0)
)

// SetLevel 设置日志级别
func SetLevel(level int) {
	currentLevel = level
}

// GetLevel 获取当前日志级别
func GetLevel() int {
	return currentLevel
}

// Debug 输出调试日志
func Debug(format string, v ...interface{}) {
	if currentLevel <= LevelDebug {
		logWithCaller(LevelDebug, format, v...)
	}
}

// Info 输出信息日志
func Info(format string, v ...interface{}) {
	if currentLevel <= LevelInfo {
		logWithCaller(LevelInfo, format, v...)
	}
}

// Warn 输出警告日志
func Warn(format string, v ...interface{}) {
	if currentLevel <= LevelWarn {
		logWithCaller(LevelWarn, format, v...)
	}
}

// Error 输出错误日志
func Error(format string, v ...interface{}) {
	if currentLevel <= LevelError {
		logWithCaller(LevelError, format, v...)
	}
}

// Fatal 输出致命错误日志并退出程序
func Fatal(format string, v ...interface{}) {
	if currentLevel <= LevelFatal {
		logWithCaller(LevelFatal, format, v...)
	}
	os.Exit(1)
}

// logWithCaller 输出带调用者信息的日志
func logWithCaller(level int, format string, v ...interface{}) {
	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	// 获取文件名
	file = filepath.Base(file)

	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05.000")

	// 获取日志级别标签
	var levelTag string
	switch level {
	case LevelDebug:
		levelTag = "DEBUG"
	case LevelInfo:
		levelTag = "INFO"
	case LevelWarn:
		levelTag = "WARN"
	case LevelError:
		levelTag = "ERROR"
	case LevelFatal:
		levelTag = "FATAL"
	default:
		levelTag = "UNKNOWN"
	}

	// 格式化日志消息
	message := fmt.Sprintf(format, v...)

	// 输出日志
	logger.Printf("%s [%s] %s:%d: %s", now, levelTag, file, line, message)
}
