package types

// CommonResponse 通用响应结构
type CommonResponse struct {
	// 状态码
	Code int `json:"code"`
	// 消息
	Message string `json:"message"`
	// 数据
	Data interface{} `json:"data,omitempty"`
}

// PageInfo 分页信息
type PageInfo struct {
	// 当前页码
	Page int64 `json:"page"`
	// 每页数量
	PageSize int64 `json:"page_size"`
	// 总数量
	Total int64 `json:"total"`
	// 总页数
	TotalPages int64 `json:"total_pages"`
}
