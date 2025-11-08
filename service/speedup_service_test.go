package service

import (
	"testing"
	"time"

	"speedtestup/api"
	"speedtestup/config"
)

func TestNewSpeedupService(t *testing.T) {
	cfg := config.NewDefaultConfig()
	speedTestCNClient := api.NewSpeedTestCNClient("192.168.1.1")
	speedupService := NewSpeedupService(speedTestCNClient, cfg)

	if speedupService == nil {
		t.Fatal("NewSpeedupService should not return nil")
	}

	if speedupService.apiClient == nil {
		t.Fatal("SpeedupService apiClient should not be nil")
	}

	if speedupService.config == nil {
		t.Fatal("SpeedupService config should not be nil")
	}

	if speedupService.selfCheck == nil {
		t.Fatal("SpeedupService selfCheck should not be nil")
	}
}

func TestSpeedupService_ShouldSelfCheck(t *testing.T) {
	cfg := config.NewDefaultConfig()
	cfg.Speedup.SelfCheck.Enabled = true
	speedTestCNClient := api.NewSpeedTestCNClient("")
	speedupService := NewSpeedupService(speedTestCNClient, cfg)

	// 初始状态下不应该执行自检
	if speedupService.ShouldSelfCheck() {
		t.Error("ShouldSelfCheck should return false when lastExecute is zero")
	}

	// 设置上次执行时间
	speedupService.lastExecute = time.Now()
	if speedupService.ShouldSelfCheck() {
		t.Error("ShouldSelfCheck should return false when interval not reached")
	}
}
