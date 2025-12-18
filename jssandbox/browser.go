package jssandbox

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
)

// BrowserSession 表示一个浏览器会话，可以保持状态（如cookies）并执行多个连续操作
type BrowserSession struct {
	ctx     context.Context
	cancel  context.CancelFunc
	sb      *Sandbox
	mu      sync.Mutex
	closed  bool
	timeout time.Duration
}

func init() {
	// 设置 chromedp 使用的 logrus logger 级别为 Fatal，抑制解析错误和警告
	// 这些错误（如 "could not unmarshal event" 和 "unknown PrivateNetworkRequestPolicy value"）
	// 通常来自 Chrome DevTools Protocol 事件解析，不影响功能，但会在日志中产生噪音
	logrus.SetLevel(logrus.FatalLevel)
}

// getOrCreateBrowserAllocator 获取或创建共享的浏览器 allocator
// 这确保所有会话使用同一个浏览器进程，避免打开多个窗口
func (sb *Sandbox) getOrCreateBrowserAllocator() context.Context {
	sb.browserMu.Lock()
	defer sb.browserMu.Unlock()

	if sb.browserInit && sb.browserAllocator != nil {
		return sb.browserAllocator
	}

	// 创建 allocator 选项，配置反检测参数
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", sb.config.Headless),                   // 根据配置决定是否使用headless模式
		chromedp.Flag("disable-blink-features", "AutomationControlled"), // 隐藏自动化特征
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"), // 设置真实User-Agent
		chromedp.Flag("disable-dev-shm-usage", true),              // 避免共享内存问题
		chromedp.Flag("no-sandbox", true),                         // 在某些环境下需要
		chromedp.Flag("disable-setuid-sandbox", true),             // 禁用setuid沙箱
		chromedp.Flag("disable-web-security", false),              // 保持web安全
		chromedp.Flag("disable-features", "VizDisplayCompositor"), // 禁用某些可能暴露的特征
	)

	// 只在headless模式下禁用GPU
	if sb.config.Headless {
		opts = append(opts, chromedp.Flag("disable-gpu", true))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(sb.ctx, opts...)
	sb.browserAllocator = allocCtx
	sb.browserCancel = cancel
	sb.browserInit = true

	return allocCtx
}

