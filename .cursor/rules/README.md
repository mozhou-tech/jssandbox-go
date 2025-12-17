# Cursor Rules for jssandbox-go

本项目是一个基于 Go 和 goja 的 JavaScript 沙盒环境，集成了浏览器操作、HTTP请求、文件系统、文档处理等功能。

## 代码风格

### Go 代码规范
- 遵循 Go 官方代码规范和最佳实践
- 使用 `gofmt` 格式化代码
- 使用 `golint` 和 `golangci-lint` 进行代码检查
- 所有导出的函数、类型、变量必须有注释
- 注释使用中文，格式：`// 函数名 功能描述`
- 函数名使用驼峰命名，首字母大写表示导出
- 私有函数和变量使用小写开头

### 代码组织
- 每个功能模块一个文件（如 `browser.go`, `filesystem.go`, `http.go`）
- 测试文件使用 `*_test.go` 命名
- 包级别函数按功能分组，使用空行分隔
- 导入包按标准库、第三方库、本地包分组，每组之间空行分隔

### 错误处理
- 所有可能失败的操作必须检查错误
- 错误信息使用中文，格式清晰
- 使用 `fmt.Errorf` 包装错误，添加上下文信息
- JavaScript 函数返回统一格式：`map[string]interface{}`，包含 `success` 和 `error` 字段
- 示例：
```go
return map[string]interface{}{
    "success": false,
    "error":   err.Error(),
}
```

### 日志使用
- 使用 `zap.Logger` 进行日志记录
- 日志级别：`Error` 用于错误，`Info` 用于重要操作，`Debug` 用于调试
- 日志包含结构化字段，使用 `zap.String()`, `zap.Error()` 等
- 示例：
```go
sb.logger.Error("操作失败", zap.String("path", filePath), zap.Error(err))
```

### Context 使用
- 所有长时间运行的操作必须支持 context 取消
- 使用 `context.Context` 传递取消信号和超时
- 浏览器操作、HTTP 请求等必须使用 context

## 项目特定规范

### 沙盒环境
- 所有 JavaScript 扩展函数通过 `sb.vm.Set()` 注册
- 函数注册在对应的 `register*()` 方法中完成
- 注册函数时考虑安全性，避免暴露危险操作
- JavaScript 函数参数验证在 Go 层完成

### 浏览器操作
- 使用 `chromedp` 进行浏览器自动化
- 浏览器会话使用 `BrowserSession` 结构体管理
- 必须注入反检测脚本（`injectStealthScript`）
- 浏览器上下文配置反检测选项
- 所有浏览器操作必须支持超时控制

### 文件系统操作
- 文件路径必须验证，防止路径遍历攻击
- 使用 `filepath` 包处理路径，确保跨平台兼容
- 大文件操作考虑使用流式处理
- 文件操作错误必须记录日志

### HTTP 请求
- 支持 GET、POST、PUT、DELETE 等方法
- 请求必须支持超时控制
- 响应数据统一格式返回
- 错误情况返回详细错误信息

### 文档处理
- 支持 Word、Excel、PPT、PDF 等格式
- 大文档支持分页读取
- 文档读取错误必须处理并返回友好错误信息

### 图片处理
- 使用 `disintegration/imaging` 库
- 图片操作支持格式转换、大小调整、裁剪等
- 图片处理错误必须处理

### 文件类型检测
- 优先使用二进制检测（`filetype` 库）
- 失败时回退到扩展名判断
- 检测结果统一格式返回

## 测试规范

### 测试文件
- 每个源文件对应一个测试文件
- 测试函数命名：`Test<函数名>` 或 `Test<类型>_<方法名>`
- 使用表驱动测试（table-driven tests）处理多个测试用例
- 测试用例使用中文描述

### 测试结构
```go
func TestFunctionName(t *testing.T) {
    ctx := context.Background()
    sb := NewSandbox(ctx)
    defer sb.Close()
    
    tests := []struct {
        name    string
        input   string
        want    interface{}
        wantErr bool
    }{
        {
            name:    "正常情况",
            input:   "test",
            want:    "expected",
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试代码
        })
    }
}
```

