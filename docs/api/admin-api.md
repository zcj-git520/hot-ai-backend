# 管理后台 API 文档

## 概述

管理后台服务提供学习路径和章节的增删改查操作。

**基础 URL**: `http://localhost:8006`

---

## 学习路径管理

### 获取学习路径列表

```
GET /api/admin/learning-paths
```

**Query 参数**:

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 20 | 每页数量 |
| difficulty | string | 否 | - | 难度筛选 (beginner/intermediate/advanced) |
| search | string | 否 | - | 关键词搜索 |
| status | int | 否 | - | 状态筛选 (0-待发布, 1-已发布, 2-已删除) |

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [...],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 获取学习路径详情

```
GET /api/admin/learning-paths/{id}
```

**Path 参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 学习路径 ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "Python 入门教程",
    "slug": "python-rumen",
    "icon": "🐍",
    "description": "从零开始学习 Python 编程",
    "difficulty": "beginner",
    "level_label": "入门",
    "learning_goals": ["掌握 Python 基础语法"],
    "target_audience": ["编程初学者"],
    "estimated_days": 30,
    "estimated_hours": 60,
    "chapter_count": 12,
    "is_featured": true,
    "status": 1,
    "chapters": [...]
  }
}
```

---

### 创建学习路径

```
POST /api/admin/learning-paths
```

**请求体**:
```json
{
  "title": "Python 入门教程",
  "icon": "🐍",
  "description": "从零开始学习 Python 编程",
  "difficulty": "beginner",
  "level_label": "入门",
  "learning_goals": ["掌握 Python 基础语法", "理解面向对象编程"],
  "target_audience": ["编程初学者", "转行开发者"],
  "estimated_days": 30,
  "estimated_hours": 60,
  "cover_image": "https://example.com/cover.jpg",
  "is_featured": true
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 标题 |
| icon | string | 否 | 图标 emoji |
| description | string | 否 | 描述 |
| difficulty | string | 否 | 难度等级 (beginner/intermediate/advanced) |
| level_label | string | 否 | 难度标签 |
| learning_goals | array | 否 | 学习目标 |
| target_audience | array | 否 | 目标受众 |
| estimated_days | int | 否 | 预计天数 |
| estimated_hours | int | 否 | 预计小时数 |
| cover_image | string | 否 | 封面图片 URL |
| is_featured | bool | 否 | 是否推荐 |

---

### 更新学习路径

```
PUT /api/admin/learning-paths/{id}
```

**请求体**: 同创建学习路径，所有字段可选

---

### 删除学习路径

```
DELETE /api/admin/learning-paths/{id}
```

软删除，将状态改为 2

---

### 发布学习路径

```
POST /api/admin/learning-paths/{id}/publish
```

将状态改为 1（已发布）

---

### 下架学习路径

```
POST /api/admin/learning-paths/{id}/unpublish
```

将状态改为 0（待发布）

---

### 设置推荐状态

```
POST /api/admin/learning-paths/{id}/featured
```

**请求体**:
```json
{
  "featured": true
}
```

---

## 章节管理

### 获取章节列表

```
GET /api/admin/learning-paths/{path_id}/chapters
```

**Path 参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| path_id | int | 是 | 学习路径 ID |

---

### 获取章节详情

```
GET /api/admin/chapters/{id}
```

**Path 参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 章节 ID |

---

### 创建章节

```
POST /api/admin/learning-paths/{path_id}/chapters
```

**请求体**:
```json
{
  "title": "Python 基础语法",
  "description": "学习 Python 的基础语法",
  "content_type": "article",
  "content": "章节内容...",
  "video_url": "",
  "external_links": [],
  "estimated_hours": 2,
  "is_free": true
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 章节标题 |
| description | string | 否 | 章节描述 |
| content_type | string | 否 | 内容类型 (article/video/exercise) |
| content | string | 否 | 章节内容 |
| video_url | string | 否 | 视频 URL |
| external_links | array | 否 | 外部链接 |
| estimated_hours | number | 否 | 预计学习小时数 |
| is_free | bool | 否 | 是否免费 |

---

### 更新章节

```
PUT /api/admin/chapters/{id}
```

**请求体**: 同创建章节，所有字段可选

---

### 删除章节

```
DELETE /api/admin/chapters/{id}
```

软删除，将状态改为 2

---

## 通用响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {...}
}
```

**错误码**:

| code | 说明 |
|------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 数据模型

### LearningPath

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 学习路径 ID |
| title | string | 标题 |
| slug | string | URL 友好标识 |
| icon | string | 图标 |
| description | string | 描述 |
| difficulty | string | 难度等级 |
| level_label | string | 难度标签 |
| learning_goals | array | 学习目标 |
| target_audience | array | 目标受众 |
| estimated_days | int | 预计天数 |
| estimated_hours | int | 预计小时数 |
| chapter_count | int | 章节数量 |
| cover_image | string | 封面图片 |
| is_featured | bool | 是否推荐 |
| status | int | 状态 (0-待发布, 1-已发布, 2-已删除) |
| published_at | datetime | 发布时间 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### PathChapter

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 章节 ID |
| path_id | int | 所属学习路径 ID |
| title | string | 章节标题 |
| slug | string | 章节 slug |
| description | string | 章节描述 |
| content_type | string | 内容类型 (article/video/exercise) |
| content | string | 章节内容 |
| video_url | string | 视频 URL |
| external_links | array | 外部链接 |
| estimated_hours | number | 预计学习小时数 |
| order_index | int | 排序索引 |
| is_free | bool | 是否免费 |
| status | int | 状态 (1-正常, 2-已删除) |