// createBrowserContext 创建配置了反检测选项的浏览器上下文
// 包括隐藏自动化特征、设置真实User-Agent等
// 使用共享的 allocator，确保只打开一个浏览器窗口
func (sb *Sandbox) createBrowserContext() (context.Context, context.CancelFunc) {
	// 获取或创建共享的 allocator（浏览器进程）
	allocCtx := sb.getOrCreateBrowserAllocator()

	// 为每个会话创建新的 context（标签页），但共享同一个 allocator（浏览器进程）
	ctx, cancel := chromedp.NewContext(allocCtx)

	return ctx, cancel
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

// createBrowserSession 创建一个新的浏览器会话，可以保持状态并执行多个连续操作
func (sb *Sandbox) createBrowserSession(timeoutSeconds float64) *BrowserSession {
	timeout := 30 * time.Second
	if timeoutSeconds > 0 {
		timeout = time.Duration(timeoutSeconds * float64(time.Second))
	}

	ctx, cancel := sb.createBrowserContext()
	ctx, cancelTimeout := context.WithTimeout(ctx, timeout)

	session := &BrowserSession{
		ctx:     ctx,
		cancel:  func() { cancelTimeout(); cancel() },
		sb:      sb,
		timeout: timeout,
	}

	// 初始化浏览器，注入反检测脚本
	go func() {
		_ = chromedp.Run(ctx, injectStealthScript())
	}()

	return session
}

// Close 关闭浏览器会话并清理资源
func (bs *BrowserSession) Close() {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if !bs.closed {
		bs.closed = true
		bs.cancel()
	}
}

// Navigate 导航到指定URL
func (bs *BrowserSession) Navigate(url string) map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	// 先执行导航
	bs.sb.logger.WithField("url", url).Debug("开始导航到页面")

	// 创建导航超时上下文（30秒）
	navCtx, navCancel := context.WithTimeout(bs.ctx, 30*time.Second)
	defer navCancel()

	err := chromedp.Run(navCtx, chromedp.Navigate(url))
	if err != nil {
		bs.sb.logger.WithError(err).WithField("url", url).Error("浏览器导航失败")
		return map[string]interface{}{
			"success": false,
			"error":   "导航失败: " + err.Error(),
		}
	}
	bs.sb.logger.WithField("url", url).Debug("导航命令已执行，等待页面加载")

	// 等待页面加载完成 - 使用更可靠的策略
	// 1. 先等待 DOMContentLoaded 事件
	waitCtx, waitCancel := context.WithTimeout(bs.ctx, 30*time.Second)
	defer waitCancel()

	// 等待页面就绪（DOMContentLoaded）
	err = chromedp.Run(waitCtx,
		chromedp.WaitReady("body", chromedp.ByQuery),
	)
	if err != nil {
		bs.sb.logger.WithError(err).WithField("url", url).Warn("等待body元素超时，尝试继续")
	} else {
		bs.sb.logger.WithField("url", url).Debug("页面body元素已就绪")
	}

	// 2. 等待一小段时间让JavaScript执行
	err = chromedp.Run(bs.ctx, chromedp.Sleep(1*time.Second))
	if err != nil {
		bs.sb.logger.WithError(err).WithField("url", url).Warn("等待页面脚本执行时出错")
	}

	// 3. 尝试等待网络空闲（可选，如果页面还在加载资源）
	// 使用JavaScript检查页面加载状态
	networkIdleCtx, networkIdleCancel := context.WithTimeout(bs.ctx, 10*time.Second)
	defer networkIdleCancel()

	// 检查页面是否加载完成
	var pageReady bool
	err = chromedp.Run(networkIdleCtx,
		chromedp.Evaluate(`
			(function() {
				if (document.readyState === 'complete') {
					return true;
				}
				// 检查是否有正在进行的网络请求（通过performance API）
				if (window.performance && window.performance.getEntriesByType) {
					var entries = window.performance.getEntriesByType('resource');
					// 如果最近1秒内没有新的资源加载，认为页面已就绪
					var now = Date.now();
					var recentLoads = entries.filter(function(entry) {
						return (now - entry.responseEnd) < 1000;
					});
					return recentLoads.length === 0;
				}
				return document.readyState === 'interactive' || document.readyState === 'complete';
			})();
		`, &pageReady),
	)
	if err != nil {
		bs.sb.logger.WithError(err).WithField("url", url).Debug("检查页面就绪状态时出错，继续执行")
	} else if pageReady {
		bs.sb.logger.WithField("url", url).Debug("页面加载完成")
	} else {
		bs.sb.logger.WithField("url", url).Debug("页面可能仍在加载中，继续执行")
	}

	// 注入反检测脚本
	err = chromedp.Run(bs.ctx, injectStealthScript())
	if err != nil {
		bs.sb.logger.WithError(err).WithField("url", url).Warn("注入反检测脚本失败，但继续执行")
		// 注入脚本失败不影响导航结果
	}

	// 验证导航是否成功 - 检查当前URL
	var currentURL string
	err = chromedp.Run(bs.ctx, chromedp.Location(&currentURL))
	if err != nil {
		bs.sb.logger.WithError(err).WithField("url", url).Warn("无法获取当前URL")
		return map[string]interface{}{
			"success": false,
			"error":   "无法获取当前URL: " + err.Error(),
		}
	}

	bs.sb.logger.WithField("url", url).WithField("currentURL", currentURL).Debug("导航完成，当前URL")

	// 检查URL是否匹配（允许重定向）
	if currentURL == "" {
		bs.sb.logger.WithField("url", url).Error("当前URL为空")
		return map[string]interface{}{
			"success": false,
			"error":   "导航后URL为空",
		}
	}

	return map[string]interface{}{
		"success": true,
		"url":     currentURL,
	}
}

// Wait 等待元素出现或等待指定时间
func (bs *BrowserSession) Wait(selectorOrSeconds interface{}) map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	var err error
	switch v := selectorOrSeconds.(type) {
	case string:
		// 等待元素出现
		err = chromedp.Run(bs.ctx,
			chromedp.WaitVisible(v, chromedp.ByQuery),
		)
	case float64:
		// 等待指定秒数
		err = chromedp.Run(bs.ctx,
			chromedp.Sleep(time.Duration(v*float64(time.Second))),
		)
	default:
		return map[string]interface{}{
			"success": false,
			"error":   "参数类型错误，需要字符串（选择器）或数字（秒数）",
		}
	}

	if err != nil {
		bs.sb.logger.WithError(err).Error("等待操作失败")
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}

// Click 点击指定选择器的元素
func (bs *BrowserSession) Click(selector string) map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	err := chromedp.Run(bs.ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Click(selector, chromedp.ByQuery),
	)

	if err != nil {
		bs.sb.logger.WithError(err).WithField("selector", selector).Error("点击元素失败")
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}

// Fill 填充表单字段
func (bs *BrowserSession) Fill(selector, value string) map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	err := chromedp.Run(bs.ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Clear(selector, chromedp.ByQuery),
		chromedp.SendKeys(selector, value, chromedp.ByQuery),
	)

	if err != nil {
		bs.sb.logger.WithError(err).WithField("selector", selector).Error("填充表单失败")
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}

