// Package i18n 中文翻译资源
package i18n

// messagesZh 中文翻译映射表
var messagesZh = map[string]string{
	// 错误信息
	ErrLoadConfig:      "加载配置失败",
	ErrCreateProvider:  "创建 LLM Provider 失败",
	ErrTranslateFailed: "命令转换失败",
	ErrExecuteFailed:   "命令执行失败",
	ErrNoInput:         "请提供自然语言描述",
	ErrEmptyCommand:    "LLM 返回空命令",
	ErrLoadHistory:     "加载历史记录失败",
	ErrSaveHistory:     "保存历史记录失败",
	ErrGetUserHome:     "获取用户主目录失败",
	ErrPipeModeDanger:  "管道模式下拒绝执行危险命令(使用 --force 强制执行)",
	ErrUserCancelled:   "用户取消执行危险命令",
	ErrHistoryNotFound: "历史记录不存在",
	ErrReadStdin:       "读取 stdin 失败",
	ErrSaveConfig:      "保存配置失败",

	// 提示信息
	PromptConfirmRisky:    "是否继续执行?(y/n): ",
	PromptContinue:        "是否继续?",
	PromptEnterAPIKey:     "请输入 API Key",
	PromptEnterModel:      "请输入模型名称",
	PromptEnterAPIBase:    "请输入 API Base URL",
	PromptSelectProvider:  "请选择 LLM 提供商",
	PromptEnableCheck:     "是否启用危险命令安全检查?",
	PromptEnableHistory:   "是否启用历史记录?",
	PromptOverwriteConfig: "是否覆盖?",
	PromptInputChoice:     "请输入序号",

	// 界面文本
	MsgHistoryEmpty:       "没有历史记录",
	MsgHistoryCount:       "历史记录(共 %d 条):",
	MsgNoHistory:          "没有历史记录",
	MsgConfigSaved:        "配置已成功保存到 %s",
	MsgConfigNotExist:     "提示: 配置文件 %s 不存在。",
	MsgConfigHint:         "您可以运行 'aicli init' 来快速设置配置。",
	MsgRetryCommand:       "重新执行历史命令 #%d:",
	MsgWelcomeInit:        "欢迎使用 aicli 配置向导!",
	MsgInitGuide:          "我们将引导您完成基本配置。",
	MsgSavingConfig:       "正在保存配置...",
	MsgNowCanUse:          "现在您可以开始使用 aicli 了!",
	MsgExampleUsage:       "示例: aicli \"查询我的公网IP\"",
	MsgOperationCancelled: "操作已取消。",
	MsgOtherSettings:      "--- 其他设置 ---",
	MsgConfigExists:       "配置文件 %s 已存在。",
	MsgInvalidChoice:      "无效的选择,默认使用 OpenAI",
	MsgDefaultUseOpenAI:   "无效的选择,默认使用 OpenAI",
	MsgDefault:            "默认",

	// 警告信息
	WarnDangerousCommand: "检测到潜在危险命令!",
	WarnAPIKeyEmpty:      "警告: API Key 为空,您可能需要在环境变量 AICLI_API_KEY 中设置。",
	WarnRisk:             "风险",
	WarnRiskLevel:        "等级",

	// 字段标签
	LabelCommand:    "命令",
	LabelInput:      "输入",
	LabelError:      "错误",
	LabelOutput:     "输出",
	LabelTimestamp:  "时间",
	LabelRisk:       "风险",
	LabelLevel:      "等级",
	LabelOS:         "操作系统",
	LabelShell:      "Shell",
	LabelWorkDir:    "工作目录",
	LabelStdin:      "标准输入",
	LabelStdinBytes: "字节",
	LabelProvider:   "提供商",
	LabelModel:      "模型",
	LabelAPIBase:    "API Base URL",
	LabelAPIKey:     "API Key",

	// Verbose 模式信息
	VerboseInput:             "自然语言输入",
	VerboseStdin:             "标准输入",
	VerboseContext:           "执行上下文",
	VerboseCommand:           "转换后的命令",
	VerboseTranslateTime:     "转换耗时",
	VerboseExecuting:         "开始执行命令...",
	VerboseExecuteTime:       "执行耗时",
	VerboseTotalTime:         "总耗时",
	VerboseConfigNotExist:    "配置文件不存在,使用默认配置",
	VerboseLoadHistoryFailed: "加载历史记录失败: %v",
	VerboseSaveHistoryFailed: "保存历史记录失败: %v",

	// Dry-run 模式
	DryRunWillExecute: "将要执行的命令: %s",

	// LLM 提示词
	LLMSystemPromptIntro: "你是一个命令行助手,专门将用户的自然语言描述转换为可执行的 shell 命令。",
	LLMSystemPromptRules: "规则:",
	LLMSystemPromptRule1: "1. 只返回命令本身,不要有任何解释或说明",
	LLMSystemPromptRule2: "2. 不要使用 markdown 代码块格式",
	LLMSystemPromptRule3: "3. 命令必须是可以直接执行的",
	LLMSystemPromptRule4: "4. 如果需要多个命令,使用 && 或 ; 连接",
	LLMSystemPromptRule5: "5. 优先使用常见且兼容性好的命令",
	LLMSystemPromptEnv:   "执行环境:",
	LLMUserPromptIntro:   "将以下自然语言描述转换为命令:",
	LLMStdinData:         "标准输入数据:",
	LLMTruncated:         "... (已截断)",
	LLMContextNoContext:  "无执行上下文",
	LLMContextFormat:     "OS: %s, Shell: %s, 工作目录: %s",

	// Cobra 命令描述
	CobraUse:   "aicli [自然语言描述]",
	CobraShort: "AI 命令行助手",
	CobraLong: `aicli 是一个让命令行支持自然语言操作的工具。
用户只需要输入自然语言,就能通过 LLM 服务将其转换为实际的命令并执行。

示例:
  aicli "查找log.txt中的ERROR日志"
  aicli "显示当前目录下所有txt文件"
  cat file.txt | aicli "统计行数"
  aicli --history
  aicli --retry 3`,

	CobraFlagConfig:      "配置文件路径",
	CobraFlagVerbose:     "显示详细输出",
	CobraFlagDryRun:      "仅显示命令不执行",
	CobraFlagForce:       "强制执行,跳过确认",
	CobraFlagNoSendStdin: "不将 stdin 数据发送到 LLM",
	CobraFlagHistory:     "显示历史记录",
	CobraFlagRetry:       "重新执行历史命令 ID",

	// Init 命令
	InitUse:   "init",
	InitShort: "初始化配置",
	InitLong:  "引导用户设置 LLM 配置并生成配置文件 ~/.aicli.json",

	// Completion 命令
	CompletionShort: "为指定的 shell 生成自动补全脚本",

	// Help 命令
	HelpShort: "显示任何命令的帮助信息",

	// Version flag
	VersionShort: "显示版本信息",

	// Help flag
	HelpFlag: "显示帮助信息",

	// 配置向导选项
	InitProviderOpenAI:    "1. OpenAI (GPT-4, GPT-3.5)",
	InitProviderAnthropic: "2. Anthropic (Claude)",
	InitProviderLocal:     "3. Local (Ollama, LocalAI)",
	InitProviderDeepSeek:  "4. DeepSeek (深度求索)",
	InitProviderOther:     "5. Other (兼容 OpenAI 协议)",
}
