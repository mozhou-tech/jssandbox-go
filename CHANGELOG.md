# 更新日志

## [未发布] - 2024-XX-XX

### 新增功能

#### 配置管理
- ✅ 添加 `Config` 结构体，支持统一配置管理
- ✅ 支持配置超时时间（默认超时、HTTP超时、浏览器超时）
- ✅ 支持配置最大文件大小和允许的文件类型
- ✅ 支持选择性启用/禁用功能模块（浏览器、文件系统、HTTP等）
- ✅ 提供链式配置方法（`WithTimeout`, `WithHTTPTimeout` 等）

#### 版本管理
- ✅ 添加版本信息支持（`Version`, `BuildTime`, `GitCommit`）
- ✅ 支持通过构建标志注入版本信息
- ✅ 提供 `GetVersion()` 和 `GetBuildInfo()` 方法

#### 错误处理
- ✅ 统一错误处理机制，定义 `SandboxError` 类型
- ✅ 支持错误代码分类（`ErrCodeTimeout`, `ErrCodeFileNotFound` 等）
- ✅ 支持错误链（`Unwrap` 方法）
- ✅ 提供便捷的错误创建方法

#### 接口抽象
- ✅ 定义 `SandboxRunner` 接口，提升可测试性和可扩展性
- ✅ `Sandbox` 结构体实现 `SandboxRunner` 接口

### 改进

#### 沙盒核心
- ✅ `NewSandbox` 现在使用默认配置
- ✅ 新增 `NewSandboxWithConfig` 支持自定义配置
- ✅ 新增 `NewSandboxWithLoggerAndConfig` 支持自定义 logger 和配置
- ✅ `RunWithTimeout` 现在支持使用配置中的默认超时时间
- ✅ `registerExtensions` 根据配置选择性注册功能模块

#### HTTP 模块
- ✅ HTTP 请求现在使用配置中的默认超时时间

#### 构建系统
- ✅ 改进 Makefile，支持版本注入
- ✅ 添加测试覆盖率支持（`make test-coverage`）
- ✅ 添加代码格式化（`make fmt`）
- ✅ 添加代码检查（`make vet`, `make lint`）
- ✅ 添加帮助信息（`make help`）

### 文档

- ✅ 添加 `STRUCTURE_EVALUATION.md` 项目结构评估报告
- ✅ 添加 `CHANGELOG.md` 更新日志

### 向后兼容性

- ✅ 所有现有 API 保持向后兼容
- ✅ `NewSandbox` 和 `NewSandboxWithLogger` 行为不变（使用默认配置）
- ✅ 现有代码无需修改即可使用新功能

### 使用示例

#### 使用自定义配置

```go
config := jssandbox.DefaultConfig().
    WithTimeout(60 * time.Second).
    WithHTTPTimeout(45 * time.Second).
    DisableBrowser()

sandbox := jssandbox.NewSandboxWithConfig(ctx, config)
defer sandbox.Close()
```

#### 获取版本信息

```go
version := jssandbox.GetVersion()
buildInfo := jssandbox.GetBuildInfo()
fmt.Printf("版本: %s, 构建时间: %s\n", buildInfo["version"], buildInfo["buildTime"])
```

#### 错误处理

```go
result, err := sandbox.RunWithTimeout(code, 10*time.Second)
if err != nil {
    if sandboxErr, ok := err.(*jssandbox.SandboxError); ok {
        if sandboxErr.IsTimeout() {
            // 处理超时错误
        }
    }
}
```

