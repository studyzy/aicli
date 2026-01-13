// Package llm 提供了 Anthropic LLM 提供商的实现
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
	providerAnthropic = "anthropic"
)

// AnthropicProvider 实现了 Anthropic Claude API 的 LLMProvider 接口
type AnthropicProvider struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// anthropicRequest 表示 Anthropic API 请求体
type anthropicRequest struct {
	Model     string             `json:"model"`
	Messages  []anthropicMessage `json:"messages"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
}

// anthropicMessage 表示 Anthropic 消息
type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// anthropicResponse 表示 Anthropic API 响应体
type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// NewAnthropicProvider 创建一个新的 Anthropic Provider
func NewAnthropicProvider(apiKey, model, baseURL string) *AnthropicProvider {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}

	return &AnthropicProvider{
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name 返回提供商名称
func (p *AnthropicProvider) Name() string {
	return providerAnthropic
}

// Translate 将自然语言转换为命令
func (p *AnthropicProvider) Translate(ctx context.Context, input string, execCtx *ExecutionContext) (string, error) {
	if input == "" {
		return "", fmt.Errorf("输入不能为空")
	}

	// 构建提示词
	prompt := BuildPrompt(input, execCtx)
	systemPrompt := GetSystemPrompt(execCtx)

	// 构建请求体
	reqBody := anthropicRequest{
		Model:     p.model,
		MaxTokens: 500,
		System:    systemPrompt,
		Messages: []anthropicMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// 序列化请求体
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建 HTTP 请求
	url := p.baseURL + "/messages"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		var errResp anthropicResponse
		json.Unmarshal(body, &errResp)
		errMsg := fmt.Sprintf("HTTP %d", resp.StatusCode)
		if errResp.Error != nil {
			errMsg = fmt.Sprintf("%s: %s", errMsg, errResp.Error.Message)
		}
		return "", fmt.Errorf("API 错误: %s", errMsg)
	}

	// 解析响应
	var apiResp anthropicResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 验证响应
	if len(apiResp.Content) == 0 {
		return "", fmt.Errorf("API 返回空响应")
	}

	// 提取文本内容
	command := ""
	for _, content := range apiResp.Content {
		if content.Type == "text" {
			command = strings.TrimSpace(content.Text)
			break
		}
	}

	if command == "" {
		return "", fmt.Errorf("API 返回空命令")
	}

	// 清理命令（移除可能的 markdown 代码块标记）
	command = cleanCommand(command)

	return command, nil
}
