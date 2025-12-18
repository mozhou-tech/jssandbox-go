package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/supacloud/jssandbox-go/jssandbox"
)

// getMapKeys 获取 map 的所有键，用于调试
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// BidInfo 招标信息结构
type BidInfo struct {
	Index       string `json:"index"`       // 序号
	Project     string `json:"project"`     // 项目名称
	Section     string `json:"section"`     // 标段名称
	Region      string `json:"region"`      // 所在地区
	PublishTime string `json:"publishTime"` // 发布时间
	URL         string `json:"url"`         // 详情链接
}

func main() {
	// 设置日志级别为 Debug，以便看到详细的调试信息
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	ctx := context.Background()
	logrus.Info("初始化沙箱环境...")
	// 配置沙箱以显示浏览器窗口
	config := jssandbox.DefaultConfig().WithHeadless(false)
	sandbox := jssandbox.NewSandboxWithConfig(ctx, config)
	defer sandbox.Close()
	logrus.Info("沙箱环境初始化完成（浏览器窗口将显示）")

	url := "http://jsggzy.jszwfw.gov.cn/jyxx/tradeInfonew.html"

	fmt.Println("开始爬取招标信息...")
	fmt.Printf("目标URL: %s\n\n", url)

	logrus.Info("准备执行 JavaScript 代码...")

	// 使用JavaScript代码进行爬取
	jsCode := fmt.Sprintf(`
		(function() {
		// 创建浏览器会话
		console.log("[DEBUG] 开始创建浏览器会话...");
		var session = createBrowserSession(60);
		console.log("[DEBUG] 浏览器会话创建成功");
		
		try {
			// 导航到目标页面
			console.log("[DEBUG] 正在导航到页面: %s");
			var navResult = session.navigate("%s");
			console.log("[DEBUG] 导航结果:", JSON.stringify(navResult));
			if (!navResult.success) {
				console.error("[ERROR] 导航失败:", navResult.error);
				throw new Error("导航失败: " + navResult.error);
			}
			console.log("[DEBUG] 导航成功");
			
			// 等待页面加载完成
			console.log("[DEBUG] 等待页面加载...");
			session.wait(3); // 等待3秒确保页面完全加载
			console.log("[DEBUG] 页面加载等待完成");
			
			// 获取页面HTML
			console.log("[DEBUG] 获取页面内容...");
			var htmlResult = session.getHTML();
			console.log("[DEBUG] HTML获取结果 - success:", htmlResult.success);
			if (!htmlResult.success) {
				console.error("[ERROR] 获取HTML失败:", htmlResult.error);
				throw new Error("获取HTML失败: " + htmlResult.error);
			}
			console.log("[DEBUG] HTML长度:", htmlResult.html ? htmlResult.html.length : 0);
			
			// 使用 goquery 解析 HTML
			console.log("[DEBUG] 开始使用 goquery 解析 HTML...");
			var doc = parseHTML(htmlResult.html);
			if (doc.error) {
				console.error("[ERROR] 解析 HTML 失败:", doc.error);
				throw new Error("解析 HTML 失败: " + doc.error);
			}
			console.log("[DEBUG] HTML 解析成功");
			
			// 构建完整 URL 的辅助函数
			function buildFullUrl(href, baseUrl, currentPath) {
				if (!href || href === '#' || href === 'javascript:void(0)' || href === 'javascript:;') {
					return '';
				}
				if (href.indexOf('http://') === 0 || href.indexOf('https://') === 0) {
					return href;
				}
				if (href.indexOf('//') === 0) {
					return 'http:' + href;
				}
				if (href.indexOf('/') === 0) {
					return baseUrl + href;
				}
				return baseUrl + currentPath + href;
			}
			
			// 获取当前页面的基础 URL
			var urlResult = session.getURL();
			var baseUrl = '';
			var currentPath = '';
			if (urlResult.success && urlResult.url) {
				var url = urlResult.url;
				var match = url.match(/^(https?:\/\/[^\/]+)/);
				if (match) {
					baseUrl = match[1];
				}
				var pathMatch = url.match(/^(https?:\/\/[^\/]+)(\/[^?#]*)/);
				if (pathMatch) {
					currentPath = pathMatch[2].substring(0, pathMatch[2].lastIndexOf('/') + 1);
				}
			}
			
			// 提取链接的辅助函数
			function extractLink(sel) {
				if (!sel || typeof sel.length !== 'function' || sel.length() === 0) {
					return '';
				}
				// 尝试查找链接
				var linkSel = sel.find('a').first();
				if (linkSel.length() > 0) {
					var href = linkSel.attr('href');
					if (href && href !== '#' && href !== 'javascript:void(0)' && href !== 'javascript:;') {
						return buildFullUrl(href, baseUrl, currentPath);
					}
				}
				// 尝试从 onclick 属性中提取
				var onclick = sel.attr('onclick');
				if (onclick) {
					var patterns = [
						/["']([^"']*\.(html|htm|jsp|aspx|php))["']/i,
						/["']([^"']*\/detail[^"']*)["']/i,
						/["']([^"']*\/info[^"']*)["']/i
					];
					for (var i = 0; i < patterns.length; i++) {
						var match = onclick.match(patterns[i]);
						if (match && match[1]) {
							return buildFullUrl(match[1], baseUrl, currentPath);
						}
					}
				}
				// 尝试从 data-url 或 data-href 属性中提取
				var dataUrl = sel.attr('data-url') || sel.attr('data-href');
				if (dataUrl) {
					return buildFullUrl(dataUrl, baseUrl, currentPath);
				}
				// 尝试从父元素中查找链接
				var parent = sel.parent();
				if (parent.length() > 0) {
					var parentLink = parent.find('a').first();
					if (parentLink.length() > 0) {
						var href = parentLink.attr('href');
						if (href && href !== '#' && href !== 'javascript:void(0)' && href !== 'javascript:;') {
							return buildFullUrl(href, baseUrl, currentPath);
						}
					}
				}
				return '';
			}
			
			// 使用 goquery 提取数据
			var bidInfos = [];
			var tableSelectors = ['table', '.table', '#dataTable', 'tbody', '[class*="table"]', '[id*="table"]', '[class*="list"]'];
			var table = null;
			
			// 查找表格
			for (var i = 0; i < tableSelectors.length; i++) {
				var sel = doc.find(tableSelectors[i]);
				if (sel.length() > 0) {
					// 检查是否有行
					var rows = sel.find('tr');
					if (rows.length() > 0) {
						table = sel;
						console.log("[DEBUG] 找到表格，选择器:", tableSelectors[i], "行数:", rows.length());
						break;
					}
				}
			}
			
			if (table && table.length() > 0) {
				// 从表格中提取数据
				var rows = table.find('tr');
				rows.each(function(rowSel, index) {
					if (index === 0) return; // 跳过表头
					
					var cells = rowSel.find('td, th');
					if (cells.length() >= 2) {
						var projectCell = cells.eq(1).length() > 0 ? cells.eq(1) : cells.eq(0);
						var url = extractLink(projectCell) || extractLink(rowSel);
						
						var info = {
							index: cells.eq(0).text().trim() || (index).toString(),
							project: cells.eq(1).text().trim() || cells.eq(0).text().trim(),
							section: cells.length() > 2 ? cells.eq(2).text().trim() : '',
							region: cells.length() > 3 ? cells.eq(3).text().trim() : '',
							publishTime: cells.length() > 4 ? cells.eq(4).text().trim() : '',
							url: url
						};
						
						if (info.project && info.project.length > 0) {
							bidInfos.push(info);
						}
					}
				});
			} else {
				// 如果没有找到表格，尝试查找链接
				console.log("[DEBUG] 未找到表格，尝试查找链接...");
				var linkSelectors = ['a[href*="trade"]', 'a[href*="bid"]', 'a[href*="tender"]', 'a[href*="detail"]', 'a[href*="info"]'];
				var allLinks = doc.find(linkSelectors.join(', '));
				
				allLinks.each(function(linkSel, index) {
					var parent = linkSel.parent();
					if (parent.length() > 0) {
						var text = parent.text().trim();
						if (text.length > 10) {
							var href = linkSel.attr('href');
							var url = buildFullUrl(href, baseUrl, currentPath);
							bidInfos.push({
								index: (bidInfos.length + 1).toString(),
								project: linkSel.text().trim() || text.substring(0, 50),
								section: '',
								region: '',
								publishTime: '',
								url: url
							});
						}
					}
				});
			}
			
			console.log("[DEBUG] 提取到的数据数量:", bidInfos.length);
			
			// 返回结果
			console.log("[DEBUG] 准备返回结果，数据数量:", bidInfos ? bidInfos.length : 0);
			return {
				success: true,
				data: bidInfos || [],
				count: bidInfos ? bidInfos.length : 0
			};
		} catch (error) {
			console.error("[ERROR] 捕获到异常:", error);
			console.error("[ERROR] 错误消息:", error.message || String(error));
			console.error("[ERROR] 错误堆栈:", error.stack || "无堆栈信息");
			return {
				success: false,
				error: error.message || String(error),
				data: [],
				count: 0
			};
		} finally {
			console.log("[DEBUG] 关闭浏览器会话");
			session.close();
		}
		})();
	`, url, url)

	logrus.Info("开始执行 JavaScript 代码...")
	result, err := sandbox.Run(jsCode)
	if err != nil {
		logrus.WithError(err).WithField("error_detail", err.Error()).Fatal("执行爬虫代码失败")
	}
	logrus.Info("JavaScript 代码执行完成，开始解析结果...")

	// 解析结果
	logrus.Debug("检查结果类型...")
	var resultMap map[string]interface{}
	if resultObj := result.ToObject(nil); resultObj != nil {
		logrus.Debug("结果是一个对象，开始导出...")
		exported := resultObj.Export()
		logrus.WithField("exported_type", fmt.Sprintf("%T", exported)).Debug("导出类型")
		if exportedMap, ok := exported.(map[string]interface{}); ok {
			resultMap = exportedMap
			logrus.WithField("keys", getMapKeys(exportedMap)).Debug("结果对象键")
		} else {
			logrus.WithField("type", fmt.Sprintf("%T", exported)).WithField("value", fmt.Sprintf("%+v", exported)).Fatal("结果格式错误: 期望对象")
		}
	} else {
		logrus.WithField("result_type", fmt.Sprintf("%T", result)).Fatal("无法解析结果对象")
	}

	logrus.Debug("检查执行结果...")
	success, ok := resultMap["success"].(bool)
	logrus.WithField("success", success).WithField("ok", ok).Debug("success 字段")
	if !ok || !success {
		logrus.Warn("执行未成功，检查错误信息...")

		// 打印详细的错误信息
		fmt.Println("\n" + strings.Repeat("=", 80))
		fmt.Println("❌ 爬取失败")
		fmt.Println(strings.Repeat("=", 80))

		if errMsg, ok := resultMap["error"].(string); ok {
			fmt.Printf("\n错误信息: %s\n", errMsg)
			logrus.WithField("error", errMsg).Error("爬取失败")

			// 检查是否是超时错误
			if strings.Contains(errMsg, "deadline exceeded") || strings.Contains(errMsg, "timeout") {
				fmt.Println("\n💡 提示: 这可能是网络超时或页面加载时间过长导致的。")
				fmt.Println("   建议:")
				fmt.Println("   1. 检查网络连接是否正常")
				fmt.Println("   2. 尝试增加浏览器会话的超时时间")
				fmt.Println("   3. 检查目标网站是否可以正常访问")
			}
		} else {
			fmt.Println("\n错误信息: 未知错误")
			logrus.WithField("full_result", fmt.Sprintf("%+v", resultMap)).Error("爬取失败: 未知错误")
		}

		// 打印完整结果用于调试
		fmt.Println("\n完整结果:")
		jsonData, _ := json.MarshalIndent(resultMap, "", "  ")
		fmt.Println(string(jsonData))
		fmt.Println(strings.Repeat("=", 80) + "\n")

		os.Exit(1)
	}
	logrus.Info("执行成功，开始处理数据...")

	data := resultMap["data"]
	var count float64
	if countVal, ok := resultMap["count"].(float64); ok {
		count = countVal
	} else if countVal, ok := resultMap["count"].(int64); ok {
		count = float64(countVal)
	} else if countVal, ok := resultMap["count"].(int); ok {
		count = float64(countVal)
	}

	fmt.Printf("成功爬取 %d 条招标信息\n\n", int(count))

	// 将数据转换为JSON格式
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logrus.WithError(err).Fatal("JSON序列化失败")
	}

	// 保存到文件
	outputFile := fmt.Sprintf("bid_info_%s.json", time.Now().Format("20060102_150405"))
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		logrus.WithError(err).Fatal("保存文件失败")
	}

	fmt.Printf("数据已保存到: %s\n\n", outputFile)

	// 打印前几条数据作为预览
	if dataArray, ok := data.([]interface{}); ok && len(dataArray) > 0 {
		fmt.Println("数据预览（前5条）:")
		fmt.Println(strings.Repeat("=", 80))
		for i, item := range dataArray {
			if i >= 5 {
				break
			}
			if itemMap, ok := item.(map[string]interface{}); ok {
				fmt.Printf("\n[%d]\n", i+1)
				if project, ok := itemMap["project"].(string); ok {
					fmt.Printf("  项目名称: %s\n", project)
				}
				if section, ok := itemMap["section"].(string); ok && section != "" {
					fmt.Printf("  标段名称: %s\n", section)
				}
				if region, ok := itemMap["region"].(string); ok && region != "" {
					fmt.Printf("  所在地区: %s\n", region)
				}
				if publishTime, ok := itemMap["publishTime"].(string); ok && publishTime != "" {
					fmt.Printf("  发布时间: %s\n", publishTime)
				}
				if url, ok := itemMap["url"].(string); ok && url != "" {
					fmt.Printf("  详情链接: %s\n", url)
				}
			}
		}
		fmt.Println("\n" + strings.Repeat("=", 80))
	}
}
