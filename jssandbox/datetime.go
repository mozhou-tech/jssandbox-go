package jssandbox

import (
	"fmt"
	"strings"
	"time"

	"github.com/dop251/goja"
)

// registerDateTime 注册日期时间增强功能到JavaScript运行时
func (sb *Sandbox) registerDateTime() {
	// 格式化日期
	sb.vm.Set("formatDate", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供日期和格式参数",
			})
		}

		dateStr := call.Arguments[0].String()
		format := call.Arguments[1].String()

		// 尝试解析日期字符串
		var t time.Time
		var err error

		// 尝试多种日期格式
		formats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05",
			"2006-01-02",
			time.RFC3339,
			time.RFC3339Nano,
		}

		for _, f := range formats {
			if t, err = time.Parse(f, dateStr); err == nil {
				break
			}
		}

		if err != nil {
			// 如果解析失败，使用当前时间
			t = time.Now()
		}

		// 转换格式字符串（Go格式 -> 常见格式）
		format = convertDateFormat(format)
		formatted := t.Format(format)

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"date":    formatted,
		})
	})

	// 解析日期字符串
	sb.vm.Set("parseDate", func(dateString string) goja.Value {
		// 尝试多种日期格式
		formats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05",
			"2006-01-02",
			time.RFC3339,
			time.RFC3339Nano,
		}

		for _, format := range formats {
			if t, err := time.Parse(format, dateString); err == nil {
				return sb.vm.ToValue(map[string]interface{}{
					"success":   true,
					"timestamp": t.Unix(),
					"date":      t.Format("2006-01-02 15:04:05"),
					"year":      t.Year(),
					"month":     int(t.Month()),
					"day":       t.Day(),
					"hour":      t.Hour(),
					"minute":    t.Minute(),
					"second":    t.Second(),
				})
			}
		}

		return sb.vm.ToValue(map[string]interface{}{
			"error": "无法解析日期字符串",
		})
	})

	// 日期加减
	sb.vm.Set("addDays", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供日期和天数参数",
			})
		}

		dateStr := call.Arguments[0].String()
		days := int(call.Arguments[1].ToInteger())

		// 解析日期
		var t time.Time
		var err error
		formats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02",
			time.RFC3339,
		}

		for _, format := range formats {
			if t, err = time.Parse(format, dateStr); err == nil {
				break
			}
		}

		if err != nil {
			t = time.Now()
		}

		// 加减天数
		newDate := t.AddDate(0, 0, days)

		return sb.vm.ToValue(map[string]interface{}{
			"success":   true,
			"date":      newDate.Format("2006-01-02 15:04:05"),
			"timestamp": newDate.Unix(),
		})
	})

	// 获取时区
	sb.vm.Set("getTimezone", func() goja.Value {
		tz, _ := time.Now().Zone()
		return sb.vm.ToValue(map[string]interface{}{
			"success":  true,
			"timezone": tz,
			"offset":   time.Now().Format("-07:00"),
		})
	})

	// 时区转换
	sb.vm.Set("convertTimezone", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供日期和时区参数",
			})
		}

		dateStr := call.Arguments[0].String()
		tzName := call.Arguments[1].String()

		// 解析日期
		var t time.Time
		var err error
		formats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02",
			time.RFC3339,
		}

		for _, format := range formats {
			if t, err = time.Parse(format, dateStr); err == nil {
				break
			}
		}

		if err != nil {
			t = time.Now()
		}

		// 加载时区
		loc, err := time.LoadLocation(tzName)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("无效的时区: %v", err),
			})
		}

		// 转换时区
		converted := t.In(loc)

		return sb.vm.ToValue(map[string]interface{}{
			"success":   true,
			"date":      converted.Format("2006-01-02 15:04:05"),
			"timestamp": converted.Unix(),
		})
	})
}

// convertDateFormat 转换日期格式字符串
// 将常见的格式占位符转换为Go的格式
func convertDateFormat(format string) string {
	// 常见的格式映射
	replacements := map[string]string{
		"YYYY": "2006",
		"MM":   "01",
		"DD":   "02",
		"HH":   "15",
		"mm":   "04",
		"ss":   "05",
		"yyyy": "2006",
		"MMM":  "Jan",
		"MMMM": "January",
	}

	result := format
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}

	return result
}
