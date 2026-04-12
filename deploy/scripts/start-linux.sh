#!/bin/bash

# AI热点追踪平台 - Linux 启动脚本
# 需要权限: chmod +x start-linux.sh

set -e

echo "========================================"
echo "AI热点追踪平台 - Linux 启动脚本"
echo "========================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo -e "${RED}[错误] 未检测到 Go，请先安装 Go 1.22+${NC}"
    exit 1
fi
echo -e "${GREEN}[完成] Go 版本: $(go version)${NC}"
echo ""

# 检查 MySQL 是否运行
if ! command -v mysql &> /dev/null; then
    echo -e "${YELLOW}[警告] 未检测到 MySQL${NC}"
    echo -e "${YELLOW}[提示] 数据库: hot_ai，用户: root${NC}"
else
    echo -e "${GREEN}[完成] MySQL 已安装${NC}"
fi
echo ""

# 检查 Nginx 是否安装
if ! command -v nginx &> /dev/null; then
    echo -e "${YELLOW}[警告] 未检测到 Nginx${NC}"
    echo -e "${YELLOW}[提示] 安装 Nginx 可获得反向代理功能${NC}"
else
    echo -e "${GREEN}[完成] Nginx 已安装${NC}"
fi
echo ""

# 检查配置文件
if [ ! -f "config/app.env" ]; then
    echo -e "${YELLOW}[提示] 未找到配置文件，正在创建默认配置...${NC}"
    mkdir -p config
    cat > config/app.env <<EOF
DB_DSN=root:root@tcp(localhost:3306)/hot_ai?charset=utf8mb4&parseTime=True&loc=Local
CONTENT_SVC_URL=http://localhost:8001
JWT_SECRET=your-jwt-secret-key-change-in-production
REDIS_URL=localhost:6379
EOF
    echo -e "${GREEN}[完成] 默认配置文件已创建${NC}"
    echo -e "${YELLOW}[警告] 请修改 config/app.env 中的数据库密码和 JWT Secret${NC}"
    echo ""
fi

# 创建 logs 目录
mkdir -p logs

echo -e "${GREEN}[1/6] 编译 Gateway 服务...${NC}"
cd "$(dirname "$0")/../.."
go build -o bin/gateway ./apps/gateway || exit 1
echo -e "${GREEN}[完成] Gateway 编译成功${NC}"
echo ""

echo -e "${GREEN}[2/6] 编译 Content Service...${NC}"
go build -o bin/content-svc ./apps/services/content-svc || exit 1
echo -e "${GREEN}[完成] Content Service 编译成功${NC}"
echo ""

echo -e "${GREEN}[3/6] 编译 Learning Path Service...${NC}"
go build -o bin/learning-path-svc ./apps/services/learning-path-svc || exit 1
echo -e "${GREEN}[完成] Learning Path Service 编译成功${NC}"
echo ""

echo -e "${GREEN}[4/6] 启动 Gateway (端口 8000)...${NC}"
nohup bin/gateway -f apps/gateway/etc/gateway.yaml > logs/gateway.log 2>&1 &
GATEWAY_PID=$!
sleep 3
if ps -p $GATEWAY_PID > /dev/null; then
    echo -e "${GREEN}[完成] Gateway 已启动 (PID: $GATEWAY_PID)${NC}"
else
    echo -e "${RED}[错误] Gateway 启动失败${NC}"
    tail -n 20 logs/gateway.log
    exit 1
fi
echo ""

echo -e "${GREEN}[5/6] 启动 Content Service (端口 8001)...${NC}"
nohup bin/content-svc -f apps/services/content-svc/etc/content-svc.yaml > logs/content-svc.log 2>&1 &
CONTENT_SVC_PID=$!
sleep 3
if ps -p $CONTENT_SVC_PID > /dev/null; then
    echo -e "${GREEN}[完成] Content Service 已启动 (PID: $CONTENT_SVC_PID)${NC}"
else
    echo -e "${RED}[错误] Content Service 启动失败${NC}"
    tail -n 20 logs/content-svc.log
    exit 1
fi
echo ""

echo -e "${GREEN}[6/6] 启动 Learning Path Service (端口 8003)...${NC}"
nohup bin/learning-path-svc -f apps/services/learning-path-svc/etc/learning-path-svc.yaml > logs/learning-path-svc.log 2>&1 &
LEARNING_PATH_PID=$!
sleep 3
if ps -p $LEARNING_PATH_PID > /dev/null; then
    echo -e "${GREEN}[完成] Learning Path Service 已启动 (PID: $LEARNING_PATH_PID)${NC}"
else
    echo -e "${RED}[错误] Learning Path Service 启动失败${NC}"
    tail -n 20 logs/learning-path-svc.log
    exit 1
fi
echo ""

# 启动 Nginx
if command -v nginx &> /dev/null; then
    echo -e "${GREEN}[7/7] 启动 Nginx (端口 80)...${NC}"
    nginx -c "$(dirname "$0")/../nginx.conf" || true
    sleep 2
    echo -e "${GREEN}[完成] Nginx 已启动${NC}"
    echo ""
else
    echo -e "${YELLOW}[跳过] Nginx 未安装，跳过启动${NC}"
    echo ""
fi

echo "========================================"
echo -e "${GREEN}所有服务启动完成！${NC}"
echo "========================================"
echo ""
echo -e "${GREEN}服务地址：${NC}"
echo "  - Gateway API: http://localhost:8000"
echo "  - Content API: http://localhost:8001"
echo "  - Learning Path API: http://localhost:8003"
if command -v nginx &> /dev/null; then
    echo "  - Nginx Proxy: http://localhost"
fi
echo ""
echo -e "${GREEN}日志文件位置：${NC}"
echo "  - Gateway: logs/gateway.log"
echo "  - Content Service: logs/content-svc.log"
echo "  - Learning Path Service: logs/learning-path-svc.log"
if command -v nginx &> /dev/null; then
    echo "  - Nginx: logs/nginx/access.log, logs/nginx/error.log"
fi
echo ""
echo -e "${GREEN}查看日志：${NC}"
echo "  - Gateway: tail -f logs/gateway.log"
echo "  - Content Service: tail -f logs/content-svc.log"
echo "  - Learning Path Service: tail -f logs/learning-path-svc.log"
echo ""
echo -e "${GREEN}停止服务：${NC}"
echo "  - 运行: ./stop-linux.sh"
echo "  - 或手动: killall gateway content-svc learning-path-svc nginx"
echo ""

# 保存 PID
echo $GATEWAY_PID > deploy/.pids/gateway.pid
echo $CONTENT_SVC_PID > deploy/.pids/content-svc.pid
echo $LEARNING_PATH_PID > deploy/.pids/learning-path-svc.pid
mkdir -p deploy/.pids
