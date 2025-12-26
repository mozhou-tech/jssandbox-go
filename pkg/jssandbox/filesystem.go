package jssandbox

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/djherbis/times"
	"github.com/dop251/goja"
	"github.com/h2non/filetype"
)

// registerFileSystem 注册文件系统操作功能到JavaScript运行时
func (sb *Sandbox) registerFileSystem() {
	// 使用操作系统默认软件打开文件
	sb.vm.Set("openFile", func(filePath string) map[string]interface{} {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", filePath)
		case "linux":
			cmd = exec.Command("xdg-open", filePath)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", filePath)
		default:
			return map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("不支持的操作系统: %s", runtime.GOOS),
			}
		}

		err := cmd.Run()
		if err != nil {
			sb.logger.WithError(err).WithField("path", filePath).Error("打开文件失败")
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}

		return map[string]interface{}{
			"success": true,
		}
	})

	// 读取文件元信息
	sb.vm.Set("getFileInfo", func(filePath string) goja.Value {
		info, err := os.Stat(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		result := map[string]interface{}{
			"name":  info.Name(),
			"size":  info.Size(),
			"mode":  info.Mode().String(),
			"isDir": info.IsDir(),
		}

		// 获取时间信息
		t, err := times.Stat(filePath)
		if err == nil {
			if t.HasBirthTime() {
				result["birthTime"] = t.BirthTime().Format("2006-01-02 15:04:05")
			}
			// ModTime 和 AccessTime 总是可用的，直接使用
			modTime := t.ModTime()
			if !modTime.IsZero() {
				result["modTime"] = modTime.Format("2006-01-02 15:04:05")
			}
			accessTime := t.AccessTime()
			if !accessTime.IsZero() {
				result["accessTime"] = accessTime.Format("2006-01-02 15:04:05")
			}
		} else {
			result["modTime"] = info.ModTime().Format("2006-01-02 15:04:05")
		}

		// 获取文件类型（优先使用filetype库检测，失败则使用扩展名）
		ext := filepath.Ext(filePath)
		result["extension"] = ext

		// 尝试使用filetype库检测
		file, err := os.Open(filePath)
		if err == nil {
			buf := make([]byte, 261)
			n, _ := file.Read(buf)
			file.Close()

			if kind, err := filetype.Match(buf[:n]); err == nil && kind != filetype.Unknown {
				result["type"] = kind.MIME.Value
				result["mime"] = kind.MIME.Value
				result["mimeType"] = kind.MIME.Type
				result["mimeSubtype"] = kind.MIME.Subtype
			} else {
				// 回退到基于扩展名的判断
				result["type"] = getFileType(ext)
			}
		} else {
			// 回退到基于扩展名的判断
			result["type"] = getFileType(ext)
		}

		return sb.vm.ToValue(result)
	})

	// 重命名文件
	sb.vm.Set("renameFile", func(oldPath, newPath string) map[string]interface{} {
		err := os.Rename(oldPath, newPath)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})

	// 读取文件内容（支持分页）
	sb.vm.Set("readFile", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文件路径",
			})
		}

		filePath := call.Arguments[0].String()
		offset := int64(0)
		limit := int64(1024 * 1024) // 默认1MB

		if len(call.Arguments) > 1 {
			options := call.Arguments[1].ToObject(sb.vm)
			if offsetVal := options.Get("offset"); offsetVal != nil && !goja.IsUndefined(offsetVal) {
				offset = offsetVal.ToInteger()
			}
			if limitVal := options.Get("limit"); limitVal != nil && !goja.IsUndefined(limitVal) {
				limit = limitVal.ToInteger()
			}
		}

		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer file.Close()

		if offset > 0 {
			file.Seek(offset, 0)
		}

		buffer := make([]byte, limit)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"data":   string(buffer[:n]),
			"length": n,
			"offset": offset,
		})
	})

	// 读取文本文件的前几行
	sb.vm.Set("readFileHead", func(filePath string, lines int) goja.Value {
		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var result []string
		count := 0

		for scanner.Scan() && count < lines {
			result = append(result, scanner.Text())
			count++
		}

		if err := scanner.Err(); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"lines": result,
			"count": count,
		})
	})

	// 读取文本文件的后几行
	sb.vm.Set("readFileTail", func(filePath string, lines int) goja.Value {
		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var allLines []string

		for scanner.Scan() {
			allLines = append(allLines, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		start := len(allLines) - lines
		if start < 0 {
			start = 0
		}

		result := allLines[start:]
		return sb.vm.ToValue(map[string]interface{}{
			"lines": result,
			"count": len(result),
		})
	})

	// 读取文件的哈希值
	sb.vm.Set("getFileHash", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文件路径",
			})
		}

		filePath := call.Arguments[0].String()
		hashType := "md5"
		if len(call.Arguments) > 1 {
			hashType = strings.ToLower(call.Arguments[1].String())
		}

		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer file.Close()

		var hash string
		switch hashType {
		case "md5":
			h := md5.New()
			io.Copy(h, file)
			hash = fmt.Sprintf("%x", h.Sum(nil))
		case "sha1":
			h := sha1.New()
			io.Copy(h, file)
			hash = fmt.Sprintf("%x", h.Sum(nil))
		case "sha256":
			h := sha256.New()
			io.Copy(h, file)
			hash = fmt.Sprintf("%x", h.Sum(nil))
		case "sha512":
			h := sha512.New()
			io.Copy(h, file)
			hash = fmt.Sprintf("%x", h.Sum(nil))
		default:
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("不支持的哈希类型: %s", hashType),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"hash": hash,
			"type": hashType,
		})
	})

	// 读取图片文件的base64编码
	sb.vm.Set("readImageBase64", func(filePath string) goja.Value {
		file, err := os.Open(filePath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": err.Error(),
			})
		}

		ext := strings.ToLower(filepath.Ext(filePath))
		mimeType := "image/png"
		switch ext {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".gif":
			mimeType = "image/gif"
		case ".webp":
			mimeType = "image/webp"
		case ".bmp":
			mimeType = "image/bmp"
		}

		base64Str := base64.StdEncoding.EncodeToString(data)
		return sb.vm.ToValue(map[string]interface{}{
			"base64":   base64Str,
			"mimeType": mimeType,
			"dataUrl":  fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str),
		})
	})

	// 写入文件
	sb.vm.Set("writeFile", func(filePath string, content string) map[string]interface{} {
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})

	// 追加文件
	sb.vm.Set("appendFile", func(filePath string, content string) map[string]interface{} {
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		defer file.Close()

		_, err = file.WriteString(content)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}

		return map[string]interface{}{
			"success": true,
		}
	})

	// 创建临时文件
	sb.vm.Set("createTempFile", func(call goja.FunctionCall) goja.Value {
		dir := ""
		pattern := ""

		if len(call.Arguments) > 0 {
			options := call.Arguments[0].ToObject(sb.vm)
			if dirVal := options.Get("dir"); dirVal != nil && !goja.IsUndefined(dirVal) {
				dir = dirVal.String()
			}
			if patternVal := options.Get("pattern"); patternVal != nil && !goja.IsUndefined(patternVal) {
				pattern = patternVal.String()
			}
		}

		// 如果没有指定目录，使用系统临时目录
		if dir == "" {
			dir = os.TempDir()
		}

		// 如果没有指定模式，使用默认模式
		if pattern == "" {
			pattern = "temp-*.tmp"
		}

		file, err := os.CreateTemp(dir, pattern)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}
		defer file.Close()

		filePath := file.Name()
		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    filePath,
		})
	})

	// 获取当前工作目录
	getCurrentDir := func() map[string]interface{} {
		dir, err := os.Getwd()
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
			"path":    dir,
		}
	}
	sb.vm.Set("getCurrentDir", getCurrentDir)
	sb.vm.Set("pwd", getCurrentDir)

	// 创建目录
	makeDir := func(call goja.FunctionCall) map[string]interface{} {
		if len(call.Arguments) < 1 {
			return map[string]interface{}{
				"success": false,
				"error":   "需要提供目录路径",
			}
		}
		dirPath := call.Arguments[0].String()
		recursive := false
		if len(call.Arguments) > 1 {
			recursive = call.Arguments[1].ToBoolean()
		}

		var err error
		if recursive {
			err = os.MkdirAll(dirPath, 0755)
		} else {
			err = os.Mkdir(dirPath, 0755)
		}

		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	}
	sb.vm.Set("makeDir", makeDir)
	sb.vm.Set("mkdir", makeDir)

	// 列出目录内容
	listDir := func(dirPath string) map[string]interface{} {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}

		var result []map[string]interface{}
		for _, entry := range entries {
			info, _ := entry.Info()
			item := map[string]interface{}{
				"name":  entry.Name(),
				"isDir": entry.IsDir(),
			}
			if info != nil {
				item["size"] = info.Size()
				item["modTime"] = info.ModTime().Format("2006-01-02 15:04:05")
			}
			result = append(result, item)
		}

		return map[string]interface{}{
			"success": true,
			"entries": result,
		}
	}
	sb.vm.Set("listDir", listDir)
	sb.vm.Set("ls", listDir)

	// 检查路径是否存在
	sb.vm.Set("pathExists", func(path string) bool {
		_, err := os.Stat(path)
		return err == nil || os.IsExist(err)
	})

	// 删除目录
	sb.vm.Set("removeDir", func(call goja.FunctionCall) map[string]interface{} {
		if len(call.Arguments) < 1 {
			return map[string]interface{}{
				"success": false,
				"error":   "需要提供目录路径",
			}
		}
		dirPath := call.Arguments[0].String()
		recursive := false
		if len(call.Arguments) > 1 {
			recursive = call.Arguments[1].ToBoolean()
		}

		var err error
		if recursive {
			err = os.RemoveAll(dirPath)
		} else {
			err = os.Remove(dirPath)
		}

		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})

	// 删除文件
	sb.vm.Set("deleteFile", func(filePath string) map[string]interface{} {
		err := os.Remove(filePath)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		return map[string]interface{}{
			"success": true,
		}
	})
}

// getFileType 根据扩展名返回文件类型
func getFileType(ext string) string {
	ext = strings.ToLower(ext)
	types := map[string]string{
		".txt":  "文本文件",
		".pdf":  "PDF文档",
		".doc":  "Word文档",
		".docx": "Word文档",
		".xls":  "Excel表格",
		".xlsx": "Excel表格",
		".ppt":  "PPT演示文稿",
		".pptx": "PPT演示文稿",
		".jpg":  "图片",
		".jpeg": "图片",
		".png":  "图片",
		".gif":  "图片",
		".bmp":  "图片",
		".mp4":  "视频",
		".avi":  "视频",
		".mov":  "视频",
		".mp3":  "音频",
		".wav":  "音频",
	}

	if t, ok := types[ext]; ok {
		return t
	}
	return "未知类型"
}
