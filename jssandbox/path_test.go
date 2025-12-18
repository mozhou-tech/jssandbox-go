package jssandbox

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dop251/goja"
)

func TestPathJoin(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = pathJoin("/usr", "local", "bin");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("pathJoin() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("pathJoin()应该返回success: true")
	}

	path := resultObj.Get("path")
	expected := filepath.Join("/usr", "local", "bin")
	if path.String() != expected {
		t.Errorf("pathJoin()结果不正确, got %s, want %s", path.String(), expected)
	}
}

func TestPathDir(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []struct {
		input    string
		expected string
	}{
		{"/usr/local/bin/app", "/usr/local/bin"},
		{"file.txt", "."},
		{"./path/to/file", "./path/to"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			code := `
				var result = pathDir("` + tc.input + `");
				result;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("pathDir() error = %v", err)
			}

			resultObj := result.ToObject(sb.vm)
			dir := resultObj.Get("dir")
			expected := filepath.Dir(tc.input)
			if dir.String() != expected {
				t.Errorf("pathDir(%s) = %s, want %s", tc.input, dir.String(), expected)
			}
		})
	}
}

func TestPathBase(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []struct {
		input    string
		expected string
	}{
		{"/usr/local/bin/app", "app"},
		{"file.txt", "file.txt"},
		{"./path/to/file", "file"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			code := `
				var result = pathBase("` + tc.input + `");
				result;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("pathBase() error = %v", err)
			}

			resultObj := result.ToObject(sb.vm)
			base := resultObj.Get("base")
			expected := filepath.Base(tc.input)
			if base.String() != expected {
				t.Errorf("pathBase(%s) = %s, want %s", tc.input, base.String(), expected)
			}
		})
	}
}

func TestPathExt(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []struct {
		input    string
		expected string
	}{
		{"file.txt", ".txt"},
		{"file.tar.gz", ".gz"},
		{"file", ""},
		{".hidden", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			code := `
				var result = pathExt("` + tc.input + `");
				result;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("pathExt() error = %v", err)
			}

			resultObj := result.ToObject(sb.vm)
			ext := resultObj.Get("ext")
			expected := filepath.Ext(tc.input)
			if ext.String() != expected {
				t.Errorf("pathExt(%s) = %s, want %s", tc.input, ext.String(), expected)
			}
		})
	}
}

func TestPathAbs(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建临时文件用于测试
	tempFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(tempFile, []byte("test"), 0644)

	code := `
		var result = pathAbs("` + tempFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("pathAbs() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("pathAbs()应该返回success: true")
	}

	absPath := resultObj.Get("path")
	if absPath.String() == "" {
		t.Error("pathAbs()返回的绝对路径不应该为空")
	}

	// 验证返回的是绝对路径
	expected, _ := filepath.Abs(tempFile)
	if absPath.String() != expected {
		t.Errorf("pathAbs()结果不正确, got %s, want %s", absPath.String(), expected)
	}
}

func TestPath_Integration(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var file = "/usr/local/bin/app.exe";
		var dir = pathDir(file);
		var base = pathBase(file);
		var ext = pathExt(file);
		var joined = pathJoin(dir.dir, base.base);
		{
			original: file,
			dir: dir.dir,
			base: base.base,
			ext: ext.ext,
			joined: joined.path
		};
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("路径处理集成测试失败: %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	ext := resultObj.Get("ext")
	if ext.String() != ".exe" {
		t.Errorf("集成测试扩展名不正确, got %s, want .exe", ext.String())
	}
}

