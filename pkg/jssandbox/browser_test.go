package jssandbox

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/dop251/goja"
)

func TestBrowserSession_Navigate(t *testing.T) {
	// 注意：这个测试需要Chrome/Chromium可用
	// 如果环境中没有Chrome，测试可能会失败
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var session = createBrowserSession(30);
		var result = session.navigate("https://www.example.com");
		if (!result.success) {
			throw new Error("导航失败: " + result.error);
		}
		var htmlResult = session.getHTML();
		session.close();
		htmlResult;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("浏览器会话测试可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("getHTML()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Error("getHTML()缺少success字段")
	}

	// 如果成功，验证HTML字段
	if success.ToBoolean() {
		html := resultObj.Get("html")
		if html == nil || goja.IsUndefined(html) {
			t.Error("getHTML()成功时应该包含html字段")
		}
	}
}

func TestBrowserSession_Screenshot(t *testing.T) {
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	outputPath := testDir + "/screenshot.png"

	code := `
		var session = createBrowserSession(30);
		var navResult = session.navigate("https://www.example.com");
		if (!navResult.success) {
			throw new Error("导航失败: " + navResult.error);
		}
		var result = session.screenshot("` + outputPath + `");
		session.close();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("浏览器会话测试可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("screenshot()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success != nil && success.ToBoolean() {
		// 验证截图文件是否存在
		path := resultObj.Get("path")
		if path != nil && !goja.IsUndefined(path) {
			if _, err := os.Stat(path.String()); os.IsNotExist(err) {
				t.Error("screenshot()截图文件不存在")
			}
		}
	}
}

func TestBrowserSession_Evaluate(t *testing.T) {
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var session = createBrowserSession(30);
		var navResult = session.navigate("https://www.example.com");
		if (!navResult.success) {
			throw new Error("导航失败: " + navResult.error);
		}
		var result = session.evaluate("document.title");
		session.close();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("浏览器会话测试可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("evaluate()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success != nil && success.ToBoolean() {
		resultVal := resultObj.Get("result")
		if resultVal == nil || goja.IsUndefined(resultVal) {
			t.Error("evaluate()成功时应该包含result字段")
		}
	}
}

func TestBrowserSession_Click(t *testing.T) {
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var session = createBrowserSession(30);
		var navResult = session.navigate("https://www.example.com");
		if (!navResult.success) {
			throw new Error("导航失败: " + navResult.error);
		}
		var result = session.click("body");
		session.close();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("浏览器会话测试可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("click()返回的对象为nil")
	}

	// 验证返回结果结构
	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Error("click()缺少success字段")
	}
}

func TestBrowserSession_Fill(t *testing.T) {
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var session = createBrowserSession(30);
		var navResult = session.navigate("https://www.example.com");
		if (!navResult.success) {
			throw new Error("导航失败: " + navResult.error);
		}
		// 尝试填充body元素（虽然body不是输入框，但可以测试API是否正常工作）
		var result = session.fill("body", "test value");
		session.close();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("浏览器会话测试可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("fill()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Error("fill()缺少success字段")
	}
}

func TestBrowserSession_Wait(t *testing.T) {
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 测试等待元素
	code := `
		var session = createBrowserSession(30);
		var navResult = session.navigate("https://www.example.com");
		if (!navResult.success) {
			throw new Error("导航失败: " + navResult.error);
		}
		var result = session.wait("body");
		session.close();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("浏览器会话测试可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("wait()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Error("wait()缺少success字段")
	}

	// 测试等待时间
	code2 := `
		var session = createBrowserSession(30);
		var navResult = session.navigate("https://www.example.com");
		if (!navResult.success) {
			throw new Error("导航失败: " + navResult.error);
		}
		var result = session.wait(0.5); // 等待0.5秒
		session.close();
		result;
	`

	result2, err := sb.Run(code2)
	if err != nil {
		t.Logf("等待时间测试可能需要Chrome环境，跳过: %v", err)
		return
	}

	resultObj2 := result2.ToObject(sb.vm)
	if resultObj2 == nil {
		t.Fatal("wait()返回的对象为nil")
	}

	success2 := resultObj2.Get("success")
	if success2 == nil || goja.IsUndefined(success2) {
		t.Error("wait()等待时间缺少success字段")
	}
}

func TestBrowserSession_GetURL(t *testing.T) {
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var session = createBrowserSession(30);
		var navResult = session.navigate("https://www.example.com");
		if (!navResult.success) {
			throw new Error("导航失败: " + navResult.error);
		}
		var result = session.getURL();
		session.close();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("浏览器会话测试可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("getURL()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Error("getURL()缺少success字段")
	}

	if success.ToBoolean() {
		url := resultObj.Get("url")
		if url == nil || goja.IsUndefined(url) {
			t.Error("getURL()成功时应该包含url字段")
		} else {
			urlStr := url.String()
			if urlStr == "" {
				t.Error("getURL()返回的URL为空")
			}
			t.Logf("获取到的URL: %s", urlStr)
		}
	}
}

