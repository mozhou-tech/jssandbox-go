# 日志管理规范

本项目使用 `logrus` 统一管理日志。

## 日志模块

所有日志功能通过 `jssandbox/log.go` 模块统一管理：

- `GetLogger()` - 获取默认的日志实例
- `SetLoggerLevel(level)` - 设置日志级别
- `SetLoggerFormatter(formatter)` - 设置日志格式化器
- `SetLoggerOutput(output)` - 设置日志输出
- `NewLogger()` - 创建新的日志实例
- `NewLoggerWithConfig(level, formatter, output)` - 使用指定配置创建新的日志实例

## 使用方式

### 在 Sandbox 中使用

Sandbox 结构体包含一个 `logger *logrus.Logger` 字段，可以通过以下方式使用：

```go
// 错误日志
sb.logger.WithError(err).Error("操作失败")

// 带字段的错误日志
sb.logger.WithError(err).WithField("path", filePath).Error("打开文件失败")

// 信息日志
sb.logger.Info("操作成功")

// 调试日志
sb.logger.Debug("调试信息")
```

### 在示例和命令行工具中使用

直接使用 logrus 的标准方法：

```go
import "github.com/sirupsen/logrus"

// 错误日志
logrus.WithError(err).Error("操作失败")

// 致命错误（会退出程序）
logrus.WithError(err).Fatal("致命错误")

// 带字段的日志
logrus.WithField("key", "value").Info("信息")
```

## 日志级别

- `Debug` - 调试信息
- `Info` - 一般信息
- `Warn` - 警告信息
- `Error` - 错误信息
- `Fatal` - 致命错误（会退出程序）

## 默认配置

默认日志配置：
- 格式：文本格式，包含完整时间戳
- 级别：Info
- 输出：标准输出
- 时间格式：2006-01-02 15:04:05

## 自定义配置

可以通过以下方式自定义日志配置：

```go
import (
    "github.com/sirupsen/logrus"
    "github.com/mozhou-tech/jssandbox-go/pkg/jssandbox"
)

// 设置日志级别
jssandbox.SetLoggerLevel(logrus.DebugLevel)

// 设置JSON格式
jssandbox.SetLoggerFormatter(&logrus.JSONFormatter{})

// 创建自定义logger
logger := jssandbox.NewLoggerWithConfig(
    logrus.DebugLevel,
    &logrus.JSONFormatter{},
    os.Stderr,
)

// 使用自定义logger创建沙盒
sandbox := jssandbox.NewSandboxWithLogger(ctx, logger)
```

