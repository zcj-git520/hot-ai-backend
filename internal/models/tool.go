package models

import "time"

// ToolCategory 工具类别
type ToolCategory struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	Featured    bool      `json:"featured"`
	Status      int8      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// Tool 工具
type Tool struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	Icon            string  `json:"icon"`
	Description     string  `json:"description"`
	OfficialURL     string  `json:"official_url"`
	DocumentationURL string `json:"documentation_url"`
	Pricing         string  `json:"pricing"`
	PricingDesc     string  `json:"pricing_description"`
	CategoryID      uint    `json:"category_id"`
	Difficulty      string  `json:"difficulty"`
	Rating          float64 `json:"rating"`
	ReviewCount     int     `json:"review_count"`
	ViewCount       int     `json:"view_count"`
	Popularity      int     `json:"popularity"`
	Tags            []string `json:"tags"`
	Featured        bool    `json:"featured"`
	Status          int8    `json:"status"`
	IsOnline        bool    `json:"is_online"`
	ExternalID      *string   `json:"external_id"`
	CreatedBy       *string   `json:"created_by"`
	UpdatedBy       *string   `json:"updated_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`

	// 审核相关字段
	ReviewStatus    string    `json:"review_status"`     // pending/rejected/approved/revision_requested
	SubmittedAt     time.Time `json:"submitted_at"`     // 提交时间
	SubmittedBy     string    `json:"submitted_by"`     // 申请人
	RevisionReason  string    `json:"revision_reason"`  // 退回原因
}

// ToolReview 用户评测
type ToolReview struct {
	ID                uint      `json:"id"`
	UserID            string    `json:"user_id"`
	ToolID            uint      `json:"tool_id"`
	UserIP            string    `json:"user_ip"`
	UserAgent         string    `json:"user_agent"`
	Rating            int8      `json:"rating"`
	EaseOfUse         int8      `json:"ease_of_use"`
	Effectiveness     int8      `json:"effectiveness"`
	ValueForMoney     int8      `json:"value_for_money"`
	Features          int8      `json:"features"`
	UpdateFrequency   int8      `json:"update_frequency"`
	Support           int8      `json:"support"`
	Pros              string    `json:"pros"`
	Cons              string    `json:"cons"`
	Comment           string    `json:"comment"`
	Images            string    `json:"images"`
	ProsJSON          string    `json:"pros_json"`
	ConsJSON          string    `json:"cons_json"`
	IsAnonymous       bool      `json:"is_anonymous"`
	Status            int8      `json:"status"`
	Reason            string    `json:"reason"`
	IsVerifiedPurchase *bool    `json:"is_verified_purchase"`
	VerifiedAt        *time.Time `json:"verified_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

// Comment 评论
type Comment struct {
	ID              uint      `json:"id"`
	UserID          string    `json:"user_id"`
	CommentableType string    `json:"commentable_type"`
	CommentableID   uint      `json:"commentable_id"`
	ParentID        *uint     `json:"parent_id"`
	Content         string    `json:"content"`
	Images          string    `json:"images"`
	IsAnonymous     bool      `json:"is_anonymous"`
	Likes           int       `json:"likes"`
	IsLiked         bool      `json:"is_liked"`
	IsSpam          bool      `json:"is_spam"`
	SpamReason      string    `json:"spam_reason"`
	SpamCount       int       `json:"spam_count"`
	Status          int8      `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

// PromptTemplate 提示词模板
type PromptTemplate struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	Slug              string    `json:"slug"`
	Description       string    `json:"description"`
	Content           string    `json:"content"`
	ToolID            *uint     `json:"tool_id"`
	CategoryID        *uint     `json:"category_id"`
	UseCases          string    `json:"use_cases"`
	Tags              string    `json:"tags"`
	ExampleResponse   string    `json:"example_response"`
	ExampleInput      string    `json:"example_input"`
	Likes             int       `json:"likes"`
	Views             int       `json:"views"`
	Favorites         int       `json:"favorites"`
	Featured          bool      `json:"featured"`
	Status            int8      `json:"status"`
	Language          string    `json:"language"`
	AuthorType        string    `json:"author_type"`
	AuthorID          string    `json:"author_id"`
	ApprovedAt        *time.Time `json:"approved_at"`
	ApprovedBy        string    `json:"approved_by"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

// UserFavorite 用户收藏
type UserFavorite struct {
	ID            uint      `json:"id"`
	UserID        string    `json:"user_id"`
	ToolID        uint      `json:"tool_id"`
	Note          string    `json:"note"`
	FavoriteListID *uint    `json:"favorite_list_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// Badge 徽章
type Badge struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Slug           string  `json:"slug"`
	Description    string  `json:"description"`
	Icon           string  `json:"icon"`
	Type           string  `json:"type"`
	ConditionType  string  `json:"condition_type"`
	ConditionValue int     `json:"condition_value"`
	IconColor      string  `json:"icon_color"`
	BackgroundColor string `json:"background_color"`
	Status         int8    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// UserBadge 用户徽章
type UserBadge struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"user_id"`
	BadgeID   uint      `json:"badge_id"`
	IssuedAt  time.Time `json:"issued_at"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	ID           uint    `json:"id"`
	Key          string  `json:"key"`
	Value        string  `json:"value"`
	ValueType    string  `json:"value_type"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToolSearchHistory 工具搜索历史
type ToolSearchHistory struct {
	ID           uint      `json:"id"`
	UserID       string    `json:"user_id"`
	SearchKeyword string  `json:"search_keyword"`
	CategoryID   *uint     `json:"category_id"`
	IsFree       *bool     `json:"is_free"`
	MinRating    *float64  `json:"min_rating"`
	CreatedAt    time.Time `json:"created_at"`
}
