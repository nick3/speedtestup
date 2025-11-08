package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"speedtestup/api"
	"speedtestup/config"
	"speedtestup/service"
	"speedtestup/utils"
)

// 版本信息变量（通过 LDFLAGS 设置）
var (
	version    = "2.0.0"
	buildDate  = "unknown"
	commitHash = "unknown"
)

func main() {
	// 解析命令行参数
	var configPath string
	var showVersion bool
	flag.StringVar(&configPath, "config", "config.json", "配置文件路径")
	flag.BoolVar(&showVersion, "version", false, "显示版本信息")
	flag.Parse()

	// 显示版本信息
	if showVersion {
		fmt.Printf("SpeedTestUp v%s\n", version)
		fmt.Printf("构建日期: %s\n", buildDate)
		fmt.Printf("提交哈希: %s\n", commitHash)
		fmt.Println("基于 luci-app-broadbandacc 架构重构的宽带提速工具")
		return
	}

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("❌ 加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 验证配置
	if !cfg.Speedup.Enabled {
		fmt.Println("❌ 提速服务未启用，请在 config.json 中设置 speedup.enabled = true")
		os.Exit(1)
	}

	// 初始化日志
	logger, err := utils.NewLogger(cfg.Logging.Level, cfg.Logging.Output, cfg.Logging.File)
	if err != nil {
		fmt.Printf("❌ 初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Info("🚀 SpeedTestUp 宽带提速服务启动")
	logger.Info("📋 配置信息:")
	logger.Info("  - 提速服务: %v", cfg.Speedup.Enabled)
	logger.Info("  - 下行提速: %v", cfg.Speedup.DownAcc)
	logger.Info("  - 上行提速: %v", cfg.Speedup.UpAcc)
	logger.Info("  - IP 绑定: %v", cfg.Speedup.IPBinding.Enabled)
	logger.Info("  - 自动恢复: %v", cfg.Speedup.AutoRecovery.Enabled)
	logger.Info("  - 7 天自检: %v", cfg.Speedup.SelfCheck.Enabled)
	logger.Info("  - 日志记录: %v", cfg.Speedup.Logging)
	logger.Info("  - 详细模式: %v", cfg.Speedup.Verbose)

	// 初始化 API 客户端
	ipAPI := api.NewIPAPI()
	speedupAPI := api.NewSpeedTestCNClient(cfg.Speedup.IPBinding.BindIP)

	// 初始化服务
	ipService := service.NewIPService(ipAPI, cfg)
	speedupService := service.NewSpeedupService(speedupAPI, cfg)
	scheduler := service.NewScheduler(ipService, speedupService, cfg)

	// 启动服务
	if err := scheduler.Start(); err != nil {
		logger.Error("❌ 启动服务失败: %v", err)
		os.Exit(1)
	}
	logger.Info("✅ 服务启动成功")

	// 执行首次提速检查
	logger.Info("🔍 执行首次提速检查...")
	if err := speedupService.Execute(); err != nil {
		logger.Warn("⚠️  首次提速检查失败: %v", err)
	} else {
		logger.Info("✅ 首次提速检查完成")
	}

	// 等待退出信号
	waitForShutdown(logger, scheduler)
}

// waitForShutdown 等待退出信号并优雅关闭
func waitForShutdown(logger *utils.Logger, scheduler *service.Scheduler) {
	// 创建信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	sig := <-sigChan
	logger.Info("📴 收到信号 %v，正在优雅关闭...", sig)

	// 关闭调度器
	if err := scheduler.Stop(); err != nil {
		logger.Error("❌ 关闭服务失败: %v", err)
		os.Exit(1)
	}

	logger.Info("✅ 服务已优雅关闭")
}
