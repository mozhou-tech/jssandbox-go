package jssandbox

import (
	"path/filepath"

	"github.com/dop251/goja"
)

// registerPath 注册路径处理增强功能到JavaScript运行时
func (sb *Sandbox) registerPath() {
	// 路径拼接
	sb.vm.Set("pathJoin", func(call goja.FunctionCall) goja.Value {
		var paths []string
		for _, arg := range call.Arguments {
			if arg != nil && !goja.IsUndefined(arg) {
				paths = append(paths, arg.String())
			}
		}

		if len(paths) == 0 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供至少一个路径参数",
			})
		}

		joined := filepath.Join(paths...)
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    joined,
		})
	})

	// 获取目录
	sb.vm.Set("pathDir", func(path string) goja.Value {
		dir := filepath.Dir(path)
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"dir":     dir,
		})
	})

	// 获取文件名
	sb.vm.Set("pathBase", func(path string) goja.Value {
		base := filepath.Base(path)
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"base":    base,
		})
	})

	// 获取扩展名
	sb.vm.Set("pathExt", func(path string) goja.Value {
		ext := filepath.Ext(path)
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"ext":     ext,
		})
	})

	// 获取绝对路径
	sb.vm.Set("pathAbs", func(path string) goja.Value {
		abs, err := filepath.Abs(path)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "获取绝对路径失败: " + err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    abs,
		})
	})
}
