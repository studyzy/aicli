# AICLI - AI Command Line Assistant

[![Test and Coverage](https://github.com/studyzy/aicli/actions/workflows/test.yml/badge.svg)](https://github.com/studyzy/aicli/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/studyzy/aicli)](https://goreportcard.com/report/github.com/studyzy/aicli)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

English | [中文](README_CN.md)

AICLI is a Go tool that brings natural-language operations to your terminal. You describe what you want to do in plain language, and AICLI uses an LLM provider to convert it into a shell command and (optionally) execute it.

## Features

- **Natural language → shell command**: describe the action you want, and get a generated command
- **Pipe-friendly**: works with stdin/stdout, so it composes well with other CLI tools
- **Safety confirmations**: detects risky commands (e.g., bulk delete/format) and asks before executing
- **Command history**: stores past prompts/commands and supports retry
- **Multiple LLM providers**: OpenAI, Anthropic, local models, and other OpenAI-compatible APIs
- **Internationalization (i18n)**: supports Chinese and English with automatic detection from OS locale
- **Cross-platform**: Linux, macOS, and Windows

## Quick start

### Install

```bash
# Option 1: build from source
git clone https://github.com/studyzy/aicli.git
cd aicli
make build
make install

# Option 2: go install
go install github.com/studyzy/aicli/cmd/aicli@latest
```

### Configure

#### Interactive setup (recommended)

Run:

```bash
aicli init
```

This will guide you through choosing an LLM provider and setting your API key.

#### Manual setup

Create `~/.aicli.json`:

```json
{
  "version": "1.0",
  "language": "en",
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

**Language setting**: The `language` field is optional. Supported values:
- `"zh"` - Chinese (中文)
- `"en"` - English

If not set, AICLI automatically detects your system locale from `LANG` or `LC_ALL` environment variables. Default is Chinese.

You can also set the API key via environment variable:

```bash
export AICLI_API_KEY="your-api-key-here"
```

## Basic usage

```bash
# Example 1: search within a file
aicli "find ERROR in log.txt"
# 💡 Executing: grep "ERROR" log.txt
# -> grep "ERROR" log.txt

# Example 2: file listing
aicli "show all .txt files in current directory"
# 💡 Executing: ls *.txt
# -> ls *.txt (or find . -name "*.txt")

# Example 3: pipe input through aicli
cat log.txt | aicli "filter lines containing ERROR"

# Example 4: chain with other commands
aicli "list all txt files" | wc -l

# Example 5: view history
aicli --history

# Example 6: retry a history item
aicli --retry 3
```

**Note**: By default, AICLI displays the translated command before execution (output to stderr). This helps you learn the actual shell commands. The command prompt won't interfere with pipes or output redirection.

### Common CLI options

```bash
# Print the generated command only (do not execute)
aicli --dry-run "delete temp files"

# Show detailed conversion process
aicli --verbose "find all go files"

# Skip safety confirmation (use carefully)
aicli --force "delete all temp files"

# Quiet mode: hide the translated command (only show output)
aicli --quiet "list files"
# or use the short form
aicli -q "list files"

# Do not send stdin to the LLM (privacy)
cat sensitive.txt | aicli --no-send-stdin "count lines"
```

### Understanding output streams

AICLI follows Unix conventions for output streams:

- **stdout**: Command execution results (for piping to other commands)
- **stderr**: Status messages, translated commands, errors

```bash
# Example: Command prompt goes to stderr, results go to stdout
$ aicli "list numbers 1 to 5"
💡 Executing: seq 1 5    # This is on stderr
1                         # These lines are on stdout
2
3
4
5

# Piping works correctly because only stdout is passed
$ aicli "list numbers 1 to 5" | aicli "calculate sum"
💡 Executing: seq 1 5           # stderr: visible to user
💡 Executing: awk '{s+=$1} END {print s}'  # stderr: visible to user
15                              # stdout: piped between commands

# Hide command prompts if needed
$ aicli -q "list numbers 1 to 5" | wc -l
5                               # Clean output, no prompts

# Redirect only output, keep command prompts visible
$ aicli "list files" > output.txt
💡 Executing: ls                # Still visible on screen
                                 # File list saved to output.txt
```

## LLM providers

Switch providers by updating your config.

### OpenAI (GPT)

```json
{
  "llm": {
    "provider": "openai",
    "api_key": "sk-xxxxx",
    "model": "gpt-4"
  }
}
```

### Anthropic (Claude)

```json
{
  "llm": {
    "provider": "anthropic",
    "api_key": "sk-ant-xxxxx",
    "model": "claude-3-sonnet-20240229"
  }
}
```

### Local models (Ollama)

```json
{
  "llm": {
    "provider": "local",
    "model": "llama2",
    "api_base": "http://localhost:11434"
  }
}
```

### DeepSeek (OpenAI-compatible API)

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

## Project structure

```text
aicli/
├── cmd/aicli/          # CLI entry
├── pkg/                # shared packages
│   ├── llm/            # LLM provider abstractions
│   ├── executor/       # command execution engine
│   ├── config/         # configuration management
│   └── safety/         # safety checks
├── internal/           # internal app logic
│   ├── app/            # core workflow
│   └── history/        # history store
├── tests/              # tests
│   └── integration/    # integration tests
└── docs/               # docs
```

## Development

```bash
# Build
make build

# Test
make test

# Coverage
make coverage
make coverage-check

# Format
make fmt

# Lint
make lint

# Clean
make clean
```

You can also run Go tests directly:

```bash
go test ./...
```

## Documentation

- [Architecture](docs/architecture.md)
- [Configuration](docs/configuration.md)
- [Internationalization Guide](docs/i18n-guide.md)
- [Development guide](docs/development.md)

## Security & privacy

- **Local config**: API keys are stored in `~/.aicli.json`. Protect the file permissions.
- **Sensitive stdin**: use `--no-send-stdin` to avoid sending stdin content to the LLM.
- **Risky command detection**: destructive operations require confirmation unless `--force` is used.
- **Log redaction**: logs should not contain full API keys or sensitive parameters.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Licensed under the [Apache License 2.0](LICENSE).

---

Note: this project is in early development; features and APIs may change.
