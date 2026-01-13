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

	// 如果有错误，合并 stderr 到错误信息
	if err != nil {
		if errOutput != "" {
			return output, fmt.Errorf("命令执行失败: %w\n%s", err, errOutput)
		}
		return output, fmt.Errorf("命令执行失败: %w", err)
	}

	// 即使成功，也可能有 stderr 输出（如警告信息），但在非交互式执行中，
	// 我们通常只关心 stdout 的结果。stderr 中的内容（如进度条）可能会污染输出。
	// 因此，在成功执行的情况下，我们忽略 stderr。
	// if errOutput != "" {
	// 	output = output + errOutput
	// }

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
	return cmd.Run()
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
	cmd.Stdout = io.MultiWriter(&stdout, os.Stdout)
	cmd.Stderr = io.MultiWriter(&stderr, os.Stderr)

	// 执行命令
	err := cmd.Run()

	// 获取输出
	output := stdout.String()
	errOutput := stderr.String()

	// 如果有错误，合并 stderr 到错误信息
	if err != nil {
		if errOutput != "" {
			return output, fmt.Errorf("命令执行失败: %w\n%s", err, errOutput)
		}
		return output, fmt.Errorf("命令执行失败: %w", err)
	}

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

	// 处理错误
	if err != nil {
		if errOutput != "" {
			return output, fmt.Errorf("命令执行失败: %w\n%s", err, errOutput)
		}
		return output, fmt.Errorf("命令执行失败: %w", err)
	}

	// if errOutput != "" {
	// 	output = output + errOutput
	// }

	return output, nil
}
