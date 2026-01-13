package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/studyzy/aicli/internal/app"
	"github.com/studyzy/aicli/internal/history"
	"github.com/studyzy/aicli/pkg/config"
	"github.com/studyzy/aicli/pkg/executor"
	"github.com/studyzy/aicli/pkg/i18n"
	"github.com/studyzy/aicli/pkg/llm"
	"github.com/studyzy/aicli/pkg/safety"
)

// version 是当前版本号
const version = "1.0.0"

var (
	flags = app.NewFlags()
)

func main() {
	// 早期初始化i18n(用于--help等不进入run的场景)
	// 使用环境变量检测,配置文件会在run中重新初始化
	i18n.Init(nil) // nil表示仅使用环境变量和默认值
	
	// 设置自定义 Help 函数，在显示帮助前更新命令描述
	// 这样可以确保 Cobra 自动生成的命令（如 completion、help）也能被国际化
	originalHelpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		// 在显示帮助前更新所有命令的描述
		updateCommandDescriptions(rootCmd)
		originalHelpFunc(cmd, args)
	})
	
	// 更新rootCmd描述（第一次，在子命令添加之前）
	updateCommandDescriptions(rootCmd)
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", i18n.T("label.error"), err)
		os.Exit(1)
	}
}

// rootCmd 是根命令
var rootCmd = &cobra.Command{
	Use:   "aicli [自然语言描述]",
	Short: "AI 命令行助手",
	Long: `aicli 是一个让命令行支持自然语言操作的工具。
用户只需要输入自然语言，就能通过 LLM 服务将其转换为实际的命令并执行。

示例:
  aicli "查找log.txt中的ERROR日志"
  aicli "显示当前目录下所有txt文件"
  cat file.txt | aicli "统计行数"
  aicli --history
  aicli --retry 3`,
	Version:               version,
	RunE:                  run,
	Args:                  cobra.ArbitraryArgs,
	DisableFlagParsing:    false,
	FParseErrWhitelist:    cobra.FParseErrWhitelist{UnknownFlags: false},
	SilenceUsage:          false,
	DisableSuggestions:    false,
	SuggestionsMinimumDistance: 2,
}

func run(cmd *cobra.Command, args []string) error {
	// 加载配置(用于语言检测)
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrLoadConfig), err)
	}

	// 初始化 i18n
	i18n.Init(cfg)

	// 更新命令描述为对应语言
	updateCommandDescriptions(cmd)

	// 如果没有参数，显示帮助
	if len(args) == 0 && !flags.History && flags.Retry < 0 {
		return cmd.Help()
	}

	// 历史记录功能
	if flags.History {
		return showHistory()
	}

	// 重试功能
	if flags.Retry >= 0 {
		return retryCommand(flags.Retry)
	}

	// 创建 LLM Provider
	provider, err := createLLMProvider(cfg)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrCreateProvider), err)
	}

	// 创建 Executor
	exec := executor.NewExecutor()

	// 创建 Safety Checker
	checker := safety.NewChecker(cfg.Safety.EnableChecks)

	// 创建应用实例
	application := app.NewApp(cfg, provider, exec, checker)

	// 加载历史记录
	hist := history.NewHistory()
	historyPath := getHistoryPath()
	if loadErr := hist.Load(historyPath); loadErr != nil {
		if flags.Verbose {
			msg := i18n.T(i18n.VerboseLoadHistoryFailed, loadErr)
			fmt.Fprintf(os.Stderr, "%s\n", msg)
		}
	}
	application.SetHistory(hist)

	// 获取自然语言输入
	input := strings.Join(args, " ")

	// 读取 stdin（如果有）
	stdin := ""
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// stdin 是管道或重定向
		stdinBytes, readErr := io.ReadAll(os.Stdin)
		if readErr != nil {
			return fmt.Errorf("%s: %w", i18n.T(i18n.ErrReadStdin), readErr)
		}
		stdin = string(stdinBytes)
	}

	// 执行应用逻辑
	output, err := application.Run(input, stdin, flags)

	// 保存历史记录（即使执行失败也保存）
	if saveErr := hist.Save(historyPath); saveErr != nil && flags.Verbose {
		msg := i18n.T(i18n.VerboseSaveHistoryFailed, saveErr)
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}

	if err != nil {
		return err
	}

	// 输出结果
	fmt.Print(output)

	return nil
}

