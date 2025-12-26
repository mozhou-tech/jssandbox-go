package jssandbox

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDocx(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	tempDir, err := os.MkdirTemp("", "docx_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	outFile := filepath.Join(tempDir, "test.docx")

	script := `
		const doc = docxNew();
		docxAddHeading(doc, "测试标题", 1);
		docxAddParagraph(doc, "这是一个普通段落。");
		
		const format = { bold: true, fontSize: 14, fontColor: "FF0000" };
		docxAddFormattedParagraph(doc, "这是一个加粗红色段落。", format);
		
		const tableConfig = {
			rows: 2,
			cols: 2,
			width: 5000,
			data: [
				["单元格1-1", "单元格1-2"],
				["单元格2-1", "单元格2-2"]
			]
		};
		const table = docxAddTable(doc, tableConfig);
		docxSetCellText(table, 0, 0, "修改后的内容");
		
		docxSave(doc, "` + outFile + `");
		
		const text = docxReadText("` + outFile + `");
		text;
	`

	result, err := sb.Run(script)
	if err != nil {
		t.Fatalf("执行脚本失败: %v", err)
	}

	text := result.String()
	t.Logf("读取到的文本:\n%s", text)

	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Errorf("文件未成功生成: %s", outFile)
	}
}
