// Package integration 包含端到端集成测试
package integration

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/studyzy/aicli/internal/app"
	"github.com/studyzy/aicli/pkg/config"
	"github.com/studyzy/aicli/pkg/executor"
	"github.com/studyzy/aicli/pkg/i18n"
	"github.com/studyzy/aicli/pkg/llm"
	"github.com/studyzy/aicli/pkg/safety"
)

func init() {
	// 初始化 i18n (默认中文)
	i18n.Init(config.Default())
}

const (
	osWindows           = "windows"
	providerMock        = "mock"
	commandEchoTest     = "echo test"
	commandCat          = "cat"
	commandFindstrError = "findstr ERROR"
	commandGrepError    = "grep ERROR"
)

// TestBasicCommandTranslation 测试基本命令转换和执行
func TestBasicCommandTranslation(t *testing.T) {
	// 使用 Mock LLM Provider
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			// 根据输入返回模拟命令
			if strings.Contains(input, "列出") || strings.Contains(input, "显示") {
				if runtime.GOOS == osWindows {
					return "dir"
				}
				return "ls"
			}
			return commandEchoTest
		},
	}

	// 创建配置
	cfg := config.Default()
	cfg.LLM.Provider = providerMock

	// 创建应用实例（不启用安全检查）
	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewChecker(false))

	// 测试执行
	flags := &app.Flags{
		Verbose: false,
		DryRun:  false,
		Force:   true, // 跳过确认
	}

	output, err := application.Run("列出当前目录", "", flags)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// 验证输出不为空
	if output == "" {
		t.Error("期望有输出，但输出为空")
	}

	t.Logf("输出: %s", output)
}

// TestCommandWithStdin 测试带 stdin 的命令
func TestCommandWithStdin(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			if strings.Contains(input, "过滤") || strings.Contains(input, "grep") {
				if runtime.GOOS == osWindows {
					return commandFindstrError
				}
				return commandGrepError
			}
			return commandCat
		},
	}

	cfg := config.Default()
	cfg.LLM.Provider = providerMock

	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewChecker(false))

	flags := &app.Flags{
		Force: true,
	}

	stdin := "line 1\nERROR: something wrong\nline 3"
	output, err := application.Run("过滤出包含ERROR的行", stdin, flags)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if !strings.Contains(output, "ERROR") {
		t.Errorf("期望输出包含 'ERROR', 实际为: %s", output)
	}
}

// TestDryRunMode 测试 dry-run 模式
func TestDryRunMode(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return commandEchoTest
		},
	}

	cfg := config.Default()
	cfg.LLM.Provider = providerMock

	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewChecker(false))

	flags := &app.Flags{
		DryRun: true,
	}

	output, err := application.Run("显示测试", "", flags)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// dry-run 模式应该只显示命令，不执行
	if !strings.Contains(output, "echo test") {
		t.Errorf("期望输出包含命令 'echo test', 实际为: %s", output)
	}
}

// TestDangerousCommandDetection 测试危险命令检测
func TestDangerousCommandDetection(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "rm -rf /tmp/test"
		},
	}

	cfg := config.Default()
	cfg.LLM.Provider = providerMock

	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewChecker(true))

	flags := &app.Flags{
		Force: false, // 不强制执行，应该被拦截
	}

	_, err := application.Run("删除测试文件", "", flags)

	// 期望返回错误，因为需要用户确认但没有提供
	if err == nil {
		t.Fatal("期望危险命令被拦截，但执行成功了")
	}

	if !strings.Contains(err.Error(), "危险") && !strings.Contains(err.Error(), "确认") {
		t.Logf("错误信息: %v", err)
	}
}

// TestVerboseMode 测试详细模式
func TestVerboseMode(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "echo verbose test"
		},
	}

	cfg := config.Default()
	cfg.LLM.Provider = providerMock

	// 捕获标准输出
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewChecker(false))

	flags := &app.Flags{
		Verbose: true,
		Force:   true,
	}

	_, err := application.Run("详细测试", "", flags)

	// 恢复标准输出
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Logf("执行可能失败: %v", err)
	}

	// 读取捕获的输出
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	captured := string(buf[:n])

	// 详细模式应该输出更多信息
	if captured != "" {
		t.Logf("捕获的输出: %s", captured)
	}
}

// TestEmptyInput 测试空输入
func TestEmptyInput(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return ""
		},
	}

	cfg := config.Default()
	cfg.LLM.Provider = providerMock

	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewChecker(false))

	flags := &app.Flags{}

	_, err := application.Run("", "", flags)
	if err == nil {
		t.Fatal("期望空输入返回错误，但成功了")
	}
}