// loadConfig 加载配置文件
func loadConfig() (*config.Config, error) {
	// 如果指定了配置文件，使用指定的
	if flags.Config != "" {
		return config.Load(flags.Config)
	}

	// 否则尝试从默认位置加载
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory / 获取用户主目录失败: %w", err)
	}

	configPath := homeDir + "/.aicli.json"

	// 如果配置文件不存在，使用默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 提示用户初始化 (这里还未初始化i18n,使用双语提示)
		fmt.Fprintf(os.Stderr, "Note / 提示: Configuration file / 配置文件 %s does not exist / 不存在。\n", configPath)
		fmt.Fprintf(os.Stderr, "You can run / 您可以运行 'aicli init' for quick setup / 来快速设置配置。\n\n")

		if flags.Verbose {
			fmt.Fprintf(os.Stderr, "Using default configuration / 使用默认配置\n")
		}
		return config.Default(), nil
	}

	return config.Load(configPath)
}

// createLLMProvider 创建 LLM Provider
// 使用工厂函数统一管理 Provider 创建
func createLLMProvider(cfg *config.Config) (llm.Provider, error) {
	return llm.NewProvider(cfg)
}

func init() {
	// 持久化标志（所有子命令都可用）
	rootCmd.PersistentFlags().StringVarP(&flags.Config, "config", "c", flags.Config, "配置文件路径")
	rootCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", flags.Verbose, "显示详细输出")

	// 本地标志（仅根命令可用）
	rootCmd.Flags().BoolVarP(&flags.DryRun, "dry-run", "n", flags.DryRun, "仅显示命令不执行")
	rootCmd.Flags().BoolVarP(&flags.Force, "force", "f", flags.Force, "强制执行，跳过确认")
	rootCmd.Flags().BoolVarP(&flags.Quiet, "quiet", "q", flags.Quiet, "静默模式，不显示翻译后的命令")
	rootCmd.Flags().BoolVar(&flags.NoSendStdin, "no-send-stdin", flags.NoSendStdin, "不将 stdin 数据发送到 LLM")
	rootCmd.Flags().BoolVar(&flags.History, "history", flags.History, "显示历史记录")
	rootCmd.Flags().IntVar(&flags.Retry, "retry", flags.Retry, "重新执行历史命令 ID")

	// 设置版本模板
	rootCmd.SetVersionTemplate(`{{printf "aicli version %s\n" .Version}}`)
	
	// 添加 completion 命令（Cobra 不会自动添加，需要手动添加）
	rootCmd.AddCommand(createCompletionCmd())
}

// createCompletionCmd 创建 completion 命令
func createCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "", // 将在 updateCommandDescriptions 中设置
		Long:  "", // 将在 updateCommandDescriptions 中设置
	}
}

// getHistoryPath 获取历史记录文件路径
func getHistoryPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".aicli_history.json" // 降级到当前目录
	}
	return homeDir + "/.aicli_history.json"
}

// showHistory 显示历史记录
func showHistory() error {
	hist := history.NewHistory()
	historyPath := getHistoryPath()

	if err := hist.Load(historyPath); err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrLoadHistory), err)
	}

	entries := hist.List()
	if len(entries) == 0 {
		fmt.Println(i18n.T(i18n.MsgNoHistory))
		return nil
	}

	fmt.Println(i18n.T(i18n.MsgHistoryCount, len(entries)) + "\n")
	for _, entry := range entries {
		status := "✓"
		if !entry.Success {
			status = "✗"
		}

		fmt.Printf("[%d] %s %s\n", entry.ID, status, entry.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("    %s: %s\n", i18n.T(i18n.LabelInput), entry.Input)
		fmt.Printf("    %s: %s\n", i18n.T(i18n.LabelCommand), entry.Command)

		if entry.Error != "" {
			fmt.Printf("    %s: %s\n", i18n.T(i18n.LabelError), entry.Error)
		}

		fmt.Println()
	}

	return nil
}

