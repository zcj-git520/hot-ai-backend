# AI 热点追踪平台 - API 接口文档

## 架构说明

本项目采用微服务架构:

```
前端 (Nuxt.js) 
    ↓
网关服务 (gateway:8000) ← 反向代理/认证
    ↓
内容服务 (content-svc:8001) ← 文章/职业业务逻辑
```

- **网关服务** (`apps/gateway`): 端口 8000,负责认证、用户管理,并代理内容请求到 content-svc
- **内容服务** (`apps/services/content-svc`): 端口 8001,负责文章和职业的业务逻辑和数据存储

## 基础信息

- **网关 URL**: `http://localhost:8000/api` (前端调用此地址)
- **内容服务 URL**: `http://localhost:8001/api` (内部服务,通过网关代理访问)
- **Content-Type**: `application/json`
- **认证方式**: Bearer Token (JWT)

## 通用响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### 错误响应
```json
{
  "code": 400,
  "message": "错误信息",
  "data": null
}
```

---

## 📰 文章模块 (Articles)

### 1. 获取文章列表

**接口**: `GET /api/articles`

**权限**: 公开(无需登录)

**查询参数**:
| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| page | int | 否 | 页码 | 1 |
| page_size | int | 否 | 每页数量 | 10 |
| category | string | 否 | 分类筛选 | 全部 |

**请求示例**:
```
GET /api/articles?page=1&page_size=10&category=AI 动态
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "articles": [
      {
        "id": "art-1",
        "title": "GPT-5 发布：AI 能力再次飞跃，这些职业将受到影响",
        "summary": "OpenAI 今日发布 GPT-5...",
        "content": "详细内容...",
        "category": "AI 动态",
        "source": "机器之心",
        "author": "AI观察",
        "cover_image": "",
        "tags": ["GPT-5", "OpenAI", "职业影响"],
        "view_count": 1234,
        "comment_count": 23,
        "published_at": "2026-03-30T10:00:00Z",
        "status": "published",
        "created_at": "2026-03-30T10:00:00Z",
        "updated_at": "2026-03-30T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 2. 获取文章详情

**接口**: `GET /api/articles/:id`

**权限**: 公开(无需登录)

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 文章ID |

**请求示例**:
```
GET /api/articles/art-1
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "art-1",
    "title": "GPT-5 发布：AI 能力再次飞跃，这些职业将受到影响",
    "summary": "OpenAI 今日发布 GPT-5...",
    "content": "详细文章内容...",
    "category": "AI 动态",
    "source": "机器之心",
    "author": "AI观察",
    "cover_image": "",
    "tags": ["GPT-5", "OpenAI", "职业影响"],
    "view_count": 1235,
    "comment_count": 23,
    "published_at": "2026-03-30T10:00:00Z",
    "status": "published",
    "created_at": "2026-03-30T10:00:00Z",
    "updated_at": "2026-03-30T10:00:00Z"
  }
}
```

---

### 3. 获取文章分类列表

**接口**: `GET /api/articles/categories`

**权限**: 公开(无需登录)

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "cat-1",
      "name": "AI 动态",
      "slug": "dynamic",
      "description": "AI 技术最新动态和行业资讯",
      "sort_order": 1,
      "is_active": true
    },
    {
      "id": "cat-2",
      "name": "职业影响",
      "slug": "career",
      "description": "AI 对各行业职业的影响分析",
      "sort_order": 2,
      "is_active": true
    }
  ]
}
```

---

## 💼 职业模块 (Professions)

### 1. 获取职业列表

**接口**: `GET /api/professions`

**权限**: 公开(无需登录)

**查询参数**:
| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| page | int | 否 | 页码 | 1 |
| page_size | int | 否 | 每页数量 | 10 |
| risk_level | string | 否 | 风险等级筛选 | 全部 |

**风险等级可选值**: `extreme`, `high`, `medium`, `low`

