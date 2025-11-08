package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	// 验证测速配置
	assert.NotEmpty(t, cfg.SpeedTest.URLs)
	assert.Equal(t, 30*time.Second, cfg.SpeedTest.Timeout)
	assert.Equal(t, 3, cfg.SpeedTest.RetryCount)

	// 验证IP检查配置
	assert.NotEmpty(t, cfg.IPCheck.URLs)
	assert.Equal(t, 10*time.Minute, cfg.IPCheck.CheckInterval)
	assert.Equal(t, 10*time.Second, cfg.IPCheck.Timeout)

	// 验证日志配置
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "", cfg.Logging.File)

	// 验证服务器配置
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Server.Host)

	// 验证通用配置
	assert.False(t, cfg.General.Verbose)
}

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	tempConfig := `{
		"speed_test": {
			"urls": ["http://example.com/test"],
			"timeout": 60000000000,
			"retry_count": 5
		},
		"ip_check": {
			"urls": ["http://example.com/ip"],
			"check_interval": 1200000000000,
			"timeout": 20000000000
		},
		"logging": {
			"level": "debug",
			"file": "/tmp/test.log"
		},
		"server": {
			"port": "9090",
			"host": "0.0.0.0"
		},
		"general": {
			"verbose": true
		}
	}`

	tempFile := "temp_test_config.json"
	err := os.WriteFile(tempFile, []byte(tempConfig), 0644)
	assert.NoError(t, err)

	defer os.Remove(tempFile)

	// 加载配置
	cfg, err := LoadConfig(tempFile)
	assert.NoError(t, err)

	// 验证加载的配置
	assert.Equal(t, "http://example.com/test", cfg.SpeedTest.URLs[0])
	assert.Equal(t, 60*time.Second, cfg.SpeedTest.Timeout)
	assert.Equal(t, 5, cfg.SpeedTest.RetryCount)

	assert.Equal(t, "http://example.com/ip", cfg.IPCheck.URLs[0])
	assert.Equal(t, 20*time.Minute, cfg.IPCheck.CheckInterval)
	assert.Equal(t, 20*time.Second, cfg.IPCheck.Timeout)

	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "/tmp/test.log", cfg.Logging.File)

	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)

	assert.True(t, cfg.General.Verbose)
}

func TestLoadConfigWithInvalidFile(t *testing.T) {
	cfg, err := LoadConfig("nonexistent_file.json")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoadConfigWithInvalidJSON(t *testing.T) {
	// 创建无效的JSON配置文件
	invalidConfig := `{
		"speed_test": {
			"urls": ["http://example.com/test"
		}
	}`

	tempFile := "temp_invalid_config.json"
	err := os.WriteFile(tempFile, []byte(invalidConfig), 0644)
	assert.NoError(t, err)

	defer os.Remove(tempFile)

	cfg, err := LoadConfig(tempFile)
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestSaveConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	tempFile := "temp_save_config.json"
	err := SaveConfig(cfg, tempFile)
	assert.NoError(t, err)

	defer os.Remove(tempFile)

	// 重新加载保存的配置
	loadedCfg, err := LoadConfig(tempFile)
	assert.NoError(t, err)
	assert.Equal(t, cfg.SpeedTest.URLs, loadedCfg.SpeedTest.URLs)
	assert.Equal(t, cfg.SpeedTest.Timeout, loadedCfg.SpeedTest.Timeout)
	assert.Equal(t, cfg.SpeedTest.RetryCount, loadedCfg.SpeedTest.RetryCount)
	assert.Equal(t, cfg.IPCheck.CheckInterval, loadedCfg.IPCheck.CheckInterval)
}

func TestValidateAndSetDefaults(t *testing.T) {
	cfg := NewDefaultConfig()

	// 测试负值的处理
	cfg.SpeedTest.Timeout = -1 * time.Second
	cfg.SpeedTest.RetryCount = -1
	cfg.IPCheck.CheckInterval = -1 * time.Minute
	cfg.IPCheck.Timeout = -1 * time.Second
	cfg.Server.Port = ""
	cfg.Server.Host = ""
	cfg.Logging.Level = ""

	// 重新加载配置以验证默认值设置
	tempConfig := `{
		"speed_test": {
			"urls": ["http://example.com/test"],
			"timeout": -1000000000,
			"retry_count": -1
		},
		"ip_check": {
			"urls": ["http://example.com/ip"],
			"check_interval": -60000000000,
			"timeout": -10000000000
		},
		"logging": {
			"level": "",
			"file": "/tmp/test.log"
		},
		"server": {
			"port": "",
			"host": ""
		},
		"general": {
			"verbose": false
		}
	}`

	tempFile := "temp_validate_config.json"
	err := os.WriteFile(tempFile, []byte(tempConfig), 0644)
	assert.NoError(t, err)

	defer os.Remove(tempFile)

	cfg, err = LoadConfig(tempFile)
	assert.NoError(t, err)

	// 验证默认值是否已设置
	assert.Equal(t, 30*time.Second, cfg.SpeedTest.Timeout)  // 默认值
	assert.Equal(t, 3, cfg.SpeedTest.RetryCount)             // 默认值
	assert.Equal(t, 10*time.Minute, cfg.IPCheck.CheckInterval) // 默认值
	assert.Equal(t, 10*time.Second, cfg.IPCheck.Timeout)     // 默认值
	assert.Equal(t, "8080", cfg.Server.Port)                 // 默认值
	assert.Equal(t, "localhost", cfg.Server.Host)            // 默认值
	assert.Equal(t, "info", cfg.Logging.Level)               // 默认值
}