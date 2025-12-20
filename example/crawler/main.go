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
	url := "http://jsggzy.jszwfw.gov.cn/jyxx/tradeInfonew.html"
	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("提取页面Title")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("目标URL: %s\n", url)
	fmt.Printf("会话超时: 120秒\n\n")
	logrus.WithField("url", url).Debug("准备访问URL")

	// 使用JavaScript代码提取页面title
	jsCode := fmt.Sprintf(`
		(function() {
		var session = createBrowserSession(120);
		try {
			console.log("[DEBUG] 开始导航...");
			var navResult = session.navigate("%s");
			console.log("[DEBUG] 导航结果:", JSON.stringify(navResult));
			if (!navResult.success) throw new Error("导航失败: " + navResult.error);
			console.log("[DEBUG] 导航成功，等待页面稳定...");
			session.wait(2);
			console.log("[DEBUG] 开始提取title...");
			
			var titleResult = session.evaluate("document.title");
			if (!titleResult.success) throw new Error("提取title失败: " + titleResult.error);
			
			var title = titleResult.result || "";
			console.log("[DEBUG] 提取到的title:", title);
			
			return { success: true, title: title };
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

	// 获取title
	title, ok := resultMap["title"].(string)
	if !ok {
		title = ""
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("页面Title: %s\n", title)
	fmt.Println(strings.Repeat("=", 80))
}

// getMapKeys 获取 map 的所有键，用于调试
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
