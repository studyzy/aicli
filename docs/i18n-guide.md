# AICLI 国际化 (i18n) 指南

> 版本: 1.0 | 日期: 2026-01-13

## 概述

AICLI 支持多语言界面和 LLM 提示词,目前支持中文(zh)和英文(en)。系统会自动检测用户的操作系统语言,并可通过配置文件进行自定义。

## 支持的语言

| 语言代码 | 语言名称 | 覆盖范围 |
|---------|---------|---------|
| `zh` | 中文 | CLI输出、提示词、错误信息、配置向导 |
| `en` | English | CLI output, prompts, errors, configuration wizard |

## 语言检测机制

AICLI 使用以下优先级自动检测语言:

```
1. 配置文件中的 language 字段 (最高优先级)
   ↓
2. LANG 环境变量
   ↓
3. LC_ALL 环境变量
   ↓
4. 默认值: zh (中文)
```

### 示例流程

```bash
# 场景 1: 使用配置文件 (最高优先级)
# ~/.aicli.json
{
  "language": "en"
}
# 结果: 使用英文,忽略环境变量

# 场景 2: 使用环境变量
# 配置文件中未设置 language
export LANG=en_US.UTF-8
aicli "list files"
# 结果: 使用英文

# 场景 3: 使用默认值
# 配置文件中未设置,且无相关环境变量
aicli "列出文件"
# 结果: 使用中文 (默认)
```

## 配置方法

### 方法 1: 配置文件 (推荐)

编辑 `~/.aicli.json`:

```json
{
  "version": "1.0",
  "language": "en",
  "llm": {
    "provider": "openai",
    "api_key": "your-api-key"
  }
}
```

**优点**:
- 永久生效
- 优先级最高
- 不受shell环境影响

### 方法 2: 环境变量

#### 临时设置(单次命令)

```bash
# 使用英文
LANG=en_US.UTF-8 aicli "show all files"

# 使用中文
LANG=zh_CN.UTF-8 aicli "显示所有文件"
```

#### 永久设置(Shell 配置)

**Bash/Zsh** (~/.bashrc 或 ~/.zshrc):

```bash
# 设置为英文
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8

# 或设置为中文
export LANG=zh_CN.UTF-8
export LC_ALL=zh_CN.UTF-8
```

**Fish** (~/.config/fish/config.fish):

```fish
# 设置为英文
set -x LANG en_US.UTF-8
set -x LC_ALL en_US.UTF-8
```

应用配置:

```bash
source ~/.bashrc  # 或 ~/.zshrc
```

### 方法 3: 系统级设置

修改系统语言会影响所有应用,包括 AICLI。

## Locale 格式支持

AICLI 支持标准 POSIX locale 格式,自动提取语言代码:

| Locale 格式 | 识别结果 |
|------------|---------|
| `zh_CN.UTF-8` | `zh` (中文) |
| `zh_CN` | `zh` (中文) |
| `zh` | `zh` (中文) |
| `en_US.UTF-8` | `en` (英文) |
| `en_US` | `en` (英文) |
| `en_GB.utf8` | `en` (英文) |
| `en` | `en` (英文) |
| `fr_FR.UTF-8` | `zh` (不支持,fallback到中文) |

## 国际化覆盖范围

### 1. CLI 界面输出

所有用户可见的输出都已国际化:

**中文示例**:
```bash
$ aicli --history
历史记录(共 5 条):

[1] ✓ 2026-01-13 10:30:00
    输入: 列出文件
    命令: ls -la
```

**英文示例**:
```bash
$ LANG=en_US.UTF-8 aicli --history
History (5 entries):

[1] ✓ 2026-01-13 10:30:00
    Input: list files
    Command: ls -la
```

### 2. LLM 提示词 (System Prompt)

LLM 收到的系统提示词会根据当前语言自动切换:

**中文提示词**:
```
你是一个命令行助手,专门将用户的自然语言描述转换为可执行的 shell 命令。

规则:
1. 只返回命令本身,不要有任何解释或说明
2. 不要使用 markdown 代码块格式
...
```

**英文提示词**:
```
You are a command-line assistant that converts natural language 
descriptions into executable shell commands.

Rules:
1. Return only the command itself, without any explanation or description
2. Do not use markdown code block format
...
```

### 3. 错误信息

**中文**:
```bash
$ aicli
错误: 请提供自然语言描述
```

**英文**:
```bash
$ LANG=en_US.UTF-8 aicli
Error: Please provide natural language description
```

### 4. 交互式确认

**中文 - 危险命令确认**:
```bash
$ aicli "删除所有临时文件"

⚠️  检测到潜在危险命令！
命令: rm -rf /tmp/*
风险: 递归删除操作 (等级: 高)

是否继续执行?(y/n): 
```

**英文 - 危险命令确认**:
```bash
$ LANG=en_US.UTF-8 aicli "delete all temp files"

⚠️  Potentially dangerous command detected!
Command: rm -rf /tmp/*
Risk: Recursive deletion (Level: High)

Continue execution? (y/n):
```

### 5. 配置向导

**中文**:
```bash
$ aicli init
欢迎使用 aicli 配置向导！
我们将引导您完成基本配置。

请选择 LLM 提供商:
1. OpenAI (GPT-4, GPT-3.5)
2. Anthropic (Claude)
...
```

