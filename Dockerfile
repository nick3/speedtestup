
# 使用 Go 官方轻量级镜像
FROM golang:1.20-alpine

# 设置工作目录
WORKDIR /app

# 复制依赖文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目代码
COPY . .

# 编译项目
RUN go build -o speedup

# 设置容器启动命令
CMD ["./speedup"]