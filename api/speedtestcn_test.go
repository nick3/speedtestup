package api

import (
	"testing"
)

func TestNewSpeedTestCNClient(t *testing.T) {
	client := NewSpeedTestCNClient("192.168.1.1")
	if client == nil {
		t.Fatal("NewSpeedTestCNClient should not return nil")
	}

	if client.client == nil {
		t.Fatal("SpeedTestCNClient client should not be nil")
	}

	if client.bindIP != "192.168.1.1" {
		t.Errorf("Expected bindIP %s, got %s", "192.168.1.1", client.bindIP)
	}
}

func TestSpeedupQueryResponse_IsSpeedupAvailable(t *testing.T) {
	resp := &SpeedupQueryResponse{
		Data: struct {
			IP            string      `json:"ip"`
			UpdatedAt     string      `json:"updatedAt"`
			CanSpeed      int         `json:"canSpeed"`
			Download      int         `json:"download"`
			DownExpire    string      `json:"downExpire"`
			DownExpireT   interface{} `json:"downExpireT"`
			TargetUpH     int         `json:"targetUpH"`
			UpHExpire     string      `json:"upHExpire"`
			UpHExpireT    interface{} `json:"upHExpireT"`
			TargetUp100   int         `json:"targetUp100"`
			Up100Expire   string      `json:"up100Expire"`
			Up100ExpireT  interface{} `json:"up100ExpireT"`
			DownUp50Expire  string      `json:"downUp50Expire"`
			DownUp50ExpireT interface{} `json:"downUp50ExpireT"`
			DownUpExpire  string      `json:"downUpExpire"`
			DownUpExpireT interface{} `json:"downUpExpireT"`
		}{
			CanSpeed: 1,
		},
	}

	if !resp.IsSpeedupAvailable() {
		t.Error("Expected IsSpeedupAvailable to return true when CanSpeed is 1")
	}

	resp.Data.CanSpeed = 0
	if resp.IsSpeedupAvailable() {
		t.Error("Expected IsSpeedupAvailable to return false when CanSpeed is 0")
	}
}

func TestSpeedupQueryResponse_GetBandwidth(t *testing.T) {
	resp := &SpeedupQueryResponse{
		Data: struct {
			IP            string      `json:"ip"`
			UpdatedAt     string      `json:"updatedAt"`
			CanSpeed      int         `json:"canSpeed"`
			Download      int         `json:"download"`
			DownExpire    string      `json:"downExpire"`
			DownExpireT   interface{} `json:"downExpireT"`
			TargetUpH     int         `json:"targetUpH"`
			UpHExpire     string      `json:"upHExpire"`
			UpHExpireT    interface{} `json:"upHExpireT"`
			TargetUp100   int         `json:"targetUp100"`
			Up100Expire   string      `json:"up100Expire"`
			Up100ExpireT  interface{} `json:"up100ExpireT"`
			DownUp50Expire  string      `json:"downUp50Expire"`
			DownUp50ExpireT interface{} `json:"downUp50ExpireT"`
			DownUpExpire  string      `json:"downUpExpire"`
			DownUpExpireT interface{} `json:"downUpExpireT"`
		}{
			Download:     100,
			TargetUpH:    2048,
			TargetUp100:  5120,
		},
	}

	if resp.GetDownloadBandwidth() != 100 {
		t.Errorf("Expected GetDownloadBandwidth to return 100, got %d", resp.GetDownloadBandwidth())
	}

	if resp.GetUpHBandwidth() != 2 {
		t.Errorf("Expected GetUpHBandwidth to return 2, got %d", resp.GetUpHBandwidth())
	}

	if resp.GetUp100Bandwidth() != 5 {
		t.Errorf("Expected GetUp100Bandwidth to return 5, got %d", resp.GetUp100Bandwidth())
	}
}

func TestSpeedupReopenResponse(t *testing.T) {
	resp := &SpeedupReopenResponse{
		Code:    0,
		Message: "success",
		Data: struct {
			Result string `json:"result"`
		}{
			Result: "ok",
		},
	}

	if resp.Code != 0 || resp.Message != "success" || resp.Data.Result != "ok" {
		t.Error("SpeedupReopenResponse fields do not match expected values")
	}
}
