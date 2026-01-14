// Package llm 提供了内置试用 LLM 提供商的实现
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	providerBuiltin = "builtin"
	// 内置试用的API配置（这里使用硬编码的配置）
	builtinAPIKey  = "sk-devin-trial-key"
	builtinAPIBase = "http://124.221.164.209:8888/v1"
	builtinModel   = "DeepSeek-V3-0324" // 使用较便宜的模型作为试用
)

// BuiltinProvider 实现了内置试用的 LLMProvider 接口
// 使用硬编码的API配置，用户无需配置即可试用
type BuiltinProvider struct {
	client *http.Client
}

// NewBuiltinProvider 创建一个新的内置试用 Provider
func NewBuiltinProvider() *BuiltinProvider {
	return &BuiltinProvider{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name 返回 Provider 名称
func (p *BuiltinProvider) Name() string {
	return providerBuiltin
}

// Translate 将自然语言转换为命令
func (p *BuiltinProvider) Translate(ctx context.Context, input string, execCtx *ExecutionContext) (string, error) {
	if input == "" {
		return "", fmt.Errorf("输入不能为空")
	}

	// 构建提示词
	prompt := BuildPrompt(input, execCtx)
	systemPrompt := GetSystemPrompt(execCtx)

	// 构建请求体（使用 OpenAI 兼容格式）
	reqBody := map[string]interface{}{
		"model": builtinModel,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.3,
		"max_tokens":  500,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", builtinAPIBase+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+builtinAPIKey)

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	command := strings.TrimSpace(result.Choices[0].Message.Content)

	// 清理可能的 markdown 代码块格式
	command = strings.TrimPrefix(command, "```bash")
	command = strings.TrimPrefix(command, "```sh")
	command = strings.TrimPrefix(command, "```")
	command = strings.TrimSuffix(command, "```")
	command = strings.TrimSpace(command)

	return command, nil
}
