package models

import (
	"time"
)

// UserActivityLog 用户活动日志
type UserActivityLog struct {
	ID         uint      `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"id"`
	UserID     string    `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_id" json:"user_id"`
	Action     string    `gorm:"column:action;type:varchar(50);not null;index:idx_action" json:"action"`
	TargetType string    `gorm:"column:target_type;type:varchar(20)" json:"target_type,omitempty"`
	TargetID   string    `gorm:"column:target_id;type:varchar(36)" json:"target_id,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (UserActivityLog) TableName() string {
	return "user_activity_logs"
}