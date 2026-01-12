// Package llm 提供了 LLM 提供商的测试
package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestOpenAIProvider_Translate_Success 测试成功的命令转换
func TestOpenAIProvider_Translate_Success(t *testing.T) {
	// 创建模拟的 OpenAI API 服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法和路径
		if r.Method != http.MethodPost {
			t.Errorf("期望 POST 请求, 实际为 %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/chat/completions") {
			t.Errorf("期望路径包含 /chat/completions, 实际为 %s", r.URL.Path)
		}

		// 验证请求头
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("期望 Authorization 头为 'Bearer test-api-key', 实际为 %s", r.Header.Get("Authorization"))
		}

		// 返回模拟的成功响应
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]string{
						"content": "ls -la",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建 OpenAI Provider 实例
	provider := NewOpenAIProvider("test-api-key", "gpt-4", server.URL)

	// 创建执行上下文
	ctx := &ExecutionContext{
		OS:      "darwin",
		Shell:   "zsh",
		WorkDir: "/tmp",
	}

	// 执行转换
	command, err := provider.Translate(context.Background(), "列出当前目录的所有文件", ctx)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	// 验证结果
	if command != "ls -la" {
		t.Errorf("期望命令为 'ls -la', 实际为 '%s'", command)
	}
}

// TestOpenAIProvider_Translate_WithStdin 测试包含 stdin 的命令转换
func TestOpenAIProvider_Translate_WithStdin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 解析请求体
		var reqBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&reqBody)

		// 验证消息中包含 stdin 信息
		messages := reqBody["messages"].([]interface{})
		found := false
		for _, msg := range messages {
			msgMap := msg.(map[string]interface{})
			content := msgMap["content"].(string)
			if strings.Contains(content, "标准输入数据") || strings.Contains(content, "test input data") {
				found = true
				break
			}
		}
		if !found {
			t.Error("请求消息中未包含 stdin 信息")
		}

		// 返回响应
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]string{
						"content": "grep ERROR",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	provider := NewOpenAIProvider("test-api-key", "gpt-4", server.URL)

	ctx := &ExecutionContext{
		OS:      "darwin",
		Shell:   "zsh",
		WorkDir: "/tmp",
		Stdin:   "test input data\nwith errors\nERROR: something went wrong",
	}

	command, err := provider.Translate(context.Background(), "过滤出包含ERROR的行", ctx)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	if command != "grep ERROR" {
		t.Errorf("期望命令为 'grep ERROR', 实际为 '%s'", command)
	}
}

// TestOpenAIProvider_Translate_APIError 测试 API 错误处理
func TestOpenAIProvider_Translate_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"message": "Invalid API key",
			},
		})
	}))
	defer server.Close()

	provider := NewOpenAIProvider("invalid-key", "gpt-4", server.URL)

	ctx := &ExecutionContext{
		OS:      "darwin",
		Shell:   "zsh",
		WorkDir: "/tmp",
	}

	_, err := provider.Translate(context.Background(), "列出文件", ctx)
	if err == nil {
		t.Fatal("期望返回错误，但成功了")
	}

	if !strings.Contains(err.Error(), "API") && !strings.Contains(err.Error(), "401") {
		t.Errorf("期望错误信息包含 API 或 401, 实际为: %v", err)
	}
}

// TestOpenAIProvider_Translate_Timeout 测试超时处理
func TestOpenAIProvider_Translate_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟长时间处理
		select {
		case <-r.Context().Done():
			return
		}
	}))
	defer server.Close()

	provider := NewOpenAIProvider("test-api-key", "gpt-4", server.URL)

	ctx := &ExecutionContext{
		OS:      "darwin",
		Shell:   "zsh",
		WorkDir: "/tmp",
	}

	// 创建带超时的上下文
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 100)
	defer cancel()

	_, err := provider.Translate(timeoutCtx, "列出文件", ctx)
	if err == nil {
		t.Fatal("期望返回超时错误，但成功了")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") && !strings.Contains(err.Error(), "timeout") {
		t.Errorf("期望超时错误, 实际为: %v", err)
	}
}

// TestOpenAIProvider_Translate_InvalidResponse 测试无效响应处理
func TestOpenAIProvider_Translate_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"invalid": "response"}`))
	}))
	defer server.Close()

	provider := NewOpenAIProvider("test-api-key", "gpt-4", server.URL)

	ctx := &ExecutionContext{
		OS:      "darwin",
		Shell:   "zsh",
		WorkDir: "/tmp",
	}

	_, err := provider.Translate(context.Background(), "列出文件", ctx)
	if err == nil {
		t.Fatal("期望返回错误，但成功了")
	}
}

// TestOpenAIProvider_Translate_EmptyCommand 测试空命令响应
func TestOpenAIProvider_Translate_EmptyCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]string{
						"content": "",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	provider := NewOpenAIProvider("test-api-key", "gpt-4", server.URL)

	ctx := &ExecutionContext{
		OS:      "darwin",
		Shell:   "zsh",
		WorkDir: "/tmp",
	}

	_, err := provider.Translate(context.Background(), "列出文件", ctx)
	if err == nil {
		t.Fatal("期望返回错误（空命令），但成功了")
	}
}
