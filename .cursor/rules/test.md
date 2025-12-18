# Golang 单元测试最佳实践

## 目录
1. [测试文件命名和组织](#测试文件命名和组织)
2. [测试函数命名](#测试函数命名)
3. [表驱动测试](#表驱动测试)
4. [子测试](#子测试)
5. [测试辅助函数](#测试辅助函数)
6. [错误处理测试](#错误处理测试)
7. [资源清理](#资源清理)
8. [测试隔离](#测试隔离)
9. [测试数据管理](#测试数据管理)
10. [性能测试](#性能测试)
11. [测试覆盖率](#测试覆盖率)
12. [常见模式和反模式](#常见模式和反模式)

---

## 测试文件命名和组织

### 规则
- 测试文件必须以 `_test.go` 结尾
- 测试文件应该与被测试文件在同一包中
- 测试文件命名：`<被测试文件>_test.go`

### 示例
```
jssandbox/
  ├── sandbox.go
  ├── sandbox_test.go      ✅ 正确
  ├── browser.go
  └── browser_test.go      ✅ 正确
```

### 包名
- 单元测试：使用 `package jssandbox`（与被测试包相同）
- 集成测试：可以使用 `package jssandbox_test`（外部测试包）

---

## 测试函数命名

### 规则
- 测试函数必须以 `Test` 开头
- 函数名应该描述被测试的功能
- 使用驼峰命名法
- 格式：`Test<被测试函数名>_<场景描述>`

### 示例
```go
// ✅ 好的命名
func TestNewSandbox(t *testing.T)
func TestSandbox_Run(t *testing.T)
func TestSandbox_RunWithTimeout(t *testing.T)
func TestHTTPRequest_GET(t *testing.T)
func TestHTTPRequest_ErrorHandling(t *testing.T)

// ❌ 不好的命名
func Test1(t *testing.T)
func TestSandbox(t *testing.T)  // 太泛泛
func testRun(t *testing.T)     // 小写开头，不会被执行
```

---

## 表驱动测试

### 原则
表驱动测试是 Go 社区推荐的标准测试模式，特别适合测试多个输入/输出组合。

### 结构
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string      // 测试用例名称
        input   interface{} // 输入
        want    interface{} // 期望输出
        wantErr bool        // 是否期望错误
    }{
        {
            name:    "正常情况",
            input:    "valid input",
            want:     "expected output",
            wantErr:  false,
        },
        {
            name:    "边界情况",
            input:    "",
            want:     "",
            wantErr:  false,
        },
        {
            name:    "错误情况",
            input:    "invalid",
            want:     nil,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Function() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Function() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 实际示例
```go
func TestSandbox_Run(t *testing.T) {
    ctx := context.Background()
    sb := NewSandbox(ctx)
    defer sb.Close()

    tests := []struct {
        name    string
        code    string
        wantErr bool
    }{
        {
            name:    "简单计算",
            code:    "1 + 1",
            wantErr: false,
        },
        {
            name:    "变量赋值",
            code:    "var x = 10; x * 2",
            wantErr: false,
        },
        {
            name:    "无效代码",
            code:    "var x = ;",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := sb.Run(tt.code)
            if (err != nil) != tt.wantErr {
                t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && result == nil {
                t.Error("Run()返回nil结果")
            }
        })
    }
}
```

---

## 子测试

### 使用场景
- 将相关测试组织在一起
- 每个子测试可以独立运行
- 使用 `-run` 标志可以运行特定子测试

### 示例
```go
func TestHTTPRequest_ErrorHandling(t *testing.T) {
    ctx := context.Background()
    sb := NewSandbox(ctx)
    defer sb.Close()

    t.Run("缺少URL参数", func(t *testing.T) {
        result, err := sb.Run("httpRequest()")
        if err != nil {
            t.Fatalf("httpRequest() error = %v", err)
        }
        // 验证错误处理...
    })

    t.Run("无效URL", func(t *testing.T) {
        code := `var response = httpRequest("http://invalid-url"); response;`
        result, err := sb.Run(code)
        // 验证错误处理...
    })
}
```

### 运行特定子测试
```bash
# 运行所有 ErrorHandling 测试
go test -run TestHTTPRequest_ErrorHandling

# 运行特定子测试
go test -run TestHTTPRequest_ErrorHandling/缺少URL参数
```

---

## 测试辅助函数

### 原则
- 辅助函数应该调用 `t.Helper()` 标记为辅助函数
- 辅助函数失败时，错误信息会指向调用者，而不是辅助函数本身

### 示例
```go
func setupSandbox(t *testing.T) *Sandbox {
    t.Helper()
    ctx := context.Background()
    sb := NewSandbox(ctx)
    return sb
}

func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v, want %v", got, want)
    }
}

// 使用
func TestSomething(t *testing.T) {
    sb := setupSandbox(t)
    defer sb.Close()
    
    result := sb.Run("1 + 1")
    assertEqual(t, result.ToInteger(), int64(2))
}
```

---

## 错误处理测试

### 原则
- 不仅要测试成功情况，还要测试错误情况
- 验证错误类型和错误消息
- 使用 `wantErr` 字段明确表达期望

### 示例
```go
func TestFunction_ErrorHandling(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
        errMsg  string // 可选的错误消息验证
    }{
        {
            name:    "语法错误",
            input:   "var x =",
            wantErr: true,
        },
        {
            name:    "运行时错误",
            input:   "undefinedVariable.foo",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := sb.Run(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
            }
            if tt.wantErr && err == nil {
                t.Error("期望错误但未返回错误")
            }
        })
    }
}
```

---

## 资源清理

### 原则
- 使用 `defer` 确保资源被清理
- 对于需要清理的资源，在创建后立即使用 `defer`
- 使用 `t.Cleanup()` 进行测试级别的清理

### 示例
```go
func TestWithResources(t *testing.T) {
    // ✅ 正确：立即 defer
    sb := NewSandbox(ctx)
    defer sb.Close()

    // ✅ 正确：临时文件
    testFile := filepath.Join(t.TempDir(), "test.txt")
    defer os.Remove(testFile)  // 虽然 t.TempDir() 会自动清理，但显式清理更清晰

    // ✅ 正确：HTTP 测试服务器
    server := httptest.NewServer(handler)
    defer server.Close()

    // ✅ 使用 t.Cleanup (Go 1.14+)
    t.Cleanup(func() {
        sb.Close()
    })
}
```

### t.TempDir() 和 t.TempFile()
```go
func TestFileOperations(t *testing.T) {
    // 自动清理的临时目录
    testDir := t.TempDir()
    testFile := filepath.Join(testDir, "test.txt")
    
    // 测试结束后自动清理，无需手动 defer
}
```

---

## 测试隔离

### 原则
- 每个测试应该是独立的
- 测试之间不应该共享状态
- 使用 `t.Parallel()` 进行并行测试（注意：需要确保测试真正独立）

### 示例
```go
func TestIsolated1(t *testing.T) {
    t.Parallel()  // 可以并行运行
    // 测试代码...
}

func TestIsolated2(t *testing.T) {
    t.Parallel()  // 可以并行运行
    // 测试代码...
}

// ⚠️ 注意：如果测试共享全局状态，不要使用 t.Parallel()
```

### 环境变量控制
```go
func TestBrowserSession(t *testing.T) {
    if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
        t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
    }
    // 测试代码...
}
```

---

## 测试数据管理

### 原则
- 测试数据应该内联在测试中（小数据）
- 大文件可以使用 `testdata/` 目录
- 使用 `embed` 包嵌入测试数据（Go 1.16+）

### 示例
```go
// 小数据：内联
func TestSmallData(t *testing.T) {
    data := "Hello, World!"
    // 使用 data...
}

// 大数据：使用 testdata/
// testdata/large-file.txt
func TestLargeData(t *testing.T) {
    data, err := os.ReadFile("testdata/large-file.txt")
    if err != nil {
        t.Fatalf("读取测试数据失败: %v", err)
    }
    // 使用 data...
}

// 使用 embed (Go 1.16+)
//go:embed testdata/large-file.txt
var testData []byte

func TestEmbeddedData(t *testing.T) {
    // 直接使用 testData
}
```

---

## 性能测试

### 基准测试
```go
func BenchmarkFunction(b *testing.B) {
    ctx := context.Background()
    sb := NewSandbox(ctx)
    defer sb.Close()

    b.ResetTimer()  // 重置计时器，排除初始化时间
    for i := 0; i < b.N; i++ {
        sb.Run("1 + 1")
    }
}
```

### 运行基准测试
```bash
# 运行所有基准测试
go test -bench=.

# 运行特定基准测试
go test -bench=BenchmarkFunction

# 显示内存分配
go test -bench=. -benchmem
```

### 示例测试
```go
func ExampleFunction() {
    result := Function("input")
    fmt.Println(result)
    // Output: expected output
}
```

---

## 测试覆盖率

### 查看覆盖率
```bash
# 生成覆盖率文件
go test -coverprofile=coverage.out

# 查看覆盖率报告
go tool cover -html=coverage.out

# 查看覆盖率百分比
go test -cover
```

### 覆盖率目标
- 核心业务逻辑：> 80%
- 工具函数：> 90%
- 边界情况和错误处理：尽可能覆盖

---

## 常见模式和反模式

### ✅ 好的实践

1. **使用表驱动测试**
```go
tests := []struct {
    name string
    // ...
}{}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // ...
    })
}
```

2. **清晰的错误消息**
```go
if got != want {
    t.Errorf("Function() = %v, want %v", got, want)
}
```

3. **使用 t.Fatalf 提前退出**
```go
if err != nil {
    t.Fatalf("初始化失败: %v", err)
}
```

4. **测试边界条件**
```go
tests := []struct {
    name string
    input string
    wantErr bool
}{
    {"空字符串", "", false},
    {"超长字符串", strings.Repeat("a", 10000), false},
    {"特殊字符", "!@#$%", false},
}
```

### ❌ 反模式

1. **测试函数名不清晰**
```go
// ❌ 不好
func Test1(t *testing.T)
func Test(t *testing.T)

// ✅ 好
func TestNewSandbox(t *testing.T)
```

2. **测试之间共享状态**
```go
// ❌ 不好
var globalVar = "test"

func Test1(t *testing.T) {
    globalVar = "changed"
}

func Test2(t *testing.T) {
    // 可能受到 Test1 的影响
}
```

3. **忽略错误**
```go
// ❌ 不好
result, _ := sb.Run(code)

// ✅ 好
result, err := sb.Run(code)
if err != nil {
    t.Fatalf("Run() error = %v", err)
}
```

4. **测试逻辑过于复杂**
```go
// ❌ 不好：测试本身难以理解
func TestComplex(t *testing.T) {
    // 100+ 行的复杂测试逻辑
}

// ✅ 好：拆分成多个小测试
func TestComplex_Step1(t *testing.T) { /* ... */ }
func TestComplex_Step2(t *testing.T) { /* ... */ }
```

5. **硬编码路径**
```go
// ❌ 不好
testFile := "/tmp/test.txt"

// ✅ 好
testFile := filepath.Join(t.TempDir(), "test.txt")
```

---

## 测试工具和库

### 标准库
- `testing`: 核心测试包
- `testing/quick`: 属性测试
- `net/http/httptest`: HTTP 测试

### 第三方库
- `testify/assert`: 断言库
- `testify/mock`: Mock 框架
- `golang.org/x/tools/cmd/cover`: 覆盖率工具

### 使用 testify/assert 示例
```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestWithTestify(t *testing.T) {
    result, err := Function("input")
    
    // 使用 assert（失败后继续执行）
    assert.NoError(t, err)
    assert.Equal(t, "expected", result)
    
    // 使用 require（失败后立即停止）
    require.NoError(t, err)
    require.Equal(t, "expected", result)
}
```

---

## 总结

### 核心原则
1. **可读性**：测试应该像文档一样易读
2. **独立性**：每个测试应该独立运行
3. **可维护性**：测试应该易于修改和扩展
4. **完整性**：覆盖正常情况、边界情况和错误情况

### 检查清单
- [ ] 测试文件命名正确（`*_test.go`）
- [ ] 测试函数命名清晰（`TestFunction_Scenario`）
- [ ] 使用表驱动测试（多个测试用例）
- [ ] 使用子测试组织相关测试（`t.Run`）
- [ ] 正确清理资源（`defer` 或 `t.Cleanup`）
- [ ] 测试错误情况（`wantErr`）
- [ ] 使用 `t.Helper()` 标记辅助函数
- [ ] 错误消息清晰（包含 got/want）
- [ ] 测试独立（不共享状态）
- [ ] 使用临时目录（`t.TempDir()`）

---

## 参考资源

- [Go Testing Package](https://pkg.go.dev/testing)
- [Effective Go - Testing](https://go.dev/doc/effective_go#testing)
- [Go Blog - The cover story](https://go.dev/blog/cover)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

