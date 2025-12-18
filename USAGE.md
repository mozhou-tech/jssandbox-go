# JavaScript沙盒使用说明

## 安装

### 1. 安装 Go 依赖

```bash
go mod download
```

### 2. 安装系统依赖（如需要）

#### Chrome/Chromium（浏览器功能需要）

浏览器自动化功能（`browserNavigate`, `browserScreenshot` 等）需要 Chrome/Chromium。chromedp 会自动下载 Chrome，也可以使用系统已安装的 Chrome。

## 作为库使用

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
    
    // 执行JavaScript代码
    result, err := sandbox.Run(`
        var date = getCurrentDate();
        var time = getCurrentTime();
        console.log("日期:", date, "时间:", time);
    `)
}
```

## 运行示例程序

```bash
go run cmd/jssandbox/main.go
```

## 功能列表

### 1. 系统操作
- `getCurrentTime()` - 获取当前时间
- `getCurrentDate()` - 获取当前日期
- `getCurrentDateTime()` - 获取当前日期时间
- `getCPUNum()` - 获取CPU数量
- `getMemorySize()` - 获取内存信息
- `getDiskSize(path)` - 获取磁盘空间信息
- `sleep(ms)` - 休眠指定毫秒数

### 2. HTTP请求
- `httpRequest(url, options)` - 通用HTTP请求
- `httpGet(url)` - GET请求
- `httpPost(url, body)` - POST请求

### 3. 文件系统操作
- `openFile(path)` - 使用系统默认程序打开文件
- `getFileInfo(path)` - 获取文件元信息
- `renameFile(oldPath, newPath)` - 重命名文件
- `readFile(path, options)` - 读取文件内容（支持分页）
- `readFileHead(path, lines)` - 读取文件前几行
- `readFileTail(path, lines)` - 读取文件后几行
- `getFileHash(path, type)` - 获取文件哈希值（md5/sha1/sha256/sha512）
- `readImageBase64(path)` - 读取图片的base64编码
- `writeFile(path, content)` - 写入文件
- `appendFile(path, content)` - 追加文件

### 4. 文档读取
- `readWord(path, options)` - 读取Word文档（支持分页）
- `readExcel(path, options)` - 读取Excel文件（支持分页）
- `readPPT(path, options)` - 读取PPT文件（支持分页）
- `readPDF(path, options)` - 读取PDF文件（支持分页）

### 5. 浏览器操作
- `browserNavigate(url)` - 导航到URL并获取HTML
- `browserScreenshot(url, outputPath)` - 截图
- `browserEvaluate(url, jsCode)` - 在页面中执行JavaScript
- `browserClick(url, selector)` - 点击元素
- `browserFill(url, selector, value)` - 填充表单

### 6. 图片处理
- `imageResize(inputPath, outputPath, width, height?)` - 调整图片大小
- `imageCrop(inputPath, outputPath, x, y, width, height)` - 裁剪图片
- `imageRotate(inputPath, outputPath, angle)` - 旋转图片
- `imageFlip(inputPath, outputPath, direction)` - 翻转图片（horizontal/vertical）
- `imageInfo(filePath)` - 获取图片信息（宽度、高度、格式）
- `imageConvert(inputPath, outputPath)` - 转换图片格式
- `imageQuality(inputPath, outputPath, quality)` - 调整图片质量（1-100，仅JPEG）

### 7. 文件类型检测
- `detectFileType(filePath)` - 检测文件类型（返回MIME、扩展名等）
- `isImage(filePath)` - 检测是否为图片
- `isAudio(filePath)` - 检测是否为音频
- `isDocument(filePath)` - 检测是否为文档
- `isFont(filePath)` - 检测是否为字体
- `isArchive(filePath)` - 检测是否为归档文件

## 示例代码

### 系统信息查询
```javascript
var cpuNum = getCPUNum();
var mem = getMemorySize();
console.log("CPU数量:", cpuNum);
console.log("总内存:", mem.totalStr);
console.log("已用内存:", mem.usedStr);
```

### HTTP请求
```javascript
var response = httpGet("https://api.example.com/data");
console.log("状态码:", response.status);
console.log("响应体:", response.body);
```

### 文件操作
```javascript
// 写入文件
writeFile("test.txt", "Hello World");

// 读取文件
var content = readFile("test.txt");
console.log(content.data);

// 获取文件信息
var info = getFileInfo("test.txt");
console.log("文件大小:", info.size);
console.log("修改时间:", info.modTime);

// 获取文件哈希
var hash = getFileHash("test.txt", "md5");
console.log("MD5:", hash.hash);
```

### 文档读取
```javascript
// 读取Word文档第一页
var word = readWord("document.docx", {page: 1, pageSize: 1000});
console.log("总页数:", word.totalPages);
console.log("内容:", word.text);

// 读取Excel第一页
var excel = readExcel("data.xlsx", {page: 1, pageSize: 100});
console.log("总行数:", excel.totalRows);
console.log("数据:", excel.rows);
```

### 浏览器自动化
```javascript
// 导航并获取HTML
var result = browserNavigate("https://www.example.com");
if (result.success) {
    console.log("HTML长度:", result.html.length);
}

// 截图
var screenshot = browserScreenshot("https://www.example.com", "screenshot.png");
console.log("截图保存到:", screenshot.path);
```

### 图片处理
```javascript
// 调整图片大小
var result = imageResize("input.jpg", "output.jpg", 800, 600);
console.log("调整大小成功:", result.success);

// 获取图片信息
var info = imageInfo("photo.png");
console.log("图片尺寸:", info.width, "x", info.height);

// 旋转图片
imageRotate("input.jpg", "output.jpg", 90);

// 翻转图片
imageFlip("input.jpg", "output.jpg", "horizontal");

// 调整JPEG质量
imageQuality("input.jpg", "output.jpg", 85);
```

### 文件类型检测
```javascript
// 检测文件类型
var type = detectFileType("file.bin");
if (!type.unknown) {
    console.log("MIME类型:", type.mime);
    console.log("扩展名:", type.extension);
}

// 检测是否为图片
var isImg = isImage("photo.jpg");
console.log("是图片:", isImg.isImage);
```

