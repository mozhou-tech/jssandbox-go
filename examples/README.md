# 示例程序

本目录包含多个使用 jssandbox 的示例程序，展示了各种功能的使用方法。

## 目录结构

```
example/
├── README.md           # 本文件
├── crawler/            # 网页爬虫示例
│   ├── main.go
│   └── README.md
└── image_editor/        # 图片编辑示例
    ├── main.go
    └── README.md
```

## 示例列表

### 1. 网页爬虫示例 (`crawler/`)

使用 jssandbox 的浏览器自动化功能实现网页爬虫，爬取江苏省公共资源交易平台的招标信息。

**快速开始：**
```bash
cd crawler
go run main.go
```

**主要功能：**
- 浏览器自动化导航
- 页面数据提取
- JSON 数据保存
- 数据预览

详细说明请参考 [crawler/README.md](crawler/README.md)

### 2. 图片编辑示例 (`image_editor/`)

使用 jssandbox 的图片处理功能进行各种图片编辑操作。

**快速开始：**
```bash
cd image_editor
go run main.go -input photo.jpg
```

**主要功能：**
- 调整图片大小
- 裁剪图片
- 旋转和翻转
- 格式转换
- 质量调整

详细说明请参考 [image_editor/README.md](image_editor/README.md)

## 运行示例

### 前置要求

- Go 1.21+
- jssandbox-go 库（已包含在项目中）
- Chrome/Chromium（仅爬虫示例需要，chromedp 会自动下载）

### 运行步骤

1. **进入示例目录**
   ```bash
   cd example/crawler    # 或 example/image_editor
   ```

2. **运行示例**
   ```bash
   go run main.go [参数]
   ```

3. **查看输出**
   - 爬虫示例会生成 JSON 文件
   - 图片编辑示例会在输出目录生成处理后的图片

## 代码结构

每个示例都遵循相同的结构：

- `main.go` - 主程序文件
- `README.md` - 详细的使用说明和文档

## 贡献

如果你有新的示例想法，欢迎：

1. 在 `example/` 目录下创建新的子目录
2. 添加 `main.go` 和 `README.md`
3. 更新本 README 文件

## 注意事项

- 所有示例都使用 `logrus` 进行日志记录
- 确保有足够的权限访问文件系统
- 网络相关示例需要稳定的网络连接
- 浏览器自动化示例需要 Chrome/Chromium 支持
