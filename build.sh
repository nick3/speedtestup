#!/bin/bash

# SpeedTestUp 构建脚本
# 基于 luci-app-broadbandacc 架构重构

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 检查 Go 是否安装
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装，请先安装 Go"
        exit 1
    fi
    log_info "Go 版本: $(go version)"
}

# 检查 Git 是否安装
check_git() {
    if ! command -v git &> /dev/null; then
        log_warn "Git 未安装，使用默认版本信息"
        VERSION_HASH="unknown"
        return 1
    fi
    return 0
}

# 获取版本信息
get_version_info() {
    VERSION="2.0.0"
    BUILD_DATE=$(date +%Y-%m-%d)

    if check_git; then
        COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
        log_info "提交哈希: $COMMIT_HASH"
    else
        COMMIT_HASH="unknown"
    fi

    log_info "版本: $VERSION"
    log_info "构建日期: $BUILD_DATE"

    echo "VERSION=$VERSION" > .version
    echo "BUILD_DATE=$BUILD_DATE" >> .version
    echo "COMMIT_HASH=$COMMIT_HASH" >> .version
}

# 清理旧文件
clean_old_files() {
    log_step "清理旧文件..."
    rm -f speedup speedup-*
    rm -f coverage.out coverage.html
    log_info "旧文件已清理"
}

# 整理依赖
tidy_deps() {
    log_step "整理依赖..."
    go mod tidy
    log_info "依赖已整理"
}

# 运行测试
run_tests() {
    log_step "运行测试..."
    go test -v ./...
    if [ $? -eq 0 ]; then
        log_info "所有测试通过"
    else
        log_error "测试失败"
        exit 1
    fi
}

# 构建项目
build_project() {
    log_step "构建项目..."

    # 读取版本信息
    source .version

    LDFLAGS="-ldflags '-X main.version=$VERSION -X main.buildDate=$BUILD_DATE -X main.commitHash=$COMMIT_HASH'"

    # 构建当前平台
    eval "go build $LDFLAGS -o speedup ."
    log_info "当前平台构建完成: speedup"

    # 构建其他平台（可选）
    if [ "$1" = "all" ]; then
        log_step "构建所有平台..."

        eval "GOOS=linux GOARCH=amd64 go build $LDFLAGS -o speedup-linux-amd64 -ldflags '-w -s' ."
        log_info "Linux amd64 构建完成"

        eval "GOOS=darwin GOARCH=amd64 go build $LDFLAGS -o speedup-darwin-amd64 -ldflags '-w -s' ."
        log_info "Darwin amd64 构建完成"

        eval "GOOS=windows GOARCH=amd64 go build $LDFLAGS -o speedup-windows-amd64.exe -ldflags '-w -s' ."
        log_info "Windows amd64 构建完成"
    fi
}

# 显示构建结果
show_build_result() {
    log_step "构建结果:"
    ls -lh speedup* 2>/dev/null | awk '{print "  " $9 " (" $5 ")"}'
}

# 主函数
main() {
    echo "========================================"
    echo "  SpeedTestUp 构建脚本"
    echo "  基于 luci-app-broadbandacc 架构重构"
    echo "========================================"
    echo ""

    check_go
    get_version_info
    clean_old_files
    tidy_deps

    # 运行测试
    if [ "$1" != "no-test" ]; then
        run_tests
    fi

    # 构建项目
    build_project $1

    # 显示结果
    show_build_result

    echo ""
    log_info "构建完成！"
    log_info "使用 ./speedup 运行程序"
    log_info "使用 ./speedup --version 查看版本信息"

    # 清理版本文件
    rm -f .version
}

# 执行主函数
main $1
