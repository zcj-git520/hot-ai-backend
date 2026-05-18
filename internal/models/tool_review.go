package models

import "time"

// ToolReviewRecord 工具审核记录
type ToolReviewRecord struct {
	ID        uint      `json:"id"`
	ToolID    uint      `json:"tool_id"`
	AdminID   string    `json:"admin_id"`
	Action    string    `json:"action"`    // approve/reject/request_revision
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}