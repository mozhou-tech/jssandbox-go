.PHONY: build test test-coverage clean install lint fmt vet help

# 构建变量
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标志
LDFLAGS = -X github.com/supacloud/jssandbox-go/jssandbox.Version=$(VERSION) \
          -X github.com/supacloud/jssandbox-go/jssandbox.BuildTime=$(BUILD_TIME) \
          -X github.com/supacloud/jssandbox-go/jssandbox.GitCommit=$(GIT_COMMIT)

# 默认目标
.DEFAULT_GOAL := help

help: ## 显示帮助信息
	@echo "可用的 make 目标:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## 构建 jssandbox 可执行文件
	@echo "构建 jssandbox (版本: $(VERSION))..."
	@mkdir -p bin
	@go build -ldflags "$(LDFLAGS)" -o bin/jssandbox ./cmd/jssandbox
	@echo "构建完成: bin/jssandbox"

test: ## 运行测试
	@echo "运行测试..."
	@go test -v ./...

test-coverage: ## 生成测试覆盖率报告
	@echo "生成测试覆盖率报告..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

clean: ## 清理构建文件
	@echo "清理构建文件..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "清理完成"

install: ## 安装到 $GOPATH/bin
	@echo "安装 jssandbox..."
	@go install -ldflags "$(LDFLAGS)" ./cmd/jssandbox
	@echo "安装完成"

lint: ## 运行代码检查（需要安装 golangci-lint）
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "警告: golangci-lint 未安装，跳过检查"; \
		echo "安装方法: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt: ## 格式化代码
	@echo "格式化代码..."
	@go fmt ./...
	@echo "格式化完成"

vet: ## 运行 go vet
	@echo "运行 go vet..."
	@go vet ./...
	@echo "检查完成"
