# HTTP API 数据爬取示例

这个示例展示了如何使用 jssandbox 的 HTTP 功能来爬取 API 接口数据。

## 功能说明

- 使用 `httpGet()` 进行简单的 GET 请求
- 使用 `httpRequest()` 进行自定义配置的 HTTP 请求（支持自定义 headers、timeout 等）
- 使用 `httpPost()` 发送 POST 请求
- 自动解析 JSON 响应
- 错误处理和状态码检查
- 数据预览和保存

## 使用方法

### 1. 编译运行

```bash
cd http_api
go run main.go
```

### 2. 指定API URL

```bash
go run main.go https://api.example.com/data
```

### 3. 输出

程序会：
- 在控制台显示4个示例的执行结果
- 将获取的数据保存到 `api_data_YYYYMMDD_HHMMSS.json` 文件
- 显示数据预览和统计信息

## 代码说明

### 示例1: 简单的GET请求

使用 `httpGet()` 方法进行最简单的 GET 请求：

```javascript
var result = httpGet("https://api.example.com/data");
```

### 示例2: 带自定义headers的GET请求

使用 `httpRequest()` 方法，可以自定义请求头：

```javascript
var result = httpRequest(url, {
    method: "GET",
    headers: {
        "User-Agent": "jssandbox-http-client/1.0",
        "Accept": "application/json"
    },
    timeout: 30
});
```

### 示例3: POST请求

使用 `httpPost()` 方法发送 POST 请求：

```javascript
var postData = {
    title: "测试标题",
    body: "这是测试内容",
    userId: 1
};
var result = httpPost(url, JSON.stringify(postData));
```

### 示例4: 完整配置的HTTP请求

使用 `httpRequest()` 进行完整配置，包括错误处理：

```javascript
var result = httpRequest(url, {
    method: "GET",
    headers: {
        "User-Agent": "jssandbox-http-client/1.0",
        "Accept": "application/json"
    },
    timeout: 30
});

// 检查错误
if (result.error) {
    console.error("请求失败:", result.error);
    return;
}

// 检查HTTP状态码
if (result.status < 200 || result.status >= 300) {
    console.warn("HTTP错误:", result.status);
    return;
}

// 解析JSON响应
var data = JSON.parse(result.body);
```

## HTTP API 方法说明

### httpGet(url)

简单的 GET 请求，返回响应对象。

**参数：**
- `url` (string): 请求的URL

**返回：**
```javascript
{
    status: 200,              // HTTP状态码
    statusText: "200 OK",     // 状态文本
    headers: {...},           // 响应头
    body: "...",              // 响应体（字符串）
    contentType: "application/json",  // Content-Type
    error: "..."              // 错误信息（如果有）
}
```

### httpPost(url, body)

发送 POST 请求，自动设置 `Content-Type: application/json`。

**参数：**
- `url` (string): 请求的URL
- `body` (string): 请求体（通常是 JSON 字符串）

**返回：** 同 `httpGet()`

### httpRequest(url, options)

通用的 HTTP 请求方法，支持完整配置。

**参数：**
- `url` (string): 请求的URL
- `options` (object): 请求选项
  - `method` (string): HTTP方法，如 "GET", "POST", "PUT", "DELETE" 等
  - `headers` (object): 请求头对象
  - `body` (string): 请求体
  - `timeout` (number): 超时时间（秒）

**返回：** 同 `httpGet()`

## 响应对象结构

所有 HTTP 方法都返回相同的响应对象结构：

```javascript
{
    status: 200,                    // HTTP状态码
    statusText: "200 OK",           // 状态文本
    headers: {                      // 响应头（对象）
        "Content-Type": "application/json",
        "Content-Length": "1234"
    },
    body: "...",                    // 响应体（字符串）
    contentType: "application/json", // Content-Type
    error: "..."                     // 错误信息（如果有错误）
}
```

## 错误处理

### 网络错误

```javascript
var result = httpGet("https://invalid-url.com");
if (result.error) {
    console.error("请求失败:", result.error);
}
```

### HTTP状态码错误

```javascript
var result = httpGet("https://api.example.com/data");
if (result.status < 200 || result.status >= 300) {
    console.error("HTTP错误:", result.status, result.statusText);
}
```

### JSON解析错误

```javascript
try {
    var data = JSON.parse(result.body);
} catch (e) {
    console.error("JSON解析失败:", e.message);
}
```

## 实际应用场景

### 1. 爬取REST API数据

```javascript
var result = httpGet("https://api.example.com/users");
var users = JSON.parse(result.body);
console.log("获取到", users.length, "个用户");
```

### 2. 带认证的API请求

```javascript
var result = httpRequest("https://api.example.com/protected", {
    method: "GET",
    headers: {
        "Authorization": "Bearer your-token-here",
        "Accept": "application/json"
    }
});
```

### 3. 发送表单数据

```javascript
var formData = "name=value&key=value";
var result = httpRequest("https://api.example.com/submit", {
    method: "POST",
    headers: {
        "Content-Type": "application/x-www-form-urlencoded"
    },
    body: formData
});
```

### 4. 分页获取数据

```javascript
var allData = [];
for (var page = 1; page <= 10; page++) {
    var result = httpGet("https://api.example.com/data?page=" + page);
    var pageData = JSON.parse(result.body);
    allData = allData.concat(pageData);
}
console.log("总共获取", allData.length, "条数据");
```

## 注意事项

1. **超时设置**：默认超时时间为30秒，可以通过 `options.timeout` 自定义
2. **JSON解析**：响应体是字符串，需要手动使用 `JSON.parse()` 解析
3. **错误处理**：始终检查 `result.error` 和 `result.status`
4. **请求头**：某些API可能需要特定的请求头（如 User-Agent、Accept等）
5. **HTTPS**：支持HTTPS请求，无需额外配置
6. **并发限制**：注意API的速率限制，避免请求过于频繁

## 配置说明

示例中禁用了浏览器功能以节省资源：

```go
config := jssandbox.DefaultConfig().
    WithHTTPTimeout(30 * time.Second).
    DisableBrowser() // HTTP爬取不需要浏览器
```

如果只需要HTTP功能，可以这样配置以提高性能。

## 依赖

- Go 1.21+
- jssandbox-go 库

## 测试API

示例中使用的测试API：
- https://jsonplaceholder.typicode.com/posts - 提供测试数据的REST API

你可以替换为任何公开的API进行测试。