**英文**:
```bash
$ LANG=en_US.UTF-8 aicli init
Welcome to aicli configuration wizard!
We will guide you through the basic configuration.

Please select LLM provider:
1. OpenAI (GPT-4, GPT-3.5)
2. Anthropic (Claude)
...
```

### 6. Verbose 模式输出

**中文**:
```bash
$ aicli --verbose "列出文件"
自然语言输入: 列出文件
执行上下文: OS: darwin, Shell: zsh, 工作目录: /Users/user
转换后的命令: ls -la
转换耗时: 1.2s
开始执行命令...
执行耗时: 0.05s
总耗时: 1.25s
```

**英文**:
```bash
$ LANG=en_US.UTF-8 aicli --verbose "list files"
Natural language input: list files
Execution context: OS: darwin, Shell: zsh, WorkDir: /Users/user
Translated command: ls -la
Translation time: 1.2s
Executing command...
Execution time: 0.05s
Total time: 1.25s
```

## 验证当前语言

### 方法 1: 通过帮助信息

```bash
# 中文环境
$ aicli --help | head -1
aicli 是一个让命令行支持自然语言操作的工具。

# 英文环境
$ LANG=en_US.UTF-8 aicli --help | head -1
aicli is a tool that brings natural language operations to the command line.
```

### 方法 2: 通过错误信息

```bash
# 中文环境
$ aicli
错误: 请提供自然语言描述

# 英文环境  
$ LANG=en_US.UTF-8 aicli
Error: Please provide natural language description
```

## 常见问题

### Q1: 为什么设置了 LANG 还是显示中文?

**原因**: 配置文件中的 `language` 字段优先级更高。

**解决方案**:
1. 检查 `~/.aicli.json` 中是否有 `language` 字段
2. 删除该字段或修改为期望的语言
3. 或使用 `--config` 指定不含 language 字段的配置文件

### Q2: 如何验证环境变量是否生效?

```bash
# 查看当前 LANG 设置
echo $LANG

# 查看当前 LC_ALL 设置
echo $LC_ALL

# 查看所有 locale 相关变量
locale
```

### Q3: 支持其他语言吗?

目前仅支持中文(zh)和英文(en)。不支持的语言会自动 fallback 到中文。

**计划支持的语言** (未来版本):
- 日语 (ja)
- 韩语 (ko)
- 法语 (fr)
- 德语 (de)
- 西班牙语 (es)

### Q4: LLM 提示词语言会影响命令质量吗?

**不会**。我们针对中英文分别优化了提示词,确保命令生成质量一致:

- 中文提示词 → 理解中文输入 → 生成准确命令
- 英文提示词 → 理解英文输入 → 生成准确命令

### Q5: 可以混用中英文吗?

可以,但不推荐。建议保持输入语言与界面语言一致以获得最佳体验:

```bash
# 推荐: 中文界面 + 中文输入
aicli "列出所有txt文件"

# 推荐: 英文界面 + 英文输入
LANG=en_US.UTF-8 aicli "list all txt files"

# 不推荐: 中文界面 + 英文输入
aicli "list all txt files"
```

### Q6: 配置向导在 i18n 初始化前运行怎么办?

`aicli init` 命令使用双语提示(中英文并列),确保在任何环境下都能理解。

## 技术实现

### 架构设计

AICLI 采用轻量级自研 i18n 框架,无外部依赖:

```
pkg/i18n/
├── i18n.go          # 核心逻辑: 语言检测、翻译
├── keys.go          # 翻译键常量 (100+ 个)
├── messages_zh.go   # 中文翻译资源
└── messages_en.go   # 英文翻译资源
```

### 使用示例(开发者)

```go
package main

import (
    "github.com/studyzy/aicli/pkg/i18n"
    "github.com/studyzy/aicli/pkg/config"
)

func main() {
    // 1. 加载配置
    cfg := config.Load()
    
    // 2. 初始化 i18n
    i18n.Init(cfg)
    
    // 3. 使用翻译
    fmt.Println(i18n.T(i18n.MsgWelcome))
    // 中文: "欢迎使用 aicli"
    // 英文: "Welcome to aicli"
    
    // 4. 带参数的翻译
    fmt.Println(i18n.T(i18n.MsgHistoryCount, 10))
    // 中文: "历史记录(共 10 条):"
    // 英文: "History (10 entries):"
}
```

### 性能影响

- **二进制大小增加**: <50KB
- **内存占用**: ~100KB (加载翻译表)
- **运行时开销**: 几乎为零 (简单 map 查找)
- **启动时间影响**: <1ms

## 贡献翻译

欢迎贡献新语言翻译！

### 步骤

1. 在 `pkg/i18n/` 创建新文件: `messages_xx.go` (xx为语言代码)
2. 复制 `messages_en.go` 的结构
3. 翻译所有键值对
4. 在 `i18n.go` 的 `NewLocalizer()` 添加新语言支持
5. 运行测试: `go test ./pkg/i18n/...`
6. 提交 Pull Request

### 翻译规范

- 保持原文的语气和风格
- 技术术语保持一致(如: "Shell", "LLM")
- 简洁明了,避免冗长
- 考虑目标用户的文化背景

## 参考资料

- [配置文件说明](configuration.md)
- [POSIX Locale 规范](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap07.html)
- [Unicode CLDR](https://cldr.unicode.org/)

---

如有问题或建议,请提交 [Issue](https://github.com/studyzy/aicli/issues)。
