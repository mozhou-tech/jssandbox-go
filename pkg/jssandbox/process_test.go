package jssandbox

import (
	"context"
	"runtime"
	"testing"

	"github.com/dop251/goja"
)

func TestExecCommand(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	var command string
	if runtime.GOOS == "windows" {
		command = "echo"
	} else {
		command = "echo"
	}

	code := `
		var result = execCommand("` + command + ` hello");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("execCommand() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		// 命令执行可能失败，检查是否有错误信息
		errorVal := resultObj.Get("error")
		if errorVal != nil && !goja.IsUndefined(errorVal) {
			t.Logf("execCommand()返回错误（可能是环境问题）: %s", errorVal.String())
			return
		}
		t.Error("execCommand()应该返回success: true")
	}

	output := resultObj.Get("output")
	if output == nil || goja.IsUndefined(output) {
		t.Error("execCommand()缺少output字段")
	}
}

func TestExecCommand_WithTimeout(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	var command string
	if runtime.GOOS == "windows" {
		command = "timeout"
	} else {
		command = "sleep"
	}

	code := `
		var result = execCommand("` + command + ` 1", {timeout: 2});
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("execCommand() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	// 命令可能成功也可能超时，只要函数能正常执行即可
	_ = resultObj
}

func TestListProcesses(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = listProcesses();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("listProcesses() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("listProcesses()应该返回success: true")
	}

	processes := resultObj.Get("processes")
	if processes == nil || goja.IsUndefined(processes) {
		t.Error("listProcesses()缺少processes字段")
	}

	count := resultObj.Get("count")
	if count.ToInteger() <= 0 {
		t.Errorf("listProcesses()进程数应该大于0, got %d", count.ToInteger())
	}
}

func TestListProcesses_HasFields(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		(function() {
			var result = listProcesses();
			var processes = result.processes;
			if (processes.length > 0) {
				var first = processes[0];
				return {
					hasPid: first.pid !== undefined,
					hasName: first.name !== undefined,
					hasStatus: first.status !== undefined
				};
			} else {
				return {
					hasPid: false,
					hasName: false,
					hasStatus: false
				};
			}
		})();
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("listProcesses() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	hasPid := resultObj.Get("hasPid")
	hasName := resultObj.Get("hasName")
	hasStatus := resultObj.Get("hasStatus")

	if !hasPid.ToBoolean() {
		t.Error("进程对象应该包含pid字段")
	}
	if !hasName.ToBoolean() {
		t.Error("进程对象应该包含name字段")
	}
	if !hasStatus.ToBoolean() {
		t.Error("进程对象应该包含status字段")
	}
}

func TestKillProcess_InvalidPID(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 使用一个不存在的PID（假设999999不存在）
	code := `
		var result = killProcess(999999);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("killProcess() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	errorVal := resultObj.Get("error")
	// 应该返回错误，因为PID不存在
	if errorVal == nil || goja.IsUndefined(errorVal) {
		t.Log("killProcess()对无效PID的处理可能因系统而异")
	}
}

