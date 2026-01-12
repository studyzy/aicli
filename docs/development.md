# aicli 开发指南

> 版本: 1.0 | 日期: 2026-01-13

## 目录

- [开发环境设置](#开发环境设置)
- [项目结构](#项目结构)
- [构建与测试](#构建与测试)
- [代码规范](#代码规范)
- [贡献流程](#贡献流程)
- [调试技巧](#调试技巧)
- [常见问题](#常见问题)

## 开发环境设置

### 前置要求

- **Go**: 1.21 或更高版本
- **Git**: 用于版本控制
- **Make**: 用于构建自动化（可选）

### 安装 Go

**macOS**:
```bash
brew install go
```

**Linux**:
```bash
# 下载并安装
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 添加到 PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**Windows**:
从 [https://go.dev/dl/](https://go.dev/dl/) 下载安装程序

### 克隆项目

```bash
git clone https://github.com/studyzy/aicli.git
cd aicli
```

### 安装依赖

```bash
go mod download
```

### 验证环境

```bash
go version
go env
```

## 项目结构

```
aicli/
├── cmd/                    # 命令行入口
│   └── aicli/
│       └── main.go         # 主程序入口
├── internal/               # 私有代码
│   ├── app/                # 应用逻辑
│   └── history/            # 历史记录
├── pkg/                    # 公共库
│   ├── config/             # 配置管理
│   ├── executor/           # 命令执行
│   ├── llm/                # LLM 提供商
│   └── safety/             # 安全检查
├── tests/                  # 集成测试
│   └── integration/
├── docs/                   # 文档
├── bin/                    # 编译输出（忽略）
├── go.mod                  # Go 模块定义
├── go.sum                  # 依赖校验
├── Makefile                # 构建脚本
└── README.md               # 项目说明

```

### 目录说明

- **cmd/**: 可执行程序入口，保持简洁
- **internal/**: 仅供项目内部使用的代码
- **pkg/**: 可被外部项目导入的库
- **tests/**: 集成测试和 E2E 测试
- **docs/**: 项目文档

## 构建与测试

### 使用 Make

项目提供了 Makefile 简化常用操作：

```bash
# 查看所有可用命令
make help

# 构建项目
make build

# 运行所有测试
make test

# 查看测试覆盖率
make coverage

# 运行代码检查
make lint

# 清理构建产物
make clean
```

### 手动构建

```bash
# 构建
go build -o bin/aicli ./cmd/aicli

# 交叉编译（Linux）
GOOS=linux GOARCH=amd64 go build -o bin/aicli-linux ./cmd/aicli

# 交叉编译（Windows）
GOOS=windows GOARCH=amd64 go build -o bin/aicli.exe ./cmd/aicli

# 交叉编译（macOS ARM64）
GOOS=darwin GOARCH=arm64 go build -o bin/aicli-darwin-arm64 ./cmd/aicli
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行指定包的测试
go test ./pkg/llm

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 运行单个测试
go test -v ./pkg/llm -run TestOpenAIProvider_Translate_Success

# 运行测试（详细输出）
go test -v ./...
```

### 代码覆盖率目标

| 包 | 目标覆盖率 | 原因 |
|----|-----------|------|
| pkg/llm | 80%+ | 核心业务逻辑 |
| pkg/executor | 80%+ | 核心业务逻辑 |
| pkg/safety | 85%+ | 安全关键 |
| pkg/config | 70%+ | 配置管理 |
| internal/app | 60%+ | 应用层 |
| internal/history | 65%+ | 历史记录 |
| **整体** | **65%+** | 项目标准 |

## 代码规范

### Go 代码规范

遵循官方 [Effective Go](https://go.dev/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)。

### 格式化

```bash
# 格式化所有代码
go fmt ./...

# 或使用 gofmt
gofmt -w .

# 使用 goimports（自动管理导入）
go install golang.org/x/tools/cmd/goimports@latest
goimports -w .
```

### 代码检查

```bash
# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 运行检查
golangci-lint run

# 运行检查（显示所有问题）
golangci-lint run --max-issues-per-linter 0 --max-same-issues 0
```

### 命名规范

- **包名**: 小写，单个单词，如 `executor`, `config`
- **文件名**: 小写加下划线，如 `shell_adapter.go`
- **变量名**: 驼峰命名，如 `shellType`, `apiKey`
- **常量**: 驼峰命名，如 `DefaultTimeout`
- **接口**: 驼峰命名，通常以 `-er` 结尾，如 `LLMProvider`, `Executor`

### 注释规范

**所有导出的函数、类型、常量必须有中文注释**:

```go
// LLMProvider 定义 LLM 服务提供商的接口
type LLMProvider interface {
    // Translate 将自然语言转换为命令
    // ctx: 上下文，用于超时控制
    // input: 用户的自然语言描述
    // execCtx: 执行上下文信息
    // 返回: 转换后的命令字符串和可能的错误
    Translate(ctx context.Context, input string, execCtx *ExecutionContext) (string, error)
}
```

### 错误处理

```go
// ✅ 好的做法：包装错误
if err != nil {
    return fmt.Errorf("加载配置失败: %w", err)
}

// ❌ 不好的做法：丢失错误上下文
if err != nil {
    return err
}
```

### 测试规范

```go
// 测试函数命名：Test + 被测试的函数/方法 + 测试场景
func TestOpenAIProvider_Translate_Success(t *testing.T) {
    // Arrange（准备）
    provider := NewOpenAIProvider("key", "model", "url")
    
    // Act（执行）
    result, err := provider.Translate(ctx, input, execCtx)
    
    // Assert（断言）
    if err != nil {
        t.Fatalf("期望成功，但返回错误: %v", err)
    }
    if result != expected {
        t.Errorf("期望 %s, 实际 %s", expected, result)
    }
}
```

## 贡献流程

### 1. Fork 项目

在 GitHub 上 Fork 项目到你的账号。

### 2. 创建分支

```bash
git checkout -b feature/my-new-feature
# 或
git checkout -b fix/issue-123
```

**分支命名规范**:
- `feature/` - 新功能
- `fix/` - Bug 修复
- `docs/` - 文档更新
- `refactor/` - 代码重构
- `test/` - 测试相关

### 3. 编写代码

- 遵循代码规范
- 编写测试（TDD）
- 添加中文注释
- 更新相关文档

### 4. 提交更改

```bash
# 添加文件
git add .

# 提交（使用有意义的提交消息）
git commit -m "feat: 添加 Anthropic Provider 支持"
```

**提交消息格式**:
```
<type>: <subject>

<body>

<footer>
```

**Type**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建、工具等

### 5. 推送到 GitHub

```bash
git push origin feature/my-new-feature
```

### 6. 创建 Pull Request

在 GitHub 上创建 PR，描述你的更改。

### 7. 代码审查

响应审查意见，根据需要修改代码。

### 8. 合并

PR 被批准后，维护者会合并你的代码。

## 调试技巧

### 使用 Verbose 模式

```bash
# 显示详细执行信息
./bin/aicli --verbose "列出文件"
```

输出示例：
```
自然语言输入: 列出文件
执行上下文: OS=darwin Shell=zsh WorkDir=/Users/user/project
转换后的命令: ls -la
转换耗时: 1.234s
开始执行命令...
执行耗时: 0.012s
总耗时: 1.246s
```

### 使用 Dry-run 模式

```bash
# 只显示命令，不执行
./bin/aicli --dry-run "删除临时文件"
```

### 使用 Delve 调试器

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug ./cmd/aicli -- "测试命令"

# 在代码中设置断点
(dlv) break main.main
(dlv) continue
```

### 添加日志

```go
import "log"

log.Printf("调试信息: %+v\n", variable)
```

### 使用测试覆盖率定位未测试代码

```bash
go test -coverprofile=coverage.out ./pkg/llm
go tool cover -func=coverage.out
```

## 常见问题

### Q1: 如何添加新的 LLM Provider？

**步骤**:

1. 在 `pkg/llm/` 创建新文件 `myprovider.go`
2. 实现 `LLMProvider` 接口
3. 在 `factory.go` 中注册
4. 创建测试文件 `myprovider_test.go`
5. 更新文档

**示例**:
```go
// pkg/llm/myprovider.go
type MyProvider struct {
    apiKey string
}

func NewMyProvider(apiKey string) *MyProvider {
    return &MyProvider{apiKey: apiKey}
}

func (p *MyProvider) Name() string {
    return "myprovider"
}

func (p *MyProvider) Translate(ctx context.Context, input string, execCtx *ExecutionContext) (string, error) {
    // 实现逻辑
}
```

```go
// pkg/llm/factory.go
case "myprovider":
    return NewMyProvider(cfg.LLM.APIKey), nil
```

### Q2: 如何添加新的安全检查模式？

编辑 `pkg/safety/patterns.go`:

```go
var DangerousPatterns = []Pattern{
    // ... 现有模式
    {
        Regex:       regexp.MustCompile(`my-dangerous-pattern`),
        Description: "危险操作描述",
        Level:       RiskHigh,
    },
}
```

### Q3: 测试失败怎么办？

```bash
# 查看具体错误
go test -v ./pkg/llm

# 只运行失败的测试
go test -v ./pkg/llm -run TestFailingTest

# 查看测试覆盖
go test -cover ./pkg/llm
```

### Q4: 如何模拟 LLM 响应进行测试？

使用 `MockLLMProvider`:

```go
mockLLM := &llm.MockLLMProvider{
    TranslateFn: func(input string) string {
        if input == "列出文件" {
            return "ls -la"
        }
        return "echo unknown"
    },
}
```

或使用 `httptest.Server`:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    response := map[string]interface{}{
        "choices": []map[string]interface{}{
            {"message": map[string]string{"content": "ls -la"}},
        },
    }
    json.NewEncoder(w).Encode(response)
}))
defer server.Close()

provider := llm.NewOpenAIProvider("key", "model", server.URL)
```

### Q5: 如何处理跨平台问题？

使用 Go 标准库的跨平台功能：

```go
import (
    "runtime"
    "path/filepath"
)

// 获取操作系统
os := runtime.GOOS // "linux", "darwin", "windows"

// 路径处理
path := filepath.Join("dir", "file.txt") // 自动使用正确的分隔符
```

## 性能优化

### 性能测试

```bash
# 运行性能测试
go test -bench=. ./pkg/llm

# 生成 CPU 性能分析
go test -cpuprofile=cpu.prof -bench=. ./pkg/llm
go tool pprof cpu.prof

# 生成内存性能分析
go test -memprofile=mem.prof -bench=. ./pkg/llm
go tool pprof mem.prof
```

### 优化建议

1. **避免不必要的内存分配**
2. **使用 `strings.Builder` 而非字符串拼接**
3. **复用 HTTP Client**
4. **使用 Context 控制超时**

## 发布流程

### 1. 更新版本号

编辑 `cmd/aicli/main.go`:

```go
const version = "0.2.0"
```

### 2. 更新 CHANGELOG

记录新功能、修复、变更。

### 3. 创建标签

```bash
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

### 4. 构建发布版本

```bash
make release
```

### 5. 创建 GitHub Release

上传编译好的二进制文件。

## 资源链接

- [Go 官方文档](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Cobra CLI 文档](https://github.com/spf13/cobra)
- [项目 GitHub](https://github.com/studyzy/aicli)

## 获取帮助

- **Issues**: [GitHub Issues](https://github.com/studyzy/aicli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/studyzy/aicli/discussions)
- **文档**: 查看 `docs/` 目录

## 贡献者

感谢所有为 aicli 做出贡献的开发者！

查看完整贡献者列表: [Contributors](https://github.com/studyzy/aicli/graphs/contributors)
