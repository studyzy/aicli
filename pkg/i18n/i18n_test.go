package i18n

import (
	"os"
	"testing"

	"github.com/studyzy/aicli/pkg/config"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.Config
		envLang  string
		envLCAll string
		want     string
	}{
		{
			name: "配置文件指定中文",
			cfg:  &config.Config{Language: "zh"},
			want: "zh",
		},
		{
			name: "配置文件指定英文",
			cfg:  &config.Config{Language: "en"},
			want: "en",
		},
		{
			name:    "LANG 环境变量为中文",
			cfg:     &config.Config{},
			envLang: "zh_CN.UTF-8",
			want:    "zh",
		},
		{
			name:    "LANG 环境变量为英文",
			cfg:     &config.Config{},
			envLang: "en_US.UTF-8",
			want:    "en",
		},
		{
			name:     "LC_ALL 环境变量",
			cfg:      &config.Config{},
			envLCAll: "en_GB.UTF-8",
			want:     "en",
		},
		{
			name: "默认值",
			cfg:  &config.Config{},
			want: "zh",
		},
		{
			name:    "配置文件优先于环境变量",
			cfg:     &config.Config{Language: "en"},
			envLang: "zh_CN.UTF-8",
			want:    "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			oldLang := os.Getenv("LANG")
			oldLCAll := os.Getenv("LC_ALL")
			defer func() {
				os.Setenv("LANG", oldLang)
				os.Setenv("LC_ALL", oldLCAll)
			}()

			os.Setenv("LANG", tt.envLang)
			os.Setenv("LC_ALL", tt.envLCAll)

			got := DetectLanguage(tt.cfg)
			if got != tt.want {
				t.Errorf("DetectLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLocale(t *testing.T) {
	tests := []struct {
		locale string
		want   string
	}{
		{"zh_CN.UTF-8", "zh"},
		{"zh_CN", "zh"},
		{"zh", "zh"},
		{"en_US.UTF-8", "en"},
		{"en_US", "en"},
		{"en", "en"},
		{"en_GB.utf8", "en"},
		{"fr_FR.UTF-8", "zh"}, // 不支持的语言 fallback
		{"", "zh"},
	}

	for _, tt := range tests {
		t.Run(tt.locale, func(t *testing.T) {
			got := parseLocale(tt.locale)
			if got != tt.want {
				t.Errorf("parseLocale(%q) = %v, want %v", tt.locale, got, tt.want)
			}
		})
	}
}

func TestNormalizeLanguage(t *testing.T) {
	tests := []struct {
		lang string
		want string
	}{
		{"zh", "zh"},
		{"ZH", "zh"},
		{"zh-CN", "zh"},
		{"zh_CN", "zh"},
		{"en", "en"},
		{"EN", "en"},
		{"en-US", "en"},
		{"en_US", "en"},
		{"fr", "zh"},    // 不支持 fallback
		{"ja", "zh"},    // 不支持 fallback
		{"", "zh"},      // 空字符串 fallback
		{"  en  ", "en"}, // 去除空格
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			got := normalizeLanguage(tt.lang)
			if got != tt.want {
				t.Errorf("normalizeLanguage(%q) = %v, want %v", tt.lang, got, tt.want)
			}
		})
	}
}

func TestLocalizer_T(t *testing.T) {
	// 保存原始翻译
	originalZh := messagesZh
	originalEn := messagesEn
	defer func() {
		messagesZh = originalZh
		messagesEn = originalEn
		globalLocalizer = nil // 重置全局 localizer
	}()

	// 临时定义测试用翻译
	messagesZh = map[string]string{
		"test.hello":      "你好",
		"test.greeting":   "你好,%s!",
		"test.multi_args": "用户 %s 有 %d 条消息",
	}

	messagesEn = map[string]string{
		"test.hello":      "Hello",
		"test.greeting":   "Hello, %s!",
		"test.multi_args": "User %s has %d messages",
	}

	tests := []struct {
		name string
		lang string
		key  string
		args []interface{}
		want string
	}{
		{
			name: "中文简单翻译",
			lang: "zh",
			key:  "test.hello",
			want: "你好",
		},
		{
			name: "英文简单翻译",
			lang: "en",
			key:  "test.hello",
			want: "Hello",
		},
		{
			name: "中文带参数",
			lang: "zh",
			key:  "test.greeting",
			args: []interface{}{"张三"},
			want: "你好,张三!",
		},
		{
			name: "英文带参数",
			lang: "en",
			key:  "test.greeting",
			args: []interface{}{"Alice"},
			want: "Hello, Alice!",
		},
		{
			name: "多个参数",
			lang: "en",
			key:  "test.multi_args",
			args: []interface{}{"Alice", 5},
			want: "User Alice has 5 messages",
		},
		{
			name: "不存在的键 fallback",
			lang: "en",
			key:  "non.existent.key",
			want: "non.existent.key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLocalizer(tt.lang)
			got := l.T(tt.key, tt.args...)
			if got != tt.want {
				t.Errorf("Localizer.T() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalizer_Lang(t *testing.T) {
	tests := []struct {
		name string
		lang string
		want string
	}{
		{"中文", "zh", "zh"},
		{"英文", "en", "en"},
		{"不支持的语言fallback", "fr", "zh"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLocalizer(tt.lang)
			if got := l.Lang(); got != tt.want {
				t.Errorf("Localizer.Lang() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTranslationCompleteness 测试中英文翻译键的完整性
func TestTranslationCompleteness(t *testing.T) {
	// 获取所有键
	zhKeys := make(map[string]bool)
	enKeys := make(map[string]bool)

	for k := range messagesZh {
		zhKeys[k] = true
	}

	for k := range messagesEn {
		enKeys[k] = true
	}

	// 检查中文有但英文没有的键
	var missingInEn []string
	for k := range zhKeys {
		if !enKeys[k] {
			missingInEn = append(missingInEn, k)
		}
	}

	// 检查英文有但中文没有的键
	var missingInZh []string
	for k := range enKeys {
		if !zhKeys[k] {
			missingInZh = append(missingInZh, k)
		}
	}

	if len(missingInEn) > 0 {
		t.Errorf("英文翻译缺失以下键: %v", missingInEn)
	}

	if len(missingInZh) > 0 {
		t.Errorf("中文翻译缺失以下键: %v", missingInZh)
	}

	// 确保至少有一些翻译
	if len(zhKeys) == 0 {
		t.Error("中文翻译为空")
	}

	if len(enKeys) == 0 {
		t.Error("英文翻译为空")
	}
}
