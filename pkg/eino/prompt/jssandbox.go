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
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/dop251/goja"
	"github.com/mozhou-tech/jssandbox-go/pkg/jssandbox"
)

type Config struct {
	// Sandbox is the JavaScript sandbox instance for executing prompt generation code
	// If nil, a new sandbox will be created with default config
	Sandbox *jssandbox.Sandbox
	// Code is the JavaScript code that generates prompt messages
	// The code should return an array of message objects
	// Each message object should have: { role: "user"|"assistant"|"system", content?: string, multiContent?: [...] }
	// Required
	Code string
	// Timeout specifies the execution timeout for the JavaScript code
	// Default: 30 seconds
	Timeout time.Duration
	// SandboxConfig is the config for creating a new sandbox if Sandbox is nil
	SandboxConfig *jssandbox.Config
}

func NewPromptTemplate(ctx context.Context, conf *Config) (prompt.ChatTemplate, error) {
	if conf == nil {
		return nil, fmt.Errorf("config is required")
	}
	if conf.Code == "" {
		return nil, fmt.Errorf("code is required")
	}

	sandbox := conf.Sandbox
	if sandbox == nil {
		sandboxConfig := conf.SandboxConfig
		if sandboxConfig == nil {
			sandboxConfig = jssandbox.DefaultConfig()
		}
		sandbox = jssandbox.NewSandboxWithConfig(ctx, sandboxConfig)
	}

	timeout := conf.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &chatTemplate{
		sandbox: sandbox,
		code:    conf.Code,
		timeout: timeout,
	}, nil
}

type chatTemplate struct {
	sandbox *jssandbox.Sandbox
	code    string
	timeout time.Duration
}

func (c *chatTemplate) Format(ctx context.Context, vs map[string]any, _ ...prompt.Option) ([]*schema.Message, error) {
	// 将变量注入到 JavaScript 运行时
	for k, v := range vs {
		// 直接设置值到 JavaScript 运行时
		c.sandbox.Set(k, v)
	}

	// 执行 JavaScript 代码
	// timeout 在初始化时已保证大于 0（默认为 30 秒），因此始终使用 RunWithTimeout
	result, err := c.sandbox.RunWithTimeout(c.code, c.timeout)
	if err != nil {
		return nil, fmt.Errorf("execute jssandbox code failed: %w", err)
	}

	// 将结果转换为消息数组
	messages, err := convertToMessages(result)
	if err != nil {
		return nil, fmt.Errorf("convert jssandbox result to messages failed: %w", err)
	}

	return messages, nil
}

// GetType returns the type of the chat template (JSSandbox).
func (c *chatTemplate) GetType() string {
	return "JSSandbox"
}

