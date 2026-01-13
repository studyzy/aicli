# 功能规范:LLM 提示词国际化

## 修改需求

### 需求:系统提示词国际化

系统提示词(System Prompt)必须根据用户语言生成对应语言的版本。

#### 场景:中文系统提示词
- **当** 用户语言设置为中文(`zh`)
- **那么** `GetSystemPrompt()` 应返回中文系统提示词
- **并且** 提示词内容包括:"你是一个命令行助手,专门将用户的自然语言描述转换为可执行的 shell 命令"
- **并且** 规则说明使用中文

#### 场景:英文系统提示词
- **当** 用户语言设置为英文(`en`)
- **那么** `GetSystemPrompt()` 应返回英文系统提示词
- **并且** 提示词内容包括:"You are a command-line assistant that converts natural language descriptions into executable shell commands"
- **并且** 规则说明使用英文

#### 场景:执行上下文的国际化描述
- **当** 系统提示词包含执行上下文信息(OS, Shell, WorkDir)
- **并且** 用户语言为英文
- **那么** 上下文说明应为:"Execution Environment:\n- Operating System: linux\n- Shell: bash\n- Working Directory: /home/user"
- **并且** 字段名使用对应语言

### 需求:用户提示词国际化

用户提示词构建时的模板文本必须国际化。

#### 场景:中文用户提示词模板
- **当** 用户语言为中文
- **并且** 调用 `BuildPrompt("列出所有文件", ctx)`
- **那么** 应生成包含中文模板的提示词:"将以下自然语言描述转换为命令:\n列出所有文件"

#### 场景:英文用户提示词模板
- **当** 用户语言为英文
- **并且** 调用 `BuildPrompt("list all files", ctx)`
- **那么** 应生成包含英文模板的提示词:"Convert the following natural language description into a command:\nlist all files"

#### 场景:标准输入数据提示国际化
- **当** 用户提示词包含标准输入数据
- **并且** 用户语言为英文
- **那么** 应显示:"Standard input data:\n..."
- **并且** 截断提示为:"... (truncated)"

### 需求:上下文描述国际化

执行上下文的调试描述必须国际化。

#### 场景:中文上下文描述
- **当** 调用 `BuildContextDescription(ctx)` 且语言为中文
- **那么** 应返回:"OS: linux, Shell: bash, 工作目录: /home/user"
- **并且** 关键字段名保持英文(OS, Shell),描述性字段使用中文

#### 场景:英文上下文描述
- **当** 调用 `BuildContextDescription(ctx)` 且语言为英文
- **那么** 应返回:"OS: linux, Shell: bash, WorkDir: /home/user"
- **并且** 所有字段名使用英文缩写

### 需求:LLM 提示词质量保证

国际化后的提示词必须保持原有的命令生成质量。

#### 场景:英文提示词生成准确命令
- **当** 用户使用英文环境
- **并且** 输入自然语言描述:"list all txt files in current directory"
- **那么** LLM 应返回正确的命令:`ls *.txt` 或 `find . -name "*.txt"`
- **并且** 命令生成准确率与中文提示词相当

#### 场景:中文提示词保持现有质量
- **当** 用户使用中文环境
- **并且** 输入自然语言描述:"列出当前目录的所有txt文件"
- **那么** LLM 应返回正确的命令
- **并且** 命令生成质量与国际化前完全一致

### 需求:提示词模板可维护性

提示词翻译必须易于维护和更新。

#### 场景:提示词模板集中管理
- **当** 开发者需要更新系统提示词
- **那么** 应在 `pkg/llm/prompt.go` 中统一修改
- **并且** 中英文提示词逻辑分离,互不影响

#### 场景:提示词规则一致性
- **当** 系统提示词包含规则列表(如"只返回命令本身")
- **那么** 中英文版本的规则条目数量应一致
- **并且** 规则语义完全对应
