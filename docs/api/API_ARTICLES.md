# 文章模块 API 接口文档

## 基础信息

- **Base URL**: `http://localhost/api` (通过 Nginx)
- **Content-Type**: `application/json`
- **认证**: 公开接口,无需登录

---

## 📋 接口列表

### 1. 获取文章分类列表

**请求**
```http
GET /api/articles/categories
```

**响应** `200 OK`
```json
[
  {
    "id": 1,
    "name": "动态",
    "code": "news",
    "color": "#3B82F6",
    "icon": "news",
    "sort_order": 1,
    "status": 1
  },
  {
    "id": 2,
    "name": "职业",
    "code": "impact",
    "color": "#F97316",
    "icon": "work",
    "sort_order": 2,
    "status": 1
  }
]
```

**分类说明**
| Code | 名称 | 颜色 | 用途 |
|------|------|------|------|
| news | 动态 | #3B82F6 | AI行业动态 |
| impact | 职业 | #F97316 | 职业影响分析 |
| learn | 学习 | #10B981 | 学习资源 |
| tool | 工具 | #8B5CF6 | 工具产品 |

---

### 2. 获取文章列表

**请求**
```http
GET /api/articles?page=1&page_size=10&category=news
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| page | int | 否 | 页码(从1开始) | 1 |
| page_size | int | 否 | 每页数量(1-100) | 10 |
| category | string | 否 | 分类代码(news/impact/learn/tool) | 全部 |

**响应** `200 OK`
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "GPT-5 发布：AI 能力再次飞跃，这些职业将受到影响",
        "summary": "OpenAI 今日发布 GPT-5，在多个基准测试中取得突破性进展...",
        "category_id": 1,
        "category_name": "动态",
        "source_id": 1,
        "source_name": "机器之心",
        "author": "AI观察",
        "published_at": "2026-03-30T10:00:00Z",
        "view_count": 1234,
        "comment_count": 23,
        "tags": [
          {"id": 1, "name": "GPT-5", "type": 1},
          {"id": 2, "name": "OpenAI", "type": 1}
        ]
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

### 3. 获取文章详情

**请求**
```http
GET /api/articles/1
```

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int64 | 是 | 文章ID |

**响应** `200 OK`
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "GPT-5 发布：AI 能力再次飞跃，这些职业将受到影响",
    "summary": "OpenAI 今日发布 GPT-5，在多个基准测试中取得突破性进展...",
    "content_mongo_id": "507f1f77bcf86cd799439011",
    "category_id": 1,
    "category_name": "动态",
    "source_id": 1,
    "source_name": "机器之心",
    "author": "AI观察",
    "published_at": "2026-03-30T10:00:00Z",
    "status": 1,
    "view_count": 1235,
    "comment_count": 23,
    "tags": [
      {"id": 1, "name": "GPT-5", "type": 1},
      {"id": 2, "name": "OpenAI", "type": 1},
      {"id": 3, "name": "职业影响", "type": 2},
      {"id": 11, "name": "大模型", "type": 1}
    ],
    "created_at": "2026-03-30T10:00:00Z",
    "updated_at": "2026-03-30T10:00:00Z"
  }
}
```

**注意**: 
- 文章内容存储在 MongoDB,通过 `content_mongo_id` 获取
- 每次访问会自动增加 `view_count`

---

## 🔧 前端调用示例

### Vue/Nuxt 示例

```typescript
// 1. 获取分类列表
const categories = await $fetch('/api/articles/categories')

// 2. 获取文章列表(带筛选)
const response = await $fetch('/api/articles', {
  query: {
    page: 1,
    page_size: 10,
    category: 'news' // 或 'impact', 'learn', 'tool'
  }
})

// 3. 获取文章详情
const article = await $fetch('/api/articles/1')
```

### React 示例

```typescript
// 使用 fetch
const response = await fetch('/api/articles?page=1&page_size=10')
const data = await response.json()

// 使用 axios
import axios from 'axios'
const { data } = await axios.get('/api/articles', {
  params: { page: 1, page_size: 10, category: 'news' }
})
```

---

## 📊 数据表关系

```
articles (文章主表)
  ├─ category_id → categories.id (分类)
  ├─ source_id → sources.id (来源媒体)
  ├─ article_stats (统计数据)
  │   ├─ view_count (阅读量)
  │   ├─ comment_count (评论数)
  │   └─ like_count (点赞数)
  └─ article_tag_relation → tags (标签,多对多)
```

---

## ⚠️ 错误响应

**400 Bad Request**
```json
{
  "code": 400,
  "message": "请求参数错误",
  "data": null
}
```

**404 Not Found**
```json
{
  "code": 404,
  "message": "文章不存在",
  "data": null
}
```

**500 Internal Server Error**
```json
{
  "code": 500,
  "message": "服务器内部错误",
  "data": null
}
```

---

## 📝 OpenAPI 文档

完整的 OpenAPI 3.0 规范文档请查看:
- [openapi.yaml](./openapi.yaml)

可以使用 Swagger UI 或 Postman 导入查看。
