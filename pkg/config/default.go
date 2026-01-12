package config

// Default 返回默认配置
func Default() *Config {
	return &Config{
		Version: "1.0",
		LLM: LLMConfig{
			Provider:  "openai",
			APIKey:    "",
			APIBase:   "https://api.openai.com/v1",
			Model:     "gpt-4",
			Timeout:   10,
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
