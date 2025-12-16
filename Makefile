.PHONY: build test

build:
	@echo "构建 jssandbox..."
	@go build -o bin/jssandbox ./cmd/jssandbox

test:
	@echo "运行测试..."
	@go test -v ./...
