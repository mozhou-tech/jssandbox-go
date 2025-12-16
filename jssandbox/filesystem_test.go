package jssandbox

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/dop251/goja"
)

func TestWriteFile(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.txt")
	content := "Hello, World!"

	code := `
		var result = writeFile("` + testFile + `", "` + content + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("writeFile() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("writeFile()应该返回success: true")
	}

	// 验证文件确实被写入
	readContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("读取文件失败: %v", err)
	}
	if string(readContent) != content {
		t.Errorf("文件内容不正确, got %s, want %s", string(readContent), content)
	}
}

func TestReadFile(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.txt")
	content := "Test content\nLine 2\nLine 3"
	os.WriteFile(testFile, []byte(content), 0644)

	code := `
		var result = readFile("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readFile() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	data := resultObj.Get("data")
	if data.String() != content {
		t.Errorf("readFile()内容不正确, got %s, want %s", data.String(), content)
	}

	length := resultObj.Get("length")
	if length.ToInteger() != int64(len(content)) {
		t.Errorf("readFile()长度不正确, got %d, want %d", length.ToInteger(), len(content))
	}
}

func TestReadFile_WithOptions(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.txt")
	content := "0123456789"
	os.WriteFile(testFile, []byte(content), 0644)

	code := `
		var result = readFile("` + testFile + `", {
			offset: 2,
			limit: 5
		});
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readFile() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	data := resultObj.Get("data")
	if data.String() != "23456" {
		t.Errorf("readFile()带选项的内容不正确, got %s, want 23456", data.String())
	}
}

func TestReadFileHead(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.txt")
	content := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"
	os.WriteFile(testFile, []byte(content), 0644)

	code := `
		var result = readFileHead("` + testFile + `", 3);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readFileHead() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	lines := resultObj.Get("lines")
	linesArray := lines.ToObject(sb.vm)

	count := resultObj.Get("count")
	if count.ToInteger() != 3 {
		t.Errorf("readFileHead()行数不正确, got %d, want 3", count.ToInteger())
	}

	// 验证第一行
	firstLine := linesArray.Get("0")
	if firstLine.String() != "Line 1" {
		t.Errorf("readFileHead()第一行不正确, got %s, want Line 1", firstLine.String())
	}
}

func TestReadFileTail(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.txt")
	content := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5"
	os.WriteFile(testFile, []byte(content), 0644)

	code := `
		var result = readFileTail("` + testFile + `", 2);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readFileTail() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	count := resultObj.Get("count")
	if count.ToInteger() != 2 {
		t.Errorf("readFileTail()行数不正确, got %d, want 2", count.ToInteger())
	}

	lines := resultObj.Get("lines")
	linesArray := lines.ToObject(sb.vm)
	lastLine := linesArray.Get("1")
	if lastLine.String() != "Line 5" {
		t.Errorf("readFileTail()最后一行不正确, got %s, want Line 5", lastLine.String())
	}
}

