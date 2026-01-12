# 任务: 自然语言命令行转换工具

**输入**: 来自 `/specs/001-nl-to-cmd/` 的设计文档
**前置条件**: plan.md(必需)、spec.md(用户故事必需)、research.md

**测试**: 根据项目章程,测试是**强制要求**且不可协商。必须:
- 遵循 TDD 流程(测试先行)
- 单元测试覆盖率达到 60%+(核心业务逻辑 80%+)
- 使用 `go test -cover` 验证覆盖率
- 所有测试任务必须在实施任务之前完成

**组织结构**: 任务按用户故事分组, 以便每个故事能够独立实施和测试.

## 格式: `[ID] [P?] [Story] 描述`
- **[P]**: 可以并行运行(不同文件, 无依赖关系)
- **[Story]**: 此任务属于哪个用户故事(例如: US1、US2、US3)
- 在描述中包含确切的文件路径

## 路径约定
- **Go 标准项目布局**: 
  - 源代码: `cmd/`, `pkg/`, `internal/` 目录组织
  - 测试: 与源代码同目录,使用 `_test.go` 后缀
  - 文档: `docs/` 目录，使用中文编写

---

## 阶段 1: 设置(共享基础设施)

**目的**: 项目初始化和基本结构

- [ ] T001 根据 plan.md 创建 Go 项目目录结构 (cmd/, pkg/, internal/, tests/, docs/)
- [ ] T002 初始化 Go modules (go.mod), 设置模块名为 github.com/studyzy/aicli
- [ ] T003 [P] 创建 Makefile 包含构建、测试、覆盖率检查命令
- [ ] T004 [P] 配置 .gitignore 文件(二进制文件、配置文件、IDE文件)
- [ ] T005 [P] 设置 GitHub Actions CI 工作流 (.github/workflows/test.yml) 验证测试覆盖率 ≥60%
- [ ] T006 [P] 创建中文 README.md 包含项目介绍、安装和基本使用示例

---

## 阶段 2: 基础(阻塞前置条件)

**目的**: 在任何用户故事可以实施之前必须完成的核心基础设施

**⚠️ 关键**: 在此阶段完成之前, 无法开始任何用户故事工作

- [ ] T007 在 pkg/config/config.go 中定义配置结构体(Config, LLMConfig, ExecutionConfig, SafetyConfig等)
- [ ] T008 在 pkg/config/config.go 中实现 Load() 函数从 ~/.aicli.json 加载配置
- [ ] T009 在 pkg/config/default.go 中定义默认配置值
- [ ] T010 [P] 在 pkg/config/config_test.go 中编写配置加载和默认值测试
- [ ] T011 在 pkg/llm/provider.go 中定义 LLMProvider 接口和 ExecutionContext 结构体
- [ ] T012 在 pkg/llm/mock.go 中实现 MockLLMProvider 用于测试
- [ ] T013 在 pkg/executor/shell.go 中定义 ShellAdapter 结构体和 ShellType 常量
- [ ] T014 在 pkg/executor/shell.go 中实现 DetectShell() 函数检测当前系统Shell
- [ ] T015 [P] 在 pkg/executor/executor_test.go 中编写 Shell 检测测试(跨平台)
- [ ] T016 在 pkg/safety/patterns.go 中定义危险命令模式切片(DangerousPatterns)
- [ ] T017 在 pkg/safety/checker.go 中实现 SafetyChecker 结构体和 IsDangerous() 方法
- [ ] T018 [P] 在 pkg/safety/safety_test.go 中编写危险命令检测测试(确保100%准确率)
- [ ] T019 在 cmd/aicli/main.go 中创建基础 Cobra 命令结构和根命令
- [ ] T020 在 internal/app/flags.go 中定义命令行标志(--dry-run, --verbose, --force等)

**检查点**: 基础就绪 - 现在可以开始并行实施用户故事

---

## 阶段 3: 用户故事 1 - 基础自然语言命令转换与执行 (优先级: P1) 🎯 MVP

**目标**: 实现核心功能 - 用户输入自然语言，通过 LLM 转换为命令并执行

