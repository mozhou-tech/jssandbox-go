package jssandbox

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dop251/goja"
)

// 注意：文档测试需要实际的文档文件
// 这些测试在缺少文档文件时会跳过

func TestReadWord(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建一个简单的测试文件路径（实际测试需要真实的Word文件）
	testFile := filepath.Join(t.TempDir(), "test.docx")

	code := `
		var result = readWord("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readWord() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("readWord()返回的对象为nil")
	}

	// 检查是否有错误（文件不存在时会返回错误）
	errorVal := resultObj.Get("error")
	if errorVal != nil && !goja.IsUndefined(errorVal) {
		t.Logf("readWord()文件不存在或无法读取（这是预期的）: %s", errorVal.String())
		return
	}

	// 如果文件存在，验证返回结构
	page := resultObj.Get("page")
	if page == nil || goja.IsUndefined(page) {
		t.Error("readWord()缺少page字段")
	}

	totalPages := resultObj.Get("totalPages")
	if totalPages == nil || goja.IsUndefined(totalPages) {
		t.Error("readWord()缺少totalPages字段")
	}
}

func TestReadWord_WithOptions(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.docx")

	code := `
		var result = readWord("` + testFile + `", {
			page: 2,
			pageSize: 500
		});
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readWord() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("readWord()返回的对象为nil")
	}

	// 验证选项是否正确应用
	page := resultObj.Get("page")
	if page != nil && !goja.IsUndefined(page) && page.ToInteger() != 2 {
		t.Errorf("readWord()页码不正确, got %d, want 2", page.ToInteger())
	}

	pageSize := resultObj.Get("pageSize")
	if pageSize != nil && !goja.IsUndefined(pageSize) && pageSize.ToInteger() != 500 {
		t.Errorf("readWord()页面大小不正确, got %d, want 500", pageSize.ToInteger())
	}
}

func TestReadExcel(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.xlsx")

	code := `
		var result = readExcel("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readExcel() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("readExcel()返回的对象为nil")
	}

	errorVal := resultObj.Get("error")
	if errorVal != nil && !goja.IsUndefined(errorVal) {
		t.Logf("readExcel()文件不存在或无法读取（这是预期的）: %s", errorVal.String())
		return
	}

	// 验证返回结构
	rows := resultObj.Get("rows")
	if rows == nil || goja.IsUndefined(rows) {
		t.Error("readExcel()缺少rows字段")
	}

	sheetName := resultObj.Get("sheetName")
	if sheetName == nil || goja.IsUndefined(sheetName) {
		t.Error("readExcel()缺少sheetName字段")
	}
}

func TestReadExcel_WithOptions(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.xlsx")

	code := `
		var result = readExcel("` + testFile + `", {
			sheet: "Sheet1",
			page: 1,
			pageSize: 50
		});
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readExcel() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("readExcel()返回的对象为nil")
	}

	errorVal := resultObj.Get("error")
	if errorVal != nil && !goja.IsUndefined(errorVal) {
		t.Logf("readExcel()文件不存在或无法读取: %s", errorVal.String())
		return
	}

	// 验证选项
	page := resultObj.Get("page")
	if page != nil && !goja.IsUndefined(page) && page.ToInteger() != 1 {
		t.Errorf("readExcel()页码不正确, got %d, want 1", page.ToInteger())
	}
}

