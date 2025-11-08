# SpeedTestUp v2.0

[![Build and Release](https://github.com/nick3/speedtestup/actions/workflows/release.yml/badge.svg)](https://github.com/nick3/speedtestup/actions/workflows/release.yml)
[![Docker Build and Push](https://github.com/nick3/speedtestup/actions/workflows/docker.yml/badge.svg)](https://github.com/nick3/speedtestup/actions/workflows/docker.yml)

SpeedTestUp v2.0 是一个基于 Go 语言开发的自动化宽带提速工具，完全兼容 OpenWrt 插件 `luci-app-broadbandacc` 的 API 和功能。该项目能够定期检测网络状态、监控 IP 变化，并在需要时自动执行提速操作。

## 特性

- ✅ **API 兼容性**：与 luci-app-broadbandacc 使用相同的接口
- ✅ **自动提速**：自动检测网络状态并触发提速操作
- ✅ **IP 监控**：实时监控公网 IP 变化
- ✅ **自动恢复**：网络异常时自动重试和恢复
- ✅ **7 天自检**：定期自检和修复提速状态
- ✅ **详细日志**：多级别日志输出，便于问题排查
- ✅ **跨平台**：支持 Linux、macOS、Windows
- ✅ **容器化**：支持 Docker 部署

## 与 luci-app-broadbandacc 的兼容性

| 功能 | luci-app-broadbandacc | SpeedTestUp v2.0 | 状态 |
|------|----------------------|------------------|------|
| IP 查询 | `ipinfo.io/ip/` | ✅ 相同 | 完全兼容 |
| 提速查询 | `speedtest.cn/speedUp/query` | ✅ 相同 | 完全兼容 |
| 重新开启提速 | `speedtest.cn/speedup/reopen` | ✅ 相同 | 完全兼容 |
| IP 绑定 | `--bind-address` | ✅ 支持 | 完全兼容 |
| 心跳检测 | 每 5 秒 | 每 10 分钟 | 优化 |
| 7 天自检 | `sleep 7d` | ✅ 相同 | 完全兼容 |
| 自动恢复 | `_start_Strategy` | ✅ 相同 | 完全兼容 |

## 快速开始

### 二进制文件运行

1. 从 [Releases](https://github.com/nick3/speedtestup/releases) 页面下载适合你系统的最新版本
2. 解压并运行程序：

```bash
# Linux/MacOS
chmod +x speedup_*
./speedup_*

# Windows
speedup-windows.exe
```

**⚠️ 重要**: 程序启动时会检查配置文件。如果提速服务未启用，程序会提示您需要设置 `speedup.enabled = true`。

### Docker 运行

```bash
# 从 GitHub Container Registry 拉取镜像
docker pull ghcr.io/nick3/speedtestup:latest

# 运行容器（使用默认配置）
docker run -d --name speedtestup ghcr.io/nick3/speedtestup:latest

# 运行容器（挂载自定义配置）
docker run -d \
  --name speedtestup \
  -v /path/to/config.json:/app/config.json \
  ghcr.io/nick3/speedtestup:latest
```

### 配置

项目使用 `config.json` 文件进行配置。默认配置示例：

```json
{
  "speedup": {
    "enabled": true,
    "down_acc": true,
    "up_acc": true,
    "check_interval": "10m",
    "reopen_schedule": "0 0 * * 1",
    "ip_binding": {
      "enabled": false,
      "interface": "wan",
      "bind_ip": ""
    },
    "auto_recovery": {
      "enabled": true,
      "max_retries": 3,
      "retry_interval": "5m"
    },
    "self_check": {
      "enabled": true,
      "interval": "168h"
    },
    "logging": false,
    "verbose": false
  },
  "logging": {
    "level": "info",
    "output": "stdout",
    "file": ""
  }
}
```

## 开发指南

### 环境要求

- Go 1.21 或更高版本
- Docker（可选，用于容器化部署）

### 本地开发

- 克隆仓库：

```bash
git clone https://github.com/nick3/speedtestup.git
cd speedtestup
```

- 安装依赖：

```bash
go mod download
```

- 本地构建：

```bash
# 构建当前平台
make build

# 构建所有平台
make build-all
```

- 本地运行：

```bash
# 运行程序
./speedup
```

### Docker 构建

```bash
# 构建镜像
docker build -t speedtestup:local .

# 运行容器
docker run -d --name speedtestup speedtestup:local

# 查看日志
docker logs -f speedtestup
```

## 架构说明

```
SpeedTestUp v2.0
├── api/               # API 封装
│   ├── ipapi.go          # IP 查询 API (ipinfo.io)
│   └── speedtestcn.go    # speedtest.cn API
├── service/           # 核心业务逻辑
│   ├── ip_service.go        # IP 服务
│   ├── speedup_service.go   # 提速服务
│   └── scheduler.go         # 调度服务
├── config/            # 配置管理
│   ├── config.go           # 配置结构
│   └── loader.go           # 配置加载
├── model/             # 数据模型
│   └── speedup.go          # 提速相关结构体
├── utils/             # 工具库
│   └── logger.go           # 日志工具
└── speedup.go         # 主程序
```

### 核心组件

#### IP 服务 (IPService)
- 获取当前公网 IP
- 验证 IP 绑定
- 检测 IP 变化

#### 提速服务 (SpeedupService)
- 执行提速操作
- 自动恢复机制
- 7 天自检

#### 调度服务 (Scheduler)
- 心跳检测（每 10 分钟）
- 7 天自检（每周一 0:0）
- 定期重启提速

## 📊 日志输出

程序运行时将输出详细的提速信息，包括：

```
[2024/11/08 15:30:45] [SpeedupService] [SUCCESS] 提速开始时间: 2024-11-08 15:30:45
[2024/11/08 15:30:45] [SpeedupService] [INFO] 出口IP地址: 192.168.1.100
[2024/11/08 15:30:45] [SpeedupService] [INFO] 一类上行带宽100M提速截至时间: 2024-11-15 15:30:45
[2024/11/08 15:30:45] [SpeedupService] [INFO] 二类上行带宽500M提速截至时间: 2024-11-15 15:30:45
[2024/11/08 15:30:45] [SpeedupService] [INFO] 下行带宽1000M提速截至时间: 2024-11-15 15:30:45
[2024/11/08 15:30:45] [SpeedupService] [SUCCESS] 上行提速已激活
[2024/11/08 15:30:45] [SpeedupService] [SUCCESS] 下行提速已激活
```

## 故障排除

### 常见问题

**Q: 提速失败怎么办？**
A: 检查网络连接和配置是否正确，程序会自动重试。

**Q: 如何查看详细日志？**
A: 修改配置中的 `logging.level` 为 `debug`。

**Q: 程序无法启动？**
A: 检查配置文件是否正确，使用 `--version` 查看版本信息。

### 调试模式

```bash
# 启用调试模式
./speedup --config config-debug.json
```

在配置文件中设置：
```json
{
  "logging": {
    "level": "debug"
  }
}
```

## 贡献指南

欢迎提交 Issue 和 Pull Request！

### 开发环境设置

```bash
# 克隆项目
git clone https://github.com/nick3/speedtestup.git
cd speedtestup

# 构建项目
make build

# 运行测试
make test
```

### 提交规范

- 使用标准 Go 代码格式
- 确保所有测试通过
- 添加必要的单元测试
- 更新相关文档

## 致谢

- [luci-app-broadbandacc](https://github.com/Diciya/luci-app-broadbandacc) - 原始架构参考
- [speedtest.cn](https://www.speedtest.cn/) - 提供提速 API 服务
- [Go](https://golang.org/) - 优秀的编程语言
- [resty](https://github.com/go-resty/resty) - HTTP 客户端库
- [cron](https://github.com/robfig/cron) - 定时任务调度库

## 许可证

MIT License

---

**基于 luci-app-broadbandacc 架构重构，打造更强大的 Go 语言版网络提速工具！** 🎉
