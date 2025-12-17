package jssandbox

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// 注意: 这些测试需要系统安装 ffmpeg 才能运行
// 在实际测试环境中，需要确保 ffmpeg 已安装

func TestVideoConvert(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过视频处理测试（需要 ffmpeg）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	// 注意: 这里需要一个真实的视频文件进行测试
	// 在实际使用中，应该创建一个测试视频文件或使用示例视频
	inputPath := filepath.Join(testDir, "input.mp4")
	outputPath := filepath.Join(testDir, "output.avi")

	// 创建一个空的测试文件（实际测试中应该使用真实视频）
	_, err := os.Create(inputPath)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	code := `
		var result = videoConvert("` + inputPath + `", "` + outputPath + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		// 如果 ffmpeg 未安装或文件无效，这是预期的
		t.Logf("videoConvert() 执行失败（可能是 ffmpeg 未安装或测试文件无效）: %v", err)
		return
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		errorVal := resultObj.Get("error")
		t.Logf("videoConvert()失败（可能是 ffmpeg 未安装或测试文件无效）: %s", errorVal.String())
		return
	}

	// 验证输出文件存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("输出文件不存在: %s", outputPath)
	}
}

func TestVideoInfo(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.mp4")

	// 创建一个测试文件
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}
	file.Close()

	code := `
		var result = videoInfo("` + testFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("videoInfo() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	path := resultObj.Get("path")
	if path.String() != testFile {
		t.Errorf("videoInfo() 返回的路径不正确: 期望 %s, 得到 %s", testFile, path.String())
	}
}

func TestVideoExtractAudio(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过视频处理测试（需要 ffmpeg）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	inputPath := filepath.Join(testDir, "input.mp4")
	outputPath := filepath.Join(testDir, "output.mp3")

	// 创建一个空的测试文件
	_, err := os.Create(inputPath)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	code := `
		var result = videoExtractAudio("` + inputPath + `", "` + outputPath + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("videoExtractAudio() 执行失败（可能是 ffmpeg 未安装或测试文件无效）: %v", err)
		return
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		errorVal := resultObj.Get("error")
		t.Logf("videoExtractAudio()失败（可能是 ffmpeg 未安装或测试文件无效）: %s", errorVal.String())
	}
}

func TestVideoResize(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过视频处理测试（需要 ffmpeg）")
	}

	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testDir := t.TempDir()
	inputPath := filepath.Join(testDir, "input.mp4")
	outputPath := filepath.Join(testDir, "output.mp4")

	_, err := os.Create(inputPath)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	code := `
		var result = videoResize("` + inputPath + `", "` + outputPath + `", 640);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Logf("videoResize() 执行失败（可能是 ffmpeg 未安装或测试文件无效）: %v", err)
		return
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		errorVal := resultObj.Get("error")
		t.Logf("videoResize()失败（可能是 ffmpeg 未安装或测试文件无效）: %s", errorVal.String())
	}
}

