-- 009_add_access_level.sql
-- 给 5 张内容表加 access_level 列
-- 0=游客可读, 1=普通用户可读, 2=会员可读
-- 默认 0 = 当前所有内容都对游客开放（但游客实际只看到 500 字预览）

ALTER TABLE articles         ADD COLUMN IF NOT EXISTS access_level TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=游客, 1=普通用户, 2=会员';
ALTER TABLE professions     ADD COLUMN IF NOT EXISTS access_level TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=游客, 1=普通用户, 2=会员';
ALTER TABLE tools           ADD COLUMN IF NOT EXISTS access_level TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=游客, 1=普通用户, 2=会员';
ALTER TABLE learning_paths  ADD COLUMN IF NOT EXISTS access_level TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=游客, 1=普通用户, 2=会员';
ALTER TABLE path_chapters   ADD COLUMN IF NOT EXISTS access_level TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=游客, 1=普通用户, 2=会员';