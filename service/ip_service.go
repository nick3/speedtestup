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
	logger, err := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Output, cfg.Logging.File)
	if err != nil {
		// 至少在标准错误输出中打印一条日志
		fmt.Printf("Failed to initialize logger for IPService: %v\n", err)
		// 使用默认日志器或panic，根据业务逻辑决定
		// 这里选择panic，因为日志器初始化失败是严重问题
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
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
	s.logger.Debug("尝试获取接口 %s 的 IP 地址", interfaceName)

	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("获取网络接口列表失败: %v", err)
	}

	// 遍历网络接口
	for _, iface := range interfaces {
		// 如果指定了接口名称，匹配接口名称
		if interfaceName != "" && iface.Name != interfaceName {
			continue
		}

		// 获取接口的地址
		addrs, err := iface.Addrs()
		if err != nil {
			s.logger.Debug("获取接口 %s 地址失败: %v", iface.Name, err)
			continue
		}

		// 遍历接口的地址
		for _, addr := range addrs {
			ip, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			// 过滤回环地址和 IPv6
			if ip.IP.IsLoopback() || ip.IP.To4() == nil {
				continue
			}

			s.logger.Debug("找到接口 %s 的 IP 地址: %s", iface.Name, ip.IP.String())
			return ip.IP.String(), nil
		}
	}

	if interfaceName != "" {
		return "", fmt.Errorf("未找到接口 %s 的有效网络 IP", interfaceName)
	}
	return "", fmt.Errorf("未找到任何有效的网络接口 IP")
}

// ResetIP 重置 IP 记录（用于测试或特殊情况）
func (s *IPService) ResetIP() {
	s.lastIP = ""
	s.logger.Info("IP 记录已重置")
}
