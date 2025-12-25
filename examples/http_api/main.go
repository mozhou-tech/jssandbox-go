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
	"github.com/mozhou-tech/jssandbox-go/jssandbox"
	"github.com/sirupsen/logrus"
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

	// 从命令行参数获取API URL，如果没有则使用默认API
	apiURL := "http://jsggzy.jszwfw.gov.cn/inteligentsearch/rest/esinteligentsearch/getFullTextDataNew"
	if len(os.Args) > 1 {
		apiURL = os.Args[1]
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("HTTP API 数据爬取示例")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("API URL: %s\n", apiURL)
	fmt.Printf("HTTP超时: 30秒\n\n")

	// 示例3: 使用 httpRequest 发送POST请求（带完整请求头）
	fmt.Println("示例: 使用 httpRequest 发送POST请求")
	fmt.Println(strings.Repeat("-", 80))
	example3(sandbox, apiURL)
}

// 示例3: POST请求（带完整请求头）
func example3(sandbox *jssandbox.Sandbox, url string) {
	jsCode := fmt.Sprintf(`
		(function() {
			try {
				console.log("[INFO] 开始发送POST请求...");
				
				var sortObj = {"infodatepx": "0"};
				var postData = {
					"token": "",
					"pn": 0,
					"rn": 20,
					"sdt": "",
					"edt": "",
					"wd": "",
					"inc_wd": "",
					"exc_wd": "",
					"fields": "title",
					"cnum": "001",
					"sort": JSON.stringify(sortObj),
					"ssort": "title",
					"cl": 200,
					"terminal": "",
					"condition": [],
					"time": [{
						"fieldName": "infodatepx",
						"startTime": "2025-12-15 00:00:00",
						"endTime": "2025-12-18 23:59:59"
					}],
					"highlights": "title",
					"statistics": null,
					"unionCondition": [],
					"accuracy": "",
					"noParticiple": "1",
					"searchRange": null,
					"isBusiness": "1"
				};
				
				var result = httpRequest("%s", {
					method: "POST",
					headers: {
						"Accept": "application/json, text/javascript, */*; q=0.01",
						"Accept-Encoding": "gzip, deflate",
						"Accept-Language": "zh,en-US;q=0.9,en;q=0.8,zh-CN;q=0.7",
						"Connection": "keep-alive",
						"Content-Type": "application/json;charset=UTF-8",
						"Origin": "http://jsggzy.jszwfw.gov.cn",
						"Referer": "http://jsggzy.jszwfw.gov.cn/jyxx/tradeInfonew.html",
						"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
						"X-Requested-With": "XMLHttpRequest"
					},
					body: JSON.stringify(postData),
					timeout: 30
				});
				
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
					console.log("[INFO] 响应数据解析成功");
					
					// 补全URL
					var baseUrl = "http://jsggzy.jszwfw.gov.cn";
					if (responseData.result && responseData.result.records && Array.isArray(responseData.result.records)) {
						console.log("[INFO] 查询到", responseData.result.records.length, "条数据");
						for (var i = 0; i < responseData.result.records.length; i++) {
							var record = responseData.result.records[i];
							if (record.linkurl && record.linkurl.indexOf("http") !== 0) {
								// 如果是相对路径，补全为完整URL
								if (record.linkurl.indexOf("/") === 0) {
									record.linkurl = baseUrl + record.linkurl;
								} else {
									record.linkurl = baseUrl + "/" + record.linkurl;
								}
							}
						}
						console.log("[INFO] URL补全完成");
					}
				} catch (e) {
					console.log("[WARN] 响应不是有效的JSON格式");
				}
				
				return {
					success: true,
					status: result.status,
					statusText: result.statusText,
					contentType: result.contentType,
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
