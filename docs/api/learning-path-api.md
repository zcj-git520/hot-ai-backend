# 学习路径 API 接口文档

## 概述

本文档描述了学习路径管理相关的 API 接口。所有接口均基于 RESTful 风格设计，使用 JSON 格式进行数据交换。

**基础路径**: `/api/v1`  
**认证方式**: JWT Bearer Token（部分接口需要）

---

## 目录

1. [学习路径接口](#学习路径接口)
   - [获取学习路径列表](#获取学习路径列表)
   - [获取推荐路径](#获取推荐路径)
   - [获取难度等级信息](#获取难度等级信息)
   - [根据 ID 获取学习路径详情](#根据 id 获取学习路径详情)
   - [根据 slug 获取学习路径详情](#根据 slug 获取学习路径详情)

2. [章节接口](#章节接口)
   - [获取路径的所有章节](#获取路径的所有章节)
   - [根据章节 ID 获取详情](#根据章节 id 获取详情)
   - [根据路径 slug 和章节 slug 获取章节详情](#根据路径 slug 和章节 slug 获取章节详情)

3. [学习进度接口](#学习进度接口)
   - [获取用户的学习进度](#获取用户的学习进度)
   - [获取用户已完成的章节列表](#获取用户已完成的章节列表)
   - [保存学习进度](#保存学习进度)

4. [综合接口](#综合接口)
   - [获取路径学习仪表盘](#获取路径学习仪表盘)

---

## 数据结构

### LearningPath - 学习路径

```json
{
  "id": 1,
  "title": "Python 入门教程",
  "slug": "python-basics",
  "icon": "🐍",
  "description": "从零开始学习 Python 编程",
  "difficulty": "beginner",
  "level_label": "入门",
  "learning_goals": ["掌握 Python 基础语法", "理解面向对象编程"],
  "target_audience": ["编程初学者", "转行开发者"],
  "estimated_days": 30,
  "estimated_hours": 60,
  "chapter_count": 12,
  "student_count": 1500,
  "cover_image": "https://example.com/cover.jpg",
  "is_featured": true,
  "published_at": "2026-01-01T00:00:00Z",
  "chapters": []
}
```

### PathChapter - 学习路径章节

```json
{
  "id": 1,
  "path_id": 1,
  "title": "Python 基础语法",
  "slug": "basic-syntax",
  "description": "学习 Python 的基础语法",
  "content_type": "article",
  "content": "文章内容...",
  "video_url": "",
  "external_links": [],
  "estimated_hours": 2,
  "order_index": 1,
  "is_free": true
}
```

### LevelInfo - 难度等级信息

```json
{
  "id": "beginner",
  "level": "beginner",
  "name": "入门",
  "icon": "🌱",
  "description": "零基础友好，讲解基础概念和工具使用",
  "color": "#4ade80",
  "min_hours": 20,
  "max_hours": 40
}
```

---

## 学习路径接口

### 获取学习路径列表

获取所有学习路径的列表，支持分页和筛选。

**接口信息**:
- **路径**: `GET /api/learning-paths`
- **认证**: 不需要

**查询参数**:

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| page | int | 否 | 页码，默认 1 | 1 |
| page_size | int | 否 | 每页数量，默认 12 | 12 |
| difficulty | string | 否 | 难度筛选：beginner/intermediate/advanced | beginner |

**成功响应** (200 OK):

```json
{
  "list": [
    {
      "id": 1,
      "title": "Python 入门教程",
      "slug": "python-basics",
      "icon": "🐍",
      "description": "从零开始学习 Python 编程",
      "difficulty": "beginner",
      "level_label": "入门",
      "estimated_days": 30,
      "estimated_hours": 60,
      "chapter_count": 12,
      "student_count": 1500,
      "cover_image": "https://example.com/cover.jpg",
      "is_featured": true
    }
  ],
  "total": 50,
  "page": 1
}
```

---

### 获取推荐路径

获取推荐的学习路径列表。

**接口信息**:
- **路径**: `GET /api/learning-paths/featured`
- **认证**: 不需要

**查询参数**:

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| limit | int | 否 | 返回数量，默认 4 | 4 |

**成功响应** (200 OK):

```json
[
  {
    "id": 1,
    "title": "Python 入门教程",
    "slug": "python-basics",
    "icon": "🐍",
    "description": "从零开始学习 Python 编程",
    "difficulty": "beginner",
    "level_label": "入门",
    "estimated_days": 30,
    "estimated_hours": 60,
    "chapter_count": 12,
    "student_count": 1500,
    "cover_image": "https://example.com/cover.jpg",
    "is_featured": true
  }
]
```

---

### 获取难度等级信息

获取所有难度等级的展示信息。

**接口信息**:
- **路径**: `GET /api/learning-paths/levels`
- **认证**: 不需要

**成功响应** (200 OK):

```json
[
  {
    "id": "beginner",
    "level": "beginner",
    "name": "入门",
    "icon": "🌱",
    "description": "零基础友好，讲解基础概念和工具使用",
    "color": "#4ade80",
    "min_hours": 20,
    "max_hours": 40
  },
  {
    "id": "intermediate",
    "level": "intermediate",
    "name": "进阶",
    "icon": "✍️",
    "description": "需要一定基础，讲解实用技巧和进阶应用",
    "color": "#60a5fa",
    "min_hours": 30,
    "max_hours": 60
  },
  {
    "id": "advanced",
    "level": "advanced",
    "name": "高级",
    "icon": "🚀",
    "description": "需要扎实基础，讲解深度应用和开发技能",
    "color": "#a855f7",
    "min_hours": 40,
    "max_hours": 100
  }
]
```

---

### 根据 ID 获取学习路径详情

根据学习路径 ID 获取详细信息，包括所有章节。

**接口信息**:
- **路径**: `GET /api/learning-paths/id/{id}`
- **认证**: 不需要

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 学习路径 ID |

**成功响应** (200 OK):

```json
{
  "id": 1,
  "title": "Python 入门教程",
  "slug": "python-basics",
  "icon": "🐍",
  "description": "从零开始学习 Python 编程",
  "difficulty": "beginner",
  "level_label": "入门",
  "learning_goals": ["掌握 Python 基础语法", "理解面向对象编程"],
  "target_audience": ["编程初学者", "转行开发者"],
  "estimated_days": 30,
  "estimated_hours": 60,
  "chapter_count": 12,
  "student_count": 1500,
  "cover_image": "https://example.com/cover.jpg",
  "is_featured": true,
  "published_at": "2026-01-01T00:00:00Z",
  "chapters": [
    {
      "id": 1,
      "path_id": 1,
      "title": "Python 基础语法",
      "slug": "basic-syntax",
      "description": "学习 Python 的基础语法",
      "content_type": "article",
      "estimated_hours": 2,
      "order_index": 1,
      "is_free": true
    }
  ]
}
```

---

### 根据 slug 获取学习路径详情

根据学习路径 slug 获取详细信息。

**接口信息**:
- **路径**: `GET /api/learning-paths/slug/{slug}`
- **认证**: 不需要

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| slug | string | 是 | 学习路径 slug |

**成功响应**: 同「根据 ID 获取学习路径详情」

---

## 章节接口

### 获取路径的所有章节

获取指定学习路径的所有章节列表。

**接口信息**:
- **路径**: `GET /api/learning-paths/{path_id}/chapters`
- **认证**: 不需要

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| path_id | int | 是 | 学习路径 ID |

**成功响应** (200 OK):

```json
[
  {
    "id": 1,
    "path_id": 1,
    "title": "Python 基础语法",
    "slug": "basic-syntax",
    "description": "学习 Python 的基础语法",
    "content_type": "article",
    "estimated_hours": 2,
    "order_index": 1,
    "is_free": true
  },
  {
    "id": 2,
    "path_id": 1,
    "title": "数据类型与变量",
    "slug": "data-types",
    "description": "学习 Python 的数据类型",
    "content_type": "article",
    "estimated_hours": 3,
    "order_index": 2,
    "is_free": true
  }
]
```

---

### 根据章节 ID 获取详情

根据章节 ID 获取章节详情。

**接口信息**:
- **路径**: `GET /api/chapters/{chapter_id}`
- **认证**: 不需要

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| chapter_id | int | 是 | 章节 ID |

**成功响应** (200 OK):

```json
{
  "id": 1,
  "path_id": 1,
  "title": "Python 基础语法",
  "slug": "basic-syntax",
  "description": "学习 Python 的基础语法",
  "content_type": "article",
  "content": "文章内容...",
  "video_url": "",
  "external_links": [],
  "estimated_hours": 2,
  "order_index": 1,
  "is_free": true
}
```

---

### 根据路径 slug 和章节 slug 获取章节详情

根据路径 slug 和章节 slug 获取章节详情，包含前一章和下一章信息。

**接口信息**:
- **路径**: `GET /api/learning-paths/{path_slug}/chapters/{chapter_slug}`
- **认证**: 不需要

**路径参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| path_slug | string | 是 | 学习路径 slug |
| chapter_slug | string | 是 | 章节 slug |

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| path_id | int | 否 | 学习路径 ID（优先使用） |

**成功响应** (200 OK):

```json
{
  "chapter": {
    "id": 1,
    "path_id": 1,
    "title": "Python 基础语法",
    "slug": "basic-syntax",
    "description": "学习 Python 的基础语法",
    "content_type": "article",
    "content": "文章内容...",
    "estimated_hours": 2,
    "order_index": 1,
    "is_free": true
  },
  "prev": null,
  "next": {
    "id": 2,
    "path_id": 1,
    "title": "数据类型与变量",
    "slug": "data-types",
    "description": "学习 Python 的数据类型",
    "content_type": "article",
    "estimated_hours": 3,
    "order_index": 2,
    "is_free": true
  }
}
```

---

## 学习进度接口

### 获取用户的学习进度

获取用户在指定学习路径上的学习进度。

**接口信息**:
- **路径**: `GET /api/learning-paths/progress`
- **认证**: **需要**

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | string | 否 | 用户 ID |
| path_id | int | 是 | 学习路径 ID |

**成功响应** (200 OK):

```json
{
  "id": 1,
  "user_id": "user-123",
  "path_id": 1,
  "chapter_id": 5,
  "status": "in_progress",
  "time_spent": 180,
  "notes": "学习笔记",
  "created_at": "2026-03-01T10:00:00Z",
  "updated_at": "2026-03-08T15:30:00Z"
}
```

---

### 获取用户已完成的章节列表

获取用户在指定学习路径上已完成的章节 ID 列表。

**接口信息**:
- **路径**: `GET /api/learning-paths/completed-chapters`
- **认证**: **需要**

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | string | 否 | 用户 ID |
| path_id | int | 是 | 学习路径 ID |

**成功响应** (200 OK):

```json
{
  "completed_chapters": [1, 2, 3, 4, 5]
}
```

---

### 保存学习进度

保存用户的学习进度。

**接口信息**:
- **路径**: `POST /api/learning-paths/save-progress`
- **认证**: **需要**

**请求参数**:

```json
{
  "user_id": "user-123",
  "session_id": "session-abc",
  "path_id": 1,
  "chapter_id": 5,
  "status": "completed",
  "time_spent": 60,
  "notes": "这一章讲得很好"
}
```

**请求字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | string | 是 | 用户 ID |
| session_id | string | 否 | 会话 ID（未登录时使用） |
| path_id | int | 是 | 学习路径 ID |
| chapter_id | int | 是 | 章节 ID |
| status | string | 是 | 进度状态：in_progress/completed |
| time_spent | int | 否 | 花费时间（分钟） |
| notes | string | 否 | 学习笔记 |

**成功响应** (200 OK):

```json
{
  "success": true,
  "message": "进度已保存"
}
```

---

## 综合接口

### 获取路径学习仪表盘

获取学习路径的综合统计信息，包括路径详情、章节列表和完成进度。

**接口信息**:
- **路径**: `GET /api/learning-paths/dashboard`
- **认证**: 不需要

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | string | 否 | 用户 ID（用于获取个人进度） |
| path_id | int | 是 | 学习路径 ID |

**成功响应** (200 OK):

```json
{
  "path": {
    "id": 1,
    "title": "Python 入门教程",
    "slug": "python-basics",
    "icon": "🐍",
    "description": "从零开始学习 Python 编程",
    "difficulty": "beginner",
    "estimated_days": 30,
    "estimated_hours": 60,
    "chapter_count": 12
  },
  "chapters": [
    {
      "id": 1,
      "title": "Python 基础语法",
      "slug": "basic-syntax",
      "order_index": 1
    },
    {
      "id": 2,
      "title": "数据类型与变量",
      "slug": "data-types",
      "order_index": 2
    }
  ],
  "completed_chapters": [1, 2],
  "progress": {
    "total_chapters": 12,
    "completed_count": 2,
    "progress_percentage": 16.67
  }
}
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未授权或认证失败 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 使用示例

### 1. 获取学习路径列表

```bash
curl -X GET "http://localhost:8080/api/learning-paths?page=1&page_size=12&difficulty=beginner"
```

### 2. 获取学习路径详情

```bash
curl -X GET "http://localhost:8080/api/learning-paths/slug/python-basics"
```

### 3. 保存学习进度（需要认证）

```bash
curl -X POST "http://localhost:8080/api/learning-paths/save-progress" \
  -H "Authorization: Bearer <your_access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "path_id": 1,
    "chapter_id": 5,
    "status": "completed",
    "time_spent": 60
  }'
```

---

## 更新日志

- **v1.0.0** (2026-04-08) - 初始版本
  - 学习路径列表接口
  - 学习路径详情接口
  - 章节管理接口
  - 学习进度管理接口
  - 推荐路径接口
