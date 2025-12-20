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
    "github.com/mozhou-tech/jssandbox-go/jssandbox"
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

### 8. 加密/解密功能
- `encryptAES(data, key)` - AES加密数据
- `decryptAES(encrypted, key)` - AES解密数据
- `hashSHA256(data)` - 计算SHA256哈希值
- `generateUUID()` - 生成UUID
- `generateRandomString(length?)` - 生成随机字符串（默认32字符）

### 9. 压缩/解压缩
- `compressZip(files, outputPath)` - 压缩文件为ZIP
- `extractZip(zipPath, outputDir)` - 解压ZIP文件
- `compressGzip(data)` - GZIP压缩字符串
- `decompressGzip(compressed)` - GZIP解压字符串

### 10. CSV处理
- `readCSV(filePath, options?)` - 读取CSV文件（支持自定义分隔符、注释符等）
- `writeCSV(filePath, data, options?)` - 写入CSV文件
- `parseCSV(csvString, options?)` - 解析CSV字符串

### 11. 环境变量和配置
- `getEnv(name)` - 获取环境变量
- `getEnvAll()` - 获取所有环境变量
- `readConfig(filePath)` - 读取配置文件（支持JSON和YAML）

### 12. 数据验证
- `validateEmail(email)` - 验证邮箱格式
- `validateURL(url)` - 验证URL格式
- `validateIP(ip)` - 验证IP地址（返回是否为IPv4/IPv6）
- `validatePhone(phone)` - 验证中国手机号

### 13. 日期时间增强
- `formatDate(date, format)` - 格式化日期
- `parseDate(dateString)` - 解析日期字符串
- `addDays(date, days)` - 日期加减天数
- `getTimezone()` - 获取当前时区
- `convertTimezone(date, timezone)` - 时区转换

### 14. 编码/解码增强
- `encodeBase64(data)` - Base64编码
- `decodeBase64(encoded)` - Base64解码
- `encodeURL(str)` - URL编码
- `decodeURL(encoded)` - URL解码
- `encodeHTML(str)` - HTML实体编码
- `decodeHTML(encoded)` - HTML实体解码

### 15. 进程管理
- `execCommand(command, options?)` - 执行系统命令
- `listProcesses()` - 列出运行中的进程
- `killProcess(pid)` - 终止进程

### 16. 网络工具
- `resolveDNS(hostname)` - DNS解析
- `ping(host, count?)` - Ping测试（默认4次）
- `checkPort(host, port, timeout?)` - 检查端口是否开放

### 17. 路径处理增强
- `pathJoin(...paths)` - 路径拼接
- `pathDir(path)` - 获取目录
- `pathBase(path)` - 获取文件名
- `pathExt(path)` - 获取扩展名
- `pathAbs(path)` - 获取绝对路径

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

### 加密/解密
```javascript
// AES加密
var encrypted = encryptAES("敏感数据", "密钥");
console.log("加密结果:", encrypted.data);

// AES解密
var decrypted = decryptAES(encrypted.data, "密钥");
console.log("解密结果:", decrypted.data);

// SHA256哈希
var hash = hashSHA256("数据");
console.log("哈希值:", hash.hash);

// 生成UUID
var uuid = generateUUID();
console.log("UUID:", uuid);

// 生成随机字符串
var random = generateRandomString(16);
console.log("随机字符串:", random.data);
```

### 压缩/解压缩
```javascript
// 压缩为ZIP
var result = compressZip(["file1.txt", "file2.txt"], "archive.zip");
console.log("压缩成功:", result.success);

// 解压ZIP
var extracted = extractZip("archive.zip", "./output");
console.log("解压文件数:", extracted.files.length);

// GZIP压缩
var compressed = compressGzip("要压缩的数据");
console.log("压缩后:", compressed.data);

// GZIP解压
var decompressed = decompressGzip(compressed.data);
console.log("解压后:", decompressed.data);
```

