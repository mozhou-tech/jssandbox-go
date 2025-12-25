package jssandbox

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dop251/goja"
	"gopkg.in/yaml.v3"
)

// registerEnv 注册环境变量和配置功能到JavaScript运行时
func (sb *Sandbox) registerEnv() {
	// 获取环境变量
	sb.vm.Set("getEnv", func(name string) goja.Value {
		value := os.Getenv(name)
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"value":   value,
			"exists":  value != "",
		})
	})

	// 获取所有环境变量
	sb.vm.Set("getEnvAll", func() goja.Value {
		env := make(map[string]string)
		for _, e := range os.Environ() {
			for i := 0; i < len(e); i++ {
				if e[i] == '=' {
					key := e[:i]
					value := e[i+1:]
					env[key] = value
					break
				}
			}
		}
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"env":     env,
		})
	})

	// 读取配置文件（支持JSON和YAML）
	sb.vm.Set("readConfig", func(filePath string) goja.Value {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("读取文件失败: %v", err),
			})
		}

		ext := filepath.Ext(filePath)
		var config map[string]interface{}

		switch ext {
		case ".yaml", ".yml":
			if err := yaml.Unmarshal(data, &config); err != nil {
				return sb.vm.ToValue(map[string]interface{}{
					"error": fmt.Sprintf("解析YAML失败: %v", err),
				})
			}
		case ".json":
			// 使用goja的JSON解析
			jsonVal, err := sb.vm.RunString("JSON.parse(" + fmt.Sprintf("%q", string(data)) + ")")
			if err != nil {
				return sb.vm.ToValue(map[string]interface{}{
					"error": fmt.Sprintf("解析JSON失败: %v", err),
				})
			}
			if configObj := jsonVal.ToObject(sb.vm); configObj != nil {
				if exported, ok := configObj.Export().(map[string]interface{}); ok {
					config = exported
				} else {
					return sb.vm.ToValue(map[string]interface{}{
						"error": "配置文件格式错误: 期望对象",
					})
				}
			}
		default:
			// 尝试YAML
			if err := yaml.Unmarshal(data, &config); err != nil {
				// 尝试JSON
				jsonVal, err2 := sb.vm.RunString("JSON.parse(" + fmt.Sprintf("%q", string(data)) + ")")
				if err2 != nil {
					return sb.vm.ToValue(map[string]interface{}{
						"error": fmt.Sprintf("无法解析配置文件: YAML错误=%v, JSON错误=%v", err, err2),
					})
				}
				if configObj := jsonVal.ToObject(sb.vm); configObj != nil {
					if exported, ok := configObj.Export().(map[string]interface{}); ok {
						config = exported
					} else {
						return sb.vm.ToValue(map[string]interface{}{
							"error": "配置文件格式错误: 期望对象",
						})
					}
				}
			}
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"config":  config,
		})
	})
}
