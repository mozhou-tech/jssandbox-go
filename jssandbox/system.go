package jssandbox

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// registerSystemOps 注册系统操作函数到JavaScript运行时
func (sb *Sandbox) registerSystemOps() {
	sb.vm.Set("getCurrentTime", func() string {
		return time.Now().Format("15:04:05")
	})

	sb.vm.Set("getCurrentDate", func() string {
		return time.Now().Format("2006-01-02")
	})

	sb.vm.Set("getCurrentDateTime", func() string {
		return time.Now().Format("2006-01-02 15:04:05")
	})

	sb.vm.Set("getCPUNum", func() int {
		return runtime.NumCPU()
	})

	sb.vm.Set("getMemorySize", func(call goja.FunctionCall) goja.Value {
		vm, _ := mem.VirtualMemory()
		if vm == nil {
			return sb.vm.ToValue(map[string]interface{}{
				"total":        0,
				"available":    0,
				"used":         0,
				"totalStr":     "0 B",
				"availableStr": "0 B",
				"usedStr":      "0 B",
			})
		}
		return sb.vm.ToValue(map[string]interface{}{
			"total":        vm.Total,
			"available":    vm.Available,
			"used":         vm.Used,
			"totalStr":     humanize.Bytes(vm.Total),
			"availableStr": humanize.Bytes(vm.Available),
			"usedStr":      humanize.Bytes(vm.Used),
		})
	})

	sb.vm.Set("getDiskSize", func(call goja.FunctionCall) goja.Value {
		path := "/"
		if len(call.Arguments) > 0 {
			path = call.Arguments[0].String()
		}

		usage, err := disk.Usage(path)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"total":    0,
				"free":     0,
				"used":     0,
				"totalStr": "0 B",
				"freeStr":  "0 B",
				"usedStr":  "0 B",
				"error":    err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"total":    usage.Total,
			"free":     usage.Free,
			"used":     usage.Used,
			"totalStr": humanize.Bytes(usage.Total),
			"freeStr":  humanize.Bytes(usage.Free),
			"usedStr":  humanize.Bytes(usage.Used),
		})
	})

	sb.vm.Set("sleep", func(ms int) {
		time.Sleep(time.Duration(ms) * time.Millisecond)
	})

	// 注册 console 对象
	consoleObj := sb.vm.NewObject()
	
	// 辅助函数：格式化 console 参数
	formatConsoleArgs := func(call goja.FunctionCall) string {
		if len(call.Arguments) == 0 {
			return ""
		}
		parts := make([]string, len(call.Arguments))
		for i, arg := range call.Arguments {
			if goja.IsUndefined(arg) {
				parts[i] = "undefined"
			} else if goja.IsNull(arg) {
				parts[i] = "null"
			} else {
				parts[i] = fmt.Sprintf("%v", arg.Export())
			}
		}
		return strings.Join(parts, " ")
	}
	
	consoleObj.Set("log", func(call goja.FunctionCall) goja.Value {
		sb.logger.Info(formatConsoleArgs(call))
		return goja.Undefined()
	})
	consoleObj.Set("error", func(call goja.FunctionCall) goja.Value {
		sb.logger.Error(formatConsoleArgs(call))
		return goja.Undefined()
	})
	consoleObj.Set("warn", func(call goja.FunctionCall) goja.Value {
		sb.logger.Warn(formatConsoleArgs(call))
		return goja.Undefined()
	})
	consoleObj.Set("info", func(call goja.FunctionCall) goja.Value {
		sb.logger.Info(formatConsoleArgs(call))
		return goja.Undefined()
	})
	consoleObj.Set("debug", func(call goja.FunctionCall) goja.Value {
		sb.logger.Debug(formatConsoleArgs(call))
		return goja.Undefined()
	})
	sb.vm.Set("console", consoleObj)
}
