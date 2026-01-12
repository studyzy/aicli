package llm

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

// TestNewLogger 测试创建日志记录器
func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)

	if logger == nil {
		t.Fatal("期望非 nil logger")
	}
	if !logger.IsEnabled() {
		t.Error("期望 logger 启用")
	}
	if logger.level != LogLevelInfo {
		t.Errorf("期望级别 INFO, 实际 %v", logger.level)
	}
}

// TestDisabledLogger 测试禁用的日志记录器
func TestDisabledLogger(t *testing.T) {
	logger := DisabledLogger()

	if logger == nil {
		t.Fatal("期望非 nil logger")
	}
	if logger.IsEnabled() {
		t.Error("期望 logger 禁用")
	}
}

// TestLogLevel_String 测试日志级别字符串转换
func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LogLevelDebug, "DEBUG"},
		{LogLevelInfo, "INFO"},
		{LogLevelWarn, "WARN"},
		{LogLevelError, "ERROR"},
		{LogLevel(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.level.String()
			if result != tt.expected {
				t.Errorf("期望 %s, 实际 %s", tt.expected, result)
			}
		})
	}
}

// TestLogger_Debug 测试 Debug 级别日志
func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelDebug, &buf)

	logger.Debug("test message: %s", "value")

	output := buf.String()
	if !strings.Contains(output, "[DEBUG]") {
		t.Errorf("期望包含 [DEBUG], 实际: %s", output)
	}
	if !strings.Contains(output, "test message: value") {
		t.Errorf("期望包含消息, 实际: %s", output)
	}
}

// TestLogger_Info 测试 Info 级别日志
func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)

	logger.Info("info message")

	output := buf.String()
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("期望包含 [INFO], 实际: %s", output)
	}
	if !strings.Contains(output, "info message") {
		t.Errorf("期望包含消息, 实际: %s", output)
	}
}

// TestLogger_Warn 测试 Warn 级别日志
func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelWarn, &buf)

	logger.Warn("warning message")

	output := buf.String()
	if !strings.Contains(output, "[WARN]") {
		t.Errorf("期望包含 [WARN], 实际: %s", output)
	}
}

// TestLogger_Error 测试 Error 级别日志
func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelError, &buf)

	logger.Error("error message")

	output := buf.String()
	if !strings.Contains(output, "[ERROR]") {
		t.Errorf("期望包含 [ERROR], 实际: %s", output)
	}
}

// TestLogger_LogLevelFiltering 测试日志级别过滤
func TestLogger_LogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelWarn, &buf) // 只记录 Warn 及以上

	logger.Debug("debug message") // 不应记录
	logger.Info("info message")   // 不应记录
	logger.Warn("warn message")   // 应记录
	logger.Error("error message") // 应记录

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Error("不应记录 debug 消息")
	}
	if strings.Contains(output, "info message") {
		t.Error("不应记录 info 消息")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("应记录 warn 消息")
	}
	if !strings.Contains(output, "error message") {
		t.Error("应记录 error 消息")
	}
}

// TestLogger_SetLevel 测试设置日志级别
func TestLogger_SetLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)

	logger.Debug("before change") // 不应记录
	buf.Reset()

	logger.SetLevel(LogLevelDebug)
	logger.Debug("after change") // 应记录

	output := buf.String()
	if !strings.Contains(output, "after change") {
		t.Error("更改级别后应记录 debug 消息")
	}
}

// TestLogger_DisabledNoOutput 测试禁用的 logger 不输出
func TestLogger_DisabledNoOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelDebug, &buf)
	logger.enabled = false

	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")

	if buf.Len() > 0 {
		t.Errorf("禁用的 logger 不应产生输出, 实际: %s", buf.String())
	}
}

