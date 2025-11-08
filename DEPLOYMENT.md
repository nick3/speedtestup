# SpeedTestUp v2.0 部署指南

基于 luci-app-broadbandacc 架构重构的部署文档。

## 📦 构建项目

### 使用构建脚本（推荐）

```bash
# 给构建脚本执行权限
chmod +x build.sh

# 构建项目（包含测试）
./build.sh

# 构建所有平台
./build.sh all

# 跳过测试构建
./build.sh no-test
```

### 使用 Makefile

```bash
# 安装依赖
make install

# 构建项目
make build

# 构建所有平台
make build-all

# 运行测试
make test

# 生成覆盖率报告
make test-coverage
```

### 手动构建

```bash
# 1. 整理依赖
go mod tidy

# 2. 运行测试
go test -v ./...

# 3. 构建项目
go build -o speedup .

# 4. 查看版本
./speedup --version
```

## 🚀 部署方式

### 方式一：二进制文件部署

```bash
# 1. 构建项目
./build.sh

# 2. 复制到系统路径
sudo cp speedup /usr/local/bin/

# 3. 设置执行权限
sudo chmod +x /usr/local/bin/speedup

# 4. 创建配置目录
sudo mkdir -p /etc/speedtestup
cp config.json /etc/speedtestup/config.json

# 5. 编辑配置
sudo vim /etc/speedtestup/config.json

# 6. 创建 systemd 服务
sudo tee /etc/systemd/system/speedtestup.service > /dev/null <<EOF
[Unit]
Description=SpeedTestUp Network Acceleration Service
After=network.target

[Service]
Type=simple
User=nobody
WorkingDirectory=/etc/speedtestup
ExecStart=/usr/local/bin/speedup --config /etc/speedtestup/config.json
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# 7. 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable speedtestup
sudo systemctl start speedtestup

# 8. 查看服务状态
sudo systemctl status speedtestup
```

### 方式二：Docker 部署

#### 构建镜像

```bash
# 构建镜像
docker build -t speedtestup:latest .

# 标记版本
docker tag speedtestup:latest speedtestup:v2.0.0

# 推送到仓库（可选）
docker push speedtestup:latest
```

#### 运行容器

```bash
# 运行容器（后台）
docker run -d \
  --name speedtestup \
  --restart unless-stopped \
  -v /etc/speedtestup:/app/config \
  -v /var/log/speedtestup:/app/logs \
  speedtestup:latest

# 运行容器（交互式）
docker run -it \
  --name speedtestup \
  -v /etc/speedtestup:/app/config \
  -v /var/log/speedtestup:/app/logs \
  speedtestup:latest bash
```

#### Docker Compose

项目根目录已包含 `docker-compose.yml` 文件：

```yaml
version: '3.8'

services:
  speedtestup:
    image: ghcr.io/nick3/speedtestup:latest
    container_name: speedtestup
    restart: unless-stopped
    volumes:
      - ./config.json:/app/config.json:ro
      - ./logs:/app/logs
    environment:
      - TZ=Asia/Shanghai
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
    networks:
      - speedtestup-network

networks:
  speedtestup-network:
    driver: bridge
```

运行：

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## 🔧 配置管理

### 配置文件位置

- **Linux/macOS**: `~/.config/speedtestup/config.json` 或 `/etc/speedtestup/config.json`
- **Windows**: `%APPDATA%\speedtestup\config.json`

### 配置验证

```bash
# 验证配置文件
./speedup --config /path/to/config.json

# 查看版本信息
./speedup --version
```

### 配置示例

```json
{
  "speedup": {
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
    }
  },
  "logging": {
    "level": "info",
    "output": "file",
    "file": "/var/log/speedtestup/speedup.log"
  }
}
```

## 📊 监控与维护

### 查看日志

#### 二进制部署

```bash
# 查看 systemd 服务日志
sudo journalctl -u speedtestup -f

# 查看日志文件
tail -f /var/log/speedtestup/speedup.log
```

#### Docker 部署

```bash
# 查看容器日志
docker logs -f speedtestup

# 查看 docker-compose 日志
docker-compose logs -f speedtestup
```

### 健康检查

```bash
# 检查服务状态
sudo systemctl status speedtestup

# 重启服务
sudo systemctl restart speedtestup

# 停止服务
sudo systemctl stop speedtestup
```

### 性能监控

```bash
# 查看进程状态
ps aux | grep speedup

# 查看资源使用
top -p $(pgrep speedup)

# 查看网络连接
netstat -tulpn | grep speedup
```

## 🔄 升级与回滚

### 升级

```bash
# 1. 备份当前版本
sudo cp /usr/local/bin/speedup /usr/local/bin/speedup.backup

# 2. 下载新版本
wget https://github.com/nick3/speedtestup/releases/download/v2.0.1/speedup-linux-amd64

# 3. 替换文件
sudo cp speedup-linux-amd64 /usr/local/bin/speedup

# 4. 重启服务
sudo systemctl restart speedtestup

# 5. 验证升级
./speedup --version
```

### 回滚

```bash
# 回滚到上一版本
sudo cp /usr/local/bin/speedup.backup /usr/local/bin/speedup
sudo systemctl restart speedtestup
```

## 🐛 故障排除

### 常见问题

#### 1. 服务无法启动

```bash
# 检查配置文件
./speedup --config /etc/speedtestup/config.json --version

# 查看错误日志
sudo journalctl -u speedtestup -n 50
```

#### 2. 网络连接失败

```bash
# 测试 API 连接
curl -v https://ipinfo.io/ip
curl -v https://tisu-api-v3.speedtest.cn/speedUp/query
```

#### 3. 权限问题

```bash
# 检查文件权限
ls -l /usr/local/bin/speedup
ls -l /etc/speedtestup/config.json

# 修复权限
sudo chown root:root /usr/local/bin/speedup
sudo chmod 755 /usr/local/bin/speedup
```

### 调试模式

```bash
# 启用调试模式
./speedup --config config-debug.json

# 前台运行（查看实时日志）
./speedup
```

## 📝 维护任务

### 定期任务

- [ ] 每周检查服务状态
- [ ] 每月清理日志文件
- [ ] 每季度检查更新
- [ ] 每年备份配置

### 自动化脚本

创建维护脚本 `maintenance.sh`：

```bash
#!/bin/bash

# 维护任务脚本

LOG_FILE="/var/log/speedtestup/maintenance.log"

# 记录日志
log_message() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> $LOG_FILE
}

# 检查服务状态
check_service() {
    if systemctl is-active --quiet speedtestup; then
        log_message "服务运行正常"
    else
        log_message "服务未运行，尝试重启"
        systemctl restart speedtestup
    fi
}

# 清理旧日志
clean_logs() {
    find /var/log/speedtestup -name "*.log" -mtime +30 -delete
    log_message "清理旧日志完成"
}

# 主函数
main() {
    check_service
    clean_logs
    log_message "维护任务完成"
}

main
```

## 📞 支持

如遇到问题，请：

1. 查看日志文件
2. 搜索 GitHub Issues
3. 提交新的 Issue
4. 联系维护者

## 📄 许可证

MIT License
