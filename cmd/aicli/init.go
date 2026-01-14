// Package main 提供aicli命令行工具
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/studyzy/aicli/pkg/config"
	"github.com/studyzy/aicli/pkg/i18n"
)

const (
	providerOpenAI    = "openai"
	providerLocal     = "local"
	providerAnthropic = "anthropic"
	providerBuiltin   = "builtin"
	apiBaseOpenAI     = "https://api.openai.com/v1"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "", // 将在 main 中通过 updateCommandDescriptions 设置
	Long:  "", // 将在 main 中通过 updateCommandDescriptions 设置
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// 初始化 i18n (在配置文件加载前使用默认配置)
	cfg := config.Default()
	i18n.Init(cfg)

	reader := bufio.NewReader(os.Stdin)

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := checkExistingConfig(reader, configPath); err != nil {
		return err
	}

	fmt.Println(i18n.T(i18n.MsgWelcomeInit))
	fmt.Println(i18n.T(i18n.MsgInitGuide))
	fmt.Println()

	if err := configureProvider(reader, cfg); err != nil {
		return err
	}

	if err := configureAPI(reader, cfg); err != nil {
		return err
	}

	if err := configureOtherSettings(reader, cfg); err != nil {
		return err
	}

	return saveConfig(cfg, configPath)
}

func getConfigPath() (string, error) {
	if flags.Config != "" {
		return flags.Config, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T(i18n.ErrGetUserHome), err)
	}
	return homeDir + "/.aicli.json", nil
}

func checkExistingConfig(reader *bufio.Reader, configPath string) error {
	if _, err := os.Stat(configPath); err == nil {
		msg := i18n.T(i18n.MsgConfigExists, configPath)
		fmt.Printf("%s\n", msg)
		if !promptBool(reader, i18n.T(i18n.PromptOverwriteConfig), false) {
			fmt.Println(i18n.T(i18n.MsgOperationCancelled))
			return fmt.Errorf("%s", i18n.T(i18n.ErrUserCancelled))
		}
	}
	return nil
}

func configureProvider(reader *bufio.Reader, cfg *config.Config) error {
	fmt.Println(i18n.T(i18n.PromptSelectProvider) + ":")
	fmt.Println(i18n.T(i18n.InitProviderOpenAI))
	fmt.Println(i18n.T(i18n.InitProviderAnthropic))
	fmt.Println(i18n.T(i18n.InitProviderLocal))
	fmt.Println(i18n.T(i18n.InitProviderDeepSeek))
	fmt.Println(i18n.T(i18n.InitProviderOther))

	providerChoice := prompt(reader, i18n.T(i18n.PromptInputChoice)+" [1-5]", "1")

	switch providerChoice {
	case "1":
		cfg.LLM.Provider = providerOpenAI
		cfg.LLM.Model = "gpt-4"
		cfg.LLM.APIBase = apiBaseOpenAI
	case "2":
		cfg.LLM.Provider = providerAnthropic
		cfg.LLM.Model = "claude-3-sonnet-20240229"
		cfg.LLM.APIBase = "https://api.anthropic.com/v1"
	case "3":
		cfg.LLM.Provider = providerLocal
		cfg.LLM.Model = "llama2"
		cfg.LLM.APIBase = "http://localhost:11434"
	case "4":
		cfg.LLM.Provider = providerOpenAI // DeepSeek 兼容 OpenAI
		cfg.LLM.Model = "deepseek-chat"
		cfg.LLM.APIBase = "https://api.deepseek.com/v1"
	case "5":
		cfg.LLM.Provider = providerOpenAI
		cfg.LLM.Model = "gpt-3.5-turbo"
		cfg.LLM.APIBase = apiBaseOpenAI
	default:
		fmt.Println(i18n.T(i18n.MsgDefaultUseOpenAI))
		cfg.LLM.Provider = providerOpenAI
		cfg.LLM.Model = "gpt-4"
		cfg.LLM.APIBase = apiBaseOpenAI
	}

	return nil
}

func configureAPI(reader *bufio.Reader, cfg *config.Config) error {
	if cfg.LLM.Provider != providerLocal {
		cfg.LLM.APIKey = prompt(reader, i18n.T(i18n.PromptEnterAPIKey), "")
		if cfg.LLM.APIKey == "" {
			fmt.Println(i18n.T(i18n.WarnAPIKeyEmpty))
		}
	}

	modelPrompt := fmt.Sprintf("%s (%s: %s)", i18n.T(i18n.PromptEnterModel), i18n.T(i18n.MsgDefault), cfg.LLM.Model)
	cfg.LLM.Model = prompt(reader, modelPrompt, cfg.LLM.Model)
	
	apiBasePrompt := fmt.Sprintf("%s (%s: %s)", i18n.T(i18n.PromptEnterAPIBase), i18n.T(i18n.MsgDefault), cfg.LLM.APIBase)
	cfg.LLM.APIBase = prompt(reader, apiBasePrompt, cfg.LLM.APIBase)

	return nil
}

func configureOtherSettings(reader *bufio.Reader, cfg *config.Config) error {
	fmt.Println("\n" + i18n.T(i18n.MsgOtherSettings))
	cfg.Safety.EnableChecks = promptBool(reader, i18n.T(i18n.PromptEnableCheck), true)
	cfg.History.Enabled = promptBool(reader, i18n.T(i18n.PromptEnableHistory), true)
	return nil
}

func saveConfig(cfg *config.Config, configPath string) error {
	fmt.Println("\n" + i18n.T(i18n.MsgSavingConfig))
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrSaveConfig), err)
	}

	fmt.Println(i18n.T(i18n.MsgConfigSaved, configPath))
	fmt.Println(i18n.T(i18n.MsgNowCanUse))
	fmt.Println(i18n.T(i18n.MsgExampleUsage))

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
