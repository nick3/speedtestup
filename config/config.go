package config

import (
	"time"
)

// Config 定义应用程序配置结构
type Config struct {
	// 提速服务配置
	Speedup SpeedupConfig `json:"speedup" yaml:"speedup"`

	// 日志配置
	Logging LoggingConfig `json:"logging" yaml:"logging"`
}

// SpeedupConfig 提速服务配置
type SpeedupConfig struct {
	// 检测间隔
	CheckInterval time.Duration `json:"check_interval" yaml:"check_interval"`

	// 重新开启提速的定时任务（cron 表达式）
	ReopenSchedule string `json:"reopen_schedule" yaml:"reopen_schedule"`

	// IP 绑定配置
	IPBinding IPBindingConfig `json:"ip_binding" yaml:"ip_binding"`

	// 自动恢复配置
	AutoRecovery AutoRecoveryConfig `json:"auto_recovery" yaml:"auto_recovery"`

	// 自检配置
	SelfCheck SelfCheckConfig `json:"self_check" yaml:"self_check"`

	// 提速功能开关
	Enabled     bool `json:"enabled" yaml:"enabled"`
	DownAcc     bool `json:"down_acc" yaml:"down_acc"`
	UpAcc       bool `json:"up_acc" yaml:"up_acc"`
	Logging     bool `json:"logging" yaml:"logging"`
	Verbose     bool `json:"verbose" yaml:"verbose"`
	MoreOptions bool `json:"more" yaml:"more"`
}

// IPBindingConfig IP 绑定配置
type IPBindingConfig struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	Interface string `json:"interface" yaml:"interface"` // 网络接口名称（如 wan、pppoe-wan）
	BindIP    string `json:"bind_ip" yaml:"bind_ip"`     // 绑定的 IP 地址
}

// AutoRecoveryConfig 自动恢复配置
type AutoRecoveryConfig struct {
	Enabled       bool          `json:"enabled" yaml:"enabled"`
	MaxRetries    int           `json:"max_retries" yaml:"max_retries"`        // 最大重试次数
	RetryInterval time.Duration `json:"retry_interval" yaml:"retry_interval"` // 重试间隔
}

// SelfCheckConfig 自检配置
type SelfCheckConfig struct {
	Enabled  bool          `json:"enabled" yaml:"enabled"`
	Interval time.Duration `json:"interval" yaml:"interval"` // 自检间隔（默认 7 天）
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string `json:"level" yaml:"level"`  // 日志级别（debug, info, warn, error）
	Output string `json:"output" yaml:"output"` // 输出方式（stdout, file）
	File   string `json:"file" yaml:"file"`     // 日志文件路径
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *Config {
	cfg := &Config{}

	// 设置默认提速配置
	cfg.Speedup.CheckInterval = 10 * time.Minute
	cfg.Speedup.ReopenSchedule = "0 0 * * 1" // 每周一 0:00

	// 设置默认 IP 绑定配置
	cfg.Speedup.IPBinding.Enabled = false
	cfg.Speedup.IPBinding.Interface = "wan"
	cfg.Speedup.IPBinding.BindIP = ""

	// 设置默认自动恢复配置
	cfg.Speedup.AutoRecovery.Enabled = true
	cfg.Speedup.AutoRecovery.MaxRetries = 3
	cfg.Speedup.AutoRecovery.RetryInterval = 5 * time.Minute

	// 设置默认自检配置
	cfg.Speedup.SelfCheck.Enabled = true
	cfg.Speedup.SelfCheck.Interval = 168 * time.Hour // 7 天

	// 设置默认功能开关
	cfg.Speedup.Enabled = false
	cfg.Speedup.DownAcc = true
	cfg.Speedup.UpAcc = true
	cfg.Speedup.Logging = false
	cfg.Speedup.Verbose = false
	cfg.Speedup.MoreOptions = false

	// 设置默认日志配置
	cfg.Logging.Level = "info"
	cfg.Logging.Output = "stdout"
	cfg.Logging.File = ""

	return cfg
}
