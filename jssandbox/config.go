package jssandbox

import "time"

// Config 沙盒配置
type Config struct {
	// DefaultTimeout 默认执行超时时间
	DefaultTimeout time.Duration
	// HTTPTimeout HTTP 请求默认超时时间
	HTTPTimeout time.Duration
	// BrowserTimeout 浏览器操作默认超时时间
	BrowserTimeout time.Duration
	// MaxFileSize 最大文件大小（字节），0 表示不限制
	MaxFileSize int64
	// AllowedFileTypes 允许的文件类型列表，空列表表示不限制
	AllowedFileTypes []string
	// EnableBrowser 是否启用浏览器功能
	EnableBrowser bool
	// EnableFileSystem 是否启用文件系统功能
	EnableFileSystem bool
	// EnableHTTP 是否启用 HTTP 功能
	EnableHTTP bool
	// EnableDocuments 是否启用文档处理功能
	EnableDocuments bool
	// EnableImageProcessing 是否启用图片处理功能
	EnableImageProcessing bool
	// EnableVideoProcessing 是否启用视频处理功能
	EnableVideoProcessing bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DefaultTimeout:        30 * time.Second,
		HTTPTimeout:          30 * time.Second,
		BrowserTimeout:       60 * time.Second,
		MaxFileSize:          100 * 1024 * 1024, // 100MB
		AllowedFileTypes:     []string{},
		EnableBrowser:        true,
		EnableFileSystem:     true,
		EnableHTTP:           true,
		EnableDocuments:      true,
		EnableImageProcessing: true,
		EnableVideoProcessing: true,
	}
}

// WithTimeout 设置默认超时时间
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.DefaultTimeout = timeout
	return c
}

// WithHTTPTimeout 设置 HTTP 超时时间
func (c *Config) WithHTTPTimeout(timeout time.Duration) *Config {
	c.HTTPTimeout = timeout
	return c
}

// WithBrowserTimeout 设置浏览器超时时间
func (c *Config) WithBrowserTimeout(timeout time.Duration) *Config {
	c.BrowserTimeout = timeout
	return c
}

// WithMaxFileSize 设置最大文件大小
func (c *Config) WithMaxFileSize(size int64) *Config {
	c.MaxFileSize = size
	return c
}

// WithAllowedFileTypes 设置允许的文件类型
func (c *Config) WithAllowedFileTypes(types []string) *Config {
	c.AllowedFileTypes = types
	return c
}

// DisableBrowser 禁用浏览器功能
func (c *Config) DisableBrowser() *Config {
	c.EnableBrowser = false
	return c
}

// DisableFileSystem 禁用文件系统功能
func (c *Config) DisableFileSystem() *Config {
	c.EnableFileSystem = false
	return c
}

// DisableHTTP 禁用 HTTP 功能
func (c *Config) DisableHTTP() *Config {
	c.EnableHTTP = false
	return c
}

