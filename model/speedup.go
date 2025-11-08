package model

import "time"

// SpeedupInfo 提速信息
type SpeedupInfo struct {
	// 基础信息
	IP          string    `json:"ip"`           // 出口 IP 地址
	UpdatedAt   time.Time `json:"updated_at"`   // 提速开始时间
	CanSpeed    int       `json:"can_speed"`    // 是否可以提速 (0: 不支持, 1: 支持)

	// 上行带宽信息
	UpHBandwidth     int       `json:"up_h_bandwidth"`      // 一类上行带宽 (Mbps)
	UpHExpire        time.Time `json:"up_h_expire"`         // 一类上行带宽提速截止时间
	Up100Bandwidth   int       `json:"up_100_bandwidth"`    // 二类上行带宽 (Mbps)
	Up100Expire      time.Time `json:"up_100_expire"`       // 二类上行带宽提速截止时间

	// 下行带宽信息
	DownloadBandwidth int       `json:"download_bandwidth"` // 下行带宽 (Mbps)
	DownExpire        time.Time `json:"down_expire"`        // 下行提速截止时间

	// 套餐信息
	Package1UpH   int       `json:"package_1_up_h"`   // 一类套餐上行带宽 (Mbps)
	Package1Down  int       `json:"package_1_down"`   // 一类套餐下行带宽 (Mbps)
	Package1Expire time.Time `json:"package_1_expire"` // 一类套餐提速截止时间

	Package2UpH   int       `json:"package_2_up_h"`   // 二类套餐上行带宽 (Mbps)
	Package2Down  int       `json:"package_2_down"`   // 二类套餐下行带宽 (Mbps)
	Package2Expire time.Time `json:"package_2_expire"` // 二类套餐提速截止时间

	// 提速状态
	UpSpeedupActive   bool `json:"up_speedup_active"`   // 上行提速是否激活
	DownSpeedupActive bool `json:"down_speedup_active"` // 下行提速是否激活
}

// Status 状态信息
type Status struct {
	Running           bool        `json:"running"`            // 程序是否在运行
	LastExecute       time.Time   `json:"last_execute"`       // 上次执行时间
	CheckInterval     string      `json:"check_interval"`     // 检查间隔
	SelfCheck         bool        `json:"self_check"`         // 是否启用 7 天自检
	AutoRecovery      bool        `json:"auto_recovery"`      // 是否启用自动恢复
	CurrentIP         string      `json:"current_ip"`         // 当前公网 IP
	IPBindingEnabled  bool        `json:"ip_binding_enabled"` // 是否启用 IP 绑定
	SpeedupInfo       *SpeedupInfo `json:"speedup_info,omitempty"` // 提速信息
}
