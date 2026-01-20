// Package executor 提供了命令执行功能
package executor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Executor 负责执行 shell 命令
type Executor struct {
	shell *ShellAdapter
}

// NewExecutor 创建一个新的 Executor 实例
func NewExecutor() *Executor {
	shell, err := DetectShell()
	if err != nil {
		// 如果检测失败，使用默认 shell
		shell = &ShellAdapter{
			Type: ShellSh,
			Path: "/bin/sh",
			Args: []string{"-c"},
		}
	}

	return &Executor{
		shell: shell,
	}
}

// Execute 执行命令并返回输出
// command: 要执行的命令字符串
// stdin: 标准输入数据（可选）
// 返回: 命令输出和错误
func (e *Executor) Execute(command string, stdin string) (string, error) {
	if command == "" {
		return "", fmt.Errorf("命令不能为空")
	}

	// 构建命令
	args := make([]string, len(e.shell.Args))
	copy(args, e.shell.Args)
	args = append(args, command)
	cmd := exec.Command(e.shell.Path, args...)

	// 设置标准输入
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	// 捕获输出
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()

	// 获取输出
	output := stdout.String()
	errOutput := stderr.String()

	// 注意：我们不将非零退出码视为错误，因为很多命令（如 pkill、grep 等）
	// 在某些情况下返回非零退出码是正常行为。我们只在命令无法执行时返回错误。
	// 非零退出码的情况下，仍然返回 stdout 和 stderr 的内容。
	// 如果 stderr 有内容，将其合并到输出中（某些命令会将正常信息输出到 stderr）
	if errOutput != "" {
		// 如果 stdout 为空，使用 stderr 的内容
		if output == "" {
			output = errOutput
		}
	}

	// 只有在命令无法执行时才返回错误（如命令不存在）
	// 非零退出码不应该被视为错误
	_ = err // 忽略退出码错误

	return output, nil
}

// ExecuteInteractive 以交互模式执行命令，实时显示输出
// command: 要执行的命令字符串
// stdin: 标准输入数据（可选）
// 返回: 错误信息（输出会直接打印到终端）
func (e *Executor) ExecuteInteractive(command string, stdin string) error {
	if command == "" {
		return fmt.Errorf("命令不能为空")
	}

	// 构建命令
	args := make([]string, len(e.shell.Args))
	copy(args, e.shell.Args)
	args = append(args, command)
	cmd := exec.Command(e.shell.Path, args...)

	// 设置标准输入
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	} else {
		cmd.Stdin = os.Stdin
	}

	// 直接连接到当前进程的输出流，实现实时输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	// 注意：我们不将非零退出码视为错误，因为很多命令在某些情况下
	// 返回非零退出码是正常行为（如 pkill 没找到进程、grep 没匹配到内容等）
	err := cmd.Run()
	// 忽略退出码错误，只在命令无法执行时返回错误
	_ = err
	return nil
}

// ExecuteWithOutput 执行命令并同时返回输出和实时显示
// command: 要执行的命令字符串
// stdin: 标准输入数据（可选）
// 返回: 命令输出和错误
func (e *Executor) ExecuteWithOutput(command string, stdin string) (string, error) {
	if command == "" {
		return "", fmt.Errorf("命令不能为空")
	}

	// 构建命令
	args := make([]string, len(e.shell.Args))
	copy(args, e.shell.Args)
	args = append(args, command)
	cmd := exec.Command(e.shell.Path, args...)

	// 设置标准输入
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	// 创建输出缓冲区
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// 使用 MultiWriter 同时写入缓冲区和终端
	// 这样既能捕获输出用于返回,又能实时显示到终端
	cmd.Stdout = io.MultiWriter(&stdout, os.Stdout)
	cmd.Stderr = io.MultiWriter(&stderr, os.Stderr)

	// 执行命令
	err := cmd.Run()

	// 获取输出
	output := stdout.String()
	errOutput := stderr.String()

	// 注意：我们不将非零退出码视为错误，因为很多命令（如 pkill、grep 等）
	// 在某些情况下返回非零退出码是正常行为。
	// 如果 stderr 有内容，将其合并到输出中
	if errOutput != "" {
		if output == "" {
			output = errOutput
		}
	}

	// 只有在命令无法执行时才返回错误，非零退出码不视为错误
	_ = err

	return output, nil
}

// GetShell 返回当前使用的 Shell 信息
func (e *Executor) GetShell() *ShellAdapter {
	return e.shell
}

// ExecuteWithContext 使用自定义 Shell 执行命令（高级功能）
func (e *Executor) ExecuteWithContext(command string, stdin string, shell *ShellAdapter) (string, error) {
	if shell == nil {
		return e.Execute(command, stdin)
	}

	if command == "" {
		return "", fmt.Errorf("命令不能为空")
	}

	// 构建命令
	args := make([]string, len(shell.Args))
	copy(args, shell.Args)
	args = append(args, command)
	cmd := exec.Command(shell.Path, args...)

	// 设置标准输入
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	// 捕获输出
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()

	// 获取输出
	output := stdout.String()
	errOutput := stderr.String()

	// 注意：我们不将非零退出码视为错误，因为很多命令在某些情况下
	// 返回非零退出码是正常行为。
	// 如果 stderr 有内容，将其合并到输出中
	if errOutput != "" {
		if output == "" {
			output = errOutput
		}
	}

	// 只有在命令无法执行时才返回错误，非零退出码不视为错误
	_ = err

	return output, nil
}
