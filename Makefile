# Makefile for aicli project

.PHONY: all build test coverage clean install help

# 变量定义
BINARY_NAME=aicli
BUILD_DIR=./bin
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# 默认目标
all: test build

# 构建二进制文件
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/aicli
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行所有测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 运行测试并生成覆盖率报告
coverage:
	@echo "Running tests with coverage..."
	@echo "Running all tests..."
	@go test -v ./...
	@echo "\nCalculating coverage for core libraries (internal/... pkg/...)..."
	@go test -coverprofile=$(COVERAGE_FILE) ./internal/... ./pkg/...
	@go tool cover -func=$(COVERAGE_FILE)
	@echo "\nGenerate HTML coverage report..."
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# 检查测试覆盖率是否达标 (≥60%)
coverage-check:
	@echo "Checking coverage threshold (≥60%)..."
	@echo "Running all tests..."
	@go test ./... > /dev/null 2>&1
	@echo "Calculating coverage for core libraries..."
	@go test -coverprofile=$(COVERAGE_FILE) ./internal/... ./pkg/... > /dev/null 2>&1
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$COVERAGE >= 60" | bc) -eq 1 ]; then \
		echo "✓ Coverage $$COVERAGE% meets threshold (≥60%)"; \
	else \
		echo "✗ Coverage $$COVERAGE% below threshold (≥60%)"; \
		exit 1; \
	fi

# 代码格式化
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 代码静态检查
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

# 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@go clean
	@echo "Clean complete"

# 安装到系统路径
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $$(go env GOPATH)/bin/
	@echo "Installed to $$(go env GOPATH)/bin/$(BINARY_NAME)"

# 更新依赖
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  all            - Run tests and build (default)"
	@echo "  build          - Build the binary"
	@echo "  test           - Run all tests"
	@echo "  coverage       - Run tests with coverage report"
	@echo "  coverage-check - Check if coverage meets threshold (≥60%)"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run static code analysis"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Install binary to GOPATH/bin"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  help           - Show this help message"
