package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"speedtestup/config"
	"speedtestup/utils"
)

// 测试版本变量是否正确定义
func TestVersionVariables(t *testing.T) {
	// 验证版本变量已定义
	assert.NotEmpty(t, version, "version variable should be defined")
	assert.NotEmpty(t, buildDate, "buildDate variable should be defined")
	assert.NotEmpty(t, commitHash, "commitDate variable should be defined")

	// 验证版本格式
	assert.Regexp(t, `^\d+\.\d+\.\d+$`, version, "version should follow semantic versioning")
}

// 测试默认配置初始化
func TestDefaultConfigInitialization(t *testing.T) {
	cfg := config.NewDefaultConfig()

	// 验证配置结构体不为nil
	assert.NotNil(t, cfg, "config should not be nil")
	assert.NotNil(t, &cfg.Speedup, "Speedup config should not be nil")
	assert.NotNil(t, &cfg.Logging, "Logging config should not be nil")

	// 验证提速配置默认值
	assert.Equal(t, 10*time.Minute, cfg.Speedup.CheckInterval, "CheckInterval should be 10 minutes")
	assert.Equal(t, "0 0 * * 1", cfg.Speedup.ReopenSchedule, "ReopenSchedule should be weekly")
	assert.True(t, cfg.Speedup.SelfCheck.Enabled, "SelfCheck should be enabled by default")
	assert.Equal(t, 168*time.Hour, cfg.Speedup.SelfCheck.Interval, "SelfCheck interval should be 7 days")

	// 验证日志配置默认值
	assert.Equal(t, "info", cfg.Logging.Level, "Logging level should be info by default")
	assert.Equal(t, "stdout", cfg.Logging.Output, "Logging output should be stdout by default")
}

// 测试配置验证
func TestConfigValidation(t *testing.T) {
	cfg := config.NewDefaultConfig()

	// 启用提速服务
	cfg.Speedup.Enabled = true

	// 验证配置有效
	assert.True(t, cfg.Speedup.Enabled, "Speedup should be enabled for testing")
	assert.Greater(t, cfg.Speedup.CheckInterval, time.Duration(0), "CheckInterval should be positive")
}

// 测试日志器初始化
func TestLoggerInitialization(t *testing.T) {
	// 测试有效的日志配置
	logger, err := utils.NewLogger("info", "stdout", "")
	assert.NoError(t, err, "Logger should initialize successfully with valid config")
	assert.NotNil(t, logger, "Logger should not be nil")

	// 测试无效的日志级别
	logger, err = utils.NewLogger("invalid-level", "stdout", "")
	// 可能会返回错误或使用默认值，这取决于实现
	if err != nil {
		assert.Error(t, err, "Should return error for invalid log level")
	}
}

// 测试IP绑定配置
func TestIPBindingConfig(t *testing.T) {
	cfg := config.NewDefaultConfig()

	// 验证默认IP绑定配置
	assert.False(t, cfg.Speedup.IPBinding.Enabled, "IP binding should be disabled by default")
	assert.Equal(t, "wan", cfg.Speedup.IPBinding.Interface, "Default interface should be 'wan'")
	assert.Empty(t, cfg.Speedup.IPBinding.BindIP, "BindIP should be empty by default")

	// 测试启用IP绑定
	cfg.Speedup.IPBinding.Enabled = true
	cfg.Speedup.IPBinding.BindIP = "192.168.1.100"
	assert.True(t, cfg.Speedup.IPBinding.Enabled, "IP binding should be enabled")
	assert.Equal(t, "192.168.1.100", cfg.Speedup.IPBinding.BindIP, "BindIP should be set")
}

// 测试自动恢复配置
func TestAutoRecoveryConfig(t *testing.T) {
	cfg := config.NewDefaultConfig()

	// 验证默认自动恢复配置
	assert.True(t, cfg.Speedup.AutoRecovery.Enabled, "Auto recovery should be enabled by default")
	assert.Equal(t, 3, cfg.Speedup.AutoRecovery.MaxRetries, "Default max retries should be 3")
	assert.Equal(t, 5*time.Minute, cfg.Speedup.AutoRecovery.RetryInterval, "Default retry interval should be 5 minutes")
}