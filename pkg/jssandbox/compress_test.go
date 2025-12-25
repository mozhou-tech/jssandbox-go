package jssandbox

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dop251/goja"
)

func TestCompressZip(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建测试文件
	tempDir := t.TempDir()
	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")
	os.WriteFile(file1, []byte("Content 1"), 0644)
	os.WriteFile(file2, []byte("Content 2"), 0644)

	zipPath := filepath.Join(tempDir, "test.zip")

	code := `
		var result = compressZip(["` + file1 + `", "` + file2 + `"], "` + zipPath + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("compressZip() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("compressZip()应该返回success: true")
	}

	// 验证ZIP文件存在
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Error("ZIP文件应该被创建")
	}
}

func TestExtractZip(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 先创建ZIP文件
	tempDir := t.TempDir()
	file1 := filepath.Join(tempDir, "file1.txt")
	os.WriteFile(file1, []byte("Content 1"), 0644)

	zipPath := filepath.Join(tempDir, "test.zip")
	extractDir := filepath.Join(tempDir, "extracted")

	// 压缩
	compressCode := `
		var result = compressZip(["` + file1 + `"], "` + zipPath + `");
		result.success;
	`

	_, err := sb.Run(compressCode)
	if err != nil {
		t.Fatalf("压缩失败: %v", err)
	}

	// 解压
	extractCode := `
		var result = extractZip("` + zipPath + `", "` + extractDir + `");
		result;
	`

	result, err := sb.Run(extractCode)
	if err != nil {
		t.Fatalf("extractZip() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("extractZip()应该返回success: true")
	}

	files := resultObj.Get("files")
	if files == nil || goja.IsUndefined(files) {
		t.Error("extractZip()缺少files字段")
	}
}

func TestCompressGzip(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = compressGzip("Hello, World!");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("compressGzip() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("compressGzip()应该返回success: true")
	}

	data := resultObj.Get("data")
	if len(data.String()) == 0 {
		t.Error("compressGzip()返回的压缩数据为空")
	}
}

func TestDecompressGzip(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	original := "Hello, World!"

	// 先压缩
	compressCode := `
		var compressed = compressGzip("` + original + `");
		compressed.data;
	`

	compressedResult, err := sb.Run(compressCode)
	if err != nil {
		t.Fatalf("压缩失败: %v", err)
	}

	compressedData := compressedResult.String()

	// 再解压
	decompressCode := `
		var result = decompressGzip("` + compressedData + `");
		result;
	`

	result, err := sb.Run(decompressCode)
	if err != nil {
		t.Fatalf("decompressGzip() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("decompressGzip()应该返回success: true")
	}

	data := resultObj.Get("data")
	if data.String() != original {
		t.Errorf("decompressGzip()解压结果不正确, got %s, want %s", data.String(), original)
	}
}

func TestCompressDecompressGzip_Integration(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []string{
		"Hello, World!",
		"测试中文",
		"123456",
		"Special chars: !@#$%^&*()",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			code := `
				var compressed = compressGzip("` + tc + `");
				var decompressed = decompressGzip(compressed.data);
				decompressed.data;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("压缩解压失败: %v", err)
			}

			if result.String() != tc {
				t.Errorf("压缩解压结果不匹配, got %s, want %s", result.String(), tc)
			}
		})
	}
}
