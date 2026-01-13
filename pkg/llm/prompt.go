// Package llm 提供了提示词模板构建功能
package llm

import (
	"fmt"
	"strings"

	"github.com/studyzy/aicli/pkg/i18n"
)

// GetSystemPrompt 返回系统提示词
// lang: 语言代码 (zh/en)
func GetSystemPrompt(ctx *ExecutionContext) string {
	lang := i18n.Lang()
	
	if lang == "en" {
		return buildEnglishSystemPrompt(ctx)
	}
	return buildChineseSystemPrompt(ctx)
}

// buildChineseSystemPrompt 构建中文系统提示词
func buildChineseSystemPrompt(ctx *ExecutionContext) string {
	var sb strings.Builder

	sb.WriteString(i18n.T(i18n.LLMSystemPromptIntro) + "\n\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRules) + "\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRule1) + "\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRule2) + "\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRule3) + "\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRule4) + "\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRule5) + "\n\n")

	if ctx != nil {
		sb.WriteString(i18n.T(i18n.LLMSystemPromptEnv) + "\n")
		sb.WriteString(fmt.Sprintf("- %s: %s\n", i18n.T(i18n.LabelOS), ctx.OS))
		sb.WriteString(fmt.Sprintf("- %s: %s\n", i18n.T(i18n.LabelShell), ctx.Shell))
		sb.WriteString(fmt.Sprintf("- %s: %s\n", i18n.T(i18n.LabelWorkDir), ctx.WorkDir))
	}

	return sb.String()
}

// buildEnglishSystemPrompt 构建英文系统提示词
func buildEnglishSystemPrompt(ctx *ExecutionContext) string {
	var sb strings.Builder

	sb.WriteString(i18n.T(i18n.LLMSystemPromptIntro) + "\n\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRules) + "\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRule1) + "\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRule2) + "\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRule3) + "\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRule4) + "\n")
	sb.WriteString(i18n.T(i18n.LLMSystemPromptRule5) + "\n\n")

	if ctx != nil {
		sb.WriteString(i18n.T(i18n.LLMSystemPromptEnv) + "\n")
		sb.WriteString(fmt.Sprintf("- %s: %s\n", i18n.T(i18n.LabelOS), ctx.OS))
		sb.WriteString(fmt.Sprintf("- %s: %s\n", i18n.T(i18n.LabelShell), ctx.Shell))
		sb.WriteString(fmt.Sprintf("- %s: %s\n", i18n.T(i18n.LabelWorkDir), ctx.WorkDir))
	}

	return sb.String()
}

// BuildPrompt 构建用户提示词
func BuildPrompt(input string, ctx *ExecutionContext) string {
	var sb strings.Builder

	// 添加用户输入
	sb.WriteString(fmt.Sprintf("%s\n%s\n", i18n.T(i18n.LLMUserPromptIntro), input))

	// 如果有标准输入，添加上下文
	if ctx != nil && ctx.Stdin != "" {
		sb.WriteString("\n" + i18n.T(i18n.LLMStdinData) + "\n")

		// 限制 stdin 的长度，避免提示词过长
		maxStdinLen := 500
		stdin := ctx.Stdin
		if len(stdin) > maxStdinLen {
			stdin = stdin[:maxStdinLen] + i18n.T(i18n.LLMTruncated)
		}

		sb.WriteString(stdin)
		sb.WriteString("\n")
	}

	return sb.String()
}

// BuildContextDescription 构建执行上下文描述（用于调试和日志）
func BuildContextDescription(ctx *ExecutionContext) string {
	if ctx == nil {
		return i18n.T(i18n.LLMContextNoContext)
	}

	var sb strings.Builder
	sb.WriteString(i18n.T(i18n.LLMContextFormat, ctx.OS, ctx.Shell, ctx.WorkDir))

	if ctx.Stdin != "" {
		stdinLen := len(ctx.Stdin)
		if stdinLen > 50 {
			sb.WriteString(fmt.Sprintf(", %s", i18n.T(i18n.LabelStdinBytes, stdinLen)))
		} else {
			sb.WriteString(fmt.Sprintf(", %s: %q", i18n.T(i18n.LabelStdin), ctx.Stdin))
		}
	}

	return sb.String()
}
