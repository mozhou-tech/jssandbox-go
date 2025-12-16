package jssandbox

import (
	"context"
	"fmt"
	"time"

	"github.com/dop251/goja"
	"go.uber.org/zap"
)

// Sandbox 表示一个JavaScript沙盒环境
type Sandbox struct {
	vm     *goja.Runtime
	logger *zap.Logger
	ctx    context.Context
}

// NewSandbox 创建一个新的沙盒实例
func NewSandbox(ctx context.Context) *Sandbox {
	vm := goja.New()
	logger, _ := zap.NewProduction()

	sb := &Sandbox{
		vm:     vm,
		logger: logger,
		ctx:    ctx,
	}

	// 注册所有扩展功能
	sb.registerExtensions()

	return sb
}

// NewSandboxWithLogger 使用自定义logger创建沙盒
func NewSandboxWithLogger(ctx context.Context, logger *zap.Logger) *Sandbox {
	vm := goja.New()
	sb := &Sandbox{
		vm:     vm,
		logger: logger,
		ctx:    ctx,
	}
	sb.registerExtensions()
	return sb
}

// registerExtensions 注册所有扩展功能到JavaScript运行时
func (sb *Sandbox) registerExtensions() {
	// 注册系统操作
	sb.registerSystemOps()
	// 注册HTTP请求
	sb.registerHTTP()
	// 注册文件系统操作
	sb.registerFileSystem()
	// 注册浏览器操作
	sb.registerBrowser()
	// 注册文档读取功能
	sb.registerDocuments()
	// 注册图片处理功能
	sb.registerImageProcessing()
	// 注册文件类型检测功能
	sb.registerFileTypeDetection()
}

// Run 执行JavaScript代码
func (sb *Sandbox) Run(code string) (goja.Value, error) {
	return sb.vm.RunString(code)
}

// RunWithTimeout 在指定超时时间内执行JavaScript代码
func (sb *Sandbox) RunWithTimeout(code string, timeout time.Duration) (goja.Value, error) {
	ctx, cancel := context.WithTimeout(sb.ctx, timeout)
	defer cancel()

	done := make(chan error, 1)
	var result goja.Value
	var err error

	go func() {
		result, err = sb.vm.RunString(code)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("执行超时: %v", timeout)
	case err := <-done:
		return result, err
	}
}

// Set 在JavaScript运行时中设置变量
func (sb *Sandbox) Set(name string, value interface{}) {
	sb.vm.Set(name, value)
}

// Get 从JavaScript运行时中获取变量
func (sb *Sandbox) Get(name string) goja.Value {
	return sb.vm.Get(name)
}

// Close 关闭沙盒并清理资源
func (sb *Sandbox) Close() error {
	if sb.logger != nil {
		sb.logger.Sync()
	}
	return nil
}
