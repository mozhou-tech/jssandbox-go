# JavaScript沙盒使用说明

## 安装

### 1. 安装 Go 依赖

```bash
go mod download
```

### 2. 安装系统依赖（如需要）

#### ffmpeg（视频处理功能需要）

如果使用视频处理功能（`videoConvert`, `videoTrim`, `videoCrop` 等），需要安装 ffmpeg：

- **macOS:**
  ```bash
  brew install ffmpeg
  ```

- **Linux (Ubuntu/Debian):**
  ```bash
  sudo apt update && sudo apt install ffmpeg
  ```

- **Linux (CentOS/RHEL):**
  ```bash
  sudo yum install ffmpeg
  # 或
  sudo dnf install ffmpeg
  ```

- **Windows:**
  - 从 [ffmpeg.org](https://ffmpeg.org/download.html) 下载
  - 或使用 Chocolatey: `choco install ffmpeg`
  - 或使用 Scoop: `scoop install ffmpeg`

验证安装：
```bash
ffmpeg -version
```

> **提示：** 如果不需要视频处理功能，可以在创建沙盒时禁用：
> ```go
> config := jssandbox.DefaultConfig().DisableVideoProcessing()
> sandbox := jssandbox.NewSandboxWithConfig(ctx, config)
> ```

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
- `isVideo(filePath)` - 检测是否为视频
- `isAudio(filePath)` - 检测是否为音频
- `isDocument(filePath)` - 检测是否为文档
- `isFont(filePath)` - 检测是否为字体
- `isArchive(filePath)` - 检测是否为归档文件

### 8. 视频处理
- `videoConvert(inputPath, outputPath)` - 视频转码/转换格式
- `videoTrim(inputPath, outputPath, startTime, duration)` - 视频裁剪（按时间）
  - startTime: 开始时间，格式 "00:00:10" 或 "10"（秒）
  - duration: 持续时间，格式 "00:00:05" 或 "5"（秒）
- `videoCrop(inputPath, outputPath, x, y, width, height)` - 视频裁剪（按尺寸和位置）
- `videoResize(inputPath, outputPath, width, height?)` - 调整视频分辨率
  - height: 可选，如果不提供则按比例缩放
- `videoExtractAudio(inputPath, outputPath, options?)` - 提取音频
  - options: 可选对象，包含 codec（编解码器，默认 "libmp3lame"）和 bitrate（比特率，默认 "192k"）
- `videoConcat(videoPaths, outputPath)` - 合并视频
  - videoPaths: 视频文件路径数组
- `videoCompress(inputPath, outputPath, options?)` - 压缩视频
  - options: 可选对象，包含 crf（质量参数，18-28，默认23）和 preset（编码速度，默认 "medium"）
- `videoInfo(filePath)` - 获取视频信息（文件大小、修改时间等）
- `videoWatermark(inputPath, outputPath, watermarkPath, options?)` - 添加水印
  - options: 可选对象，包含 position（位置，默认 "10:10"）和 scale（水印大小，默认 "100:100"）

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

// 检测是否为视频
var isVid = isVideo("video.mp4");
console.log("是视频:", isVid.isVideo);
```

### 视频处理
```javascript
// 视频转码
var result = videoConvert("input.avi", "output.mp4");
console.log("转码成功:", result.success);

// 视频裁剪（按时间）
videoTrim("input.mp4", "output.mp4", "00:00:10", "00:00:30");

// 视频裁剪（按尺寸）
videoCrop("input.mp4", "output.mp4", 100, 100, 640, 480);

// 调整视频分辨率
videoResize("input.mp4", "output.mp4", 1280, 720);

// 提取音频
var audio = videoExtractAudio("video.mp4", "audio.mp3", {
    codec: "libmp3lame",
    bitrate: "192k"
});

// 合并视频
var videos = ["video1.mp4", "video2.mp4", "video3.mp4"];
videoConcat(videos, "merged.mp4");

// 压缩视频
videoCompress("input.mp4", "output.mp4", {
    crf: 23,
    preset: "medium"
});

// 获取视频信息
var info = videoInfo("video.mp4");
console.log("文件大小:", info.size);

// 添加水印
videoWatermark("input.mp4", "output.mp4", "watermark.png", {
    position: "10:10",
    scale: "200:200"
});
```

