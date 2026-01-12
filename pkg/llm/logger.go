// Package llm 提供 LLM 日志记录功能
package llm

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel 定义日志级别
type LogLevel int

const (
	// LogLevelDebug 调试级别
	LogLevelDebug LogLevel = iota
	// LogLevelInfo 信息级别
	LogLevelInfo
	// LogLevelWarn 警告级别
	LogLevelWarn
	// LogLevelError 错误级别
	LogLevelError
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger LLM 日志记录器
type Logger struct {
	level   LogLevel
	logger  *log.Logger
	enabled bool
}

// NewLogger 创建新的日志记录器
// level: 日志级别
// writer: 日志输出目标（可以是文件或 stderr）
func NewLogger(level LogLevel, writer io.Writer) *Logger {
	return &Logger{
		level:   level,
		logger:  log.New(writer, "", 0),
		enabled: true,
	}
}

// NewFileLogger 创建写入文件的日志记录器
// level: 日志级别
// filepath: 日志文件路径
func NewFileLogger(level LogLevel, filepath string) (*Logger, error) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %w", err)
	}

	return &Logger{
		level:   level,
		logger:  log.New(file, "", 0),
		enabled: true,
	}, nil
}

// NewStderrLogger 创建输出到 stderr 的日志记录器
func NewStderrLogger(level LogLevel) *Logger {
	return NewLogger(level, os.Stderr)
}

// DisabledLogger 返回一个禁用的日志记录器
func DisabledLogger() *Logger {
	return &Logger{
		enabled: false,
	}
}

// IsEnabled 返回日志是否启用
func (l *Logger) IsEnabled() bool {
	return l != nil && l.enabled
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
	if l != nil {
		l.level = level
	}
}

// shouldLog 判断是否应该记录日志
func (l *Logger) shouldLog(level LogLevel) bool {
	return l.IsEnabled() && level >= l.level
}

// log 记录日志
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if !l.shouldLog(level) {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] %s", timestamp, level.String(), message)

	l.logger.Println(logLine)
}

// Debug 记录调试级别日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, format, args...)
}

// Info 记录信息级别日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, format, args...)
}

// Warn 记录警告级别日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LogLevelWarn, format, args...)
}

// Error 记录错误级别日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, format, args...)
}

// LogRequest 记录 LLM 请求（脱敏）
// provider: 提供商名称
// model: 模型名称
// input: 用户输入（不记录完整内容，仅记录长度）
func (l *Logger) LogRequest(provider, model, input string) {
	if !l.IsEnabled() {
		return
	}

	l.Info("LLM 请求 - Provider: %s, Model: %s, Input长度: %d字符",
		provider, model, len(input))
}

// LogResponse 记录 LLM 响应（脱敏）
// provider: 提供商名称
// command: 返回的命令
// duration: 请求耗时
func (l *Logger) LogResponse(provider, command string, duration time.Duration) {
	if !l.IsEnabled() {
		return
	}

	// 脱敏：不记录完整命令，仅记录命令的前10个字符
	safeCommand := maskCommand(command)

	l.Info("LLM 响应 - Provider: %s, 命令: %s, 耗时: %v",
		provider, safeCommand, duration)
}

// LogError 记录 LLM 错误
// provider: 提供商名称
// err: 错误信息
func (l *Logger) LogError(provider string, err error) {
	if !l.IsEnabled() {
		return
	}

	// 脱敏：移除可能包含敏感信息的错误详情
	safeError := maskError(err)

	l.Error("LLM 错误 - Provider: %s, 错误: %s", provider, safeError)
}

// maskCommand 脱敏命令（防止记录敏感参数）
func maskCommand(command string) string {
	if len(command) <= 20 {
		return command
	}
	return command[:20] + "... (truncated)"
}

// maskError 脱敏错误信息（移除可能的 API Key 等敏感信息）
func maskError(err error) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()

	// 脱敏 API Key（查找类似 sk-xxx 或 key=xxx 的模式）
	if strings.Contains(errStr, "sk-") {
		errStr = strings.ReplaceAll(errStr, "sk-", "sk-***")
	}
	if strings.Contains(errStr, "key=") {
		errStr = strings.ReplaceAll(errStr, "key=", "key=***")
	}
	if strings.Contains(errStr, "api_key") {
		errStr = strings.ReplaceAll(errStr, "api_key", "api_key=***")
	}

	return errStr
}

// ParseLogLevel 从字符串解析日志级别
func ParseLogLevel(level string) (LogLevel, error) {
	switch strings.ToLower(level) {
	case "debug":
		return LogLevelDebug, nil
	case "info":
		return LogLevelInfo, nil
	case "warn", "warning":
		return LogLevelWarn, nil
	case "error":
		return LogLevelError, nil
	default:
		return LogLevelInfo, fmt.Errorf("未知的日志级别: %s", level)
	}
}
