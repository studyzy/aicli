# aicli 架构设计文档

> 版本: 1.0 | 日期: 2026-01-13

## 概述

aicli 是一个基于 Go 语言开发的自然语言命令行工具，允许用户使用自然语言描述意图，通过大语言模型（LLM）服务将其转换为实际的 shell 命令并执行。

### 设计目标

- **简单易用**: 用户只需输入自然语言，无需记忆复杂的命令语法
- **安全可靠**: 检测危险命令并要求用户确认
- **灵活可扩展**: 支持多种 LLM 提供商，易于添加新功能
- **高性能**: 命令转换和执行总时间 < 5 秒
- **跨平台**: 支持 Linux、macOS、Windows

## 系统架构

### 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         用户界面层                            │
│                     (cmd/aicli/main.go)                      │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  命令行解析 (Cobra)                                    │   │
│  │  - 参数解析                                            │   │
│  │  - 标志处理 (--verbose, --dry-run, --force)           │   │
│  │  - 子命令 (--history, --retry)                        │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                        应用逻辑层                             │
│                    (internal/app/app.go)                     │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  应用主逻辑                                            │   │
│  │  - 输入验证                                            │   │
│  │  - 上下文构建                                          │   │
│  │  - 流程编排                                            │   │
│  │  - 错误处理                                            │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
           ↓              ↓              ↓              ↓
    ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
    │   LLM    │   │  执行器   │   │  安全检查 │   │   历史   │
    │  Provider│   │ Executor │   │  Checker │   │  History │
    └──────────┘   └──────────┘   └──────────┘   └──────────┘
         ↓              ↓              ↓              ↓
