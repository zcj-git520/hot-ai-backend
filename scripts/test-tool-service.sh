#!/bin/bash

# 工具微服务测试脚本

echo "🚀 测试工具微服务"
echo "=================="

BASE_URL="http://localhost:8890/api/tools"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local endpoint=$1
    local name=$2
    
    echo -n "测试 $name... "
    response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint" 2>/dev/null)
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✓${NC} (HTTP $http_code)"
        # 只显示结果数量
        if echo "$body" | grep -q "total"; then
            total=$(echo "$body" | grep -o '"total":[0-9]*' | head -n1 | cut -d':' -f2)
            echo "   找到 $total 条数据"
        fi
    else
        echo -e "${RED}✗${NC} (HTTP $http_code)"
        echo "   响应: $body"
    fi
}

echo ""
echo "1. 测试工具类别列表"
echo "---------------------"
test_api "/categories" "工具类别列表"

echo ""
echo "2. 测试工具列表（默认）"
echo "---------------------"
test_api "" "工具列表"

echo ""
echo "3. 测试工具列表（分页）"
echo "---------------------"
test_api "?page=1&page_size=10" "分页工具列表"

echo ""
echo "4. 测试图像类工具"
echo "---------------------"
test_api "?category_id=2" "图像类工具"

echo ""
echo "5. 测试搜索功能"
echo "---------------------"
test_api "?search=ChatGPT&sort_by=rating&order=desc" "搜索ChatGPT"

echo ""
echo "6. 测试难度筛选"
echo "---------------------"
test_api "?difficulty=beginner" "入门级工具"

echo ""
echo "7. 测试工具详情（通过slug）"
echo "---------------------"
test_api "/chatgpt" "ChatGPT详情"

echo ""
echo "8. 测试工具详情（通过ID）"
echo "---------------------"
test_api "/id/1" "工具ID详情"

echo ""
echo "================================"
echo "✅ 测试完成！"
echo "================================"
