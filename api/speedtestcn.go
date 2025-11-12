package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// SpeedTestCNClient speedtest.cn API 客户端
type SpeedTestCNClient struct {
	client  *resty.Client
	bindIP  string // 绑定的 IP 地址
}

// NewSpeedTestCNClient 创建新的 speedtest.cn API 客户端
func NewSpeedTestCNClient(bindIP string) *SpeedTestCNClient {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetHeader("User-Agent", "SpeedTestUp/1.0")

	// 如果设置了绑定 IP，创建自定义的 Transport
	if bindIP != "" {
		// 解析绑定的 IP 地址
		localAddr, err := net.ResolveTCPAddr("tcp", bindIP+":0")
		if err != nil {
			// 如果解析失败，使用默认配置
			fmt.Printf("Warning: Failed to parse bind IP %s: %v\n", bindIP, err)
		} else {
			// 创建自定义的 Dialer
			dialer := &net.Dialer{
				LocalAddr: localAddr,
			}

			// 设置自定义的 Transport
			transport := &http.Transport{
				DialContext: dialer.DialContext,
			}

			client.SetTransport(transport)
		}
	}

	return &SpeedTestCNClient{
		client: client,
		bindIP: bindIP,
	}
}

// QuerySpeedupStatus 查询提速状态
// 对应 luci-app-broadbandacc 中的 $_http_cmd
func (c *SpeedTestCNClient) QuerySpeedupStatus() (*SpeedupQueryResponse, error) {
	url := "https://tisu-api-v3.speedtest.cn/speedUp/query"

	req := c.client.R().
		SetHeader("Content-Type", "application/json")

	resp, err := req.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求提速查询接口失败: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("提速查询接口返回错误，状态码: %d", resp.StatusCode())
	}

	var data SpeedupQueryResponse
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, fmt.Errorf("解析提速查询响应失败: %v", err)
	}

	return &data, nil
}

// ReopenSpeedup 重新开启提速
// 对应 luci-app-broadbandacc 中的 $_http_cmd2
func (c *SpeedTestCNClient) ReopenSpeedup() (*SpeedupReopenResponse, error) {
	url := "https://tisu-api.speedtest.cn/api/v2/speedup/reopen"

	req := c.client.R().
		SetHeader("Content-Type", "application/json")

	resp, err := req.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求重新开启提速接口失败: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("重新开启提速接口返回错误，状态码: %d", resp.StatusCode())
	}

	var data SpeedupReopenResponse
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, fmt.Errorf("解析重新开启提速响应失败: %v", err)
	}

	return &data, nil
}

// SpeedupQueryResponse 提速查询响应结构
type SpeedupQueryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		IP          string `json:"ip"`          // 出口 IP 地址
		UpdatedAt   string `json:"updatedAt"`   // 提速开始时间
		CanSpeed    int    `json:"canSpeed"`    // 是否可以提速 (0: 不支持, 1: 支持)
		Download    int    `json:"download"`    // 下行带宽 (Mbps)
		DownExpire  string `json:"downExpire"`  // 下行提速截止时间
		DownExpireT json.Number `json:"downExpireT"` // 下行提速截止时间 (时间戳，可能是数字或字符串)

		// 上行带宽信息
		TargetUpH     int    `json:"targetUpH"`     // 一类上行带宽 (Kbps)
		UpHExpire     string `json:"upHExpire"`     // 一类上行带宽提速截止时间
		UpHExpireT    json.Number `json:"upHExpireT"`    // 一类上行带宽提速截止时间 (时间戳)

		TargetUp100   int    `json:"targetUp100"`   // 二类上行带宽 (Kbps)
		Up100Expire   string `json:"up100Expire"`   // 二类上行带宽提速截止时间
		Up100ExpireT  json.Number `json:"up100ExpireT"`  // 二类上行带宽提速截止时间 (时间戳)

		// 套餐信息
		DownUp50Expire  string `json:"downUp50Expire"`  // 一类套餐带宽上行+下行提速截止时间
		DownUp50ExpireT json.Number `json:"downUp50ExpireT"` // 一类套餐带宽上行+下行提速截止时间 (时间戳)

		DownUpExpire  string `json:"downUpExpire"`  // 二类套餐带宽上行+下行提速截止时间
		DownUpExpireT json.Number `json:"downUpExpireT"` // 二类套餐带宽上行+下行提速截止时间 (时间戳)
	} `json:"data"`
}

