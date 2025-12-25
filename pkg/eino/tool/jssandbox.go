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

package tool

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
const JSSandboxToolDescription = `JavaScript沙盒执行工具，用于在安全的沙盒环境中执行JavaScript代码
* 支持执行任意JavaScript代码，包括系统操作、HTTP请求、文件系统操作、浏览器自动化等功能
* 代码执行在隔离的沙盒环境中进行，具有超时保护机制
* 支持设置执行超时时间，防止代码长时间运行
* 返回代码执行结果，包括返回值和控制台输出

功能特性：
* 系统操作：获取时间、日期、CPU、内存、磁盘信息等
* HTTP请求：支持GET、POST等HTTP请求
* 文件系统：文件读写、文件信息获取、文件哈希等
* 浏览器自动化：页面导航、截图、点击、表单填写等
* 文档处理：Word、Excel、PPT、PDF文档读取
* 图片处理：图片格式转换、缩放、裁剪等
* 其他：加密解密、压缩解压、CSV处理、数据验证等`

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

	// 执行JavaScript代码
	var result goja.Value
	var err error
	if timeout > 0 {
		result, err = t.sandbox.RunWithTimeout(params.Code, timeout)
	} else {
		result, err = t.sandbox.Run(params.Code)
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
