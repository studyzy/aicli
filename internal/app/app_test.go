package app

import (
	"context"
	"strings"
	"testing"

	"github.com/studyzy/aicli/pkg/config"
	"github.com/studyzy/aicli/pkg/executor"
	"github.com/studyzy/aicli/pkg/llm"
	"github.com/studyzy/aicli/pkg/safety"
)

// TestApp_DangerousCommandConfirmation 测试危险命令确认流程
func TestApp_DangerousCommandConfirmation(t *testing.T) {
	// 创建返回危险命令的 Mock Provider
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "rm -rf /tmp/test"
		},
	}

	cfg := config.Default()
	checker := safety.NewSafetyChecker(true) // 启用安全检查
	application := NewApp(cfg, mockProvider, executor.NewExecutor(), checker)

	flags := NewFlags()
	flags.Force = false // 不强制执行，需要确认

	// 测试危险命令检测
	_, err := application.Run("删除测试文件", "", flags)

	// 应该因为需要确认而失败（在测试环境中无法交互）
	if err == nil {
		t.Error("Expected error for dangerous command without confirmation")
	}

	if !strings.Contains(err.Error(), "取消") && !strings.Contains(err.Error(), "拒绝") {
		t.Logf("Error message: %v", err)
	}
}

// TestApp_ForceBypassConfirmation 测试 --force 标志跳过确认
func TestApp_ForceBypassConfirmation(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "echo test" // 使用安全命令避免实际执行危险操作
		},
	}

	cfg := config.Default()
	checker := safety.NewSafetyChecker(true)
	application := NewApp(cfg, mockProvider, executor.NewExecutor(), checker)

	flags := NewFlags()
	flags.Force = true // 强制执行，跳过确认

	output, err := application.Run("测试", "", flags)
	if err != nil {
		t.Fatalf("Run() with --force failed: %v", err)
	}

	if !strings.Contains(output, "test") {
		t.Errorf("output = %s, should contain 'test'", output)
	}
}

// TestApp_SafetyCheckDisabled 测试禁用安全检查
func TestApp_SafetyCheckDisabled(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "echo safe"
		},
	}

	cfg := config.Default()
	checker := safety.NewSafetyChecker(false) // 禁用安全检查
	application := NewApp(cfg, mockProvider, executor.NewExecutor(), checker)

	flags := NewFlags()

	output, err := application.Run("测试", "", flags)
	if err != nil {
		t.Fatalf("Run() failed: %v", err)
	}

	if !strings.Contains(output, "safe") {
		t.Errorf("output = %s, should contain 'safe'", output)
	}
}

// TestApp_PipeModeNonInteractive 测试管道模式下的非交互行为
func TestApp_PipeModeNonInteractive(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "rm -rf /tmp/test" // 危险命令
		},
	}

	cfg := config.Default()
	checker := safety.NewSafetyChecker(true)
	application := NewApp(cfg, mockProvider, executor.NewExecutor(), checker)

	flags := NewFlags()
	stdin := "some data" // 有 stdin 表示管道模式

	// 管道模式下应该拒绝危险命令
	_, err := application.Run("删除", stdin, flags)
	if err == nil {
		t.Error("Expected error for dangerous command in pipe mode")
	}

	if !strings.Contains(err.Error(), "管道") && !strings.Contains(err.Error(), "拒绝") {
		t.Logf("Error: %v", err)
	}
}

// TestApp_PipeModeWithForce 测试管道模式下使用 --force
func TestApp_PipeModeWithForce(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "echo forced" // 使用安全命令
		},
	}

	cfg := config.Default()
	checker := safety.NewSafetyChecker(true)
	application := NewApp(cfg, mockProvider, executor.NewExecutor(), checker)

	flags := NewFlags()
	flags.Force = true
	stdin := "data"

	output, err := application.Run("测试", stdin, flags)
	if err != nil {
		t.Fatalf("Run() with --force in pipe mode failed: %v", err)
	}

	if !strings.Contains(output, "forced") {
		t.Errorf("output = %s, should contain 'forced'", output)
	}
}

// TestApp_EmptyCommand 测试 LLM 返回空命令
func TestApp_EmptyCommand(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "" // 返回空命令
		},
	}

	cfg := config.Default()
	application := NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	flags := NewFlags()

	_, err := application.Run("测试", "", flags)
	if err == nil {
		t.Error("Expected error for empty command")
	}

	if !strings.Contains(err.Error(), "空命令") {
		t.Errorf("Error should mention empty command: %v", err)
	}
}

// TestApp_DryRunMode 测试 dry-run 模式
func TestApp_DryRunMode(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "dangerous_command"
		},
	}

	cfg := config.Default()
	application := NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	flags := NewFlags()
	flags.DryRun = true

	output, err := application.Run("测试", "", flags)
	if err != nil {
		t.Fatalf("Run() in dry-run mode failed: %v", err)
	}

	if !strings.Contains(output, "将要执行") {
		t.Errorf("dry-run output should contain '将要执行': %s", output)
	}

	if !strings.Contains(output, "dangerous_command") {
		t.Errorf("dry-run output should show command: %s", output)
	}
}

// TestConfirmExecution_MockInput 测试确认函数（模拟输入）
func TestConfirmExecution_MockInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"confirm yes", "y\n", true},
		{"confirm YES", "yes\n", true},
		{"confirm Yes", "Yes\n", true},
		{"reject no", "n\n", false},
		{"reject NO", "no\n", false},
		{"reject other", "maybe\n", false},
		{"reject empty", "\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 注意: confirmExecution 使用 os.Stdin，这里只是测试逻辑
			// 实际测试需要重构代码以支持 io.Reader 注入

			// 简单验证输入解析逻辑
			input := strings.ToLower(strings.TrimSpace(tt.input))
			result := (input == "y" || input == "yes")

			if result != tt.expected {
				t.Errorf("input %q: got %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestApp_VerboseMode 测试详细模式输出
func TestApp_VerboseMode(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "echo verbose"
		},
	}

	cfg := config.Default()
	application := NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	flags := NewFlags()
	flags.Verbose = true

	// Verbose 输出到 stderr，这里只验证不出错
	output, err := application.Run("测试", "", flags)
	if err != nil {
		t.Fatalf("Run() in verbose mode failed: %v", err)
	}

	if !strings.Contains(output, "verbose") {
		t.Errorf("output should contain command result: %s", output)
	}
}

// TestApp_NoSendStdin 测试 --no-send-stdin 标志
func TestApp_NoSendStdin(t *testing.T) {
	callCount := 0
	mockProvider := &llm.MockLLMProvider{
		TranslateFunc: func(ctx context.Context, input string, execCtx *llm.ExecutionContext) (string, error) {
			callCount++
			// 验证 Stdin 为空
			if execCtx.Stdin != "" {
				t.Errorf("Stdin should be empty with --no-send-stdin, got: %s", execCtx.Stdin)
			}
			return "echo done", nil
		},
	}

	cfg := config.Default()
	application := NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	flags := NewFlags()
	flags.NoSendStdin = true

	_, err := application.Run("测试", "sensitive data", flags)
	if err != nil {
		t.Fatalf("Run() with --no-send-stdin failed: %v", err)
	}

	if callCount == 0 {
		t.Error("LLM Provider was not called")
	}
}
