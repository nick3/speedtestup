package service

import (
	"net"
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

// TestIPService_ValidateBinding 测试IP绑定验证
func TestIPService_ValidateBinding(t *testing.T) {
	cfg := config.NewDefaultConfig()
	ipAPI := api.NewIPAPI()
	ipService := NewIPService(ipAPI, cfg)

	// 测试IP绑定未启用的情况
	cfg.Speedup.IPBinding.Enabled = false
	err := ipService.ValidateBinding("192.168.1.1")
	if err != nil {
		t.Errorf("Expected no error when IP binding is disabled, got: %v", err)
	}

	// 测试IP绑定启用但未设置绑定IP的情况
	cfg.Speedup.IPBinding.Enabled = true
	cfg.Speedup.IPBinding.BindIP = ""
	err = ipService.ValidateBinding("192.168.1.1")
	if err != nil {
		t.Errorf("Expected no error when bind IP is not set, got: %v", err)
	}

	// 测试IP匹配的情况
	cfg.Speedup.IPBinding.Enabled = true
	cfg.Speedup.IPBinding.BindIP = "192.168.1.100"
	err = ipService.ValidateBinding("192.168.1.100")
	if err != nil {
		t.Errorf("Expected no error when IPs match, got: %v", err)
	}

	// 测试IP不匹配的情况
	cfg.Speedup.IPBinding.BindIP = "192.168.1.100"
	err = ipService.ValidateBinding("192.168.1.200")
	if err == nil {
		t.Error("Expected error when IPs don't match")
	}
}

// TestIPService_GetInterfaceIP 测试获取指定网络接口的IP
func TestIPService_GetInterfaceIP(t *testing.T) {
	cfg := config.NewDefaultConfig()
	ipAPI := api.NewIPAPI()
	ipService := NewIPService(ipAPI, cfg)

	// 测试空接口名（应该返回第一个有效IP）
	// 注意：在测试环境中，这取决于系统的网络接口
	// 我们主要测试函数不会panic
	ip, err := ipService.GetInterfaceIP("")
	if err == nil {
		t.Logf("Successfully retrieved IP for empty interface name: %s", ip)
	} else {
		t.Logf("Expected error for empty interface name (no valid interfaces): %v", err)
	}

	// 测试不存在的接口名
	_, err = ipService.GetInterfaceIP("nonexistent_interface_12345")
	if err == nil {
		t.Error("Expected error for nonexistent interface")
	} else {
		t.Logf("Got expected error for nonexistent interface: %v", err)
	}

	// 测试存在的接口（如果有的话）
	// 枚举系统接口并尝试获取第一个可用接口的IP
	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		// 跳过回环接口
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 尝试获取接口IP
		ip, err := ipService.GetInterfaceIP(iface.Name)
		if err == nil {
			t.Logf("Successfully retrieved IP for interface %s: %s", iface.Name, ip)
			break
		}
	}
}
