package jssandbox

import (
	"time"

	"github.com/dop251/goja"
)

// SandboxRunner 定义沙盒执行接口
// 这个接口使得沙盒实现可以被替换，便于测试和扩展
type SandboxRunner interface {
	// Run 执行 JavaScript 代码
	Run(code string) (goja.Value, error)
	// RunWithTimeout 在指定超时时间内执行 JavaScript 代码
	RunWithTimeout(code string, timeout time.Duration) (goja.Value, error)
	// Set 在 JavaScript 运行时中设置变量
	Set(name string, value interface{})
	// Get 从 JavaScript 运行时中获取变量
	Get(name string) goja.Value
	// Close 关闭沙盒并清理资源
	Close() error
}

// 确保 Sandbox 实现了 SandboxRunner 接口
var _ SandboxRunner = (*Sandbox)(nil)

