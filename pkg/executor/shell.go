// Package executor 提供命令执行功能
package executor

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// ShellType 表示 Shell 类型
type ShellType string

const (
	// ShellBash 表示 Bash Shell
	ShellBash ShellType = "bash"

	// ShellZsh 表示 Zsh Shell
	ShellZsh ShellType = "zsh"

	// ShellPowerShell 表示 PowerShell
	ShellPowerShell ShellType = "powershell"

	// ShellCmd 表示 Windows CMD
	ShellCmd ShellType = "cmd"

	// ShellSh 表示 POSIX sh
	ShellSh ShellType = "sh"
)

// ShellAdapter 表示 Shell 适配器
type ShellAdapter struct {
	// Type Shell 类型
	Type ShellType

	// Path Shell 可执行文件路径
	Path string

	// Args 执行命令时的参数模板
	Args []string
}

// DetectShell 检测当前系统的 Shell
func DetectShell() (*ShellAdapter, error) {
	// 优先检查环境变量
	if shellPath := os.Getenv("SHELL"); shellPath != "" {
		return detectUnixShell(shellPath)
	}

	// 根据操作系统检测
	switch runtime.GOOS {
	case "windows":
		return detectWindowsShell()
	case "darwin", "linux":
		return detectUnixShellDefault()
	default:
		return nil, fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}

// detectUnixShell 从路径检测 Unix Shell 类型
func detectUnixShell(shellPath string) (*ShellAdapter, error) {
	shellPath = strings.ToLower(shellPath)

	if strings.Contains(shellPath, "zsh") {
		return &ShellAdapter{
			Type: ShellZsh,
			Path: shellPath,
			Args: []string{"-c"},
		}, nil
	}

	if strings.Contains(shellPath, "bash") {
		return &ShellAdapter{
			Type: ShellBash,
			Path: shellPath,
			Args: []string{"-c"},
		}, nil
	}

	// 默认使用 sh
	return &ShellAdapter{
		Type: ShellSh,
		Path: shellPath,
		Args: []string{"-c"},
	}, nil
}

// detectUnixShellDefault 检测 Unix 系统的默认 Shell
func detectUnixShellDefault() (*ShellAdapter, error) {
	// 尝试常见的 Shell 路径
	shells := []struct {
		path      string
		shellType ShellType
	}{
		{"/bin/zsh", ShellZsh},
		{"/usr/bin/zsh", ShellZsh},
		{"/bin/bash", ShellBash},
		{"/usr/bin/bash", ShellBash},
		{"/bin/sh", ShellSh},
	}

	for _, shell := range shells {
		if _, err := os.Stat(shell.path); err == nil {
			return &ShellAdapter{
				Type: shell.shellType,
				Path: shell.path,
				Args: []string{"-c"},
			}, nil
		}
	}

	return nil, fmt.Errorf("未找到可用的 Shell")
}

// detectWindowsShell 检测 Windows Shell
func detectWindowsShell() (*ShellAdapter, error) {
	// 优先使用 PowerShell
	if psPath, err := findExecutable("powershell.exe"); err == nil {
		return &ShellAdapter{
			Type: ShellPowerShell,
			Path: psPath,
			Args: []string{"-NoProfile", "-Command"},
		}, nil
	}

	// 回退到 CMD
	if cmdPath, err := findExecutable("cmd.exe"); err == nil {
		return &ShellAdapter{
			Type: ShellCmd,
			Path: cmdPath,
			Args: []string{"/C"},
		}, nil
	}

	return nil, fmt.Errorf("未找到可用的 Windows Shell")
}

// findExecutable 在 PATH 中查找可执行文件
func findExecutable(name string) (string, error) {
	// 在 Windows 上，Go 的 exec.LookPath 会自动处理
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", fmt.Errorf("PATH 环境变量为空")
	}

	separator := ":"
	if runtime.GOOS == "windows" {
		separator = ";"
	}

	paths := strings.Split(pathEnv, separator)
	for _, p := range paths {
		fullPath := p + string(os.PathSeparator) + name
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("未找到可执行文件: %s", name)
}

// String 返回 Shell 的字符串表示
func (s *ShellAdapter) String() string {
	return fmt.Sprintf("%s (%s)", s.Type, s.Path)
}

// GetShellType 返回 Shell 类型字符串
func (s *ShellAdapter) GetShellType() string {
	return string(s.Type)
}
