package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Logger 日志工具
type Logger struct {
	logger  *log.Logger
	prefix  string
	level   string
	output  string
	file    string
	fileHandle io.Writer
}

// LogLevel 日志级别常量
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// NewLogger 创建新的日志实例
// 参数：level, output, file
func NewLogger(level, output, file string) (*Logger, error) {
	var outputWriter io.Writer

	// 设置输出位置
	if output == "file" && file != "" {
		fileHandle, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("打开日志文件失败: %v", err)
		}
		outputWriter = fileHandle
	} else {
		outputWriter = os.Stdout
	}

	// 设置日志格式
	logger := log.New(outputWriter, "", log.Ldate|log.Ltime|log.Lshortfile)

	l := &Logger{
		logger:  logger,
		prefix:  "",
		level:   level,
		output:  output,
		file:    file,
		fileHandle: outputWriter,
	}

	return l, nil
}

// Debug 输出调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level == LevelDebug {
		l.logger.Printf("[DEBUG] "+format, args...)
	}
}

// Info 输出信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level == LevelDebug || l.level == LevelInfo {
		l.logger.Printf("[INFO] "+format, args...)
	}
}

// Warn 输出警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level == LevelDebug || l.level == LevelInfo || l.level == LevelWarn {
		l.logger.Printf("[WARN] "+format, args...)
	}
}

// Error 输出错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+format, args...)
}

// Success 输出成功日志
func (l *Logger) Success(format string, args ...interface{}) {
	if l.level == LevelDebug || l.level == LevelInfo || l.level == LevelWarn {
		l.logger.Printf("[SUCCESS] "+format, args...)
	}
}

// WithPrefix 设置日志前缀
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		logger:  l.logger,
		prefix:  l.prefix + "[" + prefix + "] ",
		level:   l.level,
		output:  l.output,
		file:    l.file,
		fileHandle: l.fileHandle,
	}
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	if l.fileHandle != nil {
		if file, ok := l.fileHandle.(*os.File); ok {
			return file.Close()
		}
	}
	return nil
}
