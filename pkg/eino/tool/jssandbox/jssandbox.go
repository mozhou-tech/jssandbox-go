/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jssandbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dop251/goja"
	"github.com/eino-contrib/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/mozhou-tech/jssandbox-go/pkg/jssandbox"
)

// JSSandboxToolDescription 工具描述
const JSSandboxToolDescription = `JavaScript沙盒执行工具，用于在安全的沙盒环境中执行JavaScript代码。

重要限制与说明：
1. **不支持 async/await**：沙盒环境不支持异步语法，代码必须同步执行。
2. **执行隔离与返回值**：代码在匿名函数中执行。**必须使用 return 语句返回结果**，否则将返回 undefined。
3. **错误处理**：大多数操作返回包含 error 字段的对象，建议始终检查 success 或 error 字段。

主要可用函数：
- 系统/环境：getCurrentDateTime(), getCPUNum(), getMemorySize(), getDiskSize(), sleep(ms), getEnv(name), readConfig(path)
- HTTP请求：httpGet(url), httpPost(url, body), httpRequest(url, options), fetch(url, options)
- 文件系统：readFile(path, options?), writeFile(path, content), appendFile(path, content), readFileHead(path, lines), getFileInfo(path), getFileHash(path, type), readImageBase64(path)
- 文档读取：readWord(path), readExcel(path), readPPT(path), readPDF(path)
- 浏览器自动化：createBrowserSession(timeout) -> navigate(url), wait(selector/sec), click(selector), fill(selector, value), evaluate(code), screenshot(path), getHTML(), getURL(), close()
- 图片处理：imageInfo(path), imageResize(in, out, w, h?), imageCrop(in, out, x, y, w, h), imageRotate(in, out, angle), imageConvert(in, out), imageQuality(in, out, q)
- 数据验证/处理：validateEmail(email), validateURL(url), validateIP(ip), validatePhone(phone), formatDate(date, fmt), parseDate(str), addDays(date, days)
- 编码/加密：encodeBase64(data), decodeBase64(str), encryptAES(data, key), decryptAES(enc, key), hashSHA256(data), generateUUID(), generateRandomString(len)
- 压缩/CSV：compressZip(files, out), extractZip(zip, dir), readCSV(path, opts), writeCSV(path, data), parseCSV(str)
- 网络/进程：resolveDNS(host), ping(host), checkPort(host, port), execCommand(cmd), listProcesses(), killProcess(pid)
- 文本/路径：textReplace(text, old, new), textSplit(text, sep), textJoin(parts, sep), textTrim(text), textContains(text, sub), pathJoin(...paths), pathAbs(path)
- 日志记录：logger.info/debug/warn/error/fatal(...args), logger.setLevel(level), logger.withFields(fields)

示例：
const resp = httpGet('https://api.example.com/data');
if (resp.error) return '错误: ' + resp.error;
const data = JSON.parse(resp.body);
writeFile('data.json', resp.body);
return { count: data.length, status: 'success' };`

// JSSandboxTool JavaScript沙盒工具
type JSSandboxTool struct {
	sandbox *jssandbox.Sandbox
	config  *JSSandboxConfig
	info    *schema.ToolInfo
}

// JSSandboxParams 工具参数
type JSSandboxParams struct {
	Code    string `json:"code"`              // 必需：要执行的JavaScript代码
	Timeout *int   `json:"timeout,omitempty"` // 可选：超时时间（秒），默认30秒
}

// JSSandboxConfig 工具配置
type JSSandboxConfig struct {
	SandboxConfig  *jssandbox.Config // jssandbox配置
	DefaultTimeout time.Duration     // 默认超时时间
}

