package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/supacloud/jssandbox-go/jssandbox"
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
	config := jssandbox.DefaultConfig().WithHeadless(true)
	sandbox := jssandbox.NewSandboxWithConfig(ctx, config)
	defer sandbox.Close()
	logrus.Info("沙箱环境初始化完成")

	// 从命令行参数获取URL，如果没有则使用默认URL
	url := "http://jsggzy.jszwfw.gov.cn/jyxx/tradeInfonew.html"
	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("开始爬取招标信息")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("目标URL: %s\n", url)
	fmt.Printf("会话超时: 120秒\n")
	fmt.Printf("Headless模式: true\n\n")
	logrus.WithField("url", url).Debug("准备访问URL")

	// 使用JavaScript代码进行爬取
	jsCode := fmt.Sprintf(`
		(function() {
		var session = createBrowserSession(120);
		try {
			var navResult = session.navigate("%s");
			if (!navResult.success) throw new Error("导航失败: " + navResult.error);
			session.wait(3);
			
			var extractCode = "(function() {" +
				"var baseUrl = window.location.protocol + '//' + window.location.host;" +
				"var currentPath = window.location.pathname.substring(0, window.location.pathname.lastIndexOf('/') + 1);" +
				"var pathPatterns = [" +
				"new RegExp('[\"\\']([/][^\"\\']*\\\\.(html|htm|jsp|aspx|php))[\"\\']', 'i')," +
				"new RegExp('[\"\\']([/][^\"\\']*)[\"\\']', 'i')," +
				"new RegExp('[\"\\']([^\"\\']*\\\\.(html|htm|jsp|aspx|php))[\"\\']', 'i')," +
				"new RegExp('[\"\\']([^\"\\']*\\\\/(detail|info)[^\"\\']*)[\"\\']', 'i')" +
				"];" +
				"function buildFullUrl(href) {" +
				"if (!href || href === '#' || /^javascript:/.test(href)) return '';" +
				"if (/^https?:/.test(href)) return href;" +
				"if (/^\\/\\//.test(href)) return window.location.protocol + href;" +
				"if (/^\\//.test(href)) return baseUrl + href;" +
				"return baseUrl + currentPath + href;" +
				"}" +
				"function extractPathFromOnclick(onclick) {" +
				"if (!onclick) return null;" +
				"for (var i = 0; i < pathPatterns.length; i++) {" +
				"var match = onclick.match(pathPatterns[i]);" +
				"if (match && match[1]) return match[1];" +
				"}" +
				"return null;" +
				"}" +
				"function extractLink(element) {" +
				"if (!element) return '';" +
				"var link = element.querySelector('a');" +
				"if (link) {" +
				"var href = link.href || link.getAttribute('href') || '';" +
				"if (href && !/^#|^javascript:/.test(href)) return buildFullUrl(href);" +
				"var path = extractPathFromOnclick(link.getAttribute('onclick'));" +
				"if (path) return buildFullUrl(path);" +
				"}" +
				"var path = extractPathFromOnclick(element.getAttribute('onclick'));" +
				"if (path) return buildFullUrl(path);" +
				"var dataUrl = element.getAttribute('data-url') || element.getAttribute('data-href');" +
				"if (dataUrl) return buildFullUrl(dataUrl);" +
				"if (element.parentElement) {" +
				"link = element.parentElement.querySelector('a');" +
				"if (link) {" +
				"href = link.href || link.getAttribute('href') || '';" +
				"if (href && !/^#|^javascript:/.test(href)) return buildFullUrl(href);" +
				"path = extractPathFromOnclick(link.getAttribute('onclick'));" +
				"if (path) return buildFullUrl(path);" +
				"}" +
				"}" +
				"return '';" +
				"}" +
				"function getTextContent(el) {" +
				"return el ? (el.innerText || el.textContent || '').trim() : '';" +
				"}" +
				"var table = null;" +
				"var selectors = ['table', '.table', '#dataTable', 'tbody', '[class*=\"table\"]', '[id*=\"table\"]', '[class*=\"list\"]'];" +
				"for (var i = 0; i < selectors.length && !table; i++) {" +
				"var els = document.querySelectorAll(selectors[i]);" +
				"for (var j = 0; j < els.length; j++) {" +
				"if (els[j].querySelectorAll('tr').length > 0) { table = els[j]; break; }" +
				"}" +
				"}" +
				"var results = [];" +
				"if (table) {" +
				"var rows = table.querySelectorAll('tr');" +
				"for (var i = 1; i < rows.length; i++) {" +
				"var cells = rows[i].querySelectorAll('td, th');" +
				"if (cells.length >= 2) {" +
				"var projectCell = cells[1] || cells[0];" +
				"var url = extractLink(projectCell) || extractLink(rows[i]);" +
				"var project = getTextContent(cells[1]) || getTextContent(cells[0]);" +
				"if (project) {" +
				"results.push({" +
				"index: getTextContent(cells[0]) || i.toString()," +
				"project: project," +
				"section: cells.length > 2 ? getTextContent(cells[2]) : ''," +
				"region: cells.length > 3 ? getTextContent(cells[3]) : ''," +
				"publishTime: cells.length > 4 ? getTextContent(cells[4]) : ''," +
				"url: url" +
				"});" +
				"}" +
				"}" +
				"}" +
				"} else {" +
				"var linkSelectors = ['a[href*=\"trade\"]', 'a[href*=\"bid\"]', 'a[href*=\"tender\"]', 'a[href*=\"detail\"]', 'a[href*=\"info\"]'];" +
				"var allLinks = [];" +
				"for (var i = 0; i < linkSelectors.length; i++) {" +
				"var links = document.querySelectorAll(linkSelectors[i]);" +
				"for (var j = 0; j < links.length; j++) allLinks.push(links[j]);" +
				"}" +
				"for (var i = 0; i < allLinks.length; i++) {" +
				"var link = allLinks[i];" +
				"var parent = link.parentElement;" +
				"while (parent && parent !== document.body) {" +
				"if (parent.tagName === 'TR' || parent.tagName === 'LI' || parent.tagName === 'DIV' || parent.tagName === 'TD') break;" +
				"parent = parent.parentElement;" +
				"}" +
				"if (parent) {" +
				"var text = getTextContent(parent);" +
				"if (text.length > 10) {" +
				"var href = link.href || link.getAttribute('href') || '';" +
				"results.push({" +
				"index: (results.length + 1).toString()," +
				"project: getTextContent(link) || text.substring(0, 50)," +
				"section: ''," +
				"region: ''," +
				"publishTime: ''," +
				"url: buildFullUrl(href)" +
				"});" +
				"}" +
				"}" +
				"}" +
				"return results;" +
				"})();";
			
			var extractResult = session.evaluate(extractCode);
			if (!extractResult.success) throw new Error("提取数据失败: " + extractResult.error);
			var bidInfos = extractResult.result;
			
			if (!bidInfos || bidInfos.length === 0) {
				var backupCode = "(function() {" +
					"var baseUrl = window.location.protocol + '//' + window.location.host;" +
					"var currentPath = window.location.pathname.substring(0, window.location.pathname.lastIndexOf('/') + 1);" +
					"var pathPatterns = [" +
					"new RegExp('[\"\\']([/][^\"\\']*\\\\.(html|htm|jsp|aspx|php))[\"\\']', 'i')," +
					"new RegExp('[\"\\']([/][^\"\\']*)[\"\\']', 'i')," +
					"new RegExp('[\"\\']([^\"\\']*\\\\.(html|htm|jsp|aspx|php))[\"\\']', 'i')" +
					"];" +
					"function buildFullUrl(href) {" +
					"if (!href || href === '#' || /^javascript:/.test(href)) return '';" +
					"if (/^https?:/.test(href)) return href;" +
					"if (/^\\/\\//.test(href)) return window.location.protocol + href;" +
					"if (/^\\//.test(href)) return baseUrl + href;" +
					"return baseUrl + currentPath + href;" +
					"}" +
					"function extractPathFromOnclick(onclick) {" +
					"if (!onclick) return null;" +
					"for (var i = 0; i < pathPatterns.length; i++) {" +
					"var match = onclick.match(pathPatterns[i]);" +
					"if (match && match[1]) return match[1];" +
					"}" +
					"return null;" +
					"}" +
					"function extractLink(element) {" +
					"if (!element) return '';" +
					"var link = element.querySelector('a');" +
					"if (link) {" +
					"var href = link.href || link.getAttribute('href') || '';" +
					"if (href && !/^#|^javascript:/.test(href)) return buildFullUrl(href);" +
					"var path = extractPathFromOnclick(link.getAttribute('onclick'));" +
					"if (path) return buildFullUrl(path);" +
					"}" +
					"var path = extractPathFromOnclick(element.getAttribute('onclick'));" +
					"if (path) return buildFullUrl(path);" +
					"return '';" +
					"}" +
					"var keywords = ['项目', '招标', '采购', '工程', '公告'];" +
					"var elements = document.querySelectorAll('tr, .item, .list-item, [class*=\"row\"], [class*=\"item\"]');" +
					"var results = [];" +
					"for (var i = 0; i < Math.min(elements.length, 30); i++) {" +
					"var text = (elements[i].innerText || elements[i].textContent || '').trim();" +
					"if (text.length > 20 && text.length < 500) {" +
					"var hasKeyword = false;" +
					"for (var j = 0; j < keywords.length; j++) {" +
					"if (text.indexOf(keywords[j]) >= 0) { hasKeyword = true; break; }" +
					"}" +
					"if (hasKeyword) {" +
					"results.push({" +
					"index: (results.length + 1).toString()," +
					"project: text.substring(0, 100).trim()," +
					"section: ''," +
					"region: ''," +
					"publishTime: ''," +
					"url: extractLink(elements[i])" +
					"});" +
					"}" +
					"}" +
					"}" +
					"return results;" +
					"})();";
				var backupResult = session.evaluate(backupCode);
				if (backupResult.success && backupResult.result) bidInfos = backupResult.result;
			}
			
			return { success: true, data: bidInfos || [], count: bidInfos ? bidInfos.length : 0 };
		} catch (error) {
			return { success: false, error: error.message || String(error), data: [], count: 0 };
		} finally {
			session.close();
		}
		})();
	`, url)

	logrus.Info("开始执行 JavaScript 代码...")
	logrus.WithField("code_length", len(jsCode)).Debug("JavaScript代码长度")
	result, err := sandbox.Run(jsCode)
	if err != nil {
		logrus.WithError(err).WithField("error_type", fmt.Sprintf("%T", err)).Error("执行代码失败")
		fmt.Printf("\n执行错误详情:\n")
		fmt.Printf("  错误类型: %T\n", err)
		fmt.Printf("  错误信息: %v\n", err)
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

	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		logrus.WithError(err).Fatal("获取工作目录失败")
	}

	// 查找 example 目录
	var dataDir string
	currentDir := wd

	// 向上查找 example 目录，最多向上查找5级
	for i := 0; i < 5; i++ {
		testPath := filepath.Join(currentDir, "example", "data")
		examplePath := filepath.Join(currentDir, "example")

		// 检查 example 目录是否存在
		if info, err := os.Stat(examplePath); err == nil && info.IsDir() {
			dataDir = testPath
			break
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
	if dataDir == "" {
		dataDir = filepath.Join(wd, "example", "data")
		logrus.WithField("path", dataDir).Warn("未找到项目根目录，使用当前目录")
	}

	// 确保 data 目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		logrus.WithError(err).WithField("path", dataDir).Fatal("创建data目录失败")
	}

	logrus.WithField("path", dataDir).Debug("数据目录路径")

	// 将数据转换为JSON格式
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		logrus.WithError(err).Fatal("JSON序列化失败")
	}

	// 保存到文件
	outputFile := filepath.Join(dataDir, fmt.Sprintf("bid_info_%s.json", time.Now().Format("20060102_150405")))
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

// getMapKeys 获取 map 的所有键，用于调试
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
