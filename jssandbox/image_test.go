package jssandbox

import (
	"context"
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/dop251/goja"
)

func TestImageResize(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建一个测试图片
	testDir := t.TempDir()
	inputPath := filepath.Join(testDir, "input.png")
	outputPath := filepath.Join(testDir, "output.png")

	// 创建一个简单的测试图片
	img := imaging.New(100, 100, imaging.White)
	err := imaging.Save(img, inputPath)
	if err != nil {
		t.Fatalf("创建测试图片失败: %v", err)
	}

	code := `
		var result = imageResize("` + inputPath + `", "` + outputPath + `", 50, 50);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("imageResize() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		errorVal := resultObj.Get("error")
		t.Errorf("imageResize()失败: %s", errorVal.String())
		return
	}

	// 验证输出文件存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("imageResize()输出文件不存在")
	}

	// 验证图片尺寸
	outputImg, err := imaging.Open(outputPath)
	if err != nil {
		t.Fatalf("打开输出图片失败: %v", err)
	}
	bounds := outputImg.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("imageResize()尺寸不正确, got %dx%d, want 50x50", bounds.Dx(), bounds.Dy())
	}
}

func TestImageInfo(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.png")

	// 创建一个测试图片
	img := imaging.New(200, 150, imaging.White)
	err := imaging.Save(img, testFile)
	if err != nil {
		t.Fatalf("创建测试图片失败: %v", err)
	}

	code := `
		var result = imageInfo("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("imageInfo() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	width := resultObj.Get("width")
	height := resultObj.Get("height")

	if width.ToInteger() != 200 {
		t.Errorf("imageInfo()宽度不正确, got %d, want 200", width.ToInteger())
	}
	if height.ToInteger() != 150 {
		t.Errorf("imageInfo()高度不正确, got %d, want 150", height.ToInteger())
	}
}

func TestImageCrop(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	inputPath := filepath.Join(testDir, "input.png")
	outputPath := filepath.Join(testDir, "output.png")

	img := imaging.New(100, 100, imaging.White)
	err := imaging.Save(img, inputPath)
	if err != nil {
		t.Fatalf("创建测试图片失败: %v", err)
	}

	code := `
		var result = imageCrop("` + inputPath + `", "` + outputPath + `", 10, 10, 50, 50);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("imageCrop() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("imageCrop()应该成功")
	}
}

func TestImageRotate(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	inputPath := filepath.Join(testDir, "input.png")
	outputPath := filepath.Join(testDir, "output.png")

	img := imaging.New(100, 100, imaging.White)
	err := imaging.Save(img, inputPath)
	if err != nil {
		t.Fatalf("创建测试图片失败: %v", err)
	}

	code := `
		var result = imageRotate("` + inputPath + `", "` + outputPath + `", 90);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("imageRotate() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("imageRotate()应该成功")
	}
}

func TestImageFlip(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	inputPath := filepath.Join(testDir, "input.png")
	outputPath := filepath.Join(testDir, "output.png")

	img := imaging.New(100, 100, imaging.White)
	err := imaging.Save(img, inputPath)
	if err != nil {
		t.Fatalf("创建测试图片失败: %v", err)
	}

	code := `
		var result = imageFlip("` + inputPath + `", "` + outputPath + `", "horizontal");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("imageFlip() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("imageFlip()应该成功")
	}
}

func TestImage_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	t.Run("imageResize缺少参数", func(t *testing.T) {
		result, err := sb.Run("imageResize()")
		if err != nil {
			t.Fatalf("imageResize() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("imageResize()缺少参数应该返回错误")
		}
	})

	t.Run("打开不存在的图片", func(t *testing.T) {
		testDir := t.TempDir()
		outputPath := filepath.Join(testDir, "output.png")

		code := `
			var result = imageResize("/nonexistent/image.png", "` + outputPath + `", 50, 50);
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("imageResize() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		success := resultObj.Get("success")
		if success.ToBoolean() {
			t.Error("打开不存在的图片应该失败")
		}
	})
}

