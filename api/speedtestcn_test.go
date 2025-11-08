package api

import (
	"testing"
	"time"
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
			CanSpeed int `json:"canSpeed"`
		}{CanSpeed: 1},
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
			Download   int    `json:"download"`
			TargetUpH  int    `json:"targetUpH"`
			TargetUp100 int   `json:"targetUp100"`
		}{
			Download:   100,
			TargetUpH:  2048,
			TargetUp100: 5120,
		},
	}

	if resp.GetDownloadBandwidth() != 100 {
		t.Errorf("Expected download bandwidth %d, got %d", 100, resp.GetDownloadBandwidth())
	}

	if resp.GetUpHBandwidth() != 2 { // 2048 / 1024
		t.Errorf("Expected upH bandwidth %d, got %d", 2, resp.GetUpHBandwidth())
	}

	if resp.GetUp100Bandwidth() != 5 { // 5120 / 1024
		t.Errorf("Expected up100 bandwidth %d, got %d", 5, resp.GetUp100Bandwidth())
	}
}
