package models

import (
	"time"
)

// Profession 职业模型
type Profession struct {
	ID             string    `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	Name           string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Slug           string    `gorm:"column:slug;type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description    string    `gorm:"column:description;type:text" json:"description,omitempty"`
	RiskLevel      string    `gorm:"column:risk_level;type:enum('extreme','high','medium','low');not null;index:idx_risk_level" json:"risk_level"`
	AutomationRate int       `gorm:"column:automation_rate;default:0" json:"automation_rate"`
	SafeSkills     JSONArray `gorm:"column:safe_skills;type:json" json:"safe_skills,omitempty"`
	AffectedSkills JSONArray `gorm:"column:affected_skills;type:json" json:"affected_skills,omitempty"`
	TransformTips  string    `gorm:"column:transform_tips;type:text" json:"transform_tips,omitempty"`
	LearningPaths  JSONArray `gorm:"column:learning_paths;type:json" json:"learning_paths,omitempty"`
	ViewCount      int       `gorm:"column:view_count;default:0" json:"view_count"`
	IsActive       bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Profession) TableName() string {
	return "professions"
}

// RiskLevelInfo 风险等级信息
type RiskLevelInfo struct {
	ID          string `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	Level       string `gorm:"column:level;type:enum('extreme','high','medium','low');uniqueIndex;not null" json:"level"`
	Name        string `gorm:"column:name;type:varchar(50);not null" json:"name"`
	Icon        string `gorm:"column:icon;type:varchar(20)" json:"icon"`
	Description string `gorm:"column:description;type:varchar(500)" json:"description"`
	Color       string `gorm:"column:color;type:varchar(20)" json:"color"`
	SortOrder   int    `gorm:"column:sort_order;default:0" json:"sort_order"`
}

// TableName 指定表名
func (RiskLevelInfo) TableName() string {
	return "risk_level_info"
}
