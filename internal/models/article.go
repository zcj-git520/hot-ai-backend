package models

import (
	"time"
)

// Category 分类模型
type Category struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(20);not null" json:"name"`
	Code      string    `gorm:"column:code;type:varchar(20);not null;uniqueIndex:uk_code" json:"code"`
	Color     string    `gorm:"column:color;type:varchar(10);not null" json:"color"`
	Icon      string    `gorm:"column:icon;type:varchar(50)" json:"icon,omitempty"`
	SortOrder int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	Status    int       `gorm:"column:status;default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}

// Source 来源媒体模型
type Source struct {
	ID               uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name             string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Domain           string    `gorm:"column:domain;type:varchar(200);not null" json:"domain"`
	LogoURL          string    `gorm:"column:logo_url;type:varchar(500)" json:"logo_url,omitempty"`
	Description      string    `gorm:"column:description;type:text" json:"description,omitempty"`
	ReliabilityScore int       `gorm:"column:reliability_score;default:5" json:"reliability_score"`
	Status           int       `gorm:"column:status;default:1" json:"status"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Source) TableName() string {
	return "sources"
}

// Tag 标签模型
type Tag struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(50);not null;uniqueIndex:uk_name" json:"name"`
	Type      int       `gorm:"column:type;default:0" json:"type"`
	Status    int       `gorm:"column:status;default:1" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Tag) TableName() string {
	return "tags"
}

// Article 文章模型
type Article struct {
	ID             uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title          string    `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Summary        string    `gorm:"column:summary;type:text;not null" json:"summary"`
	Content        string    `gorm:"column:content;type:longtext" json:"content,omitempty"`
	TitleEn        string    `gorm:"column:title_en;type:varchar(200)" json:"title_en,omitempty"`
	SummaryEn      string    `gorm:"column:summary_en;type:text" json:"summary_en,omitempty"`
	ContentEn      string    `gorm:"column:content_en;type:longtext" json:"content_en,omitempty"`
	ContentMongoID string    `gorm:"column:content_mongo_id;type:varchar(50);uniqueIndex:uk_content_mongo_id" json:"content_mongo_id,omitempty"`
	OriginalURL    string    `gorm:"column:original_url;type:varchar(500)" json:"original_url,omitempty"`
	SourceID       uint      `gorm:"column:source_id;not null;index:idx_source_id" json:"source_id"`
	Author         string    `gorm:"column:author;type:varchar(255)" json:"author,omitempty"`
	CategoryID     uint      `gorm:"column:category_id;not null;index:idx_category_id" json:"category_id"`
	PublishedAt    time.Time `gorm:"column:published_at;not null;index:idx_published_at" json:"published_at"`
	Status         int       `gorm:"column:status;default:1;index:idx_status" json:"status"` // 0-待审核, 1-已发布, 2-已删除
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`

	// 关联数据(非数据库字段)
	SourceName   string   `gorm:"-" json:"source_name,omitempty"`
	CategoryName string   `gorm:"-" json:"category_name,omitempty"`
	Tags         []string `gorm:"-" json:"tags,omitempty"`
	ViewCount    int64    `gorm:"-" json:"view_count"`
	CommentCount int64    `gorm:"-" json:"comment_count"`
	LikeCount    int64    `gorm:"-" json:"like_count"`
}

func (Article) TableName() string {
	return "articles"
}

// ArticleStats 文章统计模型
type ArticleStats struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ArticleID    uint      `gorm:"column:article_id;not null;uniqueIndex:uk_article_id;index:idx_article_id" json:"article_id"`
	ViewCount    int64     `gorm:"column:view_count;default:0" json:"view_count"`
	CommentCount int64     `gorm:"column:comment_count;default:0" json:"comment_count"`
	LikeCount    int64     `gorm:"column:like_count;default:0" json:"like_count"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ArticleStats) TableName() string {
	return "article_stats"
}

// ArticleTagRelation 文章标签关联模型
type ArticleTagRelation struct {
	ArticleID uint      `gorm:"column:article_id;primaryKey" json:"article_id"`
	TagID     uint      `gorm:"column:tag_id;primaryKey;index:idx_tag_id" json:"tag_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ArticleTagRelation) TableName() string {
	return "article_tag_relation"
}
