# 实施计划: 自然语言命令行转换工具

**分支**: `001-nl-to-cmd` | **日期**: 2026-01-12 | **规范**: [spec.md](spec.md)
**输入**: 来自 `/specs/001-nl-to-cmd/spec.md` 的功能规范

**注意**: 此模板由 `/speckit.plan` 命令填充. 执行工作流程请参见 `.specify/templates/commands/plan.md`.

## 摘要

构建一个 Go 语言实现的命令行工具 `aicli`，允许用户使用自然语言描述意图，通过 LLM 服务将其转换为实际的 shell 命令并执行。核心功能包括：
- 接受自然语言输入并调用 LLM API 进行命令转换
- 在当前 shell 环境中执行转换后的命令
- 支持标准输入输出和管道符，无缝集成到命令行工作流
- 检测危险命令并提供确认机制
- 通过 `~/.aicli.json` 配置文件管理 LLM API 密钥和选项

技术方法：采用模块化包设计，抽象 LLM 提供商接口以支持多种 AI 服务，使用 Go 标准库处理进程执行和 IO 流，遵循 TDD 流程确保 60%+ 测试覆盖率。

## 技术背景

**语言/版本**: Go 1.21+（利用泛型和改进的错误处理）
**主要依赖**: 
- `encoding/json` - 配置文件解析和 LLM API 交互
- `os/exec` - 命令执行
- `net/http` - HTTP 客户端用于 LLM API 调用
- `github.com/spf13/cobra` - CLI 框架（命令行参数解析）
- `github.com/spf13/viper` - 配置管理
- 待研究：具体的 LLM SDK（OpenAI Go SDK、Anthropic SDK 等）

**存储**: 文件系统（配置文件 `~/.aicli.json`，可选的命令历史记录 `~/.aicli_history.json`）
**测试**: `go test` + `testing` 包，使用 `go test -cover` 验证覆盖率，Mock LLM API 响应进行单元测试
**目标平台**: 跨平台（Linux、macOS、Windows），编译为单一可执行二进制文件
**项目类型**: 单一命令行工具
**性能目标**: 
- LLM API 调用 + 命令执行总时间 < 5 秒（SC-001）
- 本地命令解析和执行启动时间 < 100ms
- 配置文件读取时间 < 10ms

**约束条件**: 
- LLM API 超时设置为 10 秒
- 配置文件大小 < 1MB
- 命令历史记录保留最近 1000 条
- 内存占用 < 50MB（不含 LLM 响应缓存）

**规模/范围**: 
- 单用户命令行工具
- 代码库规模预计 5000-10000 行
- 支持 3-5 种主流 LLM 提供商
- 命令转换准确率目标 90%+（SC-003）

## 章程检查

*门控: 必须在阶段 0 研究前通过. 阶段 1 设计后重新检查.*

根据 `.specify/memory/constitution.md` 验证以下要求:

- [x] **模块化库优先**: ✅ 功能按独立 Go 包设计
  - `pkg/llm` - LLM 提供商抽象接口和实现
  - `pkg/executor` - 命令执行引擎
  - `pkg/config` - 配置管理
  - `pkg/safety` - 危险命令检测
  - 每个包职责单一，可独立测试

- [x] **CLI 接口**: ✅ 完全符合 CLI 标准
  - 支持命令行参数（自然语言输入）
  - 支持标准输入（stdin）用于管道符
  - 正常输出到标准输出（stdout）
  - 错误输出到标准错误（stderr）
  - 使用退出码表示执行状态

- [x] **测试计划**: ✅ 已制定 TDD 流程
  - 每个包都有对应的 `_test.go` 文件
  - 使用 Mock 接口测试 LLM 调用
  - 使用测试夹具（fixtures）测试命令执行
  - 目标：整体覆盖率 60%+，核心业务逻辑（llm、executor）80%+
  - CI 集成：GitHub Actions 运行 `go test -cover` 并验证门禁

- [x] **中文文档**: ✅ 已规划完整中文文档
  - README.md：中文项目介绍、安装指南、使用示例
  - 所有导出函数和类型都有 godoc 格式的中文注释
  - docs/architecture.md：架构设计文档（中文）
  - docs/configuration.md：配置文件说明（中文）

- [x] **AI 集成规范**: ✅ 完全符合 AI 集成要求
  - API 抽象：定义 `LLMProvider` 接口，支持 OpenAI、Anthropic、本地模型
  - 配置管理：通过 `~/.aicli.json` 管理 API 密钥（不硬编码）
  - 错误处理：处理 API 超时、限流、网络错误
  - 隐私保护：提供 `--no-send-stdin` 选项，敏感信息不发送到 LLM
  - 可观测性：可选日志记录（不包含敏感数据）

