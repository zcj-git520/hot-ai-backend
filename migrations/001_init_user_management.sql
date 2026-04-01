-- AI 热点追踪平台 - 用户管理模块数据库初始化脚本
-- 版本：v1.0
-- 创建日期：2026-03-30

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS hot_ai DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE hot_ai;

-- ============================================
-- 用户表
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY COMMENT '用户 ID（UUID）',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱（唯一）',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    nickname VARCHAR(50) NOT NULL COMMENT '昵称',
    avatar VARCHAR(500) COMMENT '头像 URL',
    bio TEXT COMMENT '个人简介',
    status ENUM('active', 'inactive', 'banned') DEFAULT 'active' COMMENT '账号状态',
    email_verified BOOLEAN DEFAULT FALSE COMMENT '邮箱是否验证',
    last_login_at DATETIME COMMENT '最后登录时间',
    last_login_ip VARCHAR(45) COMMENT '最后登录 IP',
    preferences JSON COMMENT '偏好设置',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_email (email),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 角色表
-- ============================================
CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(36) PRIMARY KEY COMMENT '角色 ID',
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名称',
    description VARCHAR(255) COMMENT '角色描述',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ============================================
-- 权限表
-- ============================================
CREATE TABLE IF NOT EXISTS permissions (
    id VARCHAR(36) PRIMARY KEY COMMENT '权限 ID',
    name VARCHAR(100) NOT NULL UNIQUE COMMENT '权限名称',
    resource VARCHAR(50) NOT NULL COMMENT '资源类型',
    action VARCHAR(20) NOT NULL COMMENT '操作类型',
    description VARCHAR(255) COMMENT '权限描述',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_resource_action (resource, action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- ============================================
-- 用户角色关联表（多对多）
-- ============================================
CREATE TABLE IF NOT EXISTS user_roles (
    user_id VARCHAR(36) NOT NULL COMMENT '用户 ID',
    role_id VARCHAR(36) NOT NULL COMMENT '角色 ID',
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ============================================
-- 角色权限关联表（多对多）
-- ============================================
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id VARCHAR(36) NOT NULL COMMENT '角色 ID',
    permission_id VARCHAR(36) NOT NULL COMMENT '权限 ID',
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- ============================================
-- 收藏表
-- ============================================
CREATE TABLE IF NOT EXISTS favorites (
    id VARCHAR(36) PRIMARY KEY COMMENT '收藏 ID',
    user_id VARCHAR(36) NOT NULL COMMENT '用户 ID',
    type ENUM('article', 'profession', 'learningPath', 'tool') NOT NULL COMMENT '收藏类型',
    target_id VARCHAR(36) NOT NULL COMMENT '目标资源 ID',
    target_title VARCHAR(500) NOT NULL COMMENT '目标资源标题',
    target_summary TEXT COMMENT '目标资源摘要',
    tags JSON COMMENT '标签',
    note TEXT COMMENT '备注',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_user_type (user_id, type),
    INDEX idx_target (type, target_id),
    UNIQUE KEY uk_user_target (user_id, type, target_id) COMMENT '防止重复收藏'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏表';

-- ============================================
-- 刷新 Token 表（用于 Token 黑名单管理）
-- ============================================
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id VARCHAR(36) PRIMARY KEY COMMENT 'Token ID',
    user_id VARCHAR(36) NOT NULL COMMENT '用户 ID',
    token_hash VARCHAR(255) NOT NULL UNIQUE COMMENT 'Token 哈希值',
    expires_at DATETIME NOT NULL COMMENT '过期时间',
    revoked BOOLEAN DEFAULT FALSE COMMENT '是否已撤销',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_token_hash (token_hash),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='刷新 Token 表';

-- ============================================
-- 初始化基础数据
-- ============================================

-- 插入预定义角色
INSERT INTO roles (id, name, description) VALUES
('role_guest', 'guest', '访客 - 仅浏览公开内容'),
('role_user', 'user', '普通用户 - 浏览 + 收藏 + 评论'),
('role_editor', 'editor', '内容编辑 - 内容 CRUD + 审核'),
('role_admin', 'admin', '系统管理员 - 所有权限'),
('role_crawler', 'crawler', '爬虫服务 - 服务间调用权限')
ON DUPLICATE KEY UPDATE name=name;

-- 插入基础权限
INSERT INTO permissions (id, name, resource, action, description) VALUES
-- 文章权限
('perm_article_view', 'article:view', 'article', 'view', '查看文章'),
('perm_article_create', 'article:create', 'article', 'create', '创建文章'),
('perm_article_edit', 'article:edit', 'article', 'edit', '编辑文章'),
('perm_article_delete', 'article:delete', 'article', 'delete', '删除文章'),
('perm_article_publish', 'article:publish', 'article', 'publish', '发布文章'),

-- 职业权限
('perm_profession_view', 'profession:view', 'profession', 'view', '查看职业'),
('perm_profession_create', 'profession:create', 'profession', 'create', '创建职业'),
('perm_profession_edit', 'profession:edit', 'profession', 'edit', '编辑职业'),
('perm_profession_delete', 'profession:delete', 'profession', 'delete', '删除职业'),

-- 工具权限
('perm_tool_view', 'tool:view', 'tool', 'view', '查看工具'),
('perm_tool_create', 'tool:create', 'tool', 'create', '创建工具'),
('perm_tool_edit', 'tool:edit', 'tool', 'edit', '编辑工具'),
('perm_tool_delete', 'tool:delete', 'tool', 'delete', '删除工具'),

-- 学习路径权限
('perm_learning_path_view', 'learningPath:view', 'learningPath', 'view', '查看学习路径'),
('perm_learning_path_create', 'learningPath:create', 'learningPath', 'create', '创建学习路径'),
('perm_learning_path_edit', 'learningPath:edit', 'learningPath', 'edit', '编辑学习路径'),
('perm_learning_path_delete', 'learningPath:delete', 'learningPath', 'delete', '删除学习路径'),

-- 用户管理权限
('perm_user_manage', 'user:manage', 'user', 'manage', '用户管理'),
('perm_user_view', 'user:view', 'user', 'view', '查看用户'),
('perm_user_edit', 'user:edit', 'user', 'edit', '编辑用户'),
('perm_user_delete', 'user:delete', 'user', 'delete', '删除用户'),

-- 角色管理权限
('perm_role_manage', 'role:manage', 'role', 'manage', '角色管理'),
('perm_role_view', 'role:view', 'role', 'view', '查看角色'),
('perm_role_create', 'role:create', 'role', 'create', '创建角色'),
('perm_role_edit', 'role:edit', 'role', 'edit', '编辑角色'),
('perm_role_delete', 'role:delete', 'role', 'delete', '删除角色'),

-- 系统配置权限
('perm_system_config', 'system:config', 'system', 'config', '系统配置'),
('perm_system_monitor', 'system:monitor', 'system', 'monitor', '系统监控'),

-- 收藏权限
('perm_favorite_view', 'favorite:view', 'favorite', 'view', '查看收藏'),
('perm_favorite_create', 'favorite:create', 'favorite', 'create', '添加收藏'),
('perm_favorite_delete', 'favorite:delete', 'favorite', 'delete', '删除收藏'),

-- 评论权限
('perm_comment_view', 'comment:view', 'comment', 'view', '查看评论'),
('perm_comment_create', 'comment:create', 'comment', 'create', '发表评论'),
('perm_comment_delete', 'comment:delete', 'comment', 'delete', '删除评论')
ON DUPLICATE KEY UPDATE name=name;

-- 为 guest 角色分配权限（仅浏览公开内容）
INSERT INTO role_permissions (role_id, permission_id) VALUES
('role_guest', 'perm_article_view'),
('role_guest', 'perm_profession_view'),
('role_guest', 'perm_tool_view'),
('role_guest', 'perm_learning_path_view')
ON DUPLICATE KEY UPDATE role_id=role_id;

-- 为 user 角色分配权限（浏览 + 收藏 + 评论）
INSERT INTO role_permissions (role_id, permission_id) VALUES
('role_user', 'perm_article_view'),
('role_user', 'perm_profession_view'),
('role_user', 'perm_tool_view'),
('role_user', 'perm_learning_path_view'),
('role_user', 'perm_favorite_view'),
('role_user', 'perm_favorite_create'),
('role_user', 'perm_favorite_delete'),
('role_user', 'perm_comment_view'),
('role_user', 'perm_comment_create')
ON DUPLICATE KEY UPDATE role_id=role_id;

-- 为 editor 角色分配权限（内容 CRUD + 审核）
INSERT INTO role_permissions (role_id, permission_id) VALUES
('role_editor', 'perm_article_view'),
('role_editor', 'perm_article_create'),
('role_editor', 'perm_article_edit'),
('role_editor', 'perm_article_publish'),
('role_editor', 'perm_profession_view'),
('role_editor', 'perm_profession_create'),
('role_editor', 'perm_profession_edit'),
('role_editor', 'perm_tool_view'),
('role_editor', 'perm_tool_create'),
('role_editor', 'perm_tool_edit'),
('role_editor', 'perm_learning_path_view'),
('role_editor', 'perm_learning_path_create'),
('role_editor', 'perm_learning_path_edit'),
('role_editor', 'perm_favorite_view'),
('role_editor', 'perm_favorite_create'),
('role_editor', 'perm_favorite_delete'),
('role_editor', 'perm_comment_view'),
('role_editor', 'perm_comment_create'),
('role_editor', 'perm_comment_delete')
ON DUPLICATE KEY UPDATE role_id=role_id;

-- 为 admin 角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'role_admin', id FROM permissions
ON DUPLICATE KEY UPDATE role_id=role_id;

-- 为 crawler 角色分配服务调用权限
INSERT INTO role_permissions (role_id, permission_id) VALUES
('role_crawler', 'perm_article_create'),
('role_crawler', 'perm_profession_create'),
('role_crawler', 'perm_tool_create'),
('role_crawler', 'perm_learning_path_create')
ON DUPLICATE KEY UPDATE role_id=role_id;
