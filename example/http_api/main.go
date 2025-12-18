package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/supacloud/jssandbox-go/jssandbox"
)

func main() {
	// 设置日志级别
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
		DisableColors: true,
	})

	ctx := context.Background()
	logrus.Info("初始化沙箱环境...")

	// 配置沙箱（启用HTTP功能，禁用浏览器功能以节省资源）
	config := jssandbox.DefaultConfig().
		WithHTTPTimeout(30 * time.Second).
		DisableBrowser() // HTTP爬取不需要浏览器
	sandbox := jssandbox.NewSandboxWithConfig(ctx, config)
	defer sandbox.Close()
	logrus.Info("沙箱环境初始化完成")

	// 从命令行参数获取API URL，如果没有则使用示例API
	apiURL := "https://jsonplaceholder.typicode.com/posts"
	if len(os.Args) > 1 {
		apiURL = os.Args[1]
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("HTTP API 数据爬取示例")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("API URL: %s\n", apiURL)
	fmt.Printf("HTTP超时: 30秒\n\n")

	// 示例1: 使用 httpGet 简单GET请求
	fmt.Println("示例1: 使用 httpGet 获取数据")
	fmt.Println(strings.Repeat("-", 80))
	example1(sandbox, apiURL)

	// 示例2: 使用 httpRequest 带自定义headers的GET请求
	fmt.Println("\n示例2: 使用 httpRequest 带自定义headers")
	fmt.Println(strings.Repeat("-", 80))
	example2(sandbox, apiURL)

	// 示例3: 使用 httpPost 发送POST请求
	fmt.Println("\n示例3: 使用 httpPost 发送数据")
	fmt.Println(strings.Repeat("-", 80))
	example3(sandbox, "https://jsonplaceholder.typicode.com/posts")

	// 示例4: 使用 httpRequest 完整配置
	fmt.Println("\n示例4: 使用 httpRequest 完整配置")
	fmt.Println(strings.Repeat("-", 80))
	example4(sandbox, apiURL)
}

// 示例1: 简单的GET请求
func example1(sandbox *jssandbox.Sandbox, url string) {
	jsCode := fmt.Sprintf(`
		(function() {
			try {
				console.log("[INFO] 开始发送GET请求...");
				var result = httpGet("%s");
				
				if (result.error) {
					return {
						success: false,
						error: result.error
					};
				}
				
				console.log("[INFO] 请求成功，状态码:", result.status);
				console.log("[INFO] 响应体长度:", result.body ? result.body.length : 0);
				
				// 尝试解析JSON
				var data = null;
				try {
					data = JSON.parse(result.body);
				} catch (e) {
					console.log("[WARN] 响应不是有效的JSON格式");
				}
				
				return {
					success: true,
					status: result.status,
					statusText: result.statusText,
					contentType: result.contentType,
					data: data,
					bodyLength: result.body ? result.body.length : 0
				};
			} catch (error) {
				return {
					success: false,
					error: error.message || String(error)
				};
			}
		})();
	`, url)

	result, err := sandbox.Run(jsCode)
	if err != nil {
		logrus.WithError(err).Error("执行失败")
		return
	}

	printResult(sandbox, result)
}

// 示例2: 带自定义headers的GET请求
func example2(sandbox *jssandbox.Sandbox, url string) {
	jsCode := fmt.Sprintf(`
		(function() {
			try {
				console.log("[INFO] 开始发送带自定义headers的GET请求...");
				var result = httpRequest("%s", {
					method: "GET",
					headers: {
						"User-Agent": "jssandbox-http-client/1.0",
						"Accept": "application/json",
						"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8"
					},
					timeout: 30
				});
				
				if (result.error) {
					return {
						success: false,
						error: result.error
					};
				}
				
				console.log("[INFO] 请求成功，状态码:", result.status);
				
				// 解析JSON数组
				var data = null;
				try {
					data = JSON.parse(result.body);
					console.log("[INFO] 解析到", Array.isArray(data) ? data.length : 1, "条数据");
				} catch (e) {
					console.log("[WARN] 响应不是有效的JSON格式");
				}
				
				return {
					success: true,
					status: result.status,
					data: data,
					dataCount: Array.isArray(data) ? data.length : (data ? 1 : 0)
				};
			} catch (error) {
				return {
					success: false,
					error: error.message || String(error)
				};
			}
		})();
	`, url)

	result, err := sandbox.Run(jsCode)
	if err != nil {
		logrus.WithError(err).Error("执行失败")
		return
	}

	printResult(sandbox, result)
}

