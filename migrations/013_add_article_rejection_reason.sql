-- ======================================================
-- AI 热点追踪平台 - 添加文章拒绝原因字段
-- 版本：v1.0
-- 日期：2026-05-26
-- ======================================================

USE `hot_ai`;

-- 在 articles 表添加 rejection_reason 字段
ALTER TABLE `articles` ADD COLUMN `rejection_reason` VARCHAR(500) DEFAULT NULL COMMENT '拒绝原因' AFTER `status`;