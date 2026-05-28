-- 用户管理日志表
CREATE TABLE IF NOT EXISTS admin_operation_logs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  admin_user_id VARCHAR(36) NOT NULL COMMENT '管理员ID',
  target_user_id VARCHAR(36) NOT NULL COMMENT '目标用户ID',
  action VARCHAR(50) NOT NULL COMMENT '操作类型: change_role, disable, enable',
  detail TEXT COMMENT '操作详情(JSON)',
  ip VARCHAR(45) COMMENT 'IP地址',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  INDEX idx_admin_user (admin_user_id),
  INDEX idx_target_user (target_user_id),
  INDEX idx_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理后台操作日志';

-- 用户活动日志表
CREATE TABLE IF NOT EXISTS user_activity_logs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
  action VARCHAR(50) NOT NULL COMMENT '活动类型: login, logout, view_article, favorite',
  target_type VARCHAR(20) COMMENT '目标类型: article, profession, tool, learning_path',
  target_id VARCHAR(36) COMMENT '目标ID',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  INDEX idx_user_id (user_id),
  INDEX idx_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户活动日志';