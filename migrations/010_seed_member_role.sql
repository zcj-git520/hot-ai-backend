-- 010_seed_member_role.sql
-- 确保 roles 表里有 'admin' / 'user' / 'member' 三个角色

INSERT IGNORE INTO roles (id, name, description, created_at, updated_at) VALUES
  (UUID(), 'admin',  '平台管理员', NOW(), NOW()),
  (UUID(), 'user',   '普通注册用户', NOW(), NOW()),
  (UUID(), 'member', '付费会员', NOW(), NOW());