**请求示例**:
```
GET /api/professions?page=1&page_size=10&risk_level=high
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "professions": [
      {
        "id": "prof-1",
        "name": "文案策划",
        "slug": "copywriter",
        "description": "负责撰写广告文案、营销材料的专业人员",
        "risk_level": "extreme",
        "automation_rate": 85,
        "safe_skills": ["创意构思", "品牌策略", "用户洞察"],
        "affected_skills": ["基础文案写作", "内容生成", "SEO优化"],
        "transform_tips": "向品牌策略师转型...",
        "learning_paths": ["brand-strategy", "creative-thinking"],
        "view_count": 1234,
        "is_active": true,
        "created_at": "2026-03-30T10:00:00Z",
        "updated_at": "2026-03-30T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 2. 获取职业详情

**接口**: `GET /api/professions/:slug`

**权限**: 公开(无需登录)

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| slug | string | 是 | 职业slug标识 |

**请求示例**:
```
GET /api/professions/copywriter
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "prof-1",
    "name": "文案策划",
    "slug": "copywriter",
    "description": "负责撰写广告文案、营销材料、品牌宣传内容的专业人员",
    "risk_level": "extreme",
    "automation_rate": 85,
    "safe_skills": ["创意构思", "品牌策略", "用户洞察"],
    "affected_skills": ["基础文案写作", "内容生成", "SEO优化"],
    "transform_tips": "向品牌策略师转型,专注于创意策划和品牌故事讲述,提升对用户心理的洞察力",
    "learning_paths": ["brand-strategy", "creative-thinking"],
    "view_count": 1235,
    "is_active": true,
    "created_at": "2026-03-30T10:00:00Z",
    "updated_at": "2026-03-30T10:00:00Z"
  }
}
```

---

### 3. 搜索职业

**接口**: `GET /api/professions/search`

**权限**: 公开(无需登录)

**查询参数**:
| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| q | string | 是 | 搜索关键词 | - |
| page | int | 否 | 页码 | 1 |
| page_size | int | 否 | 每页数量 | 10 |

**请求示例**:
```
GET /api/professions/search?q=设计师&page=1&page_size=10
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "professions": [
      {
        "id": "prof-2",
        "name": "平面设计",
        "slug": "graphic-designer",
        "description": "从事视觉传达设计的设计师",
        "risk_level": "high",
        "automation_rate": 75,
        "safe_skills": ["创意思维", "品牌设计", "用户体验"],
        "affected_skills": ["基础排版", "素材处理", "模板设计"],
        "transform_tips": "转向UI/UX设计或品牌视觉系统设计",
        "learning_paths": ["ui-ux-design", "brand-design"],
        "view_count": 2456,
        "is_active": true
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 4. 获取风险等级信息

**接口**: `GET /api/professions/risk-levels`

**权限**: 公开(无需登录)

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "risk-extreme",
      "level": "extreme",
      "name": "极高风险",
      "icon": "🔴",
      "description": "高度重复、规则明确的工作内容",
      "color": "red",
      "sort_order": 1
    },
    {
      "id": "risk-high",
      "level": "high",
      "name": "高风险",
      "icon": "🟠",
      "description": "大部分工作内容可被 AI 替代",
      "color": "orange",
      "sort_order": 2
    },
    {
      "id": "risk-medium",
      "level": "medium",
      "name": "中等风险",
      "icon": "🟡",
      "description": "部分工作内容可被 AI 辅助",
      "color": "yellow",
      "sort_order": 3
    },
    {
      "id": "risk-low",
      "level": "low",
      "name": "低风险",
      "icon": "🟢",
      "description": "需要创造力和复杂决策",
      "color": "green",
      "sort_order": 4
    }
  ]
}
```

---

## 🔧 前端对接说明

### 更新前端 API 配置

在前端项目的 `.env.local` 中配置后端地址:

```env
NUXT_PUBLIC_API_URL=http://localhost:8000/api
```

### 调用示例

```typescript
// 获取文章列表
const response = await articleApi.getList({
  page: 1,
  pageSize: 10,
  category: 'AI 动态'
})

// 获取职业列表
const response = await professionApi.getList({
  page: 1,
  pageSize: 10,
  riskLevel: 'high'
})

// 搜索职业
const response = await professionApi.search('设计师')

// 获取职业详情
const response = await $fetch('/api/professions/copywriter')
```

---

## 📊 数据库表结构

### articles (文章表)
- id, title, summary, content, category, source, author, cover_image, tags, view_count, comment_count, published_at, status, created_at, updated_at

### article_categories (文章分类表)
- id, name, slug, description, sort_order, is_active, created_at, updated_at

### professions (职业表)
- id, name, slug, description, risk_level, automation_rate, safe_skills, affected_skills, transform_tips, learning_paths, view_count, is_active, created_at, updated_at

### risk_level_info (风险等级表)
- id, level, name, icon, description, color, sort_order

---

## 🚀 启动后端服务

### 方式一:分别启动各服务(推荐)

**1. 启动内容服务 (content-svc)**
```bash
cd hot-ai-backend
go run apps/services/content-svc/main.go -f apps/services/content-svc/etc/content-svc.yaml
```
内容服务将在 `http://localhost:8001` 启动。

**2. 启动网关服务 (gateway)**
```bash
cd hot-ai-backend
go run apps/gateway/main.go -f apps/gateway/etc/gateway.yaml
```
网关服务将在 `http://localhost:8000` 启动。

### 方式二:使用脚本一键启动

创建启动脚本 `start.sh` (Linux/Mac) 或 `start.bat` (Windows):

**Linux/Mac (start.sh):**
```bash
#!/bin/bash
# 启动内容服务
go run apps/services/content-svc/main.go -f apps/services/content-svc/etc/content-svc.yaml &
CONTENT_PID=$!

# 等待内容服务启动
sleep 2

# 启动网关服务
go run apps/gateway/main.go -f apps/gateway/etc/gateway.yaml &
GATEWAY_PID=$!

echo "Content service PID: $CONTENT_PID"
echo "Gateway service PID: $GATEWAY_PID"

wait
```

**Windows (start.bat):**
```bat
@echo off
start "Content Service" go run apps/services/content-svc/main.go -f apps/services/content-svc/etc/content-svc.yaml
timeout /t 2
start "Gateway Service" go run apps/gateway/main.go -f apps/gateway/etc/gateway.yaml
```

### 配置说明

**内容服务配置** (`apps/services/content-svc/etc/content-svc.yaml`):
```yaml
Name: content-svc
Host: 0.0.0.0
Port: 8001

DataSource:
  MySQL:
    DSN: root:password@tcp(localhost:3306)/hot_ai?charset=utf8mb4&parseTime=True&loc=Local
```

**网关服务配置** (`apps/gateway/etc/gateway.yaml`):
```yaml
Name: gateway
Host: 0.0.0.0
Port: 8000

Services:
  ContentSvc: http://localhost:8001  # 内容服务地址
```

首次启动时会自动创建数据库表并填充初始数据。
