// Package safety 提供命令安全检查功能
package safety

import "regexp"

// RiskLevel 表示风险等级
type RiskLevel int

const (
	// RiskLow 低风险
	RiskLow RiskLevel = iota

	// RiskMedium 中等风险
	RiskMedium

	// RiskHigh 高风险
	RiskHigh

	// RiskCritical 极高风险
	RiskCritical
)

// Pattern 表示危险命令模式
type Pattern struct {
	// Regex 正则表达式
	Regex *regexp.Regexp

	// Description 描述
	Description string

	// Level 风险等级
	Level RiskLevel
}

// DangerousPatterns 定义内置的危险命令模式
var DangerousPatterns = []Pattern{
	// 文件删除操作
	{
		Regex:       regexp.MustCompile(`\brm\s+-[rf]+\s+/\s*$`),
		Description: "删除根目录文件",
		Level:       RiskCritical,
	},
	{
		Regex:       regexp.MustCompile(`\brm\s+.*-[rf]+`),
		Description: "递归删除文件或目录",
		Level:       RiskHigh,
	},
	{
		Regex:       regexp.MustCompile(`\bdel\s+/[SQF]`),
		Description: "Windows 批量删除命令",
		Level:       RiskHigh,
	},
	{
		Regex:       regexp.MustCompile(`Remove-Item.*-Recurse`),
		Description: "PowerShell 递归删除",
		Level:       RiskHigh,
	},

	// 格式化操作
	{
		Regex:       regexp.MustCompile(`\bmkfs\.[a-z0-9]+`),
		Description: "格式化文件系统",
		Level:       RiskCritical,
	},
	{
		Regex:       regexp.MustCompile(`\bformat\s+[A-Z]:`),
		Description: "Windows 格式化磁盘",
		Level:       RiskCritical,
	},

	// 磁盘操作
	{
		Regex:       regexp.MustCompile(`\bdd\s+.*of=/dev/`),
		Description: "直接写入磁盘设备",
		Level:       RiskCritical,
	},

	// 权限修改
	{
		Regex:       regexp.MustCompile(`\bchmod\s+(777|a\+rwx)`),
		Description: "设置完全开放的文件权限",
		Level:       RiskMedium,
	},
	{
		Regex:       regexp.MustCompile(`\bchown\s+.*:/`),
		Description: "修改根目录所有权",
		Level:       RiskHigh,
	},

	// 网络危险操作
	{
		Regex:       regexp.MustCompile(`(curl|wget)\s+.*\|\s*(ba)?sh`),
		Description: "从网络下载并执行脚本",
		Level:       RiskHigh,
	},
	{
		Regex:       regexp.MustCompile(`\|\s*sudo\s+(ba)?sh`),
		Description: "以管理员权限执行管道输入",
		Level:       RiskHigh,
	},

	// 系统修改
	{
		Regex:       regexp.MustCompile(`\b(sudo|su)\s+(rm|mkfs|format|dd)`),
		Description: "以管理员权限执行危险命令",
		Level:       RiskCritical,
	},
	{
		Regex:       regexp.MustCompile(`>\s*/dev/(sd[a-z]|hd[a-z]|nvme[0-9])`),
		Description: "直接写入磁盘设备",
		Level:       RiskCritical,
	},

	// 批量操作
	{
		Regex:       regexp.MustCompile(`\brm\s+.*\*`),
		Description: "使用通配符删除文件",
		Level:       RiskMedium,
	},

	// 系统关闭/重启
	{
		Regex:       regexp.MustCompile(`\b(shutdown|reboot|halt|poweroff)\b`),
		Description: "系统关闭或重启",
		Level:       RiskMedium,
	},

	// 修改系统配置
	{
		Regex:       regexp.MustCompile(`>\s*/etc/(passwd|shadow|sudoers|hosts)`),
		Description: "修改关键系统配置文件",
		Level:       RiskHigh,
	},

	// 禁用安全功能
	{
		Regex:       regexp.MustCompile(`setenforce\s+0`),
		Description: "禁用 SELinux",
		Level:       RiskHigh,
	},
	{
		Regex:       regexp.MustCompile(`iptables\s+-F`),
		Description: "清空防火墙规则",
		Level:       RiskHigh,
	},

	// fork 炸弹和恶意命令
	{
		Regex:       regexp.MustCompile(`:\(\)\{.*:\|:&\}`),
		Description: "Fork 炸弹",
		Level:       RiskCritical,
	},
}

// String 返回风险等级的字符串表示
func (r RiskLevel) String() string {
	switch r {
	case RiskLow:
		return "低"
	case RiskMedium:
		return "中"
	case RiskHigh:
		return "高"
	case RiskCritical:
		return "极高"
	default:
		return "未知"
	}
}
