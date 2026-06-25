-- 016_default_member_only.sql
-- 把 access_level 默认值改成 2 (会员专享)，并把所有已有内容也设为 2。
-- admin 后台通过显式下调 (0 / 1) 来开放给游客或登录用户。
-- 同时给 users 加 is_member + member_expire_at 字段。

-- 1) 改默认值（仅影响新插入的行）
ALTER TABLE articles        MODIFY COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '0=游客, 1=普通用户, 2=会员';
ALTER TABLE professions    MODIFY COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '0=游客, 1=普通用户, 2=会员';
ALTER TABLE tools          MODIFY COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '0=游客, 1=普通用户, 2=会员';
ALTER TABLE learning_paths MODIFY COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '0=游客, 1=普通用户, 2=会员';
ALTER TABLE path_chapters  MODIFY COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '0=游客, 1=普通用户, 2=会员';

-- 2) 把现有所有内容都设为 2 (会员专享)，admin 再按需下调
UPDATE articles        SET access_level = 2 WHERE access_level = 0;
UPDATE professions    SET access_level = 2 WHERE access_level = 0;
UPDATE tools          SET access_level = 2 WHERE access_level = 0;
UPDATE learning_paths SET access_level = 2 WHERE access_level = 0;
UPDATE path_chapters  SET access_level = 2 WHERE access_level = 0;

-- 3) users 表加会员字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_member         TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '是否为付费会员';
ALTER TABLE users ADD COLUMN IF NOT EXISTS member_expire_at  DATETIME     NULL                COMMENT '会员到期时间（NULL = 永不过期）';
ALTER TABLE users ADD INDEX IF NOT EXISTS idx_is_member (is_member);
