package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// User 用户模型
type User struct {
	ID            string                 `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	Email         string                 `gorm:"column:email;type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash  string                 `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Nickname      string                 `gorm:"column:nickname;type:varchar(50);not null" json:"nickname"`
	Avatar        string                 `gorm:"column:avatar;type:varchar(500)" json:"avatar,omitempty"`
	Bio           string                 `gorm:"column:bio;type:text" json:"bio,omitempty"`
	Status        UserStatus             `gorm:"column:status;type:enum('active','inactive','banned');default:'active'" json:"status"`
	EmailVerified bool                   `gorm:"column:email_verified;default:false" json:"email_verified"`
	LastLoginAt   *time.Time             `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	LastLoginIP   string                 `gorm:"column:last_login_ip;type:varchar(45)" json:"last_login_ip,omitempty"`
	Preferences   UserPreferences        `gorm:"column:preferences;type:json" json:"preferences"`
	CreatedAt     time.Time              `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time              `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

// UserPreferences 用户偏好设置
type UserPreferences struct {
	// 内容偏好
	InterestedCategories []string `json:"interested_categories,omitempty"`
	FollowedProfessions  []string `json:"followed_professions,omitempty"`

	// 通知设置
	EmailNotifications EmailNotificationSettings `json:"email_notifications"`

	// 显示设置
	Theme    string `json:"theme,omitempty"`
	Language string `json:"language,omitempty"`
}

// EmailNotificationSettings 邮件通知设置
type EmailNotificationSettings struct {
	SystemUpdate bool `json:"system_update"`
	NewContent   bool `json:"new_content"`
	WeeklyDigest bool `json:"weekly_digest"`
}

// Value 实现 driver.Valuer 接口
func (p UserPreferences) Value() (driver.Value, error) {
	data, err := json.Marshal(p)
	return string(data), err
}

// Scan 实现 sql.Scanner 接口
func (p *UserPreferences) Scan(value interface{}) error {
	if value == nil {
		*p = UserPreferences{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, p)
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Role 角色模型
type Role struct {
	ID          string    `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"column:description;type:varchar(255)" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限模型
type Permission struct {
	ID          string    `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(100);uniqueIndex;not null" json:"name"`
	Resource    string    `gorm:"column:resource;type:varchar(50);not null;index:idx_resource_action" json:"resource"`
	Action      string    `gorm:"column:action;type:varchar(20);not null;index:idx_resource_action" json:"action"`
	Description string    `gorm:"column:description;type:varchar(255)" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// UserRole 用户角色关联
type UserRole struct {
	UserID string `gorm:"column:user_id;type:varchar(36);primaryKey" json:"user_id"`
	RoleID string `gorm:"column:role_id;type:varchar(36);primaryKey" json:"role_id"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色权限关联
type RolePermission struct {
	RoleID       string `gorm:"column:role_id;type:varchar(36);primaryKey" json:"role_id"`
	PermissionID string `gorm:"column:permission_id;type:varchar(36);primaryKey" json:"permission_id"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// Favorite 收藏模型
type Favorite struct {
	ID           string     `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	UserID       string     `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_id" json:"user_id"`
	Type         string     `gorm:"column:type;type:enum('article','profession','learningPath','tool');not null;index:idx_user_type" json:"type"`
	TargetID     string     `gorm:"column:target_id;type:varchar(36);not null;index:idx_target" json:"target_id"`
	TargetTitle  string     `gorm:"column:target_title;type:varchar(500);not null" json:"target_title"`
	TargetSummary string    `gorm:"column:target_summary;type:text" json:"target_summary,omitempty"`
	Tags         JSONArray  `gorm:"column:tags;type:json" json:"tags"`
	Note         string     `gorm:"column:note;type:text" json:"note,omitempty"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}

// RefreshToken 刷新 Token 模型
type RefreshToken struct {
	ID        string     `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	UserID    string     `gorm:"column:user_id;type:varchar(36);not null;index:idx_user_id" json:"user_id"`
	TokenHash string     `gorm:"column:token_hash;type:varchar(255);uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null;index:idx_expires_at" json:"expires_at"`
	Revoked   bool       `gorm:"column:revoked;default:false" json:"revoked"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// JSONArray 自定义 JSON 数组类型
type JSONArray []string

// Value 实现 driver.Valuer 接口
func (j JSONArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	data, err := json.Marshal(j)
	return string(data), err
}

// Scan 实现 sql.Scanner 接口
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}
