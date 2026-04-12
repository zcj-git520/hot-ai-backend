# 工具库模块 API 文档

## 概述

本文档描述了工具库模块的 API 接口设计。

**基础路径**: `/api/v1/tools`

**认证方式**: 部分接口需要认证（根据具体接口需求）

---

## 目录

1. [工具类别接口](#1-工具类别接口)
2. [工具列表接口](#2-工具列表接口)
3. [工具详情接口](#3-工具详情接口)

---

## 1. 工具类别接口

### 1.1 获取所有工具类别

获取所有工具类别列表。

**接口**: `GET /api/tools/categories`

**请求参数**: 无

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "写作类",
      "slug": "writing",
      "icon": "✍️",
      "description": "用于写作、文案创作的 AI 工具",
      "sort_order": 1,
      "featured": true,
      "status": 1,
      "created_at": "2026-04-11T00:00:00Z",
      "updated_at": "2026-04-11T00:00:00Z"
    }
  ]
}
```

---

## 2. 工具列表接口

### 2.1 获取工具列表

获取工具列表，支持分页和筛选。

**接口**: `GET /api/tools`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| page | int | 否 | 页码，默认 1 | 1 |
| page_size | int | 否 | 每页数量，默认 20 | 20 |
| category_id | int | 否 | 类别 ID | 1 |
| difficulty | string | 否 | 难度等级：beginner/intermediate/advanced | beginner |
| min_rating | float | 否 | 最低评分 | 4.0 |
| search | string | 否 | 搜索关键词 | ChatGPT |
| sort_by | string | 否 | 排序字段：popularity/rating/update_time | popularity |
| order | string | 否 | 排序方向：asc/desc | desc |

**请求示例**:
```bash
GET /api/tools?page=1&page_size=10&category_id=1&sort_by=rating&order=desc
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "ChatGPT",
        "slug": "chatgpt",
        "icon": "🤖",
        "description": "OpenAI 开发的大型语言模型",
        "official_url": "https://chat.openai.com",
        "documentation_url": "https://help.openai.com",
        "pricing": "{\"free\":{\"available\":true,\"limit\":\"每天消息限额\",\"features\":[\"基础对话\",\"基础写作\"]}}",
        "pricing_description": "免费版每天有消息配额，Plus 版本 $20/月",
        "category_id": 1,
        "difficulty": "beginner",
        "rating": 4.8,
        "review_count": 2300,
        "view_count": 15000,
        "popularity": 95,
        "tags": "[1,2,3,7]",
        "featured": true,
        "status": 1,
        "created_at": "2026-04-11T00:00:00Z",
        "updated_at": "2026-04-11T00:00:00Z"
      }
    ],
    "total": 200,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 3. 工具详情接口

### 3.1 获取工具详情（通过slug）

根据工具 slug 获取工具详情。

**接口**: `GET /api/tools/:slug`

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| slug | string | 是 | 工具 slug |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "ChatGPT",
    "slug": "chatgpt",
    "icon": "🤖",
    "description": "OpenAI 开发的大型语言模型",
    "official_url": "https://chat.openai.com",
    "documentation_url": "https://help.openai.com",
    "pricing": "{\"free\":{\"available\":true,\"limit\":\"每天消息限额\",\"features\":[\"基础对话\",\"基础写作\"]}}",
    "pricing_description": "免费版每天有消息配额，Plus 版本 $20/月",
    "category_id": 1,
    "difficulty": "beginner",
    "rating": 4.8,
    "review_count": 2300,
    "view_count": 15000,
    "popularity": 95,
    "tags": "[1,2,3,7]",
    "featured": true,
    "status": 1,
    "created_at": "2026-04-11T00:00:00Z",
    "updated_at": "2026-04-11T00:00:00Z"
  }
}
```

### 3.2 获取工具详情（通过ID）

根据工具 ID 获取工具详情。

**接口**: `GET /api/tools/id/:id`

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 工具 ID |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "ChatGPT",
    "slug": "chatgpt",
    "icon": "🤖",
    "description": "OpenAI 开发的大型语言模型",
    "official_url": "https://chat.openai.com",
    "documentation_url": "https://help.openai.com",
    "pricing": "{\"free\":{\"available\":true,\"limit\":\"每天消息限额\",\"features\":[\"基础对话\",\"基础写作\"]}}",
    "pricing_description": "免费版每天有消息配额，Plus 版本 $20/月",
    "category_id": 1,
    "difficulty": "beginner",
    "rating": 4.8,
    "review_count": 2300,
    "view_count": 15000,
    "popularity": 95,
    "tags": "[1,2,3,7]",
    "featured": true,
    "status": 1,
    "created_at": "2026-04-11T00:00:00Z",
    "updated_at": "2026-04-11T00:00:00Z"
  }
}
```

---

## 错误响应

### 通用错误格式

```json
{
  "code": 400,
  "message": "错误信息描述",
  "data": null
}
```

### 常见错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 使用示例

### 获取工具列表

```bash
# 获取第1页，每页10条
curl "http://localhost:8890/api/tools?page=1&page_size=10"

# 获取图像类工具
curl "http://localhost:8890/api/tools?category_id=2"

# 搜索工具
curl "http://localhost:8890/api/tools?search=ChatGPT&sort_by=rating&order=desc"
```

### 获取工具详情

```bash
# 通过slug获取
curl "http://localhost:8890/api/tools/chatgpt"

# 通过ID获取
curl "http://localhost:8890/api/tools/id/1"
```

---

## 数据字段说明

### Tool 对象字段

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 工具ID |
| name | string | 工具名称 |
| slug | string | URL友好的标识 |
| icon | string | 工具图标（emoji或URL） |
| description | string | 工具描述 |
| official_url | string | 官方网站 |
| documentation_url | string | 文档链接 |
| pricing | string | 定价信息（JSON） |
| pricing_description | string | 定价说明 |
| category_id | uint | 所属类别ID |
| difficulty | string | 难度等级 |
| rating | float | 平均评分（0-5） |
| review_count | int | 评测数量 |
| view_count | int | 浏览量 |
| popularity | int | 热度值 |
| tags | string | 标签列表（JSON） |
| featured | bool | 是否精选 |
| status | int | 状态（1-上架，0-下架） |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |

---

**文档结束**
