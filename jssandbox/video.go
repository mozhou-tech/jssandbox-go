package jssandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/dop251/goja"
	"go.uber.org/zap"
)

// registerVideoProcessing 注册视频处理功能到JavaScript运行时
func (sb *Sandbox) registerVideoProcessing() {
	// 视频转码/转换格式
	sb.vm.Set("videoConvert", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "需要提供输入路径和输出路径参数",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()

		// 检查输入文件是否存在
		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			sb.logger.Error("输入文件不存在", zap.String("path", inputPath))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("输入文件不存在: %s", inputPath),
			})
		}

		// 确保输出目录存在
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			sb.logger.Error("创建输出目录失败", zap.String("path", outputDir), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		err := ffmpeg.Input(inputPath).
			Output(outputPath).
			OverWriteOutput().
			Run()

		if err != nil {
			sb.logger.Error("视频转码失败", zap.String("input", inputPath), zap.String("output", outputPath), zap.Error(err))
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

	// 视频裁剪（按时间）
	sb.vm.Set("videoTrim", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 4 {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "需要提供输入路径、输出路径、开始时间和持续时间参数",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()
		startTime := call.Arguments[2].String() // 格式: "00:00:10" 或 "10"
		duration := call.Arguments[3].String()  // 格式: "00:00:05" 或 "5"

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("输入文件不存在: %s", inputPath),
			})
		}

		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		err := ffmpeg.Input(inputPath, ffmpeg.KwArgs{"ss": startTime}).
			Output(outputPath, ffmpeg.KwArgs{"t": duration}).
			OverWriteOutput().
			Run()

		if err != nil {
			sb.logger.Error("视频裁剪失败", zap.Error(err))
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

	// 视频裁剪（按尺寸和位置）
	sb.vm.Set("videoCrop", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 6 {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "需要提供输入路径、输出路径、x、y、width、height参数",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()
		x := int(call.Arguments[2].ToInteger())
		y := int(call.Arguments[3].ToInteger())
		width := int(call.Arguments[4].ToInteger())
		height := int(call.Arguments[5].ToInteger())

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("输入文件不存在: %s", inputPath),
			})
		}

		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		// 使用 crop 滤镜: crop=width:height:x:y
		cropFilter := fmt.Sprintf("crop=%d:%d:%d:%d", width, height, x, y)
		err := ffmpeg.Input(inputPath).
			Output(outputPath, ffmpeg.KwArgs{"vf": cropFilter}).
			OverWriteOutput().
			Run()

		if err != nil {
			sb.logger.Error("视频裁剪失败", zap.Error(err))
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

	// 调整视频分辨率
	sb.vm.Set("videoResize", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 3 {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "需要提供输入路径、输出路径和宽度参数",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()
		width := int(call.Arguments[2].ToInteger())
		height := 0

		if len(call.Arguments) > 3 {
			height = int(call.Arguments[3].ToInteger())
		}

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("输入文件不存在: %s", inputPath),
			})
		}

		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		var scaleFilter string
		if height > 0 {
			scaleFilter = fmt.Sprintf("scale=%d:%d", width, height)
		} else {
			scaleFilter = fmt.Sprintf("scale=%d:-1", width)
		}

		err := ffmpeg.Input(inputPath).
			Output(outputPath, ffmpeg.KwArgs{"vf": scaleFilter}).
			OverWriteOutput().
			Run()

		if err != nil {
			sb.logger.Error("调整视频分辨率失败", zap.Error(err))
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

	// 提取音频
	sb.vm.Set("videoExtractAudio", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "需要提供输入视频路径和输出音频路径参数",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()
		audioCodec := "libmp3lame" // 默认使用 MP3
		bitrate := "192k"           // 默认比特率

		if len(call.Arguments) > 2 {
			options := call.Arguments[2].ToObject(sb.vm)
			if codecVal := options.Get("codec"); codecVal != nil && !goja.IsUndefined(codecVal) {
				audioCodec = codecVal.String()
			}
			if bitrateVal := options.Get("bitrate"); bitrateVal != nil && !goja.IsUndefined(bitrateVal) {
				bitrate = bitrateVal.String()
			}
		}

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("输入文件不存在: %s", inputPath),
			})
		}

		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		err := ffmpeg.Input(inputPath).
			Output(outputPath, ffmpeg.KwArgs{
				"vn":      "",           // 不包含视频
				"acodec":  audioCodec,   // 音频编解码器
				"ab":      bitrate,      // 音频比特率
				"ar":      "44100",      // 采样率
				"ac":      "2",          // 声道数
			}).
			OverWriteOutput().
			Run()

		if err != nil {
			sb.logger.Error("提取音频失败", zap.Error(err))
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

	// 合并视频
	sb.vm.Set("videoConcat", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "需要提供视频文件路径数组和输出路径参数",
			})
		}

		// 获取视频文件列表
		videoList := call.Arguments[0].ToObject(sb.vm)
		videoArray := videoList.Get("length").ToInteger()
		if videoArray == 0 {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "视频文件列表不能为空",
			})
		}

		var inputPaths []string
		for i := 0; i < int(videoArray); i++ {
			pathVal := videoList.Get(strconv.Itoa(i))
			if pathVal == nil || goja.IsUndefined(pathVal) {
				continue
			}
			path := pathVal.String()
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return sb.vm.ToValue(map[string]interface{}{
					"success": false,
					"error":   fmt.Sprintf("视频文件不存在: %s", path),
				})
			}
			inputPaths = append(inputPaths, path)
		}

		if len(inputPaths) == 0 {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "没有有效的视频文件",
			})
		}

		outputPath := call.Arguments[1].String()
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		// 创建临时文件列表
		listFile, err := os.CreateTemp("", "ffmpeg-concat-*.txt")
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}
		defer os.Remove(listFile.Name())
		defer listFile.Close()

		// 写入文件列表（使用绝对路径）
		for _, path := range inputPaths {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return sb.vm.ToValue(map[string]interface{}{
					"success": false,
					"error":   err.Error(),
				})
			}
			// FFmpeg concat 格式: file 'path'
			fmt.Fprintf(listFile, "file '%s'\n", strings.ReplaceAll(absPath, "'", "'\\''"))
		}
		listFile.Close()

		// 使用 concat demuxer
		err = ffmpeg.Input(listFile.Name(), ffmpeg.KwArgs{"f": "concat", "safe": "0"}).
			Output(outputPath, ffmpeg.KwArgs{"c": "copy"}).
			OverWriteOutput().
			Run()

		if err != nil {
			sb.logger.Error("合并视频失败", zap.Error(err))
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

	// 压缩视频
	sb.vm.Set("videoCompress", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "需要提供输入路径和输出路径参数",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()
		crf := 23 // 默认 CRF 值（18-28，值越大质量越低文件越小）
		preset := "medium" // 编码速度预设

		if len(call.Arguments) > 2 {
			options := call.Arguments[2].ToObject(sb.vm)
			if crfVal := options.Get("crf"); crfVal != nil && !goja.IsUndefined(crfVal) {
				crf = int(crfVal.ToInteger())
			}
			if presetVal := options.Get("preset"); presetVal != nil && !goja.IsUndefined(presetVal) {
				preset = presetVal.String()
			}
		}

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("输入文件不存在: %s", inputPath),
			})
		}

		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		err := ffmpeg.Input(inputPath).
			Output(outputPath, ffmpeg.KwArgs{
				"c:v":    "libx264",     // 视频编解码器
				"crf":    strconv.Itoa(crf), // 质量参数
				"preset": preset,        // 编码速度
				"c:a":    "aac",         // 音频编解码器
				"b:a":    "128k",        // 音频比特率
			}).
			OverWriteOutput().
			Run()

		if err != nil {
			sb.logger.Error("压缩视频失败", zap.Error(err))
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

	// 获取视频信息
	sb.vm.Set("videoInfo", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供视频文件路径",
			})
		}

		filePath := call.Arguments[0].String()

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("文件不存在: %s", filePath),
			})
		}

		// 使用 ffprobe 获取视频信息
		// 注意: ffmpeg-go 不直接提供 probe 功能，这里返回基本信息
		// 实际项目中可能需要使用 ffprobe 命令或相关库
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		// 返回基本信息，实际视频信息需要调用 ffprobe
		return sb.vm.ToValue(map[string]interface{}{
			"path":     filePath,
			"size":     fileInfo.Size(),
			"modified": fileInfo.ModTime().Format("2006-01-02 15:04:05"),
			"note":     "详细视频信息（分辨率、时长、编码等）需要使用 ffprobe 获取",
		})
	})

	// 添加水印
	sb.vm.Set("videoWatermark", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 3 {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "需要提供输入路径、输出路径和水印图片路径参数",
			})
		}

		inputPath := call.Arguments[0].String()
		outputPath := call.Arguments[1].String()
		watermarkPath := call.Arguments[2].String()
		position := "10:10" // 默认位置（左上角）
		scale := "100:100"  // 默认水印大小

		if len(call.Arguments) > 3 {
			options := call.Arguments[3].ToObject(sb.vm)
			if posVal := options.Get("position"); posVal != nil && !goja.IsUndefined(posVal) {
				position = posVal.String()
			}
			if scaleVal := options.Get("scale"); scaleVal != nil && !goja.IsUndefined(scaleVal) {
				scale = scaleVal.String()
			}
		}

		if _, err := os.Stat(inputPath); os.IsNotExist(err) {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("输入文件不存在: %s", inputPath),
			})
		}

		if _, err := os.Stat(watermarkPath); os.IsNotExist(err) {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("水印文件不存在: %s", watermarkPath),
			})
		}

		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		// 使用 overlay 滤镜添加水印
		// 构建 filter_complex 参数
		filterComplex := fmt.Sprintf("[1:v]scale=%s[wm];[0:v][wm]overlay=%s", scale, position)

		// 对于多个输入，使用 Filter 函数
		// Filter 函数签名：Filter([]*Stream, filter string, args Args, kwargs ...KwArgs)
		// 注意：根据 API，Filter 的第二个参数是 filter 字符串，第三个是 Args，第四个是可变 KwArgs
		streams := []*ffmpeg.Stream{
			ffmpeg.Input(inputPath),
			ffmpeg.Input(watermarkPath),
		}
		err := ffmpeg.Filter(
			streams,
			filterComplex,
			ffmpeg.Args{},
			ffmpeg.KwArgs{"map": "[v]"},
		).Output(outputPath).OverWriteOutput().Run()

		if err != nil {
			sb.logger.Error("添加水印失败", zap.Error(err))
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
}

