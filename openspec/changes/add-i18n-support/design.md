# 国际化(i18n)技术设计文档

## 上下文

aicli 当前所有文本使用中文硬编码,包括 CLI 输出、LLM 提示词、错误信息等。为支持国际用户,需要添加国际化能力。

**约束条件:**
- 保持二进制体积轻量(Go 静态编译,避免资源文件)
- 不破坏现有功能和配置文件格式
- 保持性能(翻译查找 < 1ms)
- 简化实现,不依赖重型 i18n 框架

**利益相关者:**
- 国际用户:希望使用英文界面
- 中文用户:保持当前体验不变
- 开发者:易于维护和扩展翻译

## 目标 / 非目标

### 目标
- 支持中英双语自动切换
- 检测操作系统语言并自动适配
- 支持配置文件手动指定语言
- 翻译 CLI 界面、LLM 提示词、配置向导
- 轻量级实现,二进制体积增长 <100KB

### 非目标
- 不支持复数形式、日期格式化等复杂 i18n 特性
- 不支持运行时动态加载翻译资源
- 不支持中英文以外的其他语言(初期)
- 不翻译 LLM 返回的命令本身(命令是 shell 语法,与语言无关)
- 不翻译 README 和在线文档(可由社区贡献)

## 决策

### 决策 1: 轻量级自定义 i18n 实现

**选择:** 实现简单的 `pkg/i18n` 包,使用 Go map 存储翻译文本

**理由:**
- Go 标准库没有内置 i18n 支持
- 第三方库(如 `go-i18n`)功能过于复杂,会增加依赖和二进制体积
- aicli 的翻译需求简单(静态文本替换),不需要复数/格式化等高级特性
- 自定义实现可控性强,易于维护

**实现方式:**
```go
// pkg/i18n/i18n.go
type Localizer struct {
    lang     string
    messages map[string]string
}

func (l *Localizer) T(key string, args ...interface{}) string {
    msg := l.messages[key]
    if msg == "" {
        return key // fallback 到 key
    }
    return fmt.Sprintf(msg, args...)
}
```

### 决策 2: 语言检测优先级

**优先级顺序:**
1. 配置文件 `~/.aicli.json` 中的 `language` 字段(如果设置)
2. 操作系统环境变量 `LANG` 或 `LC_ALL`
3. 默认值:`zh`(保持向后兼容)

**理由:**
- 配置文件优先:允许用户显式覆盖
- 环境变量次之:自动化场景(CI/CD)可通过环境变量控制
- 默认中文:当前大部分用户为中文用户,保持兼容

**检测逻辑:**
```go
func DetectLanguage(cfg *config.Config) string {
    // 1. 配置文件优先
    if cfg.Language != "" {
        return normalizeLanguage(cfg.Language)
    }
    
    // 2. 环境变量
    lang := os.Getenv("LANG")
    if lang == "" {
        lang = os.Getenv("LC_ALL")
    }
    if lang != "" {
        return parseLocale(lang) // "zh_CN.UTF-8" -> "zh"
    }
    
    // 3. 默认中文
    return "zh"
}
```

### 决策 3: 翻译文本组织方式

**选择:** 按包分离翻译文件,每种语言一个文件

```
pkg/i18n/
├── i18n.go              # 核心逻辑
├── messages_zh.go       # 中文翻译
├── messages_en.go       # 英文翻译
└── keys.go              # 翻译键常量
```

**理由:**
- 按语言分文件:便于维护和审查
- 翻译键使用常量:避免拼写错误,IDE 自动补全
- 编译期打包:所有翻译编译到二进制中,无需外部资源文件

**翻译键命名规范:**
```go
// pkg/i18n/keys.go
const (
    ErrLoadConfig      = "error.load_config"
    ErrCreateProvider  = "error.create_provider"
    MsgHistoryEmpty    = "msg.history_empty"
    PromptConfirmRisky = "prompt.confirm_risky"
)
```

### 决策 4: LLM 提示词国际化策略

**选择:** 根据用户语言生成对应语言的 System Prompt

**理由:**
- LLM 对中文和英文提示词都能理解
- 使用用户语言可能提升命令生成质量(与用户输入语言一致)
- 避免混合语言(用户输入中文但 prompt 是英文)带来的理解偏差

**实现:**
```go
// pkg/llm/prompt.go
func GetSystemPrompt(ctx *ExecutionContext, lang string) string {
    if lang == "en" {
        return buildEnglishSystemPrompt(ctx)
    }
    return buildChineseSystemPrompt(ctx)
}
```

### 决策 5: 配置文件结构扩展

**选择:** 在 `Config` 结构中增加可选的 `Language` 字段

