# 开发指南

本项目使用 Go 编写，遵循 Go 标准项目布局。

## 环境要求

- Go 1.25+

建议先确认开发环境：

```bash
go version
go env
```

## 项目结构

```text
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
├── README.md               # English documentation
└── README_CN.md            # 中文项目说明
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
