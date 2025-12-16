package jssandbox

import (
	"image"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/dop251/goja"
	"go.uber.org/zap"
)

// registerImageProcessing 注册图片处理功能到JavaScript运行时
func (sb *Sandbox) registerImageProcessing() {
	// 调整图片大小
	sb.vm.Set("imageResize", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 3 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供输入路径、输出路径和宽度参数",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()
		width := int(call.Arguments[2].ToInteger())
		height := 0

		if len(call.Arguments) > 3 {
			height = int(call.Arguments[3].ToInteger())
		}

		img, err := imaging.Open(inputPath)
		if err != nil {
			sb.logger.Error("打开图片失败", zap.String("path", inputPath), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		var resized image.Image
		if height > 0 {
			resized = imaging.Resize(img, width, height, imaging.Lanczos)
		} else {
			resized = imaging.Resize(img, width, 0, imaging.Lanczos)
		}

		err = imaging.Save(resized, outputPath)
		if err != nil {
			sb.logger.Error("保存图片失败", zap.String("path", outputPath), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    outputPath,
		})
	})

	// 裁剪图片
	sb.vm.Set("imageCrop", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 6 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供输入路径、输出路径、x、y、width、height参数",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()
		x := int(call.Arguments[2].ToInteger())
		y := int(call.Arguments[3].ToInteger())
		width := int(call.Arguments[4].ToInteger())
		height := int(call.Arguments[5].ToInteger())

		img, err := imaging.Open(inputPath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		cropped := imaging.Crop(img, image.Rect(x, y, x+width, y+height))
		err = imaging.Save(cropped, outputPath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    outputPath,
		})
	})

	// 旋转图片
	sb.vm.Set("imageRotate", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 3 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供输入路径、输出路径和角度参数",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()
		angle := call.Arguments[2].ToFloat()

		img, err := imaging.Open(inputPath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		rotated := imaging.Rotate(img, angle, nil)
		err = imaging.Save(rotated, outputPath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    outputPath,
		})
	})

	// 翻转图片
	sb.vm.Set("imageFlip", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 3 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供输入路径、输出路径和方向参数（horizontal/vertical）",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()
		direction := call.Arguments[2].String()

		img, err := imaging.Open(inputPath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		var flipped image.Image
		if direction == "horizontal" {
			flipped = imaging.FlipH(img)
		} else if direction == "vertical" {
			flipped = imaging.FlipV(img)
		} else {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "方向参数必须是horizontal或vertical",
			})
		}

		err = imaging.Save(flipped, outputPath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    outputPath,
		})
	})

	// 获取图片信息
	sb.vm.Set("imageInfo", func(filePath string) goja.Value {
		img, err := imaging.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		return sb.vm.ToValue(map[string]interface{}{
			"width":  width,
			"height": height,
			"format": getImageFormat(filePath),
		})
	})

	// 转换图片格式
	sb.vm.Set("imageConvert", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 3 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供输入路径、输出路径参数",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()

		img, err := imaging.Open(inputPath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		err = imaging.Save(img, outputPath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    outputPath,
		})
	})

	// 调整图片质量（JPEG）
	sb.vm.Set("imageQuality", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 3 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供输入路径、输出路径和质量参数（1-100）",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()
		quality := int(call.Arguments[2].ToInteger())

		if quality < 1 || quality > 100 {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "质量参数必须在1-100之间",
			})
		}

		img, err := imaging.Open(inputPath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		// 根据输出文件扩展名选择编码选项
		ext := filepath.Ext(outputPath)
		var err2 error
		if ext == ".jpg" || ext == ".jpeg" {
			err2 = imaging.Save(img, outputPath, imaging.JPEGQuality(quality))
		} else {
			err2 = imaging.Save(img, outputPath)
		}

		if err2 != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err2.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    outputPath,
		})
	})
}

// getImageFormat 根据文件扩展名获取图片格式
func getImageFormat(filePath string) string {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".png":
		return "png"
	case ".gif":
		return "gif"
	case ".bmp":
		return "bmp"
	case ".webp":
		return "webp"
	default:
		return "unknown"
	}
}
