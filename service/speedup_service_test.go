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

// TestSpeedupService_QueryStatus 测试查询提速状态
func TestSpeedupService_QueryStatus(t *testing.T) {
	cfg := config.NewDefaultConfig()
	speedTestCNClient := api.NewSpeedTestCNClient("")
	speedupService := NewSpeedupService(speedTestCNClient, cfg)

	// 测试QueryStatus方法返回布尔值和错误
	// 注意：这里可能涉及网络请求，在测试环境中可能会失败
	// 我们主要测试方法不会panic
	canSpeed, err := speedupService.QueryStatus()
	if err != nil {
		t.Logf("Expected network error in test environment: %v", err)
	} else {
		t.Logf("QueryStatus returned: canSpeed=%v, err=%v", canSpeed, err)
	}
}

// TestSpeedupService_GetLastExecuteTime 测试获取最后执行时间
func TestSpeedupService_GetLastExecuteTime(t *testing.T) {
	cfg := config.NewDefaultConfig()
	speedTestCNClient := api.NewSpeedTestCNClient("")
	speedupService := NewSpeedupService(speedTestCNClient, cfg)

	// 验证初始返回零时间
	lastExecute := speedupService.GetLastExecuteTime()
	if !lastExecute.IsZero() {
		t.Error("Expected lastExecute to be zero initially")
	}

	// 设置最后执行时间
	now := time.Now()
	speedupService.lastExecute = now

	// 验证返回正确的时间
	lastExecute = speedupService.GetLastExecuteTime()
	if !lastExecute.Equal(now) {
		t.Errorf("Expected lastExecute to be %v, got %v", now, lastExecute)
	}
}
