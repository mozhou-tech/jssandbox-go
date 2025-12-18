package jssandbox

import (
	"context"
	"testing"

	"github.com/dop251/goja"
)

func TestEncodeBase64(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = encodeBase64("Hello, World!");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("encodeBase64() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("encodeBase64()应该返回success: true")
	}

	data := resultObj.Get("data")
	if data.String() == "" {
		t.Error("encodeBase64()返回的编码数据不应该为空")
	}
}

func TestDecodeBase64(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	original := "Hello, World!"

	// 先编码
	encodeCode := `
		var encoded = encodeBase64("` + original + `");
		encoded.data;
	`

	encodedResult, err := sb.Run(encodeCode)
	if err != nil {
		t.Fatalf("编码失败: %v", err)
	}

	encodedData := encodedResult.String()

	// 再解码
	decodeCode := `
		var result = decodeBase64("` + encodedData + `");
		result;
	`

	result, err := sb.Run(decodeCode)
	if err != nil {
		t.Fatalf("decodeBase64() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("decodeBase64()应该返回success: true")
	}

	data := resultObj.Get("data")
	if data.String() != original {
		t.Errorf("decodeBase64()解码结果不正确, got %s, want %s", data.String(), original)
	}
}

func TestEncodeDecodeBase64_Integration(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []string{
		"Hello, World!",
		"测试中文",
		"123456",
		"Special chars: !@#$%^&*()",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			code := `
				var encoded = encodeBase64("` + tc + `");
				var decoded = decodeBase64(encoded.data);
				decoded.data;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("编码解码失败: %v", err)
			}

			if result.String() != tc {
				t.Errorf("编码解码结果不匹配, got %s, want %s", result.String(), tc)
			}
		})
	}
}

func TestEncodeURL(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = encodeURL("hello world");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("encodeURL() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("encodeURL()应该返回success: true")
	}

	data := resultObj.Get("data")
	if data.String() == "" {
		t.Error("encodeURL()返回的编码数据不应该为空")
	}
}

func TestDecodeURL(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	original := "hello world"

	// 先编码
	encodeCode := `
		var encoded = encodeURL("` + original + `");
		encoded.data;
	`

	encodedResult, err := sb.Run(encodeCode)
	if err != nil {
		t.Fatalf("编码失败: %v", err)
	}

	encodedData := encodedResult.String()

	// 再解码
	decodeCode := `
		var result = decodeURL("` + encodedData + `");
		result;
	`

	result, err := sb.Run(decodeCode)
	if err != nil {
		t.Fatalf("decodeURL() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	data := resultObj.Get("data")
	if data.String() != original {
		t.Errorf("decodeURL()解码结果不正确, got %s, want %s", data.String(), original)
	}
}

func TestEncodeHTML(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = encodeHTML("<div>Hello</div>");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("encodeHTML() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("encodeHTML()应该返回success: true")
	}

	data := resultObj.Get("data")
	// 应该包含HTML实体
	if data.String() == "<div>Hello</div>" {
		t.Error("encodeHTML()应该对HTML进行编码")
	}
}

func TestDecodeHTML(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	original := "<div>Hello</div>"

	// 先编码
	encodeCode := `
		var encoded = encodeHTML("` + original + `");
		encoded.data;
	`

	encodedResult, err := sb.Run(encodeCode)
	if err != nil {
		t.Fatalf("编码失败: %v", err)
	}

	encodedData := encodedResult.String()

	// 再解码
	decodeCode := `
		var result = decodeHTML("` + encodedData + `");
		result;
	`

	result, err := sb.Run(decodeCode)
	if err != nil {
		t.Fatalf("decodeHTML() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	data := resultObj.Get("data")
	if data.String() != original {
		t.Errorf("decodeHTML()解码结果不正确, got %s, want %s", data.String(), original)
	}
}

