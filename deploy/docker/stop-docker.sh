#!/bin/bash

# AI热点追踪平台 - Docker 停止脚本

echo "========================================"
echo "AI热点追踪平台 - Docker 停止"
echo "========================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${YELLOW}[1/2] 停止所有容器...${NC}"
cd "$(dirname "$0")/../.."
docker-compose down
echo -e "${GREEN}[完成] 所有容器已停止${NC}"
echo ""

echo -e "${YELLOW}[2/2] 清理未使用的 Docker 资源...${NC}"
docker system prune -f
echo -e "${GREEN}[完成] 清理完成${NC}"
echo ""

echo "========================================"
echo -e "${GREEN}停止完成！${NC}"
echo "========================================"
echo ""
echo -e "${YELLOW}[提示] 如需保留数据卷，请使用: docker-compose down${NC}"
echo -e "${YELLOW}[提示] 如需完全清理（包括数据），请使用: docker-compose down -v${NC}"
echo ""
