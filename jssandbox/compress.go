package jssandbox

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
)

// registerCompress 注册压缩/解压缩功能到JavaScript运行时
func (sb *Sandbox) registerCompress() {
	// 压缩为ZIP
	sb.vm.Set("compressZip", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供文件列表和输出路径参数",
			})
		}

		// 解析文件列表
		var files []string
		if filesVal := call.Arguments[0]; filesVal != nil && !goja.IsUndefined(filesVal) {
			if filesObj := filesVal.ToObject(sb.vm); filesObj != nil {
				if filesArray, ok := filesObj.Export().([]interface{}); ok {
					for _, f := range filesArray {
						if str, ok := f.(string); ok {
							files = append(files, str)
						}
					}
				} else if str, ok := filesVal.Export().(string); ok {
					// 单个文件
					files = []string{str}
				}
			} else if str := filesVal.String(); str != "" {
				files = []string{str}
			}
		}

		if len(files) == 0 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "文件列表为空",
			})
		}

		outputPath := call.Arguments[1].String()

		// 创建ZIP文件
		zipFile, err := os.Create(outputPath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("创建ZIP文件失败: %v", err),
			})
		}
		defer zipFile.Close()

		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		// 添加文件到ZIP
		for _, file := range files {
			if err := sb.addFileToZip(zipWriter, file); err != nil {
				return sb.vm.ToValue(map[string]interface{}{
					"error": fmt.Sprintf("添加文件 %s 失败: %v", file, err),
				})
			}
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"path":    outputPath,
		})
	})

	// 解压ZIP
	sb.vm.Set("extractZip", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return sb.vm.ToValue(map[string]interface{}{
				"error": "需要提供ZIP文件路径和输出目录参数",
			})
		}

		zipPath := call.Arguments[0].String()
		outputDir := call.Arguments[1].String()

		// 打开ZIP文件
		zipReader, err := zip.OpenReader(zipPath)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("打开ZIP文件失败: %v", err),
			})
		}
		defer zipReader.Close()

		// 创建输出目录
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("创建输出目录失败: %v", err),
			})
		}

		// 解压所有文件
		var extractedFiles []string
		for _, file := range zipReader.File {
			path := filepath.Join(outputDir, file.Name)

			// 检查路径安全性
			if !strings.HasPrefix(path, filepath.Clean(outputDir)+string(os.PathSeparator)) {
				continue
			}

			if file.FileInfo().IsDir() {
				os.MkdirAll(path, file.FileInfo().Mode())
				continue
			}

			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return sb.vm.ToValue(map[string]interface{}{
					"error": fmt.Sprintf("创建目录失败: %v", err),
				})
			}

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
			if err != nil {
				return sb.vm.ToValue(map[string]interface{}{
					"error": fmt.Sprintf("创建文件失败: %v", err),
				})
			}

			rc, err := file.Open()
			if err != nil {
				outFile.Close()
				return sb.vm.ToValue(map[string]interface{}{
					"error": fmt.Sprintf("打开ZIP内文件失败: %v", err),
				})
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			if err != nil {
				return sb.vm.ToValue(map[string]interface{}{
					"error": fmt.Sprintf("解压文件失败: %v", err),
				})
			}

			extractedFiles = append(extractedFiles, path)
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"files":   extractedFiles,
		})
	})

	// GZIP压缩
	sb.vm.Set("compressGzip", func(data string) goja.Value {
		var buf strings.Builder
		writer := gzip.NewWriter(&buf)
		if _, err := writer.Write([]byte(data)); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("压缩失败: %v", err),
			})
		}
		if err := writer.Close(); err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("关闭压缩器失败: %v", err),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"data":    buf.String(),
		})
	})

	// GZIP解压
	sb.vm.Set("decompressGzip", func(compressed string) goja.Value {
		reader, err := gzip.NewReader(strings.NewReader(compressed))
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("创建解压器失败: %v", err),
			})
		}
		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			return sb.vm.ToValue(map[string]interface{}{
				"error": fmt.Sprintf("解压失败: %v", err),
			})
		}

		return sb.vm.ToValue(map[string]interface{}{
			"success": true,
			"data":    string(data),
		})
	})
}

// addFileToZip 添加文件到ZIP
func (sb *Sandbox) addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
