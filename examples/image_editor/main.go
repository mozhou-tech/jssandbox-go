package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mozhou-tech/jssandbox-go/pkg/jssandbox"
	"github.com/sirupsen/logrus"
)

func main() {
	// 命令行参数
	var inputFile = flag.String("input", "", "输入图片文件路径（必需）")
	var outputDir = flag.String("output", "output", "输出目录（默认: output）")
	var operations = flag.String("ops", "all", "要执行的操作: all, resize, crop, rotate, flip, convert, quality, info（默认: all）")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("图片编辑示例")
		fmt.Println("用法: go run main.go -input <图片路径> [选项]")
		fmt.Println("\n选项:")
		flag.PrintDefaults()
		fmt.Println("\n示例:")
		fmt.Println("  go run main.go -input photo.jpg -ops all")
		fmt.Println("  go run main.go -input photo.jpg -ops resize,rotate,flip")
		fmt.Println("  go run main.go -input photo.jpg -ops info")
		os.Exit(1)
	}

	// 检查输入文件是否存在
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		logrus.WithField("file", *inputFile).Fatal("输入文件不存在")
	}

	// 创建输出目录
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		logrus.WithError(err).Fatal("创建输出目录失败")
	}

	// 获取输入文件信息
	inputBase := filepath.Base(*inputFile)
	inputExt := filepath.Ext(*inputFile)
	inputName := strings.TrimSuffix(inputBase, inputExt)

	ctx := context.Background()
	sandbox := jssandbox.NewSandbox(ctx)
	defer sandbox.Close()

	fmt.Printf("开始处理图片: %s\n", *inputFile)
	fmt.Printf("输出目录: %s\n\n", *outputDir)

	// 解析要执行的操作
	opsList := strings.Split(*operations, ",")
	opsMap := make(map[string]bool)
	for _, op := range opsList {
		opsMap[strings.TrimSpace(strings.ToLower(op))] = true
	}
	allOps := opsMap["all"]

	// 执行图片处理操作
	jsCode := fmt.Sprintf(`
		var inputFile = "%s";
		var outputDir = "%s";
		var inputName = "%s";
		var inputExt = "%s";
		
		console.log("=== 图片编辑示例 ===\\n");
		
		// 1. 获取图片信息
		console.log("1. 获取图片信息...");
		var info = imageInfo(inputFile);
		if (info.error) {
			console.error("获取图片信息失败:", info.error);
		} else {
			console.log("   原始尺寸:", info.width, "x", info.height);
			console.log("   格式:", info.format);
		}
		
		var results = [];
		
		// 2. 调整图片大小
		if (%t || %t) {
			console.log("\\n2. 调整图片大小...");
			// 缩放到宽度800，保持宽高比
			var resizeResult = imageResize(inputFile, outputDir + "/" + inputName + "_resized" + inputExt, 800);
			if (resizeResult.success) {
				console.log("   ✓ 调整大小成功:", resizeResult.path);
				results.push({operation: "resize", path: resizeResult.path});
			} else {
				console.error("   ✗ 调整大小失败:", resizeResult.error);
			}
			
			// 缩放到指定尺寸 600x400
			var resizeResult2 = imageResize(inputFile, outputDir + "/" + inputName + "_resized_600x400" + inputExt, 600, 400);
			if (resizeResult2.success) {
				console.log("   ✓ 调整大小(600x400)成功:", resizeResult2.path);
				results.push({operation: "resize_600x400", path: resizeResult2.path});
			}
		}
		
		// 3. 裁剪图片
		if (%t || %t) {
			console.log("\\n3. 裁剪图片...");
			// 从左上角裁剪 400x300 的区域
			var cropResult = imageCrop(inputFile, outputDir + "/" + inputName + "_cropped" + inputExt, 0, 0, 400, 300);
			if (cropResult.success) {
				console.log("   ✓ 裁剪成功:", cropResult.path);
				results.push({operation: "crop", path: cropResult.path});
			} else {
				console.error("   ✗ 裁剪失败:", cropResult.error);
			}
		}
		
		// 4. 旋转图片
		if (%t || %t) {
			console.log("\\n4. 旋转图片...");
			// 旋转90度
			var rotateResult = imageRotate(inputFile, outputDir + "/" + inputName + "_rotated_90" + inputExt, 90);
			if (rotateResult.success) {
				console.log("   ✓ 旋转90度成功:", rotateResult.path);
				results.push({operation: "rotate_90", path: rotateResult.path});
			}
			
			// 旋转180度
			var rotateResult2 = imageRotate(inputFile, outputDir + "/" + inputName + "_rotated_180" + inputExt, 180);
			if (rotateResult2.success) {
				console.log("   ✓ 旋转180度成功:", rotateResult2.path);
				results.push({operation: "rotate_180", path: rotateResult2.path});
			}
		}
		
		// 5. 翻转图片
		if (%t || %t) {
			console.log("\\n5. 翻转图片...");
			// 水平翻转
			var flipHResult = imageFlip(inputFile, outputDir + "/" + inputName + "_flipped_h" + inputExt, "horizontal");
			if (flipHResult.success) {
				console.log("   ✓ 水平翻转成功:", flipHResult.path);
				results.push({operation: "flip_horizontal", path: flipHResult.path});
			}
			
			// 垂直翻转
			var flipVResult = imageFlip(inputFile, outputDir + "/" + inputName + "_flipped_v" + inputExt, "vertical");
			if (flipVResult.success) {
				console.log("   ✓ 垂直翻转成功:", flipVResult.path);
				results.push({operation: "flip_vertical", path: flipVResult.path});
			}
		}
		
		// 6. 转换图片格式
		if (%t || %t) {
			console.log("\\n6. 转换图片格式...");
			// 转换为PNG
			var convertResult = imageConvert(inputFile, outputDir + "/" + inputName + "_converted.png");
			if (convertResult.success) {
				console.log("   ✓ 转换为PNG成功:", convertResult.path);
				results.push({operation: "convert_png", path: convertResult.path});
			}
			
			// 转换为JPEG
			var convertResult2 = imageConvert(inputFile, outputDir + "/" + inputName + "_converted.jpg");
			if (convertResult2.success) {
				console.log("   ✓ 转换为JPEG成功:", convertResult2.path);
				results.push({operation: "convert_jpg", path: convertResult2.path});
			}
		}
		
		// 7. 调整JPEG质量
		if (%t || %t) {
			console.log("\\n7. 调整JPEG质量...");
			// 如果输入是JPEG，调整质量
			if (inputExt.toLowerCase() === ".jpg" || inputExt.toLowerCase() === ".jpeg") {
				// 高质量 (95)
				var qualityResult1 = imageQuality(inputFile, outputDir + "/" + inputName + "_quality_95.jpg", 95);
				if (qualityResult1.success) {
					console.log("   ✓ 高质量(95)保存成功:", qualityResult1.path);
					results.push({operation: "quality_95", path: qualityResult1.path});
				}
				
				// 中等质量 (75)
				var qualityResult2 = imageQuality(inputFile, outputDir + "/" + inputName + "_quality_75.jpg", 75);
				if (qualityResult2.success) {
					console.log("   ✓ 中等质量(75)保存成功:", qualityResult2.path);
					results.push({operation: "quality_75", path: qualityResult2.path});
				}
				
				// 低质量 (50)
				var qualityResult3 = imageQuality(inputFile, outputDir + "/" + inputName + "_quality_50.jpg", 50);
				if (qualityResult3.success) {
					console.log("   ✓ 低质量(50)保存成功:", qualityResult3.path);
					results.push({operation: "quality_50", path: qualityResult3.path});
				}
			} else {
				console.log("   ⚠ 输入文件不是JPEG格式，跳过质量调整");
			}
		}
		
		// 8. 组合操作示例：调整大小 + 旋转 + 翻转
		if (%t || %t) {
			console.log("\\n8. 组合操作示例...");
			// 先调整大小
			var step1 = imageResize(inputFile, outputDir + "/" + inputName + "_combo_step1" + inputExt, 500);
			if (step1.success) {
				// 再旋转
				var step2 = imageRotate(step1.path, outputDir + "/" + inputName + "_combo_step2" + inputExt, 45);
				if (step2.success) {
					// 最后翻转
					var step3 = imageFlip(step2.path, outputDir + "/" + inputName + "_combo_final" + inputExt, "horizontal");
					if (step3.success) {
						console.log("   ✓ 组合操作成功:", step3.path);
						results.push({operation: "combo", path: step3.path});
					}
				}
			}
		}
		
		console.log("\\n=== 处理完成 ===");
		console.log("共生成", results.length, "个文件");
		
		// 返回结果
		{
			success: true,
			info: info,
			results: results,
			count: results.length
		}
	`,
		*inputFile,
		*outputDir,
		inputName,
		inputExt,
		allOps || opsMap["resize"],
		allOps || opsMap["resize"],
		allOps || opsMap["crop"],
		allOps || opsMap["crop"],
		allOps || opsMap["rotate"],
		allOps || opsMap["rotate"],
		allOps || opsMap["flip"],
		allOps || opsMap["flip"],
		allOps || opsMap["convert"],
		allOps || opsMap["convert"],
		allOps || opsMap["quality"],
		allOps || opsMap["quality"],
		allOps,
		allOps,
	)

	result, err := sandbox.Run(jsCode)
	if err != nil {
		logrus.WithError(err).Fatal("执行图片处理代码失败")
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
		logrus.Fatal("图片处理失败")
	}

	// 显示图片信息
	if info, ok := resultMap["info"].(map[string]interface{}); ok {
		if info["error"] == nil {
			fmt.Printf("\n图片信息:\n")
			if width, ok := info["width"].(float64); ok {
				fmt.Printf("  宽度: %.0f 像素\n", width)
			}
			if height, ok := info["height"].(float64); ok {
				fmt.Printf("  高度: %.0f 像素\n", height)
			}
			if format, ok := info["format"].(string); ok {
				fmt.Printf("  格式: %s\n", format)
			}
		}
	}

	// 显示处理结果
	if results, ok := resultMap["results"].([]interface{}); ok {
		fmt.Printf("\n处理结果:\n")
		for i, item := range results {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if op, ok := itemMap["operation"].(string); ok {
					if path, ok := itemMap["path"].(string); ok {
						fmt.Printf("  [%d] %s: %s\n", i+1, op, path)
					}
				}
			}
		}
	}

	if count, ok := resultMap["count"].(float64); ok {
		fmt.Printf("\n总共生成了 %.0f 个文件\n", count)
	}
}
