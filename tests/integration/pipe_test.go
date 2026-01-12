package integration

import (
	"bytes"
	"context"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/studyzy/aicli/internal/app"
	"github.com/studyzy/aicli/pkg/config"
	"github.com/studyzy/aicli/pkg/executor"
	"github.com/studyzy/aicli/pkg/llm"
	"github.com/studyzy/aicli/pkg/safety"
)

// TestPipeInput_BasicPipeline 测试基本管道输入场景
func TestPipeInput_BasicPipeline(t *testing.T) {
	// 创建 Mock LLM Provider，返回使用 stdin 的命令
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			// 根据输入返回适合的命令
			if strings.Contains(input, "统计") || strings.Contains(input, "行数") {
				if runtime.GOOS == "windows" {
					return "find /c /v \"\""
				}
				return "wc -l"
			}
			if strings.Contains(input, "查找") || strings.Contains(input, "grep") {
				if runtime.GOOS == "windows" {
					return "findstr ERROR"
				}
				return "grep ERROR"
			}
			return "cat"
		},
	}

	cfg := config.Default()
	cfg.LLM.Provider = "mock"

	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	// 测试数据
	stdinData := "line 1\nline 2 ERROR\nline 3\nline 4 ERROR\nline 5"

	flags := app.NewFlags()
	flags.Verbose = false

	output, err := application.Run("查找ERROR", stdinData, flags)
	if err != nil {
		t.Fatalf("管道输入处理失败: %v", err)
	}

	// 验证输出包含 ERROR 行
	if !strings.Contains(output, "ERROR") {
		t.Errorf("期望输出包含 'ERROR', 实际为: %s", output)
	}
}

// TestPipeInput_NoStdin 测试无管道输入时的行为
func TestPipeInput_NoStdin(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			if runtime.GOOS == "windows" {
				return "echo hello"
			}
			return "echo hello"
		},
	}

	cfg := config.Default()
	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	flags := app.NewFlags()

	// 不提供 stdin
	output, err := application.Run("输出hello", "", flags)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if !strings.Contains(output, "hello") {
		t.Errorf("期望输出包含 'hello', 实际为: %s", output)
	}
}

// TestPipeInput_LargeData 测试大量 stdin 数据
func TestPipeInput_LargeData(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			if runtime.GOOS == "windows" {
				return "find /c /v \"\""
			}
			return "wc -l"
		},
	}

	cfg := config.Default()
	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	// 生成 1000 行数据
	var largeData strings.Builder
	for i := 0; i < 1000; i++ {
		largeData.WriteString("line ")
		largeData.WriteString(strings.Repeat("x", 100))
		largeData.WriteString("\n")
	}

	flags := app.NewFlags()

	output, err := application.Run("统计行数", largeData.String(), flags)
	if err != nil {
		t.Fatalf("大数据处理失败: %v", err)
	}

	// 验证有输出
	if len(output) == 0 {
		t.Error("期望有输出，但输出为空")
	}
}

// TestPipeOutput_RedirectToFile 测试输出重定向
func TestPipeOutput_RedirectToFile(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "echo test output"
		},
	}

	cfg := config.Default()
	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	flags := app.NewFlags()

	// 捕获 stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	output, err := application.Run("输出测试", "", flags)

	// 恢复 stdout
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// 读取捕获的输出
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// 验证输出
	if !strings.Contains(output, "test output") {
		t.Errorf("期望输出包含 'test output', 实际为: %s", output)
	}
}

// TestPipeChain_MultipleCommands 测试多级管道
func TestPipeChain_MultipleCommands(t *testing.T) {
	// 这个测试模拟 shell 中的多级管道: cat file | aicli "filter" | aicli "count"

	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			if strings.Contains(input, "filter") || strings.Contains(input, "筛选") {
				if runtime.GOOS == "windows" {
					return "findstr line"
				}
				return "grep line"
			}
			if strings.Contains(input, "count") || strings.Contains(input, "统计") {
				if runtime.GOOS == "windows" {
					return "find /c /v \"\""
				}
				return "wc -l"
			}
			return "cat"
		},
	}

	cfg := config.Default()
	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	flags := app.NewFlags()

	// 第一步: 筛选包含 "line" 的行
	initialData := "line 1\nother 2\nline 3\nother 4\nline 5"
	output1, err := application.Run("筛选包含line的行", initialData, flags)
	if err != nil {
		t.Fatalf("第一步失败: %v", err)
	}

	// 验证第一步输出
	if !strings.Contains(output1, "line") {
		t.Errorf("第一步输出应包含 'line', 实际为: %s", output1)
	}

	// 第二步: 统计行数
	output2, err := application.Run("统计行数", output1, flags)
	if err != nil {
		t.Fatalf("第二步失败: %v", err)
	}

	// 验证第二步有输出（行数）
	if len(output2) == 0 {
		t.Error("第二步输出不应为空")
	}
}

// TestPipeInput_NoSendStdin 测试 --no-send-stdin 标志
func TestPipeInput_NoSendStdin(t *testing.T) {
	callCount := 0
	mockProvider := &llm.MockLLMProvider{
		TranslateFunc: func(ctx context.Context, input string, execCtx *llm.ExecutionContext) (string, error) {
			callCount++
			// 验证 Stdin 字段是否为空
			if execCtx.Stdin != "" {
				t.Errorf("使用 --no-send-stdin 时，Stdin 应为空，实际为: %s", execCtx.Stdin)
			}
			return "echo done", nil
		},
	}

	cfg := config.Default()
	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	stdinData := "sensitive data\npassword: 12345"

	flags := app.NewFlags()
	flags.NoSendStdin = true

	_, err := application.Run("处理数据", stdinData, flags)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if callCount == 0 {
		t.Error("LLM Provider 未被调用")
	}
}

// TestPipeInput_EmptyStdin 测试空 stdin
func TestPipeInput_EmptyStdin(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			return "echo no input"
		},
	}

	cfg := config.Default()
	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	flags := app.NewFlags()

	output, err := application.Run("测试空输入", "", flags)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if !strings.Contains(output, "no input") {
		t.Errorf("期望输出包含 'no input', 实际为: %s", output)
	}
}

// TestPipeInput_BinaryData 测试二进制数据处理
func TestPipeInput_BinaryData(t *testing.T) {
	mockProvider := &llm.MockLLMProvider{
		TranslateFn: func(input string) string {
			// 对于二进制数据，返回简单命令
			return "echo binary processed"
		},
	}

	cfg := config.Default()
	application := app.NewApp(cfg, mockProvider, executor.NewExecutor(), safety.NewSafetyChecker(false))

	// 模拟二进制数据（包含 null 字节）
	binaryData := "text\x00binary\x01data\xff"

	flags := app.NewFlags()

	output, err := application.Run("处理二进制", binaryData, flags)
	if err != nil {
		t.Fatalf("二进制数据处理失败: %v", err)
	}

	if len(output) == 0 {
		t.Error("期望有输出")
	}
}
