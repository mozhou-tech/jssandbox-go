package jssandbox

import (
	"fmt"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
	"github.com/dop251/goja"
)

// registerDocx 注册Word (.docx) 处理功能到JavaScript运行时
func (sb *Sandbox) registerDocx() {
	// 创建新文档
	sb.vm.Set("docxNew", func() *document.Document {
		return document.New()
	})

	// 打开现有文档
	sb.vm.Set("docxOpen", func(filePath string) (goja.Value, error) {
		doc, err := document.Open(filePath)
		if err != nil {
			return nil, err
		}
		return sb.vm.ToValue(doc), nil
	})

	// 保存文档
	sb.vm.Set("docxSave", func(doc *document.Document, filePath string) error {
		if doc == nil {
			return fmt.Errorf("文档对象不能为空")
		}
		return doc.Save(filePath)
	})

	// 添加普通段落
	sb.vm.Set("docxAddParagraph", func(doc *document.Document, text string) *document.Paragraph {
		if doc == nil {
			return nil
		}
		return doc.AddParagraph(text)
	})

	// 添加标题段落
	sb.vm.Set("docxAddHeading", func(doc *document.Document, text string, level int) *document.Paragraph {
		if doc == nil {
			return nil
		}
		return doc.AddHeadingParagraph(text, level)
	})

	// 添加格式化段落
	// format: { bold: bool, italic: bool, fontSize: int, fontColor: string, fontFamily: string }
	sb.vm.Set("docxAddFormattedParagraph", func(doc *document.Document, text string, format map[string]interface{}) *document.Paragraph {
		if doc == nil {
			return nil
		}
		tf := &document.TextFormat{}
		if v, ok := format["bold"].(bool); ok {
			tf.Bold = v
		}
		if v, ok := format["italic"].(bool); ok {
			tf.Italic = v
		}
		if v, ok := format["fontSize"].(int64); ok {
			tf.FontSize = int(v)
		}
		if v, ok := format["fontColor"].(string); ok {
			tf.FontColor = v
		}
		if v, ok := format["fontFamily"].(string); ok {
			tf.FontFamily = v
		}
		return doc.AddFormattedParagraph(text, tf)
	})

	// 向段落添加格式化文本
	sb.vm.Set("docxAddFormattedText", func(para *document.Paragraph, text string, format map[string]interface{}) {
		if para == nil {
			return
		}
		tf := &document.TextFormat{}
		if v, ok := format["bold"].(bool); ok {
			tf.Bold = v
		}
		if v, ok := format["italic"].(bool); ok {
			tf.Italic = v
		}
		if v, ok := format["fontSize"].(int64); ok {
			tf.FontSize = int(v)
		}
		if v, ok := format["fontColor"].(string); ok {
			tf.FontColor = v
		}
		if v, ok := format["fontFamily"].(string); ok {
			tf.FontFamily = v
		}
		para.AddFormattedText(text, tf)
	})

	// 添加分页符
	sb.vm.Set("docxAddPageBreak", func(doc *document.Document) {
		if doc == nil {
			return
		}
		doc.AddPageBreak()
	})

	// 添加表格
	// config: { rows: int, cols: int, width: int, data: [][]string }
	sb.vm.Set("docxAddTable", func(doc *document.Document, config map[string]interface{}) (*document.Table, error) {
		if doc == nil {
			return nil, fmt.Errorf("文档对象不能为空")
		}
		tc := &document.TableConfig{}
		if v, ok := config["rows"].(int64); ok {
			tc.Rows = int(v)
		}
		if v, ok := config["cols"].(int64); ok {
			tc.Cols = int(v)
		}
		if v, ok := config["width"].(int64); ok {
			tc.Width = int(v)
		}
		if v, ok := config["data"].([]interface{}); ok {
			var data [][]string
			for _, row := range v {
				if r, ok := row.([]interface{}); ok {
					var strRow []string
					for _, cell := range r {
						strRow = append(strRow, fmt.Sprintf("%v", cell))
					}
					data = append(data, strRow)
				}
			}
			tc.Data = data
		}
		return doc.AddTable(tc)
	})

	// 设置单元格文本
	sb.vm.Set("docxSetCellText", func(table *document.Table, row, col int, text string) error {
		if table == nil {
			return fmt.Errorf("表格对象不能为空")
		}
		return table.SetCellText(row, col, text)
	})

	// 获取单元格文本
	sb.vm.Set("docxGetCellText", func(table *document.Table, row, col int) (string, error) {
		if table == nil {
			return "", fmt.Errorf("表格对象不能为空")
		}
		return table.GetCellText(row, col)
	})

	// 获取表格行数和列数
	sb.vm.Set("docxGetTableSize", func(table *document.Table) map[string]int {
		if table == nil {
			return nil
		}
		return map[string]int{
			"rows": table.GetRowCount(),
			"cols": table.GetColumnCount(),
		}
	})

	// 获取文档中的所有文本（包含表格内容）
	sb.vm.Set("docxReadText", func(filePath string) (string, error) {
		doc, err := document.Open(filePath)
		if err != nil {
			return "", err
		}
		var result string
		for _, element := range doc.Body.Elements {
			switch e := element.(type) {
			case *document.Paragraph:
				for _, run := range e.Runs {
					result += run.Text.Content
				}
				result += "\n"
			case *document.Table:
				for _, row := range e.Rows {
					var rowText []string
					for _, cell := range row.Cells {
						var cellText string
						for _, p := range cell.Paragraphs {
							for _, r := range p.Runs {
								cellText += r.Text.Content
							}
						}
						rowText = append(rowText, cellText)
					}
					result += "| " + fmt.Sprintf("%v", rowText) + " |\n"
				}
			}
		}
		return result, nil
	})
}
