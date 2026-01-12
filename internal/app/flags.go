// Package app 提供应用程序的核心逻辑
package app

// Flags 包含所有命令行标志
type Flags struct {
	// DryRun 只显示命令不执行
	DryRun bool

	// Verbose 显示详细输出
	Verbose bool

	// Force 强制执行，跳过确认
	Force bool

	// NoSendStdin 不将 stdin 数据发送到 LLM
	NoSendStdin bool

	// Config 配置文件路径
	Config string

	// History 是否显示历史记录
	History bool

	// Retry 重新执行历史命令的 ID
	Retry int

	// Version 显示版本信息
	Version bool
}

// NewFlags 创建默认的标志配置
func NewFlags() *Flags {
	return &Flags{
		DryRun:      false,
		Verbose:     false,
		Force:       false,
		NoSendStdin: false,
		Config:      "~/.aicli.json",
		History:     false,
		Retry:       -1,
		Version:     false,
	}
}
