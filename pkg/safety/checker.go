package safety

import "strings"

// Checker 提供命令安全检查功能
type Checker struct {
	patterns       []Pattern
	customPatterns []Pattern
	enableChecks   bool
}

// NewChecker 创建新的安全检查器
func NewChecker(enableChecks bool) *Checker {
	return &Checker{
		patterns:       DangerousPatterns,
		customPatterns: []Pattern{},
		enableChecks:   enableChecks,
	}
}

// AddCustomPattern 添加自定义危险模式
func (c *Checker) AddCustomPattern(pattern Pattern) {
	c.customPatterns = append(c.customPatterns, pattern)
}

// IsDangerous 检查命令是否危险
// 返回: 是否危险, 匹配的模式描述, 风险等级
func (c *Checker) IsDangerous(command string) (bool, string, RiskLevel) {
	if !c.enableChecks {
		return false, "", RiskLow
	}

	// 清理命令字符串
	command = strings.TrimSpace(command)
	if command == "" {
		return false, "", RiskLow
	}

	// 检查内置模式
	for _, pattern := range c.patterns {
		if pattern.Regex.MatchString(command) {
			return true, pattern.Description, pattern.Level
		}
	}

	// 检查自定义模式
	for _, pattern := range c.customPatterns {
		if pattern.Regex.MatchString(command) {
			return true, pattern.Description, pattern.Level
		}
	}

	return false, "", RiskLow
}

// CheckMultiple 检查多个命令（用于管道或命令链）
func (c *Checker) CheckMultiple(commands []string) (bool, []string, RiskLevel) {
	if !c.enableChecks {
		return false, nil, RiskLow
	}

	var dangerousCommands []string
	maxRiskLevel := RiskLow

	for _, cmd := range commands {
		if isDangerous, desc, level := c.IsDangerous(cmd); isDangerous {
			dangerousCommands = append(dangerousCommands, desc)
			if level > maxRiskLevel {
				maxRiskLevel = level
			}
		}
	}

	if len(dangerousCommands) > 0 {
		return true, dangerousCommands, maxRiskLevel
	}

	return false, nil, RiskLow
}

// Enable 启用安全检查
func (c *Checker) Enable() {
	c.enableChecks = true
}

// Disable 禁用安全检查
func (c *Checker) Disable() {
	c.enableChecks = false
}

// IsEnabled 返回是否启用安全检查
func (c *Checker) IsEnabled() bool {
	return c.enableChecks
}