**额外门控条件**:
- [x] **跨平台兼容性**: 命令执行层需要处理不同操作系统的 shell 差异（bash/zsh/PowerShell/CMD）
- [x] **安全性**: 危险命令检测机制必须在执行前拦截（100% 触发确认，SC-006）

## 项目结构

### 文档(此功能)

```
specs/001-nl-to-cmd/
├── plan.md              # 此文件 (/speckit.plan 命令输出)
├── research.md          # 阶段 0 输出 (/speckit.plan 命令)
├── data-model.md        # 阶段 1 输出 (/speckit.plan 命令)
├── quickstart.md        # 阶段 1 输出 (/speckit.plan 命令)
├── contracts/           # 阶段 1 输出 (/speckit.plan 命令)
│   └── llm-api.md       # LLM 提供商接口定义
└── tasks.md             # 阶段 2 输出 (/speckit.tasks 命令 - 非 /speckit.plan 创建)
```

### 源代码(仓库根目录)

选择 **Go 标准项目布局**（适合可扩展的命令行工具）：

```
aicli/
├── cmd/
│   └── aicli/
│       └── main.go                 # 主入口，初始化 Cobra 命令
├── pkg/                            # 公共库（可被外部导入）
│   ├── llm/                        # LLM 提供商抽象层
│   │   ├── provider.go             # LLMProvider 接口定义
│   │   ├── openai.go               # OpenAI 实现
│   │   ├── anthropic.go            # Anthropic 实现
│   │   ├── localmodel.go           # 本地模型实现
│   │   └── llm_test.go             # LLM 包测试
│   ├── executor/                   # 命令执行引擎
│   │   ├── executor.go             # 命令执行接口和实现
│   │   ├── shell.go                # Shell 环境检测和适配
│   │   └── executor_test.go        # 执行器测试
│   ├── config/                     # 配置管理
│   │   ├── config.go               # 配置加载和保存
│   │   ├── config_test.go          # 配置测试
│   │   └── default.go              # 默认配置
│   └── safety/                     # 安全检查
│       ├── checker.go              # 危险命令检测
│       ├── patterns.go             # 危险模式定义
│       └── safety_test.go          # 安全检查测试
├── internal/                       # 私有代码（不可外部导入）
│   ├── app/                        # 应用层
│   │   ├── app.go                  # 应用主逻辑
│   │   ├── flags.go                # 命令行参数定义
│   │   └── app_test.go             # 应用测试
│   └── history/                    # 命令历史记录
│       ├── history.go              # 历史记录管理
│       └── history_test.go         # 历史记录测试
├── tests/                          # 集成测试
│   └── integration/
│       ├── basic_test.go           # 基础场景集成测试
│       └── pipeline_test.go        # 管道场景集成测试
├── docs/                           # 文档
│   ├── architecture.md             # 架构设计（中文）
│   ├── configuration.md            # 配置说明（中文）
│   └── development.md              # 开发指南（中文）
├── go.mod                          # Go modules 依赖管理
├── go.sum                          # 依赖校验和
├── README.md                       # 项目说明（中文）
├── CONTRIBUTING.md                 # 贡献指南（中文）
├── Makefile                        # 构建脚本
└── .github/
    └── workflows/
        └── test.yml                # CI 测试流水线
```

**结构决策**: 
- 选择 Go 标准项目布局，因为：
  1. 项目具有清晰的公共 API（pkg/）和私有实现（internal/）
  2. 需要支持多个 LLM 提供商，模块化包结构便于扩展
  3. 符合 Go 社区最佳实践，便于其他开发者理解
  4. 便于单元测试和集成测试分离
- `cmd/aicli/main.go` 作为唯一入口点，保持简洁，仅负责命令行解析
- `pkg/` 目录下的包可以被其他项目复用（如测试工具）
- `internal/` 目录确保实现细节不被外部依赖

## 复杂度跟踪

*仅在章程检查有必须证明的违规时填写*

| 违规 | 为什么需要 | 拒绝更简单替代方案的原因 |
|-----------|------------|-------------------------------------|
| [例如: 第 4 个项目] | [当前需求] | [为什么 3 个项目不够] |
| [例如: 仓储模式] | [特定问题] | [为什么直接数据库访问不够] |
