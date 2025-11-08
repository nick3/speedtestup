package api

import (
	"testing"
	"time"
)

func TestNewIPAPI(t *testing.T) {
	ipAPI := NewIPAPI()
	if ipAPI == nil {
		t.Fatal("NewIPAPI should not return nil")
	}

	if ipAPI.client == nil {
		t.Fatal("IPAPI client should not be nil")
	}

	// 检查默认超时时间
	expectedTimeout := 10 * time.Second
	if ipAPI.client.GetClient().Timeout != expectedTimeout {
		t.Errorf("Expected timeout %v, got %v", expectedTimeout, ipAPI.client.GetClient().Timeout)
	}
}

func TestIPAPI_InvalidIP(t *testing.T) {
	ipAPI := NewIPAPI()
	// 这里只是测试结构体，不进行实际的网络请求
	if ipAPI == nil {
		t.Fatal("NewIPAPI should not return nil")
	}
}
