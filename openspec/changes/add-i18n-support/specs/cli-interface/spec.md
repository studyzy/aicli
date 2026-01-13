# 功能规范:命令行界面国际化

## 修改需求

### 需求:命令行输出国际化

所有命令行输出文本必须根据用户语言显示对应翻译。

#### 场景:错误信息国际化
- **当** 用户语言为英文
- **并且** 加载配置文件失败
- **那么** 应输出英文错误信息:"Error: Failed to load configuration: <详细错误>"
- **并且** 错误前缀为 "Error:"

#### 场景:成功信息国际化
- **当** 用户语言为中文
- **并且** 配置初始化成功
- **那么** 应输出:"配置已成功保存到 ~/.aicli.json"

#### 场景:详细输出模式国际化
- **当** 用户使用 `--verbose` 标志
- **并且** 语言为英文
- **那么** 应输出:"Natural language input: list all files\nExecution context: OS: linux, Shell: bash, WorkDir: /home/user"

### 需求:交互式提示国际化

所有交互式提示(如确认、输入请求)必须国际化。

#### 场景:危险命令确认提示(英文)
- **当** 检测到危险命令
- **并且** 用户语言为英文
- **那么** 应显示:"\n⚠️  Potentially dangerous command detected!\nCommand: rm -rf /\nRisk: <描述> (Level: <等级>)\n\nContinue execution? (y/n): "

#### 场景:危险命令确认提示(中文)
- **当** 检测到危险命令
- **并且** 用户语言为中文
- **那么** 应显示:"\n⚠️  检测到潜在危险命令!\n命令: rm -rf /\n风险: <描述> (等级: <等级>)\n\n是否继续执行?(y/n): "

#### 场景:用户输入提示国际化
- **当** 配置向导请求用户输入
- **并且** 语言为英文
- **那么** 提示应为:"Please enter API Key: " 或 "Please enter API Key [default]: "

### 需求:历史记录显示国际化

历史记录查看功能的所有文本必须国际化。

#### 场景:历史记录列表显示(英文)
- **当** 用户执行 `aicli --history`
- **并且** 语言为英文
- **那么** 应输出:"History (10 entries):\n\n[1] ✓ 2026-01-13 10:30:00\n    Input: list all files\n    Command: ls\n"

#### 场景:历史记录列表显示(中文)
- **当** 用户执行 `aicli --history`
- **并且** 语言为中文
- **那么** 应输出:"历史记录(共 10 条):\n\n[1] ✓ 2026-01-13 10:30:00\n    输入: 列出所有文件\n    命令: ls\n"

#### 场景:空历史记录提示
- **当** 历史记录为空
- **并且** 语言为英文
- **那么** 应输出:"No history records"

### 需求:配置向导国际化

`aicli init` 配置向导的所有交互文本必须国际化。

#### 场景:配置向导欢迎信息(英文)
- **当** 用户执行 `aicli init`
- **并且** 语言为英文
- **那么** 应显示:"Welcome to aicli configuration wizard!\nWe will guide you through the basic configuration.\n"

#### 场景:配置向导选项提示(中文)
- **当** 配置向导请求选择 LLM 提供商
- **并且** 语言为中文
- **那么** 应显示:"请选择 LLM 提供商:\n1. OpenAI (GPT-4, GPT-3.5)\n2. Anthropic (Claude)\n..."

#### 场景:配置保存成功提示
- **当** 配置文件保存成功
- **并且** 语言为英文
- **那么** 应输出:"Configuration successfully saved to ~/.aicli.json\nYou can now start using aicli!\nExample: aicli \"check my public IP\""

### 需求:帮助信息国际化

命令行帮助信息(`--help`)必须国际化。

#### 场景:根命令帮助(英文)
- **当** 用户执行 `aicli --help`
- **并且** 语言为英文
- **那么** 应显示英文帮助文本,包括:
  - Short: "AI command-line assistant"
  - Long: 功能描述和使用示例(英文)
  - Flags 说明(英文)

#### 场景:根命令帮助(中文)
- **当** 用户执行 `aicli --help`
- **并且** 语言为中文
- **那么** 应显示当前的中文帮助文本

#### 场景:命令行标志说明国际化
- **当** 帮助信息包含标志说明
- **并且** 语言为英文
- **那么** 标志说明应为:
  - `--config`: "Configuration file path"
  - `--verbose`: "Show detailed output"
  - `--dry-run`: "Show command without executing"
  - `--force`: "Force execution, skip confirmation"

### 需求:提示和通知信息国际化

所有提示性和通知性信息必须国际化。

#### 场景:配置文件不存在提示(英文)
- **当** 配置文件不存在
- **并且** 语言为英文
- **那么** 应输出:"Note: Configuration file ~/.aicli.json does not exist.\nYou can run 'aicli init' for quick setup.\n"

#### 场景:API Key 警告信息
- **当** 配置向导中 API Key 为空
- **并且** 语言为英文
- **那么** 应显示:"Warning: API Key is empty. You may need to set it in the AICLI_API_KEY environment variable."

#### 场景:重试命令提示
- **当** 用户使用 `--retry` 功能
- **并且** 语言为英文
- **那么** 应输出:"Retrying history command #3:\n  Input: list files\n  Command: ls\n"

### 需求:错误类型标准化

不同类型的错误必须使用一致的国际化前缀。

#### 场景:错误前缀一致性(英文)
- **当** 发生任何错误
- **并且** 语言为英文
- **那么** 错误信息应以 "Error:" 开头
- **并且** stderr 输出格式为:"Error: <error message>"

#### 场景:错误前缀一致性(中文)
- **当** 发生任何错误
- **并且** 语言为中文
- **那么** 错误信息应以 "错误:" 开头

### 需求:Cobra 框架集成

Cobra 命令行框架的文本必须国际化。

#### 场景:Cobra Use 字段国际化
- **当** 定义根命令
- **并且** 语言为英文
- **那么** `Use` 字段应为:"aicli [natural language description]"

#### 场景:Cobra Short 和 Long 描述
- **当** 用户查看命令帮助
- **那么** `Short` 和 `Long` 字段应根据用户语言显示对应翻译
- **并且** 使用示例应本地化

### 需求:日期时间格式保持不变

历史记录等功能中的日期时间格式必须不做国际化。

#### 场景:时间戳格式统一
- **当** 显示历史记录时间戳
- **那么** 应始终使用格式:"2006-01-02 15:04:05"
- **并且** 不根据语言切换日期格式(保持 ISO 8601 风格)
