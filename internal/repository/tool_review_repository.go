package repository

import (
	"hot-ai-backend/internal/models"
	"hot-ai-backend/internal/database"
)

// ToolReviewRepository 工具审核记录仓储
type ToolReviewRepository struct{}

// NewToolReviewRepository 创建仓储实例
func NewToolReviewRepository() *ToolReviewRepository {
	return &ToolReviewRepository{}
}

// Create 创建审核记录
func (r *ToolReviewRepository) Create(record *models.ToolReviewRecord) error {
	return database.GetDB().Create(record).Error
}

// GetByToolID 获取工具的所有审核记录
func (r *ToolReviewRepository) GetByToolID(toolID uint) ([]models.ToolReviewRecord, error) {
	var records []models.ToolReviewRecord
	err := database.GetDB().Where("tool_id = ?", toolID).Order("created_at DESC").Find(&records).Error
	return records, err
}