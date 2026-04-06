#!/bin/bash

# AI热点追踪平台 - Linux/Mac 启动脚本

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "========================================"
echo "  AI热点追踪平台 - 启动脚本"
echo "========================================"
echo ""

# 检查Nginx是否安装
if ! command -v nginx &> /dev/null; then
    echo "[错误] 未找到 Nginx，请先安装 Nginx"
    echo ""
    echo "Ubuntu/Debian: sudo apt-get install nginx"
    echo "CentOS/RHEL: sudo yum install nginx"
    echo "MacOS: brew install nginx"
    exit 1
fi

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# PID文件目录
PID_DIR="$SCRIPT_DIR/pids"
mkdir -p "$PID_DIR"

# 停止已运行的服务
echo "[0/4] 清理旧进程..."
if [ -f "$PID_DIR/content-svc.pid" ]; then
    kill $(cat "$PID_DIR/content-svc.pid") 2>/dev/null || true
    rm -f "$PID_DIR/content-svc.pid"
fi

if [ -f "$PID_DIR/gateway.pid" ]; then
    kill $(cat "$PID_DIR/gateway.pid") 2>/dev/null || true
    rm -f "$PID_DIR/gateway.pid"
fi

nginx -s stop 2>/dev/null || true
sleep 1
echo -e "${GREEN}[√]${NC} 清理完成"
echo ""

# 启动 Content Service
echo "[1/4] 启动 Content Service (端口 8001)..."
nohup go run apps/services/content-svc/main.go -f apps/services/content-svc/etc/content-svc.yaml > logs/content-svc.log 2>&1 &
CONTENT_PID=$!
echo $CONTENT_PID > "$PID_DIR/content-svc.pid"
echo -e "${GREEN}[√]${NC} Content Service 已启动 (PID: $CONTENT_PID)"
sleep 2
echo ""

# 启动 Gateway Service
echo "[2/4] 启动 Gateway Service (端口 8000)..."
nohup go run apps/gateway/main.go -f apps/gateway/etc/gateway.yaml > logs/gateway.log 2>&1 &
GATEWAY_PID=$!
echo $GATEWAY_PID > "$PID_DIR/gateway.pid"
echo -e "${GREEN}[√]${NC} Gateway Service 已启动 (PID: $GATEWAY_PID)"
sleep 2
echo ""

# 启动 Nginx
echo "[3/4] 启动 Nginx (端口 80)..."
nginx -t -c "$SCRIPT_DIR/nginx.conf" 2>&1
if [ $? -eq 0 ]; then
    nginx -c "$SCRIPT_DIR/nginx.conf"
    echo -e "${GREEN}[√]${NC} Nginx 已启动"
else
    echo -e "${RED}[×]${NC} Nginx 配置测试失败"
    exit 1
fi
echo ""

# 检查服务状态
echo "[4/4] 检查服务状态..."
sleep 2

echo ""
echo "========================================"
echo -e "  ${GREEN}所有服务已启动!${NC}"
echo "========================================"
echo ""
echo "访问地址:"
echo "  - API网关: http://localhost/api"
echo "  - 健康检查: http://localhost/health"
echo "  - Content Service直接访问: http://localhost:8001/api"
echo "  - Gateway直接访问: http://localhost:8000/api"
echo ""
echo "进程信息:"
echo "  - Content Service PID: $CONTENT_PID"
echo "  - Gateway Service PID: $GATEWAY_PID"
echo "  - Nginx: 运行中"
echo ""
echo "日志文件:"
echo "  - Content Service: logs/content-svc.log"
echo "  - Gateway: logs/gateway.log"
echo "  - Nginx Access: logs/access.log"
echo "  - Nginx Error: logs/error.log"
echo ""
echo "管理命令:"
echo "  - 停止所有服务: ./stop.sh"
echo "  - 重启Nginx: nginx -s reload"
echo "  - 查看Content Service日志: tail -f logs/content-svc.log"
echo "  - 查看Gateway日志: tail -f logs/gateway.log"
echo ""
