package jssandbox

import (
	"context"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@example.co.uk", true},
		{"invalid.email", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.email, func(t *testing.T) {
			code := `
				var result = validateEmail("` + tc.email + `");
				result.valid;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("validateEmail() error = %v", err)
			}

			if result.ToBoolean() != tc.valid {
				t.Errorf("validateEmail(%s) = %v, want %v", tc.email, result.ToBoolean(), tc.valid)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []struct {
		url   string
		valid bool
	}{
		{"https://www.example.com", true},
		{"http://example.com/path", true},
		{"ftp://files.example.com", true},
		{"invalid-url", false},
		{"www.example.com", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			code := `
				var result = validateURL("` + tc.url + `");
				result.valid;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("validateURL() error = %v", err)
			}

			if result.ToBoolean() != tc.valid {
				t.Errorf("validateURL(%s) = %v, want %v", tc.url, result.ToBoolean(), tc.valid)
			}
		})
	}
}

func TestValidateIP(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []struct {
		ip     string
		valid  bool
		isIPv4 bool
		isIPv6 bool
	}{
		{"192.168.1.1", true, true, false},
		{"127.0.0.1", true, true, false},
		{"::1", true, false, true},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", true, false, true},
		{"invalid-ip", false, false, false},
		{"256.256.256.256", false, false, false},
		{"", false, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.ip, func(t *testing.T) {
			code := `
				var result = validateIP("` + tc.ip + `");
				({
					valid: result.valid,
					isIPv4: result.isIPv4,
					isIPv6: result.isIPv6
				});
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("validateIP() error = %v", err)
			}

			resultObj := result.ToObject(sb.vm)
			valid := resultObj.Get("valid").ToBoolean()
			isIPv4 := resultObj.Get("isIPv4").ToBoolean()
			isIPv6 := resultObj.Get("isIPv6").ToBoolean()

			if valid != tc.valid {
				t.Errorf("validateIP(%s).valid = %v, want %v", tc.ip, valid, tc.valid)
			}
			if isIPv4 != tc.isIPv4 {
				t.Errorf("validateIP(%s).isIPv4 = %v, want %v", tc.ip, isIPv4, tc.isIPv4)
			}
			if isIPv6 != tc.isIPv6 {
				t.Errorf("validateIP(%s).isIPv6 = %v, want %v", tc.ip, isIPv6, tc.isIPv6)
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []struct {
		phone string
		valid bool
	}{
		{"13800138000", true},
		{"15912345678", true},
		{"18888888888", true},
		{"12345678901", false}, // 不以1开头
		{"1380013800", false},  // 长度不足
		{"23800138000", false}, // 第二位数字不在3-9范围
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.phone, func(t *testing.T) {
			code := `
				var result = validatePhone("` + tc.phone + `");
				result.valid;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("validatePhone() error = %v", err)
			}

			if result.ToBoolean() != tc.valid {
				t.Errorf("validatePhone(%s) = %v, want %v", tc.phone, result.ToBoolean(), tc.valid)
			}
		})
	}
}

