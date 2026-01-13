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
	"github.com/studyzy/aicli/pkg/llm"
	"github.com/studyzy/aicli/pkg/safety"
)

// version 是当前版本号
const version = "0.1.0-dev"

var (
	flags = app.NewFlags()
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
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
	Version: version,
	RunE:    run,
}

func run(cmd *cobra.Command, args []string) error {
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

	// 加载配置
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 创建 LLM Provider
	provider, err := createLLMProvider(cfg)
	if err != nil {
		return fmt.Errorf("创建 LLM Provider 失败: %w", err)
	}

	// 创建 Executor
	exec := executor.NewExecutor()

	// 创建 Safety Checker
	checker := safety.NewSafetyChecker(cfg.Safety.EnableChecks)

	// 创建应用实例
	application := app.NewApp(cfg, provider, exec, checker)

	// 加载历史记录
	hist := history.NewHistory()
	historyPath := getHistoryPath()
	if err := hist.Load(historyPath); err != nil {
		if flags.Verbose {
			fmt.Fprintf(os.Stderr, "加载历史记录失败: %v\n", err)
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
		stdinBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("读取 stdin 失败: %w", err)
		}
		stdin = string(stdinBytes)
	}

	// 执行应用逻辑
	output, err := application.Run(input, stdin, flags)

	// 保存历史记录（即使执行失败也保存）
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

// loadConfig 加载配置文件
func loadConfig() (*config.Config, error) {
	// 如果指定了配置文件，使用指定的
	if flags.Config != "" {
		return config.Load(flags.Config)
	}

	// 否则尝试从默认位置加载
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取用户主目录失败: %w", err)
	}

	configPath := homeDir + "/.aicli.json"

	// 如果配置文件不存在，使用默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 提示用户初始化
		fmt.Fprintf(os.Stderr, "提示: 配置文件 %s 不存在。\n", configPath)
		fmt.Fprintf(os.Stderr, "您可以运行 'aicli init' 来快速设置配置。\n\n")

		if flags.Verbose {
			fmt.Fprintf(os.Stderr, "配置文件不存在，使用默认配置\n")
		}
		return config.Default(), nil
	}

	return config.Load(configPath)
}

// createLLMProvider 创建 LLM Provider
// 使用工厂函数统一管理 Provider 创建
func createLLMProvider(cfg *config.Config) (llm.LLMProvider, error) {
	return llm.NewProvider(cfg)
}

func init() {
	// 持久化标志（所有子命令都可用）
	rootCmd.PersistentFlags().StringVarP(&flags.Config, "config", "c", flags.Config, "配置文件路径")
	rootCmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", flags.Verbose, "显示详细输出")

	// 本地标志（仅根命令可用）
	rootCmd.Flags().BoolVarP(&flags.DryRun, "dry-run", "n", flags.DryRun, "仅显示命令不执行")
	rootCmd.Flags().BoolVarP(&flags.Force, "force", "f", flags.Force, "强制执行，跳过确认")
	rootCmd.Flags().BoolVar(&flags.NoSendStdin, "no-send-stdin", flags.NoSendStdin, "不将 stdin 数据发送到 LLM")
	rootCmd.Flags().BoolVar(&flags.History, "history", flags.History, "显示历史记录")
	rootCmd.Flags().IntVar(&flags.Retry, "retry", flags.Retry, "重新执行历史命令 ID")

	// 设置版本模板
	rootCmd.SetVersionTemplate(`{{printf "aicli version %s\n" .Version}}`)
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
		return fmt.Errorf("加载历史记录失败: %w", err)
	}

	entries := hist.List()
	if len(entries) == 0 {
		fmt.Println("没有历史记录")
		return nil
	}

	fmt.Printf("历史记录（共 %d 条）：\n\n", len(entries))
	for _, entry := range entries {
		status := "✓"
		if !entry.Success {
			status = "✗"
		}

		fmt.Printf("[%d] %s %s\n", entry.ID, status, entry.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("    输入: %s\n", entry.Input)
		fmt.Printf("    命令: %s\n", entry.Command)

		if entry.Error != "" {
			fmt.Printf("    错误: %s\n", entry.Error)
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
		return fmt.Errorf("加载历史记录失败: %w", err)
	}

	entry, err := hist.Get(id)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "重新执行历史命令 #%d:\n", id)
	fmt.Fprintf(os.Stderr, "  输入: %s\n", entry.Input)
	fmt.Fprintf(os.Stderr, "  命令: %s\n\n", entry.Command)

	// 加载配置
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 创建 LLM Provider
	provider, err := createLLMProvider(cfg)
	if err != nil {
		return fmt.Errorf("创建 LLM Provider 失败: %w", err)
	}

	// 创建 Executor
	exec := executor.NewExecutor()

	// 创建 Safety Checker
	checker := safety.NewSafetyChecker(cfg.Safety.EnableChecks)

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
