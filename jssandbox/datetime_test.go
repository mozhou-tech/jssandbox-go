package jssandbox

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"
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

func TestGetCurrentTimestamp(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	before := time.Now().Unix()

	code := `
		var result = getCurrentTimestamp();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("getCurrentTimestamp() error = %v", err)
	}

	after := time.Now().Unix()

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("getCurrentTimestamp()应该返回success: true")
	}

	timestamp := resultObj.Get("timestamp")
	ts := timestamp.ToInteger()
	if ts < before || ts > after {
		t.Errorf("getCurrentTimestamp()返回的时间戳不在合理范围内, got %d, want between %d and %d", ts, before, after)
	}

	timestampMs := resultObj.Get("timestampMs")
	tsMs := timestampMs.ToInteger()
	if tsMs < before*1000 || tsMs > after*1000+1000 {
		t.Errorf("getCurrentTimestamp()返回的毫秒时间戳不在合理范围内, got %d", tsMs)
	}

	date := resultObj.Get("date")
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, date.String())
	if !matched {
		t.Errorf("getCurrentTimestamp()返回的日期格式不正确: %s", date.String())
	}
}

func TestTimestampToDate(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 测试秒级时间戳 (2024-01-01 00:00:00 UTC)
	expectedDate := "2024-01-01 00:00:00"

	code := `
		var result = timestampToDate(1704067200);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("timestampToDate() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("timestampToDate()应该返回success: true")
	}

	date := resultObj.Get("date")
	if date.String() != expectedDate {
		t.Errorf("timestampToDate()返回的日期不正确, got %s, want %s", date.String(), expectedDate)
	}

	year := resultObj.Get("year")
	if year.ToInteger() != 2024 {
		t.Errorf("timestampToDate()年份不正确, got %d, want 2024", year.ToInteger())
	}

	month := resultObj.Get("month")
	if month.ToInteger() != 1 {
		t.Errorf("timestampToDate()月份不正确, got %d, want 1", month.ToInteger())
	}

	day := resultObj.Get("day")
	if day.ToInteger() != 1 {
		t.Errorf("timestampToDate()日期不正确, got %d, want 1", day.ToInteger())
	}
}