**独立测试**: 可以通过 `aicli "列出当前目录文件"` 等简单命令完全测试，验证 LLM 调用、命令转换和执行流程

### 用户故事 1 的测试 ⚠️

**注意: 先编写这些测试, 确保在实施前它们失败**

- [ ] T021 [P] [US1] 在 pkg/llm/llm_test.go 中编写 OpenAI Provider 的单元测试(使用 Mock HTTP 响应)
- [ ] T022 [P] [US1] 在 pkg/executor/executor_test.go 中编写命令执行测试(使用 MockLLMProvider)
- [ ] T023 [P] [US1] 在 tests/integration/basic_test.go 中编写端到端集成测试(Mock LLM, 真实Shell)

### 用户故事 1 的实施

- [ ] T024 [P] [US1] 在 pkg/llm/openai.go 中实现 OpenAIProvider 结构体和 Translate() 方法
- [ ] T025 [P] [US1] 在 pkg/llm/prompt.go 中实现 Prompt 模板构建函数(包含系统提示和上下文)
- [ ] T026 [US1] 在 pkg/executor/executor.go 中实现 Executor 结构体和 Execute() 方法
- [ ] T027 [US1] 在 pkg/executor/executor.go 中实现命令执行逻辑(调用 os/exec, 处理stdout/stderr)
- [ ] T028 [US1] 在 internal/app/app.go 中实现应用主逻辑(整合 Config, LLMProvider, Executor)
- [ ] T029 [US1] 在 internal/app/app.go 中实现自然语言输入解析和LLM调用流程
- [ ] T030 [US1] 在 internal/app/app.go 中实现危险命令检测和用户确认流程
- [ ] T031 [US1] 在 cmd/aicli/main.go 中连接 Cobra 命令到 app.Run() 函数
- [ ] T032 [US1] 添加错误处理(LLM 超时、网络错误、命令执行失败)到 internal/app/app.go
- [ ] T033 [US1] 在所有导出函数和类型添加中文 godoc 注释

**检查点**: 此时, 用户故事 1 应该完全功能化且可独立测试
运行 `aicli "显示当前目录"` 应该成功转换并执行命令

---

## 阶段 4: 用户故事 2 - 标准输入输出与管道符支持 (优先级: P2)

**目标**: 实现管道符支持，让 aicli 可以从 stdin 读取数据并输出到 stdout

**独立测试**: 可以通过 `echo "test" | aicli "统计字符数"` 和 `aicli "列出txt文件" | wc -l` 验证管道流转

### 用户故事 2 的测试 ⚠️

- [ ] T034 [P] [US2] 在 tests/integration/pipeline_test.go 中编写管道输入测试(stdin → aicli)
- [ ] T035 [P] [US2] 在 tests/integration/pipeline_test.go 中编写管道输出测试(aicli → stdout)
- [ ] T036 [P] [US2] 在 tests/integration/pipeline_test.go 中编写链式管道测试(stdin → aicli → stdout)

### 用户故事 2 的实施

- [ ] T037 [P] [US2] 在 internal/app/io.go 中实现 stdin 检测函数(检查是否有管道输入)
- [ ] T038 [P] [US2] 在 internal/app/io.go 中实现 stdin 读取函数(读取所有输入数据)
- [ ] T039 [US2] 在 internal/app/app.go 中集成 stdin 数据到 ExecutionContext
- [ ] T040 [US2] 在 pkg/llm/prompt.go 中更新 Prompt 模板,将 stdin 数据作为上下文传递给 LLM
- [ ] T041 [US2] 在 pkg/executor/executor.go 中实现命令输出到 stdout 而非直接打印
- [ ] T042 [US2] 在 internal/app/app.go 中实现管道模式检测(stdin存在时自动非交互模式)
- [ ] T043 [US2] 添加 --no-send-stdin 标志到 internal/app/flags.go(隐私保护)
- [ ] T044 [US2] 在 internal/app/io_test.go 中编写 stdin 检测和读取的单元测试

**检查点**: 此时, 用户故事 1 和 2 都应该独立运行
运行 `cat file.txt | aicli "查找ERROR"` 应该正确处理管道输入

