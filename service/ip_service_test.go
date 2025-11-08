package service

import (
	"testing"

	"speedtestup/api"
	"speedtestup/config"
)

func TestNewIPService(t *testing.T) {
	cfg := config.NewDefaultConfig()
	ipAPI := api.NewIPAPI()
	ipService := NewIPService(ipAPI, cfg)

	if ipService == nil {
		t.Fatal("NewIPService should not return nil")
	}

	if ipService.apiClient == nil {
		t.Fatal("IPService apiClient should not be nil")
	}

	if ipService.config == nil {
		t.Fatal("IPService config should not be nil")
	}
}

func TestIPService_ResetIP(t *testing.T) {
	cfg := config.NewDefaultConfig()
	ipAPI := api.NewIPAPI()
	ipService := NewIPService(ipAPI, cfg)

	// 设置一个测试 IP
	ipService.lastIP = "192.168.1.1"
	if ipService.lastIP == "" {
		t.Error("Expected lastIP to be set")
	}

	// 重置 IP
	ipService.ResetIP()
	if ipService.lastIP != "" {
		t.Error("Expected lastIP to be empty after ResetIP")
	}
}