func TestTimestampToDate_WithFormat(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = timestampToDate(1704067200, "YYYY-MM-DD");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("timestampToDate() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	date := resultObj.Get("date")
	matched, _ := regexp.MatchString(`^2024-01-01$`, date.String())
	if !matched {
		t.Errorf("timestampToDate()格式化结果不正确, got %s, want 2024-01-01", date.String())
	}
}

func TestTimestampToDate_Milliseconds(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 测试毫秒级时间戳 (2024-01-01 00:00:00 UTC)
	code := `
		var result = timestampToDate(1704067200000);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("timestampToDate() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("timestampToDate()应该返回success: true")
	}

	date := resultObj.Get("date")
	matched, _ := regexp.MatchString(`^2024-01-01`, date.String())
	if !matched {
		t.Errorf("timestampToDate()毫秒时间戳转换不正确, got %s", date.String())
	}
}

func TestTimestampToDate_InvalidArgs(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = timestampToDate();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("timestampToDate() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	errMsg := resultObj.Get("error")
	if errMsg.String() == "" {
		t.Error("timestampToDate()缺少参数时应该返回错误")
	}
}

func TestFormatTimestamp(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = formatTimestamp(1704067200, "YYYY-MM-DD HH:mm:ss");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("formatTimestamp() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("formatTimestamp()应该返回success: true")
	}

	date := resultObj.Get("date")
	matched, _ := regexp.MatchString(`^2024-01-01 \d{2}:\d{2}:\d{2}$`, date.String())
	if !matched {
		t.Errorf("formatTimestamp()格式化结果不正确: %s", date.String())
	}
}

func TestFormatTimestamp_Milliseconds(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = formatTimestamp(1704067200000, "YYYY-MM-DD");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("formatTimestamp() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	date := resultObj.Get("date")
	matched, _ := regexp.MatchString(`^2024-01-01$`, date.String())
	if !matched {
		t.Errorf("formatTimestamp()毫秒时间戳格式化不正确, got %s", date.String())
	}
}

func TestFormatTimestamp_InvalidArgs(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []struct {
		name string
		code string
	}{
		{"缺少参数", `formatTimestamp();`},
		{"缺少格式参数", `formatTimestamp(1704067200);`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := `var result = ` + tc.code + ` result;`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("formatTimestamp() error = %v", err)
			}

			resultObj := result.ToObject(sb.vm)
			errMsg := resultObj.Get("error")
			if errMsg.String() == "" {
				t.Error("formatTimestamp()参数不足时应该返回错误")
			}
		})
	}
}

func TestDateToTimestamp(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []struct {
		name          string
		dateStr       string
		expectedYear  int64
		expectedMonth int64
	}{
		{"标准格式", "2024-01-01 12:00:00", 2024, 1},
		{"日期格式", "2024-01-01", 2024, 1},
		{"ISO格式", "2024-01-01T12:00:00Z", 2024, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := `
				var result = dateToTimestamp("` + tc.dateStr + `");
				result;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("dateToTimestamp() error = %v", err)
			}

			resultObj := result.ToObject(sb.vm)
			success := resultObj.Get("success")
			if !success.ToBoolean() {
				t.Error("dateToTimestamp()应该返回success: true")
			}

			timestamp := resultObj.Get("timestamp")
			if timestamp.ToInteger() <= 0 {
				t.Error("dateToTimestamp()返回的时间戳应该大于0")
			}

			timestampMs := resultObj.Get("timestampMs")
			if timestampMs.ToInteger() <= 0 {
				t.Error("dateToTimestamp()返回的毫秒时间戳应该大于0")
			}

			// 验证时间戳可以正确转换回日期
			ts := timestamp.ToInteger()
			verifyCode := fmt.Sprintf(`
				var converted = timestampToDate(%d);
				converted.year;
			`, ts)

			verifyResult, err := sb.Run(verifyCode)
			if err != nil {
				t.Fatalf("验证时间戳转换失败: %v", err)
			}

			convertedYear := verifyResult.ToInteger()
			if convertedYear != tc.expectedYear {
				t.Errorf("时间戳转换回日期年份不正确, got %d, want %d", convertedYear, tc.expectedYear)
			}
		})
	}
}