┌─────────────────────────────────────────────────────────────┐
│                         核心服务层                            │
│                         (pkg/*)                              │
│                                                               │
│  ┌────────────────┐  ┌────────────────┐  ┌───────────────┐ │
│  │  LLM 服务       │  │  命令执行服务   │  │  安全检查服务  │ │
│  │  pkg/llm/      │  │  pkg/executor/ │  │  pkg/safety/  │ │
│  │                │  │                │  │               │ │
│  │  - OpenAI      │  │  - Shell 检测  │  │  - 模式匹配   │ │
│  │  - Anthropic   │  │  - 命令执行    │  │  - 风险评估   │ │
│  │  - Local Model │  │  - IO 处理     │  │  - 确认提示   │ │
│  │  - 工厂函数     │  │                │  │               │ │
│  └────────────────┘  └────────────────┘  └───────────────┘ │
│                                                               │
│  ┌────────────────┐  ┌────────────────┐                     │
│  │  配置管理       │  │  历史记录       │                     │
│  │  pkg/config/   │  │  internal/     │                     │
│  │                │  │  history/      │                     │
│  │  - 加载/保存    │  │  - 增删查改    │                     │
│  │  - 默认值      │  │  - 持久化      │                     │
│  │  - 验证        │  │  - 搜索过滤    │                     │
│  └────────────────┘  └────────────────┘                     │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                        外部服务层                             │
│                                                               │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐            │
│  │  OpenAI    │  │ Anthropic  │  │   Ollama   │            │
│  │    API     │  │    API     │  │  (本地)     │            │
│  └────────────┘  └────────────┘  └────────────┘            │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  操作系统 Shell                                          │ │
│  │  bash / zsh / PowerShell / cmd                          │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 模块设计

### 1. 用户界面层 (cmd/aicli)

**职责**: 
- 解析命令行参数和标志
- 处理用户输入（自然语言、stdin）
- 调用应用逻辑层
- 格式化输出结果

**关键组件**:
- `main.go`: 程序入口，Cobra 命令定义
- `rootCmd`: 根命令，处理自然语言输入
- 子命令: `--history`, `--retry`

### 2. 应用逻辑层 (internal/app)

**职责**:
- 编排整个执行流程
- 构建执行上下文
- 协调各个核心服务
- 处理业务逻辑错误

**关键组件**:
- `App`: 应用主结构体
- `Run()`: 主执行逻辑
- `Flags`: 命令行标志定义
- `confirm.go`: 确认提示逻辑
- `io.go`: 输入输出处理

**执行流程**:
```
1. 验证输入
2. 构建执行上下文（OS、Shell、WorkDir、Stdin）
3. 调用 LLM 翻译命令
4. 安全检查（危险命令检测）
5. 用户确认（如需要）
6. 执行命令
7. 保存历史记录
8. 返回结果
```

### 3. LLM 服务层 (pkg/llm)

**职责**:
- 定义 LLM 提供商接口
- 实现多种 LLM 提供商
- 构建和管理 Prompt
- 处理 LLM 响应

**关键组件**:
- `LLMProvider` 接口: 统一的提供商接口
- `OpenAIProvider`: OpenAI GPT 系列实现
- `AnthropicProvider`: Anthropic Claude 系列实现
- `LocalModelProvider`: 本地模型（Ollama）实现
- `NewProvider()`: 工厂函数，根据配置创建提供商
- `BuildPrompt()`: 构建提示词
- `cleanCommand()`: 清理命令输出

**接口定义**:
```go
type LLMProvider interface {
    Translate(ctx context.Context, input string, execCtx *ExecutionContext) (string, error)
    Name() string
}
```

### 4. 命令执行层 (pkg/executor)

**职责**:
- 检测系统 Shell
- 执行 Shell 命令
- 处理 stdin/stdout/stderr
- 管理命令超时

**关键组件**:
- `Executor`: 命令执行器
- `ShellAdapter`: Shell 适配器
- `DetectShell()`: Shell 检测
- `Execute()`: 命令执行

**Shell 支持**:
- Linux/macOS: bash, zsh, sh
- Windows: PowerShell, cmd

### 5. 安全检查层 (pkg/safety)

**职责**:
- 定义危险命令模式
- 检测危险操作
- 评估风险等级
- 提供风险描述

**关键组件**:
- `SafetyChecker`: 安全检查器
- `DangerousPatterns`: 危险模式列表
- `IsDangerous()`: 危险检测方法

**检测模式**:
- 文件删除: `rm -rf`, `del /S`
- 格式化: `mkfs`, `format`
- 权限操作: `chmod 777`, `chown`
- 网络危险: `curl | sh`, `wget | bash`
- 系统修改: `sudo`, `dd if=`

### 6. 配置管理层 (pkg/config)

**职责**:
- 加载和保存配置
- 提供默认配置
- 验证配置项

**配置文件**: `~/.aicli.json`

**配置结构**:
```go
type Config struct {
    Version   string
    LLM       LLMConfig       // LLM 配置
    Execution ExecutionConfig // 执行配置
    Safety    SafetyConfig    // 安全配置
    History   HistoryConfig   // 历史配置
    Logging   LoggingConfig   // 日志配置
}
```

### 7. 历史记录层 (internal/history)

**职责**:
- 记录命令执行历史
- 持久化到文件
- 提供查询和检索
- 支持命令重试

**关键功能**:
- `Add()`: 添加历史记录
- `List()`: 列出所有记录
- `Get()`: 获取指定记录
- `Search()`: 搜索记录
- `Save()/Load()`: 持久化

## 数据流

### 命令转换与执行流程

```
用户输入 "列出当前目录的所有txt文件"
    ↓
[1. 输入验证]
    ↓
[2. 构建执行上下文]
    OS: linux, Shell: bash, WorkDir: /home/user
    ↓
[3. 调用 LLM Provider]
    System Prompt: "You are a command-line expert..."
    User Input: "列出当前目录的所有txt文件"
    Context: {OS, Shell, WorkDir}
    ↓
[4. LLM 返回命令]
    "ls *.txt"
    ↓
[5. 安全检查]
    IsDangerous("ls *.txt") → false
    ↓
[6. 执行命令]
    Executor.Execute("ls *.txt")
    ↓
[7. 保存历史记录]
    {Input, Command, Success, Output}
    ↓
[8. 返回结果]
    file1.txt
    file2.txt
    notes.txt
```

### 管道模式数据流

```
cat file.txt | aicli "统计行数"
    ↓
[读取 stdin]
    "line1\nline2\nline3"
    ↓
[构建上下文 with Stdin]
    {OS, Shell, WorkDir, Stdin: "..."}
    ↓
[LLM 转换]
    "wc -l"
    ↓
[执行命令 with stdin]
    3
```

## 错误处理

### 错误类型

1. **配置错误**
   - 配置文件不存在 → 使用默认配置
   - API Key 缺失 → 返回明确错误
   - 配置格式错误 → 返回解析错误

2. **LLM 错误**
   - 网络超时 → 重试或提示用户
   - API 错误 → 显示错误消息
   - 返回空命令 → 提示用户重新描述

3. **执行错误**
   - 命令不存在 → 显示错误信息
   - 权限不足 → 提示用户
   - 超时 → 可配置的超时时间

4. **安全错误**
   - 危险命令未确认 → 拒绝执行
   - 管道模式危险命令 → 需要 --force

### 错误传播

```
底层错误 → 包装错误 → 用户友好错误
fmt.Errorf("wrap: %w", err)
```

## 性能优化

### 目标

- LLM API 调用 + 命令执行 < 5 秒
- 配置文件加载 < 10ms
- 命令执行启动 < 100ms

### 优化策略

1. **HTTP 连接复用**: 使用长连接
2. **超时控制**: Context 超时管理
3. **并发控制**: 使用 goroutine 但避免过度并发
4. **缓存**: 配置加载缓存（未实现）

## 安全性

### 安全机制

1. **危险命令检测**: 正则表达式匹配
2. **用户确认**: 交互式确认提示
3. **强制执行**: `--force` 标志跳过确认
4. **隐私保护**: `--no-send-stdin` 不发送敏感数据
5. **日志脱敏**: 不记录完整 API Key

### 安全最佳实践

- 配置文件权限: 建议 600
- API Key 管理: 使用环境变量
- 命令审查: 执行前显示命令

## 可扩展性

### 添加新 LLM Provider

1. 实现 `LLMProvider` 接口
2. 在 `factory.go` 中注册
3. 添加测试
4. 更新文档

```go
type MyProvider struct { ... }

func (p *MyProvider) Translate(...) (string, error) {
    // 实现逻辑
}

func (p *MyProvider) Name() string {
    return "myprovider"
}
```

### 添加新功能

遵循模块化设计原则：
- 新功能放在独立的包中
- 定义清晰的接口
- 编写单元测试
- 更新文档

## 依赖关系

### 外部依赖

- `github.com/spf13/cobra`: CLI 框架
- Go 标准库: `net/http`, `encoding/json`, `os/exec`

### 内部依赖

```
cmd/aicli
    ↓
internal/app
    ↓
pkg/llm, pkg/executor, pkg/safety, pkg/config
```

## 测试策略

### 测试层次

1. **单元测试**: 每个包独立测试
2. **集成测试**: 跨包测试（tests/integration）
3. **E2E 测试**: 完整流程测试（使用 Mock LLM）

### Mock 策略

- `MockLLMProvider`: Mock LLM 响应
- `httptest.Server`: Mock HTTP API
- `os.Pipe()`: Mock stdin/stdout

### 覆盖率目标

- 整体: ≥65%
- 核心包 (llm, executor, safety): ≥80%

## 部署架构

### 单二进制部署

```
aicli (可执行文件)
~/.aicli.json (配置)
~/.aicli_history.json (历史)
```

### 跨平台支持

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## 未来规划

1. **缓存机制**: 缓存常见命令转换
2. **插件系统**: 支持自定义扩展
3. **Web UI**: 可选的 Web 界面
4. **团队协作**: 共享历史和最佳实践
5. **智能建议**: 基于历史的命令建议

## 总结

aicli 采用清晰的分层架构，核心是 LLM Provider 抽象层，通过工厂模式支持多种 LLM 服务。安全检查、命令执行、历史记录等模块各司其职，保证了系统的可维护性和可扩展性。
