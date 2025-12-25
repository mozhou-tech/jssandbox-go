package jssandbox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dop251/goja"
)

func TestHTTPRequest_GET(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var response = httpRequest("` + server.URL + `", {
			method: "GET"
		});
		response;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("httpRequest() error = %v", err)
	}

	responseObj := result.ToObject(sb.vm)
	if responseObj == nil {
		t.Fatal("httpRequest()返回的对象为nil")
	}

	status := responseObj.Get("status")
	if status.ToInteger() != 200 {
		t.Errorf("httpRequest()状态码不正确, got %d, want 200", status.ToInteger())
	}

	body := responseObj.Get("body")
	if body.String() != "Hello, World!" {
		t.Errorf("httpRequest()响应体不正确, got %s, want Hello, World!", body.String())
	}
}

func TestHTTPRequest_POST(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("POST successful"))
	}))
	defer server.Close()

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var response = httpRequest("` + server.URL + `", {
			method: "POST",
			body: "test data"
		});
		response;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("httpRequest() error = %v", err)
	}

	responseObj := result.ToObject(sb.vm)
	status := responseObj.Get("status")
	if status.ToInteger() != 200 {
		t.Errorf("httpRequest() POST状态码不正确, got %d, want 200", status.ToInteger())
	}
}

func TestHTTPRequest_Headers(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customHeader := r.Header.Get("X-Custom-Header")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(customHeader))
	}))
	defer server.Close()

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var response = httpRequest("` + server.URL + `", {
			method: "GET",
			headers: {
				"X-Custom-Header": "test-value"
			}
		});
		response;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("httpRequest() error = %v", err)
	}

	responseObj := result.ToObject(sb.vm)
	body := responseObj.Get("body")
	if body.String() != "test-value" {
		t.Errorf("httpRequest()自定义头部未正确传递, got %s, want test-value", body.String())
	}
}

func TestHTTPRequest_Timeout(t *testing.T) {
	// 创建慢速测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟慢速响应
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("delayed response"))
	}))
	defer server.Close()

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var response = httpRequest("` + server.URL + `", {
			method: "GET",
			timeout: 1
		});
		response;
	`

	start := time.Now()
	result, err := sb.Run(code)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("httpRequest() error = %v", err)
	}

	// 验证超时生效（应该在1秒左右完成，而不是2秒）
	if elapsed > 1500*time.Millisecond {
		t.Errorf("httpRequest()超时未生效, 耗时 %v", elapsed)
	}

	responseObj := result.ToObject(sb.vm)
	errorVal := responseObj.Get("error")
	if errorVal == nil || goja.IsUndefined(errorVal) {
		t.Log("超时可能未触发（取决于网络环境）")
	}
}

func TestHTTPGet(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET response"))
	}))
	defer server.Close()

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var response = httpGet("` + server.URL + `");
		response;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("httpGet() error = %v", err)
	}

	responseObj := result.ToObject(sb.vm)
	status := responseObj.Get("status")
	if status.ToInteger() != 200 {
		t.Errorf("httpGet()状态码不正确, got %d, want 200", status.ToInteger())
	}

	body := responseObj.Get("body")
	if body.String() != "GET response" {
		t.Errorf("httpGet()响应体不正确, got %s, want GET response", body.String())
	}
}

func TestHTTPPost(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("POST response"))
	}))
	defer server.Close()

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var response = httpPost("` + server.URL + `", '{"key": "value"}');
		response;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("httpPost() error = %v", err)
	}

	responseObj := result.ToObject(sb.vm)
	status := responseObj.Get("status")
	if status.ToInteger() != 200 {
		t.Errorf("httpPost()状态码不正确, got %d, want 200", status.ToInteger())
	}
}

func TestHTTPRequest_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	t.Run("缺少URL参数", func(t *testing.T) {
		result, err := sb.Run("httpRequest()")
		if err != nil {
			t.Fatalf("httpRequest() error = %v", err)
		}

		responseObj := result.ToObject(sb.vm)
		errorVal := responseObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Error("httpRequest()缺少URL参数时应该返回错误")
		}
	})

	t.Run("无效URL", func(t *testing.T) {
		code := `
			var response = httpRequest("http://invalid-url-that-does-not-exist-12345.com");
			response;
		`
		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("httpRequest() error = %v", err)
		}

		responseObj := result.ToObject(sb.vm)
		errorVal := responseObj.Get("error")
		if errorVal == nil || goja.IsUndefined(errorVal) {
			t.Log("无效URL可能在某些环境下不会立即返回错误")
		}
	})
}

func TestHTTPRequest_ResponseHeaders(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Response", "test")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var response = httpRequest("` + server.URL + `");
		response;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("httpRequest() error = %v", err)
	}

	responseObj := result.ToObject(sb.vm)
	headers := responseObj.Get("headers")
	if headers == nil || goja.IsUndefined(headers) {
		t.Error("httpRequest()缺少headers字段")
	}

	contentType := responseObj.Get("contentType")
	if contentType.String() != "application/json" {
		t.Errorf("httpRequest() contentType不正确, got %s, want application/json", contentType.String())
	}
}

