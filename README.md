# SpeedTestUp

[![Build and Release](https://github.com/nick3/speedtestup/actions/workflows/release.yml/badge.svg)](https://github.com/nick3/speedtestup/actions/workflows/release.yml)
[![Docker Build and Push](https://github.com/nick3/speedtestup/actions/workflows/docker.yml/badge.svg)](https://github.com/nick3/speedtestup/actions/workflows/docker.yml)

SpeedTestUp 是一个自动化工具，用于定期检测网络状态并自动触发提速操作。它能够监控 IP 变化并在需要时自动执行提速流程。

## 功能特性

- 🔄 自动检测 IP 变化（每 10 分钟）
- ⚡ 自动执行提速操作
- 📅 定时执行（每周一 0 点）
- 🐳 支持 Docker 部署
- 💻 支持多平台（Windows、Linux、MacOS）

## 快速开始

### 二进制文件运行

1. 从 [Releases](https://github.com/nick3/speedtestup/releases) 页面下载适合你系统的最新版本
2. 解压并运行程序：

```bash
# Linux/MacOS
chmod +x speedup-*
./speedup-*

# Windows
speedup-windows-amd64.exe
```

### Docker 运行

```bash
# 从 GitHub Container Registry 拉取镜像
docker pull ghcr.io/nick3/speedtestup:latest

# 运行容器
docker run -d --name speedtestup ghcr.io/nick3/speedtestup:latest
```

## 开发指南

### 环境要求

- Go 1.20 或更高版本
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

- 本地运行：

```bash
go run speedup.go
```

### Docker 构建

```bash
docker build -t speedtestup:local .
docker run -d speedtestup:local
```
