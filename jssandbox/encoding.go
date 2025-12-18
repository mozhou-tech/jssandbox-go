package jssandbox

import (
	"encoding/base64"
	"html"
	"net/url"

	"github.com/dop251/goja"
)

// registerEncoding 注册编码/解码增强功能到JavaScript运行时
func (sb *Sandbox) registerEncoding() {
	// Base64编码（增强版，支持文件）
	sb.vm.Set("encodeBase64", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供数据参数",
			})
		}

		var data string
		if arg := call.Arguments[0]; arg != nil && !goja.IsUndefined(arg) {
			data = arg.String()
		}

		encoded := base64.StdEncoding.EncodeToString([]byte(data))

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"data":    encoded,
		})
	})

	// Base64解码
	sb.vm.Set("decodeBase64", func(encoded string) goja.Value {
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "Base64解码失败: " + err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"data":    string(decoded),
		})
	})

	// URL编码
	sb.vm.Set("encodeURL", func(str string) goja.Value {
		encoded := url.QueryEscape(str)
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"data":    encoded,
		})
	})

	// URL解码
	sb.vm.Set("decodeURL", func(encoded string) goja.Value {
		decoded, err := url.QueryUnescape(encoded)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "URL解码失败: " + err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"data":    decoded,
		})
	})

	// HTML实体编码
	sb.vm.Set("encodeHTML", func(str string) goja.Value {
		encoded := html.EscapeString(str)
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"data":    encoded,
		})
	})

	// HTML实体解码
	sb.vm.Set("decodeHTML", func(encoded string) goja.Value {
		decoded := html.UnescapeString(encoded)
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"data":    decoded,
		})
	})
}
