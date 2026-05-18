package repository

import (
	"context"
	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/models"

	"github.com/zeromicro/go-zero/core/logx"
)

// ToolReviewRepository 工具审核记录仓储
type ToolReviewRepository struct{}

// NewToolReviewRepository 创建仓储实例
func NewToolReviewRepository() *ToolReviewRepository {
	return &ToolReviewRepository{}
}

// Create 创建审核记录
func (r *ToolReviewRepository) Create(ctx context.Context, record *models.ToolReviewRecord) error {
	return database.GetDB().WithContext(ctx).Create(record).Error
}

// GetByToolID 获取工具的所有审核记录
func (r *ToolReviewRepository) GetByToolID(ctx context.Context, toolID uint) ([]models.ToolReviewRecord, error) {
	var records []models.ToolReviewRecord
	err := database.GetDB().WithContext(ctx).Where("tool_id = ?", toolID).Order("created_at DESC").Find(&records).Error
	if err != nil {
		logx.Errorf("Query tool_reviews error: %v", err)
		return nil, err
	}
	return records, nil
}