// Evaluate 在页面中执行JavaScript代码
func (bs *BrowserSession) Evaluate(jsCode string) map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	var result interface{}
	err := chromedp.Run(bs.ctx,
		chromedp.Evaluate(jsCode, &result),
	)

	if err != nil {
		bs.sb.logger.WithError(err).Error("执行浏览器脚本失败")
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
		"result":  result,
	}
}

// GetHTML 获取当前页面的HTML内容
func (bs *BrowserSession) GetHTML() map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	var html string
	err := chromedp.Run(bs.ctx,
		chromedp.OuterHTML("html", &html),
	)

	if err != nil {
		bs.sb.logger.WithError(err).Error("获取HTML失败")
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
		"html":    html,
	}
}

// Screenshot 截取当前页面截图
func (bs *BrowserSession) Screenshot(outputPath string) map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	var buf []byte
	err := chromedp.Run(bs.ctx,
		chromedp.CaptureScreenshot(&buf),
	)

	if err != nil {
		bs.sb.logger.WithError(err).Error("截图失败")
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	// 检查截图数据是否为空
	if len(buf) == 0 {
		bs.sb.logger.Error("截图数据为空")
		return map[string]interface{}{
			"success": false,
			"error":   "截图数据为空",
		}
	}

	// 确保输出目录存在
	dir := filepath.Dir(outputPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			bs.sb.logger.WithError(err).WithField("dir", dir).Error("创建目录失败")
			return map[string]interface{}{
				"success": false,
				"error":   "创建目录失败: " + err.Error(),
			}
		}
	}

	// 保存截图
	err = os.WriteFile(outputPath, buf, 0644)
	if err != nil {
		bs.sb.logger.WithError(err).WithField("path", outputPath).Error("保存截图文件失败")
		return map[string]interface{}{
			"success": false,
			"error":   "保存文件失败: " + err.Error(),
		}
	}

	// 验证文件是否真的被写入
	if fileInfo, err := os.Stat(outputPath); err != nil {
		bs.sb.logger.WithError(err).WithField("path", outputPath).Error("验证截图文件失败")
		return map[string]interface{}{
			"success": false,
			"error":   "文件写入后验证失败: " + err.Error(),
		}
	} else if fileInfo.Size() == 0 {
		bs.sb.logger.WithField("path", outputPath).Error("截图文件大小为0")
		return map[string]interface{}{
			"success": false,
			"error":   "截图文件大小为0",
		}
	}

	return map[string]interface{}{
		"success": true,
		"path":    outputPath,
	}
}

// GetURL 获取当前页面的URL
func (bs *BrowserSession) GetURL() map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	var url string
	err := chromedp.Run(bs.ctx,
		chromedp.Location(&url),
	)

	if err != nil {
		bs.sb.logger.WithError(err).Error("获取URL失败")
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
		"url":     url,
	}
}

// WaitForURL 等待URL包含指定文本或匹配指定模式
func (bs *BrowserSession) WaitForURL(pattern string, timeoutSeconds float64) map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	timeout := 10 * time.Second
	if timeoutSeconds > 0 {
		timeout = time.Duration(timeoutSeconds * float64(time.Second))
	}

	ctx, cancel := context.WithTimeout(bs.ctx, timeout)
	defer cancel()

	startTime := time.Now()
	var url string
	for {
		// 检查超时
		if time.Since(startTime) > timeout {
			return map[string]interface{}{
				"success": false,
				"error":   "等待URL超时",
				"url":     url,
			}
		}

		// 获取当前URL
		err := chromedp.Run(ctx, chromedp.Location(&url))
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}

		// 检查URL是否包含模式
		if len(url) > 0 && len(pattern) > 0 {
			// 使用简单的字符串包含检查
			matched := false
			for i := 0; i <= len(url)-len(pattern); i++ {
				if url[i:i+len(pattern)] == pattern {
					matched = true
					break
				}
			}
			if matched {
				return map[string]interface{}{
					"success": true,
					"url":     url,
				}
			}
		}

		// 等待一小段时间后重试
		time.Sleep(100 * time.Millisecond)
	}
}

// WaitForText 等待页面中出现指定文本
func (bs *BrowserSession) WaitForText(text string, timeoutSeconds float64) map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	timeout := 10 * time.Second
	if timeoutSeconds > 0 {
		timeout = time.Duration(timeoutSeconds * float64(time.Second))
	}

	ctx, cancel := context.WithTimeout(bs.ctx, timeout)
	defer cancel()

	// 转义文本中的单引号，避免JavaScript注入
	escapedText := ""
	for _, r := range text {
		if r == '\'' {
			escapedText += "\\'"
		} else if r == '\\' {
			escapedText += "\\\\"
		} else {
			escapedText += string(r)
		}
	}
	jsCode := `document.body && document.body.innerText && document.body.innerText.includes('` + escapedText + `')`

	startTime := time.Now()
	for {
		// 检查超时
		if time.Since(startTime) > timeout {
			return map[string]interface{}{
				"success": false,
				"error":   "等待文本超时: " + text,
			}
		}

		// 执行JavaScript检查文本是否存在
		var result bool
		err := chromedp.Run(ctx, chromedp.Evaluate(jsCode, &result))
		if err != nil {
			// 如果页面还没加载完成，继续等待
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if result {
			return map[string]interface{}{
				"success": true,
			}
		}

		// 等待一小段时间后重试
		time.Sleep(100 * time.Millisecond)
	}
}