// retryCommand 重新执行历史命令
func retryCommand(id int) error {
	hist := history.NewHistory()
	historyPath := getHistoryPath()

	if err := hist.Load(historyPath); err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrLoadHistory), err)
	}

	entry, err := hist.Get(id)
	if err != nil {
		return err
	}

	msg := i18n.T(i18n.MsgRetryCommand, id)
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	fmt.Fprintf(os.Stderr, "  %s: %s\n", i18n.T(i18n.LabelInput), entry.Input)
	fmt.Fprintf(os.Stderr, "  %s: %s\n\n", i18n.T(i18n.LabelCommand), entry.Command)

	// 加载配置
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrLoadConfig), err)
	}

	// 创建 LLM Provider
	provider, err := createLLMProvider(cfg)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrCreateProvider), err)
	}

	// 创建 Executor
	exec := executor.NewExecutor()

	// 创建 Safety Checker
	checker := safety.NewChecker(cfg.Safety.EnableChecks)

	// 创建应用实例
	application := app.NewApp(cfg, provider, exec, checker)
	application.SetHistory(hist)

	// 执行命令（使用原始输入重新转换）
	output, err := application.Run(entry.Input, "", flags)

	// 保存历史记录
	if saveErr := hist.Save(historyPath); saveErr != nil && flags.Verbose {
		fmt.Fprintf(os.Stderr, "保存历史记录失败: %v\n", saveErr)
	}

	if err != nil {
		return err
	}

	// 输出结果
	fmt.Print(output)

	return nil
}

// updateCommandDescriptions 更新命令描述为对应语言
func updateCommandDescriptions(cmd *cobra.Command) {
	// 更新根命令描述（不更新 Use，因为会导致 Cobra 把它当成子命令）
	cmd.Short = i18n.T(i18n.CobraShort)
	cmd.Long = i18n.T(i18n.CobraLong)
	
	// 更新子命令描述（包括 Cobra 自动生成的命令）
	for _, subCmd := range cmd.Commands() {
		switch subCmd.Name() {
		case "init":
			subCmd.Use = i18n.T(i18n.InitUse)
			subCmd.Short = i18n.T(i18n.InitShort)
			subCmd.Long = i18n.T(i18n.InitLong)
		case "completion":
			subCmd.Short = i18n.T(i18n.CompletionShort)
		case "help":
			subCmd.Short = i18n.T(i18n.HelpShort)
		}
		// 递归更新子命令的子命令
		if subCmd.HasSubCommands() {
			updateCommandDescriptions(subCmd)
		}
	}
	
	// 更新标志说明
	if flag := cmd.PersistentFlags().Lookup("config"); flag != nil {
		flag.Usage = i18n.T(i18n.CobraFlagConfig)
	}
	if flag := cmd.PersistentFlags().Lookup("verbose"); flag != nil {
		flag.Usage = i18n.T(i18n.CobraFlagVerbose)
	}
	if flag := cmd.Flags().Lookup("dry-run"); flag != nil {
		flag.Usage = i18n.T(i18n.CobraFlagDryRun)
	}
	if flag := cmd.Flags().Lookup("force"); flag != nil {
		flag.Usage = i18n.T(i18n.CobraFlagForce)
	}
	if flag := cmd.Flags().Lookup("quiet"); flag != nil {
		flag.Usage = i18n.T(i18n.CobraFlagQuiet)
	}
	if flag := cmd.Flags().Lookup("no-send-stdin"); flag != nil {
		flag.Usage = i18n.T(i18n.CobraFlagNoSendStdin)
	}
	if flag := cmd.Flags().Lookup("history"); flag != nil {
		flag.Usage = i18n.T(i18n.CobraFlagHistory)
	}
	if flag := cmd.Flags().Lookup("retry"); flag != nil {
		flag.Usage = i18n.T(i18n.CobraFlagRetry)
	}
	if flag := cmd.Flags().Lookup("version"); flag != nil {
		flag.Usage = i18n.T(i18n.VersionShort)
	}
	if flag := cmd.Flags().Lookup("help"); flag != nil {
		flag.Usage = i18n.T(i18n.HelpFlag)
	}
}
