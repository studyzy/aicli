package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestAnthropicProvider_Translate_Success 测试成功的命令转换
func TestAnthropicProvider_Translate_Success(t *testing.T) {
	// 创建模拟的 Anthropic API 服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != http.MethodPost {
			t.Errorf("期望 POST 请求, 实际为 %s", r.Method)
		}

		// 验证请求头
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("期望 Content-Type: application/json")
		}
		if r.Header.Get("x-api-key") == "" {
			t.Error("期望存在 x-api-key 请求头")
		}
		if r.Header.Get("anthropic-version") == "" {
			t.Error("期望存在 anthropic-version 请求头")
		}

		// 返回模拟响应
		resp := map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": "ls -la",
				},
			},
			"model": "claude-3-sonnet-20240229",
			"role":  "assistant",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 创建 Anthropic Provider 实例
	provider := NewAnthropicProvider("test-api-key", "claude-3-sonnet-20240229", server.URL)

	// 执行转换
	ctx := context.Background()
	execCtx := &ExecutionContext{
		OS:      "linux",
		Shell:   "bash",
		WorkDir: "/home/user",
	}

	command, err := provider.Translate(ctx, "列出当前目录所有文件", execCtx)

	// 验证结果
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}
	if command != "ls -la" {
		t.Errorf("期望命令 'ls -la', 实际为 '%s'", command)
	}
}

// TestAnthropicProvider_Translate_WithStdin 测试包含 stdin 的命令转换
func TestAnthropicProvider_Translate_WithStdin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 读取请求体验证 stdin 是否传递
		var reqBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&reqBody)

		messages := reqBody["messages"].([]interface{})
		found := false
		for _, msg := range messages {
			m := msg.(map[string]interface{})
			content := m["content"].(string)
			if len(content) > 0 && len(content) < 1000 {
				found = true
			}
		}
		if !found {
			t.Error("期望在消息中找到上下文信息")
		}

		resp := map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": "grep ERROR",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	provider := NewAnthropicProvider("test-api-key", "claude-3-sonnet-20240229", server.URL)

	ctx := context.Background()
	execCtx := &ExecutionContext{
		OS:      "linux",
		Shell:   "bash",
		WorkDir: "/home/user",
		Stdin:   "log file content with ERROR messages",
	}

	command, err := provider.Translate(ctx, "查找错误日志", execCtx)

	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}
	if command != "grep ERROR" {
		t.Errorf("期望命令 'grep ERROR', 实际为 '%s'", command)
	}
}

// TestAnthropicProvider_Translate_APIError 测试 API 错误处理
func TestAnthropicProvider_Translate_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		resp := map[string]interface{}{
			"error": map[string]interface{}{
				"type":    "authentication_error",
				"message": "Invalid API key",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	provider := NewAnthropicProvider("invalid-key", "claude-3-sonnet-20240229", server.URL)

	ctx := context.Background()
	execCtx := &ExecutionContext{OS: "linux", Shell: "bash"}

	_, err := provider.Translate(ctx, "test command", execCtx)

	if err == nil {
		t.Fatal("期望返回错误，但成功了")
	}
}

// TestAnthropicProvider_Translate_Timeout 测试超时处理
func TestAnthropicProvider_Translate_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟慢响应
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	provider := NewAnthropicProvider("test-api-key", "claude-3-sonnet-20240229", server.URL)

	// 设置短超时
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	execCtx := &ExecutionContext{OS: "linux", Shell: "bash"}

	_, err := provider.Translate(ctx, "test command", execCtx)

	if err == nil {
		t.Fatal("期望超时错误，但成功了")
	}
}

// TestAnthropicProvider_Translate_InvalidResponse 测试无效响应处理
func TestAnthropicProvider_Translate_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"invalid": "response"}`))
	}))
	defer server.Close()

	provider := NewAnthropicProvider("test-api-key", "claude-3-sonnet-20240229", server.URL)

	ctx := context.Background()
	execCtx := &ExecutionContext{OS: "linux", Shell: "bash"}

	_, err := provider.Translate(ctx, "test command", execCtx)

	if err == nil {
		t.Fatal("期望返回错误，但成功了")
	}
}

// TestAnthropicProvider_Translate_EmptyCommand 测试空命令响应
func TestAnthropicProvider_Translate_EmptyCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": "   ",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	provider := NewAnthropicProvider("test-api-key", "claude-3-sonnet-20240229", server.URL)

	ctx := context.Background()
	execCtx := &ExecutionContext{OS: "linux", Shell: "bash"}

	_, err := provider.Translate(ctx, "test command", execCtx)

	if err == nil {
		t.Fatal("期望返回错误（空命令），但成功了")
	}
}

// TestAnthropicProvider_Name 测试提供商名称
func TestAnthropicProvider_Name(t *testing.T) {
	provider := NewAnthropicProvider("test-key", "claude-3-sonnet-20240229", "")
	if provider.Name() != "anthropic" {
		t.Errorf("期望名称 'anthropic', 实际为 '%s'", provider.Name())
	}
}
