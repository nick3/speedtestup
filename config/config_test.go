package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	// 验证配置结构体不为nil
	assert.NotNil(t, cfg, "Config should not be nil")
	assert.NotNil(t, &cfg.Speedup, "Speedup config should not be nil")
	assert.NotNil(t, &cfg.Logging, "Logging config should not be nil")

	// 验证提速配置默认值
	assert.False(t, cfg.Speedup.Enabled, "Speedup should be disabled by default")
	assert.True(t, cfg.Speedup.DownAcc, "DownAcc should be enabled by default")
	assert.True(t, cfg.Speedup.UpAcc, "UpAcc should be enabled by default")
	assert.Equal(t, 10*time.Minute, cfg.Speedup.CheckInterval, "CheckInterval should be 10 minutes")
	assert.Equal(t, "0 0 * * 1", cfg.Speedup.ReopenSchedule, "ReopenSchedule should be weekly")

	// 验证IP绑定配置默认值
	assert.False(t, cfg.Speedup.IPBinding.Enabled, "IP binding should be disabled by default")
	assert.Equal(t, "wan", cfg.Speedup.IPBinding.Interface, "Default interface should be 'wan'")
	assert.Empty(t, cfg.Speedup.IPBinding.BindIP, "BindIP should be empty by default")

	// 验证自动恢复配置默认值
	assert.True(t, cfg.Speedup.AutoRecovery.Enabled, "Auto recovery should be enabled by default")
	assert.Equal(t, 3, cfg.Speedup.AutoRecovery.MaxRetries, "Default max retries should be 3")
	assert.Equal(t, 5*time.Minute, cfg.Speedup.AutoRecovery.RetryInterval, "Default retry interval should be 5 minutes")

	// 验证自检配置默认值
	assert.True(t, cfg.Speedup.SelfCheck.Enabled, "SelfCheck should be enabled by default")
	assert.Equal(t, 168*time.Hour, cfg.Speedup.SelfCheck.Interval, "SelfCheck interval should be 7 days")

	// 验证日志配置默认值
	assert.Equal(t, "info", cfg.Logging.Level, "Logging level should be info by default")
	assert.Equal(t, "stdout", cfg.Logging.Output, "Logging output should be stdout by default")
	assert.Empty(t, cfg.Logging.File, "Logging file should be empty by default")

	// 验证其他配置默认值
	assert.False(t, cfg.Speedup.Logging, "Logging should be false by default")
	assert.False(t, cfg.Speedup.Verbose, "Verbose should be false by default")
	assert.False(t, cfg.Speedup.MoreOptions, "MoreOptions should be false by default")
}

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	tempConfig := `{
		"speedup": {
			"enabled": true,
			"down_acc": true,
			"up_acc": true,
			"check_interval": "15m",
			"reopen_schedule": "0 0 * * 2",
			"ip_binding": {
				"enabled": true,
				"interface": "eth0",
				"bind_ip": "192.168.1.100"
			},
			"auto_recovery": {
				"enabled": true,
				"max_retries": 5,
				"retry_interval": "10m"
			},
			"self_check": {
				"enabled": true,
				"interval": "168h"
			},
			"logging": true,
			"verbose": true,
			"more": true
		},
		"logging": {
			"level": "debug",
			"output": "file",
			"file": "/var/log/speedtestup.log"
		}
	}`

	tempFile := "temp_test_config.json"
	err := os.WriteFile(tempFile, []byte(tempConfig), 0644)
	assert.NoError(t, err)

	defer os.Remove(tempFile)

	// 加载配置
	cfg, err := LoadConfig(tempFile)
	assert.NoError(t, err)
	assert.NotNil(t, cfg, "Loaded config should not be nil")

	// 验证加载的配置
	assert.True(t, cfg.Speedup.Enabled, "Speedup should be enabled")
	assert.Equal(t, 15*time.Minute, cfg.Speedup.CheckInterval, "CheckInterval should be 15 minutes")
	assert.Equal(t, "0 0 * * 2", cfg.Speedup.ReopenSchedule, "ReopenSchedule should be weekly on Tuesday")

	// 验证IP绑定配置
	assert.True(t, cfg.Speedup.IPBinding.Enabled, "IP binding should be enabled")
	assert.Equal(t, "eth0", cfg.Speedup.IPBinding.Interface, "Interface should be eth0")
	assert.Equal(t, "192.168.1.100", cfg.Speedup.IPBinding.BindIP, "BindIP should be 192.168.1.100")

	// 验证自动恢复配置
	assert.Equal(t, 5, cfg.Speedup.AutoRecovery.MaxRetries, "MaxRetries should be 5")
	assert.Equal(t, 10*time.Minute, cfg.Speedup.AutoRecovery.RetryInterval, "RetryInterval should be 10 minutes")

	// 验证日志配置
	assert.Equal(t, "debug", cfg.Logging.Level, "Logging level should be debug")
	assert.Equal(t, "file", cfg.Logging.Output, "Logging output should be file")
	assert.Equal(t, "/var/log/speedtestup.log", cfg.Logging.File, "Logging file path should be set")

	// 验证其他配置
	assert.True(t, cfg.Speedup.Verbose, "Verbose should be true")
	assert.True(t, cfg.Speedup.Logging, "Logging should be true")
}

