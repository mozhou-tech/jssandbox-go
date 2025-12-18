package jssandbox

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dop251/goja"
)

func TestGetEnv(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 设置测试环境变量
	testKey := "TEST_ENV_VAR"
	testValue := "test_value"
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	code := `
		var result = getEnv("` + testKey + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("getEnv() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("getEnv()应该返回success: true")
	}

	value := resultObj.Get("value")
	if value.String() != testValue {
		t.Errorf("getEnv()返回值不正确, got %s, want %s", value.String(), testValue)
	}

	exists := resultObj.Get("exists")
	if !exists.ToBoolean() {
		t.Error("getEnv()应该返回exists: true")
	}
}

func TestGetEnv_NotExists(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = getEnv("NON_EXISTENT_ENV_VAR_12345");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("getEnv() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	exists := resultObj.Get("exists")
	if exists.ToBoolean() {
		t.Error("不存在的环境变量应该返回exists: false")
	}
}

func TestGetEnvAll(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = getEnvAll();
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("getEnvAll() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("getEnvAll()应该返回success: true")
	}

	env := resultObj.Get("env")
	if env == nil || goja.IsUndefined(env) {
		t.Error("getEnvAll()缺少env字段")
	}

	// 验证至少有一些环境变量
	envObj := env.ToObject(sb.vm)
	keys := envObj.Keys()
	if len(keys) == 0 {
		t.Error("getEnvAll()应该返回至少一个环境变量")
	}
}

func TestReadConfig_JSON(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建测试JSON配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")
	configContent := `{"name": "test", "value": 123, "enabled": true}`
	os.WriteFile(configFile, []byte(configContent), 0644)

	code := `
		var result = readConfig("` + configFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readConfig() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("readConfig()应该返回success: true")
	}

	config := resultObj.Get("config")
	if config == nil || goja.IsUndefined(config) {
		t.Error("readConfig()缺少config字段")
	}

	configObj := config.ToObject(sb.vm)
	name := configObj.Get("name")
	if name.String() != "test" {
		t.Errorf("readConfig()配置值不正确, got %s, want test", name.String())
	}
}

func TestReadConfig_YAML(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 创建测试YAML配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")
	configContent := `name: test
value: 123
enabled: true`
	os.WriteFile(configFile, []byte(configContent), 0644)

	code := `
		var result = readConfig("` + configFile + `");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("readConfig() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("readConfig()应该返回success: true")
	}

	config := resultObj.Get("config")
	if config == nil || goja.IsUndefined(config) {
		t.Error("readConfig()缺少config字段")
	}

	configObj := config.ToObject(sb.vm)
	name := configObj.Get("name")
	if name.String() != "test" {
		t.Errorf("readConfig()配置值不正确, got %s, want test", name.String())
	}
}

