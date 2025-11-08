package config

import (
	"encoding/json"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// LoadConfig 从文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	cfg := NewDefaultConfig()

	// 尝试解析JSON格式
	if err := json.Unmarshal(data, cfg); err != nil {
		// 如果JSON解析失败，尝试YAML格式
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// 验证并设置合理的默认值
	validateAndSetDefaults(cfg)

	return cfg, nil
}

// SaveConfig 将配置保存到文件
func SaveConfig(cfg *Config, filePath string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// validateAndSetDefaults 验证配置并设置合理的默认值
func validateAndSetDefaults(cfg *Config) {
	// 验证提速配置
	if cfg.Speedup.CheckInterval <= 0 {
		cfg.Speedup.CheckInterval = 10 * time.Minute
	}

	if cfg.Speedup.ReopenSchedule == "" {
		cfg.Speedup.ReopenSchedule = "0 0 * * 1"
	}

	// 验证IP绑定配置
	if cfg.Speedup.IPBinding.Interface == "" {
		cfg.Speedup.IPBinding.Interface = "wan"
	}

	// 验证自动恢复配置
	if cfg.Speedup.AutoRecovery.MaxRetries <= 0 {
		cfg.Speedup.AutoRecovery.MaxRetries = 3
	}
	if cfg.Speedup.AutoRecovery.RetryInterval <= 0 {
		cfg.Speedup.AutoRecovery.RetryInterval = 5 * time.Minute
	}

	// 验证自检配置
	if cfg.Speedup.SelfCheck.Interval <= 0 {
		cfg.Speedup.SelfCheck.Interval = 168 * time.Hour // 7 天
	}

	// 设置默认日志级别
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Logging.Output == "" {
		cfg.Logging.Output = "stdout"
	}
}
