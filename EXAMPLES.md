# 使用示例

## 基础使用

### 使用默认配置

```go
package main

import (
    "context"
    "github.com/supacloud/jssandbox-go/jssandbox"
)

func main() {
    ctx := context.Background()
    sandbox := jssandbox.NewSandbox(ctx)
    defer sandbox.Close()
    
    result, err := sandbox.Run("1 + 1")
    // ...
}
```

## 配置管理

### 自定义配置

```go
config := jssandbox.DefaultConfig().
    WithTimeout(60 * time.Second).           // 设置默认超时
    WithHTTPTimeout(45 * time.Second).       // 设置 HTTP 超时
    WithBrowserTimeout(120 * time.Second).   // 设置浏览器超时
    WithMaxFileSize(50 * 1024 * 1024).       // 设置最大文件大小 50MB
    DisableBrowser()                         // 禁用浏览器功能

sandbox := jssandbox.NewSandboxWithConfig(ctx, config)
```

### 最小化配置（仅启用必要功能）

```go
config := jssandbox.DefaultConfig().
    DisableBrowser().
    DisableFileSystem().
    DisableDocuments().
    DisableImageProcessing()

sandbox := jssandbox.NewSandboxWithConfig(ctx, config)
// 只启用系统操作和 HTTP 功能
```

## 版本信息

### 获取版本信息

```go
// 获取版本号
version := jssandbox.GetVersion()
fmt.Println("版本:", version)

// 获取完整构建信息
buildInfo := jssandbox.GetBuildInfo()
fmt.Printf("版本: %s\n", buildInfo["version"])
fmt.Printf("构建时间: %s\n", buildInfo["buildTime"])
fmt.Printf("Git提交: %s\n", buildInfo["gitCommit"])
```

## 错误处理

### 处理超时错误

```go
result, err := sandbox.RunWithTimeout(code, 10*time.Second)
if err != nil {
    if sandboxErr, ok := err.(*jssandbox.SandboxError); ok {
        if sandboxErr.IsTimeout() {
            fmt.Println("执行超时")
        } else {
            fmt.Printf("错误代码: %s, 消息: %s\n", sandboxErr.Code, sandboxErr.Message)
        }
    }
}
```

### 检查错误类型

```go
result, err := sandbox.Run(code)
if err != nil {
    var sandboxErr *jssandbox.SandboxError
    if errors.As(err, &sandboxErr) {
        switch sandboxErr.Code {
        case jssandbox.ErrCodeTimeout:
            // 处理超时
        case jssandbox.ErrCodeFileNotFound:
            // 处理文件未找到
        case jssandbox.ErrCodeHTTPError:
            // 处理 HTTP 错误
        default:
            // 处理其他错误
        }
    }
}
```

## 使用自定义 Logger

```go
import "github.com/sirupsen/logrus"

logger := logrus.New()
logger.SetLevel(logrus.DebugLevel)
sandbox := jssandbox.NewSandboxWithLogger(ctx, logger)

// 或使用自定义配置
config := jssandbox.DefaultConfig().WithTimeout(30 * time.Second)
sandbox := jssandbox.NewSandboxWithLoggerAndConfig(ctx, logger, config)
```

## 构建时注入版本信息

使用 Makefile 构建时会自动注入版本信息：

```bash
make build
```

或手动构建：

```bash
VERSION=v1.0.0 \
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S') \
GIT_COMMIT=$(git rev-parse --short HEAD) \
go build -ldflags "-X github.com/supacloud/jssandbox-go/jssandbox.Version=$VERSION \
  -X github.com/supacloud/jssandbox-go/jssandbox.BuildTime=$BUILD_TIME \
  -X github.com/supacloud/jssandbox-go/jssandbox.GitCommit=$GIT_COMMIT" \
  -o bin/jssandbox ./cmd/jssandbox
```

## Makefile 命令

```bash
# 构建（自动注入版本信息）
make build

# 运行测试
make test

# 生成测试覆盖率报告
make test-coverage

# 格式化代码
make fmt

# 代码检查
make vet
make lint

# 清理构建文件
make clean

# 显示帮助
make help
```

