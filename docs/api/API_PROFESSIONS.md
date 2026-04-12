# 职业风险模块 API 接口文档

## 基础信息

- **Base URL**: `http://localhost:8002/api` (profession-svc 微服务)
- **Content-Type**: `application/json`
- **认证**: 公开接口,无需登录

---

## 📋 接口列表

### 1. 获取职业分类列表

**请求**
```http
GET /api/professions/categories
```

**响应** `200 OK`
```json
[
  {
    "id": 1,
    "name": "技术类",
    "description": "技术开发、运维等相关职业",
    "sort_order": 1,
    "status": 1
  },
  {
    "id": 2,
    "name": "设计类",
    "description": "UI/UX设计、平面设计等相关职业",
    "sort_order": 2,
    "status": 1
  }
]
```

---

### 2. 获取风险等级信息

**请求**
```http
GET /api/professions/risk-levels
```

**响应** `200 OK`
```json
[
  {
    "id": "extreme",
    "level": "extreme",
    "name": "极高风险",
    "icon": "🔴",
    "description": "该职业面临极高的自动化替代风险",
    "color": "#EF4444",
    "min_score": 80,
    "max_score": 100
  },
  {
    "id": "high",
    "level": "high",
    "name": "高风险",
    "icon": "🟠",
    "description": "该职业面临较高的自动化替代风险",
    "color": "#F97316",
    "min_score": 60,
    "max_score": 79
  },
  {
    "id": "medium",
    "level": "medium",
    "name": "中风险",
    "icon": "🟡",
    "description": "该职业面临中等的自动化替代风险",
    "color": "#EAB308",
    "min_score": 40,
    "max_score": 59
  },
  {
    "id": "low",
    "level": "low",
    "name": "低风险",
    "icon": "🟢",
    "description": "该职业面临较低的自动化替代风险",
    "color": "#22C55E",
    "min_score": 0,
    "max_score": 39
  }
]
```

---

### 3. 获取职业列表

**请求**
```http
GET /api/professions?page=1&page_size=12&category_id=1&risk_level=high&keyword=程序员
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| page | int | 否 | 页码(从1开始) | 1 |
| page_size | int | 否 | 每页数量(1-100) | 12 |
| category_id | uint | 否 | 分类ID | 全部 |
| risk_level | string | 否 | 风险等级(extreme/high/medium/low) | 全部 |
| keyword | string | 否 | 关键词搜索(匹配职业名称和描述) | - |

**响应** `200 OK`
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "professions": [
      {
        "id": 1,
        "name": "软件工程师",
        "slug": "software-engineer",
        "icon": "💻",
        "category_id": 1,
        "category_name": "技术类",
        "description": "从事软件开发、维护的专业技术人员",
        "risk_level": "medium",
        "risk_score": 55,
        "automation_rate": 50,
        "view_count": 1234,
        "created_at": "2026-04-08T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 12
  }
}
```

---

### 4. 根据 slug 获取职业详情

**请求**
```http
GET /api/professions/software-engineer
```

**路径参数**
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| slug | string | 是 | 职业唯一标识(slug) |

