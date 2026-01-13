# 更新日志

本文件记录项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
版本号遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [1.0.0] - 2026-01-14

### 新增功能
- **实时命令输出**：命令执行时实时显示进度和输出，不再等待命令完成后才显示结果
- **命令显示**：默认向用户显示翻译后的命令（输出到 stderr），帮助用户学习 shell 命令
- **静默模式**：新增 `-q/--quiet` 参数，可以隐藏命令显示，适合脚本和管道场景
- 自然语言转 shell 命令
- 多 LLM 提供商支持（OpenAI、Anthropic、本地模型 Ollama、DeepSeek）
- 完整的管道符支持（stdin/stdout）
- 危险命令安全检查（删除、格式化等操作）
- 命令历史记录与重试功能
- 国际化（i18n）支持中英文
- 交互式配置向导（`aicli init`）
- 跨平台支持（Linux、macOS、Windows）
- 完善的测试覆盖（65%+ 代码覆盖率）

### 变更
- 改进执行器以支持实时输出的交互模式
- 增强输出流处理（stdout 用于数据，stderr 用于消息）
- 更新文档，详细说明输出流行为和使用示例

### 技术细节
- 添加 `ExecuteInteractive()` 和 `ExecuteWithOutput()` 方法到执行器
- 命令提示默认输出到 stderr（可通过 `-q` 关闭）
- 对带有进度指示器的命令（下载、安装等）提供实时输出支持

## [0.1.0-dev] - 2025-12-01

### 新增功能
- 初始开发版本
- 核心功能实现
- 基础 LLM 集成
- 命令执行引擎
- 安全检查
- 历史记录管理

---

[1.0.0]: https://github.com/studyzy/aicli/releases/tag/v1.0.0
[0.1.0-dev]: https://github.com/studyzy/aicli/commits/main
