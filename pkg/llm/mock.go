package llm

import "context"

// MockLLMProvider 是用于测试的 Mock LLM 提供商
type MockLLMProvider struct {
	// TranslateFunc 自定义翻译函数
	TranslateFunc func(ctx context.Context, input string, execCtx *ExecutionContext) (string, error)

	// TranslateFn 简化版翻译函数（只接受 input）
	TranslateFn func(input string) string

	// ProviderName 提供商名称
	ProviderName string
}

// Translate 执行翻译（调用自定义函数）
func (m *MockLLMProvider) Translate(ctx context.Context, input string, execCtx *ExecutionContext) (string, error) {
	if m.TranslateFunc != nil {
		return m.TranslateFunc(ctx, input, execCtx)
	}

	// 简化版：只使用 input
	if m.TranslateFn != nil {
		return m.TranslateFn(input), nil
	}

	return "", &TranslationError{
		Provider: m.Name(),
		Message:  "TranslateFunc not implemented",
	}
}

// Name 返回提供商名称
func (m *MockLLMProvider) Name() string {
	if m.ProviderName != "" {
		return m.ProviderName
	}
	return "mock"
}

// NewMockProvider 创建一个新的 Mock 提供商
func NewMockProvider() *MockLLMProvider {
	return &MockLLMProvider{
		ProviderName: "mock",
	}
}

// NewMockProviderWithFunc 创建带有自定义翻译函数的 Mock 提供商
func NewMockProviderWithFunc(
	fn func(ctx context.Context, input string, execCtx *ExecutionContext) (string, error),
) *MockLLMProvider {
	return &MockLLMProvider{
		TranslateFunc: fn,
		ProviderName:  "mock",
	}
}