// SpeedupReopenResponse 重新开启提速响应结构
type SpeedupReopenResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Result string `json:"result"`
	} `json:"data"`
}

// IsSpeedupAvailable 检查是否可以提速
func (r *SpeedupQueryResponse) IsSpeedupAvailable() bool {
	return r.Data.CanSpeed == 1
}

// GetDownloadBandwidth 获取下行带宽 (Mbps)
func (r *SpeedupQueryResponse) GetDownloadBandwidth() int {
	return r.Data.Download
}

// GetUpHBandwidth 获取一类上行带宽 (Mbps)
func (r *SpeedupQueryResponse) GetUpHBandwidth() int {
	return r.Data.TargetUpH / 1024
}

// GetUp100Bandwidth 获取二类上行带宽 (Mbps)
func (r *SpeedupQueryResponse) GetUp100Bandwidth() int {
	return r.Data.TargetUp100 / 1024
}

// IsDownloadSpeedupActive 检查下行提速是否激活
func (r *SpeedupQueryResponse) IsDownloadSpeedupActive() (bool, error) {
	// 检查下行提速
	if r.Data.DownExpireT != "false" {
		expireTime, err := parseTimestamp(r.Data.DownExpireT)
		if err != nil {
			return false, err
		}
		if time.Now().Before(expireTime) {
			return true, nil
		}
	}

	// 检查一类套餐下行提速
	if r.Data.DownUp50ExpireT != "false" {
		expireTime, err := parseTimestamp(r.Data.DownUp50ExpireT)
		if err != nil {
			return false, err
		}
		if time.Now().Before(expireTime) {
			return true, nil
		}
	}

	// 检查二类套餐下行提速
	if r.Data.DownUpExpireT != "false" {
		expireTime, err := parseTimestamp(r.Data.DownUpExpireT)
		if err != nil {
			return false, err
		}
		if time.Now().Before(expireTime) {
			return true, nil
		}
	}

	return false, nil
}

// IsUpSpeedupActive 检查上行提速是否激活
func (r *SpeedupQueryResponse) IsUpSpeedupActive() (bool, error) {
	// 检查一类上行提速
	if r.Data.UpHExpireT != "false" {
		expireTime, err := parseTimestamp(r.Data.UpHExpireT)
		if err != nil {
			return false, err
		}
		if time.Now().Before(expireTime) {
			return true, nil
		}
	}

	// 检查二类上行提速
	if r.Data.Up100ExpireT != "false" {
		expireTime, err := parseTimestamp(r.Data.Up100ExpireT)
		if err != nil {
			return false, err
		}
		if time.Now().Before(expireTime) {
			return true, nil
		}
	}

	// 检查一类套餐上行提速
	if r.Data.DownUp50ExpireT != "false" {
		expireTime, err := parseTimestamp(r.Data.DownUp50ExpireT)
		if err != nil {
			return false, err
		}
		if time.Now().Before(expireTime) {
			return true, nil
		}
	}

	// 检查二类套餐上行提速
	if r.Data.DownUpExpireT != "false" {
		expireTime, err := parseTimestamp(r.Data.DownUpExpireT)
		if err != nil {
			return false, err
		}
		if time.Now().Before(expireTime) {
			return true, nil
		}
	}

	return false, nil
}

// parseTimestamp 解析时间戳
// 使用上海时区（speedtest.cn 服务器时区）以确保时间判断准确
func parseTimestamp(timestampNum json.Number) (time.Time, error) {
	// 先将 json.Number 转换为字符串
	timestampStr := timestampNum.String()

	// 如果值为 "false"，表示未开通提速
	if timestampStr == "false" {
		return time.Time{}, nil // 返回零值时间
	}

	// 加载上海时区
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 如果加载失败，使用本地时区
		location = time.Local
	}

	// 使用指定时区解析时间
	timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", timestampStr, location)
	if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}

