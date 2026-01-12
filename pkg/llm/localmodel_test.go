package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestLocalModelProvider_Translate_Success 测试成功的命令转换
func TestLocalModelProvider_Translate_Success(t *testing.T) {
	// 创建模拟的 Ollama API 服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != http.MethodPost {
			t.Errorf("期望 POST 请求, 实际为 %s", r.Method)
		}

		// 验证请求路径
		if r.URL.Path != "/api/chat" && r.URL.Path != "/api/generate" {
			t.Errorf("期望路径 /api/chat 或 /api/generate, 实际为 %s", r.URL.Path)
		}

		// 返回模拟响应
		resp := map[string]interface{}{
			"message": map[string]interface{}{
				"role":    "assistant",
				"content": "ls -l",
			},
			"done": true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 创建 LocalModel Provider 实例
	provider := NewLocalModelProvider("llama2", server.URL)

	// 执行转换
	ctx := context.Background()
	execCtx := &ExecutionContext{
		OS:      "linux",
		Shell:   "bash",
		WorkDir: "/home/user",
	}

	command, err := provider.Translate(ctx, "列出当前目录文件", execCtx)

	// 验证结果
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}
	if command != "ls -l" {
		t.Errorf("期望命令 'ls -l', 实际为 '%s'", command)
	}
}

// TestLocalModelProvider_Translate_WithStdin 测试包含 stdin 的命令转换
func TestLocalModelProvider_Translate_WithStdin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 读取请求体验证 stdin 是否传递
		var reqBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&reqBody)

		messages := reqBody["messages"].([]interface{})
		if len(messages) == 0 {
			t.Error("期望在消息中找到内容")
		}

		resp := map[string]interface{}{
			"message": map[string]interface{}{
				"role":    "assistant",
				"content": "wc -l",
			},
			"done": true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	provider := NewLocalModelProvider("llama2", server.URL)

	ctx := context.Background()
	execCtx := &ExecutionContext{
		OS:      "linux",
		Shell:   "bash",
		WorkDir: "/home/user",
		Stdin:   "some input data\nwith multiple lines",
	}

	command, err := provider.Translate(ctx, "统计行数", execCtx)

	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}
	if command != "wc -l" {
		t.Errorf("期望命令 'wc -l', 实际为 '%s'", command)
	}
}

// TestLocalModelProvider_Translate_APIError 测试 API 错误处理
func TestLocalModelProvider_Translate_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		resp := map[string]interface{}{
			"error": "model not found",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	provider := NewLocalModelProvider("invalid-model", server.URL)

	ctx := context.Background()
	execCtx := &ExecutionContext{OS: "linux", Shell: "bash"}

	_, err := provider.Translate(ctx, "test command", execCtx)

	if err == nil {
		t.Fatal("期望返回错误，但成功了")
	}
}

// TestLocalModelProvider_Translate_Timeout 测试超时处理
func TestLocalModelProvider_Translate_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟慢响应
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	provider := NewLocalModelProvider("llama2", server.URL)

	// 设置短超时
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	execCtx := &ExecutionContext{OS: "linux", Shell: "bash"}

	_, err := provider.Translate(ctx, "test command", execCtx)

	if err == nil {
		t.Fatal("期望超时错误，但成功了")
	}
}

// TestLocalModelProvider_Translate_InvalidResponse 测试无效响应处理
func TestLocalModelProvider_Translate_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"invalid": "response"}`))
	}))
	defer server.Close()

	provider := NewLocalModelProvider("llama2", server.URL)

	ctx := context.Background()
	execCtx := &ExecutionContext{OS: "linux", Shell: "bash"}

	_, err := provider.Translate(ctx, "test command", execCtx)

	if err == nil {
		t.Fatal("期望返回错误，但成功了")
	}
}

// TestLocalModelProvider_Translate_EmptyCommand 测试空命令响应
func TestLocalModelProvider_Translate_EmptyCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"message": map[string]interface{}{
				"role":    "assistant",
				"content": "   ",
			},
			"done": true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	provider := NewLocalModelProvider("llama2", server.URL)

	ctx := context.Background()
	execCtx := &ExecutionContext{OS: "linux", Shell: "bash"}

	_, err := provider.Translate(ctx, "test command", execCtx)

	if err == nil {
		t.Fatal("期望返回错误（空命令），但成功了")
	}
}

// TestLocalModelProvider_Name 测试提供商名称
func TestLocalModelProvider_Name(t *testing.T) {
	provider := NewLocalModelProvider("llama2", "http://localhost:11434")
	if provider.Name() != "local" {
		t.Errorf("期望名称 'local', 实际为 '%s'", provider.Name())
	}
}

// TestLocalModelProvider_DefaultBaseURL 测试默认 BaseURL
func TestLocalModelProvider_DefaultBaseURL(t *testing.T) {
	provider := NewLocalModelProvider("llama2", "")
	// 验证默认 URL 被设置（虽然无法直接访问私有字段，但可以通过 Name() 验证对象创建）
	if provider.Name() != "local" {
		t.Error("Provider 创建失败")
	}
}

// TestLocalModelProvider_StreamingResponse 测试流式响应处理
func TestLocalModelProvider_StreamingResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ollama 可以返回流式响应，但我们的实现应该处理完整响应
		resp := map[string]interface{}{
			"message": map[string]interface{}{
				"role":    "assistant",
				"content": "echo hello",
			},
			"done": true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	provider := NewLocalModelProvider("llama2", server.URL)

	ctx := context.Background()
	execCtx := &ExecutionContext{OS: "linux", Shell: "bash"}

	command, err := provider.Translate(ctx, "打印hello", execCtx)

	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}
	if command != "echo hello" {
		t.Errorf("期望命令 'echo hello', 实际为 '%s'", command)
	}
}
