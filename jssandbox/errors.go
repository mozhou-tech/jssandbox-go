package jssandbox

import "fmt"

// ErrorCode 错误代码类型
type ErrorCode string

const (
	// ErrCodeTimeout 执行超时
	ErrCodeTimeout ErrorCode = "TIMEOUT"
	// ErrCodeInvalidInput 无效输入
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
	// ErrCodeFileNotFound 文件未找到
	ErrCodeFileNotFound ErrorCode = "FILE_NOT_FOUND"
	// ErrCodeHTTPError HTTP 请求错误
	ErrCodeHTTPError ErrorCode = "HTTP_ERROR"
	// ErrCodeFileSystemError 文件系统错误
	ErrCodeFileSystemError ErrorCode = "FILE_SYSTEM_ERROR"
	// ErrCodeBrowserError 浏览器操作错误
	ErrCodeBrowserError ErrorCode = "BROWSER_ERROR"
	// ErrCodeDocumentError 文档处理错误
	ErrCodeDocumentError ErrorCode = "DOCUMENT_ERROR"
	// ErrCodeImageError 图片处理错误
	ErrCodeImageError ErrorCode = "IMAGE_ERROR"
	// ErrCodeSystemError 系统操作错误
	ErrCodeSystemError ErrorCode = "SYSTEM_ERROR"
	// ErrCodeUnknown 未知错误
	ErrCodeUnknown ErrorCode = "UNKNOWN_ERROR"
)

// SandboxError 沙盒错误类型
type SandboxError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

// Error 实现 error 接口
func (e *SandboxError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 返回原始错误（用于 errors.Unwrap）
func (e *SandboxError) Unwrap() error {
	return e.Cause
}

// NewSandboxError 创建新的沙盒错误
func NewSandboxError(code ErrorCode, message string) *SandboxError {
	return &SandboxError{
		Code:    code,
		Message: message,
	}
}

// NewSandboxErrorWithCause 创建带原因的错误
func NewSandboxErrorWithCause(code ErrorCode, message string, cause error) *SandboxError {
	return &SandboxError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsTimeout 判断是否为超时错误
func (e *SandboxError) IsTimeout() bool {
	return e.Code == ErrCodeTimeout
}

// IsFileNotFound 判断是否为文件未找到错误
func (e *SandboxError) IsFileNotFound() bool {
	return e.Code == ErrCodeFileNotFound
}

