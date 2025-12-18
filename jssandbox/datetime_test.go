package jssandbox

import (
	"context"
	"regexp"
	"testing"
)

func TestFormatDate(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = formatDate("2024-01-01 12:00:00", "YYYY-MM-DD HH:mm:ss");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("formatDate() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("formatDate()应该返回success: true")
	}

	date := resultObj.Get("date")
	if date.String() == "" {
		t.Error("formatDate()返回的日期不应该为空")
	}
}

func TestParseDate(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []string{
		"2024-01-01",
		"2024-01-01 12:00:00",
		"2024-01-01T12:00:00Z",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			code := `
				var result = parseDate("` + tc + `");
				result;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("parseDate() error = %v", err)
			}

			resultObj := result.ToObject(sb.vm)
			success := resultObj.Get("success")
			if !success.ToBoolean() {
				t.Error("parseDate()应该返回success: true")
			}

			timestamp := resultObj.Get("timestamp")
			if timestamp.ToInteger() <= 0 {
				t.Error("parseDate()返回的时间戳应该大于0")
			}

			year := resultObj.Get("year")
			if year.ToInteger() != 2024 {
				t.Errorf("parseDate()年份不正确, got %d, want 2024", year.ToInteger())
			}
		})
	}
}

func TestAddDays(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = addDays("2024-01-01", 7);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("addDays() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("addDays()应该返回success: true")
	}

	date := resultObj.Get("date")
	// 验证日期格式
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, date.String())
	if !matched {
		t.Errorf("addDays()返回的日期格式不正确: %s", date.String())
	}
}

func TestAddDays_Negative(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = addDays("2024-01-08", -7);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("addDays() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	date := resultObj.Get("date")
	// 应该回到2024-01-01
	matched, _ := regexp.MatchString(`^2024-01-01`, date.String())
	if !matched {
		t.Errorf("addDays()减去7天后日期不正确: %s", date.String())
	}
}

func TestGetTimezone(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = getTimezone();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("getTimezone() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("getTimezone()应该返回success: true")
	}

	timezone := resultObj.Get("timezone")
	if timezone.String() == "" {
		t.Error("getTimezone()返回的时区不应该为空")
	}

	offset := resultObj.Get("offset")
	// 验证偏移量格式
	matched, _ := regexp.MatchString(`^[+-]\d{2}:\d{2}$`, offset.String())
	if !matched {
		t.Errorf("getTimezone()返回的偏移量格式不正确: %s", offset.String())
	}
}

func TestConvertTimezone(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = convertTimezone("2024-01-01 12:00:00", "UTC");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("convertTimezone() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("convertTimezone()应该返回success: true")
	}

	date := resultObj.Get("date")
	if date.String() == "" {
		t.Error("convertTimezone()返回的日期不应该为空")
	}
}

func TestDateTime_Integration(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var date = "2024-01-01";
		var parsed = parseDate(date);
		var added = addDays(date, 30);
		var formatted = formatDate(added.date, "YYYY-MM-DD");
		{
			original: date,
			parsedYear: parsed.year,
			addedDate: formatted.date
		};
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("日期时间集成测试失败: %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	parsedYear := resultObj.Get("parsedYear")
	if parsedYear.ToInteger() != 2024 {
		t.Errorf("集成测试年份不正确, got %d, want 2024", parsedYear.ToInteger())
	}
}

