package service

import (
	"fmt"
	"net"

	"speedtestup/api"
	"speedtestup/config"
	"speedtestup/utils"
)

// IPService IP 服务
type IPService struct {
	apiClient *api.IPAPI
	config    *config.IPBindingConfig
	logger    *utils.Logger
	lastIP    string
}

// NewIPService 创建新的 IP 服务实例
func NewIPService(ipAPI *api.IPAPI, cfg *config.Config) *IPService {
	logger, _ := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Output, cfg.Logging.File)
	logger = logger.WithPrefix("IPService")

	return &IPService{
		apiClient: ipAPI,
		config:    &cfg.Speedup.IPBinding,
		logger:    logger,
		lastIP:    "",
	}
}

// GetCurrentIP 获取当前公网 IP
func (s *IPService) GetCurrentIP() (string, error) {
	ip, err := s.apiClient.GetPublicIP()
	if err != nil {
		s.logger.Error("获取当前公网 IP 失败: %v", err)
		return "", err
	}
	s.logger.Info("获取当前公网 IP: %s", ip)
	return ip, nil
}

// ValidateBinding 验证 IP 绑定
func (s *IPService) ValidateBinding(ip string) error {
	if !s.config.Enabled {
		s.logger.Debug("IP 绑定未启用，跳过验证")
		return nil
	}

	if s.config.BindIP != "" && ip != s.config.BindIP {
		err := fmt.Errorf("IP 绑定验证失败: 当前 IP %s != 绑定 IP %s", ip, s.config.BindIP)
		s.logger.Error("%v", err)
		return err
	}

	s.logger.Debug("IP 绑定验证通过")
	return nil
}

// CheckIPChange 检查 IP 是否发生变化
func (s *IPService) CheckIPChange() (bool, error) {
	currentIP, err := s.GetCurrentIP()
	if err != nil {
		return false, err
	}

	// 首次获取 IP
	if s.lastIP == "" {
		s.lastIP = currentIP
		s.logger.Info("初始化 IP: %s", currentIP)
		return false, nil
	}

	// 检查 IP 是否变化
	if currentIP != s.lastIP {
		s.logger.Info("检测到 IP 变化: %s -> %s", s.lastIP, currentIP)
		s.lastIP = currentIP
		return true, nil
	}

	return false, nil
}

// GetInterfaceIP 获取指定网络接口的 IP 地址
// 对应 luci-app-broadbandacc 中的 get_bind_ip 函数
func (s *IPService) GetInterfaceIP(interfaceName string) (string, error) {
	// 在实际实现中，这里需要与系统网络接口交互
	// 由于这是跨平台的 Go 实现，我们提供一个简化的版本
	s.logger.Debug("尝试获取接口 %s 的 IP 地址", interfaceName)

	// 这里可以根据不同操作系统实现不同的接口获取逻辑
	// Linux: 使用 netlink 或 sysfs
	// Windows: 使用 Windows API
	// macOS: 使用系统调用

	// 暂时返回一个示例实现
	interfaces, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("获取网络接口列表失败: %v", err)
	}

	for _, addr := range interfaces {
		ip, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		// 过滤回环地址和 IPv6
		if ip.IP.IsLoopback() || ip.IP.To4() == nil {
			continue
		}

		// 这里可以根据接口名称过滤
		// 在实际实现中，需要更复杂的逻辑来匹配接口名称
		return ip.IP.String(), nil
	}

	return "", fmt.Errorf("未找到有效的网络接口 IP")
}

// ResetIP 重置 IP 记录（用于测试或特殊情况）
func (s *IPService) ResetIP() {
	s.lastIP = ""
	s.logger.Info("IP 记录已重置")
}
