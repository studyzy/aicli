// Package llm 提供 LLM Provider 工厂函数
package llm

import (
	"fmt"
	"strings"

	"github.com/studyzy/aicli/pkg/config"
)

const (
	providerOpenAI   = "openai"
	providerLocal    = "local"
	providerMock     = "mock"
	defaultOllamaURL = "http://localhost:11434"
)

// NewProvider 根据配置创建对应的 Provider
// 这是工厂函数，根据配置中的 Provider 类型创建相应的实现
func NewProvider(cfg *config.Config) (Provider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("配置不能为空")
	}

	providerType := strings.ToLower(cfg.LLM.Provider)

	switch providerType {
	case providerBuiltin:
		return NewBuiltinProvider(), nil

	case providerOpenAI:
		return newOpenAIFromConfig(cfg)

	case providerAnthropic, "claude":
		return newAnthropicFromConfig(cfg)

	case providerLocal, "ollama":
		return newLocalModelFromConfig(cfg)

	case providerMock:
		// Mock provider 用于测试
		return &MockLLMProvider{
			TranslateFn: func(input string) string {
				return "echo mock: " + input
			},
		}, nil

	default:
		return nil, fmt.Errorf("不支持的 LLM 提供商: %s", cfg.LLM.Provider)
	}
}

// newOpenAIFromConfig 从配置创建 OpenAI Provider
func newOpenAIFromConfig(cfg *config.Config) (Provider, error) {
	// 内置 provider 不需要 API Key
	if cfg.LLM.Provider == providerBuiltin {
		return NewBuiltinProvider(), nil
	}
	
	if cfg.LLM.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API 密钥未配置")
	}

	model := cfg.LLM.Model
	if model == "" {
		model = "gpt-4" // 默认模型
	}

	return NewOpenAIProvider(cfg.LLM.APIKey, model, cfg.LLM.APIBase), nil
}

// newAnthropicFromConfig 从配置创建 Anthropic Provider
func newAnthropicFromConfig(cfg *config.Config) (Provider, error) {
	if cfg.LLM.APIKey == "" {
		return nil, fmt.Errorf("anthropic API 密钥未配置")
	}

	model := cfg.LLM.Model
	if model == "" {
		model = "claude-3-sonnet-20240229" // 默认模型
	}

	return NewAnthropicProvider(cfg.LLM.APIKey, model, cfg.LLM.APIBase), nil
}

// newLocalModelFromConfig 从配置创建本地模型 Provider
func newLocalModelFromConfig(cfg *config.Config) (Provider, error) {
	model := cfg.LLM.Model
	if model == "" {
		model = "llama2" // 默认模型
	}

	baseURL := cfg.LLM.APIBase
	if baseURL == "" {
		baseURL = defaultOllamaURL // Ollama 默认地址
	}

	return NewLocalModelProvider(model, baseURL), nil
}

// GetSupportedProviders 返回所有支持的 Provider 列表
func GetSupportedProviders() []string {
	return []string{
		"builtin",
		"openai",
		"anthropic",
		"claude",
		"local",
		"ollama",
		"mock",
	}
}

// IsProviderSupported 检查给定的 Provider 是否被支持
func IsProviderSupported(provider string) bool {
	provider = strings.ToLower(provider)
	for _, p := range GetSupportedProviders() {
		if p == provider {
			return true
		}
	}
	return false
}
