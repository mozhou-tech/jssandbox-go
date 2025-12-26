# JavaScript沙盒 API 参考文档（大模型专用）

本文档为AI大模型提供完整的JavaScript沙盒API参考，包含所有可用函数、参数说明、返回值结构和实际使用示例。

## 目录

1. [快速开始](#快速开始)
2. [系统操作](#系统操作)
3. [HTTP请求](#http请求)
4. [文件系统操作](#文件系统操作)
5. [文档处理 (PDF)](#文档处理-pdf)
6. [文档读取](#文档读取)
7. [浏览器自动化](#浏览器自动化)
8. [图片处理](#图片处理)
9. [文件类型检测](#文件类型检测)
10. [加密/解密](#加密解密)
11. [压缩/解压缩](#压缩解压缩)
12. [CSV处理](#csv处理)
13. [环境变量和配置](#环境变量和配置)
14. [数据验证](#数据验证)
15. [日期时间增强](#日期时间增强)
16. [编码/解码](#编码解码)
17. [进程管理](#进程管理)
18. [网络工具](#网络工具)
19. [路径处理](#路径处理)
20. [文本操作](#文本操作)
21. [日志功能](#日志功能)
22. [错误处理](#错误处理)

---

## 快速开始

### 重要：语法限制

**必须注意以下限制，否则会导致代码执行失败：**

1.  **不支持 `async/await`**：沙盒环境目前不支持异步语法。请使用提供的同步函数。
2.  **不支持 top-level `await`**：代码必须是纯同步执行的。
3.  **执行隔离与返回值**：作为 Eino 等工具调用时，代码会自动包装在 `(function(){ ... })()` 匿名函数中。这意味着您必须使用 `return` 语句来返回您想要获取的结果，同时您可以在不同次调用中使用相同的 `const` 或 `let` 变量名而不会冲突。
4.  **函数同步返回**：沙盒中看似异步的操作（如 `httpGet`, `session.navigate`）实际上是同步返回结果的，无需 `await`。

### 基本用法

```javascript
// 所有函数都在全局作用域中可用，无需导入
var date = getCurrentDate();
var time = getCurrentTime();
console.log("日期:", date, "时间:", time);

// 使用 logger 进行日志记录
logger.info("应用程序启动");
logger.debug("调试信息");
```

### 返回值处理

在 Eino 工具调用中，您**必须**使用 `return` 语句返回最终结果。如果未显式使用 `return`，工具将返回 `undefined`。

```javascript
// 正确：使用 return 返回结果
var response = httpGet("https://api.ipify.org");
return response.body;

// 错误：没有 return，即使有最后一条表达式，结果也是 undefined
var response = httpGet("https://api.ipify.org");
response.body; 
```

大多数函数返回对象，包含 `success` 字段表示操作是否成功：

```javascript
var result = writeFile("test.txt", "Hello");
if (result.success) {
    return "写入成功";
} else {
    return "写入失败: " + result.error;
}
```

---

## 系统操作

### getCurrentTime()

获取当前时间（格式：HH:mm:ss）

**返回值**: `string` - 当前时间字符串

**示例**:
```javascript
var time = getCurrentTime();
console.log(time); // "14:30:25"
```

### getCurrentDate()

获取当前日期（格式：YYYY-MM-DD）

**返回值**: `string` - 当前日期字符串

**示例**:
```javascript
var date = getCurrentDate();
console.log(date); // "2024-01-15"
```

### getCurrentDateTime()

获取当前日期时间（格式：YYYY-MM-DD HH:mm:ss）

**返回值**: `string` - 当前日期时间字符串

**示例**:
```javascript
var datetime = getCurrentDateTime();
console.log(datetime); // "2024-01-15 14:30:25"
```

### getCPUNum()

获取CPU核心数量

**返回值**: `number` - CPU核心数

**示例**:
```javascript
var cpuNum = getCPUNum();
console.log("CPU核心数:", cpuNum); // 8
```

### getMemorySize()

获取内存信息

**返回值**: `object`
- `total` (number): 总内存（字节）
- `available` (number): 可用内存（字节）
- `used` (number): 已用内存（字节）
- `totalStr` (string): 总内存（人类可读格式，如 "8.0 GB"）
- `availableStr` (string): 可用内存（人类可读格式）
- `usedStr` (string): 已用内存（人类可读格式）

**示例**:
```javascript
var mem = getMemorySize();
console.log("总内存:", mem.totalStr);
console.log("已用内存:", mem.usedStr);
console.log("可用内存:", mem.availableStr);
```

### getDiskSize(path?)

获取磁盘空间信息

**参数**:
- `path` (string, 可选): 磁盘路径，默认为 "/"

**返回值**: `object`
- `total` (number): 总空间（字节）
- `free` (number): 空闲空间（字节）
- `used` (number): 已用空间（字节）
- `totalStr` (string): 总空间（人类可读格式）
- `freeStr` (string): 空闲空间（人类可读格式）
- `usedStr` (string): 已用空间（人类可读格式）
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var disk = getDiskSize("/");
console.log("总空间:", disk.totalStr);
console.log("已用空间:", disk.usedStr);
console.log("空闲空间:", disk.freeStr);
```

### sleep(ms)

休眠指定毫秒数

**参数**:
- `ms` (number): 休眠毫秒数

**返回值**: `undefined`

**示例**:
```javascript
console.log("开始");
sleep(1000); // 休眠1秒
console.log("结束");
```

---

## HTTP请求

### httpRequest(url, options?)

通用HTTP请求

**参数**:
- `url` (string): 请求URL
- `options` (object, 可选): 请求选项
  - `method` (string): HTTP方法，默认 "GET"
  - `headers` (object): 请求头
  - `body` (string): 请求体
  - `timeout` (number): 超时时间（秒），默认30

**返回值**: `object`
- `status` (number): HTTP状态码
- `statusText` (string): 状态文本
- `body` (string): 响应体
- `headers` (object): 响应头
- `error` (string, 可选): 如果请求发生网络错误（如连接超时、拒绝连接），此字段将包含错误描述。建议优先检查此字段。

**示例**:
```javascript
// 带错误检查的 GET 请求
var response = httpRequest("https://api.example.com/data");
if (response.error) {
    console.error("请求失败:", response.error);
} else {
    console.log("状态码:", response.status);
    console.log("响应体:", response.body);
}
```

### httpGet(url)

GET请求（简化版）

**参数**:
- `url` (string): 请求URL

**返回值**: 同 `httpRequest`

**示例**:
```javascript
var response = httpGet("https://api.example.com/data");
console.log("状态码:", response.status);
console.log("响应体:", response.body);
```

### httpPost(url, body?)

POST请求（简化版）

**参数**:
- `url` (string): 请求URL
- `body` (string, 可选): 请求体

**返回值**: 同 `httpRequest`

**示例**:
```javascript
var response = httpPost("https://api.example.com/data", "test data");
console.log("状态码:", response.status);
```

### fetch(url, options?) (Polyfill)

为了兼容大模型的习惯，提供了一个**同步版**的 `fetch`。

**注意**: 这个 `fetch` 是同步的，**严禁**使用 `await fetch(...)`。

**示例**:
```javascript
var response = fetch("https://api.ipify.org");
console.log("IP:", response.text());
```

---

## 文件系统操作

### writeFile(path, content)

写入文件

**参数**:
- `path` (string): 文件路径
- `content` (string): 文件内容

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var result = writeFile("test.txt", "Hello World\n第二行");
if (result.success) {
    console.log("写入成功");
} else {
    console.error("写入失败:", result.error);
}
```

### appendFile(path, content)

追加内容到文件

**参数**:
- `path` (string): 文件路径
- `content` (string): 要追加的内容

**返回值**: 同 `writeFile`

**示例**:
```javascript
appendFile("test.txt", "\n追加的内容");
```

### readFile(path, options?)

读取文件内容

**参数**:
- `path` (string): 文件路径
- `options` (object, 可选): 读取选项
  - `page` (number): 页码（从1开始），用于分页读取
  - `pageSize` (number): 每页大小（字节）

**返回值**: `object`
- `data` (string): 文件内容
- `length` (number): 内容长度（字节）
- `totalSize` (number): 文件总大小（字节）
- `page` (number): 当前页码
- `totalPages` (number): 总页数
- `error` (string, 可选): 错误信息

**示例**:
```javascript
// 读取整个文件
var file = readFile("test.txt");
console.log("内容:", file.data);
console.log("大小:", file.length, "字节");

// 分页读取
var page1 = readFile("large.txt", {page: 1, pageSize: 1024});
console.log("第1页:", page1.data);
console.log("总页数:", page1.totalPages);
```

### readFileHead(path, lines)

读取文件前几行

**参数**:
- `path` (string): 文件路径
- `lines` (number): 行数

**返回值**: `object`
- `data` (string): 文件内容
- `lines` (number): 实际读取的行数
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var head = readFileHead("log.txt", 10);
console.log("前10行:", head.data);
```

### readFileTail(path, lines)

读取文件后几行

**参数**:
- `path` (string): 文件路径
- `lines` (number): 行数

**返回值**: 同 `readFileHead`

**示例**:
```javascript
var tail = readFileTail("log.txt", 10);
console.log("后10行:", tail.data);
```

### getFileInfo(path)

获取文件元信息

**参数**:
- `path` (string): 文件路径

**返回值**: `object`
- `name` (string): 文件名
- `size` (number): 文件大小（字节）
- `modTime` (string): 修改时间
- `isDir` (boolean): 是否为目录
- `mode` (string): 文件权限
- `mime` (string): MIME类型
- `extension` (string): 文件扩展名
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var info = getFileInfo("test.txt");
console.log("文件名:", info.name);
console.log("大小:", info.size, "字节");
console.log("修改时间:", info.modTime);
console.log("MIME类型:", info.mime);
```

### renameFile(oldPath, newPath)

重命名文件

**参数**:
- `oldPath` (string): 原文件路径
- `newPath` (string): 新文件路径

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var result = renameFile("old.txt", "new.txt");
if (result.success) {
    console.log("重命名成功");
}
```

### getFileHash(path, type)

获取文件哈希值

**参数**:
- `path` (string): 文件路径
- `type` (string): 哈希类型，可选值: "md5", "sha1", "sha256", "sha512"

**返回值**: `object`
- `hash` (string): 哈希值（十六进制）
- `type` (string): 哈希类型
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var hash = getFileHash("test.txt", "md5");
console.log("MD5:", hash.hash);

var sha256 = getFileHash("test.txt", "sha256");
console.log("SHA256:", sha256.hash);
```

### readImageBase64(path)

读取图片的base64编码

**参数**:
- `path` (string): 图片文件路径

**返回值**: `object`
- `data` (string): base64编码的图片数据（包含data URI前缀）
- `mime` (string): MIME类型
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var image = readImageBase64("photo.jpg");
console.log("Base64数据:", image.data);
// 输出类似: "data:image/jpeg;base64,/9j/4AAQSkZJRg..."
```

### openFile(path)

使用系统默认程序打开文件

**参数**:
- `path` (string): 文件路径

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

**示例**:
```javascript
openFile("document.pdf"); // 使用默认PDF阅读器打开
```

---

## 文档处理 (PDF)

基于 `pdfcpu` 实现的 PDF 处理功能。

### pdfGetPageCount(path)

获取 PDF 文档的总页数。

**参数**:
- `path` (string): PDF 文件路径

**返回值**: `object`
- `success` (boolean): 是否成功
- `pages` (number): 总页数
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var result = pdfGetPageCount("document.pdf");
if (result.success) {
    console.log("总页数:", result.pages);
}
```

### pdfMerge(inFiles, outFile)

合并多个 PDF 文件为一个。

**参数**:
- `inFiles` (array): 输入 PDF 文件路径数组
- `outFile` (string): 输出 PDF 文件路径

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

### pdfSplit(inFile, outDir)

将 PDF 文件拆分为单页 PDF。

**参数**:
- `inFile` (string): 输入 PDF 文件路径
- `outDir` (string): 输出目录

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

### pdfExtractPages(inFile, outDir, pages)

提取 PDF 中的指定页面。

**参数**:
- `inFile` (string): 输入 PDF 文件路径
- `outDir` (string): 输出目录
- `pages` (array): 要提取的页码或范围，例如 `["1", "2-5", "8"]`

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

### pdfOptimize(inFile, outFile)

优化 PDF 文件大小。

**参数**:
- `inFile` (string): 输入 PDF 文件路径
- `outFile` (string): 输出 PDF 文件路径

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

### pdfValidate(inFile)

验证 PDF 文件的有效性。

**参数**:
- `inFile` (string): PDF 文件路径

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

### pdfAddTextWatermark(inFile, outFile, text, options?)

给 PDF 添加文本水印。

**参数**:
- `inFile` (string): 输入 PDF 文件路径
- `outFile` (string): 输出 PDF 文件路径
- `text` (string): 水印文本
- `options` (object, 可选): 水印选项
  - `onTop` (boolean): 是否在内容上方（默认 true）
  - `opacity` (number): 不透明度 0.0-1.0
  - `scale` (number): 缩放比例
  - `rotation` (number): 旋转角度

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

### pdfExportImages(inFile, outDir)

从 PDF 中导出所有图片。

**参数**:
- `inFile` (string): 输入 PDF 文件路径
- `outDir` (string): 输出目录

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

### pdfImportImages(imgFiles, outFile)

将多张图片导入并合并为一个 PDF。

**参数**:
- `imgFiles` (array): 输入图片文件路径数组
- `outFile` (string): 输出 PDF 文件路径

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

## 文档读取

所有文档读取函数支持分页读取，通过 `options` 参数控制。

### readWord(path, options?)

读取Word文档

**参数**:
- `path` (string): Word文件路径（.docx）
- `options` (object, 可选): 读取选项
  - `page` (number): 页码（从1开始）
  - `pageSize` (number): 每页大小（字符数）

**返回值**: `object`
- `text` (string): 文档文本内容
- `totalPages` (number): 总页数
- `page` (number): 当前页码
- `error` (string, 可选): 错误信息

**示例**:
```javascript
// 读取第一页
var word = readWord("document.docx", {page: 1, pageSize: 1000});
console.log("总页数:", word.totalPages);
console.log("内容:", word.text);
```

### readExcel(path, options?)

读取Excel文件

**参数**:
- `path` (string): Excel文件路径（.xlsx）
- `options` (object, 可选): 读取选项
  - `page` (number): 页码（从1开始）
  - `pageSize` (number): 每页行数

**返回值**: `object`
- `rows` (array): 数据行数组，每行是一个对象数组
- `totalRows` (number): 总行数
- `totalPages` (number): 总页数
- `page` (number): 当前页码
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var excel = readExcel("data.xlsx", {page: 1, pageSize: 100});
console.log("总行数:", excel.totalRows);
console.log("数据:", excel.rows);
// rows 示例: [["姓名", "年龄"], ["张三", "25"], ["李四", "30"]]
```

### readPPT(path, options?)

读取PPT文件

**参数**:
- `path` (string): PPT文件路径（.pptx）
- `options` (object, 可选): 读取选项
  - `page` (number): 页码（从1开始）
  - `pageSize` (number): 每页大小（字符数）

**返回值**: 同 `readWord`

**示例**:
```javascript
var ppt = readPPT("presentation.pptx", {page: 1, pageSize: 1000});
console.log("总页数:", ppt.totalPages);
console.log("内容:", ppt.text);
```

### readPDF(path, options?)

读取PDF文件

**参数**:
- `path` (string): PDF文件路径
- `options` (object, 可选): 读取选项
  - `page` (number): 页码（从1开始）
  - `pageSize` (number): 每页大小（字符数）

**返回值**: 同 `readWord`

**示例**:
```javascript
var pdf = readPDF("document.pdf", {page: 1, pageSize: 1000});
console.log("总页数:", pdf.totalPages);
console.log("内容:", pdf.text);
```

---

## 浏览器自动化

### 浏览器会话管理

浏览器操作通过会话（Session）进行管理，需要先创建会话，使用完毕后关闭。

### createBrowserSession(timeoutSeconds?)

创建浏览器会话

**参数**:
- `timeoutSeconds` (number, 可选): 会话超时时间（秒），默认30

**返回值**: `object` - 浏览器会话对象，包含以下方法：
- `navigate(url)` - 导航到URL
- `wait(selectorOrSeconds)` - 等待元素或指定秒数
- `click(selector)` - 点击元素
- `fill(selector, value)` - 填充表单
- `evaluate(jsCode)` - 在页面中执行JavaScript
- `getHTML()` - 获取页面HTML
- `screenshot(outputPath)` - 截图
- `getURL()` - 获取当前URL
- `waitForURL(pattern, timeout?)` - 等待URL匹配
- `close()` - 关闭会话

**示例**:
```javascript
var session = createBrowserSession(120); // 120秒超时
try {
    // 导航
    var navResult = session.navigate("https://www.example.com");
    if (!navResult.success) {
        throw new Error("导航失败: " + navResult.error);
    }
    
    // 等待页面加载
    session.wait(2); // 等待2秒
    
    // 执行JavaScript
    var evalResult = session.evaluate("document.title");
    console.log("页面标题:", evalResult.result);
    
    // 点击元素
    var clickResult = session.click("#button");
    
    // 填充表单
    var fillResult = session.fill("#input", "value");
    
    // 获取HTML
    var htmlResult = session.getHTML();
    console.log("HTML长度:", htmlResult.html.length);
    
    // 截图
    var screenshotResult = session.screenshot("screenshot.png");
    
    // 获取当前URL
    var urlResult = session.getURL();
    console.log("当前URL:", urlResult.url);
} finally {
    session.close(); // 必须关闭会话
}
```

### session.navigate(url)

导航到指定URL

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

### session.wait(selectorOrSeconds)

等待元素出现或等待指定秒数

**参数**:
- `selectorOrSeconds` (string|number): CSS选择器或秒数

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

**示例**:
```javascript
session.wait(2); // 等待2秒
session.wait("#element"); // 等待元素出现
```

### session.click(selector)

点击元素

**参数**:
- `selector` (string): CSS选择器

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

### session.fill(selector, value)

填充表单字段

**参数**:
- `selector` (string): CSS选择器
- `value` (string): 要填充的值

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

### session.evaluate(jsCode)

在页面中执行JavaScript代码

**参数**:
- `jsCode` (string): JavaScript代码

**返回值**: `object`
- `success` (boolean): 是否成功
- `result` (any): 执行结果
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var result = session.evaluate("document.title");
console.log("标题:", result.result);

var result2 = session.evaluate("document.querySelectorAll('.item').length");
console.log("元素数量:", result2.result);
```

### session.getHTML()

获取页面HTML

**返回值**: `object`
- `success` (boolean): 是否成功
- `html` (string): HTML内容
- `error` (string, 可选): 错误信息

### session.screenshot(outputPath)

截图

**参数**:
- `outputPath` (string): 输出文件路径

**返回值**: `object`
- `success` (boolean): 是否成功
- `path` (string): 截图文件路径
- `error` (string, 可选): 错误信息

### session.getURL()

获取当前URL

**返回值**: `object`
- `success` (boolean): 是否成功
- `url` (string): 当前URL
- `error` (string, 可选): 错误信息

### session.waitForURL(pattern, timeout?)

等待URL匹配指定模式

**参数**:
- `pattern` (string): URL模式（支持正则表达式）
- `timeout` (number, 可选): 超时时间（秒），默认10

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

### session.close()

关闭浏览器会话（必须调用）

**返回值**: `undefined`

---

## 图片处理

### imageInfo(filePath)

获取图片信息

**参数**:
- `filePath` (string): 图片文件路径

**返回值**: `object`
- `width` (number): 图片宽度（像素）
- `height` (number): 图片高度（像素）
- `format` (string): 图片格式（如 "jpeg", "png"）
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var info = imageInfo("photo.jpg");
console.log("尺寸:", info.width, "x", info.height);
console.log("格式:", info.format);
```

### imageResize(inputPath, outputPath, width, height?)

调整图片大小

**参数**:
- `inputPath` (string): 输入图片路径
- `outputPath` (string): 输出图片路径
- `width` (number): 目标宽度（像素）
- `height` (number, 可选): 目标高度（像素），如果省略则保持宽高比

**返回值**: `object`
- `success` (boolean): 是否成功
- `path` (string): 输出文件路径
- `error` (string, 可选): 错误信息

**示例**:
```javascript
// 按宽度缩放，保持宽高比
var result = imageResize("input.jpg", "output.jpg", 800);

// 指定宽度和高度
var result2 = imageResize("input.jpg", "output.jpg", 800, 600);
```

### imageCrop(inputPath, outputPath, x, y, width, height)

裁剪图片

**参数**:
- `inputPath` (string): 输入图片路径
- `outputPath` (string): 输出图片路径
- `x` (number): 裁剪区域左上角X坐标
- `y` (number): 裁剪区域左上角Y坐标
- `width` (number): 裁剪区域宽度
- `height` (number): 裁剪区域高度

**返回值**: `object`
- `success` (boolean): 是否成功
- `path` (string): 输出文件路径
- `error` (string, 可选): 错误信息

**示例**:
```javascript
// 从左上角裁剪400x300的区域
var result = imageCrop("input.jpg", "output.jpg", 0, 0, 400, 300);
```

### imageRotate(inputPath, outputPath, angle)

旋转图片

**参数**:
- `inputPath` (string): 输入图片路径
- `outputPath` (string): 输出图片路径
- `angle` (number): 旋转角度（度，顺时针）

**返回值**: `object`
- `success` (boolean): 是否成功
- `path` (string): 输出文件路径
- `error` (string, 可选): 错误信息

**示例**:
```javascript
// 旋转90度
imageRotate("input.jpg", "output.jpg", 90);
// 旋转180度
imageRotate("input.jpg", "output.jpg", 180);
```

### imageFlip(inputPath, outputPath, direction)

翻转图片

**参数**:
- `inputPath` (string): 输入图片路径
- `outputPath` (string): 输出图片路径
- `direction` (string): 翻转方向，"horizontal"（水平）或 "vertical"（垂直）

**返回值**: `object`
- `success` (boolean): 是否成功
- `path` (string): 输出文件路径
- `error` (string, 可选): 错误信息

**示例**:
```javascript
// 水平翻转
imageFlip("input.jpg", "output.jpg", "horizontal");
// 垂直翻转
imageFlip("input.jpg", "output.jpg", "vertical");
```

### imageConvert(inputPath, outputPath)

转换图片格式

**参数**:
- `inputPath` (string): 输入图片路径
- `outputPath` (string): 输出图片路径（扩展名决定格式）

**返回值**: `object`
- `success` (boolean): 是否成功
- `path` (string): 输出文件路径
- `error` (string, 可选): 错误信息

**示例**:
```javascript
// 转换为PNG
imageConvert("input.jpg", "output.png");
// 转换为JPEG
imageConvert("input.png", "output.jpg");
```

### imageQuality(inputPath, outputPath, quality)

调整JPEG图片质量

**参数**:
- `inputPath` (string): 输入图片路径
- `outputPath` (string): 输出图片路径
- `quality` (number): 质量（1-100，100为最高质量）

**返回值**: `object`
- `success` (boolean): 是否成功
- `path` (string): 输出文件路径
- `error` (string, 可选): 错误信息

**示例**:
```javascript
// 高质量保存
imageQuality("input.jpg", "output.jpg", 95);
// 中等质量
imageQuality("input.jpg", "output.jpg", 75);
// 低质量（文件更小）
imageQuality("input.jpg", "output.jpg", 50);
```

---

## 文件类型检测

### detectFileType(filePath)

检测文件类型

**参数**:
- `filePath` (string): 文件路径

**返回值**: `object`
- `mime` (string): MIME类型
- `extension` (string): 文件扩展名
- `unknown` (boolean): 是否未知类型

**示例**:
```javascript
var type = detectFileType("file.bin");
if (!type.unknown) {
    console.log("MIME类型:", type.mime);
    console.log("扩展名:", type.extension);
}
```

### isImage(filePath)

检测是否为图片

**参数**:
- `filePath` (string): 文件路径

**返回值**: `object`
- `isImage` (boolean): 是否为图片

**示例**:
```javascript
var result = isImage("photo.jpg");
console.log("是图片:", result.isImage);
```

### isAudio(filePath)

检测是否为音频

**返回值**: `object`
- `isAudio` (boolean): 是否为音频

### isDocument(filePath)

检测是否为文档

**返回值**: `object`
- `isDocument` (boolean): 是否为文档

### isFont(filePath)

检测是否为字体

**返回值**: `object`
- `isFont` (boolean): 是否为字体

### isArchive(filePath)

检测是否为归档文件

**返回值**: `object`
- `isArchive` (boolean): 是否为归档文件

---

## 加密/解密

### encryptAES(data, key)

AES加密数据

**参数**:
- `data` (string): 要加密的数据
- `key` (string): 加密密钥（32字节，256位）

**返回值**: `object`
- `data` (string): 加密后的数据（base64编码）
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var encrypted = encryptAES("敏感数据", "12345678901234567890123456789012");
console.log("加密结果:", encrypted.data);
```

### decryptAES(encrypted, key)

AES解密数据

**参数**:
- `encrypted` (string): 加密的数据（base64编码）
- `key` (string): 解密密钥（必须与加密密钥相同）

**返回值**: `object`
- `data` (string): 解密后的数据
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var decrypted = decryptAES(encrypted.data, "12345678901234567890123456789012");
console.log("解密结果:", decrypted.data);
```

### hashSHA256(data)

计算SHA256哈希值

**参数**:
- `data` (string): 要哈希的数据

**返回值**: `object`
- `hash` (string): SHA256哈希值（十六进制）
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var hash = hashSHA256("数据");
console.log("SHA256:", hash.hash);
```

### generateUUID()

生成UUID

**返回值**: `string` - UUID字符串

**示例**:
```javascript
var uuid = generateUUID();
console.log("UUID:", uuid); // "550e8400-e29b-41d4-a716-446655440000"
```

### generateRandomString(length?)

生成随机字符串

**参数**:
- `length` (number, 可选): 字符串长度，默认32

**返回值**: `object`
- `data` (string): 随机字符串
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var random = generateRandomString(16);
console.log("随机字符串:", random.data);
```

---

## 压缩/解压缩

### compressZip(files, outputPath)

压缩文件为ZIP

**参数**:
- `files` (array): 要压缩的文件路径数组
- `outputPath` (string): 输出ZIP文件路径

**返回值**: `object`
- `success` (boolean): 是否成功
- `path` (string): 输出文件路径
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var result = compressZip(["file1.txt", "file2.txt"], "archive.zip");
if (result.success) {
    console.log("压缩成功:", result.path);
}
```

### extractZip(zipPath, outputDir)

解压ZIP文件

**参数**:
- `zipPath` (string): ZIP文件路径
- `outputDir` (string): 输出目录

**返回值**: `object`
- `success` (boolean): 是否成功
- `files` (array): 解压出的文件路径数组
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var extracted = extractZip("archive.zip", "./output");
console.log("解压文件数:", extracted.files.length);
```

### compressGzip(data)

GZIP压缩字符串

**参数**:
- `data` (string): 要压缩的数据

**返回值**: `object`
- `data` (string): 压缩后的数据（base64编码）
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var compressed = compressGzip("要压缩的数据");
console.log("压缩后:", compressed.data);
```

### decompressGzip(compressed)

GZIP解压字符串

**参数**:
- `compressed` (string): 压缩的数据（base64编码）

**返回值**: `object`
- `data` (string): 解压后的数据
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var decompressed = decompressGzip(compressed.data);
console.log("解压后:", decompressed.data);
```

---

## CSV处理

### readCSV(filePath, options?)

读取CSV文件

**参数**:
- `filePath` (string): CSV文件路径
- `options` (object, 可选): 读取选项
  - `delimiter` (string): 分隔符，默认 ","
  - `comment` (string): 注释符，默认 ""
  - `skipEmptyLines` (boolean): 是否跳过空行，默认 false
  - `trimSpace` (boolean): 是否去除空格，默认 false

**返回值**: `object`
- `rows` (array): 数据行数组，每行是一个字符串数组
- `count` (number): 行数
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var csv = readCSV("data.csv", {
    delimiter: ",",
    skipEmptyLines: true
});
console.log("行数:", csv.count);
console.log("数据:", csv.rows);
// rows 示例: [["姓名", "年龄"], ["张三", "25"]]
```

### writeCSV(filePath, data, options?)

写入CSV文件

**参数**:
- `filePath` (string): CSV文件路径
- `data` (array): 数据数组，每行是一个字符串数组
- `options` (object, 可选): 写入选项
  - `delimiter` (string): 分隔符，默认 ","
  - `header` (array, 可选): 表头数组

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var data = [
    ["姓名", "年龄", "城市"],
    ["张三", "25", "北京"],
    ["李四", "30", "上海"]
];
writeCSV("output.csv", data);
```

### parseCSV(csvString, options?)

解析CSV字符串

**参数**:
- `csvString` (string): CSV字符串
- `options` (object, 可选): 解析选项，同 `readCSV`

**返回值**: 同 `readCSV`

**示例**:
```javascript
var parsed = parseCSV("a,b,c\n1,2,3", {delimiter: ","});
console.log("解析结果:", parsed.rows);
```

---

## 环境变量和配置

### getEnv(name)

获取环境变量

**参数**:
- `name` (string): 环境变量名

**返回值**: `object`
- `value` (string): 环境变量值
- `exists` (boolean): 是否存在

**示例**:
```javascript
var env = getEnv("PATH");
if (env.exists) {
    console.log("PATH:", env.value);
}
```

### getEnvAll()

获取所有环境变量

**返回值**: `object`
- `env` (object): 环境变量对象（键值对）

**示例**:
```javascript
var allEnv = getEnvAll();
console.log("环境变量数量:", Object.keys(allEnv.env).length);
```

### readConfig(filePath)

读取配置文件（支持JSON和YAML）

**参数**:
- `filePath` (string): 配置文件路径

**返回值**: `object`
- `config` (object): 配置对象
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var config = readConfig("config.yaml");
console.log("配置:", config.config);
```

---

## 数据验证

### validateEmail(email)

验证邮箱格式

**参数**:
- `email` (string): 邮箱地址

**返回值**: `object`
- `valid` (boolean): 是否有效

**示例**:
```javascript
var emailValid = validateEmail("test@example.com");
console.log("邮箱有效:", emailValid.valid);
```

### validateURL(url)

验证URL格式

**参数**:
- `url` (string): URL地址

**返回值**: `object`
- `valid` (boolean): 是否有效

**示例**:
```javascript
var urlValid = validateURL("https://www.example.com");
console.log("URL有效:", urlValid.valid);
```

### validateIP(ip)

验证IP地址

**参数**:
- `ip` (string): IP地址

**返回值**: `object`
- `valid` (boolean): 是否有效
- `isIPv4` (boolean): 是否为IPv4
- `isIPv6` (boolean): 是否为IPv6

**示例**:
```javascript
var ipValid = validateIP("192.168.1.1");
console.log("IP有效:", ipValid.valid);
console.log("是IPv4:", ipValid.isIPv4);
```

### validatePhone(phone)

验证中国手机号

**参数**:
- `phone` (string): 手机号

**返回值**: `object`
- `valid` (boolean): 是否有效

**示例**:
```javascript
var phoneValid = validatePhone("13800138000");
console.log("手机号有效:", phoneValid.valid);
```

---

## 日期时间增强

### formatDate(date, format)

格式化日期

**参数**:
- `date` (string|number): 日期字符串或时间戳
- `format` (string): 格式字符串，如 "YYYY-MM-DD HH:mm:ss"

**返回值**: `object`
- `date` (string): 格式化后的日期字符串
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var formatted = formatDate("2024-01-01 12:00:00", "YYYY-MM-DD HH:mm:ss");
console.log("格式化后:", formatted.date);
```

### parseDate(dateString)

解析日期字符串

**参数**:
- `dateString` (string): 日期字符串

**返回值**: `object`
- `timestamp` (number): 时间戳（秒）
- `year` (number): 年份
- `month` (number): 月份（1-12）
- `day` (number): 日期（1-31）
- `hour` (number): 小时（0-23）
- `minute` (number): 分钟（0-59）
- `second` (number): 秒（0-59）
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var parsed = parseDate("2024-01-01");
console.log("时间戳:", parsed.timestamp);
console.log("年:", parsed.year);
```

### addDays(date, days)

日期加减天数

**参数**:
- `date` (string): 日期字符串
- `days` (number): 天数（正数为加，负数为减）

**返回值**: `object`
- `date` (string): 计算后的日期字符串
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var newDate = addDays("2024-01-01", 7);
console.log("7天后:", newDate.date);
```

### getTimezone()

获取当前时区

**返回值**: `object`
- `timezone` (string): 时区名称

**示例**:
```javascript
var tz = getTimezone();
console.log("时区:", tz.timezone);
```

### convertTimezone(date, timezone)

时区转换

**参数**:
- `date` (string): 日期字符串
- `timezone` (string): 目标时区，如 "America/New_York"

**返回值**: `object`
- `date` (string): 转换后的日期字符串
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var converted = convertTimezone("2024-01-01 12:00:00", "America/New_York");
console.log("转换后:", converted.date);
```

---

## 编码/解码

### encodeBase64(data)

Base64编码

**参数**:
- `data` (string): 要编码的数据

**返回值**: `object`
- `data` (string): Base64编码后的字符串
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var encoded = encodeBase64("Hello World");
console.log("Base64:", encoded.data);
```

### decodeBase64(encoded)

Base64解码

**参数**:
- `encoded` (string): Base64编码的字符串

**返回值**: `object`
- `data` (string): 解码后的字符串
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var decoded = decodeBase64(encoded.data);
console.log("解码后:", decoded.data);
```

### encodeURL(str)

URL编码

**参数**:
- `str` (string): 要编码的字符串

**返回值**: `object`
- `data` (string): URL编码后的字符串
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var urlEncoded = encodeURL("hello world");
console.log("URL编码:", urlEncoded.data); // "hello%20world"
```

### decodeURL(encoded)

URL解码

**参数**:
- `encoded` (string): URL编码的字符串

**返回值**: `object`
- `data` (string): 解码后的字符串
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var decoded = decodeURL("hello%20world");
console.log("解码后:", decoded.data); // "hello world"
```

### encodeHTML(str)

HTML实体编码

**参数**:
- `str` (string): 要编码的字符串

**返回值**: `object`
- `data` (string): HTML编码后的字符串
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var htmlEncoded = encodeHTML("<div>内容</div>");
console.log("HTML编码:", htmlEncoded.data); // "&lt;div&gt;内容&lt;/div&gt;"
```

### decodeHTML(encoded)

HTML实体解码

**参数**:
- `encoded` (string): HTML编码的字符串

**返回值**: `object`
- `data` (string): 解码后的字符串
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var decoded = decodeHTML("&lt;div&gt;内容&lt;/div&gt;");
console.log("解码后:", decoded.data); // "<div>内容</div>"
```

---

## 进程管理

### execCommand(command, options?)

执行系统命令

**参数**:
- `command` (string): 要执行的命令
- `options` (object, 可选): 执行选项
  - `timeout` (number): 超时时间（秒），默认30
  - `dir` (string): 工作目录
  - `env` (object): 环境变量

**返回值**: `object`
- `output` (string): 命令输出
- `code` (number): 退出码
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var result = execCommand("ls", {timeout: 5});
console.log("输出:", result.output);
console.log("退出码:", result.code);
```

### listProcesses()

列出运行中的进程

**返回值**: `object`
- `count` (number): 进程数量
- `processes` (array): 进程数组，每个进程包含：
  - `pid` (number): 进程ID
  - `name` (string): 进程名称
  - `cpu` (number): CPU使用率
  - `memory` (number): 内存使用（字节）

**示例**:
```javascript
var processes = listProcesses();
console.log("进程数:", processes.count);
processes.processes.forEach(function(p) {
    console.log("进程:", p.name, "PID:", p.pid);
});
```

### killProcess(pid)

终止进程

**参数**:
- `pid` (number): 进程ID

**返回值**: `object`
- `success` (boolean): 是否成功
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var result = killProcess(12345);
if (result.success) {
    console.log("进程已终止");
}
```

---

## 网络工具

### resolveDNS(hostname)

DNS解析

**参数**:
- `hostname` (string): 主机名

**返回值**: `object`
- `ips` (array): IP地址数组
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var dns = resolveDNS("www.example.com");
console.log("IP地址:", dns.ips);
```

### ping(host, count?)

Ping测试

**参数**:
- `host` (string): 主机地址
- `count` (number, 可选): Ping次数，默认4

**返回值**: `object`
- `sent` (number): 发送的包数
- `received` (number): 接收的包数
- `averageTime` (number): 平均响应时间（毫秒）
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var ping = ping("www.example.com", 4);
console.log("发送:", ping.sent);
console.log("接收:", ping.received);
console.log("平均时间:", ping.averageTime, "ms");
```

### checkPort(host, port, timeout?)

检查端口是否开放

**参数**:
- `host` (string): 主机地址
- `port` (number): 端口号
- `timeout` (number, 可选): 超时时间（秒），默认3

**返回值**: `object`
- `open` (boolean): 端口是否开放
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var port = checkPort("localhost", 80, 3);
console.log("端口开放:", port.open);
```

---

## 路径处理

### pathJoin(...paths)

路径拼接

**参数**:
- `...paths` (string): 多个路径参数

**返回值**: `object`
- `path` (string): 拼接后的路径

**示例**:
```javascript
var path = pathJoin("/usr", "local", "bin");
console.log("拼接路径:", path.path); // "/usr/local/bin"
```

### pathDir(path)

获取目录

**参数**:
- `path` (string): 文件路径

**返回值**: `object`
- `dir` (string): 目录路径

**示例**:
```javascript
var dir = pathDir("/usr/local/bin/app");
console.log("目录:", dir.dir); // "/usr/local/bin"
```

### pathBase(path)

获取文件名

**参数**:
- `path` (string): 文件路径

**返回值**: `object`
- `base` (string): 文件名

**示例**:
```javascript
var base = pathBase("/usr/local/bin/app");
console.log("文件名:", base.base); // "app"
```

### pathExt(path)

获取扩展名

**参数**:
- `path` (string): 文件路径

**返回值**: `object`
- `ext` (string): 扩展名（包含点号）

**示例**:
```javascript
var ext = pathExt("file.txt");
console.log("扩展名:", ext.ext); // ".txt"
```

### pathAbs(path)

获取绝对路径

**参数**:
- `path` (string): 文件路径

**返回值**: `object`
- `path` (string): 绝对路径
- `error` (string, 可选): 错误信息

**示例**:
```javascript
var abs = pathAbs("./file.txt");
console.log("绝对路径:", abs.path);
```

---

## 文本操作

### textReplace(text, old, new)

替换文本

**参数**:
- `text` (string): 原文本
- `old` (string): 要替换的字符串
- `new` (string): 替换为的字符串

**返回值**: `object`
- `result` (string): 替换后的文本

**示例**:
```javascript
var result = textReplace("Hello World", "World", "JavaScript");
console.log(result.result); // "Hello JavaScript"
```

### textSplit(text, separator)

分割文本

**参数**:
- `text` (string): 原文本
- `separator` (string): 分隔符

**返回值**: `object`
- `parts` (array): 分割后的字符串数组

**示例**:
```javascript
var result = textSplit("a,b,c", ",");
console.log(result.parts); // ["a", "b", "c"]
```

### textJoin(parts, separator)

连接文本

**参数**:
- `parts` (array): 字符串数组
- `separator` (string): 分隔符

**返回值**: `object`
- `result` (string): 连接后的文本

**示例**:
```javascript
var result = textJoin(["a", "b", "c"], ",");
console.log(result.result); // "a,b,c"
```

### textTrim(text)

去除首尾空格

**参数**:
- `text` (string): 原文本

**返回值**: `object`
- `result` (string): 处理后的文本

**示例**:
```javascript
var result = textTrim("  Hello World  ");
console.log(result.result); // "Hello World"
```

### textToUpper(text)

转换为大写

**参数**:
- `text` (string): 原文本

**返回值**: `object`
- `result` (string): 转换后的文本

**示例**:
```javascript
var result = textToUpper("hello");
console.log(result.result); // "HELLO"
```

### textToLower(text)

转换为小写

**参数**:
- `text` (string): 原文本

**返回值**: `object`
- `result` (string): 转换后的文本

**示例**:
```javascript
var result = textToLower("HELLO");
console.log(result.result); // "hello"
```

### textContains(text, substr)

检查是否包含子字符串

**参数**:
- `text` (string): 原文本
- `substr` (string): 子字符串

**返回值**: `object`
- `contains` (boolean): 是否包含

**示例**:
```javascript
var result = textContains("Hello World", "World");
console.log(result.contains); // true
```

### textStartsWith(text, prefix)

检查是否以指定字符串开头

**参数**:
- `text` (string): 原文本
- `prefix` (string): 前缀

**返回值**: `object`
- `startsWith` (boolean): 是否以指定字符串开头

**示例**:
```javascript
var result = textStartsWith("Hello World", "Hello");
console.log(result.startsWith); // true
```

### textEndsWith(text, suffix)

检查是否以指定字符串结尾

**参数**:
- `text` (string): 原文本
- `suffix` (string): 后缀

**返回值**: `object`
- `endsWith` (boolean): 是否以指定字符串结尾

**示例**:
```javascript
var result = textEndsWith("Hello World", "World");
console.log(result.endsWith); // true
```

### textSubstring(text, start, end?)

截取子字符串

**参数**:
- `text` (string): 原文本
- `start` (number): 起始位置
- `end` (number, 可选): 结束位置（不包含）

**返回值**: `object`
- `result` (string): 截取后的文本

**示例**:
```javascript
var result = textSubstring("Hello World", 0, 5);
console.log(result.result); // "Hello"
```

---

## 日志功能

日志功能提供了完整的日志记录能力，支持不同级别的日志输出、结构化日志和日志级别管理。

### logger.debug(...args)

输出 Debug 级别日志

**参数**:
- `...args` (any): 可变参数，可以是字符串、对象等

**返回值**: `undefined`

**示例**:
```javascript
// 简单消息
logger.debug("调试信息");

// 多个参数
logger.debug("用户ID:", 123, "操作:", "登录");

// 结构化日志（第一个参数为对象）
logger.debug({userId: 123, action: "login"}, "用户登录");
```

### logger.info(...args)

输出 Info 级别日志

**参数**:
- `...args` (any): 可变参数

**返回值**: `undefined`

**示例**:
```javascript
logger.info("应用程序启动");
logger.info("处理了", 100, "条记录");
logger.info({count: 100, type: "records"}, "处理完成");
```

### logger.warn(...args)

输出 Warn 级别日志

**参数**:
- `...args` (any): 可变参数

**返回值**: `undefined`

**示例**:
```javascript
logger.warn("警告：内存使用率较高");
logger.warn({memory: "85%", threshold: "80%"}, "内存使用超过阈值");
```

### logger.error(...args)

输出 Error 级别日志

**参数**:
- `...args` (any): 可变参数

**返回值**: `undefined`

**示例**:
```javascript
logger.error("发生错误：文件读取失败");
logger.error({file: "data.txt", error: "permission denied"}, "文件操作失败");
```

### logger.fatal(...args)

输出 Fatal 级别日志（实际记录为 Error 级别）

**参数**:
- `...args` (any): 可变参数

**返回值**: `undefined`

**示例**:
```javascript
logger.fatal("致命错误：系统无法继续运行");
logger.fatal({component: "database", error: "connection lost"}, "关键组件故障");
```

### logger.trace(...args)

输出 Trace 级别日志

**参数**:
- `...args` (any): 可变参数

**返回值**: `undefined`

**示例**:
```javascript
logger.trace("函数调用跟踪");
logger.trace({function: "processData", step: 1}, "开始处理数据");
```

### logger.setLevel(level)

设置日志级别

**参数**:
- `level` (string): 日志级别，可选值: "trace", "debug", "info", "warn"/"warning", "error", "fatal", "panic"

**返回值**: `object`
- `success` (boolean): 是否成功
- `level` (string): 设置的日志级别
- `error` (string, 可选): 错误信息（如果级别无效）

**示例**:
```javascript
// 设置为 Debug 级别（会显示所有日志）
var result = logger.setLevel("debug");
if (result.success) {
    console.log("日志级别已设置为:", result.level);
}

// 设置为 Error 级别（只显示错误和致命错误）
logger.setLevel("error");
```

### logger.getLevel()

获取当前日志级别

**返回值**: `object`
- `success` (boolean): 是否成功
- `level` (string): 当前日志级别

**示例**:
```javascript
var levelInfo = logger.getLevel();
console.log("当前日志级别:", levelInfo.level);
```

### logger.isLevelEnabled(level)

检查某个日志级别是否启用

**参数**:
- `level` (string): 要检查的日志级别

**返回值**: `object`
- `success` (boolean): 是否成功
- `enabled` (boolean): 该级别是否启用
- `error` (string, 可选): 错误信息（如果级别无效）

**示例**:
```javascript
var check = logger.isLevelEnabled("debug");
if (check.enabled) {
    logger.debug("Debug 日志已启用，这条消息会显示");
} else {
    console.log("Debug 日志未启用");
}
```

### logger.withFields(fields)

创建带字段的日志记录器（结构化日志）

**参数**:
- `fields` (object): 要附加的字段对象

**返回值**: `object` - 返回一个新的日志记录器对象，包含所有日志级别方法（debug, info, warn, error, fatal, trace）

**示例**:
```javascript
// 创建带字段的日志记录器
var userLogger = logger.withFields({userId: 123, username: "john"});

// 使用带字段的日志记录器
userLogger.info("用户登录");
// 输出: 包含 userId 和 username 字段的日志

userLogger.error("操作失败");
// 输出: 包含 userId 和 username 字段的错误日志

// 链式使用
logger.withFields({requestId: "req-123"})
    .info("处理请求")
    .warn("请求处理较慢");
```

### 日志功能完整示例

```javascript
// 1. 设置日志级别
logger.setLevel("debug");

// 2. 基本日志输出
logger.debug("调试信息");
logger.info("应用程序启动");
logger.warn("警告信息");
logger.error("错误信息");

// 3. 多参数日志
logger.info("用户", "john", "执行了操作", "login");

// 4. 结构化日志（直接方式）
logger.info({userId: 123, action: "login"}, "用户登录");

// 5. 结构化日志（使用 withFields）
var requestLogger = logger.withFields({
    requestId: "req-123",
    method: "POST",
    path: "/api/users"
});
requestLogger.info("收到请求");
requestLogger.debug("处理中...");
requestLogger.info("请求处理完成");

// 6. 条件日志
var levelCheck = logger.isLevelEnabled("debug");
if (levelCheck.enabled) {
    logger.debug("详细的调试信息");
}

// 7. 获取当前日志级别
var currentLevel = logger.getLevel();
console.log("当前日志级别:", currentLevel.level);
```

### 日志级别说明

日志级别从低到高：
- **trace**: 最详细的跟踪信息
- **debug**: 调试信息
- **info**: 一般信息（默认级别）
- **warn**: 警告信息
- **error**: 错误信息
- **fatal**: 致命错误
- **panic**: 系统恐慌

设置某个级别后，只有该级别及更高级别的日志会被输出。例如，设置为 "warn" 后，只会输出 warn、error、fatal 和 panic 级别的日志。

### 注意事项

1. **结构化日志**: 当第一个参数是对象时，会自动提取对象字段作为结构化日志的字段。
2. **日志级别**: 默认日志级别为 "info"，可以通过 `setLevel()` 修改。
3. **console 对象**: 除了 `logger` 对象，还可以使用 `console.log()`, `console.error()` 等（在系统操作章节中说明）。
4. **性能**: 在生产环境中，建议将日志级别设置为 "warn" 或 "error" 以提高性能。

---

## 错误处理

### 通用错误处理模式

大多数函数返回包含 `success` 或 `error` 字段的对象：

```javascript
// 方式1: 检查success字段
var result = writeFile("test.txt", "content");
if (result.success) {
    console.log("操作成功");
} else {
    console.error("操作失败:", result.error);
}

// 方式2: 检查error字段
var result2 = readFile("test.txt");
if (result2.error) {
    console.error("读取失败:", result2.error);
} else {
    console.log("内容:", result2.data);
}
```

### try-catch 处理

对于可能抛出异常的操作（如浏览器会话），使用 try-catch：

```javascript
try {
    var session = createBrowserSession(30);
    var navResult = session.navigate("https://example.com");
    if (!navResult.success) {
        throw new Error("导航失败: " + navResult.error);
    }
    // ... 其他操作
    session.close();
} catch (error) {
    console.error("发生错误:", error.message);
}
```

### 超时处理

对于可能长时间运行的操作，建议设置超时：

```javascript
// 浏览器操作建议设置较长的超时时间
var session = createBrowserSession(120); // 120秒

// HTTP请求可以通过options设置超时
var response = httpRequest("https://api.example.com", {
    timeout: 60 // 60秒
});
```

---

## 完整示例

### 示例1: 文件处理和HTTP请求

```javascript
// 1. 下载文件
var response = httpGet("https://example.com/data.json");
if (response.status === 200) {
    // 2. 保存文件
    writeFile("data.json", response.body);
    
    // 3. 读取文件
    var file = readFile("data.json");
    var data = JSON.parse(file.data);
    
    // 4. 处理数据
    console.log("数据:", data);
    
    // 5. 获取文件哈希
    var hash = getFileHash("data.json", "sha256");
    console.log("文件哈希:", hash.hash);
}
```

### 示例2: 浏览器自动化

```javascript
var session = createBrowserSession(120);
try {
    // 导航
    var navResult = session.navigate("https://example.com");
    if (!navResult.success) {
        throw new Error("导航失败");
    }
    
    // 等待页面加载
    session.wait(2);
    
    // 提取数据
    var titleResult = session.evaluate("document.title");
    console.log("页面标题:", titleResult.result);
    
    // 点击按钮
    session.click("#submit-button");
    session.wait("#result");
    
    // 获取结果
    var htmlResult = session.getHTML();
    console.log("页面HTML长度:", htmlResult.html.length);
    
    // 截图
    session.screenshot("screenshot.png");
} finally {
    session.close();
}
```

### 示例3: 图片处理

```javascript
// 获取图片信息
var info = imageInfo("photo.jpg");
console.log("原始尺寸:", info.width, "x", info.height);

// 调整大小
imageResize("photo.jpg", "photo_resized.jpg", 800);

// 裁剪
imageCrop("photo.jpg", "photo_cropped.jpg", 0, 0, 400, 300);

// 旋转
imageRotate("photo.jpg", "photo_rotated.jpg", 90);

// 翻转
imageFlip("photo.jpg", "photo_flipped.jpg", "horizontal");

// 转换格式
imageConvert("photo.jpg", "photo.png");

// 调整质量（仅JPEG）
imageQuality("photo.jpg", "photo_quality.jpg", 85);
```

### 示例4: 数据验证和处理

```javascript
// 验证邮箱
var email = "test@example.com";
var emailValid = validateEmail(email);
if (!emailValid.valid) {
    console.error("邮箱格式无效");
}

// 验证URL
var url = "https://www.example.com";
var urlValid = validateURL(url);
if (urlValid.valid) {
    // 发送HTTP请求
    var response = httpGet(url);
    console.log("响应状态:", response.status);
}

// 验证IP
var ip = "192.168.1.1";
var ipValid = validateIP(ip);
if (ipValid.valid && ipValid.isIPv4) {
    // 检查端口
    var port = checkPort(ip, 80);
    console.log("端口80开放:", port.open);
}
```

### 示例5: CSV数据处理

```javascript
// 读取CSV
var csv = readCSV("data.csv", {
    delimiter: ",",
    skipEmptyLines: true
});

// 处理数据
var processed = [];
csv.rows.forEach(function(row, index) {
    if (index === 0) return; // 跳过表头
    processed.push({
        name: row[0],
        age: parseInt(row[1]),
        city: row[2]
    });
});

// 写入新CSV
var output = [["姓名", "年龄", "城市"]];
processed.forEach(function(item) {
    output.push([item.name, item.age.toString(), item.city]);
});
writeCSV("output.csv", output);
```

### 示例6: 日志记录

```javascript
// 设置日志级别
logger.setLevel("debug");

// 应用程序启动日志
logger.info("应用程序启动");
logger.info({version: "1.0.0", env: "production"}, "配置加载完成");

// 处理请求时的结构化日志
var requestLogger = logger.withFields({
    requestId: "req-12345",
    userId: 123,
    ip: "192.168.1.1"
});

requestLogger.info("收到HTTP请求");
requestLogger.debug("请求参数:", JSON.stringify({name: "test"}));

// 处理过程中的日志
logger.debug("开始处理数据");
logger.info("处理了", 100, "条记录");

// 错误日志
try {
    // 某些操作
    throw new Error("操作失败");
} catch (error) {
    logger.error({error: error.message, stack: error.stack}, "处理失败");
}

// 警告日志
var memoryUsage = 85;
if (memoryUsage > 80) {
    logger.warn({memory: memoryUsage + "%", threshold: "80%"}, "内存使用率较高");
}

// 检查日志级别
var levelInfo = logger.getLevel();
console.log("当前日志级别:", levelInfo.level);

// 条件日志输出
if (logger.isLevelEnabled("debug").enabled) {
    logger.debug("详细的调试信息");
}
```

---

## 注意事项

1. **浏览器会话必须关闭**: 使用 `createBrowserSession` 创建的会话必须调用 `close()` 方法，否则可能导致资源泄漏。

2. **文件路径**: 使用绝对路径或相对于当前工作目录的相对路径。

3. **错误处理**: 始终检查返回值的 `success` 或 `error` 字段。

4. **超时设置**: 对于长时间运行的操作（如浏览器自动化、大文件处理），建议设置合适的超时时间。

5. **异步操作**: 虽然JavaScript代码是同步执行的，但底层操作（如HTTP请求、浏览器操作）是异步的，函数会等待操作完成。

6. **资源限制**: 注意文件大小限制和内存使用，大文件处理时考虑使用分页读取。

---

## 版本信息

本文档基于 jssandbox-go 项目生成，包含所有可用的JavaScript函数和API。

最后更新: 2024年