---

## 阶段 5: 用户故事 3 - 交互式命令确认与历史记录 (优先级: P3)

**目标**: 实现历史记录功能和危险命令的交互式确认机制

**独立测试**: 可以通过执行危险命令验证确认提示，通过 `aicli --history` 查看历史

### 用户故事 3 的测试 ⚠️

- [ ] T045 [P] [US3] 在 internal/history/history_test.go 中编写历史记录保存和加载测试
- [ ] T046 [P] [US3] 在 internal/app/app_test.go 中编写危险命令确认流程测试

### 用户故事 3 的实施

- [ ] T047 [P] [US3] 在 internal/history/history.go 中定义 HistoryEntry 结构体
- [ ] T048 [P] [US3] 在 internal/history/history.go 中实现 History 结构体和 Add() 方法
- [ ] T049 [US3] 在 internal/history/history.go 中实现 Load() 和 Save() 方法(从/到 ~/.aicli_history.json)
- [ ] T050 [US3] 在 internal/history/history.go 中实现 List() 和 Get() 方法
- [ ] T051 [US3] 在 internal/app/app.go 中集成历史记录保存(每次命令执行后)
- [ ] T052 [US3] 在 internal/app/confirm.go 中实现 confirmDangerousCommand() 函数(显示命令并等待 y/n)
- [ ] T053 [US3] 在 internal/app/app.go 中集成危险命令确认流程(SafetyChecker → confirm)
- [ ] T054 [US3] 在 cmd/aicli/main.go 中添加 --history 子命令
- [ ] T055 [US3] 在 cmd/aicli/main.go 中添加 --retry <ID> 子命令
- [ ] T056 [US3] 添加 --force 标志跳过确认到 internal/app/flags.go

**检查点**: 所有用户故事现在应该独立功能化
运行 `aicli "删除临时文件"` 应该触发确认提示，`aicli --history` 应该显示历史

---

## 阶段 6: 扩展 LLM 提供商支持

**目的**: 添加更多 LLM 提供商，提高工具灵活性

- [ ] T057 [P] 在 pkg/llm/anthropic.go 中实现 AnthropicProvider 结构体和 Translate() 方法
- [ ] T058 [P] 在 pkg/llm/localmodel.go 中实现 LocalModelProvider 支持本地 LLM (如 Ollama)
- [ ] T059 在 pkg/llm/factory.go 中实现 Provider 工厂函数(根据配置创建对应Provider)
- [ ] T060 在 internal/app/app.go 中集成 Provider 工厂，支持动态选择 LLM
- [ ] T061 [P] 在 pkg/llm/anthropic_test.go 中编写 Anthropic Provider 测试
- [ ] T062 [P] 在 pkg/llm/localmodel_test.go 中编写本地模型 Provider 测试

---

## 阶段 7: 完善与横切关注点

**目的**: 影响多个用户故事的改进和文档完善

- [ ] T063 [P] 在 docs/architecture.md 中编写架构设计文档(中文，包含模块图和数据流)
- [ ] T064 [P] 在 docs/configuration.md 中编写配置文件详细说明(中文，包含所有配置项)
- [ ] T065 [P] 在 docs/development.md 中编写开发指南(中文，包含构建、测试、贡献流程)
- [ ] T066 [P] 在 README.md 中完善使用示例和常见问题解答(中文)
- [ ] T067 [P] 创建 CONTRIBUTING.md 贡献指南(中文)
- [ ] T068 在 internal/app/app.go 中实现 --verbose 模式(显示详细的LLM请求和响应)
- [ ] T069 在 internal/app/app.go 中实现 --dry-run 模式(仅显示命令不执行)
- [ ] T070 [P] 在 pkg/llm/logger.go 中实现可选日志记录(不包含敏感数据)
- [ ] T071 代码审查和重构(确保所有函数有中文注释，代码符合 gofmt)
- [ ] T072 运行 `go test -cover ./...` 验证整体覆盖率 ≥60%，核心包 ≥80%
- [ ] T073 使用 golangci-lint 进行静态代码检查并修复问题
- [ ] T074 创建示例配置文件 example.aicli.json 并添加说明注释
- [ ] T075 编写安装脚本或提供预编译二进制文件下载说明