// 示例3: POST请求
func example3(sandbox *jssandbox.Sandbox, url string) {
	jsCode := fmt.Sprintf(`
		(function() {
			try {
				console.log("[INFO] 开始发送POST请求...");
				
				var postData = {
					title: "测试标题",
					body: "这是测试内容",
					userId: 1
				};
				
				var result = httpPost("%s", JSON.stringify(postData));
				
				if (result.error) {
					return {
						success: false,
						error: result.error
					};
				}
				
				console.log("[INFO] POST请求成功，状态码:", result.status);
				
				// 解析响应
				var responseData = null;
				try {
					responseData = JSON.parse(result.body);
					console.log("[INFO] 创建的资源ID:", responseData.id);
				} catch (e) {
					console.log("[WARN] 响应不是有效的JSON格式");
				}
				
				return {
					success: true,
					status: result.status,
					data: responseData
				};
			} catch (error) {
				return {
					success: false,
					error: error.message || String(error)
				};
			}
		})();
	`, url)

	result, err := sandbox.Run(jsCode)
	if err != nil {
		logrus.WithError(err).Error("执行失败")
		return
	}

	printResult(sandbox, result)
}

// 示例4: 使用 httpRequest 完整配置（包括错误处理）
func example4(sandbox *jssandbox.Sandbox, url string) {
	jsCode := fmt.Sprintf(`
		(function() {
			try {
				console.log("[INFO] 开始发送完整配置的HTTP请求...");
				
				var result = httpRequest("%s", {
					method: "GET",
					headers: {
						"User-Agent": "jssandbox-http-client/1.0",
						"Accept": "application/json"
					},
					timeout: 30
				});
				
				if (result.error) {
					console.error("[ERROR] 请求失败:", result.error);
					return {
						success: false,
						error: result.error,
						status: result.status || 0
					};
				}
				
				// 检查HTTP状态码
				if (result.status < 200 || result.status >= 300) {
					console.warn("[WARN] HTTP状态码异常:", result.status);
					return {
						success: false,
						error: "HTTP错误: " + result.statusText,
						status: result.status
					};
				}
				
				console.log("[INFO] 请求成功");
				console.log("[INFO] Content-Type:", result.contentType);
				console.log("[INFO] 响应头数量:", Object.keys(result.headers || {}).length);
				
				// 解析JSON数据
				var data = null;
				var isArray = false;
				try {
					data = JSON.parse(result.body);
					isArray = Array.isArray(data);
					console.log("[INFO] 数据类型:", isArray ? "数组" : "对象");
					console.log("[INFO] 数据条数:", isArray ? data.length : 1);
				} catch (e) {
					console.log("[WARN] 响应不是JSON格式，返回原始文本");
					data = result.body;
				}
				
				// 如果是数组，提取前几条数据作为预览
				var preview = null;
				if (isArray && data.length > 0) {
					preview = data.slice(0, Math.min(3, data.length));
				}
				
				return {
					success: true,
					status: result.status,
					statusText: result.statusText,
					contentType: result.contentType,
					headers: result.headers,
					data: data,
					preview: preview,
					totalCount: isArray ? data.length : (data ? 1 : 0)
				};
			} catch (error) {
				console.error("[ERROR] 发生异常:", error.message);
				return {
					success: false,
					error: error.message || String(error)
				};
			}
		})();
	`, url)

	result, err := sandbox.Run(jsCode)
	if err != nil {
		logrus.WithError(err).Error("执行失败")
		return
	}

	printResult(sandbox, result)
}

