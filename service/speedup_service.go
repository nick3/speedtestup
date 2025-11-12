package service

import (
	"fmt"
	"time"

	"speedtestup/api"
	"speedtestup/config"
	"speedtestup/utils"
)

// SpeedupService 提速服务
type SpeedupService struct {
	apiClient   *api.SpeedTestCNClient
	config      *config.AutoRecoveryConfig
	selfCheck   *config.SelfCheckConfig
	logger      *utils.Logger
	lastExecute time.Time
}

// NewSpeedupService 创建新的提速服务实例
func NewSpeedupService(speedTestCNClient *api.SpeedTestCNClient, cfg *config.Config) *SpeedupService {
	logger, err := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Output, cfg.Logging.File)
	if err != nil {
		// 无法初始化 logger 是一个严重问题，至少需要 panic 或返回错误
		fmt.Printf("Failed to initialize logger for SpeedupService: %v\n", err)
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	logger = logger.WithPrefix("SpeedupService")

	return &SpeedupService{
		apiClient:   speedTestCNClient,
		config:      &cfg.Speedup.AutoRecovery,
		selfCheck:   &cfg.Speedup.SelfCheck,
		logger:      logger,
		lastExecute: time.Time{},
	}
}

// Execute 执行提速（带自动恢复）
// 对应 luci-app-broadbandacc 中的 isp_bandwidth 函数
func (s *SpeedupService) Execute() error {
	s.logger.Info("开始执行提速操作...")

	// 1. 先重新开启提速
	s.logger.Debug("调用重新开启提速接口...")
	reopenResp, err := s.apiClient.ReopenSpeedup()
	if err != nil {
		s.logger.Error("重新开启提速失败: %v", err)
		return s.handleError(err, "重新开启提速")
	}

	// 检查重新开启提速的响应
	if reopenResp.Code != 0 {
		switch reopenResp.Code {
		case 10021:
			s.logger.Error("请求接口异常，请重启插件再试")
			return fmt.Errorf("接口异常，错误码: %d", reopenResp.Code)
		case 10002:
			s.logger.Warn("操作过于频繁，接口提速已受理")
		default:
			s.logger.Warn("重新开启提速返回错误码: %d, 消息: %s", reopenResp.Code, reopenResp.Message)
		}
	} else {
		s.logger.Info("重新开启提速接口连接正常")
	}

	// 2. 查询提速状态
	s.logger.Debug("查询提速状态...")
	queryResp, err := s.apiClient.QuerySpeedupStatus()
	if err != nil {
		s.logger.Error("查询提速状态失败: %v", err)
		return s.handleError(err, "查询提速状态")
	}

	// 3. 解析提速信息并输出
	s.parseAndLogSpeedupInfo(queryResp)

	// 4. 检查提速是否成功
	if !queryResp.IsSpeedupAvailable() {
		// 如果CanSpeed为0，检查是否有有效的带宽数据
		// 这可能表示已经处于提速状态
		hasBandwidth := queryResp.Data.Download > 0 || queryResp.Data.TargetUpH > 0 || queryResp.Data.TargetUp100 > 0

		if hasBandwidth {
			s.logger.Warn("当前已处于提速状态（CanSpeed=0但检测到带宽数据）")
			s.logger.Info("可能原因：当前线路已提速，或接口返回CanSpeed=0表示无需重复提速")
		} else {
			s.logger.Error("网络不支持提速")
			return fmt.Errorf("网络不支持提速")
		}
	}

	// 5. 检查上行和下行提速状态
	upActive, err := queryResp.IsUpSpeedupActive()
	if err != nil {
		s.logger.Error("检查上行提速状态失败: %v", err)
		return err
	}

	downActive, err := queryResp.IsDownloadSpeedupActive()
	if err != nil {
		s.logger.Error("检查下行提速状态失败: %v", err)
		return err
	}

	// 6. 输出提速结果
	if upActive {
		s.logger.Success("上行提速已激活")
	} else {
		s.logger.Warn("上行提速未激活")
	}

	if downActive {
		s.logger.Success("下行提速已激活")
	} else {
		s.logger.Warn("下行提速未激活")
	}

	s.lastExecute = time.Now()
	return nil
}

// handleError 处理错误（带自动恢复）
func (s *SpeedupService) handleError(err error, operation string) error {
	if !s.config.Enabled {
		s.logger.Error("%s失败，自动恢复未启用: %v", operation, err)
		return err
	}

	s.logger.Warn("%s失败，开始自动恢复流程 (最大重试次数: %d)", operation, s.config.MaxRetries)

	for i := 1; i <= s.config.MaxRetries; i++ {
		s.logger.Info("自动恢复尝试 %d/%d，等待 %v", i, s.config.MaxRetries, s.config.RetryInterval)
		time.Sleep(s.config.RetryInterval)

		if err := s.Execute(); err == nil {
			s.logger.Success("自动恢复成功")
			return nil
		}
	}

	s.logger.Error("自动恢复失败，已达到最大重试次数: %v", err)
	return fmt.Errorf("自动恢复失败: %v", err)
}

// parseAndLogSpeedupInfo 解析并记录提速信息
func (s *SpeedupService) parseAndLogSpeedupInfo(resp *api.SpeedupQueryResponse) {
	s.logger.Info("提速开始时间: %s", resp.Data.UpdatedAt)
	s.logger.Info("出口IP地址: %s", resp.Data.IP)

	// 上行带宽信息
	if resp.Data.TargetUpH > 0 {
		s.logger.Info("一类上行带宽%dM提速截至时间: %s", resp.GetUpHBandwidth(), resp.Data.UpHExpire)
	}
	if resp.Data.TargetUp100 > 0 {
		s.logger.Info("二类上行带宽%dM提速截至时间: %s", resp.GetUp100Bandwidth(), resp.Data.Up100Expire)
	}

	// 下行带宽信息
	if resp.Data.Download > 0 {
		s.logger.Info("下行带宽%dM提速截至时间: %s", resp.Data.Download, resp.Data.DownExpire)
	}

	// 套餐信息
	if resp.Data.DownUp50Expire != "" {
		s.logger.Info("一类套餐带宽%dM上行+%dM下行提速截至时间: %s",
			resp.GetUpHBandwidth(), resp.Data.Download, resp.Data.DownUp50Expire)
	}
	if resp.Data.DownUpExpire != "" {
		s.logger.Info("二类套餐带宽%dM上行+%dM下行提速截至时间: %s",
			resp.GetUp100Bandwidth(), resp.Data.Download, resp.Data.DownUpExpire)
	}
}

// QueryStatus 查询提速状态
func (s *SpeedupService) QueryStatus() (bool, error) {
	resp, err := s.apiClient.QuerySpeedupStatus()
	if err != nil {
		return false, err
	}
	return resp.IsSpeedupAvailable(), nil
}

// GetLastExecuteTime 获取上次执行时间
func (s *SpeedupService) GetLastExecuteTime() time.Time {
	return s.lastExecute
}

// ShouldSelfCheck 检查是否应该执行自检
func (s *SpeedupService) ShouldSelfCheck() bool {
	if !s.selfCheck.Enabled {
		return false
	}

	if s.lastExecute.IsZero() {
		return false
	}

	return time.Since(s.lastExecute) >= s.selfCheck.Interval
}

// ExecuteSelfCheck 执行自检
func (s *SpeedupService) ExecuteSelfCheck() error {
	s.logger.Info("开始执行 7 天自检...")
	if err := s.Execute(); err != nil {
		s.logger.Error("7 天自检失败: %v", err)
		return err
	}
	s.logger.Success("7 天自检完成")
	return nil
}
