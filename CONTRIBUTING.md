# Contributing

Thanks for your interest in AICLI! We welcome contributions of all kinds—code, documentation, bug reports, and feature requests.

English | [中文](CONTRIBUTING_CN.md)

## Table of contents

- [Code of Conduct](#code-of-conduct)
- [How to contribute](#how-to-contribute)
  - [Report bugs](#report-bugs)
  - [Request features](#request-features)
  - [Submit code changes](#submit-code-changes)
- [Development workflow](#development-workflow)
- [Coding guidelines](#coding-guidelines)
- [Commit message convention](#commit-message-convention)
- [Testing requirements](#testing-requirements)
- [Documentation requirements](#documentation-requirements)
- [Pull request checklist](#pull-request-checklist)

## Code of Conduct

We are committed to providing a friendly, safe, and welcoming environment for everyone.

Examples of positive behavior:
- Use welcoming and inclusive language
- Respect differing viewpoints and experiences
- Accept constructive criticism gracefully
- Focus on what is best for the community

Unacceptable behavior:
- Harassment, threats, or discriminatory language
- Sexualized language or imagery
- Publishing others’ private information without explicit permission

## How to contribute

### Report bugs

Before filing a new issue, please search existing issues to avoid duplicates.

When reporting a bug, include:
- A clear title
- Steps to reproduce
- Expected behavior vs. actual behavior
- Logs / screenshots if applicable
- Environment info (OS, Go version, AICLI version, LLM provider)

### Request features

When requesting a feature, include:
- What you want to achieve
- Why it’s useful (use cases)
- Any ideas on implementation (optional)
- Alternatives you considered

### Submit code changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Run tests and linters (see below)
5. Commit using the convention below
6. Push to your fork and open a Pull Request

## Development workflow

```bash
# Clone your fork
git clone https://github.com/your-username/aicli.git
cd aicli

# Add upstream remote
git remote add upstream https://github.com/studyzy/aicli.git

# Sync main branch
git checkout main
git pull upstream main

# Create a feature branch
git checkout -b feature/my-new-feature
```

Recommended branch naming:
- `feature/feature-name` for new features
- `fix/bug-description` for bug fixes
- `docs/doc-update` for documentation updates
- `refactor/refactor-description` for refactors
- `test/test-description` for tests

## Coding guidelines

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Prefer the standard library where possible
- Keep code simple, clear, and testable
- Wrap errors with context (use `%w`)

## Commit message convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/).

Format:

```text
<type>(<scope>): <subject>

<body>

<footer>
```

Common types:
- `feat`: a new feature
- `fix`: a bug fix
- `docs`: documentation changes
- `refactor`: code refactoring
- `test`: test changes
- `chore`: tooling/build changes
- `ci`: CI changes

## Testing requirements

Minimum expectations:
- New code should include tests
- Target coverage depends on the package; keep business-critical logic well tested

Run locally:

```bash
make test
make lint

# or
go test ./...
```

## Documentation requirements

When adding user-visible changes:
- Update `README.md` / `README_CN.md` as needed
- Update relevant docs under `docs/`
- Add usage examples when appropriate

## Pull request checklist

Before opening a PR:
- [ ] Code follows project style
- [ ] Tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Linter passes (`make lint` or `golangci-lint run`)
- [ ] Documentation updated
- [ ] PR description clearly explains the change

## Need help?

- Read the development guide: `docs/development.md`
- Search existing issues: https://github.com/studyzy/aicli/issues
- Ask in Discussions: https://github.com/studyzy/aicli/discussions
