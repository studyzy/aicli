// Package i18n 定义所有翻译键常量
package i18n

// 翻译键常量 - 按模块分组

// 错误信息键
const (
	ErrLoadConfig      = "error.load_config"
	ErrCreateProvider  = "error.create_provider"
	ErrTranslateFailed = "error.translate_failed"
	ErrExecuteFailed   = "error.execute_failed"
	ErrNoInput         = "error.no_input"
	ErrEmptyCommand    = "error.empty_command"
	ErrLoadHistory     = "error.load_history"
	ErrSaveHistory     = "error.save_history"
	ErrGetUserHome     = "error.get_user_home"
	ErrPipeModeDanger  = "error.pipe_mode_danger"
	ErrUserCancelled   = "error.user_cancelled"
	ErrHistoryNotFound = "error.history_not_found"
	ErrReadStdin       = "error.read_stdin"
	ErrSaveConfig      = "error.save_config"
)

// 提示信息键
const (
	PromptConfirmRisky    = "prompt.confirm_risky"
	PromptContinue        = "prompt.continue"
	PromptEnterAPIKey     = "prompt.enter_api_key"
	PromptEnterModel      = "prompt.enter_model"
	PromptEnterAPIBase    = "prompt.enter_api_base"
	PromptSelectProvider  = "prompt.select_provider"
	PromptEnableCheck     = "prompt.enable_check"
	PromptEnableHistory   = "prompt.enable_history"
	PromptOverwriteConfig = "prompt.overwrite_config"
	PromptInputChoice     = "prompt.input_choice"
)

// 界面文本键
const (
	MsgHistoryEmpty       = "msg.history_empty"
	MsgHistoryCount       = "msg.history_count"
	MsgNoHistory          = "msg.no_history"
	MsgConfigSaved        = "msg.config_saved"
	MsgConfigNotExist     = "msg.config_not_exist"
	MsgConfigHint         = "msg.config_hint"
	MsgRetryCommand       = "msg.retry_command"
	MsgWelcomeInit        = "msg.welcome_init"
	MsgInitGuide          = "msg.init_guide"
	MsgSavingConfig       = "msg.saving_config"
	MsgNowCanUse          = "msg.now_can_use"
	MsgExampleUsage       = "msg.example_usage"
	MsgOperationCancelled = "msg.operation_cancelled"
	MsgOtherSettings      = "msg.other_settings"
	MsgConfigExists       = "msg.config_exists"
	MsgInvalidChoice      = "msg.invalid_choice"
	MsgDefaultUseOpenAI   = "msg.default_use_openai"
	MsgDefault            = "msg.default"
	MsgTranslatedCommand  = "msg.translated_command"
	MsgTrialAPINotice     = "msg.trial_api_notice" // 试用 API 提示
)

// 警告信息键
const (
	WarnDangerousCommand = "warn.dangerous_command"
	WarnAPIKeyEmpty      = "warn.api_key_empty"
	WarnRisk             = "warn.risk"
	WarnRiskLevel        = "warn.risk_level"
)

// 字段标签键
const (
	LabelCommand     = "label.command"
	LabelInput       = "label.input"
	LabelError       = "label.error"
	LabelOutput      = "label.output"
	LabelTimestamp   = "label.timestamp"
	LabelRisk        = "label.risk"
	LabelLevel       = "label.level"
	LabelOS          = "label.os"
	LabelShell       = "label.shell"
	LabelWorkDir     = "label.workdir"
	LabelStdin       = "label.stdin"
	LabelStdinBytes  = "label.stdin_bytes"
	LabelProvider    = "label.provider"
	LabelModel       = "label.model"
	LabelAPIBase     = "label.api_base"
	LabelAPIKey      = "label.api_key"
)

// Verbose 模式信息键
const (
	VerboseInput            = "verbose.input"
	VerboseStdin            = "verbose.stdin"
	VerboseContext          = "verbose.context"
	VerboseCommand          = "verbose.command"
	VerboseTranslateTime    = "verbose.translate_time"
	VerboseExecuting        = "verbose.executing"
	VerboseExecuteTime      = "verbose.execute_time"
	VerboseTotalTime        = "verbose.total_time"
	VerboseConfigNotExist   = "verbose.config_not_exist"
	VerboseLoadHistoryFailed = "verbose.load_history_failed"
	VerboseSaveHistoryFailed = "verbose.save_history_failed"
)

// Dry-run 模式键
const (
	DryRunWillExecute = "dry_run.will_execute"
)

// LLM 提示词键
const (
	LLMSystemPromptIntro      = "llm.system_prompt_intro"
	LLMSystemPromptRules      = "llm.system_prompt_rules"
	LLMSystemPromptRule1      = "llm.system_prompt_rule1"
	LLMSystemPromptRule2      = "llm.system_prompt_rule2"
	LLMSystemPromptRule3      = "llm.system_prompt_rule3"
	LLMSystemPromptRule4      = "llm.system_prompt_rule4"
	LLMSystemPromptRule5      = "llm.system_prompt_rule5"
	LLMSystemPromptEnv        = "llm.system_prompt_env"
	LLMUserPromptIntro        = "llm.user_prompt_intro"
	LLMStdinData              = "llm.stdin_data"
	LLMTruncated              = "llm.truncated"
	LLMContextNoContext       = "llm.context_no_context"
	LLMContextFormat          = "llm.context_format"
)

// Cobra 命令描述键
const (
	CobraUse             = "cobra.use"
	CobraShort           = "cobra.short"
	CobraLong            = "cobra.long"
	CobraFlagConfig      = "cobra.flag_config"
	CobraFlagVerbose     = "cobra.flag_verbose"
	CobraFlagDryRun      = "cobra.flag_dry_run"
	CobraFlagForce       = "cobra.flag_force"
	CobraFlagNoSendStdin = "cobra.flag_no_send_stdin"
	CobraFlagHistory     = "cobra.flag_history"
	CobraFlagRetry       = "cobra.flag_retry"
	CobraFlagQuiet       = "cobra.flag_quiet"
)

// Init 命令键
const (
	InitUse   = "init.use"
	InitShort = "init.short"
	InitLong  = "init.long"
)

// Completion 命令键
const (
	CompletionShort = "completion.short"
)

// Help 命令键
const (
	HelpShort = "help.short"
)

// Version flag 键
const (
	VersionShort = "version.short"
)

// Help flag 键
const (
	HelpFlag = "help.flag"
)

// 配置向导选项键
const (
	InitProviderOpenAI    = "init.provider_openai"
	InitProviderAnthropic = "init.provider_anthropic"
	InitProviderLocal     = "init.provider_local"
	InitProviderDeepSeek  = "init.provider_deepseek"
	InitProviderOther     = "init.provider_other"
)
