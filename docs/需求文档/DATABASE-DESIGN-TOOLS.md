# 工具库模块 - 数据库设计文档

**版本**: v1.0
**创建日期**: 2026-04-11
**数据库**: MySQL 8.0+

---

## 目录

1. [表结构设计](#1-表结构设计)
2. [索引设计](#2-索引设计)
3. [初始化数据](#3-初始化数据)
4. [数据迁移](#4-数据迁移)
5. [ER图](#5-er图)

---

## 1. 表结构设计

### 1.1 工具表 (tools)

```sql
CREATE TABLE `tools` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '工具ID，主键',
  `name` VARCHAR(200) NOT NULL COMMENT '工具名称',
  `slug` VARCHAR(200) NOT NULL COMMENT 'URL友好的标识，唯一',
  `icon` VARCHAR(500) DEFAULT NULL COMMENT '工具图标URL',
  `description` TEXT COMMENT '工具描述',
  `official_url` VARCHAR(500) DEFAULT NULL COMMENT '官方网站',
  `documentation_url` VARCHAR(500) DEFAULT NULL COMMENT '文档链接',
  `pricing` JSON DEFAULT NULL COMMENT '定价信息，JSON格式',
  `pricing_description` VARCHAR(1000) DEFAULT NULL COMMENT '定价说明',
  `category_id` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属类别ID，0表示未分类',
  `difficulty` VARCHAR(20) DEFAULT 'beginner' COMMENT '难度等级：beginner/intermediate/advanced',
  `rating` DECIMAL(2,1) DEFAULT 0.00 COMMENT '平均评分，范围0-5',
  `review_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '评测数量',
  `view_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '浏览量',
  `popularity` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '热度值',
  `tags` JSON DEFAULT NULL COMMENT '标签列表，JSON数组',
  `featured` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否精选展示',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：1-上架，0-下架',
  `external_id` VARCHAR(100) DEFAULT NULL COMMENT '外部系统ID（如爬虫数据源ID）',
  `created_by` VARCHAR(50) DEFAULT NULL COMMENT '创建者ID',
  `updated_by` VARCHAR(50) DEFAULT NULL COMMENT '最后更新者ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_slug` (`slug`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_rating` (`rating`),
  KEY `idx_popularity` (`popularity`),
  KEY `idx_status` (`status`),
  KEY `idx_difficulty` (`difficulty`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_external_id` (`external_id`),
  KEY `idx_status_created_at` (`status`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI工具表';
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT UNSIGNED | 主键，自增 |
| name | VARCHAR(200) | 工具名称 |
| slug | VARCHAR(200) | URL友好的标识，唯一 |
| icon | VARCHAR(500) | 图标URL（可选，支持emoji或图片） |
| description | TEXT | 工具描述 |
| official_url | VARCHAR(500) | 官方网站 |
| documentation_url | VARCHAR(500) | 文档链接 |
| pricing | JSON | 定价信息（JSON结构） |
| pricing_description | VARCHAR(1000) | 定价说明 |
| category_id | INT UNSIGNED | 所属类别ID |
| difficulty | VARCHAR(20) | 难度等级 |
| rating | DECIMAL(2,1) | 平均评分（0.0-5.0） |
| review_count | INT UNSIGNED | 评测数量 |
| view_count | INT UNSIGNED | 浏览量 |
| popularity | INT UNSIGNED | 热度值 |
| tags | JSON | 标签列表 |
| featured | BOOLEAN | 是否精选 |
| status | TINYINT UNSIGNED | 状态（1-上架，0-下架） |
| external_id | VARCHAR(100) | 外部ID（爬虫数据源） |
| created_by | VARCHAR(50) | 创建者 |
| updated_by | VARCHAR(50) | 更新者 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | 软删除时间 |

**pricing JSON 结构**：
```json
{
  "free": {
    "available": true,
    "limit": "每天消息限额",
    "features": ["基础对话", "基础写作"]
  },
  "subscription": {
    "plans": [
      {
        "name": "基础版",
        "price": "¥0/月",
        "features": ["无限对话", "高级功能"]
      },
      {
        "name": "专业版",
        "price": "¥99/月",
        "features": ["GPT-4访问", "优先支持"]
      }
    ]
  },
  "usage_based": {
    "description": "按使用量付费",
    "example": "每次请求 ¥0.01"
  }
}
```

---

### 1.2 工具类别表 (tool_categories)

```sql
CREATE TABLE `tool_categories` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '类别ID，主键',
  `name` VARCHAR(100) NOT NULL COMMENT '类别名称',
  `slug` VARCHAR(100) NOT NULL COMMENT 'URL友好的标识，唯一',
  `icon` VARCHAR(500) DEFAULT NULL COMMENT '图标',
  `description` VARCHAR(1000) DEFAULT NULL COMMENT '类别描述',
  `sort_order` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `featured` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否精选展示',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_slug` (`slug`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工具类别表';
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT UNSIGNED | 主键，自增 |
| name | VARCHAR(100) | 类别名称 |
| slug | VARCHAR(100) | URL友好的标识，唯一 |
| icon | VARCHAR(500) | 图标 |
| description | VARCHAR(1000) | 类别描述 |
| sort_order | INT UNSIGNED | 排序顺序 |
| featured | BOOLEAN | 是否精选 |
| status | TINYINT UNSIGNED | 状态（1-启用，0-禁用） |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | 软删除时间 |

---

### 1.3 工具标签表 (tool_tags)

```sql
CREATE TABLE `tool_tags` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '标签ID，主键',
  `name` VARCHAR(100) NOT NULL COMMENT '标签名称',
  `slug` VARCHAR(100) NOT NULL COMMENT 'URL友好的标识，唯一',
  `color` VARCHAR(7) DEFAULT NULL COMMENT '标签颜色（十六进制）',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_slug` (`slug`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工具标签表';
```

---

### 1.4 工具-标签关联表 (tool_tag_relations)

```sql
CREATE TABLE `tool_tag_relations` (
  `tool_id` BIGINT UNSIGNED NOT NULL COMMENT '工具ID',
  `tag_id` INT UNSIGNED NOT NULL COMMENT '标签ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

  PRIMARY KEY (`tool_id`, `tag_id`),
  KEY `idx_tag_id` (`tag_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_tool_tag_tool` FOREIGN KEY (`tool_id`) REFERENCES `tools` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_tool_tag_tag` FOREIGN KEY (`tag_id`) REFERENCES `tool_tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工具标签关联表';
```

---

### 1.5 用户评测表 (tool_reviews)

```sql
CREATE TABLE `tool_reviews` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评测ID，主键',
  `user_id` VARCHAR(100) NOT NULL COMMENT '用户ID（登录用户/访客）',
  `tool_id` BIGINT UNSIGNED NOT NULL COMMENT '工具ID',
  `user_ip` VARCHAR(45) DEFAULT NULL COMMENT '用户IP地址（访客评测）',
  `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '用户User-Agent（访客评测）',
  `rating` TINYINT UNSIGNED NOT NULL COMMENT '评分，1-5',
  `ease_of_use` TINYINT UNSIGNED NOT NULL COMMENT '易用性，1-5',
  `effectiveness` TINYINT UNSIGNED NOT NULL COMMENT '效果质量，1-5',
  `value_for_money` TINYINT UNSIGNED NOT NULL COMMENT '性价比，1-5',
  `features` TINYINT UNSIGNED NOT NULL COMMENT '功能丰富度，1-5',
  `update_frequency` TINYINT UNSIGNED NOT NULL COMMENT '更新频率，1-5',
  `support` TINYINT UNSIGNED NOT NULL COMMENT '客服支持，1-5',
  `pros` TEXT COMMENT '优点（JSON数组或文本）',
  `cons` TEXT COMMENT '缺点（JSON数组或文本）',
  `comment` TEXT COMMENT '详细评论',
  `images` JSON DEFAULT NULL COMMENT '图片URL列表（JSON数组）',
  `pros_json` JSON DEFAULT NULL COMMENT '优点（JSON数组）',
  `cons_json` JSON DEFAULT NULL COMMENT '缺点（JSON数组）',
  `is_anonymous` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否匿名',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：1-审核通过，0-待审核，2-已拒绝',
  `reason` VARCHAR(500) DEFAULT NULL COMMENT '拒绝原因',
  `is_verified_purchase` BOOLEAN DEFAULT NULL COMMENT '是否为已购买用户（可选）',
  `verified_at` TIMESTAMP NULL DEFAULT NULL COMMENT '验证时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',

  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_tool_id` (`tool_id`),
  KEY `idx_rating` (`rating`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_tool_status_created` (`tool_id`, `status`, `created_at`),
  KEY `idx_user_tool` (`user_id`, `tool_id`),
  CONSTRAINT `fk_review_tool` FOREIGN KEY (`tool_id`) REFERENCES `tools` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户评测表';
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT UNSIGNED | 主键，自增 |
| user_id | VARCHAR(100) | 用户ID（登录用户/访客） |
| tool_id | BIGINT UNSIGNED | 工具ID |
| user_ip | VARCHAR(45) | 用户IP（访客评测） |
| user_agent | VARCHAR(500) | User-Agent（访客评测） |
| rating | TINYINT UNSIGNED | 总评分（1-5） |
| ease_of_use | TINYINT UNSIGNED | 易用性（1-5） |
| effectiveness | TINYINT UNSIGNED | 效果质量（1-5） |
| value_for_money | TINYINT UNSIGNED | 性价比（1-5） |
| features | TINYINT UNSIGNED | 功能丰富度（1-5） |
| update_frequency | TINYINT UNSIGNED | 更新频率（1-5） |
| support | TINYINT UNSIGNED | 客服支持（1-5） |
| pros | TEXT | 优点 |
| cons | TEXT | 缺点 |
| comment | TEXT | 详细评论 |
| images | JSON | 图片URL列表 |
| pros_json | JSON | 优点（JSON数组） |
| cons_json | JSON | 缺点（JSON数组） |
| is_anonymous | BOOLEAN | 是否匿名 |
| status | TINYINT UNSIGNED | 状态（1-通过，0-待审核，2-拒绝） |
| is_verified_purchase | BOOLEAN | 是否已购买验证 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |

---

### 1.6 评论表 (comments)

```sql
CREATE TABLE `comments` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID，主键',
  `user_id` VARCHAR(100) NOT NULL COMMENT '用户ID',
  `commentable_type` VARCHAR(50) NOT NULL COMMENT '评论对象类型：tool_review/tool',
  `commentable_id` BIGINT UNSIGNED NOT NULL COMMENT '评论对象ID',
  `parent_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '父评论ID',
  `content` TEXT NOT NULL COMMENT '评论内容',
  `images` JSON DEFAULT NULL COMMENT '图片URL列表（JSON数组）',
  `is_anonymous` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否匿名',
  `likes` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
  `is_liked` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '当前用户是否已点赞',
  `is_spam` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否为垃圾评论',
  `spam_reason` VARCHAR(500) DEFAULT NULL COMMENT '标记为垃圾评论的原因',
  `spam_count` INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '垃圾标记次数',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：1-显示，0-隐藏',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',

  PRIMARY KEY (`id`),
  KEY `idx_commentable` (`commentable_type`, `commentable_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_status` (`status`),
  KEY `idx_commentable_status_created` (`commentable_type`, `commentable_id`, `status`, `created_at`),
  KEY `idx_parent_status` (`parent_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';
```

---

### 1.7 提示词模板表 (prompt_templates)

```sql
CREATE TABLE `prompt_templates` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '模板ID，主键',
  `name` VARCHAR(200) NOT NULL COMMENT '模板名称',
  `slug` VARCHAR(200) NOT NULL COMMENT 'URL友好的标识，唯一',
  `description` TEXT COMMENT '模板描述',
  `content` LONGTEXT NOT NULL COMMENT '提示词内容',
  `tool_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '适用工具ID（可为空）',
  `category_id` INT UNSIGNED DEFAULT NULL COMMENT '模板类别ID',
  `use_cases` JSON DEFAULT NULL COMMENT '使用场景（JSON数组）',
  `tags` JSON DEFAULT NULL COMMENT '标签（JSON数组）',
  `example_response` TEXT COMMENT '示例回复',
  `example_input` TEXT COMMENT '示例输入',
  `likes` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
  `views` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '浏览量',
  `favorites` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '收藏数',
  `featured` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否精选',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
  `language` VARCHAR(10) DEFAULT 'zh-CN' COMMENT '语言：zh-CN/en-US',
  `author_type` VARCHAR(20) DEFAULT 'community' COMMENT '作者类型：system/community/user',
  `author_id` VARCHAR(100) DEFAULT NULL COMMENT '作者ID',
  `approved_at` TIMESTAMP NULL DEFAULT NULL COMMENT '审核时间',
  `approved_by` VARCHAR(50) DEFAULT NULL COMMENT '审核者ID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_slug` (`slug`),
  KEY `idx_tool_id` (`tool_id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_language` (`language`),
  KEY `idx_status` (`status`),
  KEY `idx_featured` (`featured`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_author` (`author_type`, `author_id`),
  CONSTRAINT `fk_prompt_tool` FOREIGN KEY (`tool_id`) REFERENCES `tools` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_prompt_category` FOREIGN KEY (`category_id`) REFERENCES `prompt_template_categories` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='提示词模板表';
```

---

### 1.8 提示词模板分类表 (prompt_template_categories)

```sql
CREATE TABLE `prompt_template_categories` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '类别ID，主键',
  `name` VARCHAR(100) NOT NULL COMMENT '分类名称',
  `slug` VARCHAR(100) NOT NULL COMMENT 'URL友好的标识，唯一',
  `description` TEXT COMMENT '分类描述',
  `icon` VARCHAR(500) DEFAULT NULL COMMENT '图标',
  `sort_order` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `featured` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否精选',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_slug` (`slug`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='提示词模板分类表';
```

---

### 1.9 用户收藏表 (user_favorites)

```sql
CREATE TABLE `user_favorites` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '收藏ID，主键',
  `user_id` VARCHAR(100) NOT NULL COMMENT '用户ID',
  `tool_id` BIGINT UNSIGNED NOT NULL COMMENT '工具ID',
  `note` VARCHAR(500) DEFAULT NULL COMMENT '收藏备注',
  `favorite_list_id` BIGINT UNSIGNED NULL DEFAULT NULL COMMENT '收藏夹ID（可选）',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_tool` (`user_id`, `tool_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_tool_id` (`tool_id`),
  KEY `idx_favorite_list` (`favorite_list_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_user_created` (`user_id`, `created_at`),
  CONSTRAINT `fk_favorite_tool` FOREIGN KEY (`tool_id`) REFERENCES `tools` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户收藏表';
```

---

### 1.10 徽章表 (badges)

```sql
CREATE TABLE `badges` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '徽章ID，主键',
  `name` VARCHAR(100) NOT NULL COMMENT '徽章名称',
  `slug` VARCHAR(100) NOT NULL COMMENT 'URL友好的标识，唯一',
  `description` TEXT COMMENT '徽章描述',
  `icon` VARCHAR(500) DEFAULT NULL COMMENT '徽章图标',
  `type` VARCHAR(50) NOT NULL COMMENT '类型：review/contribution/social',
  `condition_type` VARCHAR(50) NOT NULL COMMENT '获取条件类型',
  `condition_value` INT NOT NULL COMMENT '获取条件值',
  `icon_color` VARCHAR(7) DEFAULT NULL COMMENT '图标颜色',
  `background_color` VARCHAR(7) DEFAULT NULL COMMENT '背景颜色',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_slug` (`slug`),
  KEY `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='徽章表';
```

**condition_type 说明**：
| 类型 | 说明 | 示例 |
|------|------|------|
| review_count | 评测数量 | 1, 10, 50, 100 |
| review_star | 评测星级 | 1, 2, 3, 4, 5 |
| template_count | 模板数量 | 5, 10, 20, 50 |
| comment_count | 评论数量 | 10, 50, 100, 200 |
| login_days | 登录天数 | 7, 30, 90, 180 |
| like_count | 点赞数量 | 100, 500, 1000 |

---

### 1.11 用户徽章表 (user_badges)

```sql
CREATE TABLE `user_badges` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户徽章ID，主键',
  `user_id` VARCHAR(100) NOT NULL COMMENT '用户ID',
  `badge_id` INT UNSIGNED NOT NULL COMMENT '徽章ID',
  `issued_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '获得时间',
  `unique_index` (`user_id`, `badge_id`),

  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_badge_id` (`badge_id`),
  KEY `idx_issued_at` (`issued_at`),
  KEY `idx_user_issued` (`user_id`, `issued_at`),
  CONSTRAINT `fk_user_badge_badge` FOREIGN KEY (`badge_id`) REFERENCES `badges` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户徽章表';
```

---

### 1.12 用户偏好表 (user_preferences)

```sql
CREATE TABLE `user_preferences` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '偏好ID，主键',
  `user_id` VARCHAR(100) NOT NULL COMMENT '用户ID',
  `key` VARCHAR(100) NOT NULL COMMENT '偏好键',
  `value` VARCHAR(500) DEFAULT NULL COMMENT '偏好值',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_key` (`user_id`, `key`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户偏好表';
```

**常见偏好键值对**：
```json
{
  "favorite_categories": "[1,2,3]",
  "favorite_difficulty": "beginner",
  "notification_enabled": "true",
  "review_notifications": "true",
  "template_notifications": "true"
}
```

---

### 1.13 工具搜索历史表 (tool_search_history)

```sql
CREATE TABLE `tool_search_history` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '搜索ID，主键',
  `user_id` VARCHAR(100) DEFAULT NULL COMMENT '用户ID（访客不记录）',
  `search_keyword` VARCHAR(200) NOT NULL COMMENT '搜索关键词',
  `category_id` INT UNSIGNED DEFAULT NULL COMMENT '筛选类别ID',
  `is_free` BOOLEAN DEFAULT NULL COMMENT '是否筛选免费工具',
  `min_rating` DECIMAL(2,1) DEFAULT NULL COMMENT '最低评分',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '搜索时间',

  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_keyword` (`search_keyword`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_user_created` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工具搜索历史表';
```

---

### 1.14 热门工具统计表 (tool_hot_stats)

```sql
CREATE TABLE `tool_hot_stats` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '统计ID，主键',
  `tool_id` BIGINT UNSIGNED NOT NULL COMMENT '工具ID',
  `stat_date` DATE NOT NULL COMMENT '统计日期',
  `daily_views` INT UNSIGNED DEFAULT 0 COMMENT '每日浏览量',
  `daily_reviews` INT UNSIGNED DEFAULT 0 COMMENT '每日新增评测',
  `daily_likes` INT UNSIGNED DEFAULT 0 COMMENT '每日新增点赞',
  `popularity_score` INT UNSIGNED DEFAULT 0 COMMENT '热度得分',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tool_date` (`tool_id`, `stat_date`),
  KEY `idx_stat_date` (`stat_date`),
  KEY `idx_tool_id` (`tool_id`),
  CONSTRAINT `fk_hot_stat_tool` FOREIGN KEY (`tool_id`) REFERENCES `tools` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='热门工具统计表';
```

---

### 1.15 通知表 (notifications)

```sql
CREATE TABLE `notifications` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '通知ID，主键',
  `user_id` VARCHAR(100) NOT NULL COMMENT '用户ID',
  `type` VARCHAR(50) NOT NULL COMMENT '通知类型：review/comment/template',
  `title` VARCHAR(200) NOT NULL COMMENT '通知标题',
  `content` TEXT COMMENT '通知内容',
  `action_url` VARCHAR(500) DEFAULT NULL COMMENT '跳转链接',
  `action_target_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '关联对象ID',
  `action_target_type` VARCHAR(50) DEFAULT NULL COMMENT '关联对象类型',
  `is_read` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否已读',
  `read_at` TIMESTAMP NULL DEFAULT NULL COMMENT '阅读时间',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_is_read` (`is_read`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_user_created` (`user_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知表';
```

---

### 1.16 系统配置表 (system_config)

```sql
CREATE TABLE `system_config` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '配置ID，主键',
  `key` VARCHAR(100) NOT NULL COMMENT '配置键，唯一',
  `value` TEXT COMMENT '配置值',
  `value_type` VARCHAR(20) DEFAULT 'string' COMMENT '值类型：string/int/boolean/json',
  `description` VARCHAR(500) DEFAULT NULL COMMENT '配置说明',
  `category` VARCHAR(50) DEFAULT 'general' COMMENT '配置分类',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_key` (`key`),
  KEY `idx_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';
```

**示例配置数据**：
```sql
INSERT INTO `system_config` (`key`, `value`, `value_type`, `description`, `category`) VALUES
('tool_review_approve', '1', 'int', '是否需要审核评测', 'review'),
('template_review_approve', '1', 'int', '是否需要审核模板', 'template'),
('enable_anonymous_review', '1', 'boolean', '是否允许匿名评测', 'review'),
('max_review_per_day', '10', 'int', '每日最多评测数', 'review'),
('max_template_per_day', '5', 'int', '每日最多提交模板数', 'template');
```

---

## 2. 索引设计

### 2.1 索引分类

| 类型 | 用途 | 示例 |
|------|------|------|
| **主键索引** | 主键唯一标识 | `id` |
| **唯一索引** | 唯一性约束 | `uk_slug` |
| **普通索引** | 加速查询 | `idx_category_id` |
| **联合索引** | 加速多条件查询 | `idx_user_tool` |
| **复合索引** | 加速排序和范围查询 | `idx_status_created_at` |

### 2.2 索引使用建议

1. **工具列表查询**：
   ```sql
   -- 推荐索引：idx_category_id + idx_rating + idx_status + idx_created_at
   SELECT * FROM tools
   WHERE category_id = 1 AND status = 1
   ORDER BY rating DESC
   LIMIT 20;
   ```

2. **用户评测查询**：
   ```sql
   -- 推荐索引：idx_tool_id + idx_status + idx_created_at
   SELECT * FROM tool_reviews
   WHERE tool_id = 1 AND status = 1
   ORDER BY created_at DESC
   LIMIT 20;
   ```

3. **搜索查询**：
   ```sql
   -- 推荐索引：idx_name + idx_status
   SELECT * FROM tools
   WHERE name LIKE '%ChatGPT%' AND status = 1
   LIMIT 20;
   ```

---

## 3. 初始化数据

### 3.1 工具类别数据

```sql
INSERT INTO `tool_categories` (`id`, `name`, `slug`, `icon`, `description`, `sort_order`, `featured`, `status`) VALUES
(1, '写作类', 'writing', '✍️', '用于写作、文案创作的 AI 工具', 1, TRUE, 1),
(2, '图像类', 'image', '🎨', '用于图像生成、编辑的 AI 工具', 2, TRUE, 1),
(3, '视频类', 'video', '🎬', '用于视频生成、编辑的 AI 工具', 3, TRUE, 1),
(4, '音频类', 'audio', '🔊', '用于音频生成、编辑的 AI 工具', 4, TRUE, 1),
(5, '代码类', 'code', '💻', '用于编程、代码相关的 AI 工具', 5, TRUE, 1),
(6, '办公类', 'office', '📊', '用于办公、文档处理的 AI 工具', 6, TRUE, 1),
(7, '其他类', 'other', '🔧', '其他 AI 工具', 7, FALSE, 1);
```

### 3.2 提示词模板分类数据

```sql
INSERT INTO `prompt_template_categories` (`id`, `name`, `slug`, `description`, `icon`, `sort_order`, `featured`, `status`) VALUES
(1, '写作类', 'writing', '用于写作的提示词模板', '✍️', 1, TRUE, 1),
(2, '代码类', 'code', '用于编程的提示词模板', '💻', 2, TRUE, 1),
(3, '设计类', 'design', '用于设计的提示词模板', '🎨', 3, TRUE, 1),
(4, '营销类', 'marketing', '用于营销的提示词模板', '📢', 4, TRUE, 1),
(5, '学习类', 'learning', '用于学习的提示词模板', '📚', 5, FALSE, 1),
(6, '翻译类', 'translation', '用于翻译的提示词模板', '🌐', 6, FALSE, 1);
```

### 3.3 徽章数据

```sql
INSERT INTO `badges` (`id`, `name`, `slug`, `description`, `icon`, `type`, `condition_type`, `condition_value`, `icon_color`, `background_color`, `status`) VALUES
-- 评测徽章
(1, '首次评测者', 'first-review', '完成第一次工具评测', '📝', 'review', 'review_count', 1, '#FFD700', '#FFF8DC', 1),
(2, '评测达人', 'review-star', '获得5星好评评测', '⭐', 'review', 'review_star', 5, '#FFD700', '#FFF8DC', 1),
(3, '精选评测', 'featured-review', '评测被标记为精选', '💎', 'review', 'featured', 1, '#FFD700', '#FFF8DC', 1),
(4, '深度评测', 'pro-review', '提交50+字详细评测', '🏆', 'review', 'review_count', 50, '#FFD700', '#FFF8DC', 1),
(5, '活跃评测', 'active-review', '每周提交3条评测', '🔥', 'review', 'weekly_review', 3, '#FF4500', '#FFE4E1', 1),
-- 创作徽章
(6, '模板达人', 'template-star', '提交5个模板', '✨', 'contribution', 'template_count', 5, '#87CEEB', '#E0FFFF', 1),
(7, '模板大师', 'template-master', '提交10个模板', '🌟', 'contribution', 'template_count', 10, '#87CEEB', '#E0FFFF', 1),
(8, '模板专家', 'template-expert', '提交20个模板', '💫', 'contribution', 'template_count', 20, '#87CEEB', '#E0FFFF', 1),
-- 社区徽章
(9, '活跃评论', 'active-comment', '发布50条评论', '💬', 'social', 'comment_count', 50, '#32CD32', '#F0FFF0', 1),
(9, '社区达人', 'community-star', '每日登录30天', '📅', 'social', 'login_days', 30, '#32CD32', '#F0FFF0', 1),
(10, '热评', 'hot-comment', '获得100个赞的评论', '👍', 'social', 'like_count', 100, '#FF6347', '#FFE4E1', 1);
```

### 3.4 系统配置数据

```sql
INSERT INTO `system_config` (`key`, `value`, `value_type`, `description`, `category`) VALUES
('tool_review_approve', '1', 'int', '是否需要审核评测', 'review'),
('template_review_approve', '1', 'int', '是否需要审核模板', 'template'),
('enable_anonymous_review', '1', 'boolean', '是否允许匿名评测', 'review'),
('max_review_per_day', '10', 'int', '每日最多评测数', 'review'),
('max_template_per_day', '5', 'int', '每日最多提交模板数', 'template'),
('enable_anonymous_template', '1', 'boolean', '是否允许匿名提交模板', 'template'),
('enable_user_badges', '1', 'boolean', '是否启用徽章系统', 'badge'),
('enable_search_history', '1', 'boolean', '是否记录搜索历史', 'search'),
('enable_tool_views', '1', 'boolean', '是否统计工具浏览量', 'tool');
```

---

## 4. 数据迁移

### 4.1 创建数据库和表

```sql
-- 创建数据库
CREATE DATABASE IF NOT EXISTS `aihot_tools` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE `aihot_tools`;

-- 执行建表SQL（见上文各表结构）
```

### 4.2 执行顺序

```
1. tool_categories
2. prompt_template_categories
3. tool_tags
4. tools
5. tool_tag_relations
6. prompt_templates
7. user_favorites
8. badges
9. user_badges
10. system_config
11. 工具类别数据
12. 提示词模板分类数据
13. 徽章数据
14. 系统配置数据
```

---

## 5. ER图

### 5.1 实体关系图

```
┌─────────────────────────────────────────────────────────────────────┐
│                         ER 图                                       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                 │
│  users  (1) ◄─────────────────────────────────────────────────► (N) tools
│       │                                                  │      │
│       │ 1:N                                            │      │
│       │                                                │      │
│       ▼                                                │      │
│  ┌──────────────────┐                               │      │
│  │ tool_reviews     │◄──────────────────────────────┘      │
│  └──────────────────┘                                        │
│       │                                                      │
│       │ 1:N                                                  │
│       ▼                                                      │
│  ┌──────────────────┐                                       │
│  │ comments         │                                       │
│  └──────────────────┘                                       │
│                                                                 │
│  tools  (1) ◄─────────────────────────────────────────────► (N) prompt_templates
│       │                                                  │      │
│       │ 1:N                                            │      │
│       │                                                │      │
│       ▼                                                │      │
│  ┌──────────────────┐                               │      │
│  │ user_favorites   │                               │      │
│  └──────────────────┘                               │      │
│                                                                 │
│  tools_categories (1) ◄───────────────── (N) tools
│                                                     │
│                                                     │
│  tool_tags (1) ◄───────────────── (N) tool_tag_relations ◄───────────────── (N) tools
│                                                                 │
│  user_preferences (1) ◄───────────────── (N) user_preferences
│                                                                 │
│  tool_search_history (1) ◄───────────────── (N) tool_search_history
│                                                                 │
│  tool_hot_stats (1) ◄───────────────── (N) tool_hot_stats
│                                                                 │
│  badges (1) ◄───────────────────────────────────────── (N) user_badges
│                     │                                                │
│                     │ 1:N                                             │
│                     ▼                                                │
│               users (N) ◄───────────────── (N) user_preferences
│                                                                 │
└─────────────────────────────────────────────────────────────────────┘
```

### 5.2 核心关系说明

1. **工具与标签**：一对多（一个工具可以有多个标签，一个标签可以属于多个工具）
2. **用户与评测**：一对多（一个用户可以写多个评测，一个评测属于一个工具）
3. **用户与收藏**：一对多（一个用户可以收藏多个工具，一个工具可以被多个用户收藏）
4. **工具与模板**：一对多（一个工具可以有多个提示词模板，一个模板可以属于一个工具）
5. **徽章与用户**：一对多（一个徽章可以颁发给多个用户，一个用户可以获得多个徽章）
6. **评论与用户**：一对多（一个用户可以写多个评论，一个评论属于一个用户）

---

## 附录

### 附录 A：数据库性能优化建议

1. **定期维护索引**：
   ```sql
   -- 分析表
   ANALYZE TABLE tools, tool_reviews, comments;

   -- 优化表
   OPTIMIZE TABLE tools, tool_reviews, comments;
   ```

2. **分区表**（后期优化）：
   - 工具表：按创建时间分区
   - 评测表：按创建时间分区
   - 评论表：按创建时间分区

3. **读写分离**（后期优化）：
   - 主库：处理写入操作
   - 从库：处理读取操作

### 附录 B：数据库监控指标

| 指标 | 告警阈值 | 监控方式 |
|------|---------|---------|
| 数据库连接数 | > 80% | Prometheus |
| 慢查询数量 | > 10条/分钟 | 日志分析 |
| 锁等待时间 | > 1秒 | 性能监控 |
| 磁盘使用率 | > 80% | 监控系统 |
| 表空间增长 | > 10MB/天 | 自动化脚本 |

---

**文档结束**

*此数据库设计文档将作为工具库模块开发的基准文件。*
