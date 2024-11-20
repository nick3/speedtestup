# 构建阶段
FROM golang:1.20-alpine AS builder

# 设置必要的编译环境
RUN apk add --no-cache ca-certificates && \
    apk add --no-cache tzdata

# 设置工作目录
WORKDIR /build

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 优化编译选项
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o speedup

# 运行阶段
FROM alpine:3.18

# 添加必要的 CA 证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# 复制编译好的二进制文件
COPY --from=builder /build/speedup /app/speedup

# 设置工作目录
WORKDIR /app

# 设置运行权限
RUN chmod +x /app/speedup

# 设置容器启动命令
CMD ["./speedup"]