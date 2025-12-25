package jssandbox

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
)

// Sandbox 表示一个JavaScript沙盒环境
type Sandbox struct {
	vm     *goja.Runtime
	logger *logrus.Logger
	ctx    context.Context
	config *Config
	// 浏览器相关的共享资源
	browserAllocator context.Context
	browserCancel    context.CancelFunc
	browserMu        sync.Mutex
	browserInit      bool
}

// NewSandbox 创建一个新的沙盒实例（使用默认配置）
func NewSandbox(ctx context.Context) *Sandbox {
	return NewSandboxWithConfig(ctx, DefaultConfig())
}

// NewSandboxWithConfig 使用指定配置创建新的沙盒实例
func NewSandboxWithConfig(ctx context.Context, config *Config) *Sandbox {
	vm := goja.New()
	logger := GetLogger()

	sb := &Sandbox{
		vm:     vm,
		logger: logger,
		ctx:    ctx,
		config: config,
	}

	// 注册所有扩展功能
	sb.registerExtensions()

	return sb
}

// NewSandboxWithLogger 使用自定义logger创建沙盒（使用默认配置）
func NewSandboxWithLogger(ctx context.Context, logger *logrus.Logger) *Sandbox {
	return NewSandboxWithLoggerAndConfig(ctx, logger, DefaultConfig())
}

// NewSandboxWithLoggerAndConfig 使用自定义logger和配置创建沙盒
func NewSandboxWithLoggerAndConfig(ctx context.Context, logger *logrus.Logger, config *Config) *Sandbox {
	vm := goja.New()
	sb := &Sandbox{
		vm:     vm,
		logger: logger,
		ctx:    ctx,
		config: config,
	}
	sb.registerExtensions()
	return sb
}

// registerExtensions 注册所有扩展功能到JavaScript运行时
// 根据配置选择性注册功能模块
func (sb *Sandbox) registerExtensions() {
	// 注册系统操作（始终启用）
	sb.registerSystemOps()

	// 注册基础工具功能（始终启用）
	sb.registerLogger()     // 日志功能
	sb.registerCrypto()     // 加密/解密
	sb.registerCompress()   // 压缩/解压缩
	sb.registerCSV()        // CSV处理
	sb.registerEnv()        // 环境变量和配置
	sb.registerValidation() // 数据验证
	sb.registerDateTime()   // 日期时间增强
	sb.registerEncoding()   // 编码/解码增强
	sb.registerProcess()    // 进程管理
	sb.registerNetwork()    // 网络工具
	sb.registerPath()       // 路径处理增强
	sb.registerText()       // 文本操作

	// 根据配置选择性注册功能
	if sb.config.EnableHTTP {
		sb.registerHTTP()
	}
	if sb.config.EnableFileSystem {
		sb.registerFileSystem()
	}
	if sb.config.EnableBrowser {
		sb.registerBrowser()
	}
	if sb.config.EnableDocuments {
		sb.registerDocuments()
	}
	if sb.config.EnableImageProcessing {
		sb.registerImageProcessing()
	}
	if sb.config.EnableGoQuery {
		sb.registerGoQuery()
	}
	// 文件类型检测始终启用（文件系统功能依赖它）
	sb.registerFileTypeDetection()
}

// Run 执行JavaScript代码
func (sb *Sandbox) Run(code string) (goja.Value, error) {
	return sb.vm.RunString(code)
}

// RunWithTimeout 在指定超时时间内执行JavaScript代码
// 如果 timeout 为 0，则使用配置中的默认超时时间
func (sb *Sandbox) RunWithTimeout(code string, timeout time.Duration) (goja.Value, error) {
	if timeout == 0 {
		timeout = sb.config.DefaultTimeout
	}

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
		return nil, NewSandboxError(ErrCodeTimeout, fmt.Sprintf("执行超时: %v", timeout))
	case err := <-done:
		if err != nil {
			return nil, NewSandboxErrorWithCause(ErrCodeUnknown, "执行JavaScript代码失败", err)
		}
		return result, nil
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

// Delete 从JavaScript运行时中删除变量，返回是否删除成功
func (sb *Sandbox) Delete(name string) error {
	return sb.vm.GlobalObject().Delete(name)
}

// Close 关闭沙盒并清理资源
func (sb *Sandbox) Close() error {
	// 关闭浏览器 allocator（如果已初始化）
	sb.browserMu.Lock()
	if sb.browserInit && sb.browserCancel != nil {
		sb.browserCancel()
		sb.browserInit = false
	}
	sb.browserMu.Unlock()
	// logrus 不需要显式同步
	return nil
}
