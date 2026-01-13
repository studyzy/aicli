# aicli 配置文件说明

> 版本: 1.0 | 日期: 2026-01-13

## 概述

aicli 使用 JSON 格式的配置文件来管理 LLM API 密钥、执行选项、安全设置等。配置文件默认位于用户主目录下的 `~/.aicli.json`。

## 配置文件位置

### 默认位置

- Linux/macOS: `~/.aicli.json`
- Windows: `%USERPROFILE%\.aicli.json`

### 自定义位置

使用 `--config` 标志指定自定义配置文件：

```bash
aicli --config /path/to/custom.aicli.json "列出文件"
```

## 配置文件结构

完整的配置文件示例：

```json
{
  "version": "1.0",
  "language": "zh",
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
    "dangerous_patterns": ["rm -rf", "format", "mkfs"],
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

## 配置项详解

### 1. version (版本)

**类型**: `string`  
**必需**: 否  
**默认值**: `"1.0"`

配置文件的版本号，用于未来的配置迁移。

### 2. language (语言设置)

**类型**: `string`  
**必需**: 否  
**默认值**: 自动检测 (从 `LANG` 或 `LC_ALL` 环境变量)  
**可选值**: `zh` (中文), `en` (英文)

界面显示语言和 LLM 提示词语言。

**语言检测优先级**:
1. 配置文件中的 `language` 字段 (最高优先级)
2. `LANG` 环境变量
3. `LC_ALL` 环境变量  
4. 默认值: `zh` (中文)

**示例**:

```json
{
  "language": "en"
}
```

或通过环境变量:

```bash
# 使用英文
export LANG=en_US.UTF-8
aicli "list files"

# 使用中文
export LANG=zh_CN.UTF-8
aicli "列出文件"
```

详细说明请参阅 [国际化指南](i18n-guide.md)。

### 3. llm (LLM 配置)

#### llm.provider (提供商)

**类型**: `string`  
**必需**: 是  
**默认值**: `"openai"`  
**可选值**: `openai`, `anthropic`, `claude`, `local`, `ollama`, `mock`

指定使用的 LLM 提供商。

**示例**:

```json
{
  "llm": {
    "provider": "openai"
  }
}
```

#### llm.api_key (API 密钥)

**类型**: `string`  
**必需**: 是（除了 `local` 和 `mock` 提供商）  
**默认值**: 无

LLM 服务的 API 密钥。

**安全建议**:
- 设置配置文件权限为 600: `chmod 600 ~/.aicli.json`
- 或使用环境变量: `export AICLI_API_KEY=sk-xxxxx`

**示例**:

```json
{
  "llm": {
    "provider": "openai",
    "api_key": "sk-xxxxxxxxxxxxx"
  }
}
```

#### llm.api_base (API 基础URL)

**类型**: `string`  
**必需**: 否  
**默认值**: 根据提供商自动设置

自定义 API 端点，用于兼容 OpenAI API 的服务。

**默认值**:
- OpenAI: `https://api.openai.com/v1`
- Anthropic: `https://api.anthropic.com/v1`
- Local: `http://localhost:11434`

**使用场景**:
- 使用 DeepSeek 等兼容服务
- 使用企业内部代理
- 使用自托管服务

**示例 (DeepSeek)**:

```json
{
  "llm": {
    "provider": "openai",
    "api_key": "sk-xxxxx",
    "api_base": "https://api.deepseek.com/v1",
    "model": "deepseek-chat"
  }
}
```

#### llm.model (模型)

**类型**: `string`  
**必需**: 否  
**默认值**: 根据提供商自动设置

指定使用的具体模型。

**默认值**:
- OpenAI: `gpt-4`
- Anthropic: `claude-3-sonnet-20240229`
- Local: `llama2`

**常用模型**:

