# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置必要的编译环境
RUN apk add --no-cache ca-certificates make git && \
    apk add --no-cache tzdata

# 设置工作目录
WORKDIR /build

# 设置 Go 环境
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖（利用 Docker 缓存）
RUN go mod download

# 复制所有源代码
COPY . .

# 使用 Makefile 构建
RUN make build

# 运行阶段
FROM alpine:3.18

# 添加必要的 CA 证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# 复制编译好的二进制文件（注意路径与 Makefile 中的输出路径一致）
COPY --from=builder /build/bin/speedup /app/speedup

# 复制配置文件（可选）
COPY --from=builder /build/config.json /app/config.json

# 设置工作目录
WORKDIR /app

# 设置运行权限
RUN chmod +x /app/speedup

# 创建非 root 用户（安全最佳实践）
RUN adduser -D -s /bin/sh appuser
RUN chown -R appuser:appuser /app
USER appuser

# 设置容器启动命令
# 如果没有外部配置文件，使用内置配置
CMD ["./speedup"]