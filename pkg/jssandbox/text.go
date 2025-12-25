package jssandbox

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dop251/goja"
)

// registerText 注册文本操作功能到JavaScript运行时
func (sb *Sandbox) registerText() {
	// 文本内容替换
	sb.vm.Set("replaceText", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 3 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文本、查找字符串和替换字符串参数",
			})
		}

		text := call.Arguments[0].String()
		search := call.Arguments[1].String()
		replace := call.Arguments[2].String()

		// 支持全局替换（默认）或单次替换
		all := true
		if len(call.Arguments) > 3 {
			if allVal := call.Arguments[3]; !goja.IsUndefined(allVal) {
				all = allVal.ToBoolean()
			}
		}

		var result string
		if all {
			result = strings.ReplaceAll(text, search, replace)
		} else {
			result = strings.Replace(text, search, replace, 1)
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"result":  result,
		})
	})

	// 正则表达式替换
	sb.vm.Set("replaceRegex", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 3 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文本、正则表达式和替换字符串参数",
			})
		}

		text := call.Arguments[0].String()
		pattern := call.Arguments[1].String()
		replace := call.Arguments[2].String()

		re, err := regexp.Compile(pattern)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("正则表达式编译失败: %v", err),
			})
		}

		result := re.ReplaceAllString(text, replace)

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"result":  result,
		})
	})

	// 匹配 Markdown 标题（H1-H6）
	sb.vm.Set("matchMarkdownHeaders", func(text string) goja.Value {
		// 匹配 # 标题格式
		headerRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
		lines := strings.Split(text, "\n")
		var headers []map[string]interface{}

		for i, line := range lines {
			matches := headerRegex.FindStringSubmatch(line)
			if matches != nil {
				level := len(matches[1])
				content := strings.TrimSpace(matches[2])
				headers = append(headers, map[string]interface{}{
					"level":   level,
					"content": content,
					"line":    i + 1,
					"raw":     line,
				})
			}
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"headers": headers,
			"count":   len(headers),
		})
	})

	// 匹配指定级别的 Markdown 标题
	sb.vm.Set("matchMarkdownHeaderByLevel", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文本和标题级别参数",
			})
		}

		text := call.Arguments[0].String()
		level := int(call.Arguments[1].ToInteger())

		if level < 1 || level > 6 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "标题级别必须在 1-6 之间",
			})
		}

		pattern := fmt.Sprintf(`^(#{%d})\s+(.+)$`, level)
		headerRegex := regexp.MustCompile(pattern)
		lines := strings.Split(text, "\n")
		var headers []map[string]interface{}

		for i, line := range lines {
			matches := headerRegex.FindStringSubmatch(line)
			if matches != nil {
				content := strings.TrimSpace(matches[2])
				headers = append(headers, map[string]interface{}{
					"level":   level,
					"content": content,
					"line":    i + 1,
					"raw":     line,
				})
			}
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"headers": headers,
			"count":   len(headers),
		})
	})

	// 匹配 Markdown 中的所有图片
	sb.vm.Set("matchMarkdownImages", func(text string) goja.Value {
		// 匹配 ![alt](url) 和 ![alt](url "title") 格式
		// 使用一个正则表达式匹配所有图片，然后检查是否有标题
		// 匹配格式：![alt](url) 或 ![alt](url "title")
		// 使用非贪婪匹配确保正确解析
		imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+?)(?:\s+"([^"]+)")?\)`)
		matches := imageRegex.FindAllStringSubmatch(text, -1)

		var images []map[string]interface{}
		for _, match := range matches {
			alt := match[1]
			url := strings.TrimSpace(match[2]) // URL 可能包含末尾空格
			title := ""
			if len(match) > 3 && match[3] != "" {
				title = match[3]
			}

			images = append(images, map[string]interface{}{
				"alt":   alt,
				"url":   url,
				"title": title,
				"raw":   match[0],
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"images":  images,
			"count":   len(images),
		})
	})

	// 匹配 Markdown 中的代码块
	sb.vm.Set("matchMarkdownCodeBlocks", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文本参数",
			})
		}

		text := call.Arguments[0].String()
		includeInline := false
		if len(call.Arguments) > 1 {
			includeInline = call.Arguments[1].ToBoolean()
		}

		var codeBlocks []map[string]interface{}

		// 匹配围栏代码块 ```language\ncode\n```
		// 使用普通字符串，因为原始字符串字面量中不能直接使用反引号
		fencedRegex := regexp.MustCompile("(?s)```(\\w+)?\\n(.*?)```")
		fencedMatches := fencedRegex.FindAllStringSubmatch(text, -1)
		for _, match := range fencedMatches {
			language := ""
			if len(match) > 1 {
				language = match[1]
			}
			code := ""
			if len(match) > 2 {
				code = strings.TrimSpace(match[2])
			}

			codeBlocks = append(codeBlocks, map[string]interface{}{
				"type":     "fenced",
				"language": language,
				"code":     code,
				"raw":      match[0],
			})
		}

		// 如果包含行内代码
		if includeInline {
			// 匹配行内代码 `code`，但排除围栏代码块中的内容
			// 先移除已匹配的围栏代码块，避免匹配到围栏代码块内的反引号
			textForInline := text
			for _, match := range fencedMatches {
				textForInline = strings.Replace(textForInline, match[0], "", 1)
			}

			// 匹配行内代码 `code`
			inlineRegex := regexp.MustCompile("`([^`]+)`")
			inlineMatches := inlineRegex.FindAllStringSubmatch(textForInline, -1)
			for _, match := range inlineMatches {
				if len(match) > 1 {
					codeBlocks = append(codeBlocks, map[string]interface{}{
						"type": "inline",
						"code": match[1],
						"raw":  match[0],
					})
				}
			}
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success":    true,
			"codeBlocks": codeBlocks,
			"count":      len(codeBlocks),
		})
	})

	// 模板替换（支持 {{key}} 格式）
	sb.vm.Set("replaceTemplate", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供模板文本和数据对象参数",
			})
		}

		template := call.Arguments[0].String()
		var data map[string]interface{}

		if dataVal := call.Arguments[1]; !goja.IsUndefined(dataVal) {
			if dataObj := dataVal.ToObject(sb.vm); dataObj != nil {
				data = dataObj.Export().(map[string]interface{})
			} else {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "数据参数必须是对象",
				})
			}
		}

		// 匹配 {{key}} 或 {{ key }} 格式
		templateRegex := regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)
		result := templateRegex.ReplaceAllStringFunc(template, func(match string) string {
			keyMatch := templateRegex.FindStringSubmatch(match)
			if keyMatch != nil && len(keyMatch) > 1 {
				key := strings.TrimSpace(keyMatch[1])
				if value, ok := data[key]; ok {
					return fmt.Sprintf("%v", value)
				}
				// 如果键不存在，返回原始匹配（可选：也可以返回空字符串）
				return match
			}
			return match
		})

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"result":  result,
		})
	})

	// 高级模板替换（支持自定义分隔符和默认值）
	sb.vm.Set("replaceTemplateAdvanced", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供模板文本和数据对象参数",
			})
		}

		template := call.Arguments[0].String()
		var data map[string]interface{}

		if dataVal := call.Arguments[1]; !goja.IsUndefined(dataVal) {
			if dataObj := dataVal.ToObject(sb.vm); dataObj != nil {
				data = dataObj.Export().(map[string]interface{})
			} else {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "数据参数必须是对象",
				})
			}
		}

		// 自定义分隔符（可选）
		leftDelim := "{{"
		rightDelim := "}}"
		if len(call.Arguments) > 2 {
			leftDelim = call.Arguments[2].String()
		}
		if len(call.Arguments) > 3 {
			rightDelim = call.Arguments[3].String()
		}

		// 转义分隔符以便在正则中使用
		leftEscaped := regexp.QuoteMeta(leftDelim)
		rightEscaped := regexp.QuoteMeta(rightDelim)

		// 匹配 {{key}} 或 {{key|default}} 格式
		pattern := fmt.Sprintf(`%s\s*(\w+)(?:\s*\|\s*([^%s]+))?\s*%s`, leftEscaped, rightEscaped, rightEscaped)
		templateRegex := regexp.MustCompile(pattern)

		result := templateRegex.ReplaceAllStringFunc(template, func(match string) string {
			matches := templateRegex.FindStringSubmatch(match)
			if matches != nil && len(matches) > 1 {
				key := strings.TrimSpace(matches[1])
				defaultValue := ""
				if len(matches) > 2 {
					defaultValue = strings.TrimSpace(matches[2])
				}

				if value, ok := data[key]; ok {
					return fmt.Sprintf("%v", value)
				}
				// 如果键不存在，返回默认值或原始匹配
				if defaultValue != "" {
					return defaultValue
				}
				return match
			}
			return match
		})

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"result":  result,
		})
	})

	// 提取 Markdown 标题结构（树形结构）
	sb.vm.Set("extractMarkdownStructure", func(text string) goja.Value {
		lines := strings.Split(text, "\n")
		headerRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

		type HeaderNode struct {
			Level    int
			Content  string
			Line     int
			Children []*HeaderNode
		}

		var root []*HeaderNode
		var stack []*HeaderNode

		for i, line := range lines {
			matches := headerRegex.FindStringSubmatch(line)
			if matches != nil {
				level := len(matches[1])
				content := strings.TrimSpace(matches[2])
				node := &HeaderNode{
					Level:    level,
					Content:  content,
					Line:     i + 1,
					Children: []*HeaderNode{},
				}

				// 找到合适的父节点
				for len(stack) > 0 && stack[len(stack)-1].Level >= level {
					stack = stack[:len(stack)-1]
				}

				if len(stack) == 0 {
					root = append(root, node)
				} else {
					parent := stack[len(stack)-1]
					parent.Children = append(parent.Children, node)
				}

				stack = append(stack, node)
			}
		}

		// 转换为可导出的格式
		var convertNode func(*HeaderNode) map[string]interface{}
		convertNode = func(node *HeaderNode) map[string]interface{} {
			children := make([]map[string]interface{}, len(node.Children))
			for i, child := range node.Children {
				children[i] = convertNode(child)
			}
			return map[string]interface{}{
				"level":    node.Level,
				"content":  node.Content,
				"line":     node.Line,
				"children": children,
			}
		}

		structure := make([]map[string]interface{}, len(root))
		for i, node := range root {
			structure[i] = convertNode(node)
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success":   true,
			"structure": structure,
		})
	})
}
