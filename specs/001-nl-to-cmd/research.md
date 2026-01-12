# 研究文档: 自然语言命令行转换工具

**功能**: 001-nl-to-cmd
**日期**: 2026-01-12
**目的**: 解决技术选型和实施方案中的未知项

## 研究任务概览

基于技术背景中的待研究项和功能需求，需要解决以下关键技术决策：

1. LLM SDK 选择和集成方案
2. 跨平台 Shell 命令执行方案
3. 危险命令检测模式定义
4. 测试策略（Mock LLM响应）
5. 配置文件格式和管理方案

---

## 1. LLM SDK 选择和集成方案

### 决策: 使用直接 HTTP 调用 + 官方 SDK（按需）

**理由**:
- **灵活性**: 直接 HTTP 调用允许支持任何遵循 OpenAI/Anthropic API 规范的 LLM 服务
- **轻量级**: 避免引入多个重量级 SDK，减少依赖和二进制大小
- **可测试性**: HTTP 客户端易于 Mock 和测试
- **扩展性**: 可以轻松添加对新 LLM 提供商的支持

**评估的替代方案**:
| 方案 | 优点 | 缺点 | 结论 |
|------|------|------|------|
| OpenAI官方SDK (`github.com/sashabaranov/go-openai`) | 类型安全，维护良好 | 仅支持OpenAI | 可选使用 |
| Anthropic官方SDK | 官方支持 | 依赖较新，稳定性待验证 | 可选使用 |
| LangChain Go | 功能丰富 | 过于复杂，引入大量依赖 | ❌ 不采用 |
| 直接HTTP调用 | 完全控制，最灵活 | 需要手动处理协议细节 | ✅ 主要方案 |

**实施方案**:
```go
// pkg/llm/provider.go - 定义统一接口
type LLMProvider interface {
    // Translate 将自然语言转换为命令
    // context: 上下文信息（当前目录、操作系统等）
    // input: 自然语言描述
    // 返回: 转换后的命令字符串
    Translate(ctx context.Context, input string, context ExecutionContext) (string, error)
}

type ExecutionContext struct {
    OS          string   // 操作系统（linux/darwin/windows）
    Shell       string   // Shell类型（bash/zsh/powershell）
    WorkingDir  string   // 当前工作目录
    StdinData   string   // 标准输入数据（如果有）
}

// 实现示例：OpenAI
type OpenAIProvider struct {
    apiKey  string
    model   string
    baseURL string
    client  *http.Client
}
```

**依赖选择**:
- 核心：`net/http`（标准库）
- JSON处理：`encoding/json`（标准库）
- 可选：`github.com/sashabaranov/go-openai` v1.20+（如果用户选择OpenAI）

---

## 2. 跨平台 Shell 命令执行方案

### 决策: 检测系统Shell + 适配器模式

**理由**:
- 不同操作系统有不同的默认 Shell（bash/zsh/PowerShell/cmd.exe）
- 命令语法存在差异（如文件路径分隔符、环境变量引用）
- 需要在 LLM Prompt 中指定目标 Shell，以获得正确的命令语法

**Shell 检测策略**:
| 操作系统 | 检测方法 | 默认Shell | 备选Shell |
|----------|----------|-----------|-----------|
| Linux | `$SHELL` 环境变量 | bash | zsh, sh |
| macOS | `$SHELL` 环境变量 | zsh (10.15+) | bash |
| Windows | 检查PowerShell版本 | powershell | cmd.exe |

