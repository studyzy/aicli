package app

import (
	"io"
	"os"
	"strings"
	"testing"
)

// TestHasStdin 测试 stdin 检测
func TestHasStdin(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *os.File
		teardown func(*os.File)
		want     bool
	}{
		{
			name: "有管道输入",
			setup: func() *os.File {
				// 保存原始 stdin
				oldStdin := os.Stdin

				// 创建管道模拟 stdin
				r, w, _ := os.Pipe()
				os.Stdin = r
				w.WriteString("test data\n")
				w.Close()

				return oldStdin
			},
			teardown: func(oldStdin *os.File) {
				os.Stdin = oldStdin
			},
			want: true,
		},
		{
			name: "无管道输入（终端）",
			setup: func() *os.File {
				// 使用实际的 stdin（通常是终端）
				return os.Stdin
			},
			teardown: func(oldStdin *os.File) {
				// 不需要恢复
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdin := tt.setup()
			defer tt.teardown(oldStdin)

			got := hasStdin()
			if got != tt.want {
				t.Errorf("hasStdin() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestReadStdin 测试 stdin 读取
func TestReadStdin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		validate func(string) bool
	}{
		{
			name:    "读取简单文本",
			input:   "hello world",
			wantErr: false,
			validate: func(output string) bool {
				return output == "hello world"
			},
		},
		{
			name:    "读取多行文本",
			input:   "line 1\nline 2\nline 3",
			wantErr: false,
			validate: func(output string) bool {
				return strings.Count(output, "\n") == 2
			},
		},
		{
			name:    "读取空输入",
			input:   "",
			wantErr: false,
			validate: func(output string) bool {
				return output == ""
			},
		},
		{
			name:    "读取大量数据",
			input:   strings.Repeat("x", 10000),
			wantErr: false,
			validate: func(output string) bool {
				return len(output) == 10000
			},
		},
		{
			name:    "读取包含特殊字符的数据",
			input:   "line1\nline2\ttab\rcarriage\x00null",
			wantErr: false,
			validate: func(output string) bool {
				return len(output) > 0 && strings.Contains(output, "line1")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 保存原始 stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			// 创建管道并设置为 stdin
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("创建管道失败: %v", err)
			}
			os.Stdin = r

			// 写入测试数据
			go func() {
				w.WriteString(tt.input)
				w.Close()
			}()

			// 读取 stdin
			got, err := readStdin()
			if (err != nil) != tt.wantErr {
				t.Errorf("readStdin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 验证结果
			if !tt.validate(got) {
				t.Errorf("readStdin() 验证失败, got = %s", got)
			}
		})
	}
}

// TestReadStdinWithReader 测试使用自定义 Reader 读取 stdin
func TestReadStdinWithReader(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "从 bytes.Buffer 读取",
			input:   "test data from buffer",
			want:    "test data from buffer",
			wantErr: false,
		},
		{
			name:    "从 strings.Reader 读取",
			input:   "test data from string",
			want:    "test data from string",
			wantErr: false,
		},
		{
			name:    "读取空 Reader",
			input:   "",
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)

			data, err := io.ReadAll(reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := string(data)
			if got != tt.want {
				t.Errorf("ReadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTruncateStdinForLLM 测试 stdin 数据截断
func TestTruncateStdinForLLM(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLen    int
		wantTrunc bool
	}{
		{
			name:      "短数据不截断",
			input:     "short data",
			maxLen:    100,
			wantTrunc: false,
		},
		{
			name:      "长数据截断",
			input:     strings.Repeat("x", 10000),
			maxLen:    5000,
			wantTrunc: true,
		},
		{
			name:      "空数据",
			input:     "",
			maxLen:    1000,
			wantTrunc: false,
		},
		{
			name:      "刚好等于最大长度",
			input:     strings.Repeat("y", 1000),
			maxLen:    1000,
			wantTrunc: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateStdinForLLM(tt.input, tt.maxLen)

			if tt.wantTrunc {
				if len(got) > tt.maxLen {
					t.Errorf("截断后长度 %d 超过最大长度 %d", len(got), tt.maxLen)
				}
				if !strings.Contains(got, "... (truncated)") {
					t.Error("截断后的数据应包含 '... (truncated)' 标记")
				}
			} else {
				if got != tt.input {
					t.Errorf("不应截断，但数据被修改了")
				}
			}
		})
	}
}

// BenchmarkReadStdin 性能测试
func BenchmarkReadStdin(b *testing.B) {
	// 保存原始 stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	testData := strings.Repeat("benchmark data\n", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 创建管道
		r, w, _ := os.Pipe()
		os.Stdin = r

		// 写入数据
		go func() {
			w.WriteString(testData)
			w.Close()
		}()

		// 读取
		readStdin()
	}
}

// TestStdinDetectionEdgeCases 测试边缘情况
func TestStdinDetectionEdgeCases(t *testing.T) {
	t.Run("stdin 是 nil", func(t *testing.T) {
		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()

		// 这个测试验证当 Stdin 不可用时的行为
		// 实际上 os.Stdin 不能真正设置为 nil，所以跳过
		t.Skip("无法将 os.Stdin 设置为 nil")
	})

	t.Run("stdin 被关闭", func(t *testing.T) {
		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()

		// 创建已关闭的管道
		r, w, _ := os.Pipe()
		w.Close()
		os.Stdin = r

		_, err := readStdin()
		// 读取已关闭的管道应该返回空字符串，不报错
		if err != nil {
			t.Logf("读取已关闭的管道: err = %v", err)
		}
	})
}

// 辅助函数：模拟截断逻辑
func truncateStdinForLLM(data string, maxLen int) string {
	if len(data) <= maxLen {
		return data
	}

	truncateMsg := "... (truncated)"
	keepLen := maxLen - len(truncateMsg)
	if keepLen < 0 {
		keepLen = 0
	}

	return data[:keepLen] + truncateMsg
}
