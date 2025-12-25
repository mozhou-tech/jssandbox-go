package jssandbox

import (
	"fmt"
	"net"
	"time"

	"github.com/dop251/goja"
)

// registerNetwork 注册网络工具功能到JavaScript运行时
func (sb *Sandbox) registerNetwork() {
	// DNS解析
	sb.vm.Set("resolveDNS", func(hostname string) goja.Value {
		ips, err := net.LookupIP(hostname)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("DNS解析失败: %v", err),
			})
		}

		var ipList []string
		for _, ip := range ips {
			ipList = append(ipList, ip.String())
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"ips":     ipList,
		})
	})

	// Ping测试（简化版，使用TCP连接测试）
	sb.vm.Set("ping", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供主机地址参数",
			})
		}

		host := call.Arguments[0].String()
		count := 4
		timeout := 3 * time.Second

		if len(call.Arguments) > 1 {
			count = int(call.Arguments[1].ToInteger())
			if count <= 0 {
				count = 4
			}
			if count > 10 {
				count = 10
			}
		}

		var successCount int
		var totalTime time.Duration

		for i := 0; i < count; i++ {
			start := time.Now()
			conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "80"), timeout)
			if err != nil {
				continue
			}
			conn.Close()
			duration := time.Since(start)
			successCount++
			totalTime += duration
		}

		avgTime := time.Duration(0)
		if successCount > 0 {
			avgTime = totalTime / time.Duration(successCount)
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success":     true,
			"host":        host,
			"sent":        count,
			"received":    successCount,
			"lost":        count - successCount,
			"averageTime": avgTime.Milliseconds(),
		})
	})

	// 检查端口是否开放
	sb.vm.Set("checkPort", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供主机和端口参数",
			})
		}

		host := call.Arguments[0].String()
		port := int(call.Arguments[1].ToInteger())

		timeout := 3 * time.Second
		if len(call.Arguments) > 2 {
			timeout = time.Duration(call.Arguments[2].ToInteger()) * time.Second
		}

		address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
		conn, err := net.DialTimeout("tcp", address, timeout)

		open := err == nil
		if conn != nil {
			conn.Close()
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"host":    host,
			"port":    port,
			"open":    open,
		})
	})
}