func TestLoadConfigWithInvalidFile(t *testing.T) {
	cfg, err := LoadConfig("nonexistent_file.json")
	assert.Error(t, err, "Should return error for nonexistent file")
	assert.Nil(t, cfg, "Config should be nil for nonexistent file")
}

func TestLoadConfigWithInvalidJSON(t *testing.T) {
	// 创建无效的JSON配置文件
	invalidConfig := `{
		"speedup": {
			"enabled": true
		}
	` // 缺少闭合括号

	tempFile := "temp_invalid_config.json"
	err := os.WriteFile(tempFile, []byte(invalidConfig), 0644)
	assert.NoError(t, err)

	defer os.Remove(tempFile)

	cfg, err := LoadConfig(tempFile)
	assert.Error(t, err, "Should return error for invalid JSON")
	assert.Nil(t, cfg, "Config should be nil for invalid JSON")
}

func TestLoadConfigWithMinimalConfig(t *testing.T) {
	// 测试最小配置
	minimalConfig := `{
		"speedup": {
			"enabled": true
		},
		"logging": {
			"level": "info"
		}
	}`

	tempFile := "temp_minimal_config.json"
	err := os.WriteFile(tempFile, []byte(minimalConfig), 0644)
	assert.NoError(t, err)

	defer os.Remove(tempFile)

	cfg, err := LoadConfig(tempFile)
	assert.NoError(t, err)
	assert.NotNil(t, cfg, "Config should not be nil")

	// 验证最小配置有效
	assert.True(t, cfg.Speedup.Enabled, "Speedup should be enabled")
	assert.Equal(t, "info", cfg.Logging.Level, "Logging level should be info")

	// 验证默认值被应用
	assert.Equal(t, 10*time.Minute, cfg.Speedup.CheckInterval, "Default CheckInterval should be applied")
	assert.Equal(t, "0 0 * * 1", cfg.Speedup.ReopenSchedule, "Default ReopenSchedule should be applied")
}

func TestConfigWithNegativeValues(t *testing.T) {
	// 测试负值处理
	negativeConfig := `{
		"speedup": {
			"enabled": true,
			"check_interval": "-1h",
			"auto_recovery": {
				"max_retries": -1,
				"retry_interval": "-1m"
			},
			"self_check": {
				"interval": "-1h"
			}
		},
		"logging": {
			"level": ""
		}
	}`

	tempFile := "temp_negative_config.json"
	err := os.WriteFile(tempFile, []byte(negativeConfig), 0644)
	assert.NoError(t, err)

	defer os.Remove(tempFile)

	cfg, err := LoadConfig(tempFile)
	assert.NoError(t, err, "Should load config even with negative values")

	// 验证负值被处理为默认值
	assert.Equal(t, 10*time.Minute, cfg.Speedup.CheckInterval, "Negative CheckInterval should default to 10 minutes")
	assert.Equal(t, 3, cfg.Speedup.AutoRecovery.MaxRetries, "Negative MaxRetries should default to 3")
	assert.Equal(t, 5*time.Minute, cfg.Speedup.AutoRecovery.RetryInterval, "Negative RetryInterval should default to 5 minutes")
	assert.Equal(t, 168*time.Hour, cfg.Speedup.SelfCheck.Interval, "Negative SelfCheck interval should default to 168 hours")
	assert.Equal(t, "info", cfg.Logging.Level, "Empty log level should default to info")
}