```go
// pkg/config/config.go
type Config struct {
    Version   string          `json:"version"`
    Language  string          `json:"language,omitempty"` // 新增: "zh" | "en" | ""
    LLM       LLMConfig       `json:"llm"`
    Execution ExecutionConfig `json:"execution"`
    Safety    SafetyConfig    `json:"safety"`
    // ...
}
```

**理由:**
- `omitempty` 标签:旧配置文件不包含此字段时不报错
- 空值表示自动检测:不强制用户设置
- 向后兼容:现有配置文件无需修改

## 考虑的替代方案

### 替代方案 1: 使用 go-i18n 库

**优点:**
- 成熟的 i18n 库,功能完善
- 支持复数、格式化等高级特性
- 社区维护,无需自己实现

**缺点:**
- 增加外部依赖,二进制体积增长 ~300KB
- 需要学习库的 API 和配置格式
- 功能过于复杂,aicli 用不到大部分特性

**拒绝理由:** 过度工程,不符合"简单优先"原则

### 替代方案 2: 命令行参数 `--lang` 指定语言

**优点:**
- 灵活,每次运行可以指定不同语言
- 便于测试

**缺点:**
- 用户体验差,每次都要手动指定
- 与自动检测冲突(用户期望自动适配)

**拒绝理由:** 用户体验不佳,大多数场景下自动检测即可

### 替代方案 3: 仅翻译 CLI 界面,LLM 提示词保持中文

**优点:**
- 实现简单,翻译量减少
- LLM 对中英文提示词理解能力相近

**缺点:**
- 用户语言与 prompt 语言不一致,可能影响理解
- 不够彻底,国际化不完整

**拒绝理由:** LLM 提示词国际化对用户体验提升明显,值得投入

## 风险 / 权衡

### 风险 1: 翻译质量和一致性

**风险:** 中英文翻译可能不准确或不一致

**缓解措施:**
- 编写翻译规范文档(`docs/i18n-guide.md`)
- Code Review 时检查翻译质量
- 收集用户反馈,持续改进翻译

### 风险 2: LLM 提示词切换影响命令生成质量

**风险:** 英文 prompt 可能导致命令生成质量下降(如果 LLM 对中文 prompt 训练更好)

**缓解措施:**
- 充分测试英文 prompt 的命令生成效果
- 提供配置选项强制使用中文 prompt(如果需要)
- 通过 Prompt Engineering 优化英文 prompt

### 风险 3: 二进制体积增长

**风险:** 翻译文本嵌入二进制,可能显著增加体积

**缓解措施:**
- 使用简洁的翻译文本,避免冗余
- 监控二进制体积,设定阈值(<100KB 增长)
- 如果体积超标,考虑压缩或懒加载方案

### 权衡: 初期仅支持中英双语

**权衡:** 不支持日语、韩语、法语等其他语言

**原因:**
- 减少初期工作量和翻译成本
- 中英文覆盖大部分目标用户
- 架构支持未来扩展(添加新语言仅需增加翻译文件)

## 迁移计划

### 阶段 1: 基础设施搭建(Week 1)
1. 实现 `pkg/i18n` 包核心逻辑
2. 添加语言检测功能
3. 编写单元测试

### 阶段 2: 翻译资源准备(Week 1-2)
1. 梳理所有需要翻译的文本
2. 编写 `messages_zh.go` 和 `messages_en.go`
3. 定义翻译键常量

### 阶段 3: 代码改造(Week 2-3)
1. 改造 `cmd/aicli/main.go` 和 `cmd/aicli/init.go`
2. 改造 `internal/app/app.go` 和 `internal/app/confirm.go`
3. 改造 `pkg/llm/prompt.go`
4. 更新配置结构和默认值

### 阶段 4: 测试和优化(Week 3)
1. 集成测试(中英文环境下完整流程)
2. 性能测试(确保无明显性能下降)
3. 二进制体积检查
4. 修复 Bug 和优化翻译

### 回滚计划
- 如果国际化引入严重 Bug,可通过配置文件强制中文
- 代码层面保留 fallback 机制(翻译失败时显示原始 key)
- Git 分支管理,可快速回退

## 待决问题

1. **是否需要在 CLI 帮助信息中显示当前语言?**
   - 例如 `aicli --help` 末尾显示 "Language: en"
   - 决策:暂不实现,减少干扰

2. **是否支持通过环境变量 `AICLI_LANG` 覆盖系统语言?**
   - 方便测试和特殊场景
   - 决策:暂不实现,优先级低

3. **LLM 返回的错误信息是否翻译?**
   - 例如 OpenAI API 返回的错误信息
   - 决策:不翻译,这些是 API 原始错误,保持原样便于调试

4. **是否需要 i18n 文档和贡献指南?**
   - 方便社区贡献其他语言翻译
   - 决策:后续补充,初期先完成核心功能
