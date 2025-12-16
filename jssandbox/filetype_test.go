package jssandbox

import (
	"context"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/dop251/goja"
)

func TestDetectFileType(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	t.Run("检测PNG图片", func(t *testing.T) {
		testDir := t.TempDir()
		testFile := filepath.Join(testDir, "test.png")

		// 创建一个PNG图片
		img := imaging.New(10, 10, color.White)
		err := imaging.Save(img, testFile)
		if err != nil {
			t.Fatalf("创建测试图片失败: %v", err)
		}

		code := `
			var result = detectFileType("` + testFile + `");
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("detectFileType() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		if resultObj == nil {
			t.Fatal("detectFileType()返回的对象为nil")
		}

		unknown := resultObj.Get("unknown")
		if unknown != nil && !goja.IsUndefined(unknown) && unknown.ToBoolean() {
			t.Error("detectFileType()应该能识别PNG图片")
		}

		mime := resultObj.Get("mime")
		if mime == nil || goja.IsUndefined(mime) {
			t.Error("detectFileType()缺少mime字段")
		} else if mime.String() != "image/png" {
			t.Errorf("detectFileType() MIME类型不正确, got %s, want image/png", mime.String())
		}
	})

	t.Run("检测文本文件", func(t *testing.T) {
		testDir := t.TempDir()
		testFile := filepath.Join(testDir, "test.txt")
		os.WriteFile(testFile, []byte("test content"), 0644)

		code := `
			var result = detectFileType("` + testFile + `");
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("detectFileType() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		if resultObj == nil {
			t.Fatal("detectFileType()返回的对象为nil")
		}

		// 文本文件可能无法识别，这是正常的
		unknown := resultObj.Get("unknown")
		if unknown != nil && !goja.IsUndefined(unknown) {
			t.Logf("文本文件类型检测结果: unknown=%v", unknown.ToBoolean())
		}
	})
}

func TestIsImage(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.png")

	// 创建一个PNG图片
	img := imaging.New(10, 10, color.White)
	err := imaging.Save(img, testFile)
	if err != nil {
		t.Fatalf("创建测试图片失败: %v", err)
	}

	code := `
		var result = isImage("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("isImage() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("isImage()返回的对象为nil")
	}

	isImage := resultObj.Get("isImage")
	if isImage == nil || goja.IsUndefined(isImage) {
		t.Error("isImage()缺少isImage字段")
	} else if !isImage.ToBoolean() {
		t.Error("isImage()应该返回true")
	}
}

func TestIsVideo(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建一个非视频文件
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	code := `
		var result = isVideo("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("isVideo() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("isVideo()返回的对象为nil")
	}

	isVideo := resultObj.Get("isVideo")
	if isVideo == nil || goja.IsUndefined(isVideo) {
		t.Error("isVideo()缺少isVideo字段")
	} else if isVideo.ToBoolean() {
		t.Error("isVideo()对文本文件应该返回false")
	}
}

func TestIsAudio(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建一个非音频文件
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	code := `
		var result = isAudio("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("isAudio() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("isAudio()返回的对象为nil")
	}

	isAudio := resultObj.Get("isAudio")
	if isAudio == nil || goja.IsUndefined(isAudio) {
		t.Error("isAudio()缺少isAudio字段")
	} else if isAudio.ToBoolean() {
		t.Error("isAudio()对文本文件应该返回false")
	}
}

func TestFileType_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	t.Run("检测不存在的文件", func(t *testing.T) {
		code := `
			var result = detectFileType("/nonexistent/file.txt");
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("detectFileType() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		if resultObj == nil {
			t.Fatal("detectFileType()返回的对象为nil")
		}

		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("detectFileType()对不存在的文件应该返回错误")
		}
	})
}
