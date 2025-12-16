package jssandbox

import (
	"os"

	"github.com/dop251/goja"
	"github.com/h2non/filetype"
	"go.uber.org/zap"
)

// registerFileTypeDetection 注册文件类型检测功能到JavaScript运行时
func (sb *Sandbox) registerFileTypeDetection() {
	// 检测文件类型
	sb.vm.Set("detectFileType", func(filePath string) goja.Value {
		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer file.Close()

		// 读取文件头部（前261字节用于类型检测）
		buf := make([]byte, 261)
		n, err := file.Read(buf)
		if err != nil && n == 0 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		// 检测文件类型
		kind, err := filetype.Match(buf[:n])
		if err != nil {
			sb.logger.Error("文件类型检测失败", zap.String("path", filePath), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"unknown": true,
				"error":   err.Error(),
			})
		}

		if kind == filetype.Unknown {
			return sb.vm.ToValue(map[string]interface{}{
				"unknown": true,
				"message": "无法识别的文件类型",
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"unknown":   false,
			"mime":      kind.MIME.Value,
			"extension": kind.Extension,
			"type":      kind.MIME.Type,
			"subtype":   kind.MIME.Subtype,
		})
	})

	// 检测是否为图片
	sb.vm.Set("isImage", func(filePath string) goja.Value {
		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"isImage": false,
				"error":   err.Error(),
			})
		}
		defer file.Close()

		buf := make([]byte, 261)
		n, err := file.Read(buf)
		if err != nil && n == 0 {
			return sb.vm.ToValue(map[string]interface{}{
				"isImage": false,
				"error":   err.Error(),
			})
		}

		isImage := filetype.IsImage(buf[:n])
		return sb.vm.ToValue(map[string]interface{}{
			"isImage": isImage,
		})
	})

	// 检测是否为视频
	sb.vm.Set("isVideo", func(filePath string) goja.Value {
		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"isVideo": false,
				"error":   err.Error(),
			})
		}
		defer file.Close()

		buf := make([]byte, 261)
		n, err := file.Read(buf)
		if err != nil && n == 0 {
			return sb.vm.ToValue(map[string]interface{}{
				"isVideo": false,
				"error":   err.Error(),
			})
		}

		isVideo := filetype.IsVideo(buf[:n])
		return sb.vm.ToValue(map[string]interface{}{
			"isVideo": isVideo,
		})
	})

	// 检测是否为音频
	sb.vm.Set("isAudio", func(filePath string) goja.Value {
		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"isAudio": false,
				"error":   err.Error(),
			})
		}
		defer file.Close()

		buf := make([]byte, 261)
		n, err := file.Read(buf)
		if err != nil && n == 0 {
			return sb.vm.ToValue(map[string]interface{}{
				"isAudio": false,
				"error":   err.Error(),
			})
		}

		isAudio := filetype.IsAudio(buf[:n])
		return sb.vm.ToValue(map[string]interface{}{
			"isAudio": isAudio,
		})
	})

	// 检测是否为文档
	sb.vm.Set("isDocument", func(filePath string) goja.Value {
		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"isDocument": false,
				"error":      err.Error(),
			})
		}
		defer file.Close()

		buf := make([]byte, 261)
		n, err := file.Read(buf)
		if err != nil && n == 0 {
			return sb.vm.ToValue(map[string]interface{}{
				"isDocument": false,
				"error":      err.Error(),
			})
		}

		isDocument := filetype.IsDocument(buf[:n])
		return sb.vm.ToValue(map[string]interface{}{
			"isDocument": isDocument,
		})
	})

	// 检测是否为字体
	sb.vm.Set("isFont", func(filePath string) goja.Value {
		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"isFont": false,
				"error":  err.Error(),
			})
		}
		defer file.Close()

		buf := make([]byte, 261)
		n, err := file.Read(buf)
		if err != nil && n == 0 {
			return sb.vm.ToValue(map[string]interface{}{
				"isFont": false,
				"error":  err.Error(),
			})
		}

		isFont := filetype.IsFont(buf[:n])
		return sb.vm.ToValue(map[string]interface{}{
			"isFont": isFont,
		})
	})

	// 检测是否为归档文件
	sb.vm.Set("isArchive", func(filePath string) goja.Value {
		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"isArchive": false,
				"error":     err.Error(),
			})
		}
		defer file.Close()

		buf := make([]byte, 261)
		n, err := file.Read(buf)
		if err != nil && n == 0 {
			return sb.vm.ToValue(map[string]interface{}{
				"isArchive": false,
				"error":     err.Error(),
			})
		}

		isArchive := filetype.IsArchive(buf[:n])
		return sb.vm.ToValue(map[string]interface{}{
			"isArchive": isArchive,
		})
	})
}

