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

# 复制所有源代码和 Makefile
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

# 设置工作目录
WORKDIR /app

# 设置运行权限
RUN chmod +x /app/speedup

# 设置容器启动命令
CMD ["./speedup"]