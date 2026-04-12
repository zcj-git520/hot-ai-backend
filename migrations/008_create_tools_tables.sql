-- ============================================================================
-- AI 热点追踪平台 - 工具库模块数据库设计
-- ============================================================================
-- 版本: v1.0
-- 创建日期: 2026-04-11
-- 说明: 工具库模块的完整数据库表结构设计
-- ============================================================================
-- 删除已存在的表（按依赖顺序倒序）

SET FOREIGN_KEY_CHECKS = 0;

-- 通知表（无外键依赖）
DROP TABLE IF EXISTS `notifications`;

-- 用户偏好表
DROP TABLE IF EXISTS `user_preferences`;

-- 用户徽章表
DROP TABLE IF EXISTS `user_badges`;

-- 徽章表
DROP TABLE IF EXISTS `badges`;

-- 热门工具统计表
DROP TABLE IF EXISTS `tool_hot_stats`;

-- 工具搜索历史表
DROP TABLE IF EXISTS `tool_search_history`;

-- 用户收藏表（依赖 tools）
DROP TABLE IF EXISTS `user_favorites`;

-- 评论表（依赖 tool_reviews 和 tools）
DROP TABLE IF EXISTS `comments`;

-- 用户评测表（依赖 tools）
DROP TABLE IF EXISTS `tool_reviews`;

-- 提示词模板表（依赖 tools 和 prompt_template_categories）
DROP TABLE IF EXISTS `prompt_templates`;

-- 提示词模板分类表
DROP TABLE IF EXISTS `prompt_template_categories`;

-- 工具-标签关联表（依赖 tools 和 tool_tags）
DROP TABLE IF EXISTS `tool_tag_relations`;

-- 工具表（依赖 tool_categories）
DROP TABLE IF EXISTS `tools`;

-- 工具标签表
DROP TABLE IF EXISTS `tool_tags`;

-- 工具类别表
DROP TABLE IF EXISTS `tool_categories`;

-- 系统配置表（无依赖，最后删除）
DROP TABLE IF EXISTS `system_config`;

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================================================
-- 1. 工具类别表 (tool_categories)
-- ============================================================================
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

-- ============================================================================
-- 2. 工具标签表 (tool_tags)
-- ============================================================================
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

-- ============================================================================
-- 3. 工具表 (tools)
-- ============================================================================
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
  KEY `idx_status_created_at` (`status`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='AI工具表';

-- ============================================================================
-- 4. 工具-标签关联表 (tool_tag_relations)
-- ============================================================================
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

-- ============================================================================
-- 5. 用户评测表 (tool_reviews)
-- ============================================================================
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

-- ============================================================================
-- 6. 评论表 (comments)
-- ============================================================================
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

-- ============================================================================
-- 7. 提示词模板表 (prompt_templates)
-- ============================================================================
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

-- ============================================================================
-- 8. 提示词模板分类表 (prompt_template_categories)
-- ============================================================================
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

-- ============================================================================
-- 9. 用户收藏表 (user_favorites)
-- ============================================================================
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

-- ============================================================================
-- 10. 徽章表 (badges)
-- ============================================================================
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

-- ============================================================================
-- 11. 用户徽章表 (user_badges)
-- ============================================================================
CREATE TABLE `user_badges` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户徽章ID，主键',
  `user_id` VARCHAR(100) NOT NULL COMMENT '用户ID',
  `badge_id` INT UNSIGNED NOT NULL COMMENT '徽章ID',
  `issued_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '获得时间',
  UNIQUE KEY `uk_user_badge` (`user_id`, `badge_id`),

  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_badge_id` (`badge_id`),
  KEY `idx_issued_at` (`issued_at`),
  KEY `idx_user_issued` (`user_id`, `issued_at`),
  CONSTRAINT `fk_user_badge_badge` FOREIGN KEY (`badge_id`) REFERENCES `badges` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户徽章表';

-- ============================================================================
-- 12. 用户偏好表 (user_preferences)
-- ============================================================================
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

-- ============================================================================
-- 13. 工具搜索历史表 (tool_search_history)
-- ============================================================================
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

-- ============================================================================
-- 14. 热门工具统计表 (tool_hot_stats)
-- ============================================================================
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

-- ============================================================================
-- 15. 通知表 (notifications)
-- ============================================================================
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

-- ============================================================================
-- 16. 系统配置表 (system_config)
-- ============================================================================
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