### CSV处理
```javascript
// 读取CSV
var csv = readCSV("data.csv", {delimiter: ",", skipEmptyLines: true});
console.log("行数:", csv.count);
console.log("数据:", csv.rows);

// 写入CSV
var data = [
    ["姓名", "年龄", "城市"],
    ["张三", "25", "北京"],
    ["李四", "30", "上海"]
];
writeCSV("output.csv", data);

// 解析CSV字符串
var parsed = parseCSV("a,b,c\n1,2,3", {delimiter: ","});
console.log("解析结果:", parsed.rows);
```

### 环境变量和配置
```javascript
// 获取环境变量
var env = getEnv("PATH");
console.log("PATH:", env.value);

// 获取所有环境变量
var allEnv = getEnvAll();
console.log("环境变量数量:", Object.keys(allEnv.env).length);

// 读取配置文件
var config = readConfig("config.yaml");
console.log("配置:", config.config);
```

### 数据验证
```javascript
// 验证邮箱
var emailValid = validateEmail("test@example.com");
console.log("邮箱有效:", emailValid.valid);

// 验证URL
var urlValid = validateURL("https://www.example.com");
console.log("URL有效:", urlValid.valid);

// 验证IP地址
var ipValid = validateIP("192.168.1.1");
console.log("IP有效:", ipValid.valid);
console.log("是IPv4:", ipValid.isIPv4);

// 验证手机号
var phoneValid = validatePhone("13800138000");
console.log("手机号有效:", phoneValid.valid);
```

### 日期时间增强
```javascript
// 格式化日期
var formatted = formatDate("2024-01-01 12:00:00", "YYYY-MM-DD HH:mm:ss");
console.log("格式化后:", formatted.date);

// 解析日期
var parsed = parseDate("2024-01-01");
console.log("时间戳:", parsed.timestamp);
console.log("年:", parsed.year);

// 日期加减
var newDate = addDays("2024-01-01", 7);
console.log("7天后:", newDate.date);

// 获取时区
var tz = getTimezone();
console.log("时区:", tz.timezone);

// 时区转换
var converted = convertTimezone("2024-01-01 12:00:00", "America/New_York");
console.log("转换后:", converted.date);
```

### 编码/解码
```javascript
// Base64编码
var encoded = encodeBase64("Hello World");
console.log("Base64:", encoded.data);

// Base64解码
var decoded = decodeBase64(encoded.data);
console.log("解码后:", decoded.data);

// URL编码
var urlEncoded = encodeURL("hello world");
console.log("URL编码:", urlEncoded.data);

// HTML编码
var htmlEncoded = encodeHTML("<div>内容</div>");
console.log("HTML编码:", htmlEncoded.data);
```

### 进程管理
```javascript
// 执行命令
var result = execCommand("ls", {timeout: 5});
console.log("输出:", result.output);
console.log("退出码:", result.code);

// 列出进程
var processes = listProcesses();
console.log("进程数:", processes.count);
processes.processes.forEach(function(p) {
    console.log("进程:", p.name, "PID:", p.pid);
});

// 终止进程
killProcess(12345);
```

### 网络工具
```javascript
// DNS解析
var dns = resolveDNS("www.example.com");
console.log("IP地址:", dns.ips);

// Ping测试
var ping = ping("www.example.com", 4);
console.log("发送:", ping.sent);
console.log("接收:", ping.received);
console.log("平均时间:", ping.averageTime, "ms");

// 检查端口
var port = checkPort("localhost", 80, 3);
console.log("端口开放:", port.open);
```

### 路径处理
```javascript
// 路径拼接
var path = pathJoin("/usr", "local", "bin");
console.log("拼接路径:", path.path);

// 获取目录
var dir = pathDir("/usr/local/bin/app");
console.log("目录:", dir.dir);

// 获取文件名
var base = pathBase("/usr/local/bin/app");
console.log("文件名:", base.base);

// 获取扩展名
var ext = pathExt("file.txt");
console.log("扩展名:", ext.ext);

// 获取绝对路径
var abs = pathAbs("./file.txt");
console.log("绝对路径:", abs.path);
```

