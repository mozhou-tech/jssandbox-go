package jssandbox

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestExcel(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	tempDir, err := os.MkdirTemp("", "excel_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	outFile := filepath.Join(tempDir, "test.xlsx")

	script := `
		const f = excelNew();
		const sheet = "Sheet1";
		
		// 设置单元格值
		excelSetCellValue(f, sheet, "A1", "姓名");
		excelSetCellValue(f, sheet, "B1", "年龄");
		excelSetCellValue(f, sheet, "A2", "张三");
		excelSetCellValue(f, sheet, "B2", 25);
		
		// 设置整行
		excelSetSheetRow(f, sheet, "A3", ["李四", 30, "上海"]);
		
		// 合并单元格
		excelMergeCell(f, sheet, "A4", "B4");
		excelSetCellValue(f, sheet, "A4", "合并单元格");
		
		// 保存
		excelSave(f, "` + outFile + `");
		excelClose(f);
		
		// 重新打开并读取
		const f2 = excelOpen("` + outFile + `");
		const val1 = excelGetCellValue(f2, sheet, "A1");
		const val2 = excelGetCellValue(f2, sheet, "A3");
		const rows = excelGetRows(f2, sheet);
		excelClose(f2);
		
		// 测试 readExcel
		const excelRes = readExcel("` + outFile + `", { page: 1, pageSize: 2 });
		
		({ val1, val2, rows, excelRes });
	`

	result, err := sb.Run(script)
	if err != nil {
		t.Fatalf("执行脚本失败: %v", err)
	}

	res := result.Export().(map[string]interface{})
	if res["val1"] != "姓名" {
		t.Errorf("期望 A1 为 '姓名', 得到 '%v'", res["val1"])
	}
	if res["val2"] != "李四" {
		t.Errorf("期望 A3 为 '李四', 得到 '%v'", res["val2"])
	}

	rows := res["rows"].([][]string)
	if len(rows) < 3 {
		t.Errorf("行数不足: %d", len(rows))
	}

	excelRes := res["excelRes"].(map[string]interface{})
	if tr, ok := excelRes["totalRows"].(int); ok {
		if tr < 3 {
			t.Errorf("readExcel totalRows 错误: %v", tr)
		}
	} else if tr, ok := excelRes["totalRows"].(int64); ok {
		if tr < 3 {
			t.Errorf("readExcel totalRows 错误: %v", tr)
		}
	} else {
		t.Errorf("readExcel totalRows 类型错误: %T", excelRes["totalRows"])
	}

	if rows, ok := excelRes["rows"].([][]string); ok {
		if len(rows) != 2 {
			t.Errorf("readExcel 分页大小错误: %d", len(rows))
		}
	} else {
		t.Errorf("readExcel rows 类型错误: %T", excelRes["rows"])
	}
}
