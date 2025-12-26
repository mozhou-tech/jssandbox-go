package jssandbox

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/xuri/excelize/v2"
)

// registerExcel 注册 Excel (.xlsx) 处理功能到 JavaScript 运行时
func (sb *Sandbox) registerExcel() {
	// 创建新 Excel 文件
	sb.vm.Set("excelNew", func() *excelize.File {
		return excelize.NewFile()
	})

	// 打开现有 Excel 文件
	sb.vm.Set("excelOpen", func(filePath string) (goja.Value, error) {
		f, err := excelize.OpenFile(filePath)
		if err != nil {
			return nil, err
		}
		return sb.vm.ToValue(f), nil
	})

	// 保存 Excel 文件
	sb.vm.Set("excelSave", func(f *excelize.File, filePath string) error {
		if f == nil {
			return fmt.Errorf("Excel 对象不能为空")
		}
		return f.SaveAs(filePath)
	})

	// 关闭 Excel 文件
	sb.vm.Set("excelClose", func(f *excelize.File) error {
		if f == nil {
			return nil
		}
		return f.Close()
	})

	// 设置单元格值
	sb.vm.Set("excelSetCellValue", func(f *excelize.File, sheet, axis string, value interface{}) error {
		if f == nil {
			return fmt.Errorf("Excel 对象不能为空")
		}
		return f.SetCellValue(sheet, axis, value)
	})

	// 获取单元格值
	sb.vm.Set("excelGetCellValue", func(f *excelize.File, sheet, axis string) (string, error) {
		if f == nil {
			return "", fmt.Errorf("Excel 对象不能为空")
		}
		return f.GetCellValue(sheet, axis)
	})

	// 新建工作表
	sb.vm.Set("excelNewSheet", func(f *excelize.File, sheet string) (int, error) {
		if f == nil {
			return 0, fmt.Errorf("Excel 对象不能为空")
		}
		return f.NewSheet(sheet)
	})

	// 获取所有行
	sb.vm.Set("excelGetRows", func(f *excelize.File, sheet string) ([][]string, error) {
		if f == nil {
			return nil, fmt.Errorf("Excel 对象不能为空")
		}
		return f.GetRows(sheet)
	})

	// 获取所有列
	sb.vm.Set("excelGetCols", func(f *excelize.File, sheet string) ([][]string, error) {
		if f == nil {
			return nil, fmt.Errorf("Excel 对象不能为空")
		}
		return f.GetCols(sheet)
	})

	// 设置整行数据
	// slice: []interface{}
	sb.vm.Set("excelSetSheetRow", func(f *excelize.File, sheet, axis string, slice []interface{}) error {
		if f == nil {
			return fmt.Errorf("Excel 对象不能为空")
		}
		return f.SetSheetRow(sheet, axis, &slice)
	})

	// 设置活动工作表
	sb.vm.Set("excelSetActiveSheet", func(f *excelize.File, index int) {
		if f == nil {
			return
		}
		f.SetActiveSheet(index)
	})

	// 删除工作表
	sb.vm.Set("excelDeleteSheet", func(f *excelize.File, sheet string) error {
		if f == nil {
			return fmt.Errorf("Excel 对象不能为空")
		}
		return f.DeleteSheet(sheet)
	})

	// 复制工作表
	sb.vm.Set("excelCopySheet", func(f *excelize.File, from, to int) error {
		if f == nil {
			return fmt.Errorf("Excel 对象不能为空")
		}
		return f.CopySheet(from, to)
	})

	// 设置列宽
	sb.vm.Set("excelSetColWidth", func(f *excelize.File, sheet, startCol, endCol string, width float64) error {
		if f == nil {
			return fmt.Errorf("Excel 对象不能为空")
		}
		return f.SetColWidth(sheet, startCol, endCol, width)
	})

	// 设置行高
	sb.vm.Set("excelSetRowHeight", func(f *excelize.File, sheet string, row int, height float64) error {
		if f == nil {
			return fmt.Errorf("Excel 对象不能为空")
		}
		return f.SetRowHeight(sheet, row, height)
	})

	// 合并单元格
	sb.vm.Set("excelMergeCell", func(f *excelize.File, sheet, hcell, vcell string) error {
		if f == nil {
			return fmt.Errorf("Excel 对象不能为空")
		}
		return f.MergeCell(sheet, hcell, vcell)
	})

	// 取消合并单元格
	sb.vm.Set("excelUnmergeCell", func(f *excelize.File, sheet, hcell, vcell string) error {
		if f == nil {
			return fmt.Errorf("Excel 对象不能为空")
		}
		return f.UnmergeCell(sheet, hcell, vcell)
	})

	// readExcel 读取 Excel 文件（支持分页和指定工作表）
	// options: { sheet: string, page: int, pageSize: int }
	sb.vm.Set("readExcel", func(filePath string, options map[string]interface{}) (map[string]interface{}, error) {
		f, err := excelize.OpenFile(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		sheetName := "Sheet1"
		if s, ok := options["sheet"].(string); ok {
			sheetName = s
		} else {
			// 如果没指定，尝试获取第一个工作表
			sheets := f.GetSheetList()
			if len(sheets) > 0 {
				sheetName = sheets[0]
			}
		}

		rows, err := f.GetRows(sheetName)
		if err != nil {
			return nil, err
		}

		totalRows := len(rows)
		page := 1
		pageSize := totalRows

		if p, ok := options["page"].(int64); ok {
			page = int(p)
		}
		if ps, ok := options["pageSize"].(int64); ok {
			pageSize = int(ps)
		}

		start := (page - 1) * pageSize
		if start < 0 {
			start = 0
		}
		end := start + pageSize
		if end > totalRows {
			end = totalRows
		}

		var resultRows [][]string
		if start < totalRows {
			resultRows = rows[start:end]
		}

		totalPages := 0
		if pageSize > 0 {
			totalPages = (totalRows + pageSize - 1) / pageSize
		}

		return map[string]interface{}{
			"rows":       resultRows,
			"totalRows":  totalRows,
			"totalPages": totalPages,
			"page":       page,
			"sheet":      sheetName,
		}, nil
	})
}
