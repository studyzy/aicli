// Package llm 提供了提示词模板构建功能
package llm

import (
	"fmt"
	"strings"
)

// GetSystemPrompt 返回系统提示词
func GetSystemPrompt(ctx *ExecutionContext) string {
	var sb strings.Builder

	sb.WriteString("你是一个命令行助手，专门将用户的自然语言描述转换为可执行的 shell 命令。\n\n")
	sb.WriteString("规则：\n")
	sb.WriteString("1. 只返回命令本身，不要有任何解释或说明\n")
	sb.WriteString("2. 不要使用 markdown 代码块格式\n")
	sb.WriteString("3. 命令必须是可以直接执行的\n")
	sb.WriteString("4. 如果需要多个命令，使用 && 或 ; 连接\n")
	sb.WriteString("5. 优先使用常见且兼容性好的命令\n\n")

	if ctx != nil {
		sb.WriteString("执行环境：\n")
		sb.WriteString(fmt.Sprintf("- 操作系统: %s\n", ctx.OS))
		sb.WriteString(fmt.Sprintf("- Shell: %s\n", ctx.Shell))
		sb.WriteString(fmt.Sprintf("- 工作目录: %s\n", ctx.WorkDir))
	}

	return sb.String()
}

// BuildPrompt 构建用户提示词
func BuildPrompt(input string, ctx *ExecutionContext) string {
	var sb strings.Builder

	// 添加用户输入
	sb.WriteString(fmt.Sprintf("将以下自然语言描述转换为命令：\n%s\n", input))

	// 如果有标准输入，添加上下文
	if ctx != nil && ctx.Stdin != "" {
		sb.WriteString("\n标准输入数据：\n")

		// 限制 stdin 的长度，避免提示词过长
		maxStdinLen := 500
		stdin := ctx.Stdin
		if len(stdin) > maxStdinLen {
			stdin = stdin[:maxStdinLen] + "... (已截断)"
		}

		sb.WriteString(stdin)
		sb.WriteString("\n")
	}

	return sb.String()
}

// BuildContextDescription 构建执行上下文描述（用于调试和日志）
func BuildContextDescription(ctx *ExecutionContext) string {
	if ctx == nil {
		return "无执行上下文"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("OS: %s, Shell: %s, WorkDir: %s", ctx.OS, ctx.Shell, ctx.WorkDir))

	if ctx.Stdin != "" {
		stdinLen := len(ctx.Stdin)
		if stdinLen > 50 {
			sb.WriteString(fmt.Sprintf(", Stdin: %d bytes", stdinLen))
		} else {
			sb.WriteString(fmt.Sprintf(", Stdin: %q", ctx.Stdin))
		}
	}

	return sb.String()
}
