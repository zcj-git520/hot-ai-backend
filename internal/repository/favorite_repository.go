package repository

import (
	"context"
	"errors"

	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/models"

	"gorm.io/gorm"
)

// FavoriteRepository 收藏夹仓储接口
type FavoriteRepository interface {
	Create(ctx context.Context, favorite *models.Favorite) error
	GetByID(ctx context.Context, id string) (*models.Favorite, error)
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error

	// 查询
	ListByUserID(ctx context.Context, userID string, offset, limit int, favType *string) ([]*models.Favorite, int64, error)
	FindByTarget(ctx context.Context, userID, targetType, targetID string) (*models.Favorite, error)
	GetStats(ctx context.Context, userID string) (*FavoriteStats, error)
}

// FavoriteStats 收藏统计
type FavoriteStats struct {
	Articles      int64 `json:"articles"`
	Professions   int64 `json:"professions"`
	LearningPaths int64 `json:"learning_paths"`
	Tools         int64 `json:"tools"`
}

type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏夹仓储实例
func NewFavoriteRepository() FavoriteRepository {
	return &favoriteRepository{
		db: database.GetDB(),
	}
}

// Create 创建收藏
func (r *favoriteRepository) Create(ctx context.Context, favorite *models.Favorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

// GetByID 根据 ID 获取收藏
func (r *favoriteRepository) GetByID(ctx context.Context, id string) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.WithContext(ctx).First(&favorite, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &favorite, nil
}

// Delete 删除收藏
func (r *favoriteRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Favorite{}, "id = ?", id).Error
}

// BatchDelete 批量删除收藏
func (r *favoriteRepository) BatchDelete(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Delete(&models.Favorite{}, "id IN ?", ids).Error
}

// ListByUserID 获取用户收藏列表
func (r *favoriteRepository) ListByUserID(ctx context.Context, userID string, offset, limit int, favType *string) ([]*models.Favorite, int64, error) {
	var favorites []*models.Favorite
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Favorite{}).Where("user_id = ?", userID)

	if favType != nil && *favType != "" {
		query = query.Where("type = ?", *favType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return favorites, 0, nil
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&favorites).Error
	return favorites, total, err
}

// FindByTarget 根据目标资源查找收藏
func (r *favoriteRepository) FindByTarget(ctx context.Context, userID, targetType, targetID string) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND target_id = ?", userID, targetType, targetID).
		First(&favorite).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &favorite, nil
}

// GetStats 获取用户收藏统计
func (r *favoriteRepository) GetStats(ctx context.Context, userID string) (*FavoriteStats, error) {
	stats := &FavoriteStats{}

	type Result struct {
		Type  string `gorm:"column:type"`
		Count int64  `gorm:"column:count"`
	}

	var results []Result
	err := r.db.WithContext(ctx).
		Model(&models.Favorite{}).
		Select("type, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		switch result.Type {
		case "article":
			stats.Articles = result.Count
		case "profession":
			stats.Professions = result.Count
		case "learningPath":
			stats.LearningPaths = result.Count
		case "tool":
			stats.Tools = result.Count
		}
	}

	return stats, nil
}
