package jssandbox

import (
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"

	"github.com/dop251/goja"
	"github.com/google/uuid"
)

// registerCrypto 注册加密/解密功能到JavaScript运行时
func (sb *Sandbox) registerCrypto() {
	// AES加密
	sb.vm.Set("encryptAES", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供数据和密钥参数",
			})
		}

		data := call.Arguments[0].String()
		key := call.Arguments[1].String()

		// 将密钥转换为32字节（AES-256）
		keyHash := sha256.Sum256([]byte(key))
		block, err := aes.NewCipher(keyHash[:])
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("创建加密器失败: %v", err),
			})
		}

		// 使用GCM模式
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("创建GCM失败: %v", err),
			})
		}

		// 生成随机nonce
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(cryptorand.Reader, nonce); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("生成nonce失败: %v", err),
			})
		}

		// 加密
		ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
		encrypted := base64.StdEncoding.EncodeToString(ciphertext)

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"data":    encrypted,
		})
	})

	// AES解密
	sb.vm.Set("decryptAES", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供加密数据和密钥参数",
			})
		}

		encrypted := call.Arguments[0].String()
		key := call.Arguments[1].String()

		// 解码base64
		ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("解码base64失败: %v", err),
			})
		}

		// 将密钥转换为32字节（AES-256）
		keyHash := sha256.Sum256([]byte(key))
		block, err := aes.NewCipher(keyHash[:])
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("创建解密器失败: %v", err),
			})
		}

		// 使用GCM模式
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("创建GCM失败: %v", err),
			})
		}

		// 提取nonce
		nonceSize := gcm.NonceSize()
		if len(ciphertext) < nonceSize {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "密文长度不足",
			})
		}

		nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

		// 解密
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("解密失败: %v", err),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"data":    string(plaintext),
		})
	})

	// SHA256哈希
	sb.vm.Set("hashSHA256", func(data string) goja.Value {
		hash := sha256.Sum256([]byte(data))
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"hash":    hex.EncodeToString(hash[:]),
		})
	})

	// 生成UUID
	sb.vm.Set("generateUUID", func() string {
		return uuid.New().String()
	})

	// 生成随机字符串
	sb.vm.Set("generateRandomString", func(call goja.FunctionCall) goja.Value {
		length := 32
		if len(call.Arguments) > 0 {
			length = int(call.Arguments[0].ToInteger())
			if length <= 0 {
				length = 32
			}
			if length > 1024 {
				length = 1024
			}
		}

		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"data":    string(b),
		})
	})
}