func TestDateToTimestamp_InvalidDate(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = dateToTimestamp("invalid-date");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("dateToTimestamp() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	errMsg := resultObj.Get("error")
	if errMsg.String() == "" {
		t.Error("dateToTimestamp()无效日期应该返回错误")
	}
}

func TestDateToTimestamp_InvalidArgs(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = dateToTimestamp();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("dateToTimestamp() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	errMsg := resultObj.Get("error")
	if errMsg.String() == "" {
		t.Error("dateToTimestamp()缺少参数时应该返回错误")
	}
}

func TestGetTimestampInfo(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = getTimestampInfo(1704067200);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("getTimestampInfo() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("getTimestampInfo()应该返回success: true")
	}

	timestamp := resultObj.Get("timestamp")
	if timestamp.ToInteger() != 1704067200 {
		t.Errorf("getTimestampInfo()时间戳不正确, got %d, want 1704067200", timestamp.ToInteger())
	}

	timestampMs := resultObj.Get("timestampMs")
	if timestampMs.ToInteger() != 1704067200000 {
		t.Errorf("getTimestampInfo()毫秒时间戳不正确, got %d, want 1704067200000", timestampMs.ToInteger())
	}

	date := resultObj.Get("date")
	matched, _ := regexp.MatchString(`^2024-01-01`, date.String())
	if !matched {
		t.Errorf("getTimestampInfo()日期不正确, got %s", date.String())
	}

	iso8601 := resultObj.Get("iso8601")
	if iso8601.String() == "" {
		t.Error("getTimestampInfo()ISO8601格式不应该为空")
	}

	year := resultObj.Get("year")
	if year.ToInteger() != 2024 {
		t.Errorf("getTimestampInfo()年份不正确, got %d, want 2024", year.ToInteger())
	}

	month := resultObj.Get("month")
	if month.ToInteger() != 1 {
		t.Errorf("getTimestampInfo()月份不正确, got %d, want 1", month.ToInteger())
	}

	day := resultObj.Get("day")
	if day.ToInteger() != 1 {
		t.Errorf("getTimestampInfo()日期不正确, got %d, want 1", day.ToInteger())
	}

	weekday := resultObj.Get("weekday")
	if weekday.String() == "" {
		t.Error("getTimestampInfo()星期不应该为空")
	}

	yearday := resultObj.Get("yearday")
	if yearday.ToInteger() != 1 {
		t.Errorf("getTimestampInfo()一年中的第几天不正确, got %d, want 1", yearday.ToInteger())
	}
}

func TestGetTimestampInfo_Milliseconds(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = getTimestampInfo(1704067200000);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("getTimestampInfo() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	year := resultObj.Get("year")
	if year.ToInteger() != 2024 {
		t.Errorf("getTimestampInfo()毫秒时间戳年份不正确, got %d, want 2024", year.ToInteger())
	}
}

func TestGetTimestampInfo_InvalidArgs(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = getTimestampInfo();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("getTimestampInfo() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	errMsg := resultObj.Get("error")
	if errMsg.String() == "" {
		t.Error("getTimestampInfo()缺少参数时应该返回错误")
	}
}

func TestConvertDateFormat(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"YYYY格式", "YYYY-MM-DD", "2006-01-02"},
		{"yyyy格式", "yyyy-MM-dd", "2006-01-02"},
		{"时间格式", "HH:mm:ss", "15:04:05"},
		{"12小时制", "hh:mm:ss A", "03:04:05 PM"},
		{"月份名称", "MMMM DD, YYYY", "January 02, 2006"},
		{"月份缩写", "MMM DD, YYYY", "Jan 02, 2006"},
		{"毫秒", "HH:mm:ss.SSS", "15:04:05.000"},
		{"混合格式", "YYYY-MM-DD HH:mm:ss", "2006-01-02 15:04:05"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convertDateFormat(tc.input)
			if result != tc.expected {
				t.Errorf("convertDateFormat(%s) = %s, want %s", tc.input, result, tc.expected)
			}
		})
	}
}

func TestDateTime_TimestampIntegration(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		// 获取当前时间戳
		var now = getCurrentTimestamp();
		
		// 时间戳转日期
		var dateFromTs = timestampToDate(now.timestamp, "YYYY-MM-DD HH:mm:ss");
		
		// 日期转时间戳
		var tsFromDate = dateToTimestamp(dateFromTs.date);
		
		// 格式化时间戳
		var formatted = formatTimestamp(now.timestamp, "YYYY-MM-DD");
		
		// 获取详细信息
		var info = getTimestampInfo(now.timestamp);
		
		{
			originalTimestamp: now.timestamp,
			dateFromTimestamp: dateFromTs.date,
			timestampFromDate: tsFromDate.timestamp,
			formatted: formatted.date,
			infoYear: info.year
		};
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("时间戳集成测试失败: %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	originalTimestamp := resultObj.Get("originalTimestamp")
	timestampFromDate := resultObj.Get("timestampFromDate")

	// 验证时间戳转换的准确性（允许1秒误差）
	diff := originalTimestamp.ToInteger() - timestampFromDate.ToInteger()
	if diff < 0 {
		diff = -diff
	}
	if diff > 1 {
		t.Errorf("时间戳转换不准确, 原始: %d, 转换后: %d, 差值: %d", 
			originalTimestamp.ToInteger(), timestampFromDate.ToInteger(), diff)
	}

	infoYear := resultObj.Get("infoYear")
	if infoYear.ToInteger() <= 0 {
		t.Error("集成测试年份应该大于0")
	}
}

