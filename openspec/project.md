# 项目上下文

## 目的

aicli 是一个基于 Go 语言开发的自然语言命令行工具,将自然语言描述转换为 shell 命令并执行。

**核心目标**:
- **简化 CLI 操作**: 用户通过自然语言描述意图,无需记忆复杂命令语法
- **安全可靠**: 检测危险命令并要求确认,保护用户系统安全
- **多 LLM 支持**: 支持 OpenAI、Anthropic、Ollama 等多种 LLM 提供商
- **管道友好**: 与标准输入输出无缝集成,便于组合使用
- **跨平台**: 支持 Linux、macOS、Windows 三大平台

## 技术栈

### 核心技术
- **Go 1.25**: 主要编程语言,单二进制部署
- **github.com/spf13/cobra**: CLI 框架,处理命令行参数和子命令
- Go 标准库:
  - `net/http`: HTTP 客户端,调用 LLM API
  - `encoding/json`: JSON 序列化/反序列化
  - `os/exec`: Shell 命令执行
  - `context`: 超时和取消管理

### 开发工具
- **Make**: 构建自动化
- **golangci-lint**: 代码静态分析和 lint
- **Go test**: 单元测试和覆盖率测试
- **GitHub Actions**: CI/CD 流程

### 外部服务
- **OpenAI API**: GPT 系列模型
- **Anthropic API**: Claude 系列模型
- **Ollama**: 本地模型支持

## 项目约定

### 代码风格

#### 命名约定
- **包名**: 小写单词,简短且有意义 (如 `llm`, `executor`, `safety`)
- **文件名**: 小写下划线分隔 (如 `anthropic_provider.go`)
- **类型名**: 大驼峰 (如 `AnthropicProvider`, `ExecutionContext`)
- **函数名**: 大驼峰(导出)/小驼峰(私有) (如 `Translate()`, `cleanCommand()`)
- **常量**: 小驼峰带前缀 (如 `providerAnthropic`, `defaultTimeout`)

#### 格式化规则
- 使用 `go fmt` 格式化所有代码
- 行长度限制: 120 字符 (golangci-lint 强制)
- 函数长度限制: 80 行/80 语句 (测试代码除外)
- 圈复杂度限制: 20

#### 注释规范
- 包注释: 每个包的主文件顶部必须有包注释
  ```go
  // Package llm 提供了 LLM 提供商的抽象和实现
  package llm
  ```
- 导出类型/函数: 必须有注释,格式为 "TypeName/FuncName + 描述"
  ```go
  // AnthropicProvider 实现了 Anthropic Claude API 的 LLMProvider 接口
  type AnthropicProvider struct { ... }
  ```
- 复杂逻辑: 添加行内注释说明算法和设计决策

### 架构模式

#### 分层架构
```
用户界面层 (cmd/aicli)
    ↓
应用逻辑层 (internal/app)
    ↓
核心服务层 (pkg/llm, pkg/executor, pkg/safety, pkg/config)
    ↓
外部服务层 (LLM APIs, OS Shell)
```

#### 设计模式
- **工厂模式**: `llm.NewProvider()` 根据配置创建不同的 LLM 提供商
- **接口抽象**: `LLMProvider` 接口统一不同 LLM 的调用方式
- **策略模式**: 安全检查器可插拔不同的检查策略
- **依赖注入**: 通过结构体字段注入依赖,便于测试

#### 目录结构约定
- `cmd/`: 可执行程序入口
- `pkg/`: 可被外部导入的公共包
- `internal/`: 项目内部包,不对外暴露
- `tests/`: 集成测试和 E2E 测试
- `docs/`: 项目文档

### 测试策略

#### 测试层次
1. **单元测试**: 每个包独立测试,使用 `*_test.go` 文件
2. **集成测试**: 跨包测试,放在 `tests/integration/`
3. **E2E 测试**: 完整流程测试,使用 Mock LLM

#### Mock 策略
- `MockLLMProvider`: Mock LLM 响应,避免调用真实 API
- `httptest.Server`: Mock HTTP 服务器
- `os.Pipe()`: Mock stdin/stdout

#### 覆盖率要求
- **整体覆盖率**: ≥60% (Makefile 中强制检查)
- **核心包**: llm、executor、safety 包要求 ≥80%
- 使用 `make coverage` 生成覆盖率报告
- CI 流程中自动检查覆盖率

#### 测试命名
- 测试函数: `TestFunctionName` 或 `TestTypeName_MethodName`
- 表驱动测试: 使用 `tests` slice 定义测试用例
  ```go
  func TestTranslate(t *testing.T) {
      tests := []struct {
          name    string
          input   string
          want    string
          wantErr bool
      }{
          // test cases...
      }
  }
  ```

