-- ======================================================
-- AI热点追踪平台 - 增加文章翻译字段
-- 版本: v1.0
-- 日期: 2026-04-21
-- ======================================================

USE `hot_ai`;

-- 为 articles 表添加翻译相关字段
ALTER TABLE `articles`
ADD COLUMN `title_en` varchar(200) DEFAULT NULL COMMENT '英文标题' AFTER `summary`,
ADD COLUMN `summary_en` text DEFAULT NULL COMMENT '英文摘要' AFTER `title_en`,
ADD COLUMN `content_en` longtext DEFAULT NULL COMMENT '英文正文' AFTER `summary_en`;
