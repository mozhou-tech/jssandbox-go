package jssandbox

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/dop251/goja"
)

func TestGetCurrentTime(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	result, err := sb.Run("getCurrentTime()")
	if err != nil {
		t.Fatalf("getCurrentTime() error = %v", err)
	}

	timeStr := result.String()
	// 验证时间格式 HH:MM:SS
	matched, _ := regexp.MatchString(`^\d{2}:\d{2}:\d{2}$`, timeStr)
	if !matched {
		t.Errorf("getCurrentTime()返回格式不正确: %s", timeStr)
	}
}

func TestGetCurrentDate(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	result, err := sb.Run("getCurrentDate()")
	if err != nil {
		t.Fatalf("getCurrentDate() error = %v", err)
	}

	dateStr := result.String()
	// 验证日期格式 YYYY-MM-DD
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, dateStr)
	if !matched {
		t.Errorf("getCurrentDate()返回格式不正确: %s", dateStr)
	}

	// 验证日期是今天
	expectedDate := time.Now().Format("2006-01-02")
	if dateStr != expectedDate {
		t.Errorf("getCurrentDate()返回的日期不是今天, got %s, want %s", dateStr, expectedDate)
	}
}

func TestGetCurrentDateTime(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	result, err := sb.Run("getCurrentDateTime()")
	if err != nil {
		t.Fatalf("getCurrentDateTime() error = %v", err)
	}

	dateTimeStr := result.String()
	// 验证日期时间格式 YYYY-MM-DD HH:MM:SS
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, dateTimeStr)
	if !matched {
		t.Errorf("getCurrentDateTime()返回格式不正确: %s", dateTimeStr)
	}
}

func TestGetCPUNum(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	result, err := sb.Run("getCPUNum()")
	if err != nil {
		t.Fatalf("getCPUNum() error = %v", err)
	}

	cpuNum := result.ToInteger()
	if cpuNum <= 0 {
		t.Errorf("getCPUNum()返回无效值: %d", cpuNum)
	}
}

func TestGetMemorySize(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	result, err := sb.Run("getMemorySize()")
	if err != nil {
		t.Fatalf("getMemorySize() error = %v", err)
	}

	memObj := result.ToObject(sb.vm)
	if memObj == nil {
		t.Fatal("getMemorySize()返回的对象为nil")
	}

	// 检查必要的字段
	total := memObj.Get("total")
	if total == nil || goja.IsUndefined(total) {
		t.Error("getMemorySize()缺少total字段")
	}

	totalStr := memObj.Get("totalStr")
	if totalStr == nil || goja.IsUndefined(totalStr) {
		t.Error("getMemorySize()缺少totalStr字段")
	}

	// 验证total是数字
	if total.ToInteger() <= 0 {
		t.Errorf("getMemorySize()的total值无效: %d", total.ToInteger())
	}
}

func TestGetDiskSize(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	t.Run("默认路径", func(t *testing.T) {
		result, err := sb.Run("getDiskSize()")
		if err != nil {
			t.Fatalf("getDiskSize() error = %v", err)
		}

		diskObj := result.ToObject(sb.vm)
		if diskObj == nil {
			t.Fatal("getDiskSize()返回的对象为nil")
		}

		// 检查是否有错误字段
		errorVal := diskObj.Get("error")
		if errorVal != nil && !goja.IsUndefined(errorVal) {
			t.Logf("getDiskSize()返回错误（可能是权限问题）: %s", errorVal.String())
			return
		}

		// 验证字段存在
		total := diskObj.Get("total")
		if total == nil || goja.IsUndefined(total) {
			t.Error("getDiskSize()缺少total字段")
		}
	})

	t.Run("指定路径", func(t *testing.T) {
		result, err := sb.Run(`getDiskSize("/")`)
		if err != nil {
			t.Fatalf("getDiskSize() error = %v", err)
		}

		diskObj := result.ToObject(sb.vm)
		if diskObj == nil {
			t.Fatal("getDiskSize()返回的对象为nil")
		}

		// 检查是否有错误字段
		errorVal := diskObj.Get("error")
		if errorVal != nil && !goja.IsUndefined(errorVal) {
			t.Logf("getDiskSize()返回错误: %s", errorVal.String())
			return
		}

		freeStr := diskObj.Get("freeStr")
		if freeStr == nil || goja.IsUndefined(freeStr) {
			t.Error("getDiskSize()缺少freeStr字段")
		}
	})
}

func TestSleep(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	start := time.Now()
	_, err := sb.Run("sleep(100)")
	if err != nil {
		t.Fatalf("sleep() error = %v", err)
	}
	elapsed := time.Since(start)

	// 验证至少休眠了接近100ms（允许一些误差）
	if elapsed < 90*time.Millisecond {
		t.Errorf("sleep()休眠时间不足, got %v, want >= 90ms", elapsed)
	}
	if elapsed > 200*time.Millisecond {
		t.Errorf("sleep()休眠时间过长, got %v, want <= 200ms", elapsed)
	}
}

func TestSystemFunctionsIntegration(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var info = {
			time: getCurrentTime(),
			date: getCurrentDate(),
			datetime: getCurrentDateTime(),
			cpu: getCPUNum(),
			mem: getMemorySize(),
			disk: getDiskSize()
		};
		info;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("系统函数集成测试失败: %v", err)
	}

	infoObj := result.ToObject(sb.vm)
	if infoObj == nil {
		t.Fatal("集成测试返回的对象为nil")
	}

	// 验证所有字段都存在
	fields := []string{"time", "date", "datetime", "cpu", "mem", "disk"}
	for _, field := range fields {
		val := infoObj.Get(field)
		if val == nil || goja.IsUndefined(val) {
			t.Errorf("集成测试缺少字段: %s", field)
		}
	}
}

