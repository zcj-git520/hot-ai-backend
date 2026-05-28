package models

import (
	"time"
)

// AdminOperationLog 管理后台操作日志
type AdminOperationLog struct {
	ID           uint      `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"id"`
	AdminUserID  string    `gorm:"column:admin_user_id;type:varchar(36);not null;index:idx_admin_user" json:"admin_user_id"`
	TargetUserID string    `gorm:"column:target_user_id;type:varchar(36);not null;index:idx_target_user" json:"target_user_id"`
	Action       string    `gorm:"column:action;type:varchar(50);not null;index:idx_action" json:"action"`
	Detail       string    `gorm:"column:detail;type:text" json:"detail"`
	IP           string    `gorm:"column:ip;type:varchar(45)" json:"ip"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (AdminOperationLog) TableName() string {
	return "admin_operation_logs"
}