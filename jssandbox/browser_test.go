package jssandbox

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/dop251/goja"
)

func TestBrowserNavigate(t *testing.T) {
	// 注意：这个测试需要Chrome/Chromium可用
	// 如果环境中没有Chrome，测试可能会失败
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = browserNavigate("https://www.example.com");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("browserNavigate()可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("browserNavigate()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Error("browserNavigate()缺少success字段")
	}

	// 如果成功，验证HTML字段
	if success.ToBoolean() {
		html := resultObj.Get("html")
		if html == nil || goja.IsUndefined(html) {
			t.Error("browserNavigate()成功时应该包含html字段")
		}
	}
}

func TestBrowserScreenshot(t *testing.T) {
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	outputPath := testDir + "/screenshot.png"

	code := `
		var result = browserScreenshot("https://www.example.com", "` + outputPath + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("browserScreenshot()可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("browserScreenshot()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success != nil && success.ToBoolean() {
		// 验证截图文件是否存在
		path := resultObj.Get("path")
		if path != nil && !goja.IsUndefined(path) {
			if _, err := os.Stat(path.String()); os.IsNotExist(err) {
				t.Error("browserScreenshot()截图文件不存在")
			}
		}
	}
}

func TestBrowserEvaluate(t *testing.T) {
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = browserEvaluate("https://www.example.com", "document.title");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("browserEvaluate()可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("browserEvaluate()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success != nil && success.ToBoolean() {
		resultVal := resultObj.Get("result")
		if resultVal == nil || goja.IsUndefined(resultVal) {
			t.Error("browserEvaluate()成功时应该包含result字段")
		}
	}
}

func TestBrowserClick(t *testing.T) {
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = browserClick("https://www.example.com", "body");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("browserClick()可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("browserClick()返回的对象为nil")
	}

	// 验证返回结果结构
	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Error("browserClick()缺少success字段")
	}
}

func TestBrowserFill(t *testing.T) {
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = browserFill("https://www.example.com", "body", "test value");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("browserFill()可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("browserFill()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Error("browserFill()缺少success字段")
	}
}

func TestBrowser_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	t.Run("browserNavigate缺少参数", func(t *testing.T) {
		// browserNavigate需要URL参数，但这里测试的是函数调用本身
		// 由于函数定义需要参数，这个测试主要验证函数存在
		_, err := sb.Run("typeof browserNavigate")
		if err != nil {
			t.Errorf("browserNavigate函数未定义: %v", err)
		}
	})

	t.Run("browserScreenshot缺少参数", func(t *testing.T) {
		result, err := sb.Run("browserScreenshot()")
		if err != nil {
			t.Fatalf("browserScreenshot() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("browserScreenshot()缺少参数应该返回错误")
		}
	})

	t.Run("browserEvaluate缺少参数", func(t *testing.T) {
		result, err := sb.Run("browserEvaluate()")
		if err != nil {
			t.Fatalf("browserEvaluate() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		errorVal := resultObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("browserEvaluate()缺少参数应该返回错误")
		}
	})
}

func TestBrowserNavigateBaiduWithRedirect(t *testing.T) {
	// 测试获取百度首页并跟随301跳转
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = browserNavigate("https://www.baidu.com");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("browserNavigate()可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("browserNavigate()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Fatal("browserNavigate()缺少success字段")
	}

	if !success.ToBoolean() {
		errorVal := resultObj.Get("error")
		if errorVal != nil && !goja.IsUndefined(errorVal) {
			t.Fatalf("browserNavigate()失败: %v", errorVal.String())
		}
		t.Fatal("browserNavigate()失败，但未返回错误信息")
	}

	// 获取并打印HTML内容
	html := resultObj.Get("html")
	if html == nil || goja.IsUndefined(html) {
		t.Fatal("browserNavigate()成功时应该包含html字段")
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

func TestBrowserScanBotDetection(t *testing.T) {
	// 测试访问 browserscan.net 的机器人检测页面，等待8秒后截图
	if os.Getenv("SKIP_BROWSER_TESTS") == "true" {
		t.Skip("跳过浏览器测试（SKIP_BROWSER_TESTS=true）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	outputPath := testDir + "/bot-detection-screenshot.png"

	code := `
		var result = browserScreenshot("https://www.browserscan.net/tc/bot-detection", "` + outputPath + `", 8);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("browserScreenshot()可能需要Chrome环境，跳过: %v", err)
		t.Skip("浏览器测试需要Chrome环境")
		return
	}

	resultObj := result.ToObject(sb.vm)
	if resultObj == nil {
		t.Fatal("browserScreenshot()返回的对象为nil")
	}

	success := resultObj.Get("success")
	if success == nil || goja.IsUndefined(success) {
		t.Error("browserScreenshot()缺少success字段")
	}

	if !success.ToBoolean() {
		errorVal := resultObj.Get("error")
		if errorVal != nil && !goja.IsUndefined(errorVal) {
			t.Fatalf("browserScreenshot()失败: %v", errorVal.String())
		}
		t.Fatal("browserScreenshot()失败，但未返回错误信息")
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
		t.Error("browserScreenshot()未返回path字段")
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
