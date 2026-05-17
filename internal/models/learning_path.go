package models

import (
	"time"
)

// LearningPath 学习路径模型（核心表）
type LearningPath struct {
	ID             uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title          string    `gorm:"column:title;type:varchar(100);not null" json:"title"`
	Slug           string    `gorm:"column:slug;type:varchar(100);not null;uniqueIndex:uk_slug" json:"slug"`
	Icon           string    `gorm:"column:icon;type:varchar(20)" json:"icon,omitempty"`
	Description    string    `gorm:"column:description;type:text" json:"description,omitempty"`
	Difficulty     string    `gorm:"column:difficulty;type:enum('beginner','intermediate','advanced');default:'beginner';index:idx_difficulty" json:"difficulty"`
	LevelLabel     string    `gorm:"column:level_label;type:varchar(10);default:'入门'" json:"level_label"`
	LearningGoals  string    `gorm:"column:learning_goals;type:json" json:"learning_goals,omitempty"`
	TargetAudience string    `gorm:"column:target_audience;type:json" json:"target_audience,omitempty"`
	EstimatedDays  int       `gorm:"column:estimated_days;default:30" json:"estimated_days"`
	EstimatedHours int       `gorm:"column:estimated_hours;default:60" json:"estimated_hours"`
	ChapterCount   int       `gorm:"column:chapter_count;default:0" json:"chapter_count"`
	StudentCount   int       `gorm:"column:student_count;default:0" json:"student_count"`
	CoverImage     string    `gorm:"column:cover_image;type:varchar(255)" json:"cover_image,omitempty"`
	IsFeatured     int       `gorm:"column:is_featured;default:0" json:"is_featured"`
	IsActive       int       `gorm:"column:is_active;default:1" json:"is_active"`
	SortOrder      int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	Status         int       `gorm:"column:status;default:1;index:idx_status" json:"status"`
	PublishedAt    *time.Time `gorm:"column:published_at" json:"published_at,omitempty"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`

	// 关联数据（非数据库字段）
	Chapters       []PathChapter `gorm:"-" json:"chapters,omitempty"`
	ManagementData *LearningPathManagement `gorm:"-" json:"management_data,omitempty"`
}

func (LearningPath) TableName() string {
	return "learning_paths"
}

// PathChapter 学习路径章节模型
type PathChapter struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PathID        uint      `gorm:"column:path_id;not null;index:idx_path_id" json:"path_id"`
	Title         string    `gorm:"column:title;type:varchar(100);not null" json:"title"`
	Slug          string    `gorm:"column:slug;type:varchar(100);not null" json:"slug"`
	Description   string    `gorm:"column:description;type:text" json:"description,omitempty"`
	ContentType   string    `gorm:"column:content_type;type:enum('article','video','practice','external');default:'article'" json:"content_type"`
	Content       string    `gorm:"column:content;type:text" json:"content,omitempty"`
	VideoURL      string    `gorm:"column:video_url;type:varchar(500)" json:"video_url,omitempty"`
	ExternalLinks string    `gorm:"column:external_links;type:json" json:"external_links,omitempty"`
	EstimatedHours int      `gorm:"column:estimated_hours;default:1" json:"estimated_hours"`
	OrderIndex    int       `gorm:"column:order_index;not null;default:0;index:idx_order_index" json:"order_index"`
	IsFree        int       `gorm:"column:is_free;default:1" json:"is_free"`
	Status        int       `gorm:"column:status;default:1;index:idx_status" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (PathChapter) TableName() string {
	return "path_chapters"
}

// LearningProgress 学习进度模型
type LearningProgress struct {
	ID          uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      string    `gorm:"column:user_id;type:varchar(36);index:idx_user_id" json:"user_id,omitempty"`
	SessionID   string    `gorm:"column:session_id;type:varchar(100);index:idx_session_id" json:"session_id,omitempty"`
	PathID      uint      `gorm:"column:path_id;not null;index:idx_path_id" json:"path_id"`
	ChapterID   uint      `gorm:"column:chapter_id;not null;index:idx_chapter_id" json:"chapter_id"`
	Status      string    `gorm:"column:status;type:enum('in_progress','completed');default:'in_progress'" json:"status"`
	CompletedAt time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
	TimeSpent   int       `gorm:"column:time_spent;default:0" json:"time_spent"` // 分钟
	Notes       string    `gorm:"column:notes;type:text" json:"notes,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (LearningProgress) TableName() string {
	return "learning_progress"
}

// LearningPathManagement 学习路径管理数据模型
type LearningPathManagement struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PathID        uint      `gorm:"column:path_id;not null;uniqueIndex:uk_path_id;index:idx_path_id" json:"path_id"`
	ViewCount     int       `gorm:"column:view_count;default:0" json:"view_count"`
	StartCount    int       `gorm:"column:start_count;default:0" json:"start_count"`
	CompleteCount int       `gorm:"column:complete_count;default:0" json:"complete_count"`
	FavoriteCount int       `gorm:"column:favorite_count;default:0" json:"favorite_count"`
	MetaTitle     string    `gorm:"column:meta_title;type:varchar(255)" json:"meta_title,omitempty"`
	MetaDescription string  `gorm:"column:meta_description;type:varchar(500)" json:"meta_description,omitempty"`
	MetaKeywords  string    `gorm:"column:meta_keywords;type:varchar(500)" json:"meta_keywords,omitempty"`
	ReviewerID    string    `gorm:"column:reviewer_id;type:varchar(36)" json:"reviewer_id,omitempty"`
	ReviewedAt    time.Time `gorm:"column:reviewed_at" json:"reviewed_at,omitempty"`
	ReviewNotes   string    `gorm:"column:review_notes;type:text" json:"review_notes,omitempty"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (LearningPathManagement) TableName() string {
	return "learning_path_management"
}

// LevelInfo 难度等级展示信息（静态数据）
type LevelInfo struct {
	ID          string `json:"id"`
	Level       string `json:"level"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Color       string `json:"color"`
	MinHours    int    `json:"min_hours"`
	MaxHours    int    `json:"max_hours"`
}