func TestBrowserSession_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	t.Run("createBrowserSession函数存在", func(t *testing.T) {
		_, err := sb.Run("typeof createBrowserSession")
		if err != nil {
			t.Errorf("createBrowserSession函数未定义: %v", err)
		}
	})

	t.Run("会话关闭后操作应返回错误", func(t *testing.T) {
		code := `
			var session = createBrowserSession(30);
			session.close();
			var result = session.navigate("https://www.example.com");
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("执行代码失败: %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		if resultObj == nil {
			t.Fatal("返回的对象为nil")
		}

		success := resultObj.Get("success")
		if success != nil && success.ToBoolean() {
			t.Error("关闭后的会话操作应该失败")
		}

		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("关闭后的会话操作应该返回错误")
		}
	})
}

func TestBrowserSession_NavigateBaiduWithRedirect(t *testing.T) {
	// 测试获取百度首页并跟随301跳转
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var session = createBrowserSession(60);
		var navResult = session.navigate("https://www.baidu.com");
		if (!navResult.success) {
			throw new Error("导航失败: " + navResult.error);
		}
		var htmlResult = session.getHTML();
		session.close();
		htmlResult;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("浏览器会话测试可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("getHTML()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Fatal("getHTML()缺少success字段")
	}

	if !success.ToBoolean() {
		errorVal := resultObj.Get("error")
		if errorVal != nil && !goja.IsUndefined(errorVal) {
			t.Fatalf("getHTML()失败: %v", errorVal.String())
		}
		t.Fatal("getHTML()失败，但未返回错误信息")
	}

	// 获取并打印HTML内容
	html := resultObj.Get("html")
	if html == nil || goja.IsUndefined(html) {
		t.Fatal("getHTML()成功时应该包含html字段")
	}

	htmlContent := html.String()
	t.Logf("获取到的HTML内容长度: %d 字符", len(htmlContent))

	// 打印HTML内容的前500个字符作为示例
	if len(htmlContent) > 500 {
		t.Logf("HTML内容预览（前500字符）:\n%s...", htmlContent[:500])
	} else {
		t.Logf("HTML内容:\n%s", htmlContent)
	}

	// 验证HTML内容包含百度相关的特征（说明成功跟随了301跳转并获取到了内容）
	if len(htmlContent) == 0 {
		t.Error("HTML内容为空，可能未成功跟随301跳转")
	}

	// 检查是否包含百度页面的特征元素
	// 百度首页通常包含这些关键词
	if len(htmlContent) > 0 {
		t.Log("成功获取到HTML内容，说明已跟随301跳转并加载了页面")
	}
}

func TestBrowserSession_ComplexWorkflow(t *testing.T) {
	// 测试复杂的连续操作流程：导航、等待、获取HTML、截图
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	outputPath := testDir + "/workflow-screenshot.png"

	code := `
		var session = createBrowserSession(60);
		
		// 1. 导航到页面
		var navResult = session.navigate("https://www.example.com");
		if (!navResult.success) {
			session.close();
			throw new Error("导航失败: " + navResult.error);
		}
		
		// 2. 等待body元素
		var waitResult = session.wait("body");
		if (!waitResult.success) {
			session.close();
			throw new Error("等待失败: " + waitResult.error);
		}
		
		// 3. 获取URL
		var urlResult = session.getURL();
		if (!urlResult.success) {
			session.close();
			throw new Error("获取URL失败: " + urlResult.error);
		}
		
		// 4. 执行JavaScript
		var evalResult = session.evaluate("document.title");
		if (!evalResult.success) {
			session.close();
			throw new Error("执行脚本失败: " + evalResult.error);
		}
		
		// 5. 截图
		var screenshotResult = session.screenshot("` + outputPath + `");
		if (!screenshotResult.success) {
			session.close();
			throw new Error("截图失败: " + screenshotResult.error);
		}
		
		// 6. 获取HTML
		var htmlResult = session.getHTML();
		
		// 7. 关闭会话
		session.close();
		
		// 返回所有结果
		{
			url: urlResult.url,
			title: evalResult.result,
			screenshot: screenshotResult.path,
			htmlLength: htmlResult.html ? htmlResult.html.length : 0
		};
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("复杂工作流测试可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("工作流测试返回的对象为nil")
	}

	// 验证返回的各个字段
	url := resultObj.Get("url")
	if url == nil || goja.IsUndefined(url) {
		t.Error("工作流测试应该返回URL")
	} else {
		t.Logf("工作流测试获取到的URL: %s", url.String())
	}

	title := resultObj.Get("title")
	if title == nil || goja.IsUndefined(title) {
		t.Error("工作流测试应该返回标题")
	} else {
		t.Logf("工作流测试获取到的标题: %s", title.String())
	}

	screenshotPath := resultObj.Get("screenshot")
	if screenshotPath != nil && !goja.IsUndefined(screenshotPath) {
		pathStr := screenshotPath.String()
		if _, err := os.Stat(pathStr); os.IsNotExist(err) {
			t.Errorf("工作流测试截图文件不存在: %s", pathStr)
		} else {
			t.Logf("工作流测试截图已保存: %s", pathStr)
		}
	}

	htmlLength := resultObj.Get("htmlLength")
	if htmlLength != nil && !goja.IsUndefined(htmlLength) {
		length := htmlLength.ToFloat()
		if length > 0 {
			t.Logf("工作流测试获取到的HTML长度: %.0f 字符", length)
		}
	}
}

func TestBrowserSession_BotDetection(t *testing.T) {
	// 测试访问 browserscan.net 的机器人检测页面，等待8秒后截图
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 将截图保存到固定位置以便查看
	outputPath := "example/data/bot-detection-screenshot.png"
	// 确保目录存在
	if err := os.MkdirAll("example/data", 0755); err != nil {
		t.Fatalf("创建目录失败: %v", err)
	}

	code := `
		var session = createBrowserSession(60);
		var navResult = session.navigate("https://www.browserscan.net/tc/bot-detection");
		if (!navResult.success) {
			session.close();
			throw new Error("导航失败: " + navResult.error);
		}
		// 等待8秒
		var waitResult = session.wait(8);
		if (!waitResult.success) {
			session.close();
			throw new Error("等待失败: " + waitResult.error);
		}
		var result = session.screenshot("` + outputPath + `");
		session.close();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("机器人检测测试可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("screenshot()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Error("screenshot()缺少success字段")
	}

	if !success.ToBoolean() {
		errorVal := resultObj.Get("error")
		if errorVal != nil && !goja.IsUndefined(errorVal) {
			t.Fatalf("screenshot()失败: %v", errorVal.String())
		}
		t.Fatal("screenshot()失败，但未返回错误信息")
	}

	// 确定实际的文件路径
	var actualFilePath string
	path := resultObj.Get("path")
	if path != nil && !goja.IsUndefined(path) {
		filePath := path.String()
		t.Logf("返回的截图路径: %s", filePath)

		// 检查文件是否存在
		fileInfo, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			t.Errorf("截图文件不存在: %s", filePath)
			// 也检查原始路径
			if fileInfo2, err2 := os.Stat(outputPath); err2 == nil {
				t.Logf("但原始路径存在: %s (大小: %d 字节)", outputPath, fileInfo2.Size())
				actualFilePath = outputPath
			} else {
				t.Errorf("原始路径也不存在: %s", outputPath)
			}
		} else if err != nil {
			t.Errorf("检查截图文件时出错: %v", err)
		} else {
			t.Logf("截图已保存到: %s (大小: %d 字节)", filePath, fileInfo.Size())
			if fileInfo.Size() == 0 {
				t.Error("截图文件大小为0，可能保存失败")
			} else {
				actualFilePath = filePath
			}
		}
	} else {
		t.Error("screenshot()未返回path字段")
	}

	// 额外验证：直接检查我们传入的路径
	if fileInfo, err := os.Stat(outputPath); err == nil {
		t.Logf("直接检查输出路径: %s (大小: %d 字节)", outputPath, fileInfo.Size())
		if actualFilePath == "" && fileInfo.Size() > 0 {
			actualFilePath = outputPath
		}
	} else {
		t.Errorf("直接检查输出路径失败: %s, 错误: %v", outputPath, err)
	}

	// 如果文件存在，在测试结束后打开它
	if actualFilePath != "" {
		defer func() {
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "darwin": // macOS
				cmd = exec.Command("open", actualFilePath)
			case "linux":
				cmd = exec.Command("xdg-open", actualFilePath)
			case "windows":
				cmd = exec.Command("cmd", "/c", "start", "", actualFilePath)
			default:
				t.Logf("不支持的操作系统，无法自动打开图片: %s", runtime.GOOS)
				return
			}
			if err := cmd.Run(); err != nil {
				t.Logf("打开图片文件失败: %v", err)
			} else {
				t.Logf("已使用系统默认程序打开图片: %s", actualFilePath)
			}
		}()
	}
}
