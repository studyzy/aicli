// Package app 提供了应用主逻辑
package app

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/studyzy/aicli/internal/history"
	"github.com/studyzy/aicli/pkg/config"
	"github.com/studyzy/aicli/pkg/executor"
	"github.com/studyzy/aicli/pkg/llm"
	"github.com/studyzy/aicli/pkg/safety"
)

// App 是应用主结构
type App struct {
	config   *config.Config
	llm      llm.LLMProvider
	executor *executor.Executor
	safety   *safety.SafetyChecker
	history  *history.History
}

// NewApp 创建一个新的应用实例
func NewApp(cfg *config.Config, provider llm.LLMProvider, exec *executor.Executor, checker *safety.SafetyChecker) *App {
	return &App{
		config:   cfg,
		llm:      provider,
		executor: exec,
		safety:   checker,
		history:  history.NewHistory(),
	}
}

// SetHistory 设置历史记录管理器
func (a *App) SetHistory(h *history.History) {
	a.history = h
}

// Run 执行应用主逻辑
// input: 用户的自然语言输入
// stdin: 标准输入数据（来自管道）
// flags: 命令行标志
// 返回: 命令输出和错误
func (a *App) Run(input string, stdin string, flags *Flags) (string, error) {
	// 验证输入
	if input == "" {
		return "", fmt.Errorf("请提供自然语言描述")
	}

	// 详细模式：显示输入
	if flags.Verbose {
		fmt.Fprintf(os.Stderr, "自然语言输入: %s\n", input)
		if stdin != "" {
			fmt.Fprintf(os.Stderr, "标准输入: %d 字节\n", len(stdin))
		}
	}

	// 构建执行上下文
	execCtx := a.buildExecutionContext(stdin, flags)

	// 详细模式：显示上下文
	if flags.Verbose {
		fmt.Fprintf(os.Stderr, "执行上下文: %s\n", llm.BuildContextDescription(execCtx))
	}

	// 调用 LLM 转换命令
	startTime := time.Now()

	ctx := context.Background()
	if a.config.LLM.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(a.config.LLM.Timeout)*time.Second)
		defer cancel()
	}

	command, err := a.llm.Translate(ctx, input, execCtx)
	if err != nil {
		return "", fmt.Errorf("命令转换失败: %w", err)
	}

	translateTime := time.Since(startTime)

	// 验证命令不为空
	if command == "" {
		return "", fmt.Errorf("LLM 返回空命令")
	}

	// 详细模式：显示转换结果
	if flags.Verbose {
		fmt.Fprintf(os.Stderr, "转换后的命令: %s\n", command)
		fmt.Fprintf(os.Stderr, "转换耗时: %v\n", translateTime)
	}

	// 安全检查
	if a.safety != nil && a.safety.IsEnabled() {
		isDangerous, description, riskLevel := a.safety.IsDangerous(command)
		if isDangerous {
			// 管道模式下不进行交互式确认
			if a.isPipeMode(stdin) {
				if !flags.Force {
					return "", fmt.Errorf("管道模式下拒绝执行危险命令（使用 --force 强制执行）")
				}
			} else {
				// 非管道模式：如果没有强制执行，需要用户确认
				if !flags.Force {
					if !confirmDangerousCommand(command, description, riskLevel.String()) {
						return "", fmt.Errorf("用户取消执行危险命令")
					}
				}
			}
		}
	}

	// Dry-run 模式：只显示命令不执行
	if flags.DryRun {
		return fmt.Sprintf("将要执行的命令: %s", command), nil
	}

	// 执行命令
	if flags.Verbose {
		fmt.Fprintf(os.Stderr, "开始执行命令...\n")
	}

	execStartTime := time.Now()
	output, err := a.executor.Execute(command, stdin)
	execTime := time.Since(execStartTime)

	// 保存历史记录
	if a.history != nil {
		entry := &history.HistoryEntry{
			Input:     input,
			Command:   command,
			Timestamp: time.Now(),
			Success:   err == nil,
			ExitCode:  0,
		}

		if err != nil {
			entry.Error = err.Error()
			entry.ExitCode = 1
		} else {
			// 截断输出（避免历史文件过大）
			if len(output) > 500 {
				entry.Output = output[:500] + "... (truncated)"
			} else {
				entry.Output = output
			}
		}

		a.history.Add(entry)
	}

	if err != nil {
		return output, fmt.Errorf("命令执行失败: %w", err)
	}

	// 详细模式：显示执行时间
	if flags.Verbose {
		fmt.Fprintf(os.Stderr, "执行耗时: %v\n", execTime)
		fmt.Fprintf(os.Stderr, "总耗时: %v\n", translateTime+execTime)
	}

	return output, nil
}

// buildExecutionContext 构建执行上下文
func (a *App) buildExecutionContext(stdin string, flags *Flags) *llm.ExecutionContext {
	// 获取当前工作目录
	workDir, _ := os.Getwd()

	// 获取 Shell 信息
	shell := a.executor.GetShell()

	ctx := &llm.ExecutionContext{
		OS:      runtime.GOOS,
		Shell:   shell.GetShellType(),
		WorkDir: workDir,
	}

	// 添加 stdin（除非禁用）
	if !flags.NoSendStdin && stdin != "" {
		ctx.Stdin = stdin
	}

	return ctx
}

// isPipeMode 检测是否处于管道模式
// 管道模式下（有 stdin 输入）不应该进行交互式确认
func (a *App) isPipeMode(stdin string) bool {
	return stdin != ""
}

// GetHistory 获取历史记录管理器
func (a *App) GetHistory() *history.History {
	return a.history
}