// Clear 清空指定输入框的内容
func (bs *BrowserSession) Clear(selector string) map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	err := chromedp.Run(bs.ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Clear(selector, chromedp.ByQuery),
	)

	if err != nil {
		bs.sb.logger.WithError(err).WithField("selector", selector).Error("清空输入框失败")
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}

// Submit 提交表单（通过点击提交按钮或按Enter键）
func (bs *BrowserSession) Submit(selector string) map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.closed {
		return map[string]interface{}{
			"success": false,
			"error":   "会话已关闭",
		}
	}

	// 如果selector为空，尝试提交当前表单（通过JavaScript模拟Enter键）
	if selector == "" {
		jsCode := `
			var event = new KeyboardEvent('keydown', {
				key: 'Enter',
				code: 'Enter',
				keyCode: 13,
				which: 13,
				bubbles: true
			});
			document.activeElement && document.activeElement.dispatchEvent(event);
		`
		err := chromedp.Run(bs.ctx,
			chromedp.Evaluate(jsCode, nil),
		)
		if err != nil {
			bs.sb.logger.WithError(err).Error("提交表单失败")
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
	} else {
		// 点击提交按钮
		err := chromedp.Run(bs.ctx,
			chromedp.WaitVisible(selector, chromedp.ByQuery),
			chromedp.Click(selector, chromedp.ByQuery),
		)
		if err != nil {
			bs.sb.logger.WithError(err).WithField("selector", selector).Error("点击提交按钮失败")
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}

// registerBrowser 注册浏览器操作功能到JavaScript运行时
func (sb *Sandbox) registerBrowser() {
	// 注册浏览器会话管理功能
	sb.vm.Set("createBrowserSession", func(call goja.FunctionCall) goja.Value {
		timeoutSeconds := 30.0
		if len(call.Arguments) > 0 {
			timeoutSeconds = call.Arguments[0].ToFloat()
		}

		session := sb.createBrowserSession(timeoutSeconds)

		// 创建一个JavaScript对象来表示会话
		sessionObj := sb.vm.NewObject()
		sessionObj.Set("navigate", func(url string) goja.Value {
			result := session.Navigate(url)
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("wait", func(selectorOrSeconds goja.Value) goja.Value {
			var arg interface{}
			if selectorOrSeconds != nil && !goja.IsUndefined(selectorOrSeconds) {
				if selectorOrSeconds.ExportType().Kind().String() == "string" {
					arg = selectorOrSeconds.String()
				} else {
					arg = selectorOrSeconds.ToFloat()
				}
			}
			result := session.Wait(arg)
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("click", func(selector string) goja.Value {
			result := session.Click(selector)
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("fill", func(selector, value string) goja.Value {
			result := session.Fill(selector, value)
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("evaluate", func(jsCode string) goja.Value {
			result := session.Evaluate(jsCode)
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("getHTML", func() goja.Value {
			result := session.GetHTML()
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("screenshot", func(outputPath string) goja.Value {
			result := session.Screenshot(outputPath)
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("getURL", func() goja.Value {
			result := session.GetURL()
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("waitForURL", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) < 1 {
				return sb.vm.ToValue(map[string]interface{}{
					"success": false,
					"error":   "需要提供URL模式参数",
				})
			}
			pattern := call.Arguments[0].String()
			timeoutSeconds := 10.0
			if len(call.Arguments) > 1 {
				timeoutSeconds = call.Arguments[1].ToFloat()
			}
			result := session.WaitForURL(pattern, timeoutSeconds)
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("waitForText", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) < 1 {
				return sb.vm.ToValue(map[string]interface{}{
					"success": false,
					"error":   "需要提供文本参数",
				})
			}
			text := call.Arguments[0].String()
			timeoutSeconds := 10.0
			if len(call.Arguments) > 1 {
				timeoutSeconds = call.Arguments[1].ToFloat()
			}
			result := session.WaitForText(text, timeoutSeconds)
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("clear", func(selector string) goja.Value {
			result := session.Clear(selector)
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("submit", func(selector string) goja.Value {
			result := session.Submit(selector)
			return sb.vm.ToValue(result)
		})
		sessionObj.Set("close", func() {
			session.Close()
		})

		return sessionObj
	})
}
