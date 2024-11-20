.PHONY: build build-all

# 系统信息
GOOS_LINUX := linux
GOOS_MAC := darwin
GOOS_WINDOWS := windows
GOARCH := amd64

# 可执行文件后缀（Windows 需要 .exe）
BINARY_SUFFIX_WINDOWS := .exe
BINARY_SUFFIX_LINUX :=
BINARY_SUFFIX_MAC :=

# 获取版本信息
GIT_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
GO_VERSION := $(shell go version | cut -d ' ' -f 3)
VERSION := $(shell git describe --tags --always)

# 编译参数
LDFLAGS := -X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GitCommit=${GIT_COMMIT}' -X 'main.GoVersion=${GO_VERSION}' -w -s

# 默认编译当前平台
build:
	go build -ldflags "${LDFLAGS}" -o bin/speedup

# 编译所有平台
build-all:
	# Linux
	GOOS=${GOOS_LINUX} GOARCH=${GOARCH} go build -ldflags "${LDFLAGS}" \
		-o bin/speedup_${GOOS_LINUX}${BINARY_SUFFIX_LINUX}
	# MacOS
	GOOS=${GOOS_MAC} GOARCH=${GOARCH} go build -ldflags "${LDFLAGS}" \
		-o bin/speedup_${GOOS_MAC}${BINARY_SUFFIX_MAC}
	# Windows
	GOOS=${GOOS_WINDOWS} GOARCH=${GOARCH} go build -ldflags "${LDFLAGS}" \
		-o bin/speedup_${GOOS_WINDOWS}${BINARY_SUFFIX_WINDOWS}