### 测试要求
- 测试必须覆盖正常情况和错误情况
- 边界条件必须测试
- 并发安全必须测试（如适用）
- 使用 `t.Fatal()` 或 `t.Fatalf()` 处理致命错误
- 使用 `t.Error()` 或 `t.Errorf()` 处理非致命错误

## 代码示例

### 注册 JavaScript 函数
```go
func (sb *Sandbox) registerExample() {
    sb.vm.Set("exampleFunc", func(arg string) map[string]interface{} {
        // 参数验证
        if arg == "" {
            return map[string]interface{}{
                "success": false,
                "error":   "参数不能为空",
            }
        }
        
        // 执行业务逻辑
        result, err := doSomething(arg)
        if err != nil {
            sb.logger.Error("操作失败", zap.String("arg", arg), zap.Error(err))
            return map[string]interface{}{
                "success": false,
                "error":   err.Error(),
            }
        }
        
        return map[string]interface{}{
            "success": true,
            "data":    result,
        }
    })
}
```

### 错误处理模式
```go
// 标准错误处理
if err != nil {
    sb.logger.Error("操作描述", zap.String("key", value), zap.Error(err))
    return map[string]interface{}{
        "success": false,
        "error":   err.Error(),
    }
}
```

### Context 和超时
```go
// 带超时的操作
ctx, cancel := context.WithTimeout(sb.ctx, timeout)
defer cancel()

// 在 goroutine 中执行，支持取消
done := make(chan error, 1)
go func() {
    // 执行操作
    done <- err
}()

select {
case <-ctx.Done():
    return nil, fmt.Errorf("操作超时: %v", timeout)
case err := <-done:
    return result, err
}
```

## 安全考虑

### 路径安全
- 所有文件路径必须验证，防止路径遍历（`../`）
- 使用 `filepath.Clean()` 清理路径
- 限制文件操作范围（如需要）

### 资源限制
- 大文件操作考虑内存限制
- 长时间运行的操作必须支持超时
- 并发操作考虑资源竞争

### 输入验证
- JavaScript 函数参数必须验证
- 防止注入攻击
- 文件类型验证

## 性能优化

### 资源管理
- 及时关闭文件、网络连接等资源
- 使用 `defer` 确保资源释放
- 浏览器会话及时关闭

### 并发处理
- 需要并发时使用 goroutine
- 使用 `sync.Mutex` 保护共享资源
- 考虑使用 channel 进行通信

## 依赖管理

### 添加依赖
- 使用 `go get` 添加依赖
- 更新 `go.mod` 和 `go.sum`
- 重要依赖需要文档说明

### 当前主要依赖
- `github.com/dop251/goja` - JavaScript 运行时
- `github.com/chromedp/chromedp` - 浏览器自动化
- `go.uber.org/zap` - 日志库
- `github.com/disintegration/imaging` - 图片处理
- `github.com/h2non/filetype` - 文件类型检测

## 提交规范

### Commit 消息
- 使用中文描述提交内容
- 格式：`<类型>: <描述>`
- 类型：`feat`（新功能）、`fix`（修复）、`docs`（文档）、`test`（测试）、`refactor`（重构）

### 代码审查
- 所有代码变更必须通过测试
- 运行 `go test ./...` 确保测试通过
- 运行 `go build ./...` 确保编译通过
- 检查代码格式和 lint 问题

## 文档要求

### 代码注释
- 所有导出的函数、类型、变量必须有注释
- 注释说明功能、参数、返回值
- 复杂逻辑必须有注释说明

### README 更新
- 新功能添加到 README.md
- 更新功能状态（✅/❌/🚧）
- 重要变更更新文档

## 注意事项

1. **不要**在 JavaScript 函数中直接 panic，必须返回错误
2. **不要**忽略错误，所有错误必须处理
3. **不要**硬编码路径，使用相对路径或配置
4. **必须**在所有文件操作后关闭文件
5. **必须**在所有浏览器操作后清理资源
6. **必须**在测试中使用 `defer sb.Close()` 清理资源
7. **必须**验证所有用户输入
8. **必须**记录重要操作的日志