| 提供商 | 推荐模型 | 说明 |
|--------|----------|------|
| OpenAI | `gpt-4` | 最强性能，成本较高 |
| OpenAI | `gpt-3.5-turbo` | 性价比高 |
| Anthropic | `claude-3-sonnet-20240229` | 平衡性能和速度 |
| Anthropic | `claude-3-opus-20240229` | 最高性能 |
| Local | `llama2` | 本地运行，无成本 |
| Local | `codellama` | 代码优化 |

#### llm.timeout (超时时间)

**类型**: `int`  
**必需**: 否  
**默认值**: `10`  
**单位**: 秒

LLM API 请求的超时时间。

**建议**:
- 网络稳定: 10 秒
- 网络不稳定: 20-30 秒
- 本地模型: 60 秒（首次推理较慢）

#### llm.max_tokens (最大令牌数)

**类型**: `int`  
**必需**: 否  
**默认值**: `500`

LLM 响应的最大令牌数。

**说明**: 命令通常很短，500 足够。增加此值不会提高质量，但会增加成本。

### 4. execution (执行配置)

#### execution.auto_confirm (自动确认)

**类型**: `bool`  
**必需**: 否  
**默认值**: `false`

是否自动确认危险命令（不推荐）。

**警告**: 设置为 `true` 会跳过所有安全确认！

#### execution.dry_run_default (默认 Dry-run)

**类型**: `bool`  
**必需**: 否  
**默认值**: `false`

是否默认启用 dry-run 模式（仅显示命令不执行）。

**使用场景**: 在测试或学习阶段启用。

#### execution.timeout (执行超时)

**类型**: `int`  
**必需**: 否  
**默认值**: `30`  
**单位**: 秒

命令执行的最大时长。

**建议**:
- 快速命令: 10-30 秒
- 长时间任务: 300+ 秒

#### execution.shell (Shell 类型)

**类型**: `string`  
**必需**: 否  
**默认值**: `"auto"`  
**可选值**: `auto`, `bash`, `zsh`, `powershell`, `cmd`

指定使用的 Shell 类型。

**说明**: 
- `auto`: 自动检测系统默认 Shell（推荐）
- 其他值: 强制使用指定 Shell

### 5. safety (安全配置)

#### safety.enable_checks (启用检查)

**类型**: `bool`  
**必需**: 否  
**默认值**: `true`

是否启用危险命令检测。

**警告**: 禁用会失去安全保护！

#### safety.dangerous_patterns (危险模式)

**类型**: `array<string>`  
**必需**: 否  
**默认值**: 内置模式列表

自定义危险命令模式列表（追加到内置模式）。

**示例**:

```json
{
  "safety": {
    "dangerous_patterns": [
      "rm -rf",
      "format",
      "mkfs",
      "dd if=",
      "chmod 777"
    ]
  }
}
```

#### safety.require_confirmation (需要确认)

**类型**: `bool`  
**必需**: 否  
**默认值**: `true`

检测到危险命令时是否需要用户确认。

**说明**: 即使设为 `false`，仍会显示警告。

### 6. history (历史配置)

#### history.enabled (启用历史)

**类型**: `bool`  
**必需**: 否  
**默认值**: `true`

是否记录命令历史。

#### history.max_entries (最大条目数)

**类型**: `int`  
**必需**: 否  
**默认值**: `1000`

保留的最大历史记录数。

**说明**: 超过限制时，最旧的记录会被删除。

#### history.file (历史文件路径)

**类型**: `string`  
**必需**: 否  
**默认值**: `"~/.aicli_history.json"`

历史记录文件的路径。

### 7. logging (日志配置)

#### logging.enabled (启用日志)

**类型**: `bool`  
**必需**: 否  
**默认值**: `false`

是否启用日志记录。

**说明**: 目前日志功能未完全实现。

#### logging.level (日志级别)

**类型**: `string`  
**必需**: 否  
**默认值**: `"info"`  
**可选值**: `debug`, `info`, `warn`, `error`

日志记录的详细程度。