// convertToMessages 将 JavaScript 返回的结果转换为消息数组
func convertToMessages(result goja.Value) ([]*schema.Message, error) {
	if result == nil || goja.IsUndefined(result) || goja.IsNull(result) {
		return nil, fmt.Errorf("result is null or undefined")
	}

	// 导出为 Go 类型
	exported := result.Export()

	// 处理数组情况
	var messagesData []interface{}
	switch v := exported.(type) {
	case []interface{}:
		messagesData = v
	case []map[string]interface{}:
		// 转换为 []interface{}
		messagesData = make([]interface{}, len(v))
		for i, m := range v {
			messagesData[i] = m
		}
	default:
		// 尝试 JSON 序列化再反序列化
		jsonBytes, err := json.Marshal(exported)
		if err != nil {
			return nil, fmt.Errorf("marshal result failed: %w", err)
		}
		if err := json.Unmarshal(jsonBytes, &messagesData); err != nil {
			return nil, fmt.Errorf("unmarshal result to messages array failed: %w", err)
		}
	}

	messages := make([]*schema.Message, 0, len(messagesData))
	for i, msgData := range messagesData {
		msg, err := convertMessage(msgData)
		if err != nil {
			return nil, fmt.Errorf("convert message at index %d failed: %w", i, err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// convertMessage 将单个消息数据转换为 schema.Message
func convertMessage(msgData interface{}) (*schema.Message, error) {
	msgMap, ok := msgData.(map[string]interface{})
	if !ok {
		// 尝试 JSON 序列化再反序列化
		jsonBytes, err := json.Marshal(msgData)
		if err != nil {
			return nil, fmt.Errorf("marshal message data failed: %w", err)
		}
		if err := json.Unmarshal(jsonBytes, &msgMap); err != nil {
			return nil, fmt.Errorf("unmarshal message data failed: %w", err)
		}
	}

	msg := &schema.Message{}

	// 转换 role
	roleStr, ok := msgMap["role"].(string)
	if !ok {
		return nil, fmt.Errorf("message role is required and must be a string")
	}
	switch roleStr {
	case "user":
		msg.Role = schema.User
	case "assistant":
		msg.Role = schema.Assistant
	case "system":
		msg.Role = schema.System
	default:
		return nil, fmt.Errorf("unknown role: %s", roleStr)
	}

	// 转换 content
	if content, ok := msgMap["content"].(string); ok && content != "" {
		msg.Content = content
	}

	// 转换 multiContent
	if multiContentData, ok := msgMap["multiContent"].([]interface{}); ok {
		msg.MultiContent = make([]schema.ChatMessagePart, 0, len(multiContentData))
		for i, partData := range multiContentData {
			part, err := convertMessagePart(partData)
			if err != nil {
				return nil, fmt.Errorf("convert message part at index %d failed: %w", i, err)
			}
			msg.MultiContent = append(msg.MultiContent, *part)
		}
	}

	return msg, nil
}

// convertMessagePart 将消息部分数据转换为 schema.ChatMessagePart
func convertMessagePart(partData interface{}) (*schema.ChatMessagePart, error) {
	partMap, ok := partData.(map[string]interface{})
	if !ok {
		// 尝试 JSON 序列化再反序列化
		jsonBytes, err := json.Marshal(partData)
		if err != nil {
			return nil, fmt.Errorf("marshal part data failed: %w", err)
		}
		if err := json.Unmarshal(jsonBytes, &partMap); err != nil {
			return nil, fmt.Errorf("unmarshal part data failed: %w", err)
		}
	}

	part := &schema.ChatMessagePart{}

	// 转换 type
	typeStr, ok := partMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("part type is required and must be a string")
	}

	switch typeStr {
	case "image_url":
		part.Type = schema.ChatMessagePartTypeImageURL
		if imageURLData, ok := partMap["imageURL"].(map[string]interface{}); ok {
			part.ImageURL = &schema.ChatMessageImageURL{}
			if url, ok := imageURLData["url"].(string); ok {
				part.ImageURL.URL = url
			}
			if mimeType, ok := imageURLData["mimeType"].(string); ok {
				part.ImageURL.MIMEType = mimeType
			}
		}
	case "audio_url":
		part.Type = schema.ChatMessagePartTypeAudioURL
		if audioURLData, ok := partMap["audioURL"].(map[string]interface{}); ok {
			part.AudioURL = &schema.ChatMessageAudioURL{}
			if url, ok := audioURLData["url"].(string); ok {
				part.AudioURL.URL = url
			}
			if mimeType, ok := audioURLData["mimeType"].(string); ok {
				part.AudioURL.MIMEType = mimeType
			}
		}
	case "video_url":
		part.Type = schema.ChatMessagePartTypeVideoURL
		if videoURLData, ok := partMap["videoURL"].(map[string]interface{}); ok {
			part.VideoURL = &schema.ChatMessageVideoURL{}
			if url, ok := videoURLData["url"].(string); ok {
				part.VideoURL.URL = url
			}
			if mimeType, ok := videoURLData["mimeType"].(string); ok {
				part.VideoURL.MIMEType = mimeType
			}
		}
	case "text":
		part.Type = schema.ChatMessagePartTypeText
		if text, ok := partMap["text"].(string); ok {
			part.Text = text
		}
	default:
		return nil, fmt.Errorf("unknown part type: %s", typeStr)
	}

	return part, nil
}
