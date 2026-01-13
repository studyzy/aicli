package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/studyzy/aicli/pkg/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化配置",
	Long:  `引导用户设置 LLM 配置并生成配置文件 ~/.aicli.json`,
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	var configPath string
	if flags.Config != "" {
		configPath = flags.Config
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户主目录失败: %w", err)
		}
		configPath = homeDir + "/.aicli.json"
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("配置文件 %s 已存在。\n", configPath)
		if !promptBool(reader, "是否覆盖？", false) {
			fmt.Println("操作已取消。")
			return nil
		}
	}

	fmt.Println("欢迎使用 aicli 配置向导！")
	fmt.Println("我们将引导您完成基本配置。")
	fmt.Println()

	cfg := config.Default()

	// 1. 选择 LLM 提供商
	fmt.Println("请选择 LLM 提供商:")
	fmt.Println("1. OpenAI (GPT-4, GPT-3.5)")
	fmt.Println("2. Anthropic (Claude)")
	fmt.Println("3. Local (Ollama, LocalAI)")
	fmt.Println("4. DeepSeek (深度求索)")
	fmt.Println("5. Other (兼容 OpenAI 协议)")

	providerChoice := prompt(reader, "请输入序号 [1-5]", "1")

	switch providerChoice {
	case "1":
		cfg.LLM.Provider = "openai"
		cfg.LLM.Model = "gpt-4"
		cfg.LLM.APIBase = "https://api.openai.com/v1"
	case "2":
		cfg.LLM.Provider = "anthropic"
		cfg.LLM.Model = "claude-3-sonnet-20240229"
		cfg.LLM.APIBase = "https://api.anthropic.com/v1"
	case "3":
		cfg.LLM.Provider = "local"
		cfg.LLM.Model = "llama2"
		cfg.LLM.APIBase = "http://localhost:11434"
	case "4":
		cfg.LLM.Provider = "openai" // DeepSeek 兼容 OpenAI
		cfg.LLM.Model = "deepseek-chat"
		cfg.LLM.APIBase = "https://api.deepseek.com/v1"
	case "5":
		cfg.LLM.Provider = "openai"
		cfg.LLM.Model = "gpt-3.5-turbo"
		cfg.LLM.APIBase = "https://api.openai.com/v1"
	default:
		fmt.Println("无效的选择，默认使用 OpenAI")
		cfg.LLM.Provider = "openai"
	}

	// 2. 配置 API Key
	if cfg.LLM.Provider != "local" {
		cfg.LLM.APIKey = prompt(reader, "请输入 API Key", "")
		if cfg.LLM.APIKey == "" {
			fmt.Println("警告: API Key 为空，您可能需要在环境变量 AICLI_API_KEY 中设置。")
		}
	}

	// 3. 配置 Model
	cfg.LLM.Model = prompt(reader, fmt.Sprintf("请输入模型名称 (默认: %s)", cfg.LLM.Model), cfg.LLM.Model)

	// 4. 配置 API Base
	cfg.LLM.APIBase = prompt(reader, fmt.Sprintf("请输入 API Base URL (默认: %s)", cfg.LLM.APIBase), cfg.LLM.APIBase)

	// 5. 其他配置
	fmt.Println("\n--- 其他设置 ---")
	cfg.Safety.EnableChecks = promptBool(reader, "是否启用危险命令安全检查？", true)
	cfg.History.Enabled = promptBool(reader, "是否启用历史记录？", true)
	
	// 语言偏好（虽然目前配置结构体里没有，但可以作为扩展点，或者只是提示用户）
	// 目前先跳过，因为 Config 结构体没有 Language 字段

	// 6. 保存配置
	fmt.Println("\n正在保存配置...")
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	fmt.Printf("配置已成功保存到 %s\n", configPath)
	fmt.Println("现在您可以开始使用 aicli 了！")
	fmt.Println("示例: aicli \"查询我的公网IP\"")

	return nil
}

func prompt(reader *bufio.Reader, label string, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", label, defaultValue)
	} else {
		fmt.Printf("%s: ", label)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}

func promptBool(reader *bufio.Reader, label string, defaultValue bool) bool {
	defaultStr := "y"
	if !defaultValue {
		defaultStr = "n"
	}

	fmt.Printf("%s (y/n) [%s]: ", label, defaultStr)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)

	if input == "" {
		return defaultValue
	}

	return input == "y" || input == "yes"
}