#### logging.file (日志文件)

**类型**: `string`  
**必需**: 否  
**默认值**: `""`

日志文件路径。空字符串表示不记录到文件。

## 配置优先级

当同一个配置项有多个来源时，优先级顺序为：

1. **命令行标志** (最高优先级)
   ```bash
   aicli --force "删除文件"
   ```

2. **环境变量**
   ```bash
   export AICLI_API_KEY=sk-xxxxx
   export AICLI_MODEL=gpt-3.5-turbo
   ```

3. **配置文件**
   ```json
   {"llm": {"api_key": "sk-xxxxx"}}
   ```

4. **默认值** (最低优先级)

## 环境变量

支持的环境变量：

| 环境变量 | 对应配置项 | 说明 |
|----------|------------|------|
| `AICLI_API_KEY` | `llm.api_key` | LLM API 密钥 |
| `AICLI_MODEL` | `llm.model` | 模型名称 |
| `AICLI_PROVIDER` | `llm.provider` | 提供商名称 |

## 常见配置示例

### OpenAI GPT-4

```json
{
  "llm": {
    "provider": "openai",
    "api_key": "sk-xxxxx",
    "model": "gpt-4"
  }
}
```

### Anthropic Claude

```json
{
  "llm": {
    "provider": "anthropic",
    "api_key": "sk-ant-xxxxx",
    "model": "claude-3-sonnet-20240229"
  }
}
```

### 本地 Ollama

```json
{
  "llm": {
    "provider": "local",
    "model": "llama2",
    "api_base": "http://localhost:11434"
  }
}
```

### DeepSeek

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

### 安全优先配置

```json
{
  "llm": {
    "provider": "openai",
    "api_key": "sk-xxxxx"
  },
  "execution": {
    "dry_run_default": true
  },
  "safety": {
    "enable_checks": true,
    "require_confirmation": true
  }
}
```

## 配置文件管理

### 创建配置文件

```bash
# 复制示例配置
cp example.aicli.json ~/.aicli.json

# 编辑配置
nano ~/.aicli.json

# 设置权限（保护 API 密钥）
chmod 600 ~/.aicli.json
```

### 验证配置

```bash
# 使用 verbose 模式查看配置加载情况
aicli --verbose "测试命令"
```

### 多配置管理

为不同场景创建多个配置文件：

```bash
# 开发配置
~/.aicli.dev.json

# 生产配置
~/.aicli.prod.json

# 使用时指定
aicli --config ~/.aicli.dev.json "命令"
```

## 故障排查

### 配置文件未加载

**问题**: aicli 没有读取配置文件

**排查**:
1. 确认文件位置正确: `ls -la ~/.aicli.json`
2. 检查文件权限: 应可读
3. 验证 JSON 格式: 使用 `jq . ~/.aicli.json`

### API 密钥错误

**问题**: 提示 API 密钥未配置或无效

**排查**:
1. 确认密钥已设置: `grep api_key ~/.aicli.json`
2. 验证密钥格式正确（不含空格、换行）
3. 尝试使用环境变量: `export AICLI_API_KEY=sk-xxxxx`

### 性能问题

**问题**: 命令转换很慢

**调整**:
1. 减少 `llm.timeout` 检测问题
2. 尝试不同的 `llm.model`（如 gpt-3.5-turbo）
3. 检查网络连接

## 安全建议

1. **保护配置文件**: `chmod 600 ~/.aicli.json`
2. **不要提交到 Git**: 已包含在 `.gitignore`
3. **使用环境变量**: 用于 CI/CD 环境
4. **定期轮换密钥**: 定期更新 API 密钥
5. **最小权限原则**: 仅启用必需的功能

## 更新配置

配置文件可以随时编辑，下次运行 aicli 时自动生效，无需重启或重新加载。

## 参考

- [架构设计文档](./architecture.md)
- [开发指南](./development.md)
- [示例配置文件](../example.aicli.json)
