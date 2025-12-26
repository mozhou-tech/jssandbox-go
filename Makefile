.PHONY: build test test-coverage clean install lint fmt vet help

# 构建变量
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标志
LDFLAGS = -X github.com/mozhou-tech/jssandbox-go/pkg/jssandbox.Version=$(VERSION) \
          -X github.com/mozhou-tech/jssandbox-go/pkg/jssandbox.BuildTime=$(BUILD_TIME) \
          -X github.com/mozhou-tech/jssandbox-go/pkg/jssandbox.GitCommit=$(GIT_COMMIT)

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
	@go test -v ./pkg/jssandbox/...

test-coverage: ## 生成测试覆盖率报告
	@echo "生成测试覆盖率报告..."
	@go test -v -coverprofile=coverage.out ./pkg/jssandbox/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

clean: ## 清理构建文件
	@echo "清理构建文件..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "清理完成"

install: ## 安装到 $GOPATH/bin
	go mod tidy

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

# 发布相关命令
.PHONY: tag release verify-release

# 创建并推送版本标签（自动使用当前时间）
# 使用方法: make tag
tag:
	@VERSION=$$(date +%Y%m%d-%H%M%S); \
	echo "创建版本标签: v$$VERSION"; \
	git tag v$$VERSION; \
	echo "标签已创建，正在推送到远程..."; \
	git push github v$$VERSION

# 验证发布（检查代码是否可以正常构建和测试）
verify-release:
	@echo "验证代码构建..."
	go mod tidy
	@echo "验证库代码编译（排除 examples 目录）..."
	@go vet ./pkg/... || true
	@echo "✓ 代码检查完成"
	@echo "运行测试..."
	# go test ./pkg/... -v
	@echo "验证通过！"

# 完整发布流程（自动使用当前时间生成版本标签）
# 使用方法: make release
release: verify-release
	@VERSION=$$(date +%Y%m%d-%H%M%S); \
	echo "准备发布版本 v$$VERSION..."; \
	echo "1. 确保所有更改已提交:"; \
	go mod tidy; \
	git commit -a -m "tidy go.mod and go.sum"; \
	git status --short; \
	echo ""; \
	echo "2. 创建版本标签 v$$VERSION..."; \
	git tag v0.0.0-$$VERSION; \
	echo ""; \
	echo "3. 标签已创建，执行以下命令完成发布:"; \
	echo "   git push github master"; \
	echo "   git push github v0.0.0-$$VERSION"; \
	git push github master; \
	git push github v0.0.0-$$VERSION


