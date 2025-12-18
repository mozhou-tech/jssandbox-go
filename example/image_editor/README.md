# 图片编辑示例

这个示例展示了如何使用 jssandbox 进行各种图片处理操作。

## 功能说明

该示例演示了以下图片处理功能：

1. **获取图片信息** - 获取图片的宽度、高度和格式
2. **调整图片大小** - 缩放图片到指定尺寸（支持保持宽高比）
3. **裁剪图片** - 从指定位置裁剪指定大小的区域
4. **旋转图片** - 按指定角度旋转图片
5. **翻转图片** - 水平或垂直翻转图片
6. **转换图片格式** - 将图片转换为不同格式（PNG、JPEG等）
7. **调整JPEG质量** - 调整JPEG图片的压缩质量（1-100）
8. **组合操作** - 演示多个操作的组合使用

## 使用方法

### 基本用法

```bash
# 执行所有操作
go run main.go -input photo.jpg

# 指定输出目录
go run main.go -input photo.jpg -output my_output

# 只执行特定操作
go run main.go -input photo.jpg -ops resize,rotate,flip

# 只查看图片信息
go run main.go -input photo.jpg -ops info
```

### 命令行参数

- `-input` (必需): 输入图片文件路径
- `-output` (可选): 输出目录，默认为 `output`
- `-ops` (可选): 要执行的操作，可选值：
  - `all` - 执行所有操作（默认）
  - `resize` - 调整大小
  - `crop` - 裁剪
  - `rotate` - 旋转
  - `flip` - 翻转
  - `convert` - 格式转换
  - `quality` - 调整质量（仅JPEG）
  - `info` - 获取图片信息
  - 可以组合多个操作，用逗号分隔，如：`resize,rotate,flip`

### 示例

```bash
# 示例1: 处理一张照片，执行所有操作
go run main.go -input ~/Pictures/photo.jpg

# 示例2: 只调整大小和旋转
go run main.go -input photo.png -ops resize,rotate

# 示例3: 查看图片信息
go run main.go -input image.jpg -ops info

# 示例4: 自定义输出目录
go run main.go -input photo.jpg -output processed_images
```

## 输出文件说明

根据执行的操作，会在输出目录中生成以下文件：

- `{原文件名}_resized.{扩展名}` - 调整大小后的图片（宽度800，保持宽高比）
- `{原文件名}_resized_600x400.{扩展名}` - 调整大小后的图片（600x400）
- `{原文件名}_cropped.{扩展名}` - 裁剪后的图片
- `{原文件名}_rotated_90.{扩展名}` - 旋转90度后的图片
- `{原文件名}_rotated_180.{扩展名}` - 旋转180度后的图片
- `{原文件名}_flipped_h.{扩展名}` - 水平翻转后的图片
- `{原文件名}_flipped_v.{扩展名}` - 垂直翻转后的图片
- `{原文件名}_converted.png` - 转换为PNG格式
- `{原文件名}_converted.jpg` - 转换为JPEG格式
- `{原文件名}_quality_95.jpg` - 高质量JPEG（质量95）
- `{原文件名}_quality_75.jpg` - 中等质量JPEG（质量75）
- `{原文件名}_quality_50.jpg` - 低质量JPEG（质量50）
- `{原文件名}_combo_final.{扩展名}` - 组合操作后的最终图片

## 支持的图片格式

- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- BMP (.bmp)
- WebP (.webp)

## 注意事项

1. **文件路径**: 可以使用相对路径或绝对路径
2. **输出目录**: 如果输出目录不存在，会自动创建
3. **JPEG质量调整**: 仅对JPEG格式的图片有效
4. **裁剪区域**: 确保裁剪区域不超过原图尺寸
5. **文件覆盖**: 如果输出文件已存在，会被覆盖

## 代码说明

### 主要功能

1. **imageInfo(filePath)** - 获取图片信息
   ```javascript
   var info = imageInfo("photo.jpg");
   console.log(info.width, info.height, info.format);
   ```

2. **imageResize(inputPath, outputPath, width, height?)** - 调整大小
   ```javascript
   // 宽度800，保持宽高比
   imageResize("input.jpg", "output.jpg", 800);
   // 指定宽度和高度
   imageResize("input.jpg", "output.jpg", 800, 600);
   ```

3. **imageCrop(inputPath, outputPath, x, y, width, height)** - 裁剪
   ```javascript
   // 从(0,0)位置裁剪400x300的区域
   imageCrop("input.jpg", "output.jpg", 0, 0, 400, 300);
   ```

4. **imageRotate(inputPath, outputPath, angle)** - 旋转
   ```javascript
   // 旋转90度
   imageRotate("input.jpg", "output.jpg", 90);
   ```

5. **imageFlip(inputPath, outputPath, direction)** - 翻转
   ```javascript
   // 水平翻转
   imageFlip("input.jpg", "output.jpg", "horizontal");
   // 垂直翻转
   imageFlip("input.jpg", "output.jpg", "vertical");
   ```

6. **imageConvert(inputPath, outputPath)** - 格式转换
   ```javascript
   // 转换为PNG
   imageConvert("input.jpg", "output.png");
   ```

7. **imageQuality(inputPath, outputPath, quality)** - 调整质量
   ```javascript
   // 设置JPEG质量为85
   imageQuality("input.jpg", "output.jpg", 85);
   ```

## 自定义

你可以修改代码中的参数来适应不同的需求：

- **调整大小参数**: 修改 `imageResize` 的宽度和高度参数
- **裁剪区域**: 修改 `imageCrop` 的 x, y, width, height 参数
- **旋转角度**: 修改 `imageRotate` 的角度参数
- **质量设置**: 修改 `imageQuality` 的质量值（1-100）

## 依赖

- Go 1.21+
- jssandbox-go 库
- 无需额外依赖（图片处理功能已内置）

