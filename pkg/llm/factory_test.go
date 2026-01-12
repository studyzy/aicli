package llm

import (
	"testing"

	"github.com/studyzy/aicli/pkg/config"
)

// TestNewProvider_OpenAI 测试创建 OpenAI Provider
func TestNewProvider_OpenAI(t *testing.T) {
	cfg := &config.Config{
		LLM: config.LLMConfig{
			Provider: "openai",
			APIKey:   "test-key",
			Model:    "gpt-4",
			APIBase:  "https://api.openai.com/v1",
		},
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("创建 OpenAI Provider 失败: %v", err)
	}

	if provider.Name() != "openai" {
		t.Errorf("期望 Provider 名称为 'openai', 实际为 '%s'", provider.Name())
	}
}

// TestNewProvider_OpenAI_MissingAPIKey 测试缺少 API Key 的情况
func TestNewProvider_OpenAI_MissingAPIKey(t *testing.T) {
	cfg := &config.Config{
		LLM: config.LLMConfig{
			Provider: "openai",
			APIKey:   "", // 缺少 API Key
			Model:    "gpt-4",
		},
	}

	_, err := NewProvider(cfg)
	if err == nil {
		t.Fatal("期望返回错误（缺少 API Key），但成功了")
	}
}

// TestNewProvider_Anthropic 测试创建 Anthropic Provider
func TestNewProvider_Anthropic(t *testing.T) {
	cfg := &config.Config{
		LLM: config.LLMConfig{
			Provider: "anthropic",
			APIKey:   "test-key",
			Model:    "claude-3-sonnet-20240229",
		},
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("创建 Anthropic Provider 失败: %v", err)
	}

	if provider.Name() != "anthropic" {
		t.Errorf("期望 Provider 名称为 'anthropic', 实际为 '%s'", provider.Name())
	}
}

// TestNewProvider_Anthropic_AliasClaude 测试 claude 别名
func TestNewProvider_Anthropic_AliasClaude(t *testing.T) {
	cfg := &config.Config{
		LLM: config.LLMConfig{
			Provider: "claude", // 使用别名
			APIKey:   "test-key",
			Model:    "claude-3-sonnet-20240229",
		},
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("创建 Claude Provider 失败: %v", err)
	}

	if provider.Name() != "anthropic" {
		t.Errorf("期望 Provider 名称为 'anthropic', 实际为 '%s'", provider.Name())
	}
}

// TestNewProvider_LocalModel 测试创建本地模型 Provider
func TestNewProvider_LocalModel(t *testing.T) {
	cfg := &config.Config{
		LLM: config.LLMConfig{
			Provider: "local",
			Model:    "llama2",
			APIBase:  "http://localhost:11434",
		},
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("创建本地模型 Provider 失败: %v", err)
	}

	if provider.Name() != "local" {
		t.Errorf("期望 Provider 名称为 'local', 实际为 '%s'", provider.Name())
	}
}

// TestNewProvider_LocalModel_AliasOllama 测试 ollama 别名
func TestNewProvider_LocalModel_AliasOllama(t *testing.T) {
	cfg := &config.Config{
		LLM: config.LLMConfig{
			Provider: "ollama", // 使用别名
			Model:    "mistral",
		},
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("创建 Ollama Provider 失败: %v", err)
	}

	if provider.Name() != "local" {
		t.Errorf("期望 Provider 名称为 'local', 实际为 '%s'", provider.Name())
	}
}

// TestNewProvider_Mock 测试创建 Mock Provider
func TestNewProvider_Mock(t *testing.T) {
	cfg := &config.Config{
		LLM: config.LLMConfig{
			Provider: "mock",
		},
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("创建 Mock Provider 失败: %v", err)
	}

	// Mock provider 没有 Name() 方法，但应该能正常创建
	if provider == nil {
		t.Error("期望创建 Mock Provider，但返回 nil")
	}
}

