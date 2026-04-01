package types

// FavoriteItem 收藏项
type FavoriteItem struct {
	// 收藏 ID
	ID string `json:"id"`
	// 收藏类型 (article/profession/learningPath/tool)
	Type string `json:"type"`
	// 目标资源 ID
	TargetID string `json:"target_id"`
	// 目标资源标题
	TargetTitle string `json:"target_title"`
	// 目标资源摘要
	TargetSummary string `json:"target_summary,omitempty"`
	// 标签
	Tags []string `json:"tags"`
	// 备注
	Note string `json:"note,omitempty"`
	// 创建时间
	CreatedAt string `json:"created_at"`
}

// CreateFavoriteRequest 创建收藏请求
type CreateFavoriteRequest struct {
	// 收藏类型 (article/profession/learningPath/tool)
	Type string `json:"type"`
	// 目标资源 ID
	TargetID string `json:"target_id"`
	// 目标资源标题
	TargetTitle string `json:"target_title"`
	// 目标资源摘要
	TargetSummary string `json:"target_summary,omitempty"`
	// 标签
	Tags []string `json:"tags,omitempty"`
	// 备注
	Note string `json:"note,omitempty"`
}

// FavoriteListResponse 收藏列表响应
type FavoriteListResponse struct {
	// 收藏列表
	Items []FavoriteItem `json:"items"`
	// 分页信息
	PageInfo PageInfo `json:"page_info"`
	// 收藏统计
	Stats FavoriteStats `json:"stats"`
}

// FavoriteStats 收藏统计
type FavoriteStats struct {
	// 文章收藏数
	Articles int64 `json:"articles"`
	// 职业收藏数
	Professions int64 `json:"professions"`
	// 学习路径收藏数
	LearningPaths int64 `json:"learning_paths"`
	// 工具收藏数
	Tools int64 `json:"tools"`
}
