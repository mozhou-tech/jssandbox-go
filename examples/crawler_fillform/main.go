package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mozhou-tech/jssandbox-go/jssandbox"
	"github.com/sirupsen/logrus"
)

func main() {
	// 设置日志级别为 Debug，显示详细信息
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
		DisableColors: true,
	})

	ctx := context.Background()
	logrus.Info("初始化沙箱环境...")

	// 配置沙箱（headless模式，不显示浏览器窗口）
	config := jssandbox.DefaultConfig().WithHeadless(false)
	sandbox := jssandbox.NewSandboxWithConfig(ctx, config)
	defer sandbox.Close()
	logrus.Info("沙箱环境初始化完成")

	// 从命令行参数获取URL，如果没有则使用默认URL
	url := "http://jsggzy.jszwfw.gov.cn"
	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("使用 GoQuery 提取表格内容（先点击'交易信息'链接）")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("目标URL: %s\n", url)
	fmt.Printf("会话超时: 120秒\n")
	fmt.Printf("流程: 导航 -> 点击'交易信息'链接 -> 提取表格内容\n\n")
	logrus.WithField("url", url).Debug("准备访问URL")

	// 使用JavaScript代码，结合 goquery 提取表格内容
	jsCode := fmt.Sprintf(`
		(function() {
		var session = createBrowserSession(120);
		try {
			console.log("[DEBUG] 开始导航...");
			var navResult = session.navigate("%s");
			console.log("[DEBUG] 导航结果:", JSON.stringify(navResult));
			if (!navResult.success) throw new Error("导航失败: " + navResult.error);
			console.log("[DEBUG] 导航成功，等待页面稳定...");
			session.wait(3);
			
			console.log("[DEBUG] 查找并点击'交易信息'链接...");
			// 使用JavaScript查找包含"交易信息"文本的链接并点击
			var clickResult = session.evaluate(
				'(function() {' +
				'  var links = document.querySelectorAll("a");' +
				'  for (var i = 0; i < links.length; i++) {' +
				'    var link = links[i];' +
				'    var text = link.innerText || link.textContent || "";' +
				'    if (text.indexOf("交易信息") !== -1) {' +
				'      link.click();' +
				'      return true;' +
				'    }' +
				'  }' +
				'  return false;' +
				'})();'
			);
			
			if (!clickResult.success || !clickResult.result) {
				// 如果直接点击失败，尝试通过href查找并导航
				var findLinkResult = session.evaluate(
					'(function() {' +
					'  var links = document.querySelectorAll("a");' +
					'  for (var i = 0; i < links.length; i++) {' +
					'    var link = links[i];' +
					'    var text = link.innerText || link.textContent || "";' +
					'    if (text.indexOf("交易信息") !== -1) {' +
					'      return link.href || "";' +
					'    }' +
					'  }' +
					'  return "";' +
					'})();'
				);
				
				if (findLinkResult.success && findLinkResult.result) {
					var linkHref = findLinkResult.result;
					console.log("[DEBUG] 找到链接，导航到:", linkHref);
					var navResult = session.navigate(linkHref);
					if (!navResult.success) {
						throw new Error("导航到'交易信息'链接失败: " + navResult.error);
					}
					console.log("[DEBUG] 已导航到'交易信息'页面");
				} else {
					// 尝试使用CSS选择器
					var selectors = ['a[href*="tradeInfo"]', 'a[href*="jyxx"]', 'a[href*="交易信息"]'];
					var found = false;
					for (var i = 0; i < selectors.length; i++) {
						try {
							var cssClickResult = session.click(selectors[i]);
							if (cssClickResult.success) {
								console.log("[DEBUG] 使用选择器点击成功:", selectors[i]);
								found = true;
								break;
							}
						} catch (e) {
							console.log("[DEBUG] 选择器点击失败:", selectors[i]);
						}
					}
					if (!found) {
						throw new Error("未找到'交易信息'链接");
					}
				}
			} else {
				console.log("[DEBUG] 已点击'交易信息'链接");
			}
			
			console.log("[DEBUG] 等待页面加载...");
			session.wait(5);
			
			console.log("[DEBUG] 获取页面HTML...");
			var htmlResult = session.getHTML();
			if (!htmlResult.success) throw new Error("获取HTML失败: " + htmlResult.error);
			
			var html = htmlResult.html || "";
			console.log("[DEBUG] HTML长度:", html.length);
			
			console.log("[DEBUG] 使用 GoQuery 解析HTML...");
			var doc = parseHTML(html);
			if (doc.error) throw new Error("解析HTML失败: " + doc.error);
			
			console.log("[DEBUG] 查找表格元素...");
			// 尝试多种表格选择器
			var table = null;
			var selectors = ['table', '.table', '#dataTable', 'tbody', '[class*="table"]', '[id*="table"]'];
			for (var i = 0; i < selectors.length; i++) {
				var found = doc.find(selectors[i]);
				if (found.length() > 0) {
					table = found;
					console.log("[DEBUG] 找到表格，选择器:", selectors[i], "行数:", found.find('tr').length());
					break;
				}
			}
			
			if (!table || table.length() === 0) {
				console.log("[DEBUG] 未找到表格，尝试查找所有包含行的容器...");
				var containers = doc.find('div, section, article');
				var rows = [];
				containers.each(function(container, index) {
					var trs = container.find('tr');
					if (trs.length() > 0 && trs.length() < 100) {
						rows.push({
							container: container,
							rowCount: trs.length()
						});
					}
				});
				if (rows.length > 0) {
					// 选择行数最多的容器
					rows.sort(function(a, b) { return b.rowCount - a.rowCount; });
					table = rows[0].container;
					console.log("[DEBUG] 找到包含表格的容器，行数:", rows[0].rowCount);
				}
			}
			
			if (!table || table.length() === 0) {
				throw new Error("未找到表格元素");
			}
			
			console.log("[DEBUG] 提取表格数据...");
			var tableData = [];
			var rows = table.find('tr');
			var rowIndex = 0;
			
			rows.each(function(row) {
				var cells = row.find('td, th');
				if (cells.length() === 0) {
					rowIndex++;
					return;
				}
				
				var rowData = {
					index: rowIndex,
					cells: []
				};
				
				cells.each(function(cell) {
					var text = cell.text();
					// 去除首尾空白
					if (text && typeof text.trim === 'function') {
						text = text.trim();
					} else if (text) {
						text = String(text).replace(/^\s+|\s+$/g, '');
					}
					rowData.cells.push(text || '');
				});
				
				// 只添加有内容的行
				var hasContent = false;
				for (var i = 0; i < rowData.cells.length; i++) {
					if (rowData.cells[i] && rowData.cells[i].length > 0) {
						hasContent = true;
						break;
					}
				}
				if (hasContent) {
					tableData.push(rowData);
				}
				rowIndex++;
			});
			
			console.log("[DEBUG] 提取到", tableData.length, "行数据");
			
			// 提取表头（如果有）
			var headers = [];
			var firstRow = table.find('tr').first();
			if (firstRow.length() > 0) {
				var headerCells = firstRow.find('th, td');
				headerCells.each(function(cell) {
					var text = cell.text();
					if (text && typeof text.trim === 'function') {
						text = text.trim();
					} else if (text) {
						text = String(text).replace(/^\s+|\s+$/g, '');
					}
					headers.push(text || '');
				});
			}
			
			return {
				success: true,
				rowCount: tableData.length,
				headers: headers,
				data: tableData
			};
		} catch (error) {
			console.error("[ERROR] 发生异常:", error.message || String(error));
			return { success: false, error: error.message || String(error) };
		} finally {
			console.log("[DEBUG] 关闭浏览器会话");
			session.close();
		}
		})();
	`, url)

	logrus.Info("开始执行 JavaScript 代码...")
	logrus.WithField("code_length", len(jsCode)).Debug("JavaScript代码长度")

	// 使用超时执行，避免程序卡住（150秒超时，给浏览器足够时间）
	result, err := sandbox.RunWithTimeout(jsCode, 150*time.Second)
	if err != nil {
		logrus.WithError(err).WithField("error_type", fmt.Sprintf("%T", err)).Error("执行代码失败")
		fmt.Printf("\n执行错误详情:\n")
		fmt.Printf("  错误类型: %T\n", err)
		fmt.Printf("  错误信息: %v\n", err)

		// 检查是否是超时错误
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "超时") {
			fmt.Printf("\n💡 提示: 执行超时，可能是页面加载时间过长或网络问题\n")
			fmt.Printf("   建议:\n")
			fmt.Printf("   1. 检查网络连接\n")
			fmt.Printf("   2. 尝试增加超时时间\n")
			fmt.Printf("   3. 检查目标网站是否可以正常访问\n")
		}

		logrus.Fatal("执行代码失败")
	}
	logrus.Debug("JavaScript代码执行完成，开始解析结果")

	// 解析结果
	logrus.Debug("开始解析结果对象...")
	var resultMap map[string]interface{}
	if resultObj := result.ToObject(nil); resultObj != nil {
		logrus.Debug("结果对象获取成功，开始导出...")
		exported := resultObj.Export()
		logrus.WithField("exported_type", fmt.Sprintf("%T", exported)).Debug("导出类型")

		if exportedMap, ok := exported.(map[string]interface{}); ok {
			resultMap = exportedMap
			logrus.WithField("keys", getMapKeys(exportedMap)).Debug("结果对象包含的键")
		} else {
			logrus.WithField("type", fmt.Sprintf("%T", exported)).Error("结果格式错误: 期望对象")
			logrus.Fatal("结果格式错误: 期望对象")
		}
	} else {
		logrus.WithField("result_type", fmt.Sprintf("%T", result)).Error("无法解析结果对象")
		logrus.Fatal("无法解析结果对象")
	}

	// 打印完整结果用于调试
	logrus.WithField("result", fmt.Sprintf("%+v", resultMap)).Debug("完整结果对象")

	// 检查执行结果
	success, ok := resultMap["success"].(bool)
	logrus.WithField("success", success).WithField("success_ok", ok).Debug("检查success字段")

	if !ok || !success {
		fmt.Println("\n" + strings.Repeat("=", 80))
		fmt.Println("❌ 爬取失败")
		fmt.Println(strings.Repeat("=", 80))

		if errMsg, ok := resultMap["error"].(string); ok {
			fmt.Printf("\n错误信息: %s\n", errMsg)
			logrus.WithField("error", errMsg).Error("爬取失败")
		} else {
			fmt.Println("\n错误信息: 未知错误")
			logrus.WithField("result", resultMap).Error("爬取失败: 未知错误")
		}

		fmt.Println("\n完整结果对象:")
		fmt.Printf("%+v\n", resultMap)
		fmt.Println(strings.Repeat("=", 80) + "\n")
		os.Exit(1)
	}

	logrus.Info("执行成功，开始处理结果...")

	// 获取表格数据
	rowCount, _ := resultMap["rowCount"].(float64)
	headers, _ := resultMap["headers"].([]interface{})
	data, _ := resultMap["data"].([]interface{})

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("提取结果:\n")
	fmt.Printf("  表格行数: %.0f\n", rowCount)

	if len(headers) > 0 {
		fmt.Printf("  表头: %s\n", strings.Join(convertToStringSlice(headers), " | "))
	}

	fmt.Printf("\n表格内容预览 (前10行):\n")
	fmt.Println(strings.Repeat("-", 80))

	// 显示前10行数据
	displayCount := 10
	if len(data) < displayCount {
		displayCount = len(data)
	}

	for i := 0; i < displayCount; i++ {
		if row, ok := data[i].(map[string]interface{}); ok {
			if cells, ok := row["cells"].([]interface{}); ok {
				cellTexts := convertToStringSlice(cells)
				fmt.Printf("第%d行: %s\n", i+1, strings.Join(cellTexts, " | "))
			}
		}
	}

	if len(data) > displayCount {
		fmt.Printf("\n... 还有 %d 行数据未显示\n", len(data)-displayCount)
	}

	fmt.Println(strings.Repeat("=", 80))
}

// convertToStringSlice 将 interface{} 切片转换为字符串切片
func convertToStringSlice(slice []interface{}) []string {
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if str, ok := v.(string); ok {
			result = append(result, str)
		} else {
			result = append(result, fmt.Sprintf("%v", v))
		}
	}
	return result
}

// getMapKeys 获取 map 的所有键，用于调试
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
