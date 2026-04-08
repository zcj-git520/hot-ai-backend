# 学习路径微服务 API 文档

> 版本：v1.0
> 服务端口：8003
> 最后更新：2026-04-08

## 目录

1. [学习路径管理](#1-学习路径管理)
2. [章节管理](#2-章节管理)
3. [学习进度管理](#3-学习进度管理)
4. [错误处理](#4-错误处理)

---

## 1. 学习路径管理

### 1.1 获取学习路径列表

**接口描述**：分页获取所有学习路径列表，支持按难度筛选

**HTTP 请求**:
```
GET /api/learning-paths?page=1&page_size=12&difficulty=all
```

**Query 参数**:
| 参数 | 类型 | 是否必填 | 默认值 | 说明 |
|------|------|----------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 12 | 每页数量 |
| difficulty | string | 否 | all | 难度等级：all/beginner/intermediate/advanced |

**响应示例**:
```json
{
  "list": [
    {
      "id": 1,
      "title": "Python 从入门到精通",
      "slug": "python-in-depth",
      "icon": "🐍",
      "description": "系统学习 Python 编程语言",
      "difficulty": "beginner",
      "level_label": "入门",
      "learning_goals": ["掌握 Python 基础语法", "理解编程思想"],
      "target_audience": ["编程零基础学员", "想转行 IT 的职场人士"],
      "estimated_days": 30,
      "estimated_hours": 60,
      "chapter_count": 8,
      "student_count": 1250,
      "cover_image": "https://example.com/python-cover.png",
      "is_featured": 1,
      "is_active": 1,
      "sort_order": 10,
      "status": 1,
      "published_at": "2026-04-01T10:00:00Z",
      "created_at": "2026-04-01T10:00:00Z",
      "updated_at": "2026-04-01T10:00:00Z"
    }
  ],
  "total": 12,
  "page": 1
}
```

**状态码**:
- 200: 成功
- 500: 服务器错误

---

### 1.2 根据 ID 获取路径详情

**接口描述**：获取指定学习路径的完整信息，包括章节列表

**HTTP 请求**:
```
GET /api/learning-paths/:id
```

**路径参数**:
| 参数 | 类型 | 是否必填 | 说明 |
|------|------|----------|------|
| id | uint | 是 | 路径 ID |

**响应示例**:
```json
{
  "id": 1,
  "title": "Python 从入门到精通",
  "slug": "python-in-depth",
  "icon": "🐍",
  "description": "系统学习 Python 编程语言",
  "difficulty": "beginner",
  "level_label": "入门",
  "learning_goals": ["掌握 Python 基础语法", "理解编程思想"],
  "target_audience": ["编程零基础学员", "想转行 IT 的职场人士"],
  "chapters": [
    {
      "id": 1,
      "path_id": 1,
      "title": "Python 简介与环境搭建",
      "slug": "python-intro",
      "description": "介绍 Python 的历史和应用场景",
      "content_type": "article",
      "content": "Python 是一种解释型、面向对象...",
      "estimated_hours": 2,
      "order_index": 0,
      "is_free": 1,
      "status": 1
    }
  ],
  "estimated_days": 30,
  "estimated_hours": 60,
  "chapter_count": 8,
  "student_count": 1250
}
```

---

### 1.3 根据 Slug 获取路径详情

**接口描述**：使用 URL 友好的 slug 获取路径详情

**HTTP 请求**:
```
GET /api/learning-paths/slug/:slug
```

**路径参数**:
| 参数 | 类型 | 是否必填 | 说明 |
|------|------|----------|------|
| slug | string | 是 | 路径 slug |

**响应示例**：同上

---

### 1.4 获取推荐路径

**接口描述**：获取推荐的学习路径列表

**HTTP 请求**:
```
GET /api/learning-paths/featured?limit=4
```

**Query 参数**:
| 参数 | 类型 | 是否必填 | 默认值 | 说明 |
|------|------|----------|--------|------|
| limit | int | 否 | 4 | 返回数量 |

**响应示例**:
```json
[
  {
    "id": 1,
    "title": "Python 从入门到精通",
    "slug": "python-in-depth",
    "icon": "🐍",
    "difficulty": "beginner",
    "level_label": "入门",
    "estimated_days": 30,
    "student_count": 1250,
    "chapter_count": 8,
    "is_featured": 1
  }
]
```

---

### 1.5 获取难度等级信息

**接口描述**：获取所有难度等级的详细信息

**HTTP 请求**:
```
GET /api/learning-paths/levels
```

**响应示例**:
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

## 2. 章节管理

### 2.1 获取路径的所有章节

**接口描述**：获取指定学习路径的所有章节列表

**HTTP 请求**:
```
GET /api/learning-paths/:path_id/chapters
```

**路径参数**:
| 参数 | 类型 | 是否必填 | 说明 |
|------|------|----------|------|
| path_id | uint | 是 | 路径 ID |

**响应示例**:
```json
[
  {
    "id": 1,
    "path_id": 1,
    "title": "Python 简介与环境搭建",
    "slug": "python-intro",
    "description": "介绍 Python 的历史和应用场景",
    "content_type": "article",
    "content": "Python 是一种解释型、面向对象...",
    "video_url": null,
    "external_links": null,
    "estimated_hours": 2,
    "order_index": 0,
    "is_free": 1,
    "status": 1,
    "created_at": "2026-04-01T10:00:00Z",
    "updated_at": "2026-04-01T10:00:00Z"
  }
]
```

---

### 2.2 根据章节 ID 获取详情

**接口描述**：获取指定章节的详细信息

**HTTP 请求**:
```
GET /api/learning-paths/chapters/:chapter_id
```

**路径参数**:
| 参数 | 类型 | 是否必填 | 说明 |
|------|------|----------|------|
| chapter_id | uint | 是 | 章节 ID |

**响应示例**:
```json
{
  "id": 1,
  "path_id": 1,
  "title": "Python 简介与环境搭建",
  "slug": "python-intro",
  "description": "介绍 Python 的历史和应用场景",
  "content_type": "article",
  "content": "Python 是一种解释型、面向对象...",
  "video_url": null,
  "external_links": [
    {
      "title": "Python 官网",
      "url": "https://www.python.org"
    }
  ],
  "estimated_hours": 2,
  "order_index": 0,
  "is_free": 1,
  "status": 1
}
```

---

### 2.3 根据 Slug 获取章节详情

**接口描述**：使用路径 slug 和章节 slug 获取章节详情，同时返回前一章和下一章信息

**HTTP 请求**:
```
GET /api/learning-paths/:path_slug/:chapter_slug?path_id=1
```

**路径参数**:
| 参数 | 类型 | 是否必填 | 说明 |
|------|------|----------|------|
| path_slug | string | 是 | 路径 slug |
| chapter_slug | string | 是 | 章节 slug |

**Query 参数**:
| 参数 | 类型 | 是否必填 | 说明 |
|------|------|----------|------|
| path_id | uint | 否 | 路径 ID（如果提供则优先使用） |

**响应示例**:
```json
{
  "chapter": {
    "id": 1,
    "path_id": 1,
    "title": "Python 简介与环境搭建",
    "slug": "python-intro",
    "description": "介绍 Python 的历史和应用场景",
    "content_type": "article",
    "content": "Python 是一种解释型、面向对象...",
    "estimated_hours": 2,
    "order_index": 0,
    "is_free": 1
  },
  "prev": null,
  "next": {
    "id": 2,
    "path_id": 1,
    "title": "变量和数据类型",
    "slug": "variables-and-types",
    "estimated_hours": 3,
    "order_index": 1
  }
}
```

---

## 3. 学习进度管理

### 3.1 获取用户学习进度

**接口描述**：获取用户在指定路径的学习进度

**HTTP 请求**:
```
GET /api/learning-paths/progress?user_id=xxx&path_id=1
```

**Query 参数**:
| 参数 | 类型 | 是否必填 | 说明 |
|------|------|----------|------|
| user_id | string | 否 | 用户 ID |
| path_id | uint | 是 | 路径 ID |

**响应示例**:
```json
{
  "id": 1,
  "user_id": "user123",
  "path_id": 1,
  "chapter_id": 1,
  "status": "in_progress",
  "completed_at": null,
  "time_spent": 90,
  "notes": "理解了基础概念",
  "created_at": "2026-04-01T10:00:00Z",
  "updated_at": "2026-04-01T10:00:00Z"
}
```

---

### 3.2 获取已完成的章节列表

**接口描述**：获取用户在指定路径中已完成的章节 ID 列表

**HTTP 请求**:
```
GET /api/learning-paths/completed-chapters?user_id=xxx&path_id=1
```

**Query 参数**:
| 参数 | 类型 | 是否必填 | 说明 |
|------|------|----------|------|
| user_id | string | 否 | 用户 ID |
| path_id | uint | 是 | 路径 ID |

**响应示例**:
```json
{
  "completed_chapters": [1, 2, 3, 5]
}
```

---

### 3.3 保存学习进度

**接口描述**：保存用户的学习进度，支持章节完成状态和学习笔记

**HTTP 请求**:
```
POST /api/learning-paths/progress
Content-Type: application/json
```

**请求体**:
```json
{
  "user_id": "user123",
  "session_id": "session456",
  "path_id": 1,
  "chapter_id": 1,
  "status": "completed",
  "time_spent": 90,
  "notes": "理解了基础概念，需要多加练习"
}
```

**字段说明**:
| 字段 | 类型 | 是否必填 | 说明 |
|------|------|----------|------|
| user_id | string | 否 | 用户 ID（登录用户） |
| session_id | string | 否 | 会话 ID（未登录用户） |
| path_id | uint | 是 | 路径 ID |
| chapter_id | uint | 是 | 章节 ID |
| status | enum | 是 | 状态：in_progress/completed |
| time_spent | int | 否 | 学习时长（分钟） |
| notes | string | 否 | 学习笔记 |

**响应示例**:
```json
{
  "success": true,
  "message": "进度已保存"
}
```

**状态码**:
- 200: 成功
- 400: 请求参数错误
- 500: 服务器错误

---

### 3.4 获取学习仪表盘

**接口描述**：获取路径学习的综合统计信息，包括路径详情、章节列表和完成进度

**HTTP 请求**:
```
GET /api/learning-paths/dashboard?user_id=xxx&path_id=1
```

**Query 参数**:
| 参数 | 类型 | 是否必填 | 说明 |
|------|------|----------|------|
| user_id | string | 否 | 用户 ID |
| path_id | uint | 是 | 路径 ID |

**响应示例**:
```json
{
  "path": {
    "id": 1,
    "title": "Python 从入门到精通",
    "slug": "python-in-depth",
    "difficulty": "beginner",
    "chapter_count": 8
  },
  "chapters": [
    {
      "id": 1,
      "title": "Python 简介与环境搭建",
      "estimated_hours": 2
    }
  ],
  "completed_chapters": [1, 2, 3],
  "progress": {
    "total_chapters": 8,
    "completed_count": 3,
    "progress_percentage": 37.5
  }
}
```

---

## 4. 错误处理

### 错误响应格式

所有错误响应都遵循以下格式：

```json
{
  "code": 400,
  "message": "缺少路径ID",
  "details": "path_id is required"
}
```

### 错误码说明

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权，需要登录 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
| 503 | 服务不可用 |

### 常见错误场景

1. **缺少必要参数**
   ```json
   {
     "code": 400,
     "message": "缺少路径ID",
     "details": "path_id is required"
   }
   ```

2. **资源不存在**
   ```json
   {
     "code": 404,
     "message": "路径不存在",
     "details": "learning path not found"
   }
   ```

3. **服务器错误**
   ```json
   {
     "code": 500,
     "message": "服务器内部错误",
     "details": "database connection failed"
   }
   ```

---

## 5. 使用示例

### 5.1 获取所有学习路径

```bash
curl -X GET http://localhost:8003/api/learning-paths?page=1&page_size=12
```

### 5.2 开始学习路径

```bash
# 1. 获取路径详情
curl -X GET http://localhost:8003/api/learning-paths/slug/python-in-depth

# 2. 开始学习（记录开始学习）
curl -X POST -H "Content-Type: application/json" \
  http://localhost:8003/api/learning-paths/progress \
  -d '{
    "user_id": "user123",
    "path_id": 1,
    "chapter_id": 1,
    "status": "in_progress"
  }'
```

### 5.3 完成章节学习

```bash
curl -X POST -H "Content-Type: application/json" \
  http://localhost:8003/api/learning-paths/progress \
  -d '{
    "user_id": "user123",
    "path_id": 1,
    "chapter_id": 1,
    "status": "completed",
    "time_spent": 90,
    "notes": "理解了基础概念"
  }'
```

### 5.4 查看学习进度

```bash
# 查看仪表盘
curl -X GET http://localhost:8003/api/learning-paths/dashboard?user_id=user123\u0026path_id=1

# 查看已完成的章节
curl -X GET http://localhost:8003/api/learning-paths/completed-chapters?user_id=user123\u0026path_id=1
```

---

## 6. 最佳实践

### 6.1 前端集成建议

1. **路径列表页**：使用 `/api/learning-paths` 获取列表，支持筛选和分页
2. **路径详情页**：使用 `/api/learning-paths/slug/:slug` 获取详情和章节列表
3. **学习页面**：
   - 使用 `/api/learning-paths/:path_slug/:chapter_slug` 获取章节内容
   - 定时调用 `/api/learning-paths/progress` 保存学习进度（每 30 秒）
   - 章节完成时调用保存进度接口，status 设为 "completed"
4. **仪表盘**：使用 `/api/learning-paths/dashboard` 获取综合学习数据

### 6.2 性能优化

1. **缓存策略**：
   - 路径列表和详情可以缓存 5-10 分钟
   - 用户进度数据建议不缓存或缓存 1 分钟

2. **分页加载**：
   - 路径列表使用分页，每页 12-20 条
   - 章节列表一般不需要分页（单路径章节数量有限）

3. **懒加载**：
   - 章节内容可以懒加载，进入章节时再请求
   - 进度数据可以在用户登录后加载

### 6.3 异常处理

1. **网络错误**：实现重试机制，最多重试 3 次
2. **404 错误**：提示用户"内容不存在或已被删除"
3. **500 错误**：提示用户"服务器繁忙，请稍后重试"
4. **进度保存失败**：本地缓存进度，待网络恢复后同步

---

## 7. 更新日志

### v1.0 (2026-04-08)
- ✨ 初始版本发布
- ✨ 实现学习路径、章节、进度管理功能
- ✨ 支持用户认证和无状态会话
- ✨ 提供完整的API文档

---

## 8. 相关文档

- [数据库表结构](../migrations/005_create_learning_paths.sql)
- [后端架构说明](./MICROSERVICES.md)
- [部署文档](./DEPLOYMENT.md)
- [Nginx配置](./NGINX_GUIDE.md)

---

**文档维护者**：AI热点追踪平台团队
**最后更新**：2026-04-08