// printResult 打印结果
func printResult(sandbox *jssandbox.Sandbox, result goja.Value) {
	// 将goja.Value转换为map
	var resultMap map[string]interface{}

	if resultObj := result.ToObject(nil); resultObj != nil {
		exported := resultObj.Export()
		if exportedMap, ok := exported.(map[string]interface{}); ok {
			resultMap = exportedMap
		} else {
			fmt.Printf("结果类型: %T\n", exported)
			return
		}
	} else {
		fmt.Printf("无法转换为对象\n")
		return
	}

	// 检查是否成功
	success, _ := resultMap["success"].(bool)
	if !success {
		if errMsg, ok := resultMap["error"].(string); ok {
			fmt.Printf("❌ 失败: %s\n", errMsg)
		} else {
			fmt.Printf("❌ 失败: 未知错误\n")
		}
		if status, ok := resultMap["status"].(float64); ok {
			fmt.Printf("   HTTP状态码: %.0f\n", status)
		}
		return
	}

	// 打印成功信息
	fmt.Printf("✅ 请求成功\n")
	if status, ok := resultMap["status"].(float64); ok {
		fmt.Printf("   HTTP状态码: %.0f\n", status)
	}
	if statusText, ok := resultMap["statusText"].(string); ok {
		fmt.Printf("   状态文本: %s\n", statusText)
	}
	if contentType, ok := resultMap["contentType"].(string); ok {
		fmt.Printf("   Content-Type: %s\n", contentType)
	}

	// 打印数据统计
	if totalCount, ok := resultMap["totalCount"].(float64); ok {
		fmt.Printf("   数据条数: %.0f\n", totalCount)
	} else if dataCount, ok := resultMap["dataCount"].(float64); ok {
		fmt.Printf("   数据条数: %.0f\n", dataCount)
	}

	// 打印预览数据
	if preview, ok := resultMap["preview"].([]interface{}); ok && len(preview) > 0 {
		fmt.Printf("\n   数据预览（前%d条）:\n", len(preview))
		for i, item := range preview {
			if itemMap, ok := item.(map[string]interface{}); ok {
				fmt.Printf("   [%d] ", i+1)
				if id, ok := itemMap["id"].(float64); ok {
					fmt.Printf("ID: %.0f, ", id)
				}
				if title, ok := itemMap["title"].(string); ok {
					titlePreview := title
					if len(titlePreview) > 50 {
						titlePreview = titlePreview[:50] + "..."
					}
					fmt.Printf("标题: %s\n", titlePreview)
				}
			}
		}
	}

	// 保存完整数据到文件（如果有数据）
	if data, ok := resultMap["data"]; ok && data != nil {
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err == nil {
			// 查找 example/data 目录
			dataDir := findDataDir()
			if dataDir == "" {
				fmt.Printf("\n   ⚠️  无法找到 example/data 目录，跳过保存\n")
				return
			}

			// 确保目录存在
			if err := os.MkdirAll(dataDir, 0755); err != nil {
				fmt.Printf("\n   ⚠️  创建目录失败: %v\n", err)
				return
			}

			// 保存文件
			filename := fmt.Sprintf("api_data_%s.json", time.Now().Format("20060102_150405"))
			filePath := filepath.Join(dataDir, filename)
			err = os.WriteFile(filePath, jsonData, 0644)
			if err == nil {
				fmt.Printf("\n   完整数据已保存到: %s\n", filePath)
			} else {
				fmt.Printf("\n   ⚠️  保存文件失败: %v\n", err)
			}
		}
	}
}

// findDataDir 查找 example/data 目录
func findDataDir() string {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// 查找 example 目录
	currentDir := wd

	// 向上查找 example 目录，最多向上查找5级
	for i := 0; i < 5; i++ {
		testPath := filepath.Join(currentDir, "example", "data")
		examplePath := filepath.Join(currentDir, "example")

		// 检查 example 目录是否存在
		if info, err := os.Stat(examplePath); err == nil && info.IsDir() {
			return testPath
		}

		// 向上查找
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// 已经到达根目录
			break
		}
		currentDir = parent
	}

	// 如果没找到，使用当前目录下的 example/data
	return filepath.Join(wd, "example", "data")
}
