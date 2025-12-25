package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mozhou-tech/jssandbox-go/pkg/jssandbox"
	"github.com/sirupsen/logrus"
)

// BidInfo 招标信息结构
type BidInfo struct {
	Index       string `json:"index"`       // 序号
	Project     string `json:"project"`     // 项目名称
	Section     string `json:"section"`     // 标段名称
	Region      string `json:"region"`      // 所在地区
	PublishTime string `json:"publishTime"` // 发布时间
}

func main() {
	ctx := context.Background()
	sandbox := jssandbox.NewSandbox(ctx)
	defer sandbox.Close()

	url := "http://jsggzy.jszwfw.gov.cn/jyxx/tradeInfonew.html"

	fmt.Println("开始爬取招标信息...")
	fmt.Printf("目标URL: %s\n\n", url)

	// 使用JavaScript代码进行爬取
	jsCode := `
		// 创建浏览器会话
		var session = createBrowserSession(60);
		
		try {
			// 导航到目标页面
			console.log("正在导航到页面...");
			var navResult = session.navigate("` + url + `");
			if (!navResult.success) {
				throw new Error("导航失败: " + navResult.error);
			}
			
			// 等待页面加载完成
			console.log("等待页面加载...");
			session.wait(3); // 等待3秒确保页面完全加载
			
			// 获取页面HTML
			console.log("获取页面内容...");
			var htmlResult = session.getHTML();
			if (!htmlResult.success) {
				throw new Error("获取HTML失败: " + htmlResult.error);
			}
			
			// 在页面中执行JavaScript提取表格数据
			var extractCode = "(function() {" +
				"var results = [];" +
				"var table = document.querySelector('table') || " +
				"document.querySelector('.table') || " +
				"document.querySelector('#dataTable') || " +
				"document.querySelector('tbody');" +
				"if (!table) {" +
				"var containers = document.querySelectorAll('[class*=\"table\"], [id*=\"table\"], [class*=\"list\"], [id*=\"list\"]');" +
				"for (var i = 0; i < containers.length; i++) {" +
				"if (containers[i].querySelectorAll('tr').length > 0) {" +
				"table = containers[i];" +
				"break;" +
				"}" +
				"}" +
				"}" +
				"if (table) {" +
				"var rows = table.querySelectorAll('tr');" +
				"for (var i = 1; i < rows.length; i++) {" +
				"var row = rows[i];" +
				"var cells = row.querySelectorAll('td, th');" +
				"if (cells.length >= 4) {" +
				"var info = {" +
				"index: cells[0] ? cells[0].innerText.trim() : ''," +
				"project: cells[1] ? cells[1].innerText.trim() : ''," +
				"section: cells[2] ? cells[2].innerText.trim() : ''," +
				"region: cells[3] ? cells[3].innerText.trim() : ''," +
				"publishTime: cells[4] ? cells[4].innerText.trim() : ''" +
				"};" +
				"if (info.project) {" +
				"results.push(info);" +
				"}" +
				"}" +
				"}" +
				"} else {" +
				"var links = document.querySelectorAll('a[href*=\"trade\"], a[href*=\"bid\"], a[href*=\"tender\"]');" +
				"links.forEach(function(link, index) {" +
				"var parent = link.closest('tr, li, div');" +
				"if (parent) {" +
				"var text = parent.innerText || parent.textContent || '';" +
				"if (text.length > 10) {" +
				"results.push({" +
				"index: (index + 1).toString()," +
				"project: link.innerText.trim() || text.substring(0, 50)," +
				"section: ''," +
				"region: ''," +
				"publishTime: ''" +
				"});" +
				"}" +
				"}" +
				"});" +
				"}" +
				"return results;" +
				"})();";
			
			var extractResult = session.evaluate(extractCode);
			if (!extractResult.success) {
				throw new Error("提取数据失败: " + extractResult.error);
			}
			
			var bidInfos = extractResult.result;
			
			// 如果提取的数据为空，尝试更通用的方法
			if (!bidInfos || bidInfos.length === 0) {
				console.log("尝试备用提取方法...");
				var backupCode = "(function() {" +
					"var results = [];" +
					"var elements = document.querySelectorAll('tr, .item, .list-item, [class*=\"row\"]');" +
					"for (var i = 0; i < Math.min(elements.length, 20); i++) {" +
					"var el = elements[i];" +
					"var text = el.innerText || el.textContent || '';" +
					"if (text.length > 20 && text.length < 500) {" +
					"if (text.indexOf('项目') >= 0 || text.indexOf('招标') >= 0 || " +
					"text.indexOf('采购') >= 0 || text.indexOf('工程') >= 0) {" +
					"results.push({" +
					"index: (results.length + 1).toString()," +
					"project: text.substring(0, 100).trim()," +
					"section: ''," +
					"region: ''," +
					"publishTime: ''" +
					"});" +
					"}" +
					"}" +
					"}" +
					"return results;" +
					"})();";
				var backupResult = session.evaluate(backupCode);
				if (backupResult.success && backupResult.result) {
					bidInfos = backupResult.result;
				}
			}
			
			// 返回结果
			{
				success: true,
				data: bidInfos || [],
				count: bidInfos ? bidInfos.length : 0
			}
		} catch (error) {
			{
				success: false,
				error: error.message || String(error),
				data: [],
				count: 0
			}
		} finally {
			session.close();
		}
	`

	result, err := sandbox.Run(jsCode)
	if err != nil {
		logrus.WithError(err).Fatal("执行爬虫代码失败")
	}

	// 解析结果
	var resultMap map[string]interface{}
	if resultObj := result.ToObject(nil); resultObj != nil {
		exported := resultObj.Export()
		if exportedMap, ok := exported.(map[string]interface{}); ok {
			resultMap = exportedMap
		} else {
			logrus.WithField("type", fmt.Sprintf("%T", exported)).Fatal("结果格式错误: 期望对象")
		}
	} else {
		logrus.Fatal("无法解析结果对象")
	}

	success, ok := resultMap["success"].(bool)
	if !ok || !success {
		if errMsg, ok := resultMap["error"].(string); ok {
			logrus.WithField("error", errMsg).Fatal("爬取失败")
		}
		logrus.Fatal("爬取失败: 未知错误")
	}

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
			}
		}
		fmt.Println("\n" + strings.Repeat("=", 80))
	}
}
