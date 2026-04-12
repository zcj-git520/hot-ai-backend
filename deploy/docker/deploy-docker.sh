#!/bin/bash

# AI热点追踪平台 - Docker 部署脚本

set -e

echo "========================================"
echo "AI热点追踪平台 - Docker 部署"
echo "========================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}[错误] 未检测到 Docker${NC}"
    echo -e "${YELLOW}[提示] 请先安装 Docker${NC}"
    exit 1
fi

# 检查 Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}[错误] 未检测到 Docker Compose${NC}"
    echo -e "${YELLOW}[提示] 请先安装 Docker Compose${NC}"
    exit 1
fi

echo -e "${GREEN}[完成] Docker 版本: $(docker --version)${NC}"
echo -e "${GREEN}[完成] Docker Compose 版本: $(docker-compose --version)${NC}"
echo ""

# 检查环境变量文件
if [ ! -f "docker/.env" ]; then
    echo -e "${YELLOW}[提示] 未找到环境变量文件，正在创建默认配置...${NC}"
    mkdir -p docker
    cat > docker/.env <<EOF
# MySQL 配置
MYSQL_ROOT_PASSWORD=root123
MYSQL_PASSWORD=hotai123
MYSQL_DATABASE=hot_ai
MYSQL_USER=hot_ai

# JWT Secret (生产环境请修改)
JWT_SECRET=your-jwt-secret-key-change-in-production

# Grafana 密码 (生产环境请修改)
GRAFANA_PASSWORD=admin123

# MongoDB 配置
MONGO_PASSWORD=root123

# Redis 配置
REDIS_PASSWORD=

# 数据库连接 (用于容器间通信)
DB_USER=hot_ai
DB_PASSWORD=hotai123
EOF
    echo -e "${GREEN}[完成] 默认环境变量文件已创建${NC}"
    echo -e "${YELLOW}[警告] 请修改 docker/.env 中的敏感信息${NC}"
    echo ""
fi

# 创建日志目录
mkdir -p logs

# 检查 Docker 网络是否存在
if ! docker network ls | grep -q hot-ai-network; then
    echo -e "${YELLOW}[提示] 创建 Docker 网络...${NC}"
    docker network create hot-ai-network
fi

echo -e "${GREEN}[1/3] 启动所有容器...${NC}"
cd "$(dirname "$0")/../.."
docker-compose up -d
echo -e "${GREEN}[完成] 所有容器已启动${NC}"
echo ""

echo -e "${GREEN}[2/3] 等待服务就绪...${NC}"
echo -e "${YELLOW}[提示] 服务启动需要 30-60 秒...${NC}"

# 等待 MySQL 就绪
echo -e "${YELLOW}等待 MySQL...${NC}"
for i in {1..30}; do
    if docker exec hot-ai-mysql mysqladmin ping -h localhost -u root -proot123 --silent 2>/dev/null; then
        echo -e "${GREEN}[完成] MySQL 已就绪${NC}"
        break
    fi
    sleep 2
done

# 等待其他服务就绪
sleep 10

echo -e "${GREEN}[完成] 服务已就绪${NC}"
echo ""

echo -e "${GREEN}[3/3] 检查服务状态...${NC}"
docker-compose ps
echo ""

echo "========================================"
echo -e "${GREEN}部署完成！${NC}"
echo "========================================"
echo ""
echo -e "${GREEN}服务地址：${NC}"
echo "  - MySQL: localhost:3306"
echo "  - MongoDB: localhost:27017"
echo "  - Redis: localhost:6379"
echo "  - Consul: localhost:8500"
echo "  - Prometheus: localhost:9090"
echo "  - Grafana: localhost:3001 (admin/admin123)"
echo "  - Gateway API: http://localhost:8000"
echo "  - Content API: http://localhost:8001"
echo "  - Nginx Proxy: http://localhost"
echo ""
echo -e "${GREEN}常用命令：${NC}"
echo "  - 查看日志: docker-compose logs -f [service-name]"
echo "  - 停止服务: docker-compose down"
echo "  - 重启服务: docker-compose restart [service-name]"
echo "  - 进入容器: docker exec -it hot-ai-gateway bash"
echo "  - 清理所有数据: docker-compose down -v"
echo ""
echo -e "${YELLOW}[警告] 生产环境请修改 docker/.env 中的所有密码${NC}"
echo ""