// NewJSSandboxTool 创建新的JavaScript沙盒工具实例
func NewJSSandboxTool(ctx context.Context, cfg *JSSandboxConfig) (*JSSandboxTool, error) {
	if cfg == nil {
		cfg = &JSSandboxConfig{
			SandboxConfig:  jssandbox.DefaultConfig(),
			DefaultTimeout: 30 * time.Second,
		}
	}

	// 创建沙盒实例
	sandboxConfig := cfg.SandboxConfig
	if sandboxConfig == nil {
		sandboxConfig = jssandbox.DefaultConfig()
	}
	sandbox := jssandbox.NewSandboxWithConfig(ctx, sandboxConfig)

	return &JSSandboxTool{
		sandbox: sandbox,
		config:  cfg,
		info: &schema.ToolInfo{
			Name: "jssandbox",
			Desc: JSSandboxToolDescription,
			ParamsOneOf: schema.NewParamsOneOfByJSONSchema(
				&jsonschema.Schema{
					Type:     string(schema.Object),
					Required: []string{"code"},
					Properties: orderedmap.New[string, *jsonschema.Schema](
						orderedmap.WithInitialData[string, *jsonschema.Schema](
							orderedmap.Pair[string, *jsonschema.Schema]{
								Key: "code",
								Value: &jsonschema.Schema{
									Type:        string(schema.String),
									Description: "要执行的JavaScript代码。代码将在安全的沙盒环境中执行，支持系统操作、HTTP请求、文件系统操作、浏览器自动化等功能。",
								},
							},
							orderedmap.Pair[string, *jsonschema.Schema]{
								Key: "timeout",
								Value: &jsonschema.Schema{
									Type:        string(schema.Integer),
									Description: "可选参数。代码执行超时时间（秒），默认30秒。如果代码执行时间超过此值，将返回超时错误。",
								},
							},
						),
					),
				},
			),
		},
	}, nil
}

// Info 返回工具信息
func (t *JSSandboxTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return t.info, nil
}

// InvokableRun 执行工具
func (t *JSSandboxTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	param := &JSSandboxParams{}
	if err := json.Unmarshal([]byte(argumentsInJSON), param); err != nil {
		return "", fmt.Errorf("failed to extract input: %w", err)
	}
	return t.Execute(ctx, param)
}

// Execute 执行JavaScript代码
func (t *JSSandboxTool) Execute(ctx context.Context, params *JSSandboxParams) (string, error) {
	// 验证必需参数
	if params.Code == "" {
		return "", errors.New("parameter `code` is required")
	}

	// 确定超时时间
	var timeout time.Duration
	if params.Timeout != nil && *params.Timeout > 0 {
		timeout = time.Duration(*params.Timeout) * time.Second
	} else {
		timeout = t.config.DefaultTimeout
		if timeout == 0 {
			timeout = 30 * time.Second
		}
	}

	// 包装代码在匿名函数中，以支持 top-level return 并提供执行隔离
	wrappedCode := "(function(){\n" + params.Code + "\n})()"

	// 执行JavaScript代码
	var result goja.Value
	var err error
	if timeout > 0 {
		result, err = t.sandbox.RunWithTimeout(wrappedCode, timeout)
	} else {
		result, err = t.sandbox.Run(wrappedCode)
	}

	if err != nil {
		return "", fmt.Errorf("执行JavaScript代码失败: %w", err)
	}

	// 将结果转换为字符串
	resultStr := valueToString(result)
	return resultStr, nil
}

// valueToString 将goja.Value转换为字符串
func valueToString(v goja.Value) string {
	if v == nil {
		return "undefined"
	}

	// 检查是否为undefined或null
	if goja.IsUndefined(v) {
		return "undefined"
	}
	if goja.IsNull(v) {
		return "null"
	}

	// 尝试导出为Go类型并JSON序列化（适用于对象、数组等复杂类型）
	exported := v.Export()

	// 对于简单类型，直接使用String()方法
	switch exported.(type) {
	case string, int, int64, float64, bool:
		// 简单类型，尝试JSON序列化以获得更好的格式
		if jsonStr, err := json.Marshal(exported); err == nil {
			return string(jsonStr)
		}
		return v.String()
	default:
		// 复杂类型（对象、数组等），使用JSON序列化
		if jsonStr, err := json.MarshalIndent(exported, "", "  "); err == nil {
			return string(jsonStr)
		}
		// 如果JSON序列化失败，使用String()方法
		return v.String()
	}
}

// Close 关闭工具并清理资源
func (t *JSSandboxTool) Close() error {
	if t.sandbox != nil {
		return t.sandbox.Close()
	}
	return nil
}
