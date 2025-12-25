package jssandbox

import (
	"context"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogger_Registration(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetOutput(io.Discard) // 在测试中禁用日志输出
	sb := NewSandboxWithLogger(ctx, logger)
	defer sb.Close()

	// 测试 logger 对象是否已注册
	_, err := sb.Run("typeof logger")
	if err != nil {
		t.Fatalf("logger对象未注册: %v", err)
	}

	// 测试 logger 对象的方法
	_, err = sb.Run("typeof logger.debug")
	if err != nil {
		t.Fatalf("logger.debug方法未注册: %v", err)
	}

	_, err = sb.Run("typeof logger.info")
	if err != nil {
		t.Fatalf("logger.info方法未注册: %v", err)
	}

	_, err = sb.Run("typeof logger.warn")
	if err != nil {
		t.Fatalf("logger.warn方法未注册: %v", err)
	}

	_, err = sb.Run("typeof logger.error")
	if err != nil {
		t.Fatalf("logger.error方法未注册: %v", err)
	}
}

func TestLogger_BasicLogging(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	sb := NewSandboxWithLogger(ctx, logger)
	defer sb.Close()

	// 测试基本日志功能（不应该出错）
	_, err := sb.Run("logger.debug('test debug message')")
	if err != nil {
		t.Errorf("logger.debug() 失败: %v", err)
	}

	_, err = sb.Run("logger.info('test info message')")
	if err != nil {
		t.Errorf("logger.info() 失败: %v", err)
	}

	_, err = sb.Run("logger.warn('test warn message')")
	if err != nil {
		t.Errorf("logger.warn() 失败: %v", err)
	}

	_, err = sb.Run("logger.error('test error message')")
	if err != nil {
		t.Errorf("logger.error() 失败: %v", err)
	}
}

func TestLogger_WithFields(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	sb := NewSandboxWithLogger(ctx, logger)
	defer sb.Close()

	// 测试带字段的日志
	_, err := sb.Run(`
		var fieldLogger = logger.withFields({userId: 123, action: 'login'});
		fieldLogger.info('User logged in');
	`)
	if err != nil {
		t.Errorf("logger.withFields() 失败: %v", err)
	}
}

func TestLogger_SetLevel(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	sb := NewSandboxWithLogger(ctx, logger)
	defer sb.Close()

	// 测试设置日志级别
	_, err := sb.Run("logger.setLevel('debug')")
	if err != nil {
		t.Errorf("logger.setLevel() 失败: %v", err)
	}

	// 测试获取日志级别
	result, err := sb.Run("logger.getLevel()")
	if err != nil {
		t.Errorf("logger.getLevel() 失败: %v", err)
	}

	if result == nil {
		t.Error("logger.getLevel() 返回 nil")
	}
}

func TestLogger_IsLevelEnabled(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	sb := NewSandboxWithLogger(ctx, logger)
	defer sb.Close()

	// 测试检查日志级别是否启用
	result, err := sb.Run("logger.isLevelEnabled('info')")
	if err != nil {
		t.Errorf("logger.isLevelEnabled() 失败: %v", err)
	}

	if result == nil {
		t.Error("logger.isLevelEnabled() 返回 nil")
	}
}

