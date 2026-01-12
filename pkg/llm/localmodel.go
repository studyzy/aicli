// Package llm 提供了本地模型（如 Ollama）LLM 提供商的实现
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

// LocalModelProvider 实现了本地 LLM 服务（如 Ollama）的 LLMProvider 接口
type LocalModelProvider struct {
	model   string
	baseURL string
	client  *http.Client
}

// ollamaRequest 表示 Ollama API 请求体
type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

// ollamaMessage 表示 Ollama 消息
type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ollamaResponse 表示 Ollama API 响应体
type ollamaResponse struct {
	Message *struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done  bool   `json:"done"`
	Error string `json:"error,omitempty"`
}

// NewLocalModelProvider 创建一个新的本地模型 Provider
// model: 模型名称，如 "llama2", "mistral", "codellama" 等
// baseURL: Ollama 服务地址，默认为 http://localhost:11434
func NewLocalModelProvider(model, baseURL string) *LocalModelProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	return &LocalModelProvider{
		model:   model,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 60 * time.Second, // 本地模型可能需要更长时间
		},
	}
}

// Name 返回提供商名称
func (p *LocalModelProvider) Name() string {
	return "local"
}

// Translate 将自然语言转换为命令
func (p *LocalModelProvider) Translate(ctx context.Context, input string, execCtx *ExecutionContext) (string, error) {
	if input == "" {
		return "", fmt.Errorf("输入不能为空")
	}

	// 构建提示词
	prompt := BuildPrompt(input, execCtx)
	systemPrompt := GetSystemPrompt(execCtx)

	// 构建请求体
	reqBody := ollamaRequest{
		Model:  p.model,
		Stream: false, // 不使用流式响应
		Messages: []ollamaMessage{
			{
				Role:    "system",
				Content: systemPrompt,
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
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建 HTTP 请求
	url := p.baseURL + "/api/chat"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

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
		var errResp ollamaResponse
		json.Unmarshal(body, &errResp)
		errMsg := fmt.Sprintf("HTTP %d", resp.StatusCode)
		if errResp.Error != "" {
			errMsg = fmt.Sprintf("%s: %s", errMsg, errResp.Error)
		}
		return "", fmt.Errorf("API 错误: %s", errMsg)
	}

	// 解析响应
	var apiResp ollamaResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 验证响应
	if apiResp.Message == nil {
		return "", fmt.Errorf("API 返回空响应")
	}

	command := strings.TrimSpace(apiResp.Message.Content)
	if command == "" {
		return "", fmt.Errorf("API 返回空命令")
	}

	// 清理命令（移除可能的 markdown 代码块标记）
	command = cleanCommand(command)

	return command, nil
}
