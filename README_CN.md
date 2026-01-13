# AICLI - AI 命令行助手

[![Test and Coverage](https://github.com/studyzy/aicli/actions/workflows/test.yml/badge.svg)](https://github.com/studyzy/aicli/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/studyzy/aicli)](https://goreportcard.com/report/github.com/studyzy/aicli)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

[English](README.md) | 中文

一个让命令行支持自然语言操作的 Go 语言工具。用户只需要输入自然语言，就能通过 LLM 服务将其转换为实际的命令并执行。

## ✨ 核心功能

- 🗣️ **自然语言转命令**：使用自然语言描述想要执行的操作，自动转换为 shell 命令
- 🔗 **管道符支持**：完美支持标准输入输出，可与其他命令组合使用
- 🛡️ **安全确认机制**：自动检测危险命令（删除、格式化等），执行前需要用户确认
- 📜 **历史记录**：保存命令历史，支持查看和重新执行
- 🔌 **多 LLM 提供商**：支持 OpenAI、Anthropic、本地模型等多种 LLM 服务
- 🌍 **跨平台**：支持 Linux、macOS 和 Windows 系统

## 🚀 快速开始

### 安装

```bash
# 方式 1: 从源码编译
git clone https://github.com/studyzy/aicli.git
cd aicli
make build
make install

# 方式 2: 使用 go install
go install github.com/studyzy/aicli/cmd/aicli@latest
```

### 配置

#### 自动配置（推荐）

运行 `init` 命令进行交互式配置：

```bash
aicli init
```

这将引导您选择 LLM 提供商并设置 API 密钥。

#### 手动配置

手动创建配置文件 `~/.aicli.json`：

```json
{
  "version": "1.0",
  "llm": {
    "provider": "openai",
    "api_key": "your-api-key-here",
    "model": "gpt-4",
    "timeout": 10
  },
  "execution": {
    "auto_confirm": false,
    "timeout": 30
  },
  "safety": {
    "enable_checks": true,
    "require_confirmation": true
  }
}
```

也可以通过环境变量配置 API 密钥：

```bash
export AICLI_API_KEY="your-api-key-here"
```

### 基础使用

```bash
# 示例 1: 查找文件中的特定内容
aicli "查找log.txt中的ERROR日志"
# 转换为: grep "ERROR" log.txt

# 示例 2: 文件操作
aicli "显示当前目录下所有txt文件"
# 转换为: ls *.txt 或 find . -name "*.txt"

# 示例 3: 管道符支持
cat log.txt | aicli "过滤出包含ERROR的行"
# 将 stdin 内容传递给 LLM，生成对应的过滤命令

# 示例 4: 链式处理
aicli "列出所有txt文件" | wc -l
# aicli 的输出可以被其他命令使用

# 示例 5: 查看历史
aicli --history

# 示例 6: 重新执行历史命令
aicli --retry 3
```

### 命令行选项

```bash
# 仅显示转换后的命令，不执行
aicli --dry-run "删除临时文件"

# 显示详细的转换过程
aicli --verbose "查找所有go文件"

# 强制执行，跳过危险命令确认
aicli --force "删除所有临时文件"

# 不将 stdin 数据发送到 LLM（隐私保护）
cat sensitive.txt | aicli --no-send-stdin "统计行数"
```

## 📖 使用示例

### 文件操作

```bash
# 查找文件
aicli "查找当前目录下所有大于10MB的文件"
# → find . -type f -size +10M

# 文件内容处理
aicli "统计main.go文件的行数"
# → wc -l main.go

aicli "提取config.json中的所有键名"
# → jq 'keys' config.json
```

### 系统管理

```bash
# 系统信息
aicli "显示系统内存使用情况"
# → free -h  (Linux) 或 vm_stat (macOS)

aicli "查看当前网络连接"
# → netstat -an 或 ss -tuln

# 进程管理
aicli "查找占用8080端口的进程"
# → lsof -i :8080

aicli "显示CPU占用最高的5个进程"
# → ps aux --sort=-%cpu | head -6
```

### 文本处理

```bash
# 日志分析
aicli "统计access.log中每个IP的访问次数"
# → awk '{print $1}' access.log | sort | uniq -c | sort -rn

aicli "提取nginx日志中所有404错误的URL"
# → grep ' 404 ' nginx.log | awk '{print $7}' | sort -u

# 数据处理
cat data.csv | aicli "计算第三列的平均值"
# → awk -F',' '{sum+=$3; count++} END {print sum/count}'

aicli "将所有txt文件合并为一个文件" | less
# → cat *.txt
```

### Git 操作

```bash
# Git 日常操作
aicli "显示最近5次提交"
# → git log -5 --oneline

aicli "查看当前分支的所有修改文件"
# → git status --short

aicli "找出提交次数最多的作者"
# → git shortlog -sn
```

### 网络操作

```bash
# 网络请求
aicli "测试api.example.com的连接"
# → curl -I https://api.example.com

aicli "下载并解压example.tar.gz"
# → wget example.tar.gz && tar -xzf example.tar.gz

# 端口检测
aicli "扫描localhost的3000-3010端口"
# → nc -zv localhost 3000-3010
```

### 数据转换

```bash
# JSON 处理
aicli "美化这个JSON文件"
cat data.json | aicli "格式化JSON"
# → jq .

# CSV 操作
aicli "提取users.csv的第2和第4列"
# → awk -F',' '{print $2,$4}' users.csv

# 批量重命名
aicli "将所有jpg文件改为小写扩展名"
# → rename 's/\.JPG$/\.jpg/' *.JPG
```

## 🎯 高级功能

### 多 LLM 提供商支持

aicli 支持多种 LLM 服务，只需修改配置文件即可切换：

#### OpenAI (GPT)

```json
{
  "llm": {
    "provider": "openai",
    "api_key": "sk-xxxxx",
    "model": "gpt-4"
  }
}
```

#### Anthropic (Claude)

```json
{
  "llm": {
    "provider": "anthropic",
    "api_key": "sk-ant-xxxxx",
    "model": "claude-3-sonnet-20240229"
  }
}
```

#### 本地模型 (Ollama)

```json
{
  "llm": {
    "provider": "local",
    "model": "llama2",
    "api_base": "http://localhost:11434"
  }
}
```

#### DeepSeek (兼容 OpenAI API)

```json
{
  "llm": {
    "provider": "openai",
    "api_key": "sk-xxxxx",
    "model": "deepseek-chat",
    "api_base": "https://api.deepseek.com/v1"
  }
}
```

### 历史记录功能

```bash
# 查看所有历史记录
aicli --history

# 输出示例：
# 历史记录（共 5 条）：
# 
# [5] ✓ 2026-01-13 15:30:00
#     输入: 列出所有txt文件
#     命令: ls *.txt
# 
# [4] ✗ 2026-01-13 15:28:00
#     输入: 删除临时文件
#     命令: rm -rf /tmp/*
#     错误: 用户取消执行

# 重新执行历史命令
aicli --retry 5
# 重新执行历史命令 #5，使用原始输入重新转换

# 搜索历史（未实现，计划中）
# aicli --history --search "git"
```

### 安全特性

```bash
# 危险命令检测示例
aicli "删除所有临时文件"

# 输出：
# ⚠️  检测到潜在危险命令！
# 命令: rm -rf /tmp/*
# 风险: 批量删除文件 (等级: High)
# 
# 是否继续执行？(y/n): 

# 使用 --force 跳过确认（谨慎使用）
aicli --force "删除所有临时文件"

# dry-run 模式（仅查看命令，不执行）
aicli --dry-run "格式化硬盘"
# 输出: 将要执行的命令: mkfs.ext4 /dev/sda1
```

## ❓ 常见问题 (FAQ)

### Q1: API 密钥如何保护？

**A**: 
1. 配置文件设置权限：`chmod 600 ~/.aicli.json`
2. 使用环境变量：`export AICLI_API_KEY=sk-xxxxx`
3. 不要提交配置文件到 Git（已包含在 `.gitignore`）

### Q2: 如何处理敏感数据？

**A**: 使用 `--no-send-stdin` 标志防止将 stdin 数据发送到 LLM：

```bash
cat sensitive-data.txt | aicli --no-send-stdin "统计行数"
# stdin 数据不会发送给 LLM
```

### Q3: 命令转换不准确怎么办？

**A**: 
1. 使用更具体的描述
2. 尝试不同的 LLM 模型（如 gpt-4 vs gpt-3.5-turbo）
3. 使用 `--dry-run` 预览命令
4. 使用 `--verbose` 查看转换过程

```bash
# 不够具体
aicli "查找文件"

# 更具体
aicli "查找当前目录及子目录下所有.go文件，并按修改时间排序"
```

### Q4: 支持哪些操作系统和 Shell？

**A**: 
- **操作系统**: Linux、macOS、Windows
- **Shell**: bash、zsh、PowerShell、cmd
- 自动检测系统 Shell，也可在配置中指定

### Q5: 如何使用本地模型（不依赖外部 API）？

**A**: 使用 Ollama 等本地模型服务：

```bash
# 1. 安装 Ollama
curl https://ollama.ai/install.sh | sh

# 2. 下载模型
ollama pull llama2

# 3. 配置 aicli
{
  "llm": {
    "provider": "local",
    "model": "llama2"
  }
}
```

### Q6: 能否在脚本中使用 aicli？

**A**: 可以，但建议使用 `--force` 跳过确认：

```bash
#!/bin/bash
# 使用 aicli 的脚本示例

# 获取命令但不执行
CMD=$(aicli --dry-run "列出所有go文件" | grep "将要执行" | cut -d: -f2-)
echo "生成的命令: $CMD"

# 或直接执行（跳过确认）
aicli --force "列出所有go文件"
```

### Q7: 如何提高命令转换速度？

**A**: 
1. 使用更快的模型（如 gpt-3.5-turbo）
2. 减少 `llm.timeout` 配置
3. 使用本地模型（Ollama）
4. 检查网络连接

### Q8: 项目的测试覆盖率如何？

**A**: 
```bash
make coverage

# 当前覆盖率：
# pkg/llm:        80.9%
# pkg/safety:     91.4%
# internal/app:   70.3%
# 整体:           65%+
```

### Q9: 如何贡献代码？

**A**: 参阅 [CONTRIBUTING_CN.md](CONTRIBUTING_CN.md)，简要流程：
1. Fork 项目
2. 创建特性分支
3. 编写代码和测试
4. 提交 Pull Request

### Q10: DeepSeek 等其他兼容 OpenAI API 的服务如何配置？

**A**: 只需修改 `api_base` 和 `model`：

```json
{
  "llm": {
    "provider": "openai",
    "api_key": "your-deepseek-key",
    "model": "deepseek-chat",
    "api_base": "https://api.deepseek.com/v1"
  }
}
```

## 🏗️ 项目结构

```
aicli/
├── cmd/aicli/          # 主程序入口
├── pkg/                # 公共库
│   ├── llm/           # LLM 提供商抽象层
│   ├── executor/      # 命令执行引擎
│   ├── config/        # 配置管理
│   └── safety/        # 安全检查
├── internal/           # 私有代码
│   ├── app/           # 应用主逻辑
│   └── history/       # 历史记录管理
├── tests/              # 测试
│   └── integration/   # 集成测试
└── docs/               # 文档
```

## 🧪 开发

### 构建

```bash
# 编译项目
make build

# 运行测试
make test

# 生成覆盖率报告
make coverage

# 检查覆盖率是否达标 (≥60%)
make coverage-check

# 代码格式化
make fmt

# 静态代码检查
make lint

# 清理构建产物
make clean
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/llm/...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成详细的覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 贡献指南

请参阅 [CONTRIBUTING_CN.md](CONTRIBUTING_CN.md) 了解如何参与项目开发。

## 📚 文档

- [架构设计](docs/architecture.md)
- [配置说明](docs/configuration.md)
- [开发指南](docs/development.md)

## 🔐 安全与隐私

- **本地配置**：API 密钥存储在本地配置文件 `~/.aicli.json` 中，请妥善保管文件权限
- **敏感数据保护**：使用 `--no-send-stdin` 选项可避免将标准输入数据发送到 LLM
- **危险命令检测**：自动识别删除、格式化等危险操作，需要用户确认后才执行
- **日志脱敏**：日志中不会记录完整的 API 密钥和敏感命令参数

## 📊 项目状态

- ✅ **阶段 1**: 项目设置完成
- ✅ **阶段 2**: 基础功能完成
- ✅ **阶段 3**: 用户故事 1 完成（基础命令转换与执行）
- ✅ **阶段 4**: 用户故事 2 完成（管道符支持）
- ✅ **阶段 5**: 用户故事 3 完成（交互式确认与历史记录）
- ✅ **阶段 6**: 扩展 LLM 提供商支持完成（Anthropic, 本地模型）
- ✅ **阶段 7**: 完善与文档完成

**当前版本**: v0.1.0-dev  
**测试覆盖率**: ![Coverage](https://img.shields.io/badge/coverage-65%25-brightgreen.svg)（目标：≥60% ✅）

### 核心包覆盖率

| 包 | 覆盖率 | 状态 |
|----|--------|------|
| pkg/llm | 80.9% | ✅ |
| pkg/safety | 91.4% | ✅ |
| pkg/config | 77.4% | ✅ |
| internal/app | 70.3% | ✅ |
| internal/history | 83.5% | ✅ |

## 📄 许可证

本项目采用 [Apache License 2.0](LICENSE) 许可证。

## 🙏 致谢

感谢所有为本项目做出贡献的开发者。

---

**注意**：本项目处于早期开发阶段，功能和 API 可能会发生变化。使用前请仔细阅读文档，谨慎使用涉及文件操作和系统管理的命令。
