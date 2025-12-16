package jssandbox

import (
	"context"
	"os"
	"path/filepath"
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

// createBrowserContext 创建配置了反检测选项的浏览器上下文
// 包括隐藏自动化特征、设置真实User-Agent等
func (sb *Sandbox) createBrowserContext() (context.Context, context.CancelFunc) {
	// 创建 allocator 选项，配置反检测参数
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),                                 // 保持headless模式
		chromedp.Flag("disable-blink-features", "AutomationControlled"), // 隐藏自动化特征
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"), // 设置真实User-Agent
		chromedp.Flag("disable-dev-shm-usage", true),              // 避免共享内存问题
		chromedp.Flag("no-sandbox", true),                         // 在某些环境下需要
		chromedp.Flag("disable-gpu", true),                        // 在headless模式下禁用GPU
		chromedp.Flag("disable-setuid-sandbox", true),             // 禁用setuid沙箱
		chromedp.Flag("disable-web-security", false),              // 保持web安全
		chromedp.Flag("disable-features", "VizDisplayCompositor"), // 禁用某些可能暴露的特征
	)

	allocCtx, cancel := chromedp.NewExecAllocator(sb.ctx, opts...)
	ctx, cancel2 := chromedp.NewContext(allocCtx)

	return ctx, func() {
		cancel2()
		cancel()
	}
}

// injectStealthScript 注入反检测脚本，隐藏webdriver特征
// 在页面加载后立即执行，修改navigator对象
func injectStealthScript() chromedp.Action {
	stealthScript := `
		// 隐藏 webdriver 特征
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined
		});
		
		// 添加 Chrome 对象
		if (!window.chrome) {
			window.chrome = {
				runtime: {}
			};
		}
		
		// 修改 navigator.plugins 使其看起来更真实
		if (navigator.plugins.length === 0) {
			Object.defineProperty(navigator, 'plugins', {
				get: () => [1, 2, 3, 4, 5]
			});
		}
		
		// 修改 navigator.languages
		Object.defineProperty(navigator, 'languages', {
			get: () => ['zh-CN', 'zh', 'en-US', 'en']
		});
		
		// 覆盖 permissions API
		if (window.navigator.permissions && window.navigator.permissions.query) {
			const originalQuery = window.navigator.permissions.query;
			window.navigator.permissions.query = (parameters) => (
				parameters.name === 'notifications' ?
					Promise.resolve({ state: Notification.permission }) :
					originalQuery(parameters)
			);
		}
	`
	return chromedp.Evaluate(stealthScript, nil)
}

// registerBrowser 注册浏览器操作功能到JavaScript运行时
func (sb *Sandbox) registerBrowser() {
	sb.vm.Set("browserNavigate", func(url string) goja.Value {
		ctx, cancel := sb.createBrowserContext()
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var html string
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			injectStealthScript(), // 在页面加载后注入反检测脚本
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

		// 可选的等待时间（秒）
		waitSeconds := 0.0
		if len(call.Arguments) > 2 {
			waitSeconds = call.Arguments[2].ToFloat()
		}

		ctx, cancel := sb.createBrowserContext()
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var buf []byte
		actions := []chromedp.Action{
			chromedp.Navigate(url),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			injectStealthScript(), // 在页面加载后注入反检测脚本
		}

		// 如果指定了等待时间，添加等待操作
		if waitSeconds > 0 {
			actions = append(actions, chromedp.Sleep(time.Duration(waitSeconds*float64(time.Second))))
		}

		actions = append(actions, chromedp.CaptureScreenshot(&buf))

		err := chromedp.Run(ctx, actions...)

		if err != nil {
			sb.logger.Error("截图失败", zap.String("url", url), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}

		// 检查截图数据是否为空
		if len(buf) == 0 {
			sb.logger.Error("截图数据为空", zap.String("url", url))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "截图数据为空",
			})
		}

		// 确保输出目录存在
		dir := filepath.Dir(outputPath)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				sb.logger.Error("创建目录失败", zap.String("dir", dir), zap.Error(err))
				return sb.vm.ToValue(map[string]interface{}{
					"success": false,
					"error":   "创建目录失败: " + err.Error(),
				})
			}
		}

		// 保存截图
		err = os.WriteFile(outputPath, buf, 0644)
		if err != nil {
			sb.logger.Error("保存截图文件失败", zap.String("path", outputPath), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "保存文件失败: " + err.Error(),
			})
		}

		// 验证文件是否真的被写入
		if fileInfo, err := os.Stat(outputPath); err != nil {
			sb.logger.Error("验证截图文件失败", zap.String("path", outputPath), zap.Error(err))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "文件写入后验证失败: " + err.Error(),
			})
		} else if fileInfo.Size() == 0 {
			sb.logger.Error("截图文件大小为0", zap.String("path", outputPath))
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   "截图文件大小为0",
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

		ctx, cancel := sb.createBrowserContext()
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var result interface{}
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			injectStealthScript(), // 在页面加载后注入反检测脚本
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

		ctx, cancel := sb.createBrowserContext()
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			injectStealthScript(), // 在页面加载后注入反检测脚本
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

		ctx, cancel := sb.createBrowserContext()
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			injectStealthScript(), // 在页面加载后注入反检测脚本
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
