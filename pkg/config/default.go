package config

// Default 返回默认配置
func Default() *Config {
	return &Config{
		Version: "1.0",
		LLM: LLMConfig{
			Provider:  "builtin", // 使用内置试用 API
			APIKey:    "",
			APIBase:   "",
			Model:     "",
			Timeout:   30, // 增加超时时间以适应网络延迟
			MaxTokens: 500,
		},
		Execution: ExecutionConfig{
			AutoConfirm:   false,
			DryRunDefault: false,
			Timeout:       30,
			Shell:         "auto",
		},
		Safety: SafetyConfig{
			EnableChecks:        true,
			DangerousPatterns:   []string{},
			RequireConfirmation: true,
		},
		History: HistoryConfig{
			Enabled:    true,
			MaxEntries: 1000,
			File:       "~/.aicli_history.json",
		},
		Logging: LoggingConfig{
			Enabled: false,
			Level:   "info",
			File:    "",
		},
	}
}

// DefaultConfigPath 返回默认配置文件路径
func DefaultConfigPath() string {
	return "~/.aicli.json"
}