### Git 工作流

#### 分支策略
- `main`: 主分支,始终保持可发布状态
- `feature/*`: 功能分支 (如 `feature/add-ollama-support`)
- `bugfix/*`: Bug 修复分支 (如 `bugfix/fix-timeout-issue`)
- `release/*`: 发布分支 (如 `release/v0.2.0`)

#### 提交规范 (Conventional Commits)
格式: `<type>(<scope>): <subject>`

**类型 (type)**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式化 (不影响功能)
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建工具、依赖更新等

**示例**:
```
feat(llm): 添加 DeepSeek 提供商支持
fix(executor): 修复 Windows 平台命令执行超时问题
docs(readme): 更新安装说明
test(safety): 增加危险命令检测测试用例
```

#### Pull Request 流程
1. Fork 项目并创建特性分支
2. 完成开发并通过本地测试 (`make test`)
3. 确保代码通过 lint (`make lint`)
4. 提交符合规范的 commit
5. 推送到 fork 仓库并创建 PR
6. 等待 CI 检查通过和代码审查
7. 根据反馈修改代码
8. 合并到主分支

## 领域上下文

### 自然语言处理与命令转换
- **Prompt 工程**: 使用精心设计的 System Prompt 引导 LLM 生成正确命令
- **上下文信息**: 向 LLM 提供操作系统、Shell 类型、工作目录、stdin 等上下文
- **命令清理**: 从 LLM 响应中提取纯命令 (去除 markdown 代码块、解释文本等)

### Shell 兼容性
- **Unix/Linux**: bash、zsh、sh (使用 `-c` 参数执行)
- **macOS**: 优先 zsh,fallback 到 bash
- **Windows**: PowerShell、cmd (使用 `/c` 参数执行)

### 安全考虑
- **危险命令模式**: 正则表达式检测删除、格式化、权限修改等危险操作
- **用户确认**: 交互式终端下要求用户确认危险命令
- **强制执行**: `--force` 标志跳过确认 (慎用)
- **隐私保护**: `--no-send-stdin` 避免发送敏感数据到 LLM

### 性能要求
- LLM API 调用 + 命令执行总耗时 < 5 秒
- 配置文件加载 < 10ms
- 命令执行启动 < 100ms

## 重要约束

### 技术约束
- **Go 版本**: 需要 Go 1.25 或更高版本
- **网络依赖**: 需要访问 LLM API (除非使用本地模型)
- **Shell 可用性**: 需要系统安装标准 Shell (bash/zsh/PowerShell)

### 业务约束
- **API 成本**: 每次调用 LLM 产生 API 费用,需用户自备 API Key
- **隐私**: 用户输入和 stdin 会发送到第三方 LLM 服务 (除非使用本地模型)
- **准确性**: LLM 生成的命令可能不准确,需用户验证

### 安全约束
- **配置文件权限**: `~/.aicli.json` 建议设置为 600,避免泄露 API Key
- **危险命令确认**: 默认启用安全检查,禁止直接执行危险命令
- **日志脱敏**: 日志中不记录完整 API Key

### 兼容性约束
- **跨平台差异**: Windows 和 Unix 的命令语法不同,需测试多平台
- **Shell 差异**: bash 和 zsh 语法有细微差别,需兼容处理

## 外部依赖

### LLM API 服务
- **OpenAI API**
  - Endpoint: `https://api.openai.com/v1/chat/completions`
  - 认证: Bearer Token
  - 模型: gpt-4, gpt-3.5-turbo 等

- **Anthropic API**
  - Endpoint: `https://api.anthropic.com/v1/messages`
  - 认证: `x-api-key` Header
  - 模型: claude-3-opus, claude-3-sonnet 等

- **Ollama (本地)**
  - Endpoint: `http://localhost:11434` (可配置)
  - 无需认证
  - 模型: llama2, mistral 等

- **其他 OpenAI 兼容 API**
  - DeepSeek: `https://api.deepseek.com/v1`
  - 自定义 `api_base` 支持

### 系统依赖
- **Shell**: bash, zsh (Unix/macOS), PowerShell/cmd (Windows)
- **文件系统**: 配置文件存储在 `~/.aicli.json`, 历史记录存储在 `~/.aicli_history.json`

### 开发依赖
- **golangci-lint**: 代码静态分析 (可选,建议安装)
- **Make**: 构建自动化 (可选,也可直接使用 `go` 命令)
