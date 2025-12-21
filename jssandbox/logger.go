package jssandbox

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
)

// registerLogger 注册日志功能到JavaScript运行时
func (sb *Sandbox) registerLogger() {
	// 辅助函数：格式化日志参数
	formatLogArgs := func(call goja.FunctionCall) string {
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

	// 辅助函数：从对象中提取字段
	extractFields := func(obj goja.Value) logrus.Fields {
		fields := make(logrus.Fields)
		if obj == nil || goja.IsUndefined(obj) || goja.IsNull(obj) {
			return fields
		}

		if objObj := obj.ToObject(sb.vm); objObj != nil {
			exported := objObj.Export()
			if exportedMap, ok := exported.(map[string]interface{}); ok {
				for k, v := range exportedMap {
					fields[k] = v
				}
			}
		}
		return fields
	}

	// 创建 logger 对象
	loggerObj := sb.vm.NewObject()

	// Debug 级别日志
	loggerObj.Set("debug", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			sb.logger.Debug()
			return goja.Undefined()
		}

		// 检查第一个参数是否是对象（用于结构化日志）
		firstArg := call.Arguments[0]
		if firstArgObj := firstArg.ToObject(sb.vm); firstArgObj != nil {
			// 尝试作为结构化日志处理
			fields := extractFields(firstArg)
			if len(fields) > 0 && len(call.Arguments) > 1 {
				// 有字段和消息
				msg := formatLogArgs(goja.FunctionCall{
					Arguments: call.Arguments[1:],
				})
				sb.logger.WithFields(fields).Debug(msg)
			} else if len(fields) > 0 {
				// 只有字段，没有消息
				sb.logger.WithFields(fields).Debug()
			} else {
				// 普通消息
				sb.logger.Debug(formatLogArgs(call))
			}
		} else {
			// 普通消息
			sb.logger.Debug(formatLogArgs(call))
		}
		return goja.Undefined()
	})

	// Info 级别日志
	loggerObj.Set("info", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			sb.logger.Info()
			return goja.Undefined()
		}

		firstArg := call.Arguments[0]
		if firstArgObj := firstArg.ToObject(sb.vm); firstArgObj != nil {
			fields := extractFields(firstArg)
			if len(fields) > 0 && len(call.Arguments) > 1 {
				msg := formatLogArgs(goja.FunctionCall{
					Arguments: call.Arguments[1:],
				})
				sb.logger.WithFields(fields).Info(msg)
			} else if len(fields) > 0 {
				sb.logger.WithFields(fields).Info()
			} else {
				sb.logger.Info(formatLogArgs(call))
			}
		} else {
			sb.logger.Info(formatLogArgs(call))
		}
		return goja.Undefined()
	})

	// Warn 级别日志
	loggerObj.Set("warn", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			sb.logger.Warn()
			return goja.Undefined()
		}

		firstArg := call.Arguments[0]
		if firstArgObj := firstArg.ToObject(sb.vm); firstArgObj != nil {
			fields := extractFields(firstArg)
			if len(fields) > 0 && len(call.Arguments) > 1 {
				msg := formatLogArgs(goja.FunctionCall{
					Arguments: call.Arguments[1:],
				})
				sb.logger.WithFields(fields).Warn(msg)
			} else if len(fields) > 0 {
				sb.logger.WithFields(fields).Warn()
			} else {
				sb.logger.Warn(formatLogArgs(call))
			}
		} else {
			sb.logger.Warn(formatLogArgs(call))
		}
		return goja.Undefined()
	})

	// Error 级别日志
	loggerObj.Set("error", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			sb.logger.Error()
			return goja.Undefined()
		}

		firstArg := call.Arguments[0]
		if firstArgObj := firstArg.ToObject(sb.vm); firstArgObj != nil {
			fields := extractFields(firstArg)
			if len(fields) > 0 && len(call.Arguments) > 1 {
				msg := formatLogArgs(goja.FunctionCall{
					Arguments: call.Arguments[1:],
				})
				sb.logger.WithFields(fields).Error(msg)
			} else if len(fields) > 0 {
				sb.logger.WithFields(fields).Error()
			} else {
				sb.logger.Error(formatLogArgs(call))
			}
		} else {
			sb.logger.Error(formatLogArgs(call))
		}
		return goja.Undefined()
	})

	// Fatal 级别日志（注意：这不会真正退出程序，只是记录日志）
	loggerObj.Set("fatal", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			sb.logger.Error("FATAL: (fatal级别日志)")
			return goja.Undefined()
		}

		firstArg := call.Arguments[0]
		if firstArgObj := firstArg.ToObject(sb.vm); firstArgObj != nil {
			fields := extractFields(firstArg)
			if len(fields) > 0 && len(call.Arguments) > 1 {
				msg := formatLogArgs(goja.FunctionCall{
					Arguments: call.Arguments[1:],
				})
				sb.logger.WithFields(fields).Error("FATAL: " + msg)
			} else if len(fields) > 0 {
				sb.logger.WithFields(fields).Error("FATAL:")
			} else {
				sb.logger.Error("FATAL: " + formatLogArgs(call))
			}
		} else {
			sb.logger.Error("FATAL: " + formatLogArgs(call))
		}
		return goja.Undefined()
	})

	// Trace 级别日志
	loggerObj.Set("trace", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			sb.logger.Trace()
			return goja.Undefined()
		}

		firstArg := call.Arguments[0]
		if firstArgObj := firstArg.ToObject(sb.vm); firstArgObj != nil {
			fields := extractFields(firstArg)
			if len(fields) > 0 && len(call.Arguments) > 1 {
				msg := formatLogArgs(goja.FunctionCall{
					Arguments: call.Arguments[1:],
				})
				sb.logger.WithFields(fields).Trace(msg)
			} else if len(fields) > 0 {
				sb.logger.WithFields(fields).Trace()
			} else {
				sb.logger.Trace(formatLogArgs(call))
			}
		} else {
			sb.logger.Trace(formatLogArgs(call))
		}
		return goja.Undefined()
	})

	// 设置日志级别
	loggerObj.Set("setLevel", func(level string) goja.Value {
		levelStr := strings.ToLower(level)
		var logLevel logrus.Level
		switch levelStr {
		case "trace":
			logLevel = logrus.TraceLevel
		case "debug":
			logLevel = logrus.DebugLevel
		case "info":
			logLevel = logrus.InfoLevel
		case "warn", "warning":
			logLevel = logrus.WarnLevel
		case "error":
			logLevel = logrus.ErrorLevel
		case "fatal":
			logLevel = logrus.FatalLevel
		case "panic":
			logLevel = logrus.PanicLevel
		default:
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("无效的日志级别: %s", level),
			})
		}

		sb.logger.SetLevel(logLevel)
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"level":   levelStr,
		})
	})

	// 获取当前日志级别
	loggerObj.Set("getLevel", func() goja.Value {
		level := sb.logger.GetLevel()
		levelStr := level.String()
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"level":   levelStr,
		})
	})

	// 检查是否启用某个日志级别
	loggerObj.Set("isLevelEnabled", func(level string) goja.Value {
		levelStr := strings.ToLower(level)
		var logLevel logrus.Level
		switch levelStr {
		case "trace":
			logLevel = logrus.TraceLevel
		case "debug":
			logLevel = logrus.DebugLevel
		case "info":
			logLevel = logrus.InfoLevel
		case "warn", "warning":
			logLevel = logrus.WarnLevel
		case "error":
			logLevel = logrus.ErrorLevel
		case "fatal":
			logLevel = logrus.FatalLevel
		case "panic":
			logLevel = logrus.PanicLevel
		default:
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("无效的日志级别: %s", level),
			})
		}

		enabled := sb.logger.IsLevelEnabled(logLevel)
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"enabled": enabled,
		})
	})

	// 带字段的日志方法（返回一个可以链式调用的对象）
	loggerObj.Set("withFields", func(fieldsObj goja.Value) goja.Value {
		fields := extractFields(fieldsObj)
		entry := sb.logger.WithFields(fields)

		// 创建一个新的对象，包含所有日志级别方法
		fieldLoggerObj := sb.vm.NewObject()

		fieldLoggerObj.Set("debug", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				entry.Debug()
			} else {
				entry.Debug(formatLogArgs(call))
			}
			return goja.Undefined()
		})

		fieldLoggerObj.Set("info", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				entry.Info()
			} else {
				entry.Info(formatLogArgs(call))
			}
			return goja.Undefined()
		})

		fieldLoggerObj.Set("warn", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				entry.Warn()
			} else {
				entry.Warn(formatLogArgs(call))
			}
			return goja.Undefined()
		})

		fieldLoggerObj.Set("error", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				entry.Error()
			} else {
				entry.Error(formatLogArgs(call))
			}
			return goja.Undefined()
		})

		fieldLoggerObj.Set("fatal", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				entry.Error("FATAL:")
			} else {
				entry.Error("FATAL: " + formatLogArgs(call))
			}
			return goja.Undefined()
		})

		fieldLoggerObj.Set("trace", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				entry.Trace()
			} else {
				entry.Trace(formatLogArgs(call))
			}
			return goja.Undefined()
		})

		return fieldLoggerObj
	})

	// 注册 logger 对象到全局
	sb.vm.Set("logger", loggerObj)
}
