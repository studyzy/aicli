package i18n

import (
	"os"
	"testing"

	"github.com/studyzy/aicli/pkg/config"
)

// TestI18nIntegration 集成测试:验证完整的i18n流程
func TestI18nIntegration(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		key      string
		wantZh   string
		wantEn   string
	}{
		{
			name:   "错误信息",
			lang:   "zh",
			key:    ErrLoadConfig,
			wantZh: "加载配置失败",
			wantEn: "Failed to load configuration",
		},
		{
			name:   "LLM系统提示词介绍",
			lang:   "zh",
			key:    LLMSystemPromptIntro,
			wantZh: "你是一个命令行助手,专门将用户的自然语言描述转换为可执行的 shell 命令。",
			wantEn: "You are a command-line assistant that converts natural language descriptions into executable shell commands.",
		},
		{
			name:   "历史记录为空",
			lang:   "zh",
			key:    MsgNoHistory,
			wantZh: "没有历史记录",
			wantEn: "No history records",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试中文
			cfg := &config.Config{Language: "zh"}
			Init(cfg)
			got := T(tt.key)
			if got != tt.wantZh {
				t.Errorf("中文翻译错误: got %q, want %q", got, tt.wantZh)
			}

			// 测试英文
			cfg = &config.Config{Language: "en"}
			Init(cfg)
			got = T(tt.key)
			if got != tt.wantEn {
				t.Errorf("英文翻译错误: got %q, want %q", got, tt.wantEn)
			}
		})
	}
}

// TestLanguageDetectionFlow 测试语言检测流程
func TestLanguageDetectionFlow(t *testing.T) {
	// 保存原始环境变量
	oldLang := os.Getenv("LANG")
	defer os.Setenv("LANG", oldLang)

	tests := []struct {
		name    string
		cfg     *config.Config
		envLang string
		want    string
	}{
		{
			name:    "配置文件优先",
			cfg:     &config.Config{Language: "en"},
			envLang: "zh_CN.UTF-8",
			want:    "en",
		},
		{
			name:    "环境变量检测中文",
			cfg:     &config.Config{},
			envLang: "zh_CN.UTF-8",
			want:    "zh",
		},
		{
			name:    "环境变量检测英文",
			cfg:     &config.Config{},
			envLang: "en_US.UTF-8",
			want:    "en",
		},
		{
			name:    "默认中文",
			cfg:     &config.Config{},
			envLang: "",
			want:    "zh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LANG", tt.envLang)
			Init(tt.cfg)
			got := Lang()
			if got != tt.want {
				t.Errorf("Lang() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFormattedTranslation 测试带参数的翻译
func TestFormattedTranslation(t *testing.T) {
	cfg := &config.Config{Language: "zh"}
	Init(cfg)

	// 测试带参数的翻译
	result := T(MsgHistoryCount, 10)
	expected := "历史记录(共 10 条):"
	if result != expected {
		t.Errorf("格式化翻译错误: got %q, want %q", result, expected)
	}

	// 切换到英文
	cfg = &config.Config{Language: "en"}
	Init(cfg)
	result = T(MsgHistoryCount, 10)
	expected = "History (10 entries):"
	if result != expected {
		t.Errorf("格式化翻译错误: got %q, want %q", result, expected)
	}
}
