#!/bin/bash

# 学习路径微服务启动脚本
# 支持直接运行或 Docker 模式

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 显示帮助信息
show_help() {
    echo -e "${GREEN}学习路径微服务启动脚本${NC}"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -d, --docker    使用 Docker 模式启动（默认）"
    echo "  -l, --local     本地直接运行"
    echo "  -b, --build     重新构建服务"
    echo "  -h, --help      显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0              # 使用 Docker 启动"
    echo "  $0 -l           # 本地直接运行"
    echo "  $0 -d -b        # 重新构建并 Docker 启动"
}

# 默认配置
MODE="docker"
BUILD=0
MYSQL_DSN="${MYSQL_DSN:-hot_ai:hotai123@tcp(localhost:3306)/hot_ai?charset=utf8mb4&parseTime=True&loc=Local}"

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--docker)
            MODE="docker"
            shift
            ;;
        -l|--local)
            MODE="local"
            shift
            ;;
        -b|--build)
            BUILD=1
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${YELLOW}未知选项: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 检查依赖
check_dependencies() {
    echo -e "${BLUE}检查依赖...${NC}"

    if [ "$MODE" == "local" ]; then
        if ! command -v go &> /dev/null; then
            echo -e "${YELLOW}错误: 未安装 Go${NC}"
            exit 1
        fi
        echo -e "${GREEN}✓ Go 环境正常${NC}"

        # 检查数据库
        if ! command -v mysql &> /dev/null; then
            echo -e "${YELLOW}警告: 未安装 MySQL 客户端，无法测试数据库连接${NC}"
        else
            # 测试数据库连接
            DB_HOST=$(echo "$MYSQL_DSN" | grep -oP '@tcp\(\K[^)]+')
            DB_USER=$(echo "$MYSQL_DSN" | cut -d':' -f1)
            DB_PASS=$(echo "$MYSQL_DSN" | cut -d':' -f2 | cut -d'@' -f1)

            if mysql -h "${DB_HOST%:*}" -P "${DB_HOST#*:}" -u "$DB_USER" -p"$DB_PASS" -e "SELECT 1" &> /dev/null; then
                echo -e "${GREEN}✓ 数据库连接正常${NC}"
            else
                echo -e "${YELLOW}警告: 数据库连接失败，请检查 MYSQL_DSN 配置${NC}"
            fi
        fi
    else
        if ! command -v docker &> /dev/null; then
            echo -e "${YELLOW}错误: 未安装 Docker${NC}"
            exit 1
        fi
        echo -e "${GREEN}✓ Docker 环境正常${NC}"

        if ! command -v docker-compose &> /dev/null; then
            echo -e "${YELLOW}错误: 未安装 Docker Compose${NC}"
            exit 1
        fi
        echo -e "${GREEN}✓ Docker Compose 环境正常${NC}"
    fi
}

# 本地运行
run_local() {
    echo -e "${BLUE}正在启动学习路径微服务（本地模式）...${NC}"
    echo -e "${YELLOW}数据库 DSN: $MYSQL_DSN${NC}"

    export MYSQL_DSN="$MYSQL_DSN"

    cd "$(dirname "$0")"

    echo -e "${BLUE}编译并启动服务...${NC}"
    go run apps/services/learning-path-svc/main.go \
        -f apps/services/learning-path-svc/etc/learning-path-svc.yaml
}

# Docker 运行
run_docker() {
    echo -e "${BLUE}正在启动学习路径微服务（Docker 模式）...${NC}"

    cd "$(dirname "$0")"

    # 构建镜像
    if [ $BUILD -eq 1 ]; then
        echo -e "${BLUE}正在构建 Docker 镜像...${NC}"
        docker-compose build learning-path-svc
    fi

    # 启动服务
    echo -e "${BLUE}正在启动服务...${NC}"
    docker-compose up -d learning-path-svc

    echo -e "${GREEN}✓ 服务已启动！${NC}"
    echo -e "${YELLOW}服务端口: 8003${NC}"
    echo -e "${YELLOW}日志查看: docker logs -f hot-ai-learning-path-svc${NC}"
    echo -e "${YELLOW}停止服务: docker-compose stop learning-path-svc${NC}"
}

# 主流程
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}  学习路径微服务启动工具${NC}"
echo -e "${GREEN}================================${NC}"
echo "启动模式: $MODE"
if [ $BUILD -eq 1 ]; then
    echo -e "${YELLOW}构建模式: 启用${NC}"
fi
echo ""

# 检查依赖
check_dependencies

echo ""

# 根据模式启动
if [ "$MODE" == "local" ]; then
    run_local
else
    run_docker
fi
