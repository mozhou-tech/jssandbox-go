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

	// 获取当前时间戳（秒）
	sb.vm.Set("getCurrentTimestamp", func() goja.Value {
		now := time.Now()
		return sb.vm.ToValue(map[string]interface{}{
			"success":     true,
			"timestamp":   now.Unix(),
			"timestampMs": now.UnixMilli(),
			"date":        now.Format("2006-01-02 15:04:05"),
		})
	})

	// 时间戳转日期格式
	sb.vm.Set("timestampToDate", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供时间戳参数",
			})
		}

		// 获取时间戳（可能是秒或毫秒）
		timestamp := call.Arguments[0].ToFloat()
		var t time.Time

		// 判断是秒还是毫秒（大于 1e10 认为是毫秒）
		if timestamp > 1e10 {
			t = time.UnixMilli(int64(timestamp))
		} else {
			t = time.Unix(int64(timestamp), 0)
		}

		// 如果提供了格式参数，使用指定格式
		format := "2006-01-02 15:04:05"
		if len(call.Arguments) >= 2 {
			format = call.Arguments[1].String()
			format = convertDateFormat(format)
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success":     true,
			"date":        t.Format(format),
			"timestamp":   t.Unix(),
			"timestampMs": t.UnixMilli(),
			"year":        t.Year(),
			"month":       int(t.Month()),
			"day":         t.Day(),
			"hour":        t.Hour(),
			"minute":      t.Minute(),
			"second":      t.Second(),
		})
	})

	// 格式化时间戳为指定格式
	sb.vm.Set("formatTimestamp", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供时间戳和格式参数",
			})
		}

		// 获取时间戳
		timestamp := call.Arguments[0].ToFloat()
		format := call.Arguments[1].String()

		// 判断是秒还是毫秒
		var t time.Time
		if timestamp > 1e10 {
			t = time.UnixMilli(int64(timestamp))
		} else {
			t = time.Unix(int64(timestamp), 0)
		}

		// 转换格式字符串
		format = convertDateFormat(format)
		formatted := t.Format(format)

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"date":    formatted,
		})
	})

	// 日期转时间戳
	sb.vm.Set("dateToTimestamp", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供日期参数",
			})
		}

		dateStr := call.Arguments[0].String()

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
			if t, err := time.Parse(format, dateStr); err == nil {
				return sb.vm.ToValue(map[string]interface{}{
					"success":     true,
					"timestamp":   t.Unix(),
					"timestampMs": t.UnixMilli(),
				})
			}
		}

		return sb.vm.ToValue(map[string]interface{}{
			"error": "无法解析日期字符串",
		})
	})

	// 获取时间戳的详细信息
	sb.vm.Set("getTimestampInfo", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供时间戳参数",
			})
		}

		timestamp := call.Arguments[0].ToFloat()
		var t time.Time

		// 判断是秒还是毫秒
		if timestamp > 1e10 {
			t = time.UnixMilli(int64(timestamp))
		} else {
			t = time.Unix(int64(timestamp), 0)
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success":     true,
			"timestamp":   t.Unix(),
			"timestampMs": t.UnixMilli(),
			"date":        t.Format("2006-01-02 15:04:05"),
			"iso8601":     t.Format(time.RFC3339),
			"year":        t.Year(),
			"month":       int(t.Month()),
			"day":         t.Day(),
			"hour":        t.Hour(),
			"minute":      t.Minute(),
			"second":      t.Second(),
			"weekday":     t.Weekday().String(),
			"yearday":     t.YearDay(),
		})
	})
}

// convertDateFormat 转换日期格式字符串
// 将常见的格式占位符转换为Go的格式
func convertDateFormat(format string) string {
	// 常见的格式映射（按长度从长到短排序，避免替换冲突）
	replacements := []struct {
		old string
		new string
	}{
		{"YYYY", "2006"},
		{"yyyy", "2006"},
		{"MMMM", "January"},
		{"MMM", "Jan"},
		{"MM", "01"},
		{"DD", "02"},
		{"dd", "02"},
		{"HH", "15"},
		{"hh", "03"}, // 12小时制
		{"mm", "04"},
		{"ss", "05"},
		{"SSS", "000"}, // 毫秒
		{"A", "PM"},    // AM/PM
		{"a", "pm"},    // am/pm
		{"Z", "-0700"}, // 时区偏移
		{"z", "MST"},   // 时区名称
	}

	result := format
	for _, r := range replacements {
		result = strings.ReplaceAll(result, r.old, r.new)
	}

	return result
}