// TestNewProvider_Unsupported 测试不支持的 Provider
func TestNewProvider_Unsupported(t *testing.T) {
	cfg := &config.Config{
		LLM: config.LLMConfig{
			Provider: "unsupported-provider",
		},
	}

	_, err := NewProvider(cfg)
	if err == nil {
		t.Fatal("期望返回错误（不支持的 Provider），但成功了")
	}
}

// TestNewProvider_NilConfig 测试 nil 配置
func TestNewProvider_NilConfig(t *testing.T) {
	_, err := NewProvider(nil)
	if err == nil {
		t.Fatal("期望返回错误（nil 配置），但成功了")
	}
}

// TestNewProvider_DefaultModel 测试默认模型
func TestNewProvider_DefaultModel(t *testing.T) {
	tests := []struct {
		name        string
		provider    string
		apiKey      string
		expectError bool
		expectName  string
	}{
		{
			name:        "OpenAI 默认模型",
			provider:    "openai",
			apiKey:      "test-key",
			expectError: false,
			expectName:  "openai",
		},
		{
			name:        "Anthropic 默认模型",
			provider:    "anthropic",
			apiKey:      "test-key",
			expectError: false,
			expectName:  "anthropic",
		},
		{
			name:        "Local 默认模型",
			provider:    "local",
			apiKey:      "", // 本地模型不需要 API Key
			expectError: false,
			expectName:  "local",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				LLM: config.LLMConfig{
					Provider: tt.provider,
					APIKey:   tt.apiKey,
					Model:    "", // 不指定模型，使用默认值
				},
			}

			provider, err := NewProvider(cfg)
			if tt.expectError {
				if err == nil {
					t.Fatal("期望返回错误，但成功了")
				}
				return
			}

			if err != nil {
				t.Fatalf("期望成功，但返回错误: %v", err)
			}

			if provider.Name() != tt.expectName {
				t.Errorf("期望 Provider 名称为 '%s', 实际为 '%s'", tt.expectName, provider.Name())
			}
		})
	}
}

// TestGetSupportedProviders 测试获取支持的 Provider 列表
func TestGetSupportedProviders(t *testing.T) {
	providers := GetSupportedProviders()

	if len(providers) == 0 {
		t.Fatal("期望返回至少一个 Provider")
	}

	expectedProviders := map[string]bool{
		"openai":    true,
		"anthropic": true,
		"local":     true,
		"mock":      true,
	}

	for provider := range expectedProviders {
		found := false
		for _, p := range providers {
			if p == provider {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("期望支持 Provider '%s', 但未在列表中找到", provider)
		}
	}
}

// TestIsProviderSupported 测试检查 Provider 是否支持
func TestIsProviderSupported(t *testing.T) {
	tests := []struct {
		provider string
		expected bool
	}{
		{"openai", true},
		{"OpenAI", true}, // 测试大小写不敏感
		{"anthropic", true},
		{"claude", true},
		{"local", true},
		{"ollama", true},
		{"mock", true},
		{"unsupported", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			result := IsProviderSupported(tt.provider)
			if result != tt.expected {
				t.Errorf("IsProviderSupported('%s') = %v, 期望 %v", tt.provider, result, tt.expected)
			}
		})
	}
}

// TestNewProvider_CaseInsensitive 测试大小写不敏感
func TestNewProvider_CaseInsensitive(t *testing.T) {
	tests := []string{"OpenAI", "OPENAI", "openai", "OpEnAi"}

	for _, providerName := range tests {
		t.Run(providerName, func(t *testing.T) {
			cfg := &config.Config{
				LLM: config.LLMConfig{
					Provider: providerName,
					APIKey:   "test-key",
					Model:    "gpt-4",
				},
			}

			provider, err := NewProvider(cfg)
			if err != nil {
				t.Fatalf("创建 Provider 失败: %v", err)
			}

			if provider.Name() != "openai" {
				t.Errorf("期望 Provider 名称为 'openai', 实际为 '%s'", provider.Name())
			}
		})
	}
}
