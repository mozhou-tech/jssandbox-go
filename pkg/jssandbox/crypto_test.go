package jssandbox

import (
	"context"
	"regexp"
	"testing"

	"github.com/dop251/goja"
)

func TestEncryptAES(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = encryptAES("Hello, World!", "my-secret-key");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("encryptAES() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("encryptAES()应该返回success: true")
	}

	data := resultObj.Get("data")
	if data == nil || goja.IsUndefined(data) {
		t.Error("encryptAES()缺少data字段")
	}

	encrypted := data.String()
	if len(encrypted) == 0 {
		t.Error("encryptAES()返回的加密数据为空")
	}
}

func TestDecryptAES(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 先加密
	encryptCode := `
		var encrypted = encryptAES("Hello, World!", "my-secret-key");
		encrypted.data;
	`

	encryptedResult, err := sb.Run(encryptCode)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	encryptedData := encryptedResult.String()

	// 再解密
	decryptCode := `
		var result = decryptAES("` + encryptedData + `", "my-secret-key");
		result;
	`

	result, err := sb.Run(decryptCode)
	if err != nil {
		t.Fatalf("decryptAES() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("decryptAES()应该返回success: true")
	}

	data := resultObj.Get("data")
	if data.String() != "Hello, World!" {
		t.Errorf("decryptAES()解密结果不正确, got %s, want Hello, World!", data.String())
	}
}

func TestEncryptDecryptAES_Integration(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	testCases := []string{
		"Hello, World!",
		"测试中文",
		"123456",
		"Special chars: !@#$%^&*()",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			code := `
				var encrypted = encryptAES("` + tc + `", "test-key");
				var decrypted = decryptAES(encrypted.data, "test-key");
				decrypted.data;
			`

			result, err := sb.Run(code)
			if err != nil {
				t.Fatalf("加密解密失败: %v", err)
			}

			if result.String() != tc {
				t.Errorf("加密解密结果不匹配, got %s, want %s", result.String(), tc)
			}
		})
	}
}

func TestHashSHA256(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = hashSHA256("Hello, World!");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("hashSHA256() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("hashSHA256()应该返回success: true")
	}

	hash := resultObj.Get("hash")
	if hash == nil || goja.IsUndefined(hash) {
		t.Error("hashSHA256()缺少hash字段")
	}

	hashStr := hash.String()
	// SHA256哈希是64个十六进制字符
	matched, _ := regexp.MatchString(`^[a-f0-9]{64}$`, hashStr)
	if !matched {
		t.Errorf("hashSHA256()返回的哈希格式不正确: %s", hashStr)
	}
}

func TestGenerateUUID(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		generateUUID();
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("generateUUID() error = %v", err)
	}

	uuid := result.String()
	// UUID格式: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, uuid)
	if !matched {
		t.Errorf("generateUUID()返回的UUID格式不正确: %s", uuid)
	}
}

func TestGenerateRandomString(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	t.Run("默认长度", func(t *testing.T) {
		code := `
			var result = generateRandomString();
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("generateRandomString() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		data := resultObj.Get("data")
		if len(data.String()) != 32 {
			t.Errorf("generateRandomString()默认长度应该是32, got %d", len(data.String()))
		}
	})

	t.Run("指定长度", func(t *testing.T) {
		code := `
			var result = generateRandomString(16);
			result;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("generateRandomString() error = %v", err)
		}

		resultObj := result.ToObject(sb.vm)
		data := resultObj.Get("data")
		if len(data.String()) != 16 {
			t.Errorf("generateRandomString(16)长度应该是16, got %d", len(data.String()))
		}
	})

	t.Run("多次生成应该不同", func(t *testing.T) {
		code := `
			var r1 = generateRandomString(32);
			var r2 = generateRandomString(32);
			r1.data !== r2.data;
		`

		result, err := sb.Run(code)
		if err != nil {
			t.Fatalf("generateRandomString() error = %v", err)
		}

		if !result.ToBoolean() {
			t.Error("多次生成的随机字符串应该不同")
		}
	})
}
