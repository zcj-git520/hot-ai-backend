-- 快速修复 articles 表字符集问题
-- 执行此脚本前请备份数据库

-- 检查当前字符集
SELECT 
    TABLE_NAME,
    COLUMN_NAME,
    CHARACTER_SET_NAME,
    COLLATION_NAME
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'articles'
  AND COLUMN_NAME IN ('title', 'summary', 'content', 'title_en', 'summary_en', 'content_en');

-- 修复 articles 表及所有文本字段
ALTER TABLE `articles` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 验证修复结果
SELECT 
    TABLE_NAME,
    COLUMN_NAME,
    CHARACTER_SET_NAME,
    COLLATION_NAME
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'articles'
  AND COLUMN_NAME IN ('title', 'summary', 'content', 'title_en', 'summary_en', 'content_en');