---

## 依赖关系与执行顺序

### 阶段依赖关系

- **设置(阶段 1)**: 无依赖关系 - 可立即开始
- **基础(阶段 2)**: 依赖于设置完成 - 阻塞所有用户故事
- **用户故事(阶段 3-5)**: 都依赖于基础阶段完成
  - 用户故事 1(P1): 可在基础后立即开始(MVP 核心)
  - 用户故事 2(P2): 可在用户故事 1 后开始，或与 US1 并行(有依赖但可独立测试)
  - 用户故事 3(P3): 可在基础后开始，与 US1/US2 独立
- **扩展(阶段 6)**: 依赖于用户故事 1 完成(需要 LLM 接口稳定)
- **完善(阶段 7)**: 依赖于所有用户故事完成

### 用户故事依赖关系

- **用户故事 1(P1)**: 可在基础(阶段 2)后开始 - 无其他故事依赖
  - T024-T033 依赖于 T007-T020 (基础)
- **用户故事 2(P2)**: 依赖于用户故事 1 的核心执行流程
  - T037-T044 依赖于 T026-T028 (Executor 和 App 基础)
  - 但可独立测试管道功能
- **用户故事 3(P3)**: 依赖于用户故事 1 的基础流程，但历史记录独立
  - T047-T056 可以并行开发，仅需要 App 基础结构

### 每个用户故事内部

- 测试必须在实施前编写并失败(TDD)
- Provider 接口在 Provider 实现之前
- Executor 在 App 集成之前
- 核心逻辑在错误处理和日志之前
- 故事完成后才移至下一个优先级

### 并行机会

- **阶段 1**: T003, T004, T005, T006 可以并行(不同文件)
- **阶段 2**: T010(配置测试), T015(executor测试), T018(safety测试) 可以并行
- **阶段 2**: T007-T009(config), T011-T012(llm接口), T013-T015(executor), T016-T018(safety) 四组可以并行
- **用户故事 1**: T021, T022, T023(测试) 可以并行
- **用户故事 1**: T024-T025(LLM相关) 可以与 T026-T027(Executor相关) 并行
- **用户故事 2**: T034, T035, T036(测试) 可以并行
- **用户故事 2**: T037-T038(IO相关) 可以与 T041(executor输出) 并行
- **用户故事 3**: T045, T046(测试) 可以并行
- **用户故事 3**: T047-T050(history模块) 可以与 T052(confirm函数) 并行
- **阶段 6**: T057, T058(不同Provider) 可以并行, T061, T062(测试) 可以并行
- **阶段 7**: T063, T064, T065, T066, T067(文档) 全部可以并行

---

## 并行示例

### 阶段 2 基础 - 四组并行

```bash
# 组 1: 配置模块
任务 T007: "定义配置结构体"
任务 T008: "实现配置加载"
任务 T009: "定义默认配置"
任务 T010: "编写配置测试"

# 组 2: LLM 接口模块
任务 T011: "定义 LLMProvider 接口"
任务 T012: "实现 MockLLMProvider"

# 组 3: Executor 模块
任务 T013: "定义 ShellAdapter"
任务 T014: "实现 DetectShell"
任务 T015: "编写 Shell 测试"

# 组 4: Safety 模块
任务 T016: "定义危险模式"
任务 T017: "实现 SafetyChecker"
任务 T018: "编写安全测试"
```

### 用户故事 1 - 测试并行

```bash
# 一起启动所有测试:
任务 T021: "LLM 单元测试"
任务 T022: "Executor 单元测试"
任务 T023: "集成测试"
```

### 用户故事 1 - LLM 和 Executor 并行

```bash
# LLM 组:
任务 T024: "实现 OpenAIProvider"
任务 T025: "实现 Prompt 模板"

# Executor 组 (同时进行):
任务 T026: "实现 Executor 结构体"
任务 T027: "实现命令执行逻辑"
```

