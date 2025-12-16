package jssandbox

import (
	"fmt"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/dop251/goja"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/presentation"
	"github.com/unidoc/unipdf/v3/model"
	"go.uber.org/zap"
)

// registerDocuments 注册文档读取功能到JavaScript运行时
func (sb *Sandbox) registerDocuments() {
	// 读取Word文件内容（支持分页）
	sb.vm.Set("readWord", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文件路径",
			})
		}

		filePath := call.Arguments[0].String()
		page := 1
		pageSize := 1000 // 每页字符数

		if len(call.Arguments) > 1 {
			options := call.Arguments[1].ToObject(sb.vm)
			if pageVal := options.Get("page"); pageVal != nil && !goja.IsUndefined(pageVal) {
				page = int(pageVal.ToInteger())
			}
			if pageSizeVal := options.Get("pageSize"); pageSizeVal != nil && !goja.IsUndefined(pageSizeVal) {
				pageSize = int(pageSizeVal.ToInteger())
			}
		}

		doc, err := document.Open(filePath)
		if err != nil {
			sb.logger.Error("打开Word文件失败", zap.String("path", filePath), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer doc.Close()

		var allText strings.Builder
		for _, para := range doc.Paragraphs() {
			allText.WriteString(para.Text())
			allText.WriteString("\n")
		}

		fullText := allText.String()
		totalChars := len(fullText)
		totalPages := (totalChars + pageSize - 1) / pageSize

		start := (page - 1) * pageSize
		end := start + pageSize
		if end > totalChars {
			end = totalChars
		}

		var pageText string
		if start < totalChars {
			pageText = fullText[start:end]
		}

		return sb.vm.ToValue(map[string]interface{}{
			"text":      pageText,
			"page":      page,
			"pageSize":  pageSize,
			"totalPages": totalPages,
			"totalChars": totalChars,
		})
	})

	// 读取Excel文件内容（支持分页）
	sb.vm.Set("readExcel", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文件路径",
			})
		}

		filePath := call.Arguments[0].String()
		sheetName := ""
		page := 1
		pageSize := 100 // 每页行数

		if len(call.Arguments) > 1 {
			options := call.Arguments[1].ToObject(sb.vm)
			if sheetVal := options.Get("sheet"); sheetVal != nil && !goja.IsUndefined(sheetVal) {
				sheetName = sheetVal.String()
			}
			if pageVal := options.Get("page"); pageVal != nil && !goja.IsUndefined(pageVal) {
				page = int(pageVal.ToInteger())
			}
			if pageSizeVal := options.Get("pageSize"); pageSizeVal != nil && !goja.IsUndefined(pageSizeVal) {
				pageSize = int(pageSizeVal.ToInteger())
			}
		}

		f, err := excelize.OpenFile(filePath)
		if err != nil {
			sb.logger.Error("打开Excel文件失败", zap.String("path", filePath), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer f.Close()

		if sheetName == "" {
			sheetName = f.GetSheetName(0)
		}

		rows, err := f.GetRows(sheetName)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		totalRows := len(rows)
		totalPages := (totalRows + pageSize - 1) / pageSize

		start := (page - 1) * pageSize
		end := start + pageSize
		if end > totalRows {
			end = totalRows
		}

		var pageRows [][]string
		if start < totalRows {
			pageRows = rows[start:end]
		}

		// 获取所有sheet名称
		sheetList := f.GetSheetList()

		return sb.vm.ToValue(map[string]interface{}{
			"rows":       pageRows,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
			"totalRows":  totalRows,
			"sheetName":  sheetName,
			"sheets":     sheetList,
		})
	})

	// 读取PPT文件内容（支持分页）
	sb.vm.Set("readPPT", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文件路径",
			})
		}

		filePath := call.Arguments[0].String()
		page := 1
		pageSize := 5 // 每页幻灯片数

		if len(call.Arguments) > 1 {
			options := call.Arguments[1].ToObject(sb.vm)
			if pageVal := options.Get("page"); pageVal != nil && !goja.IsUndefined(pageVal) {
				page = int(pageVal.ToInteger())
			}
			if pageSizeVal := options.Get("pageSize"); pageSizeVal != nil && !goja.IsUndefined(pageSizeVal) {
				pageSize = int(pageSizeVal.ToInteger())
			}
		}

		ppt, err := presentation.Open(filePath)
		if err != nil {
			sb.logger.Error("打开PPT文件失败", zap.String("path", filePath), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer ppt.Close()

		var slides []map[string]interface{}
		for i, slide := range ppt.Slides() {
			var text strings.Builder
			for _, para := range slide.Paragraphs() {
				text.WriteString(para.Text())
				text.WriteString("\n")
			}
			slides = append(slides, map[string]interface{}{
				"index": i + 1,
				"text":  text.String(),
			})
		}

		totalSlides := len(slides)
		totalPages := (totalSlides + pageSize - 1) / pageSize

		start := (page - 1) * pageSize
		end := start + pageSize
		if end > totalSlides {
			end = totalSlides
		}

		var pageSlides []map[string]interface{}
		if start < totalSlides {
			pageSlides = slides[start:end]
		}

		return sb.vm.ToValue(map[string]interface{}{
			"slides":     pageSlides,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
			"totalSlides": totalSlides,
		})
	})

	// 读取PDF文件内容（支持分页）
	sb.vm.Set("readPDF", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文件路径",
			})
		}

		filePath := call.Arguments[0].String()
		page := 1
		pageSize := 1 // 每页PDF页数

		if len(call.Arguments) > 1 {
			options := call.Arguments[1].ToObject(sb.vm)
			if pageVal := options.Get("page"); pageVal != nil && !goja.IsUndefined(pageVal) {
				page = int(pageVal.ToInteger())
			}
			if pageSizeVal := options.Get("pageSize"); pageSizeVal != nil && !goja.IsUndefined(pageSizeVal) {
				pageSize = int(pageSizeVal.ToInteger())
			}
		}

		f, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer f.Close()

		pdfReader, err := model.NewPdfReader(f)
		if err != nil {
			sb.logger.Error("读取PDF文件失败", zap.String("path", filePath), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		totalPages := (int(numPages) + pageSize - 1) / pageSize

		startPage := (page - 1) * pageSize + 1
		endPage := startPage + pageSize - 1
		if endPage > int(numPages) {
			endPage = int(numPages)
		}

		var pages []map[string]interface{}
		for i := startPage; i <= endPage; i++ {
			pageObj, err := pdfReader.GetPage(i)
			if err != nil {
				continue
			}

			// 提取文本内容
			text, err := extractTextFromPDFPage(pageObj)
			if err != nil {
				text = fmt.Sprintf("无法提取第%d页文本: %v", i, err)
			}

			pages = append(pages, map[string]interface{}{
				"page": i,
				"text": text,
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"pages":      pages,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
			"totalPDFPages": int(numPages),
		})
	})
}

// extractTextFromPDFPage 从PDF页面提取文本
func extractTextFromPDFPage(page *model.PdfPage) (string, error) {
	// 这是一个简化的实现，实际应该使用更完善的PDF文本提取库
	// 这里返回一个占位符，实际使用时需要实现真正的文本提取逻辑
	return fmt.Sprintf("PDF页面内容（需要实现文本提取）"), nil
}

