# 贡献指南

感谢你对 aicli 项目感兴趣！我们欢迎所有形式的贡献，无论是代码、文档、Bug 报告还是功能建议。

## 目录

- [行为准则](#行为准则)
- [如何贡献](#如何贡献)
- [开发流程](#开发流程)
- [代码规范](#代码规范)
- [提交规范](#提交规范)
- [测试要求](#测试要求)
- [文档要求](#文档要求)

## 行为准则

### 我们的承诺

为了营造开放和友好的环境，我们承诺让每个人都能无障碍地参与这个项目，无论其经验水平、性别、性别认同、性取向、残疾、外貌、体型、种族、民族、年龄、宗教或国籍。

### 我们的标准

**积极行为示例**:
- 使用友好和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 专注于对社区最有利的事情
- 对其他社区成员表示同情

**不可接受的行为**:
- 使用性化语言或图像，以及不受欢迎的性关注或性骚扰
- 发表攻击性、侮辱性或贬损性评论，以及人身或政治攻击
- 公开或私下骚扰
- 未经明确许可，发布他人的私人信息
- 其他可能被合理认为在专业环境中不适当的行为

## 如何贡献

### 报告 Bug

在创建 Bug 报告之前，请先搜索现有的 Issues，避免重复报告。

**创建 Bug 报告时，请包含**:
- **简洁的标题**: 清楚地描述问题
- **详细描述**: 提供尽可能多的相关信息
- **重现步骤**: 详细说明如何重现问题
  1. 第一步
  2. 第二步
  3. ...
- **期望行为**: 描述你期望发生什么
- **实际行为**: 描述实际发生了什么
- **屏幕截图**: 如果适用，添加屏幕截图
- **环境信息**:
  - OS: [例如 macOS 13.0]
  - Go 版本: [例如 1.21.0]
  - aicli 版本: [例如 0.1.0]
  - LLM Provider: [例如 OpenAI]

**Bug 报告模板**:
```markdown
**描述**
简洁清晰地描述 bug。

**重现步骤**
1. 执行命令 '...'
2. 输入 '...'
3. 查看错误

**期望行为**
应该输出正确的命令。

**实际行为**
返回错误：...

**环境**
- OS: macOS 13.0
- Go 版本: 1.21.0
- aicli 版本: 0.1.0

**额外上下文**
添加任何其他相关信息。
```

### 建议新功能

**创建功能请求时，请包含**:
- **功能描述**: 清楚地描述你想要的功能
- **使用场景**: 解释为什么这个功能有用
- **可能的实现**: 如果有想法，描述如何实现
- **替代方案**: 是否考虑过其他解决方案

### 提交代码

1. **Fork 项目**
2. **创建特性分支** (`git checkout -b feature/AmazingFeature`)
3. **提交更改** (`git commit -m 'feat: 添加某某功能'`)
4. **推送到分支** (`git push origin feature/AmazingFeature`)
5. **开启 Pull Request**

## 开发流程

### 1. 设置开发环境

```bash
# 克隆你的 Fork
git clone https://github.com/your-username/aicli.git
cd aicli

# 添加上游仓库
git remote add upstream https://github.com/studyzy/aicli.git

# 安装依赖
go mod download
```

### 2. 创建特性分支

```bash
# 从 main 分支创建新分支
git checkout main
git pull upstream main
git checkout -b feature/my-new-feature
```

**分支命名规范**:
- `feature/feature-name` - 新功能
- `fix/bug-description` - Bug 修复
- `docs/doc-update` - 文档更新
- `refactor/refactor-description` - 代码重构
- `test/test-description` - 测试相关

### 3. 编写代码

遵循[代码规范](#代码规范)编写高质量代码。

### 4. 编写测试

**TDD 流程**:
1. 先编写测试（测试应该失败）
2. 编写实现代码（使测试通过）
3. 重构代码（保持测试通过）

```bash
# 运行测试
go test ./...

# 查看覆盖率
go test -cover ./...
```

### 5. 运行代码检查

```bash
# 格式化代码
go fmt ./...

# 运行 linter
golangci-lint run

# 检查所有测试通过
go test ./...
```

### 6. 提交更改

遵循[提交规范](#提交规范)提交代码。

### 7. 推送并创建 PR

```bash
git push origin feature/my-new-feature
```

在 GitHub 上创建 Pull Request。

### 8. 响应审查

- 及时响应审查意见
- 根据建议修改代码
- 保持 PR 范围小而专注

## 代码规范

### Go 代码规范

遵循官方 [Effective Go](https://go.dev/doc/effective_go) 指南。

**关键原则**:
- 代码应该简单、清晰、易读
- 优先使用标准库
- 避免过度设计
- 编写可测试的代码

### 命名规范

```go
// ✅ 好的命名
type LLMProvider interface { ... }
func NewOpenAIProvider(...) { ... }
var defaultTimeout = 10

// ❌ 不好的命名
type Llm_provider interface { ... }
func new_openai_provider(...) { ... }
var DEFAULT_TIMEOUT = 10
```

### 注释规范

**所有导出的内容必须有中文注释**:

```go
// Config 表示应用配置
type Config struct {
    // LLM LLM 服务配置
    LLM LLMConfig `json:"llm"`
}

// Load 从文件加载配置
// path: 配置文件路径
// 返回: 配置对象和可能的错误
func Load(path string) (*Config, error) {
    // 实现...
}
```

### 错误处理

```go
// ✅ 包装错误，提供上下文
if err := doSomething(); err != nil {
    return fmt.Errorf("执行某操作失败: %w", err)
}

// ❌ 直接返回错误，丢失上下文
if err := doSomething(); err != nil {
    return err
}
```

### 代码格式化

```bash
# 格式化所有代码
go fmt ./...

# 或使用 goimports
goimports -w .
```

## 提交规范

### Conventional Commits

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范。

**格式**:
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建、工具等
- `ci`: CI 配置

**示例**:

```
feat(llm): 添加 Anthropic Provider 支持

实现了 Anthropic Claude API 的集成，包括：
- AnthropicProvider 实现
- 完整的单元测试
- 文档更新

Closes #123
```

```
fix(executor): 修复 Windows 下 Shell 检测问题

在 Windows 系统上，Shell 检测逻辑错误导致无法执行命令。
修复了检测逻辑，现在可以正确识别 PowerShell 和 CMD。

Fixes #456
```

### 提交消息最佳实践

- **使用中文**: 提交消息使用中文
- **现在时态**: "添加功能" 而非 "添加了功能"
- **简洁明了**: 第一行不超过 50 个字符
- **详细描述**: Body 部分提供详细说明
- **引用 Issue**: 如果相关，引用 Issue 编号

## 测试要求

### 测试覆盖率

**最低要求**:
- 新代码：80%+
- 核心包 (llm, executor, safety): 80%+
- 整体项目: 65%+

### 测试类型

**单元测试**:
```go
func TestNewConfig(t *testing.T) {
    cfg := NewConfig()
    if cfg == nil {
        t.Fatal("期望非 nil 配置")
    }
}
```

**表格驱动测试**:
```go
func TestIsDangerous(t *testing.T) {
    tests := []struct {
        name     string
        command  string
        expected bool
    }{
        {"安全命令", "ls -la", false},
        {"危险命令", "rm -rf /", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := IsDangerous(tt.command)
            if result != tt.expected {
                t.Errorf("期望 %v, 实际 %v", tt.expected, result)
            }
        })
    }
}
```

**Mock 测试**:
```go
func TestWithMock(t *testing.T) {
    mockLLM := &MockLLMProvider{
        TranslateFn: func(input string) string {
            return "echo test"
        },
    }
    // 使用 mock 进行测试
}
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行指定包
go test ./pkg/llm

# 详细输出
go test -v ./...

# 查看覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 文档要求

### 代码文档

**所有导出的函数、类型、常量必须有中文注释**:

```go
// LLMProvider 定义 LLM 服务提供商的接口
//
// 所有 LLM 提供商必须实现此接口以集成到 aicli 中。
type LLMProvider interface {
    // Translate 将自然语言转换为命令
    //
    // 参数:
    //   ctx: 上下文，用于超时控制
    //   input: 用户的自然语言描述
    //   execCtx: 执行上下文信息（OS、Shell 等）
    //
    // 返回:
    //   转换后的命令字符串和可能的错误
    Translate(ctx context.Context, input string, execCtx *ExecutionContext) (string, error)
}
```

### 文档更新

**添加新功能时**:
1. 更新 `README.md`（如果用户可见）
2. 更新相关的 `docs/` 文档
3. 添加使用示例
4. 更新 `example.aicli.json`（如果需要）

### 文档风格

- 使用中文
- 简洁明了
- 提供代码示例
- 包含实际用例

## Pull Request 检查清单

在提交 PR 之前，请确保：

- [ ] 代码遵循项目的代码规范
- [ ] 所有测试通过 (`go test ./...`)
- [ ] 测试覆盖率达标（新代码 80%+）
- [ ] 代码已格式化 (`go fmt ./...`)
- [ ] Linter 检查通过 (`golangci-lint run`)
- [ ] 所有导出的内容有中文注释
- [ ] 提交消息遵循规范
- [ ] 相关文档已更新
- [ ] PR 描述清楚说明了更改内容
- [ ] 已经在本地测试过功能

## Pull Request 模板

```markdown
## 更改类型
- [ ] 新功能
- [ ] Bug 修复
- [ ] 文档更新
- [ ] 代码重构
- [ ] 性能优化
- [ ] 其他 (请说明):

## 描述
简洁清晰地描述这个 PR 的内容。

## 相关 Issue
Closes #123

## 测试
描述你如何测试了这些更改。

## 检查清单
- [ ] 代码遵循项目规范
- [ ] 所有测试通过
- [ ] 测试覆盖率达标
- [ ] 文档已更新
- [ ] 提交消息符合规范

## 屏幕截图（如适用）
添加屏幕截图帮助说明更改。

## 额外说明
添加任何其他相关信息。
```

## 发布流程

维护者负责发布流程：

1. 更新版本号
2. 更新 CHANGELOG
3. 创建 Git 标签
4. 构建发布版本
5. 创建 GitHub Release
6. 发布公告

## 获取帮助

如果你有任何问题：

- 查看 [开发指南](docs/development.md)
- 搜索现有 [Issues](https://github.com/studyzy/aicli/issues)
- 在 [Discussions](https://github.com/studyzy/aicli/discussions) 提问
- 联系维护者

## 致谢

感谢所有贡献者！你们让 aicli 变得更好。

## 许可证

通过贡献代码，你同意你的贡献将在 MIT 许可证下授权。
