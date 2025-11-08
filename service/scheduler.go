package service

import (
	"fmt"
	"sync"
	"time"

	"speedtestup/config"
	"speedtestup/utils"

	"github.com/robfig/cron/v3"
)

// Scheduler 调度服务
type Scheduler struct {
	cron          *cron.Cron
	ipService     *IPService
	speedupService *SpeedupService
	config        *config.SpeedupConfig
	logger        *utils.Logger
	lastIP        string
	running       bool
	mu            sync.Mutex
}

// NewScheduler 创建新的调度器实例
func NewScheduler(ipService *IPService, speedupService *SpeedupService, cfg *config.Config) *Scheduler {
	logger, err := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Output, cfg.Logging.File)
	if err != nil {
		// 无法初始化 logger 是一个严重问题，至少需要 panic 或返回错误
		fmt.Printf("Failed to initialize logger for Scheduler: %v\n", err)
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	logger = logger.WithPrefix("Scheduler")

	return &Scheduler{
		cron:           cron.New(),
		ipService:      ipService,
		speedupService: speedupService,
		config:         &cfg.Speedup,
		logger:         logger,
		lastIP:         "",
		running:        false,
	}
}

// Start 启动调度器
// 对应 luci-app-broadbandacc 中的 main 函数逻辑
func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		s.logger.Warn("调度器已在运行")
		return nil
	}

	s.logger.Info("启动调度器...")

	// 1. 启动心跳检测（对应 _keepalive 函数）
	s.startHeartbeat()

	// 2. 启动 7 天自检（对应 Weekly_cycle 函数）
	s.startSelfCheck()

	// 3. 启动重新开启提速的定时任务（对应每周一 0:0 的任务）
	s.startReopenSchedule()

	// 4. 启动 cron 调度器
	s.cron.Start()

	s.running = true
	s.logger.Success("调度器启动成功")

	// 5. 执行首次提速
	s.logger.Info("执行首次提速...")
	if err := s.speedupService.Execute(); err != nil {
		s.logger.Error("首次提速失败: %v", err)
	} else {
		s.logger.Success("首次提速成功")
	}

	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		s.logger.Warn("调度器未在运行")
		return nil
	}

	s.logger.Info("停止调度器...")
	s.cron.Stop()
	s.running = false
	s.logger.Success("调度器已停止")
	return nil
}

// startHeartbeat 启动心跳检测
// 对应 luci-app-broadbandacc 中的 _keepalive 函数
func (s *Scheduler) startHeartbeat() {
	// 将检查间隔转换为 cron 表达式
	// 例如: 10m -> "*/10 * * * *"
	interval := s.config.CheckInterval
	if interval < time.Minute {
		interval = time.Minute
	}

	cronExpr := fmt.Sprintf("*/%d * * * *", int(interval.Minutes()))
	s.logger.Debug("配置心跳检测间隔: %v (Cron: %s)", interval, cronExpr)

	// 添加心跳检测任务
	_, err := s.cron.AddFunc(cronExpr, s.heartbeatCheck)
	if err != nil {
		s.logger.Error("添加心跳检测任务失败: %v", err)
		return
	}

	s.logger.Debug("心跳检测任务已添加")
}

// startSelfCheck 启动 7 天自检
// 对应 luci-app-broadbandacc 中的 Weekly_cycle 函数
func (s *Scheduler) startSelfCheck() {
	if !s.config.SelfCheck.Enabled {
		s.logger.Debug("7 天自检未启用，跳过")
		return
	}

	// 使用 cron 表达式 "0 0 * * 1" 表示每周一 0:0
	cronExpr := "0 0 * * 1"
	s.logger.Debug("配置 7 天自检 (Cron: %s)", cronExpr)

	_, err := s.cron.AddFunc(cronExpr, s.selfCheckTask)
	if err != nil {
		s.logger.Error("添加 7 天自检任务失败: %v", err)
		return
	}

	s.logger.Debug("7 天自检任务已添加")
}

// startReopenSchedule 启动重新开启提速的定时任务
func (s *Scheduler) startReopenSchedule() {
	cronExpr := s.config.ReopenSchedule
	if cronExpr == "" {
		cronExpr = "0 0 * * 1" // 默认每周一 0:0
	}

	s.logger.Debug("配置重新开启提速任务 (Cron: %s)", cronExpr)

	_, err := s.cron.AddFunc(cronExpr, s.reopenSpeedupTask)
	if err != nil {
		s.logger.Error("添加重新开启提速任务失败: %v", err)
		return
	}

	s.logger.Debug("重新开启提速任务已添加")
}

// heartbeatCheck 心跳检测任务
// 对应 luci-app-broadbandacc 中的 _keepalive 函数
func (s *Scheduler) heartbeatCheck() {
	s.logger.Debug("开始心跳检测...")

	// 1. 检查 IP 是否变化
	ipChanged, err := s.ipService.CheckIPChange()
	if err != nil {
		s.logger.Error("心跳检测失败: %v", err)
		return
	}

	// 2. 如果 IP 发生变化，重新执行提速
	if ipChanged {
		s.logger.Info("IP 发生变化，重新执行提速...")
		if err := s.speedupService.Execute(); err != nil {
			s.logger.Error("IP 变化后提速失败: %v", err)
		} else {
			s.logger.Success("IP 变化后提速成功")
		}
		return
	}

	// 3. 如果设置了 IP 绑定，验证绑定状态
	if s.config.IPBinding.Enabled {
		currentIP, err := s.ipService.GetCurrentIP()
		if err != nil {
			s.logger.Error("获取当前 IP 失败: %v", err)
			return
		}

		if err := s.ipService.ValidateBinding(currentIP); err != nil {
			s.logger.Warn("IP 绑定验证失败: %v", err)
			// 可以选择重新执行提速或记录警告
		}
	}

	s.logger.Debug("心跳检测完成")
}

// selfCheckTask 7 天自检任务
// 对应 luci-app-broadbandacc 中的 Weekly_cycle 函数
func (s *Scheduler) selfCheckTask() {
	s.logger.Info("执行 7 天自检任务...")

	if err := s.speedupService.ExecuteSelfCheck(); err != nil {
		s.logger.Error("7 天自检失败: %v", err)
	} else {
		s.logger.Success("7 天自检完成")
	}
}

// reopenSpeedupTask 重新开启提速任务
func (s *Scheduler) reopenSpeedupTask() {
	s.logger.Info("执行重新开启提速任务...")

	if err := s.speedupService.Execute(); err != nil {
		s.logger.Error("重新开启提速失败: %v", err)
	} else {
		s.logger.Success("重新开启提速成功")
	}
}

// IsRunning 检查调度器是否在运行
func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// GetStatus 获取调度器状态
func (s *Scheduler) GetStatus() map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	return map[string]interface{}{
		"running":       s.running,
		"last_execute":  s.speedupService.GetLastExecuteTime(),
		"check_interval": s.config.CheckInterval.String(),
		"self_check":    s.config.SelfCheck.Enabled,
		"auto_recovery": s.config.AutoRecovery.Enabled,
	}
}
