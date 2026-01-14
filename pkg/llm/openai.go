// Package llm 提供了 OpenAI LLM 提供商的实现
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

	"github.com/studyzy/aicli/pkg/i18n"
)

// OpenAIProvider 实现了 OpenAI API 的 LLMProvider 接口
type OpenAIProvider struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// openAIRequest 表示 OpenAI API 请求体
type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

// openAIMessage 表示 OpenAI 消息
type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIResponse 表示 OpenAI API 响应体
type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// NewOpenAIProvider 创建一个新的 OpenAI Provider
func NewOpenAIProvider(apiKey, model, baseURL string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &OpenAIProvider{
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name 返回提供商名称
func (p *OpenAIProvider) Name() string {
	return providerOpenAI
}

// Translate 将自然语言转换为命令
func (p *OpenAIProvider) Translate(ctx context.Context, input string, execCtx *ExecutionContext) (string, error) {
	if input == "" {
		return "", fmt.Errorf("输入不能为空")
	}

	// 构建提示词
	prompt := BuildPrompt(input, execCtx)

	// 构建请求体
	reqBody := openAIRequest{
		Model: p.model,
		Messages: []openAIMessage{
			{
				Role:    "system",
				Content: GetSystemPrompt(execCtx),
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// 序列化请求体
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T(i18n.ErrSerializeRequest), err)
	}

	// 创建 HTTP 请求
	url := p.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T(i18n.ErrCreateRequest), err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T(i18n.ErrAPIRequest), err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T(i18n.ErrReadResponse), err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		var errResp openAIResponse
		json.Unmarshal(body, &errResp)
		errMsg := fmt.Sprintf("HTTP %d", resp.StatusCode)
		if errResp.Error != nil {
			errMsg = fmt.Sprintf("%s: %s", errMsg, errResp.Error.Message)
		}
		return "", fmt.Errorf("%s: %s", i18n.T(i18n.ErrAPIError), errMsg)
	}

	// 解析响应
	var apiResp openAIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T(i18n.ErrParseResponse), err)
	}

	// 验证响应
	if len(apiResp.Choices) == 0 {
		return "", fmt.Errorf("%s", i18n.T(i18n.ErrEmptyResponse))
	}

	command := strings.TrimSpace(apiResp.Choices[0].Message.Content)
	if command == "" {
		return "", fmt.Errorf("%s", i18n.T(i18n.ErrEmptyCommandResp))
	}

	// 清理命令（移除可能的 markdown 代码块标记）
	command = cleanCommand(command)

	return command, nil
}

// cleanCommand 清理命令字符串
func cleanCommand(cmd string) string {
	// 移除 markdown 代码块
	cmd = strings.TrimPrefix(cmd, "```bash")
	cmd = strings.TrimPrefix(cmd, "```sh")
	cmd = strings.TrimPrefix(cmd, "```shell")
	cmd = strings.TrimPrefix(cmd, "```")
	cmd = strings.TrimSuffix(cmd, "```")

	// 移除前后空白
	cmd = strings.TrimSpace(cmd)

	return cmd
}
