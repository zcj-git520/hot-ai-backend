-- ======================================================
-- AI 热点追踪平台 - 学习路径模块表结构
-- 版本：v1.0
-- 日期：2026-04-08
-- ======================================================

-- 删除已存在的表（按依赖顺序倒序）
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `learning_progress`;
DROP TABLE IF EXISTS `path_chapters`;
DROP TABLE IF EXISTS `learning_paths`;
DROP TABLE IF EXISTS learning_path_management;
SET FOREIGN_KEY_CHECKS = 1;

USE `hot_ai`;

-- ============================================
-- 1. 学习路径表（核心表）
-- ============================================
CREATE TABLE `learning_paths` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '路径 ID',
    `title` varchar(100) NOT NULL COMMENT '路径标题',
    `slug` varchar(100) NOT NULL COMMENT 'URL 友好标识',
    `icon` varchar(20) DEFAULT NULL COMMENT 'Emoji 图标',
    `description` text COMMENT '路径描述',
    `difficulty` enum('beginner','intermediate','advanced') NOT NULL DEFAULT 'beginner' COMMENT '难度等级：beginner-入门，intermediate-进阶，advanced-高级',
    `level_label` varchar(10) NOT NULL DEFAULT '入门' COMMENT '难度标签（中文）',
    `learning_goals` json DEFAULT NULL COMMENT '学习目标列表',
    `target_audience` json DEFAULT NULL COMMENT '适合人群列表',
    `estimated_days` int DEFAULT '30' COMMENT '预计学习天数',
    `estimated_hours` int DEFAULT '60' COMMENT '预计学习小时数',
    `chapter_count` int DEFAULT '0' COMMENT '章节数量',
    `student_count` int DEFAULT '0' COMMENT '学习人数',
    `cover_image` varchar(255) DEFAULT NULL COMMENT '封面图片 URL',
    `is_featured` tinyint DEFAULT '0' COMMENT '是否推荐：0-否，1-是',
    `is_active` tinyint DEFAULT '1' COMMENT '是否激活：0-否，1-是',
    `sort_order` int DEFAULT '0' COMMENT '排序权重',
    `status` tinyint DEFAULT '1' COMMENT '状态：0-待审核，1-已发布，2-已删除',
    `published_at` datetime DEFAULT NULL COMMENT '发布时间',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学习路径表';

-- ============================================
-- 2. 学习路径章节表
-- ============================================
CREATE TABLE `path_chapters` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '章节 ID',
    `path_id` int NOT NULL COMMENT '路径 ID',
    `title` varchar(100) NOT NULL COMMENT '章节标题',
    `slug` varchar(100) NOT NULL COMMENT 'URL 友好标识',
    `description` text COMMENT '章节描述',
    `content_type` enum('article','video','practice','external') NOT NULL DEFAULT 'article' COMMENT '内容类型：article-图文，video-视频，practice-实践，external-外部资源',
    `content` text COMMENT '内容正文（Markdown 格式）',
    `video_url` varchar(500) DEFAULT NULL COMMENT '视频 URL',
    `external_links` json DEFAULT NULL COMMENT '外部资源链接列表',
    `estimated_hours` int DEFAULT '1' COMMENT '预计学习小时数',
    `order_index` int NOT NULL DEFAULT '0' COMMENT '排序索引',
    `is_free` tinyint DEFAULT '1' COMMENT '是否免费：0-否，1-是',
    `status` tinyint DEFAULT '1' COMMENT '状态：0-草稿，1-已发布，2-已删除',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学习路径章节表';

-- ============================================
-- 3. 学习进度表
-- ============================================
CREATE TABLE `learning_progress` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '进度 ID',
    `user_id` varchar(36) DEFAULT NULL COMMENT '用户 ID（未登录时为 NULL）',
    `session_id` varchar(100) DEFAULT NULL COMMENT '会话 ID（未登录时使用）',
    `path_id` int NOT NULL COMMENT '路径 ID',
    `chapter_id` int NOT NULL COMMENT '章节 ID',
    `status` enum('in_progress','completed') NOT NULL DEFAULT 'in_progress' COMMENT '学习状态：in_progress-进行中，completed-已完成',
    `completed_at` datetime DEFAULT NULL COMMENT '完成时间',
    `time_spent` int DEFAULT '0' COMMENT '学习时长（分钟）',
    `notes` text COMMENT '学习笔记',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学习进度表';

-- ============================================
-- 4. 学习路径管理数据表
-- ============================================
CREATE TABLE `learning_path_management` (
    `id` int NOT NULL AUTO_INCREMENT COMMENT '管理数据 ID',
    `path_id` int NOT NULL COMMENT '路径 ID',
    `view_count` int DEFAULT '0' COMMENT '查看次数',
    `start_count` int DEFAULT '0' COMMENT '开始学习次数',
    `complete_count` int DEFAULT '0' COMMENT '完成次数',
    `favorite_count` int DEFAULT '0' COMMENT '收藏次数',
    `meta_title` varchar(255) DEFAULT NULL COMMENT 'SEO 标题',
    `meta_description` varchar(500) DEFAULT NULL COMMENT 'SEO 描述',
    `meta_keywords` varchar(500) DEFAULT NULL COMMENT 'SEO 关键词',
    `reviewer_id` varchar(36) DEFAULT NULL COMMENT '审核人 ID',
    `reviewed_at` datetime DEFAULT NULL COMMENT '审核时间',
    `review_notes` text COMMENT '审核备注',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学习路径管理数据表';