**响应** `200 OK`
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "软件工程师",
    "slug": "software-engineer",
    "icon": "💻",
    "category_id": 1,
    "category_name": "技术类",
    "description": "从事软件开发、维护的专业技术人员",
    "risk_level": "medium",
    "risk_score": 55,
    "automation_rate": 50,
    "view_count": 1235,
    "impact_analysis": {
      "affected_tasks": ["代码生成", "单元测试编写", "Bug修复"],
      "safe_tasks": ["系统架构设计", "需求分析", "技术选型"],
      "safe_skills": ["系统设计能力", "业务理解能力", "团队协作能力"],
      "impact_timeline": {
        "short_term": "AI辅助编程工具将提高开发效率",
        "mid_term": "部分重复性编码工作将被自动化",
        "long_term": "需要向系统设计和架构方向转型"
      },
      "impact_summary": "AI将在短期内提升开发效率，但长期来看，软件工程师需要向更高层次的系统设计和业务理解方向发展"
    },
    "transition_advice": {
      "transition_paths": ["系统架构师", "技术经理", "产品经理"],
      "recommended_skills": ["云原生架构", "分布式系统设计", "业务领域知识"],
      "recommended_tools": ["Kubernetes", "Terraform", "Docker"],
      "recommended_paths": [
        {
          "title": "系统架构师路径",
          "steps": ["学习分布式系统", "掌握云原生技术", "积累架构设计经验"]
        }
      ],
      "related_articles": [1, 2, 3],
      "advice_summary": "建议向系统架构和技术管理方向发展，加强业务理解能力"
    },
    "market_data": {
      "market_trend": "growing",
      "market_trend_description": "数字化转型推动软件工程师需求持续增长",
      "salary_impact": "positive",
      "salary_change_rate": 8.5,
      "avg_salary": 25000,
      "job_demand_trend": "上升",
      "supply_demand_ratio": 1.2,
      "data_source": "拉勾网、BOSS直聘",
      "data_update_date": "2026-04-01"
    },
    "created_at": "2026-04-08T10:00:00Z",
    "updated_at": "2026-04-08T10:00:00Z"
  }
}
```

**注意**: 
- 每次访问会自动增加 `view_count`
- 详情包含完整的风险分析、转型建议和市场数据

---

### 5. 搜索职业

**请求**
```http
GET /api/professions/search?q=程序员&page=1&page_size=12
```

**查询参数**
| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| q | string | 是 | 搜索关键词 | - |
| page | int | 否 | 页码(从1开始) | 1 |
| page_size | int | 否 | 每页数量(1-100) | 12 |

**响应** `200 OK`
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "professions": [
      {
        "id": 1,
        "name": "软件工程师",
        "slug": "software-engineer",
        "icon": "💻",
        "category_id": 1,
        "category_name": "技术类",
        "description": "从事软件开发、维护的专业技术人员",
        "risk_level": "medium",
        "risk_score": 55,
        "automation_rate": 50,
        "view_count": 1234
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 12
  }
}
```

---

## 🔧 前端调用示例

### Vue/Nuxt 示例

```typescript
// 1. 获取分类列表
const categories = await $fetch('/api/professions/categories')

// 2. 获取风险等级
const riskLevels = await $fetch('/api/professions/risk-levels')

// 3. 获取职业列表(带筛选)
const response = await $fetch('/api/professions', {
  query: {
    page: 1,
    page_size: 12,
    category_id: 1,
    risk_level: 'high'
  }
})

// 4. 获取职业详情
const profession = await $fetch('/api/professions/software-engineer')

// 5. 搜索职业
const searchResult = await $fetch('/api/professions/search', {
  query: {
    q: '程序员',
    page: 1,
    page_size: 12
  }
})
```

### React 示例

```typescript
// 使用 fetch
const response = await fetch('/api/professions?page=1&page_size=12')
const data = await response.json()

// 使用 axios
import axios from 'axios'
const { data } = await axios.get('/api/professions', {
  params: { page: 1, page_size: 12, risk_level: 'high' }
})
```

---

## 📊 数据表关系

```
professions (职业主表)
  ├─ category_id → profession_categories.id (分类)
  ├─ profession_impact_analysis (影响分析, 1:1)
  │   ├─ affected_tasks (JSON数组: 受影响的任务)
  │   ├─ safe_tasks (JSON数组: 安全的任务)
  │   ├─ safe_skills (JSON数组: 安全技能)
  │   └─ impact_timeline (JSON对象: 影响时间线)
  ├─ profession_transition_advice (转型建议, 1:1)
  │   ├─ transition_paths (JSON数组: 转型路径)
  │   ├─ recommended_skills (JSON数组: 推荐技能)
  │   ├─ recommended_tools (JSON数组: 推荐工具)
  │   └─ recommended_paths (JSON数组: 推荐学习路径)
  └─ profession_market_data (市场数据, 1:1)
      ├─ market_trend (enum: growing/stable/declining)
      ├─ salary_impact (enum: positive/neutral/negative)
      └─ supply_demand_ratio (供需比)
```

---

## ⚠️ 错误响应

**400 Bad Request**
```json
{
  "code": 400,
  "message": "缺少搜索关键词",
  "data": null
}
```

**404 Not Found**
```json
{
  "code": 404,
  "message": "职业不存在",
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

---

## 🎨 风险等级说明

| 等级 | 分数范围 | 颜色 | 图标 | 说明 |
|------|---------|------|------|------|
| 极高风险 | 80-100 | #EF4444 | 🔴 | AI替代风险极高,建议尽快转型 |
| 高风险 | 60-79 | #F97316 | 🟠 | AI替代风险较高,需要关注技能提升 |
| 中风险 | 40-59 | #EAB308 | 🟡 | AI替代风险中等,建议持续学习新技能 |
| 低风险 | 0-39 | #22C55E | 🟢 | AI替代风险较低,但仍需保持学习 |
