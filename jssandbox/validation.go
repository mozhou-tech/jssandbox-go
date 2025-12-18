package jssandbox

import (
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/dop251/goja"
)

// registerValidation 注册数据验证功能到JavaScript运行时
func (sb *Sandbox) registerValidation() {
	// 验证邮箱格式
	sb.vm.Set("validateEmail", func(email string) goja.Value {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		valid := emailRegex.MatchString(email)
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"valid":   valid,
		})
	})

	// 验证URL格式
	sb.vm.Set("validateURL", func(urlStr string) goja.Value {
		parsedURL, err := url.Parse(urlStr)
		valid := err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"valid":   valid,
		})
	})

	// 验证IP地址
	sb.vm.Set("validateIP", func(ip string) goja.Value {
		parsedIP := net.ParseIP(ip)
		valid := parsedIP != nil
		isIPv4 := false
		isIPv6 := false
		if valid {
			if parsedIP.To4() != nil {
				isIPv4 = true
			} else {
				isIPv6 = true
			}
		}
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"valid":   valid,
			"isIPv4":  isIPv4,
			"isIPv6":  isIPv6,
		})
	})

	// 验证中国手机号
	sb.vm.Set("validatePhone", func(phone string) goja.Value {
		// 中国手机号正则：1开头，11位数字
		phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
		valid := phoneRegex.MatchString(strings.TrimSpace(phone))
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"valid":   valid,
		})
	})
}
