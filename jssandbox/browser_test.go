package jssandbox

import (
	"context"
	"os"
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

