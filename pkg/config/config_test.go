package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Version != "1.0" {
		t.Errorf("期望版本为 1.0, 实际为 %s", cfg.Version)
	}

	if cfg.LLM.Provider != "openai" {
		t.Errorf("期望默认 LLM 提供商为 openai, 实际为 %s", cfg.LLM.Provider)
	}

	if cfg.LLM.Timeout != 10 {
		t.Errorf("期望 LLM 超时为 10, 实际为 %d", cfg.LLM.Timeout)
	}

	if cfg.Execution.Timeout != 30 {
		t.Errorf("期望执行超时为 30, 实际为 %d", cfg.Execution.Timeout)
	}

	if !cfg.Safety.EnableChecks {
		t.Error("期望安全检查默认启用")
	}

	if !cfg.History.Enabled {
		t.Error("期望历史记录默认启用")
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	// 加载不存在的文件应该返回默认配置
	cfg, err := Load("/nonexistent/path/config.json")
	if err != nil {
		t.Fatalf("加载不存在的文件应返回默认配置，但返回错误: %v", err)
	}

	if cfg.LLM.Provider != "openai" {
		t.Errorf("期望默认提供商为 openai, 实际为 %s", cfg.LLM.Provider)
	}
}

func TestLoadValidConfig(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.json")

	configContent := `{
		"version": "1.0",
		"llm": {
			"provider": "anthropic",
			"api_key": "test-key",
			"model": "claude-3",
			"timeout": 15
		},
		"execution": {
			"auto_confirm": true,
			"timeout": 60
		}
	}`

	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("创建测试配置文件失败: %v", err)
	}

	// 加载配置
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	if cfg.LLM.Provider != "anthropic" {
		t.Errorf("期望提供商为 anthropic, 实际为 %s", cfg.LLM.Provider)
	}

	if cfg.LLM.APIKey != "test-key" {
		t.Errorf("期望 API 密钥为 test-key, 实际为 %s", cfg.LLM.APIKey)
	}

	if cfg.LLM.Timeout != 15 {
		t.Errorf("期望超时为 15, 实际为 %d", cfg.LLM.Timeout)
	}

	if !cfg.Execution.AutoConfirm {
		t.Error("期望 AutoConfirm 为 true")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	// 创建包含无效 JSON 的临时文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(configPath, []byte("invalid json"), 0600); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 加载应该失败
	_, err := Load(configPath)
	if err == nil {
		t.Error("期望加载无效 JSON 失败，但成功了")
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "save_test.json")

	cfg := Default()
	cfg.LLM.Provider = "local" // 使用本地模型不需要 API 密钥
	cfg.LLM.Model = "test-model"

	// 保存配置
	if err := cfg.Save(configPath); err != nil {
		t.Fatalf("保存配置失败: %v", err)
	}

	// 重新加载并验证
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("重新加载配置失败: %v", err)
	}

	if loaded.LLM.Provider != "local" {
		t.Errorf("期望提供商为 local, 实际为 %s", loaded.LLM.Provider)
	}

	if loaded.LLM.Model != "test-model" {
		t.Errorf("期望模型为 test-model, 实际为 %s", loaded.LLM.Model)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "默认配置应该有效（除了缺少 API 密钥）",
			config:  Default(),
			wantErr: true, // 缺少 API 密钥
		},
		{
			name: "完整配置应该有效",
			config: &Config{
				Version: "1.0",
				LLM: LLMConfig{
					Provider: "openai",
					APIKey:   "test-key",
					Model:    "gpt-4",
					Timeout:  10,
				},
				Execution: ExecutionConfig{
					Timeout: 30,
				},
			},
			wantErr: false,
		},
		{
			name: "本地模型不需要 API 密钥",
			config: &Config{
				Version: "1.0",
				LLM: LLMConfig{
					Provider: "local",
					APIKey:   "",
					Model:    "llama2",
					Timeout:  10,
				},
				Execution: ExecutionConfig{
					Timeout: 30,
				},
			},
			wantErr: false,
		},
		{
			name: "缺少提供商应该无效",
			config: &Config{
				Version: "1.0",
				LLM: LLMConfig{
					Provider: "",
					APIKey:   "test-key",
					Model:    "gpt-4",
					Timeout:  10,
				},
				Execution: ExecutionConfig{
					Timeout: 30,
				},
			},
			wantErr: true,
		},
		{
			name: "缺少模型应该无效",
			config: &Config{
				Version: "1.0",
				LLM: LLMConfig{
					Provider: "openai",
					APIKey:   "test-key",
					Model:    "",
					Timeout:  10,
				},
				Execution: ExecutionConfig{
					Timeout: 30,
				},
			},
			wantErr: true,
		},
		{
			name: "超时为 0 应该无效",
			config: &Config{
				Version: "1.0",
				LLM: LLMConfig{
					Provider: "openai",
					APIKey:   "test-key",
					Model:    "gpt-4",
					Timeout:  0,
				},
				Execution: ExecutionConfig{
					Timeout: 30,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnvironmentVariableOverride(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "env_test.json")

	configContent := `{
		"version": "1.0",
		"llm": {
			"provider": "openai",
			"api_key": "file-key",
			"model": "gpt-4",
			"timeout": 10
		},
		"execution": {
			"timeout": 30
		}
	}`

	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("创建测试配置文件失败: %v", err)
	}

	// 设置环境变量
	os.Setenv("AICLI_API_KEY", "env-key")
	defer os.Unsetenv("AICLI_API_KEY")

	// 加载配置
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 环境变量应该覆盖文件中的值
	if cfg.LLM.APIKey != "env-key" {
		t.Errorf("期望 API 密钥从环境变量覆盖为 env-key, 实际为 %s", cfg.LLM.APIKey)
	}
}

func TestApplyDefaults(t *testing.T) {
	// 创建只有部分字段的配置
	cfg := &Config{
		LLM: LLMConfig{
			Provider: "openai",
			APIKey:   "test-key",
			Model:    "gpt-4",
			// Timeout 和 MaxTokens 未设置
		},
	}

	cfg.applyDefaults()

	// 应该应用默认值
	if cfg.LLM.Timeout != 10 {
		t.Errorf("期望默认超时为 10, 实际为 %d", cfg.LLM.Timeout)
	}

	if cfg.LLM.MaxTokens != 500 {
		t.Errorf("期望默认 MaxTokens 为 500, 实际为 %d", cfg.LLM.MaxTokens)
	}

	if cfg.Execution.Timeout != 30 {
		t.Errorf("期望默认执行超时为 30, 实际为 %d", cfg.Execution.Timeout)
	}
}
