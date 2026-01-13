# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-14

### Added
- **Real-time command output**: Commands now display their progress and output in real-time instead of waiting until completion
- **Command display**: By default, translated commands are shown to users (output to stderr) to help them learn shell commands
- **Quiet mode**: Added `-q/--quiet` flag to suppress command display for clean output in scripts and pipes
- Natural language to shell command translation
- Multiple LLM provider support (OpenAI, Anthropic, local models via Ollama, DeepSeek)
- Pipe support for stdin/stdout
- Safety checks for dangerous commands (delete, format operations)
- Command history with retry functionality
- Internationalization (i18n) support for Chinese and English
- Interactive configuration wizard (`aicli init`)
- Cross-platform support (Linux, macOS, Windows)
- Comprehensive test suite with 65%+ code coverage

### Changed
- Improved executor to support interactive mode with real-time output
- Enhanced output stream handling (stdout for data, stderr for messages)
- Updated documentation with detailed examples of output stream behavior

### Technical Details
- Added `ExecuteInteractive()` and `ExecuteWithOutput()` methods to executor
- Command prompts output to stderr by default (can be disabled with `-q`)
- Real-time output for commands with progress indicators (downloads, installations, etc.)

## [0.1.0-dev] - 2025-12-01

### Added
- Initial development version
- Core functionality implementation
- Basic LLM integration
- Command execution engine
- Safety checks
- History management

---

[1.0.0]: https://github.com/studyzy/aicli/releases/tag/v1.0.0
[0.1.0-dev]: https://github.com/studyzy/aicli/commits/main