**实施方案**:
```go
// pkg/executor/shell.go
type ShellAdapter struct {
    Type    ShellType  // bash, zsh, powershell, cmd
    Path    string     // Shell可执行文件路径
    Args    []string   // 执行命令的参数模板
}

type ShellType string
const (
    ShellBash       ShellType = "bash"
    ShellZsh        ShellType = "zsh"
    ShellPowerShell ShellType = "powershell"
    ShellCmd        ShellType = "cmd"
)

// DetectShell 检测当前系统的Shell
func DetectShell() (*ShellAdapter, error) {
    // 1. 读取 $SHELL 环境变量（Unix-like）
    // 2. 检测 PowerShell（Windows）
    // 3. 回退到系统默认
}

// Execute 执行命令
func (s *ShellAdapter) Execute(command string) (stdout, stderr string, err error) {
    cmd := exec.Command(s.Path, append(s.Args, command)...)
    // 设置 Stdin/Stdout/Stderr
    // 执行并返回结果
}
```

**命令执行安全考虑**:
- 使用 `os/exec.Command` 而不是 `sh -c`（避免命令注入）
- 设置执行超时（默认 30 秒）
- 支持用户中断（context.Context 取消）

---

## 3. 危险命令检测模式

### 决策: 关键词+模式匹配的混合策略

**理由**:
- 需要在用户确认前检测到潜在危险操作
- 100%准确率要求（SC-006）意味着宁可误报也不能漏报
- 模式需要覆盖常见的危险操作

**危险命令分类**:
| 类别 | 关键词/模式 | 风险等级 | 示例 |
|------|-------------|----------|------|
| 文件删除 | `rm`, `del`, `Remove-Item` | 高 | `rm -rf /` |
| 格式化 | `mkfs`, `format` | 极高 | `mkfs.ext4 /dev/sda1` |
| 系统修改 | `chmod 777`, `chown`, `sudo` | 高 | `sudo rm -rf /*` |
| 网络危险 | `wget ... \| sh`, `curl ... \| bash` | 高 | 执行远程脚本 |
| 权限提升 | `sudo`, `su`, `runas` | 中 | `sudo vim /etc/passwd` |
| 批量操作 | `rm *`, `del /S` | 高 | 删除大量文件 |

**实施方案**:
```go
// pkg/safety/patterns.go
var DangerousPatterns = []Pattern{
    {
        Regex:       regexp.MustCompile(`rm\s+.*-[rf]`),
        Description: "文件删除操作",
        Level:       RiskHigh,
    },
    {
        Regex:       regexp.MustCompile(`(curl|wget).*\|\s*(ba)?sh`),
        Description: "从网络下载并执行脚本",
        Level:       RiskHigh,
    },
    // ... 更多模式
}

// pkg/safety/checker.go
type SafetyChecker struct {
    patterns []Pattern
}

func (c *SafetyChecker) IsDangerous(command string) (bool, string) {
    // 检查每个模式
    // 返回: 是否危险, 原因描述
}
```

**误报处理**:
- 提供 `--force` 标志跳过确认（用户自行承担风险）
- 在确认提示中显示具体的危险原因

---

## 4. 测试策略（Mock LLM响应）

### 决策: 接口Mock + 测试夹具（Fixtures）

**理由**:
- LLM API 调用昂贵且不稳定，不适合单元测试
- 需要测试各种边界情况（超时、错误、不同响应）
- 测试必须可重现和快速

**测试方法**:
```go
// pkg/llm/mock.go - 用于测试
type MockLLMProvider struct {
    TranslateFunc func(ctx context.Context, input string, context ExecutionContext) (string, error)
}

func (m *MockLLMProvider) Translate(ctx context.Context, input string, context ExecutionContext) (string, error) {
    if m.TranslateFunc != nil {
        return m.TranslateFunc(ctx, input, context)
    }
    return "", errors.New("not implemented")
}

// 测试示例
func TestExecutor_ExecuteWithMock(t *testing.T) {
    mockLLM := &MockLLMProvider{
        TranslateFunc: func(ctx context.Context, input string, context ExecutionContext) (string, error) {
            if input == "列出所有txt文件" {
                return "ls *.txt", nil
            }
            return "", errors.New("unknown input")
        },
    }
    
    executor := NewExecutor(mockLLM, nil)
    result, err := executor.Execute(context.Background(), "列出所有txt文件")
    assert.NoError(t, err)
    assert.Contains(t, result, ".txt")
}
```

