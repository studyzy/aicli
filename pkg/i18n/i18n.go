// Package i18n 提供国际化(i18n)支持
package i18n

import (
	"fmt"
	"os"
	"strings"

	"github.com/studyzy/aicli/pkg/config"
)

// Localizer 国际化器,负责翻译文本
type Localizer struct {
	lang     string
	messages map[string]string
}

// globalLocalizer 全局 Localizer 实例
var globalLocalizer *Localizer

// Init 初始化全局 Localizer
func Init(cfg *config.Config) {
	lang := DetectLanguage(cfg)
	globalLocalizer = NewLocalizer(lang)
}

// NewLocalizer 创建一个新的 Localizer
func NewLocalizer(lang string) *Localizer {
	lang = normalizeLanguage(lang)

	var messages map[string]string
	switch lang {
	case "en":
		messages = messagesEn
	case "zh":
		messages = messagesZh
	default:
		// fallback 到中文
		messages = messagesZh
	}

	return &Localizer{
		lang:     lang,
		messages: messages,
	}
}

// T 翻译文本(Translate)
// key: 翻译键
// args: 格式化参数(使用 fmt.Sprintf 风格)
// 返回翻译后的文本,如果找不到翻译则返回 key 本身
func (l *Localizer) T(key string, args ...interface{}) string {
	msg, ok := l.messages[key]
	if !ok {
		// fallback: 返回 key 本身
		return key
	}

	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}

	return msg
}

// Lang 返回当前语言代码
func (l *Localizer) Lang() string {
	return l.lang
}

// T 使用全局 Localizer 翻译文本
func T(key string, args ...interface{}) string {
	if globalLocalizer == nil {
		// fallback: 返回 key
		return key
	}
	return globalLocalizer.T(key, args...)
}

// Lang 返回全局 Localizer 的语言代码
func Lang() string {
	if globalLocalizer == nil {
		return "zh"
	}
	return globalLocalizer.Lang()
}

// DetectLanguage 检测用户语言偏好
// 优先级: 配置文件 > 环境变量 > 默认值(zh)
func DetectLanguage(cfg *config.Config) string {
	// 1. 配置文件优先
	if cfg != nil && cfg.Language != "" {
		return normalizeLanguage(cfg.Language)
	}

	// 2. 环境变量
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}
	if lang != "" {
		return parseLocale(lang)
	}

	// 3. 默认中文
	return "zh"
}

// parseLocale 解析 locale 字符串,提取语言代码
// 例如: "zh_CN.UTF-8" -> "zh", "en_US.UTF-8" -> "en"
func parseLocale(locale string) string {
	// 移除编码后缀 (.UTF-8, .utf8 等)
	if idx := strings.Index(locale, "."); idx > 0 {
		locale = locale[:idx]
	}

	// 提取语言代码(下划线前的部分)
	parts := strings.Split(locale, "_")
	if len(parts) > 0 {
		return normalizeLanguage(parts[0])
	}

	return "zh" // fallback
}

// normalizeLanguage 规范化语言代码
// 支持: zh, zh-CN, zh_CN -> zh
//       en, en-US, en_US -> en
func normalizeLanguage(lang string) string {
	lang = strings.ToLower(lang)
	lang = strings.TrimSpace(lang)

	// 提取主语言代码
	if idx := strings.IndexAny(lang, "-_"); idx > 0 {
		lang = lang[:idx]
	}

	// 只支持 zh 和 en
	if lang == "zh" || lang == "en" {
		return lang
	}

	// 不支持的语言 fallback 到中文
	return "zh"
}
