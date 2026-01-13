package executor

import (
	"runtime"
	"strings"
	"testing"
)

func TestDetectShell(t *testing.T) {
	shell, err := DetectShell()
	if err != nil {
		t.Fatalf("检测 Shell 失败: %v", err)
	}

	if shell == nil {
		t.Fatal("检测到的 Shell 为 nil")
	}

	if shell.Type == "" {
		t.Error("Shell 类型为空")
	}

	if shell.Path == "" {
		t.Error("Shell 路径为空")
	}

	if len(shell.Args) == 0 {
		t.Error("Shell 参数为空")
	}

	t.Logf("检测到的 Shell: %s", shell.String())
}

func TestDetectShellByOS(t *testing.T) {
	shell, err := DetectShell()
	if err != nil {
		t.Fatalf("检测 Shell 失败: %v", err)
	}

	// 根据操作系统验证 Shell 类型
	switch runtime.GOOS {
	case osWindows:
		if shell.Type != ShellPowerShell && shell.Type != ShellCmd {
			t.Errorf("Windows 上期望 PowerShell 或 CMD, 实际为 %s", shell.Type)
		}
	case "darwin", "linux":
		validTypes := []ShellType{ShellBash, ShellZsh, ShellSh}
		valid := false
		for _, vt := range validTypes {
			if shell.Type == vt {
				valid = true
				break
			}
		}
		if !valid {
			t.Errorf("Unix 系统上期望 bash/zsh/sh, 实际为 %s", shell.Type)
		}
	}
}

func TestShellAdapterString(t *testing.T) {
	shell := &ShellAdapter{
		Type: ShellBash,
		Path: "/bin/bash",
		Args: []string{"-c"},
	}

	str := shell.String()
	if str == "" {
		t.Error("Shell 字符串表示为空")
	}

	if str != "bash (/bin/bash)" {
		t.Errorf("期望字符串为 'bash (/bin/bash)', 实际为 '%s'", str)
	}
}

func TestShellAdapterGetShellType(t *testing.T) {
	shell := &ShellAdapter{
		Type: ShellZsh,
		Path: "/bin/zsh",
		Args: []string{"-c"},
	}

	shellType := shell.GetShellType()
	if shellType != "zsh" {
		t.Errorf("期望 Shell 类型为 'zsh', 实际为 '%s'", shellType)
	}
}

func TestDetectUnixShell(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected ShellType
	}{
		{"Bash Shell", "/bin/bash", ShellBash},
		{"Zsh Shell", "/usr/bin/zsh", ShellZsh},
		{"Zsh 大写", "/usr/bin/ZSH", ShellZsh},
		{"Bash 大写", "/BIN/BASH", ShellBash},
		{"其他 Shell", "/bin/ksh", ShellSh},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell, err := detectUnixShell(tt.path)
			if err != nil {
				t.Fatalf("检测 Shell 失败: %v", err)
			}

			if shell.Type != tt.expected {
				t.Errorf("期望类型 %s, 实际为 %s", tt.expected, shell.Type)
			}

			// 路径会被转换为小写，所以比较时应该忽略大小写
			expectedPath := strings.ToLower(tt.path)
			actualPath := strings.ToLower(shell.Path)
			if actualPath != expectedPath {
				t.Errorf("期望路径 %s, 实际为 %s", expectedPath, actualPath)
			}
		})
	}
}

func TestShellTypes(t *testing.T) {
	// 验证所有定义的 Shell 类型常量
	types := []ShellType{
		ShellBash,
		ShellZsh,
		ShellPowerShell,
		ShellCmd,
		ShellSh,
	}

	for _, st := range types {
		if string(st) == "" {
			t.Errorf("Shell 类型 %v 为空字符串", st)
		}
	}
}

// TestDetectShellCrossPlatform 测试跨平台 Shell 检测
func TestDetectShellCrossPlatform(t *testing.T) {
	shell, err := DetectShell()
	if err != nil {
		// 某些测试环境可能没有 Shell
		t.Skipf("跳过测试，无法检测 Shell: %v", err)
	}

	// 验证基本属性
	if shell.Type == "" {
		t.Error("Shell 类型不应为空")
	}

	if shell.Path == "" {
		t.Error("Shell 路径不应为空")
	}

	// 验证参数格式
	switch shell.Type {
	case ShellBash, ShellZsh, ShellSh:
		if len(shell.Args) != 1 || shell.Args[0] != "-c" {
			t.Errorf("Unix Shell 参数应为 ['-c'], 实际为 %v", shell.Args)
		}
	case ShellPowerShell:
		if len(shell.Args) < 2 {
			t.Errorf("PowerShell 参数数量不足: %v", shell.Args)
		}
	case ShellCmd:
		if len(shell.Args) != 1 || shell.Args[0] != "/C" {
			t.Errorf("CMD 参数应为 ['/C'], 实际为 %v", shell.Args)
		}
	}
}

