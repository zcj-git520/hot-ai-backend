-- ======================================================
-- AI热点追踪平台 - 数据库表结构
-- 版本: v1.0
-- 日期: 2026-04-06
-- ======================================================

CREATE DATABASE IF NOT EXISTS hot_ai CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE hot_ai;

-- 1. 分类表
CREATE TABLE `categories` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(20) NOT NULL COMMENT '分类名称',
  `code` varchar(20) NOT NULL COMMENT '分类标识',
  `color` varchar(10) NOT NULL COMMENT '颜色值',
  `icon` varchar(50) DEFAULT NULL COMMENT '图标名称',
  `sort_order` int DEFAULT '0' COMMENT '排序顺序',
  `status` tinyint DEFAULT '1' COMMENT '状态: 0-禁用, 1-启用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- 2. 文章来源表
CREATE TABLE `sources` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '来源ID',
  `name` varchar(100) NOT NULL COMMENT '来源名称',
  `domain` varchar(200) NOT NULL COMMENT '域名',
  `logo_url` varchar(500) DEFAULT NULL COMMENT 'Logo地址',
  `description` text DEFAULT NULL COMMENT '来源描述',
  `reliability_score` tinyint DEFAULT '5' COMMENT '可信度评分 1-10',
  `status` tinyint DEFAULT '1' COMMENT '状态: 0-禁用, 1-启用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='来源媒体表';

-- 3. 标签主表
CREATE TABLE `tags` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '标签ID',
  `name` varchar(50) NOT NULL COMMENT '标签名称',
  `type` bigint DEFAULT '0' COMMENT '类型',
  `status` tinyint DEFAULT '1' COMMENT '状态: 0-禁用, 1-启用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签主表';

-- 4. 资讯主表（正文存储在MongoDB）
CREATE TABLE `articles` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '资讯唯一标识',
  `title` varchar(200) NOT NULL COMMENT '资讯标题',
  `summary` text NOT NULL COMMENT 'AI生成摘要',
  `content_mongo_id` varchar(50) NOT NULL COMMENT 'MongoDB正文ObjectId',
  `source_id` bigint NOT NULL COMMENT '来源媒体ID',
  `author` varchar(50) DEFAULT NULL COMMENT '作者名称',
  `category_id` int NOT NULL COMMENT '分类ID',
  `published_at` datetime NOT NULL COMMENT '原文发布时间',
  `status` tinyint DEFAULT '1' COMMENT '状态: 0-待审核, 1-已发布, 2-已删除',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_content_mongo_id` (`content_mongo_id`),
  KEY `idx_source_id` (`source_id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_published_at` (`published_at`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_articles_source` FOREIGN KEY (`source_id`) REFERENCES `sources` (`id`),
  CONSTRAINT `fk_articles_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='资讯主表';

-- 5. 文章统计表
CREATE TABLE `article_stats` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `article_id` bigint NOT NULL COMMENT '文章ID',
  `view_count` bigint DEFAULT '0' COMMENT '阅读量',
  `comment_count` bigint DEFAULT '0' COMMENT '评论数量',
  `like_count` bigint DEFAULT '0' COMMENT '点赞数量',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_article_id` (`article_id`),
  KEY `idx_article_id` (`article_id`),
  CONSTRAINT `fk_stats_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章统计表';

-- 6. 文章-标签关联表
CREATE TABLE `article_tag_relation` (
  `article_id` bigint NOT NULL COMMENT '文章ID',
  `tag_id` bigint NOT NULL COMMENT '标签ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`article_id`, `tag_id`),
  KEY `idx_tag_id` (`tag_id`),
  CONSTRAINT `fk_article_tag_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_article_tag_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章标签关联表';
