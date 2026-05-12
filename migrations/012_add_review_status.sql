-- ======================================================
-- AI 热点追踪平台 - 添加审核状态字段
-- 版本：v1.0
-- 日期：2026-05-12
-- ======================================================

USE `hot_ai`;

-- 在 learning_paths 表添加 review_status 字段
ALTER TABLE `learning_paths` ADD COLUMN `review_status` tinyint NOT NULL DEFAULT 0 COMMENT '0-无需审核, 1-审核中, 2-审核拒绝' AFTER `status`;
