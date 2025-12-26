package jssandbox

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPDFOps(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 准备测试数据
	tempDir, err := os.MkdirTemp("", "pdf-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	imgFile := "example/data/bot-detection-screenshot.png"
	pdfFile := filepath.Join(tempDir, "test.pdf")

	// 1. 测试图片导入为PDF
	t.Run("ImportImages", func(t *testing.T) {
		code := `
			var result = pdfImportImages(["` + imgFile + `"], "` + pdfFile + `");
			result;
		`
		val, err := sb.Run(code)
		assert.NoError(t, err)
		result := val.Export().(map[string]interface{})
		assert.True(t, result["success"].(bool))
		assert.FileExists(t, pdfFile)
	})

	// 2. 测试获取页数
	t.Run("GetPageCount", func(t *testing.T) {
		code := `
			var result = pdfGetPageCount("` + pdfFile + `");
			result;
		`
		val, err := sb.Run(code)
		assert.NoError(t, err)
		result := val.Export().(map[string]interface{})
		assert.True(t, result["success"].(bool))
		// goja 导出的数字可能是 int 或 float64
		pages := result["pages"]
		if p, ok := pages.(int64); ok {
			assert.Equal(t, int64(1), p)
		} else if p, ok := pages.(int); ok {
			assert.Equal(t, 1, p)
		} else {
			assert.Fail(t, "页数类型错误")
		}
	})

	// 3. 测试添加水印
	t.Run("AddTextWatermark", func(t *testing.T) {
		watermarkedFile := filepath.Join(tempDir, "watermarked.pdf")
		code := `
			var result = pdfAddTextWatermark("` + pdfFile + `", "` + watermarkedFile + `", "CONFIDENTIAL", {
				opacity: 0.5,
				scale: 0.5,
				rotation: 45
			});
			result;
		`
		val, err := sb.Run(code)
		assert.NoError(t, err)
		result := val.Export().(map[string]interface{})
		assert.True(t, result["success"].(bool))
		assert.FileExists(t, watermarkedFile)
	})

	// 4. 测试拆分PDF
	t.Run("Split", func(t *testing.T) {
		splitDir := filepath.Join(tempDir, "split")
		code := `
			var result = pdfSplit("` + pdfFile + `", "` + splitDir + `");
			result;
		`
		val, err := sb.Run(code)
		assert.NoError(t, err)
		result := val.Export().(map[string]interface{})
		assert.True(t, result["success"].(bool))

		files, _ := os.ReadDir(splitDir)
		assert.Greater(t, len(files), 0)
	})

	// 5. 测试合并PDF
	t.Run("Merge", func(t *testing.T) {
		mergedFile := filepath.Join(tempDir, "merged.pdf")
		code := `
			var result = pdfMerge(["` + pdfFile + `", "` + pdfFile + `"], "` + mergedFile + `");
			result;
		`
		val, err := sb.Run(code)
		assert.NoError(t, err)
		result := val.Export().(map[string]interface{})
		assert.True(t, result["success"].(bool))
		assert.FileExists(t, mergedFile)

		// 检查页数是否为2
		code = `pdfGetPageCount("` + mergedFile + `").pages`
		val, err = sb.Run(code)
		assert.NoError(t, err)
		assert.Equal(t, 2, int(val.ToInteger()))
	})
}
