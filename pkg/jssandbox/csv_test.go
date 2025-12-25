package jssandbox

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dop251/goja"
)

func TestReadCSV(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建测试CSV文件
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "test.csv")
	csvContent := "name,age,city\nJohn,25,New York\nJane,30,London"
	os.WriteFile(csvFile, []byte(csvContent), 0644)

	code := `
		var result = readCSV("` + csvFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readCSV() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("readCSV()应该返回success: true")
	}

	count := resultObj.Get("count")
	if count.ToInteger() != 2 {
		t.Errorf("readCSV()行数不正确, got %d, want 2", count.ToInteger())
	}

	rows := resultObj.Get("rows")
	if rows == nil || goja.IsUndefined(rows) {
		t.Error("readCSV()缺少rows字段")
	}
}

func TestReadCSV_WithOptions(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建使用分号分隔的CSV文件
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "test.csv")
	csvContent := "name;age;city\nJohn;25;New York"
	os.WriteFile(csvFile, []byte(csvContent), 0644)

	code := `
		var result = readCSV("` + csvFile + `", {delimiter: ";"});
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readCSV() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("readCSV()应该返回success: true")
	}

	count := resultObj.Get("count")
	if count.ToInteger() != 1 {
		t.Errorf("readCSV()行数不正确, got %d, want 1", count.ToInteger())
	}
}

func TestWriteCSV(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "output.csv")

	code := `
		var data = [
			["name", "age", "city"],
			["John", "25", "New York"],
			["Jane", "30", "London"]
		];
		var result = writeCSV("` + csvFile + `", data);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("writeCSV() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("writeCSV()应该返回success: true")
	}

	// 验证文件被创建
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		t.Error("CSV文件应该被创建")
	}

	// 验证文件内容
	content, err := os.ReadFile(csvFile)
	if err != nil {
		t.Fatalf("读取CSV文件失败: %v", err)
	}

	if len(content) == 0 {
		t.Error("CSV文件内容不应该为空")
	}
}

func TestParseCSV(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	csvString := "name,age,city\nJohn,25,New York\nJane,30,London"

	code := `
		var result = parseCSV("` + csvString + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("parseCSV() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("parseCSV()应该返回success: true")
	}

	count := resultObj.Get("count")
	if count.ToInteger() != 2 {
		t.Errorf("parseCSV()行数不正确, got %d, want 2", count.ToInteger())
	}
}

func TestCSV_ReadWriteIntegration(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "test.csv")

	// 写入CSV
	writeCode := `
		var data = [
			["name", "age"],
			["John", "25"]
		];
		writeCSV("` + csvFile + `", data);
	`

	_, err := sb.Run(writeCode)
	if err != nil {
		t.Fatalf("写入CSV失败: %v", err)
	}

	// 读取CSV
	readCode := `
		var result = readCSV("` + csvFile + `");
		result.count;
	`

	result, err := sb.Run(readCode)
	if err != nil {
		t.Fatalf("读取CSV失败: %v", err)
	}

	if result.ToInteger() != 1 {
		t.Errorf("读取CSV行数不正确, got %d, want 1", result.ToInteger())
	}
}

