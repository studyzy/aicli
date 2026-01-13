// Package config 提供配置管理功能
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	providerLocal = "local"
)

// Config 是应用程序的主配置结构体
type Config struct {
	Version   string          `json:"version"`
	LLM       LLMConfig       `json:"llm"`
	Execution ExecutionConfig `json:"execution"`
	Safety    SafetyConfig    `json:"safety"`
	History   HistoryConfig   `json:"history"`
	Logging   LoggingConfig   `json:"logging"`
}

// LLMConfig 包含 LLM 服务的配置
type LLMConfig struct {
	Provider  string `json:"provider"`   // LLM 提供商名称 (openai, anthropic, local)
	APIKey    string `json:"api_key"`    // API 密钥
	APIBase   string `json:"api_base"`   // API 基础 URL
	Model     string `json:"model"`      // 模型名称
	Timeout   int    `json:"timeout"`    // 超时时间（秒）
	MaxTokens int    `json:"max_tokens"` // 最大 token 数
}

// ExecutionConfig 包含命令执行的配置
type ExecutionConfig struct {
	AutoConfirm   bool   `json:"auto_confirm"`    // 是否自动确认命令
	DryRunDefault bool   `json:"dry_run_default"` // 默认是否只显示命令不执行
	Timeout       int    `json:"timeout"`         // 命令执行超时（秒）
	Shell         string `json:"shell"`           // Shell 类型 (auto, bash, zsh, powershell, cmd)
}

// SafetyConfig 包含安全检查的配置
type SafetyConfig struct {
	EnableChecks        bool     `json:"enable_checks"`        // 是否启用安全检查
	DangerousPatterns   []string `json:"dangerous_patterns"`   // 额外的危险模式
	RequireConfirmation bool     `json:"require_confirmation"` // 是否需要确认
}

// HistoryConfig 包含历史记录的配置
type HistoryConfig struct {
	Enabled    bool   `json:"enabled"`     // 是否启用历史记录
	MaxEntries int    `json:"max_entries"` // 最大保存条目数
	File       string `json:"file"`        // 历史记录文件路径
}

// LoggingConfig 包含日志的配置
type LoggingConfig struct {
	Enabled bool   `json:"enabled"` // 是否启用日志
	Level   string `json:"level"`   // 日志级别 (debug, info, warn, error)
	File    string `json:"file"`    // 日志文件路径（空表示标准输出）
}

// Load 从指定路径加载配置文件
// 如果文件不存在，返回默认配置
func Load(path string) (*Config, error) {
	// 展开 ~ 路径
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("获取用户主目录失败: %w", err)
		}
		path = filepath.Join(homeDir, path[2:])
	}

	// 如果文件不存在，返回默认配置
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Default(), nil
	}

	// 读取文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析 JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 合并默认值
	config.applyDefaults()

	// 从环境变量覆盖 API 密钥
	if apiKey := os.Getenv("AICLI_API_KEY"); apiKey != "" {
		config.LLM.APIKey = apiKey
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// Save 保存配置到指定路径
func (c *Config) Save(path string) error {
	// 展开 ~ 路径
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户主目录失败: %w", err)
		}
		path = filepath.Join(homeDir, path[2:])
	}

	// 序列化为 JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// Validate 验证配置的有效性
func (c *Config) Validate() error {
	// 验证 LLM 配置
	if c.LLM.Provider == "" {
		return fmt.Errorf("LLM 提供商不能为空")
	}

	// 对于非本地模型，需要 API 密钥
	if c.LLM.Provider != providerLocal && c.LLM.APIKey == "" {
		return fmt.Errorf("LLM API 密钥不能为空（或设置环境变量 AICLI_API_KEY）")
	}

	if c.LLM.Model == "" {
		return fmt.Errorf("LLM 模型不能为空")
	}

	// 验证超时配置
	if c.LLM.Timeout <= 0 {
		return fmt.Errorf("LLM 超时时间必须大于 0")
	}

	if c.Execution.Timeout <= 0 {
		return fmt.Errorf("命令执行超时时间必须大于 0")
	}

	return nil
}

// applyDefaults 应用默认值到未设置的字段
func (c *Config) applyDefaults() {
	defaults := Default()

	if c.Version == "" {
		c.Version = defaults.Version
	}

	// LLM 默认值
	if c.LLM.Timeout == 0 {
		c.LLM.Timeout = defaults.LLM.Timeout
	}
	if c.LLM.MaxTokens == 0 {
		c.LLM.MaxTokens = defaults.LLM.MaxTokens
	}
	if c.LLM.APIBase == "" {
		c.LLM.APIBase = defaults.LLM.APIBase
	}

	// Execution 默认值
	if c.Execution.Timeout == 0 {
		c.Execution.Timeout = defaults.Execution.Timeout
	}
	if c.Execution.Shell == "" {
		c.Execution.Shell = defaults.Execution.Shell
	}

	// History 默认值
	if c.History.MaxEntries == 0 {
		c.History.MaxEntries = defaults.History.MaxEntries
	}
	if c.History.File == "" {
		c.History.File = defaults.History.File
	}

	// Logging 默认值
	if c.Logging.Level == "" {
		c.Logging.Level = defaults.Logging.Level
	}
}
