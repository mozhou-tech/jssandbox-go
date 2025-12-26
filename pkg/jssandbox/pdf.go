package jssandbox

import (
	"fmt"
	"os"

	"github.com/dop251/goja"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// registerDocuments 注册文档处理（PDF等）功能到JavaScript运行时
func (sb *Sandbox) registerDocuments() {
	// 获取PDF页数
	sb.vm.Set("pdfGetPageCount", func(filePath string) goja.Value {
		n, err := api.PageCountFile(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"pages":   n,
		})
	})

	// 合并PDF
	sb.vm.Set("pdfMerge", func(inFiles []string, outFile string) map[string]interface{} {
		err := api.MergeCreateFile(inFiles, outFile, false, nil)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})

	// 拆分PDF (每页一个文件)
	sb.vm.Set("pdfSplit", func(inFile string, outDir string) map[string]interface{} {
		// 确保输出目录存在
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("创建输出目录失败: %v", err),
			}
		}
		err := api.SplitFile(inFile, outDir, 1, nil)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})

	// 提取指定页面
	// pages: []string, e.g. ["1", "2-5", "8"]
	sb.vm.Set("pdfExtractPages", func(inFile string, outDir string, pages []string) map[string]interface{} {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("创建输出目录失败: %v", err),
			}
		}
		err := api.ExtractPagesFile(inFile, outDir, pages, nil)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})

	// 优化PDF
	sb.vm.Set("pdfOptimize", func(inFile string, outFile string) map[string]interface{} {
		err := api.OptimizeFile(inFile, outFile, nil)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})

	// 验证PDF
	sb.vm.Set("pdfValidate", func(inFile string) map[string]interface{} {
		err := api.ValidateFile(inFile, nil)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})

	// 添加文本水印
	// options: { onTop: true, opacity: 0.5, scale: 0.5, rotation: 45 }
	sb.vm.Set("pdfAddTextWatermark", func(inFile, outFile string, text string, options map[string]interface{}) map[string]interface{} {
		wm, err := pdfcpu.ParseTextWatermarkDetails(text, "", true, types.POINTS)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("解析水印设置失败: %v", err),
			}
		}

		// 自定义设置
		if onTop, ok := options["onTop"].(bool); ok {
			wm.OnTop = onTop
		}
		if opacity, ok := options["opacity"].(float64); ok {
			wm.Opacity = opacity
		}
		if scale, ok := options["scale"].(float64); ok {
			wm.Scale = scale
		}
		if rotation, ok := options["rotation"].(float64); ok {
			wm.Rotation = rotation
		}

		err = api.AddWatermarksFile(inFile, outFile, nil, wm, nil)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})

	// 导出图片
	sb.vm.Set("pdfExportImages", func(inFile string, outDir string) map[string]interface{} {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("创建输出目录失败: %v", err),
			}
		}
		err := api.ExtractImagesFile(inFile, outDir, nil, nil)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})

	// 将图片导入为PDF
	sb.vm.Set("pdfImportImages", func(imgFiles []string, outFile string) map[string]interface{} {
		err := api.ImportImagesFile(imgFiles, outFile, nil, nil)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})
}
