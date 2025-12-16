package jssandbox

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/dop251/goja"
	"go.uber.org/zap"
)

// registerHTTP 注册HTTP请求功能到JavaScript运行时
func (sb *Sandbox) registerHTTP() {
	sb.vm.Set("httpRequest", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供URL参数",
			})
		}

		url := call.Arguments[0].String()
		method := "GET"
		headers := make(map[string]string)
		body := ""
		timeout := 30

		if len(call.Arguments) > 1 {
			options := call.Arguments[1].ToObject(sb.vm)
			if methodVal := options.Get("method"); methodVal != nil && !goja.IsUndefined(methodVal) {
				method = methodVal.String()
			}
			if headersVal := options.Get("headers"); headersVal != nil && !goja.IsUndefined(headersVal) {
				headersObj := headersVal.ToObject(sb.vm)
				for _, key := range headersObj.Keys() {
					headers[key] = headersObj.Get(key).String()
				}
			}
			if bodyVal := options.Get("body"); bodyVal != nil && !goja.IsUndefined(bodyVal) {
				body = bodyVal.String()
			}
			if timeoutVal := options.Get("timeout"); timeoutVal != nil && !goja.IsUndefined(timeoutVal) {
				timeout = int(timeoutVal.ToInteger())
			}
		}

		client := &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}

		var reqBody io.Reader
		if body != "" {
			reqBody = bytes.NewBufferString(body)
		}

		req, err := http.NewRequest(method, url, reqBody)
		if err != nil {
			sb.logger.Error("创建HTTP请求失败", zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			sb.logger.Error("执行HTTP请求失败", zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			sb.logger.Error("读取响应体失败", zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"status":  resp.StatusCode,
				"headers": resp.Header,
				"error":   err.Error(),
			})
		}

		// 构建响应头对象
		respHeaders := make(map[string]string)
		for k, v := range resp.Header {
			if len(v) > 0 {
				respHeaders[k] = v[0]
			}
		}

		return sb.vm.ToValue(map[string]interface{}{
			"status":      resp.StatusCode,
			"statusText":  resp.Status,
			"headers":     respHeaders,
			"body":        string(respBody),
			"contentType": resp.Header.Get("Content-Type"),
		})
	})

	// 便捷方法
	sb.vm.Set("httpGet", func(url string) goja.Value {
		httpRequestVal := sb.vm.Get("httpRequest")
		if callable, ok := goja.AssertFunction(httpRequestVal); ok {
			result, _ := callable(goja.Undefined(), sb.vm.ToValue(url), sb.vm.ToValue(map[string]interface{}{
				"method": "GET",
			}))
			return result
		}
		return sb.vm.ToValue(map[string]interface{}{
			"error": "httpRequest 不是一个函数",
		})
	})

	sb.vm.Set("httpPost", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供URL参数",
			})
		}
		url := call.Arguments[0].String()
		body := ""
		if len(call.Arguments) > 1 {
			body = call.Arguments[1].String()
		}
		httpRequestVal := sb.vm.Get("httpRequest")
		if callable, ok := goja.AssertFunction(httpRequestVal); ok {
			result, _ := callable(goja.Undefined(), sb.vm.ToValue(url), sb.vm.ToValue(map[string]interface{}{
				"method": "POST",
				"body":   body,
				"headers": map[string]string{
					"Content-Type": "application/json",
				},
			}))
			return result
		}
		return sb.vm.ToValue(map[string]interface{}{
			"error": "httpRequest 不是一个函数",
		})
	})
}
