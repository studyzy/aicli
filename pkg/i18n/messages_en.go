// Package i18n English translation resources
package i18n

// messagesEn English translation map
var messagesEn = map[string]string{
	// Error messages
	ErrLoadConfig:      "Failed to load configuration",
	ErrCreateProvider:  "Failed to create LLM Provider",
	ErrTranslateFailed: "Failed to translate command",
	ErrExecuteFailed:   "Failed to execute command",
	ErrNoInput:         "Please provide natural language description",
	ErrEmptyCommand:    "LLM returned empty command",
	ErrLoadHistory:     "Failed to load history",
	ErrSaveHistory:     "Failed to save history",
	ErrGetUserHome:     "Failed to get user home directory",
	ErrPipeModeDanger:  "Refusing to execute dangerous command in pipe mode (use --force to override)",
	ErrUserCancelled:   "User cancelled dangerous command execution",
	ErrHistoryNotFound: "History record not found",
	ErrReadStdin:       "Failed to read stdin",
	ErrSaveConfig:      "Failed to save configuration",

	// Prompts
	PromptConfirmRisky:    "Continue execution? (y/n): ",
	PromptContinue:        "Continue?",
	PromptEnterAPIKey:     "Please enter API Key",
	PromptEnterModel:      "Please enter model name",
	PromptEnterAPIBase:    "Please enter API Base URL",
	PromptSelectProvider:  "Please select LLM provider",
	PromptEnableCheck:     "Enable dangerous command safety checks?",
	PromptEnableHistory:   "Enable history recording?",
	PromptOverwriteConfig: "Overwrite?",
	PromptInputChoice:     "Please enter your choice",

	// UI messages
	MsgHistoryEmpty:       "No history records",
	MsgHistoryCount:       "History (%d entries):",
	MsgNoHistory:          "No history records",
	MsgConfigSaved:        "Configuration successfully saved to %s",
	MsgConfigNotExist:     "Note: Configuration file %s does not exist.",
	MsgConfigHint:         "You can run 'aicli init' for quick setup.",
	MsgRetryCommand:       "Retrying history command #%d:",
	MsgWelcomeInit:        "Welcome to aicli configuration wizard!",
	MsgInitGuide:          "We will guide you through the basic configuration.",
	MsgSavingConfig:       "Saving configuration...",
	MsgNowCanUse:          "You can now start using aicli!",
	MsgExampleUsage:       "Example: aicli \"check my public IP\"",
	MsgOperationCancelled: "Operation cancelled.",
	MsgOtherSettings:      "--- Other Settings ---",
	MsgConfigExists:       "Configuration file %s already exists.",
	MsgInvalidChoice:      "Invalid choice, defaulting to OpenAI",
	MsgDefaultUseOpenAI:   "Invalid choice, defaulting to OpenAI",
	MsgDefault:            "default",

	// Warnings
	WarnDangerousCommand: "Potentially dangerous command detected!",
	WarnAPIKeyEmpty:      "Warning: API Key is empty. You may need to set it in the AICLI_API_KEY environment variable.",
	WarnRisk:             "Risk",
	WarnRiskLevel:        "Level",

	// Field labels
	LabelCommand:    "Command",
	LabelInput:      "Input",
	LabelError:      "Error",
	LabelOutput:     "Output",
	LabelTimestamp:  "Timestamp",
	LabelRisk:       "Risk",
	LabelLevel:      "Level",
	LabelOS:         "Operating System",
	LabelShell:      "Shell",
	LabelWorkDir:    "Working Directory",
	LabelStdin:      "Standard Input",
	LabelStdinBytes: "bytes",
	LabelProvider:   "Provider",
	LabelModel:      "Model",
	LabelAPIBase:    "API Base URL",
	LabelAPIKey:     "API Key",

	// Verbose mode
	VerboseInput:             "Natural language input",
	VerboseStdin:             "Standard input",
	VerboseContext:           "Execution context",
	VerboseCommand:           "Translated command",
	VerboseTranslateTime:     "Translation time",
	VerboseExecuting:         "Executing command...",
	VerboseExecuteTime:       "Execution time",
	VerboseTotalTime:         "Total time",
	VerboseConfigNotExist:    "Configuration file does not exist, using default configuration",
	VerboseLoadHistoryFailed: "Failed to load history: %v",
	VerboseSaveHistoryFailed: "Failed to save history: %v",

	// Dry-run mode
	DryRunWillExecute: "Command to be executed: %s",

	// LLM prompts
	LLMSystemPromptIntro: "You are a command-line assistant that converts natural language descriptions into executable shell commands.",
	LLMSystemPromptRules: "Rules:",
	LLMSystemPromptRule1: "1. Return only the command itself, without any explanation or description",
	LLMSystemPromptRule2: "2. Do not use markdown code block format",
	LLMSystemPromptRule3: "3. The command must be directly executable",
	LLMSystemPromptRule4: "4. If multiple commands are needed, connect them with && or ;",
	LLMSystemPromptRule5: "5. Prefer commonly used and compatible commands",
	LLMSystemPromptEnv:   "Execution Environment:",
	LLMUserPromptIntro:   "Convert the following natural language description into a command:",
	LLMStdinData:         "Standard input data:",
	LLMTruncated:         "... (truncated)",
	LLMContextNoContext:  "No execution context",
	LLMContextFormat:     "OS: %s, Shell: %s, WorkDir: %s",

	// Cobra command descriptions
	CobraUse:   "aicli [natural language description]",
	CobraShort: "AI command-line assistant",
	CobraLong: `aicli is a tool that brings natural language operations to the command line.
Users can simply enter natural language, and it will be converted into actual commands through LLM services and executed.

Examples:
  aicli "find ERROR logs in log.txt"
  aicli "show all txt files in current directory"
  cat file.txt | aicli "count lines"
  aicli --history
  aicli --retry 3`,

	CobraFlagConfig:      "Configuration file path",
	CobraFlagVerbose:     "Show detailed output",
	CobraFlagDryRun:      "Show command without executing",
	CobraFlagForce:       "Force execution, skip confirmation",
	CobraFlagNoSendStdin: "Do not send stdin data to LLM",
	CobraFlagHistory:     "Show history records",
	CobraFlagRetry:       "Retry history command ID",

	// Init command
	InitUse:   "init",
	InitShort: "Initialize configuration",
	InitLong:  "Guide user to set up LLM configuration and generate configuration file ~/.aicli.json",

	// Completion command
	CompletionShort: "Generate the autocompletion script for the specified shell",

	// Help command
	HelpShort: "Help about any command",

	// Version flag
	VersionShort: "version for aicli",

	// Help flag
	HelpFlag: "help for aicli",

	// Configuration wizard options
	InitProviderOpenAI:    "1. OpenAI (GPT-4, GPT-3.5)",
	InitProviderAnthropic: "2. Anthropic (Claude)",
	InitProviderLocal:     "3. Local (Ollama, LocalAI)",
	InitProviderDeepSeek:  "4. DeepSeek",
	InitProviderOther:     "5. Other (OpenAI-compatible API)",
}
