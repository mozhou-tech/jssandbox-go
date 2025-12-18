package jssandbox

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
)

// registerGoQuery 注册 goquery HTML 解析功能到 JavaScript 运行时
func (sb *Sandbox) registerGoQuery() {
	// 创建文档对象
	sb.vm.Set("parseHTML", func(html string) goja.Value {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			sb.logger.WithError(err).Error("解析 HTML 失败")
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		// 创建一个 JavaScript 对象来表示文档
		docObj := sb.vm.NewObject()
		docObj.Set("_doc", doc) // 内部存储 goquery 文档对象

		// find 方法：查找匹配选择器的元素
		docObj.Set("find", func(selector string) goja.Value {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "文档对象无效",
				})
			}

			// 从 JavaScript 对象中获取 goquery 文档或选择器
			docInterface := docVal.Export()
			var selection *goquery.Selection
			
			if doc, ok := docInterface.(*goquery.Document); ok {
				selection = doc.Find(selector)
			} else if sel, ok := docInterface.(*goquery.Selection); ok {
				selection = sel.Find(selector)
			} else {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "无法获取文档对象",
				})
			}

			return createSelectionObject(sb, selection)
		})

		// text 方法：获取所有匹配元素的文本内容
		docObj.Set("text", func() string {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return ""
			}

			docInterface := docVal.Export()
			if doc, ok := docInterface.(*goquery.Document); ok {
				return doc.Text()
			}
			if sel, ok := docInterface.(*goquery.Selection); ok {
				return sel.Text()
			}
			return ""
		})

		// html 方法：获取第一个匹配元素的 HTML
		docObj.Set("html", func() string {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return ""
			}

			docInterface := docVal.Export()
			if doc, ok := docInterface.(*goquery.Document); ok {
				html, _ := doc.Html()
				return html
			}
			if sel, ok := docInterface.(*goquery.Selection); ok {
				html, _ := sel.Html()
				return html
			}
			return ""
		})

		// attr 方法：获取第一个匹配元素的属性值
		docObj.Set("attr", func(name string) string {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return ""
			}

			docInterface := docVal.Export()
			if sel, ok := docInterface.(*goquery.Selection); ok {
				val, _ := sel.Attr(name)
				return val
			}
			return ""
		})

		// each 方法：遍历所有匹配的元素
		docObj.Set("each", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) < 1 {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "需要提供回调函数",
				})
			}

			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "文档对象无效",
				})
			}

			docInterface := docVal.Export()
			sel, ok := docInterface.(*goquery.Selection)
			if !ok {
				if doc, ok := docInterface.(*goquery.Document); ok {
					sel = doc.Selection
				} else {
					return sb.vm.ToValue(map[string]interface{}{
						"error": "无法获取选择器对象",
					})
				}
			}

			callback := call.Arguments[0]
			callable, ok := goja.AssertFunction(callback)
			if !ok {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "参数必须是函数",
				})
			}

			var results []interface{}
			sel.Each(func(i int, s *goquery.Selection) {
				selObj := createSelectionObject(sb, s)
				result, err := callable(goja.Undefined(), selObj, sb.vm.ToValue(i))
				if err == nil && result != nil && !goja.IsUndefined(result) {
					results = append(results, result.Export())
				}
			})

			return sb.vm.ToValue(results)
		})

		// length 属性：获取匹配元素的数量
		docObj.Set("length", func() int {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return 0
			}

			docInterface := docVal.Export()
			if sel, ok := docInterface.(*goquery.Selection); ok {
				return sel.Length()
			}
			if doc, ok := docInterface.(*goquery.Document); ok {
				return doc.Selection.Length()
			}
			return 0
		})

		// first 方法：获取第一个匹配的元素
		docObj.Set("first", func() goja.Value {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "文档对象无效",
				})
			}

			docInterface := docVal.Export()
			var sel *goquery.Selection
			if s, ok := docInterface.(*goquery.Selection); ok {
				sel = s.First()
			} else if doc, ok := docInterface.(*goquery.Document); ok {
				sel = doc.Selection.First()
			} else {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "无法获取选择器对象",
				})
			}

			return createSelectionObject(sb, sel)
		})

		// last 方法：获取最后一个匹配的元素
		docObj.Set("last", func() goja.Value {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "文档对象无效",
				})
			}

			docInterface := docVal.Export()
			var sel *goquery.Selection
			if s, ok := docInterface.(*goquery.Selection); ok {
				sel = s.Last()
			} else if doc, ok := docInterface.(*goquery.Document); ok {
				sel = doc.Selection.Last()
			} else {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "无法获取选择器对象",
				})
			}

			return createSelectionObject(sb, sel)
		})

		// eq 方法：获取指定索引的元素
		docObj.Set("eq", func(index int) goja.Value {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "文档对象无效",
				})
			}

			docInterface := docVal.Export()
			var sel *goquery.Selection
			if s, ok := docInterface.(*goquery.Selection); ok {
				sel = s.Eq(index)
			} else if doc, ok := docInterface.(*goquery.Document); ok {
				sel = doc.Selection.Eq(index)
			} else {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "无法获取选择器对象",
				})
			}

			return createSelectionObject(sb, sel)
		})

		// children 方法：获取所有子元素
		docObj.Set("children", func(selector ...string) goja.Value {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "文档对象无效",
				})
			}

			docInterface := docVal.Export()
			var sel *goquery.Selection
			if s, ok := docInterface.(*goquery.Selection); ok {
				if len(selector) > 0 {
					sel = s.ChildrenFiltered(selector[0])
				} else {
					sel = s.Children()
				}
			} else if doc, ok := docInterface.(*goquery.Document); ok {
				if len(selector) > 0 {
					sel = doc.Selection.ChildrenFiltered(selector[0])
				} else {
					sel = doc.Selection.Children()
				}
			} else {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "无法获取选择器对象",
				})
			}

			return createSelectionObject(sb, sel)
		})

		// parent 方法：获取父元素
		docObj.Set("parent", func() goja.Value {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "文档对象无效",
				})
			}

			docInterface := docVal.Export()
			var sel *goquery.Selection
			if s, ok := docInterface.(*goquery.Selection); ok {
				sel = s.Parent()
			} else if doc, ok := docInterface.(*goquery.Document); ok {
				sel = doc.Selection.Parent()
			} else {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "无法获取选择器对象",
				})
			}

			return createSelectionObject(sb, sel)
		})

		// siblings 方法：获取所有兄弟元素
		docObj.Set("siblings", func(selector ...string) goja.Value {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "文档对象无效",
				})
			}

			docInterface := docVal.Export()
			var sel *goquery.Selection
			if s, ok := docInterface.(*goquery.Selection); ok {
				if len(selector) > 0 {
					sel = s.SiblingsFiltered(selector[0])
				} else {
					sel = s.Siblings()
				}
			} else if doc, ok := docInterface.(*goquery.Document); ok {
				if len(selector) > 0 {
					sel = doc.Selection.SiblingsFiltered(selector[0])
				} else {
					sel = doc.Selection.Siblings()
				}
			} else {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "无法获取选择器对象",
				})
			}

			return createSelectionObject(sb, sel)
		})

		// next 方法：获取下一个兄弟元素
		docObj.Set("next", func() goja.Value {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "文档对象无效",
				})
			}

			docInterface := docVal.Export()
			var sel *goquery.Selection
			if s, ok := docInterface.(*goquery.Selection); ok {
				sel = s.Next()
			} else if doc, ok := docInterface.(*goquery.Document); ok {
				sel = doc.Selection.Next()
			} else {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "无法获取选择器对象",
				})
			}

			return createSelectionObject(sb, sel)
		})

		// prev 方法：获取上一个兄弟元素
		docObj.Set("prev", func() goja.Value {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "文档对象无效",
				})
			}

			docInterface := docVal.Export()
			var sel *goquery.Selection
			if s, ok := docInterface.(*goquery.Selection); ok {
				sel = s.Prev()
			} else if doc, ok := docInterface.(*goquery.Document); ok {
				sel = doc.Selection.Prev()
			} else {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "无法获取选择器对象",
				})
			}

			return createSelectionObject(sb, sel)
		})

		// hasClass 方法：检查元素是否有指定的类
		docObj.Set("hasClass", func(className string) bool {
			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return false
			}

			docInterface := docVal.Export()
			if sel, ok := docInterface.(*goquery.Selection); ok {
				return sel.HasClass(className)
			}
			return false
		})

		// map 方法：将每个元素映射为值
		docObj.Set("map", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) < 1 {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "需要提供回调函数",
				})
			}

			docVal := docObj.Get("_doc")
			if docVal == nil || goja.IsUndefined(docVal) {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "文档对象无效",
				})
			}

			docInterface := docVal.Export()
			sel, ok := docInterface.(*goquery.Selection)
			if !ok {
				if doc, ok := docInterface.(*goquery.Document); ok {
					sel = doc.Selection
				} else {
					return sb.vm.ToValue(map[string]interface{}{
						"error": "无法获取选择器对象",
					})
				}
			}

			callback := call.Arguments[0]
			callable, ok := goja.AssertFunction(callback)
			if !ok {
				return sb.vm.ToValue(map[string]interface{}{
					"error": "参数必须是函数",
				})
			}

			var results []interface{}
			sel.Each(func(i int, s *goquery.Selection) {
				selObj := createSelectionObject(sb, s)
				result, err := callable(goja.Undefined(), selObj, sb.vm.ToValue(i))
				if err == nil && result != nil && !goja.IsUndefined(result) {
					results = append(results, result.Export())
				}
			})

			return sb.vm.ToValue(results)
		})

		return docObj
	})
}

