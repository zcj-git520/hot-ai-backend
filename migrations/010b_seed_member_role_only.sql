-- 010b_seed_member_role_only.sql (prod-corrected, replaces 010)
-- prod roles table has no updated_at, only created_at.
-- admin/user already exist; only member is missing.

INSERT IGNORE INTO roles (id, name, description, created_at) VALUES
  (UUID(), 'member', '付费会员', NOW());