// TestLogger_LogRequest 测试记录请求
func TestLogger_LogRequest(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)

	logger.LogRequest("openai", "gpt-4", "测试输入内容")

	output := buf.String()
	if !strings.Contains(output, "LLM 请求") {
		t.Error("期望包含 'LLM 请求'")
	}
	if !strings.Contains(output, "openai") {
		t.Error("期望包含 provider 名称")
	}
	if !strings.Contains(output, "gpt-4") {
		t.Error("期望包含 model 名称")
	}
	if strings.Contains(output, "测试输入内容") {
		t.Error("不应记录完整输入内容（隐私保护）")
	}
}

// TestLogger_LogResponse 测试记录响应
func TestLogger_LogResponse(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)

	logger.LogResponse("openai", "ls -la /very/long/path/to/file", 1*time.Second)

	output := buf.String()
	if !strings.Contains(output, "LLM 响应") {
		t.Error("期望包含 'LLM 响应'")
	}
	if !strings.Contains(output, "openai") {
		t.Error("期望包含 provider 名称")
	}
	if !strings.Contains(output, "1s") {
		t.Error("期望包含耗时")
	}
}

// TestLogger_LogError 测试记录错误
func TestLogger_LogError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelError, &buf)

	err := errors.New("API error with key=sk-12345")
	logger.LogError("openai", err)

	output := buf.String()
	if !strings.Contains(output, "LLM 错误") {
		t.Error("期望包含 'LLM 错误'")
	}
	if !strings.Contains(output, "openai") {
		t.Error("期望包含 provider 名称")
	}
	if strings.Contains(output, "sk-12345") {
		t.Error("不应记录完整 API Key（脱敏失败）")
	}
}

// TestMaskCommand 测试命令脱敏
func TestMaskCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{"短命令", "ls -la", "ls -la"},
		{"长命令", "this is a very long command that should be truncated", "this is a very long ... (truncated)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskCommand(tt.command)
			if result != tt.expected {
				t.Errorf("期望 %s, 实际 %s", tt.expected, result)
			}
		})
	}
}

// TestMaskError 测试错误脱敏
func TestMaskError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contains string
		notContains string
	}{
		{
			"nil 错误",
			nil,
			"",
			"",
		},
		{
			"包含 sk- 的错误",
			errors.New("authentication failed with sk-12345"),
			"sk-***",
			"sk-12345",
		},
		{
			"包含 key= 的错误",
			errors.New("invalid key=abcdef"),
			"key=***",
			"key=abcdef",
		},
		{
			"包含 api_key 的错误",
			errors.New("missing api_key parameter"),
			"api_key=***",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskError(tt.err)
			
			if tt.contains != "" && !strings.Contains(result, tt.contains) {
				t.Errorf("期望包含 %s, 实际: %s", tt.contains, result)
			}
			if tt.notContains != "" && strings.Contains(result, tt.notContains) {
				t.Errorf("不应包含 %s, 实际: %s", tt.notContains, result)
			}
		})
	}
}

// TestParseLogLevel 测试解析日志级别
func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
		hasError bool
	}{
		{"debug", LogLevelDebug, false},
		{"DEBUG", LogLevelDebug, false},
		{"info", LogLevelInfo, false},
		{"INFO", LogLevelInfo, false},
		{"warn", LogLevelWarn, false},
		{"warning", LogLevelWarn, false},
		{"WARN", LogLevelWarn, false},
		{"error", LogLevelError, false},
		{"ERROR", LogLevelError, false},
		{"invalid", LogLevelInfo, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseLogLevel(tt.input)
			
			if tt.hasError {
				if err == nil {
					t.Error("期望返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误, 实际: %v", err)
				}
				if result != tt.expected {
					t.Errorf("期望 %v, 实际 %v", tt.expected, result)
				}
			}
		})
	}
}

// TestNewStderrLogger 测试创建 stderr logger
func TestNewStderrLogger(t *testing.T) {
	logger := NewStderrLogger(LogLevelInfo)
	
	if logger == nil {
		t.Fatal("期望非 nil logger")
	}
	if !logger.IsEnabled() {
		t.Error("期望 logger 启用")
	}
}