// TestExecutor_Execute_Success 测试成功执行命令
func TestExecutor_Execute_Success(t *testing.T) {
	executor := NewExecutor()

	// 执行简单的跨平台命令
	const cmdHello = "echo hello"
	cmd := cmdHello

	output, err := executor.Execute(cmd, "")
	if err != nil {
		t.Fatalf("执行命令失败: %v", err)
	}

	if !strings.Contains(output, "hello") {
		t.Errorf("期望输出包含 'hello', 实际为: %s", output)
	}
}

// TestExecutor_Execute_WithStdin 测试带 stdin 的命令执行
func TestExecutor_Execute_WithStdin(t *testing.T) {
	executor := NewExecutor()

	// 使用 cat (Unix) 或 findstr (Windows) 读取 stdin
	const cmdCat = "cat"
	const cmdFindstr = "findstr .*"
	var cmd string
	if runtime.GOOS == osWindows {
		cmd = cmdFindstr
	} else {
		cmd = cmdCat
	}

	stdin := "test input data\nline 2\nline 3"
	output, err := executor.Execute(cmd, stdin)
	if err != nil {
		t.Fatalf("执行命令失败: %v", err)
	}

	if !strings.Contains(output, "test input data") {
		t.Errorf("期望输出包含 'test input data', 实际为: %s", output)
	}
}

// TestExecutor_Execute_CommandFailed 测试命令执行失败
func TestExecutor_Execute_CommandFailed(t *testing.T) {
	executor := NewExecutor()

	// 执行不存在的命令
	_, err := executor.Execute("nonexistent-command-12345", "")
	if err == nil {
		t.Fatal("期望返回错误，但成功了")
	}

	if !strings.Contains(err.Error(), "执行") && !strings.Contains(err.Error(), "exit") {
		t.Logf("错误信息: %v", err)
	}
}

// TestExecutor_Execute_EmptyCommand 测试空命令
func TestExecutor_Execute_EmptyCommand(t *testing.T) {
	executor := NewExecutor()

	_, err := executor.Execute("", "")
	if err == nil {
		t.Fatal("期望返回错误（空命令），但成功了")
	}
}

// TestExecutor_GetShell 测试获取 Shell 信息
func TestExecutor_GetShell(t *testing.T) {
	executor := NewExecutor()

	shell := executor.GetShell()
	if shell == nil {
		t.Fatal("Shell 不应为 nil")
	}

	if shell.Type == "" {
		t.Error("Shell 类型不应为空")
	}

	if shell.Path == "" {
		t.Error("Shell 路径不应为空")
	}
}

// TestExecutor_Execute_IgnoreStderr 测试成功执行时忽略 stderr
func TestExecutor_Execute_IgnoreStderr(t *testing.T) {
	if runtime.GOOS == osWindows {
		t.Skip("跳过 Windows 测试，因为命令语法不同")
	}

	executor := NewExecutor()

	// 这个命令会输出 "stdout output" 到 stdout，输出 "stderr output" 到 stderr
	cmd := "echo stdout output; echo stderr output >&2"

	output, err := executor.Execute(cmd, "")
	if err != nil {
		t.Fatalf("执行命令失败: %v", err)
	}

	// 验证 stdout 包含期望的内容
	if !strings.Contains(output, "stdout output") {
		t.Errorf("期望输出包含 'stdout output', 实际为: %q", output)
	}

	// 验证 stderr 不包含在 output 中
	if strings.Contains(output, "stderr output") {
		t.Errorf("期望输出不包含 'stderr output', 实际为: %q", output)
	}
}

// TestExecutor_Execute_IncludeStderrOnError 测试失败执行时包含 stderr
func TestExecutor_Execute_IncludeStderrOnError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("跳过 Windows 测试，因为命令语法不同")
	}

	executor := NewExecutor()

	// 这个命令会输出 "stderr output" 到 stderr，并以非零状态退出
	cmd := "echo stderr output >&2; exit 1"

	output, err := executor.Execute(cmd, "")
	if err == nil {
		t.Fatal("期望返回错误，但成功了")
	}

	// 验证错误信息包含 stderr 内容
	if !strings.Contains(err.Error(), "stderr output") {
		t.Errorf("期望错误信息包含 'stderr output', 实际为: %q", err.Error())
	}

	// output 应该是空的（或者包含 stdout 的内容，如果有的话）
	if output != "" {
		t.Logf("Output: %q", output)
	}
}