---

## 实施策略

### 仅 MVP(仅用户故事 1)

1. 完成阶段 1: 设置 (T001-T006)
2. 完成阶段 2: 基础 (T007-T020) - 关键阻塞
3. 完成阶段 3: 用户故事 1 (T021-T033)
4. **停止并验证**: 独立测试用户故事 1
   - 运行 `aicli "显示当前目录"` 验证基本功能
   - 运行 `go test -cover ./...` 验证覆盖率 ≥60%
5. 如准备好则部署/演示 MVP

### 增量交付

1. 完成设置 + 基础 → 基础就绪
2. 添加用户故事 1 → 独立测试 → 部署/演示 (MVP! )
   - 此时用户可以使用自然语言执行基本命令
3. 添加用户故事 2 → 独立测试 → 部署/演示
   - 此时用户可以在管道中使用 aicli
4. 添加用户故事 3 → 独立测试 → 部署/演示
   - 此时用户可以查看历史和安全确认
5. 每个故事在不破坏先前故事的情况下增加价值

### 并行团队策略

有多个开发人员时: 

1. 团队一起完成设置 + 基础 (阶段 1-2)
2. 基础完成后: 
   - 开发人员 A: 用户故事 1 (核心功能，最高优先级)
   - 开发人员 B: 用户故事 3 (历史记录，相对独立)
   - 开发人员 C: 阶段 6 扩展 LLM 提供商 (与 US1 并行)
3. US1 完成后:
   - 开发人员 A: 用户故事 2 (管道支持，依赖 US1)
   - 开发人员 B: 完成 US3 并开始阶段 7 文档
   - 开发人员 C: 完成阶段 6 并协助测试
4. 故事独立完成和集成

---

## 测试覆盖率目标

根据项目章程和 research.md 中的测试策略:

| 包 | 目标覆盖率 | 原因 |
|----|-----------|------|
| pkg/llm | 80%+ | 核心业务逻辑，LLM 集成关键 |
| pkg/executor | 80%+ | 核心业务逻辑，命令执行关键 |
| pkg/safety | 85%+ | 安全关键，100%准确率要求 |
| pkg/config | 70%+ | 配置管理，重要但非核心 |
| internal/app | 60%+ | 应用层，集成逻辑 |
| internal/history | 65%+ | 历史记录功能 |
| **整体** | **65%+** | **超过章程要求的 60%** |

**验证命令**: 
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 注意事项

- **[P] 任务** = 不同文件, 无依赖关系，可以并行执行
- **[Story] 标签** 将任务映射到特定用户故事以实现可追溯性
- **每个用户故事应该独立可完成和可测试** - 这是增量交付的关键
- **在实施前验证测试失败** (TDD 红-绿-重构循环)
- **在每个任务或逻辑组后提交** - 保持提交粒度合理
- **在任何检查点停止以独立验证故事** - 确保质量
- **所有代码必须有中文注释** - 遵循项目章程
- **避免**: 
  - 模糊任务(每个任务必须有明确的输出)
  - 相同文件冲突(标记为 [P] 的任务不能修改同一文件)
  - 破坏独立性的跨故事依赖(US2 和 US3 应该不依赖彼此)

---

## 快速开始验证

完成相应阶段后的验证命令:

**阶段 1-2 完成后**:
```bash
go build ./cmd/aicli
./aicli --version  # 应该显示版本信息
go test ./pkg/config -v  # 配置测试应该通过
```

**用户故事 1 完成后 (MVP)**:
```bash
./aicli "显示当前目录"
./aicli "列出所有go文件"
./aicli "查找main.go文件"
go test -cover ./...  # 覆盖率应该 ≥60%
```

**用户故事 2 完成后**:
```bash
echo "测试内容" | ./aicli "统计字符数"
./aicli "列出txt文件" | wc -l
cat README.md | ./aicli "提取标题"
```

**用户故事 3 完成后**:
```bash
./aicli "删除临时文件"  # 应该触发确认
./aicli --history  # 显示历史记录
./aicli --retry 1  # 重新执行历史命令
```
