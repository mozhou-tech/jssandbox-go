package jssandbox

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestNewSandbox(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	if sb == nil {
		t.Fatal("NewSandbox返回nil")
	}
	if sb.vm == nil {
		t.Fatal("vm为nil")
	}
	if sb.logger == nil {
		t.Fatal("logger为nil")
	}
	if sb.ctx == nil {
		t.Fatal("ctx为nil")
	}
}

func TestNewSandboxWithLogger(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetOutput(io.Discard) // 在测试中禁用日志输出
	sb := NewSandboxWithLogger(ctx, logger)
	defer sb.Close()

	if sb == nil {
		t.Fatal("NewSandboxWithLogger返回nil")
	}
	if sb.logger != logger {
		t.Fatal("logger未正确设置")
	}
}

func TestSandbox_Run(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "简单计算",
			code:    "1 + 1",
			wantErr: false,
		},
		{
			name:    "变量赋值",
			code:    "var x = 10; x * 2",
			wantErr: false,
		},
		{
			name:    "无效代码",
			code:    "var x = ;",
			wantErr: true,
		},
		{
			name:    "返回对象",
			code:    "({a: 1, b: 2})",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sb.Run(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Error("Run()返回nil结果")
			}
		})
	}
}

func TestSandbox_RunWithTimeout(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	t.Run("正常执行", func(t *testing.T) {
		result, err := sb.RunWithTimeout("1 + 1", 5*time.Second)
		if err != nil {
			t.Errorf("RunWithTimeout() error = %v", err)
		}
		if result == nil {
			t.Error("RunWithTimeout()返回nil结果")
		}
	})

	t.Run("超时测试", func(t *testing.T) {
		// 创建一个会长时间运行的代码
		code := `
			var start = Date.now();
			while (Date.now() - start < 2000) {
				// 等待2秒
			}
			"done";
		`
		result, err := sb.RunWithTimeout(code, 100*time.Millisecond)
		if err == nil {
			t.Error("RunWithTimeout()应该超时但没有返回错误")
		}
		if result != nil {
			t.Error("RunWithTimeout()超时后应该返回nil")
		}
	})
}

func TestSandbox_Set(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	sb.Set("testVar", 42)
	result, err := sb.Run("testVar")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	value := result.ToInteger()
	if value != 42 {
		t.Errorf("Set()设置的值不正确, got %d, want 42", value)
	}
}

func TestSandbox_Get(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 先设置一个变量
	sb.Run("var x = 100;")

	// 获取变量
	value := sb.Get("x")
	if value == nil {
		t.Fatal("Get()返回nil")
	}

	if value.ToInteger() != 100 {
		t.Errorf("Get()获取的值不正确, got %d, want 100", value.ToInteger())
	}
}

func TestSandbox_Close(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)

	err := sb.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestSandbox_RegisterExtensions(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 测试系统函数是否注册
	_, err := sb.Run("getCurrentTime()")
	if err != nil {
		t.Errorf("系统函数未注册: %v", err)
	}

	// 测试HTTP函数是否注册
	_, err = sb.Run("typeof httpGet")
	if err != nil {
		t.Errorf("HTTP函数未注册: %v", err)
	}

	// 测试文件系统函数是否注册
	_, err = sb.Run("typeof writeFile")
	if err != nil {
		t.Errorf("文件系统函数未注册: %v", err)
	}
}

func TestSandbox_ComplexScript(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		function fibonacci(n) {
			if (n <= 1) return n;
			return fibonacci(n - 1) + fibonacci(n - 2);
		}
		fibonacci(10);
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if result.ToInteger() != 55 {
		t.Errorf("复杂脚本执行结果不正确, got %d, want 55", result.ToInteger())
	}
}

func TestSandbox_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 测试语法错误
	_, err := sb.Run("var x =")
	if err == nil {
		t.Error("应该返回语法错误")
	}

	// 测试运行时错误
	_, err = sb.Run("undefinedVariable.foo")
	if err == nil {
		t.Error("应该返回运行时错误")
	}
}

func TestSandbox_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 取消上下文
	cancel()

	// 尝试执行代码（应该仍然可以执行，因为goja不直接使用context）
	// 但RunWithTimeout应该能检测到取消
	_, err := sb.RunWithTimeout("1 + 1", 1*time.Second)
	if err != nil {
		// 如果context被取消，可能会返回错误
		t.Logf("RunWithTimeout在context取消后返回错误: %v", err)
	}
}
