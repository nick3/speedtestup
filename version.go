package main

import "fmt"

var (
	// 版本信息
	Version   = "dev"     // 语义化版本号
	BuildTime = "unknown" // 构建时间
	GitCommit = "unknown" // Git commit hash
	GoVersion = "unknown" // Go 版本
)

// 获取版本信息的函数
func GetVersionInfo() string {
	return fmt.Sprintf(
		"Version: %s\nBuild Time: %s\nGit Commit: %s\nGo Version: %s",
		Version, BuildTime, GitCommit, GoVersion,
	)
}