// createSelectionObject 创建一个表示 goquery.Selection 的 JavaScript 对象
func createSelectionObject(sb *Sandbox, sel *goquery.Selection) goja.Value {
	selObj := sb.vm.NewObject()
	selObj.Set("_sel", sel) // 内部存储 goquery Selection 对象

	// find 方法
	selObj.Set("find", func(selector string) goja.Value {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "选择器对象无效",
			})
		}

		selInterface := selVal.Export()
		sel, ok := selInterface.(*goquery.Selection)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无法获取选择器对象",
			})
		}

		newSel := sel.Find(selector)
		return createSelectionObject(sb, newSel)
	})

	// text 方法
	selObj.Set("text", func() string {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return ""
		}

		selInterface := selVal.Export()
		if sel, ok := selInterface.(*goquery.Selection); ok {
			return sel.Text()
		}
		return ""
	})

	// html 方法
	selObj.Set("html", func() string {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return ""
		}

		selInterface := selVal.Export()
		if sel, ok := selInterface.(*goquery.Selection); ok {
			html, _ := sel.Html()
			return html
		}
		return ""
	})

	// attr 方法
	selObj.Set("attr", func(name string) string {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return ""
		}

		selInterface := selVal.Export()
		if sel, ok := selInterface.(*goquery.Selection); ok {
			val, _ := sel.Attr(name)
			return val
		}
		return ""
	})

	// each 方法
	selObj.Set("each", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供回调函数",
			})
		}

		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "选择器对象无效",
			})
		}

		selInterface := selVal.Export()
		sel, ok := selInterface.(*goquery.Selection)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无法获取选择器对象",
			})
		}

		callback := call.Arguments[0]
		callable, ok := goja.AssertFunction(callback)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "参数必须是函数",
			})
		}

		var results []interface{}
		sel.Each(func(i int, s *goquery.Selection) {
			selObj := createSelectionObject(sb, s)
			result, err := callable(goja.Undefined(), selObj, sb.vm.ToValue(i))
			if err == nil && result != nil && !goja.IsUndefined(result) {
				results = append(results, result.Export())
			}
		})

		return sb.vm.ToValue(results)
	})

	// length 属性
	selObj.Set("length", func() int {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return 0
		}

		selInterface := selVal.Export()
		if sel, ok := selInterface.(*goquery.Selection); ok {
			return sel.Length()
		}
		return 0
	})

	// first 方法
	selObj.Set("first", func() goja.Value {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "选择器对象无效",
			})
		}

		selInterface := selVal.Export()
		sel, ok := selInterface.(*goquery.Selection)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无法获取选择器对象",
			})
		}

		return createSelectionObject(sb, sel.First())
	})

	// last 方法
	selObj.Set("last", func() goja.Value {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "选择器对象无效",
			})
		}

		selInterface := selVal.Export()
		sel, ok := selInterface.(*goquery.Selection)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无法获取选择器对象",
			})
		}

		return createSelectionObject(sb, sel.Last())
	})

	// eq 方法
	selObj.Set("eq", func(index int) goja.Value {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "选择器对象无效",
			})
		}

		selInterface := selVal.Export()
		sel, ok := selInterface.(*goquery.Selection)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无法获取选择器对象",
			})
		}

		return createSelectionObject(sb, sel.Eq(index))
	})

	// children 方法
	selObj.Set("children", func(selector ...string) goja.Value {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "选择器对象无效",
			})
		}

		selInterface := selVal.Export()
		sel, ok := selInterface.(*goquery.Selection)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无法获取选择器对象",
			})
		}

		if len(selector) > 0 {
			return createSelectionObject(sb, sel.ChildrenFiltered(selector[0]))
		}
		return createSelectionObject(sb, sel.Children())
	})

	// parent 方法
	selObj.Set("parent", func() goja.Value {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "选择器对象无效",
			})
		}

		selInterface := selVal.Export()
		sel, ok := selInterface.(*goquery.Selection)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无法获取选择器对象",
			})
		}

		return createSelectionObject(sb, sel.Parent())
	})

	// siblings 方法
	selObj.Set("siblings", func(selector ...string) goja.Value {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "选择器对象无效",
			})
		}

		selInterface := selVal.Export()
		sel, ok := selInterface.(*goquery.Selection)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无法获取选择器对象",
			})
		}

		if len(selector) > 0 {
			return createSelectionObject(sb, sel.SiblingsFiltered(selector[0]))
		}
		return createSelectionObject(sb, sel.Siblings())
	})

	// next 方法
	selObj.Set("next", func() goja.Value {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "选择器对象无效",
			})
		}

		selInterface := selVal.Export()
		sel, ok := selInterface.(*goquery.Selection)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无法获取选择器对象",
			})
		}

		return createSelectionObject(sb, sel.Next())
	})

	// prev 方法
	selObj.Set("prev", func() goja.Value {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "选择器对象无效",
			})
		}

		selInterface := selVal.Export()
		sel, ok := selInterface.(*goquery.Selection)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无法获取选择器对象",
			})
		}

		return createSelectionObject(sb, sel.Prev())
	})

	// hasClass 方法
	selObj.Set("hasClass", func(className string) bool {
		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return false
		}

		selInterface := selVal.Export()
		if sel, ok := selInterface.(*goquery.Selection); ok {
			return sel.HasClass(className)
		}
		return false
	})

	// map 方法
	selObj.Set("map", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供回调函数",
			})
		}

		selVal := selObj.Get("_sel")
		if selVal == nil || goja.IsUndefined(selVal) {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "选择器对象无效",
			})
		}

		selInterface := selVal.Export()
		sel, ok := selInterface.(*goquery.Selection)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无法获取选择器对象",
			})
		}

		callback := call.Arguments[0]
		callable, ok := goja.AssertFunction(callback)
		if !ok {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "参数必须是函数",
			})
		}

		var results []interface{}
		sel.Each(func(i int, s *goquery.Selection) {
			selObj := createSelectionObject(sb, s)
			result, err := callable(goja.Undefined(), selObj, sb.vm.ToValue(i))
			if err == nil && result != nil && !goja.IsUndefined(result) {
				results = append(results, result.Export())
			}
		})

		return sb.vm.ToValue(results)
	})

	return selObj
}