**测试覆盖率目标**:
- `pkg/llm`: 80%（核心业务逻辑）
- `pkg/executor`: 80%（核心业务逻辑）
- `pkg/config`: 70%
- `pkg/safety`: 85%（安全关键）
- `internal/app`: 60%
- 整体: 65%+（超过章程要求的60%）

**集成测试策略**:
```bash
# tests/integration/basic_test.go
# 使用真实的Shell执行简单命令，但Mock LLM
# 测试端到端流程（除了真实LLM调用）
```

---

## 5. 配置文件格式和管理方案

### 决策: JSON格式 + `~/.aicli.json`

**理由**:
- JSON 易于解析，Go 标准库原生支持
- 用户熟悉 JSON 格式，易于手动编辑
- 支持嵌套结构，便于扩展

**配置文件结构**:
```json
{
  "version": "1.0",
  "llm": {
    "provider": "openai",
    "api_key": "sk-xxxxxxxxxxxxx",
    "api_base": "https://api.openai.com/v1",
    "model": "gpt-4",
    "timeout": 10,
    "max_tokens": 500
  },
  "execution": {
    "auto_confirm": false,
    "dry_run_default": false,
    "timeout": 30,
    "shell": "auto"
  },
  "safety": {
    "enable_checks": true,
    "dangerous_patterns": ["rm -rf", "format"],
    "require_confirmation": true
  },
  "history": {
    "enabled": true,
    "max_entries": 1000,
    "file": "~/.aicli_history.json"
  },
  "logging": {
    "enabled": false,
    "level": "info",
    "file": ""
  }
}
```

**配置管理实施**:
```go
// pkg/config/config.go
type Config struct {
    Version   string          `json:"version"`
    LLM       LLMConfig       `json:"llm"`
    Execution ExecutionConfig `json:"execution"`
    Safety    SafetyConfig    `json:"safety"`
    History   HistoryConfig   `json:"history"`
    Logging   LoggingConfig   `json:"logging"`
}

// Load 从文件加载配置
func Load(path string) (*Config, error) {
    // 1. 展开 ~ 路径
    // 2. 读取文件
    // 3. 解析 JSON
    // 4. 验证必填字段
    // 5. 应用默认值
}

// Save 保存配置到文件
func (c *Config) Save(path string) error {
    // 序列化为 JSON 并写入文件
}
```

**配置优先级**:
1. 命令行参数（`--api-key`, `--model` 等）
2. 环境变量（`AICLI_API_KEY`, `AICLI_MODEL` 等）
3. 配置文件 `~/.aicli.json`
4. 默认值

**敏感信息处理**:
- API Key 存储在配置文件中（用户负责文件权限）
- 提供环境变量方式（`AICLI_API_KEY`）避免硬编码
- 日志中不记录完整的 API Key

---

## 6. LLM Prompt 设计策略

### 决策: 结构化 System Prompt + 上下文注入

**理由**:
- 需要指导 LLM 生成正确格式的命令（不带解释）
- 需要提供执行环境信息（操作系统、Shell类型等）
- 需要约束输出格式（仅返回命令，不返回 Markdown 代码块）

**Prompt 模板**:
```
System: You are a command-line expert assistant. Your task is to translate natural language descriptions into executable shell commands.

Context:
- Operating System: {OS}
- Shell: {Shell}
- Working Directory: {WorkingDir}
- Input from stdin: {StdinPreview}

Rules:
1. Return ONLY the command, no explanations or markdown formatting
2. Return a single command (use && or | for chaining if needed)
3. Use the correct syntax for the specified shell
4. If the request is ambiguous, return the most likely command
5. If the request is impossible or dangerous, return "ERROR: {reason}"

User: {NaturalLanguageInput}