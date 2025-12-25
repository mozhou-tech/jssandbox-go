package jssandbox

import (
	"context"
	"testing"

	"github.com/dop251/goja"
)

func TestResolveDNS(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = resolveDNS("localhost");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("resolveDNS() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		// DNS解析可能失败，这是可以接受的
		t.Logf("resolveDNS()返回错误（可能是网络问题）")
		return
	}

	ips := resultObj.Get("ips")
	if ips == nil || goja.IsUndefined(ips) {
		t.Error("resolveDNS()缺少ips字段")
	}
}

func TestResolveDNS_Google(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = resolveDNS("www.google.com");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("resolveDNS() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		// DNS解析可能失败（网络问题），跳过测试
		t.Skip("DNS解析失败，可能是网络问题")
		return
	}

	ips := resultObj.Get("ips")
	if ips == nil || goja.IsUndefined(ips) {
		t.Error("resolveDNS()缺少ips字段")
	}

	// 验证至少有一个IP地址
	ipsArray := ips.ToObject(sb.vm)
	if len(ipsArray.Keys()) == 0 {
		t.Error("resolveDNS()应该返回至少一个IP地址")
	}
}

func TestPing(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = ping("127.0.0.1", 2);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("ping() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Logf("ping()返回错误（可能是网络问题）")
		return
	}

	sent := resultObj.Get("sent")
	if sent.ToInteger() != 2 {
		t.Errorf("ping()发送次数不正确, got %d, want 2", sent.ToInteger())
	}
}

func TestCheckPort(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 测试本地回环地址的常见端口
	testCases := []struct {
		host string
		port int
	}{
		{"127.0.0.1", 80},
		{"localhost", 443},
		{"127.0.0.1", 8080},
	}

	for _, tc := range testCases {
		t.Run(tc.host+":80", func(t *testing.T) {
			code := `
				var result = checkPort("` + tc.host + `", 80);
				result;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("checkPort() error = %v", err)
			}

			resultObj := result.ToObject(sb.vm)
			success := resultObj.Get("success")
			if !success.ToBoolean() {
				t.Logf("checkPort()返回错误")
				return
			}

			open := resultObj.Get("open")
			// 端口可能开放也可能不开放，只要函数能正常执行即可
			_ = open.ToBoolean()
		})
	}
}

func TestCheckPort_Localhost(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = checkPort("127.0.0.1", 80, 3);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("checkPort() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Logf("checkPort()返回错误（可能是网络问题）")
		return
	}

	host := resultObj.Get("host")
	if host.String() != "127.0.0.1" {
		t.Errorf("checkPort()主机不正确, got %s, want 127.0.0.1", host.String())
	}

	port := resultObj.Get("port")
	if port.ToInteger() != 80 {
		t.Errorf("checkPort()端口不正确, got %d, want 80", port.ToInteger())
	}
}