func TestGetFileInfo(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.txt")
	content := "Test content"
	os.WriteFile(testFile, []byte(content), 0644)

	code := `
		var result = getFileInfo("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("getFileInfo() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)

	// 验证必要字段
	name := resultObj.Get("name")
	if name.String() != "test.txt" {
		t.Errorf("getFileInfo()文件名不正确, got %s, want test.txt", name.String())
	}

	size := resultObj.Get("size")
	if size.ToInteger() != int64(len(content)) {
		t.Errorf("getFileInfo()文件大小不正确, got %d, want %d", size.ToInteger(), len(content))
	}

	isDir := resultObj.Get("isDir")
	if isDir.ToBoolean() {
		t.Error("getFileInfo()isDir应该为false")
	}

	extension := resultObj.Get("extension")
	if extension.String() != ".txt" {
		t.Errorf("getFileInfo()扩展名不正确, got %s, want .txt", extension.String())
	}
}

func TestGetFileHash(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.txt")
	content := "Test content for hash"
	os.WriteFile(testFile, []byte(content), 0644)

	// 计算预期的MD5
	hash := md5.Sum([]byte(content))
	expectedHash := hex.EncodeToString(hash[:])

	code := `
		var result = getFileHash("` + testFile + `", "md5");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("getFileHash() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	hashVal := resultObj.Get("hash")
	if hashVal.String() != expectedHash {
		t.Errorf("getFileHash() MD5不正确, got %s, want %s", hashVal.String(), expectedHash)
	}

	hashType := resultObj.Get("type")
	if hashType.String() != "md5" {
		t.Errorf("getFileHash()类型不正确, got %s, want md5", hashType.String())
	}
}

func TestGetFileHash_AllTypes(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.txt")
	content := "Test content"
	os.WriteFile(testFile, []byte(content), 0644)

	hashTypes := []string{"md5", "sha1", "sha256", "sha512"}

	for _, hashType := range hashTypes {
		code := `
			var result = getFileHash("` + testFile + `", "` + hashType + `");
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("getFileHash() %s error = %v", hashType, err)
		}

		resultObj := result.ToObject(sb.vm)
		hashVal := resultObj.Get("hash")
		if hashVal.String() == "" {
			t.Errorf("getFileHash() %s返回空哈希值", hashType)
		}

		typeVal := resultObj.Get("type")
		if typeVal.String() != hashType {
			t.Errorf("getFileHash() %s类型不正确, got %s", hashType, typeVal.String())
		}
	}
}

func TestGetFileHash_InvalidType(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	code := `
		var result = getFileHash("` + testFile + `", "invalid");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("getFileHash() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	errorVal := resultObj.Get("error")
	if errorVal == nil || goja.IsUndefined(errorVal) {
		t.Error("getFileHash()无效类型应该返回错误")
	}
}

func TestRenameFile(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "old.txt")
	newFile := filepath.Join(t.TempDir(), "new.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	code := `
		var result = renameFile("` + testFile + `", "` + newFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("renameFile() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("renameFile()应该返回success: true")
	}

	// 验证文件确实被重命名
	if _, err := os.Stat(newFile); os.IsNotExist(err) {
		t.Error("重命名后新文件不存在")
	}
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("重命名后旧文件仍然存在")
	}
}

func TestAppendFile(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.txt")
	os.WriteFile(testFile, []byte("Original"), 0644)

	code := `
		var result = appendFile("` + testFile + `", "Appended");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("appendFile() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("appendFile()应该返回success: true")
	}

	// 验证内容被追加
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("读取文件失败: %v", err)
	}
	if string(content) != "OriginalAppended" {
		t.Errorf("appendFile()内容不正确, got %s, want OriginalAppended", string(content))
	}
}

func TestReadImageBase64(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建一个简单的PNG文件（最小有效的1x1 PNG）
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}

	testFile := filepath.Join(t.TempDir(), "test.png")
	os.WriteFile(testFile, pngData, 0644)

	code := `
		var result = readImageBase64("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readImageBase64() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	base64Val := resultObj.Get("base64")
	if base64Val.String() == "" {
		t.Error("readImageBase64()返回空base64")
	}

	mimeType := resultObj.Get("mimeType")
	if mimeType.String() != "image/png" {
		t.Errorf("readImageBase64() MIME类型不正确, got %s, want image/png", mimeType.String())
	}

	dataUrl := resultObj.Get("dataUrl")
	if dataUrl.String() == "" {
		t.Error("readImageBase64()返回空dataUrl")
	}
}

func TestFileSystem_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	t.Run("读取不存在的文件", func(t *testing.T) {
		code := `
			var result = readFile("/nonexistent/file.txt");
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("readFile() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("读取不存在文件应该返回错误")
		}
	})

	t.Run("获取不存在文件的信息", func(t *testing.T) {
		code := `
			var result = getFileInfo("/nonexistent/file.txt");
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("getFileInfo() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("获取不存在文件信息应该返回错误")
		}
	})
}

