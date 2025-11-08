package api

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-resty/resty/v2"
)

// 缓存正则表达式，避免重复编译
var (
	// 匹配空白字符的正则表达式
	whitespaceRegex = regexp.MustCompile(`\s+`)
	// 验证 IPv4 地址的正则表达式
	ipRegex = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
)

// IPAPI IP 查询 API
type IPAPI struct {
	client *resty.Client
}

// NewIPAPI 创建新的 IP API 实例
func NewIPAPI() *IPAPI {
	return &IPAPI{
		client: resty.New().
			SetTimeout(10 * time.Second).
			SetHeader("User-Agent", "SpeedTestUp/1.0"),
	}
}

// GetPublicIP 获取公网 IP（使用 ipinfo.io）
func (a *IPAPI) GetPublicIP() (string, error) {
	// 根据 luci-app-broadbandacc，使用 ipinfo.io/ip/ 获取公网 IP
	resp, err := a.client.R().
		Get("https://ipinfo.io/ip/")
	if err != nil {
		return "", fmt.Errorf("获取公网 IP 失败: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("获取公网 IP 失败，状态码: %d", resp.StatusCode())
	}

	ip := string(resp.Body())
	// 清理 IP 字符串（移除可能的换行符和空格）
	ip = whitespaceRegex.ReplaceAllString(ip, "")

	// 验证 IP 格式
	if !ipRegex.MatchString(ip) {
		return "", fmt.Errorf("获取的 IP 格式无效: %s", ip)
	}

	return ip, nil
}
