#!/bin/bash

# AI热点追踪平台 - Linux/Mac 停止脚本

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "========================================"
echo "  AI热点追踪平台 - 停止服务"
echo "========================================"
echo ""

PID_DIR="$SCRIPT_DIR/pids"

# 停止 Nginx
echo "[1/3] 停止 Nginx..."
nginx -s stop 2>/dev/null
if [ $? -eq 0 ]; then
    echo "[√] Nginx 已停止"
else
    echo "[!] Nginx 未运行或已经停止"
fi
echo ""

# 停止 Content Service
echo "[2/3] 停止 Content Service..."
if [ -f "$PID_DIR/content-svc.pid" ]; then
    CONTENT_PID=$(cat "$PID_DIR/content-svc.pid")
    kill $CONTENT_PID 2>/dev/null
    rm -f "$PID_DIR/content-svc.pid"
    echo "[√] Content Service (PID: $CONTENT_PID) 已停止"
else
    echo "[!] 未找到 Content Service PID 文件"
    # 尝试通过进程名杀死
    pkill -f "content-svc/main.go" 2>/dev/null || true
fi
echo ""

# 停止 Gateway Service
echo "[3/3] 停止 Gateway Service..."
if [ -f "$PID_DIR/gateway.pid" ]; then
    GATEWAY_PID=$(cat "$PID_DIR/gateway.pid")
    kill $GATEWAY_PID 2>/dev/null
    rm -f "$PID_DIR/gateway.pid"
    echo "[√] Gateway Service (PID: $GATEWAY_PID) 已停止"
else
    echo "[!] 未找到 Gateway Service PID 文件"
    # 尝试通过进程名杀死
    pkill -f "gateway/main.go" 2>/dev/null || true
fi
echo ""

echo "所有服务已停止"
