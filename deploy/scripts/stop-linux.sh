#!/bin/bash

# AI热点追踪平台 - Linux 停止脚本
# 需要权限: chmod +x stop-linux.sh

echo "========================================"
echo "AI热点追踪平台 - Linux 停止脚本"
echo "========================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# 颜色输出函数
success() {
    echo -e "${GREEN}[完成]$NC $1"
}

error() {
    echo -e "${RED}[错误]$NC $1"
}

# 停止 Nginx
echo -e "${YELLOW}[1/4] 停止 Nginx...${NC}"
if command -v nginx &> /dev/null; then
    nginx -s quit || killall nginx || true
    sleep 1
    if pgrep nginx > /dev/null; then
        killall nginx 2>/dev/null || true
    fi
    success "Nginx 已停止"
else
    echo "[跳过] Nginx 未安装"
fi
echo ""

# 停止 Learning Path Service
echo -e "${YELLOW}[2/4] 停止 Learning Path Service...${NC}"
if pgrep learning-path-svc > /dev/null; then
    pkill learning-path-svc
    sleep 1
    if pgrep learning-path-svc > /dev/null; then
        killall learning-path-svc 2>/dev/null || true
    fi
    success "Learning Path Service 已停止"
else
    echo "[跳过] Learning Path Service 未运行"
fi
echo ""

# 停止 Content Service
echo -e "${YELLOW}[3/4] 停止 Content Service...${NC}"
if pgrep content-svc > /dev/null; then
    pkill content-svc
    sleep 1
    if pgrep content-svc > /dev/null; then
        killall content-svc 2>/dev/null || true
    fi
    success "Content Service 已停止"
else
    echo "[跳过] Content Service 未运行"
fi
echo ""

# 停止 Gateway
echo -e "${YELLOW}[4/4] 停止 Gateway...${NC}"
if pgrep gateway > /dev/null; then
    pkill gateway
    sleep 1
    if pgrep gateway > /dev/null; then
        killall gateway 2>/dev/null || true
    fi
    success "Gateway 已停止"
else
    echo "[跳过] Gateway 未运行"
fi
echo ""

# 清理 PID 文件
if [ -d "deploy/.pids" ]; then
    rm -f deploy/.pids/*.pid
    success "PID 文件已清理"
fi

echo "========================================"
echo -e "${GREEN}所有服务已停止！${NC}"
echo "========================================"
