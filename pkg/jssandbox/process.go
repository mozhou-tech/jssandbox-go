package jssandbox

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/shirou/gopsutil/v3/process"
)

// registerProcess 注册进程管理功能到JavaScript运行时
func (sb *Sandbox) registerProcess() {
	// 执行系统命令
	sb.vm.Set("execCommand", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供命令参数",
			})
		}

		var cmd *exec.Cmd
		var args []string

		// 解析命令和参数
		if cmdVal := call.Arguments[0]; cmdVal != nil && !goja.IsUndefined(cmdVal) {
			if cmdObj := cmdVal.ToObject(sb.vm); cmdObj != nil {
				// 如果是数组
				if cmdArray, ok := cmdObj.Export().([]interface{}); ok {
					if len(cmdArray) > 0 {
						command := fmt.Sprintf("%v", cmdArray[0])
						for i := 1; i < len(cmdArray); i++ {
							args = append(args, fmt.Sprintf("%v", cmdArray[i]))
						}
						cmd = exec.Command(command, args...)
					}
				} else {
					// 如果是字符串，尝试解析
					cmdStr := cmdVal.String()
					if runtime.GOOS == "windows" {
						cmd = exec.Command("cmd", "/c", cmdStr)
					} else {
						parts := strings.Fields(cmdStr)
						if len(parts) > 0 {
							cmd = exec.Command(parts[0], parts[1:]...)
						}
					}
				}
			} else {
				cmdStr := cmdVal.String()
				if runtime.GOOS == "windows" {
					cmd = exec.Command("cmd", "/c", cmdStr)
				} else {
					parts := strings.Fields(cmdStr)
					if len(parts) > 0 {
						cmd = exec.Command(parts[0], parts[1:]...)
					}
				}
			}
		}

		if cmd == nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "无效的命令",
			})
		}

		// 设置超时（默认30秒）
		timeout := 30 * time.Second
		if len(call.Arguments) > 1 {
			if options := call.Arguments[1].ToObject(sb.vm); options != nil {
				if timeoutVal := options.Get("timeout"); timeoutVal != nil && !goja.IsUndefined(timeoutVal) {
					timeout = time.Duration(timeoutVal.ToInteger()) * time.Second
				}
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error":   err.Error(),
				"output":  string(output),
				"success": false,
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"output":  string(output),
			"code":    cmd.ProcessState.ExitCode(),
		})
	})

	// 列出运行中的进程
	sb.vm.Set("listProcesses", func() goja.Value {
		processes, err := process.Processes()
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("获取进程列表失败: %v", err),
			})
		}

		var processList []map[string]interface{}
		for _, p := range processes {
			name, _ := p.Name()
			pid := p.Pid
			status, _ := p.Status()
			memInfo, _ := p.MemoryInfo()
			cpuPercent, _ := p.CPUPercent()

			procInfo := map[string]interface{}{
				"pid":        pid,
				"name":       name,
				"status":     status,
				"cpuPercent": cpuPercent,
			}

			if memInfo != nil {
				procInfo["memory"] = memInfo.RSS
			}

			processList = append(processList, procInfo)
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success":   true,
			"processes": processList,
			"count":     len(processList),
		})
	})

	// 终止进程
	sb.vm.Set("killProcess", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供进程ID参数",
			})
		}

		pid := int32(call.Arguments[0].ToInteger())
		p, err := process.NewProcess(pid)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("找不到进程: %v", err),
			})
		}

		err = p.Kill()
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("终止进程失败: %v", err),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
		})
	})
}
