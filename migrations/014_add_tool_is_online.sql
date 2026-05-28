-- ======================================================
-- AI 热点追踪平台 - 添加工具上线状态字段
-- 版本：v1.0
-- 日期：2026-05-26
-- ======================================================

USE `hot_ai`;

-- 在 tools 表添加 is_online 字段
ALTER TABLE `tools` ADD COLUMN `is_online` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否上线：0-未上线，1-已上线' AFTER `review_status`;