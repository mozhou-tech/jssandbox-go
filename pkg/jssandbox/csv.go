package jssandbox

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/dop251/goja"
)

// registerCSV 注册CSV处理功能到JavaScript运行时
func (sb *Sandbox) registerCSV() {
	// 读取CSV文件
	sb.vm.Set("readCSV", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文件路径参数",
			})
		}

		filePath := call.Arguments[0].String()
		delimiter := ','
		comment := '#'
		skipEmptyLines := false

		// 解析选项
		if len(call.Arguments) > 1 {
			if options := call.Arguments[1].ToObject(sb.vm); options != nil {
				if delimVal := options.Get("delimiter"); delimVal != nil && !goja.IsUndefined(delimVal) {
					delimStr := delimVal.String()
					if len(delimStr) > 0 {
						delimiter = rune(delimStr[0])
					}
				}
				if commentVal := options.Get("comment"); commentVal != nil && !goja.IsUndefined(commentVal) {
					commentStr := commentVal.String()
					if len(commentStr) > 0 {
						comment = rune(commentStr[0])
					}
				}
				if skipVal := options.Get("skipEmptyLines"); skipVal != nil && !goja.IsUndefined(skipVal) {
					skipEmptyLines = skipVal.ToBoolean()
				}
			}
		}

		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("打开文件失败: %v", err),
			})
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.Comma = delimiter
		reader.Comment = comment

		records, err := reader.ReadAll()
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("读取CSV失败: %v", err),
			})
		}

		// 转换为数组
		var rows [][]string
		for _, record := range records {
			if skipEmptyLines {
				// 检查是否为空行
				empty := true
				for _, field := range record {
					if strings.TrimSpace(field) != "" {
						empty = false
						break
					}
				}
				if empty {
					continue
				}
			}
			rows = append(rows, record)
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"rows":    rows,
			"count":   len(rows),
		})
	})

	// 写入CSV文件
	sb.vm.Set("writeCSV", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文件路径和数据参数",
			})
		}

		filePath := call.Arguments[0].String()
		delimiter := ','

		// 解析数据
		var rows [][]string
		if dataVal := call.Arguments[1]; dataVal != nil && !goja.IsUndefined(dataVal) {
			if dataObj := dataVal.ToObject(sb.vm); dataObj != nil {
				if dataArray, ok := dataObj.Export().([]interface{}); ok {
					for _, row := range dataArray {
						if rowArray, ok := row.([]interface{}); ok {
							var strRow []string
							for _, cell := range rowArray {
								strRow = append(strRow, fmt.Sprintf("%v", cell))
							}
							rows = append(rows, strRow)
						}
					}
				}
			}
		}

		// 解析选项
		if len(call.Arguments) > 2 {
			if options := call.Arguments[2].ToObject(sb.vm); options != nil {
				if delimVal := options.Get("delimiter"); delimVal != nil && !goja.IsUndefined(delimVal) {
					delimStr := delimVal.String()
					if len(delimStr) > 0 {
						delimiter = rune(delimStr[0])
					}
				}
			}
		}

		file, err := os.Create(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("创建文件失败: %v", err),
			})
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		writer.Comma = delimiter

		if err := writer.WriteAll(rows); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("写入CSV失败: %v", err),
			})
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("刷新CSV失败: %v", err),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    filePath,
		})
	})

	// 解析CSV字符串
	sb.vm.Set("parseCSV", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供CSV字符串参数",
			})
		}

		csvString := call.Arguments[0].String()
		delimiter := ','

		// 解析选项
		if len(call.Arguments) > 1 {
			if options := call.Arguments[1].ToObject(sb.vm); options != nil {
				if delimVal := options.Get("delimiter"); delimVal != nil && !goja.IsUndefined(delimVal) {
					delimStr := delimVal.String()
					if len(delimStr) > 0 {
						delimiter = rune(delimStr[0])
					}
				}
			}
		}

		reader := csv.NewReader(strings.NewReader(csvString))
		reader.Comma = delimiter

		records, err := reader.ReadAll()
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("解析CSV失败: %v", err),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"rows":    records,
			"count":   len(records),
		})
	})
}
