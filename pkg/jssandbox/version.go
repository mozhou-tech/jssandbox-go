package jssandbox

// Version 版本信息
// 这些变量在构建时通过 ldflags 注入
var (
	// Version 版本号，例如 "v1.0.0"
	Version = "dev"
	// BuildTime 构建时间，格式: "2006-01-02_15:04:05"
	BuildTime = "unknown"
	// GitCommit Git 提交哈希
	GitCommit = "unknown"
)

// GetVersion 返回版本信息字符串
func GetVersion() string {
	return Version
}

// GetBuildInfo 返回构建信息
func GetBuildInfo() map[string]string {
	return map[string]string{
		"version":   Version,
		"buildTime": BuildTime,
		"gitCommit": GitCommit,
	}
}