func TestReadPPT(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.pptx")

	code := `
		var result = readPPT("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readPPT() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("readPPT()返回的对象为nil")
	}

	errorVal := resultObj.Get("error")
	if errorVal != nil && !goja.IsUndefined(errorVal) {
		t.Logf("readPPT()文件不存在或无法读取（这是预期的）: %s", errorVal.String())
		return
	}

	// 验证返回结构
	slides := resultObj.Get("slides")
	if slides == nil || goja.IsUndefined(slides) {
		t.Error("readPPT()缺少slides字段")
	}

	totalSlides := resultObj.Get("totalSlides")
	if totalSlides == nil || goja.IsUndefined(totalSlides) {
		t.Error("readPPT()缺少totalSlides字段")
	}
}

func TestReadPPT_WithOptions(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.pptx")

	code := `
		var result = readPPT("` + testFile + `", {
			page: 1,
			pageSize: 10
		});
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readPPT() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("readPPT()返回的对象为nil")
	}

	errorVal := resultObj.Get("error")
	if errorVal != nil && !goja.IsUndefined(errorVal) {
		t.Logf("readPPT()文件不存在或无法读取: %s", errorVal.String())
		return
	}

	pageSize := resultObj.Get("pageSize")
	if pageSize != nil && !goja.IsUndefined(pageSize) && pageSize.ToInteger() != 10 {
		t.Errorf("readPPT()页面大小不正确, got %d, want 10", pageSize.ToInteger())
	}
}

func TestReadPDF(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.pdf")

	code := `
		var result = readPDF("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readPDF() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("readPDF()返回的对象为nil")
	}

	errorVal := resultObj.Get("error")
	if errorVal != nil && !goja.IsUndefined(errorVal) {
		t.Logf("readPDF()文件不存在或无法读取（这是预期的）: %s", errorVal.String())
		return
	}

	// 验证返回结构
	pages := resultObj.Get("pages")
	if pages == nil || goja.IsUndefined(pages) {
		t.Error("readPDF()缺少pages字段")
	}

	totalPDFPages := resultObj.Get("totalPDFPages")
	if totalPDFPages == nil || goja.IsUndefined(totalPDFPages) {
		t.Error("readPDF()缺少totalPDFPages字段")
	}
}

func TestReadPDF_WithOptions(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testFile := filepath.Join(t.TempDir(), "test.pdf")

	code := `
		var result = readPDF("` + testFile + `", {
			page: 1,
			pageSize: 5
		});
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readPDF() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("readPDF()返回的对象为nil")
	}

	errorVal := resultObj.Get("error")
	if errorVal != nil && !goja.IsUndefined(errorVal) {
		t.Logf("readPDF()文件不存在或无法读取: %s", errorVal.String())
		return
	}

	pageSize := resultObj.Get("pageSize")
	if pageSize != nil && !goja.IsUndefined(pageSize) && pageSize.ToInteger() != 5 {
		t.Errorf("readPDF()页面大小不正确, got %d, want 5", pageSize.ToInteger())
	}
}

func TestDocuments_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	t.Run("readWord缺少参数", func(t *testing.T) {
		result, err := sb.Run("readWord()")
		if err != nil {
			t.Fatalf("readWord() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("readWord()缺少参数应该返回错误")
		}
	})

	t.Run("readExcel缺少参数", func(t *testing.T) {
		result, err := sb.Run("readExcel()")
		if err != nil {
			t.Fatalf("readExcel() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("readExcel()缺少参数应该返回错误")
		}
	})

	t.Run("readPPT缺少参数", func(t *testing.T) {
		result, err := sb.Run("readPPT()")
		if err != nil {
			t.Fatalf("readPPT() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("readPPT()缺少参数应该返回错误")
		}
	})

	t.Run("readPDF缺少参数", func(t *testing.T) {
		result, err := sb.Run("readPDF()")
		if err != nil {
			t.Fatalf("readPDF() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("readPDF()缺少参数应该返回错误")
		}
	})
}

func TestDocuments_InvalidFile(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建一个无效的文件
	invalidFile := filepath.Join(t.TempDir(), "invalid.txt")
	os.WriteFile(invalidFile, []byte("not a document"), 0644)

	t.Run("用文本文件测试readWord", func(t *testing.T) {
		code := `
			var result = readWord("` + invalidFile + `");
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("readWord() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("readWord()处理无效文件应该返回错误")
		}
	})

	t.Run("用文本文件测试readExcel", func(t *testing.T) {
		code := `
			var result = readExcel("` + invalidFile + `");
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("readExcel() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("readExcel()处理无效文件应该返回错误")
		}
	})
}

