package main

import (
	"context"
	"fmt"

	"github.com/mozhou-tech/jssandbox-go/jssandbox"
	"github.com/sirupsen/logrus"
)

func main() {
	// 创建沙盒实例
	ctx := context.Background()
	sandbox := jssandbox.NewSandbox(ctx)
	defer sandbox.Close()

	// 示例1: 系统操作
	fmt.Println("=== 系统操作示例 ===")
	sandbox.Run(`
		console.log("当前时间:", getCurrentTime());
		console.log("当前日期:", getCurrentDate());
		console.log("CPU数量:", getCPUNum());
		var mem = getMemorySize();
		console.log("内存信息:", mem.totalStr, "/", mem.usedStr);
	`)

	// 示例2: HTTP请求
	fmt.Println("\n=== HTTP请求示例 ===")
	result, _ := sandbox.Run(`
		var response = httpGet("https://httpbin.org/get");
		console.log("状态码:", response.status);
		console.log("响应体:", response.body.substring(0, 100));
	`)
	fmt.Println("HTTP请求结果:", result)

	// 示例3: 文件操作
	fmt.Println("\n=== 文件操作示例 ===")
	sandbox.Run(`
		// 写入文件
		writeFile("test.txt", "Hello, World!\\n这是测试文件");
		console.log("文件已写入");
		
		// 读取文件
		var content = readFile("test.txt");
		console.log("文件内容:", content.data);
		
		// 获取文件信息
		var info = getFileInfo("test.txt");
		console.log("文件大小:", info.size, "字节");
		console.log("修改时间:", info.modTime);
	`)

	// 示例4: 浏览器操作
	fmt.Println("\n=== 浏览器操作示例 ===")
	sandbox.Run(`
		var result = browserNavigate("https://www.example.com");
		if (result.success) {
			console.log("页面加载成功，HTML长度:", result.html.length);
		}
	`)

	// 示例5: 作为大模型执行器
	fmt.Println("\n=== 大模型代码执行示例 ===")
	aiGeneratedCode := `
		// 这是大模型生成的代码
		var date = getCurrentDate();
		var time = getCurrentTime();
		var result = {
			date: date,
			time: time,
			message: "代码执行成功"
		};
		result;
	`

	result, err := sandbox.Run(aiGeneratedCode)
	if err != nil {
		logrus.WithError(err).Fatal("执行失败")
	}
	fmt.Println("执行结果:", result)
}
