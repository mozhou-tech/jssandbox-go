package jssandbox

import (
	"context"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func init() {
	// 设置 chromedp 使用的 logrus logger 级别为 Fatal，抑制解析错误和警告
	// 这些错误（如 "could not unmarshal event" 和 "unknown PrivateNetworkRequestPolicy value"）
	// 通常来自 Chrome DevTools Protocol 事件解析，不影响功能，但会在日志中产生噪音
	logrus.SetLevel(logrus.FatalLevel)
}

// registerBrowser 注册浏览器操作功能到JavaScript运行时
func (sb *Sandbox) registerBrowser() {
	sb.vm.Set("browserNavigate", func(url string) goja.Value {
		ctx, cancel := chromedp.NewContext(sb.ctx)
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var html string
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			chromedp.OuterHTML("html", &html),
		)

		if err != nil {
			sb.logger.Error("浏览器导航失败", zap.String("url", url), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"html":    html,
		})
	})

	sb.vm.Set("browserScreenshot", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供URL参数",
			})
		}

		url := call.Arguments[0].String()
		outputPath := "screenshot.png"
		if len(call.Arguments) > 1 {
			outputPath = call.Arguments[1].String()
		}

		ctx, cancel := chromedp.NewContext(sb.ctx)
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var buf []byte
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			chromedp.CaptureScreenshot(&buf),
		)

		if err != nil {
			sb.logger.Error("截图失败", zap.String("url", url), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		// 保存截图
		err = os.WriteFile(outputPath, buf, 0644)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    outputPath,
		})
	})

	sb.vm.Set("browserEvaluate", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供URL和JavaScript代码",
			})
		}

		url := call.Arguments[0].String()
		jsCode := call.Arguments[1].String()

		ctx, cancel := chromedp.NewContext(sb.ctx)
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var result interface{}
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			chromedp.Evaluate(jsCode, &result),
		)

		if err != nil {
			sb.logger.Error("执行浏览器脚本失败", zap.String("url", url), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"result":  result,
		})
	})

	sb.vm.Set("browserClick", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供URL和选择器",
			})
		}

		url := call.Arguments[0].String()
		selector := call.Arguments[1].String()

		ctx, cancel := chromedp.NewContext(sb.ctx)
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(selector, chromedp.ByQuery),
			chromedp.Click(selector, chromedp.ByQuery),
		)

		if err != nil {
			sb.logger.Error("点击元素失败", zap.String("url", url), zap.String("selector", selector), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
		})
	})

	sb.vm.Set("browserFill", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 3 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供URL、选择器和值",
			})
		}

		url := call.Arguments[0].String()
		selector := call.Arguments[1].String()
		value := call.Arguments[2].String()

		ctx, cancel := chromedp.NewContext(sb.ctx)
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(selector, chromedp.ByQuery),
			chromedp.SendKeys(selector, value, chromedp.ByQuery),
		)

		if err != nil {
			sb.logger.Error("填充表单失败", zap.String("url", url), zap.String("selector", selector), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
		})
	})
}
