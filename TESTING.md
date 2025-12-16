# 测试说明

本项目为所有功能模块提供了完整的测试用例。

## 测试文件结构

```
jssandbox/
├── sandbox_test.go      # 核心沙盒功能测试
├── system_test.go       # 系统操作测试
├── http_test.go         # HTTP请求测试
├── filesystem_test.go   # 文件系统操作测试
├── browser_test.go      # 浏览器操作测试
└── documents_test.go    # 文档读取测试
```

## 运行测试

### 运行所有测试
```bash
go test ./jssandbox -v
```

### 运行特定测试
```bash
# 运行核心功能测试
go test ./jssandbox -v -run TestSandbox

# 运行系统操作测试
go test ./jssandbox -v -run TestGetCurrent

# 运行HTTP测试
go test ./jssandbox -v -run TestHTTP

# 运行文件系统测试
go test ./jssandbox -v -run TestWriteFile

# 运行浏览器测试（需要Chrome环境）
go test ./jssandbox -v -run TestBrowser

# 运行文档测试
go test ./jssandbox -v -run TestReadWord
```

### 跳过需要外部依赖的测试
```bash
# 跳过浏览器测试
SKIP_BROWSER_TESTS=true go test ./jssandbox -v
```

## 测试覆盖范围

### 1. 核心沙盒功能 (sandbox_test.go)
- ✅ NewSandbox 创建沙盒
- ✅ NewSandboxWithLogger 使用自定义logger
- ✅ Run 执行JavaScript代码
- ✅ RunWithTimeout 超时控制
- ✅ Set/Get 变量设置和获取
- ✅ Close 资源清理
- ✅ 扩展功能注册验证
- ✅ 复杂脚本执行
- ✅ 错误处理

### 2. 系统操作 (system_test.go)
- ✅ getCurrentTime 获取当前时间
- ✅ getCurrentDate 获取当前日期
- ✅ getCurrentDateTime 获取日期时间
- ✅ getCPUNum 获取CPU数量
- ✅ getMemorySize 获取内存信息
- ✅ getDiskSize 获取磁盘空间
- ✅ sleep 休眠功能
- ✅ 系统函数集成测试

### 3. HTTP请求 (http_test.go)
- ✅ httpRequest GET请求
- ✅ httpRequest POST请求
- ✅ httpRequest 自定义请求头
- ✅ httpRequest 超时控制
- ✅ httpGet 便捷方法
- ✅ httpPost 便捷方法
- ✅ 错误处理（无效URL、缺少参数）
- ✅ 响应头解析

### 4. 文件系统操作 (filesystem_test.go)
- ✅ writeFile 写入文件
- ✅ readFile 读取文件
- ✅ readFile 分页读取（offset/limit）
- ✅ readFileHead 读取前几行
- ✅ readFileTail 读取后几行
- ✅ getFileInfo 获取文件信息
- ✅ getFileHash 计算文件哈希（MD5/SHA1/SHA256/SHA512）
- ✅ renameFile 重命名文件
- ✅ appendFile 追加文件
- ✅ readImageBase64 读取图片base64
- ✅ 错误处理（文件不存在等）

### 5. 浏览器操作 (browser_test.go)
- ✅ browserNavigate 页面导航
- ✅ browserScreenshot 截图
- ✅ browserEvaluate 执行脚本
- ✅ browserClick 点击元素
- ✅ browserFill 填充表单
- ✅ 错误处理（缺少参数）

**注意**: 浏览器测试需要Chrome/Chromium环境，如果没有安装可以设置 `SKIP_BROWSER_TESTS=true` 跳过。

### 6. 文档读取 (documents_test.go)
- ✅ readWord 读取Word文档
- ✅ readWord 分页选项
- ✅ readExcel 读取Excel文件
- ✅ readExcel 分页和sheet选项
- ✅ readPPT 读取PPT文件
- ✅ readPPT 分页选项
- ✅ readPDF 读取PDF文件
- ✅ readPDF 分页选项
- ✅ 错误处理（缺少参数、无效文件）

**注意**: 文档测试需要实际的文档文件，如果没有测试文件会返回错误（这是预期的）。

## 测试覆盖率

运行测试覆盖率分析：
```bash
go test ./jssandbox -cover
```

生成详细覆盖率报告：
```bash
go test ./jssandbox -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 测试最佳实践

1. **单元测试**: 每个功能都有对应的单元测试
2. **集成测试**: 测试函数之间的协作
3. **错误处理**: 验证错误情况的处理
4. **边界条件**: 测试边界情况和异常输入
5. **临时文件**: 使用 `t.TempDir()` 创建临时文件，测试后自动清理
6. **Mock服务器**: HTTP测试使用 `httptest.NewServer` 创建测试服务器

## 注意事项

1. **浏览器测试**: 需要Chrome/Chromium环境，可以通过环境变量跳过
2. **文档测试**: 需要实际的文档文件，测试会处理文件不存在的情况
3. **网络测试**: HTTP测试使用本地测试服务器，不依赖外部网络
4. **文件系统测试**: 使用临时目录，测试后自动清理

## 持续集成

测试可以在CI/CD环境中运行：
```yaml
# 示例 GitHub Actions
- name: Run tests
  run: go test ./jssandbox -v -coverprofile=coverage.out
  
- name: Skip browser tests
  run: SKIP_BROWSER_TESTS=true go test ./jssandbox -v
```

