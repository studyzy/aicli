// Package llm 提供 LLM 服务的抽象接口
package llm

import "context"

// LLMProvider 定义 LLM 服务提供商的接口
type LLMProvider interface {
	// Translate 将自然语言转换为命令
	// ctx: 上下文，用于超时控制
	// input: 用户的自然语言描述
	// execCtx: 执行上下文信息（操作系统、Shell等）
	// 返回: 转换后的命令字符串和可能的错误
	Translate(ctx context.Context, input string, execCtx *ExecutionContext) (string, error)

	// Name 返回提供商名称
	Name() string
}

// ExecutionContext 包含命令执行的上下文信息
type ExecutionContext struct {
	// OS 操作系统类型（linux/darwin/windows）
	OS string

	// Shell Shell 类型（bash/zsh/powershell/cmd）
	Shell string

	// WorkDir 当前工作目录
	WorkDir string

	// Stdin 标准输入数据（如果有）
	Stdin string
}

// TranslationError 表示翻译过程中的错误
type TranslationError struct {
	Provider string // 提供商名称
	Message  string // 错误消息
	Err      error  // 原始错误
}

func (e *TranslationError) Error() string {
	if e.Err != nil {
		return e.Provider + ": " + e.Message + ": " + e.Err.Error()
	}
	return e.Provider + ": " + e.Message
}

func (e *TranslationError) Unwrap() error {
	return e.Err
}
