package service

import (
	"context"
	"errors"
	"fmt"

	"hot-ai-backend/internal/models"
	"hot-ai-backend/internal/repository"
	"hot-ai-backend/internal/utils/passwordutil"

	"github.com/google/uuid"
)

// UserService 用户服务
type UserService struct {
	userRepo     repository.UserRepository
	favoriteRepo repository.FavoriteRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository, favoriteRepo repository.FavoriteRepository) *UserService {
	return &UserService{
		userRepo:     userRepo,
		favoriteRepo: favoriteRepo,
	}
}

// UserProfileRequest 用户资料请求
type UserProfileRequest struct {
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Bio      string `json:"bio,omitempty"`
}

// UserPreferencesRequest 偏好设置请求
type UserPreferencesRequest struct {
	InterestedCategories []string                      `json:"interested_categories,omitempty"`
	FollowedProfessions  []string                      `json:"followed_professions,omitempty"`
	EmailNotifications   *models.EmailNotificationSettings `json:"email_notifications,omitempty"`
	Theme                string                        `json:"theme,omitempty"`
	Language             string                        `json:"language,omitempty"`
}

// UserDetailResponse 用户详情响应
type UserDetailResponse struct {
	ID            string                  `json:"id"`
	Email         string                  `json:"email"`
	Nickname      string                  `json:"nickname"`
	Avatar        string                  `json:"avatar,omitempty"`
	Bio           string                  `json:"bio,omitempty"`
	Status        models.UserStatus       `json:"status"`
	EmailVerified bool                    `json:"email_verified"`
	Preferences   models.UserPreferences  `json:"preferences"`
	CreatedAt     string                  `json:"created_at"`
	UpdatedAt     string                  `json:"updated_at"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// GetUserProfile 获取用户资料
func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*UserDetailResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &UserDetailResponse{
		ID:            user.ID,
		Email:         user.Email,
		Nickname:      user.Nickname,
		Avatar:        user.Avatar,
		Bio:           user.Bio,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		Preferences:   user.Preferences,
		CreatedAt:     user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// UpdateUserProfile 更新用户资料
func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, req *UserProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	return s.userRepo.Update(ctx, user)
}

// UpdateUserPreferences 更新用户偏好
func (s *UserService) UpdateUserPreferences(ctx context.Context, userID string, req *UserPreferencesRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 更新偏好设置
	if req.InterestedCategories != nil {
		user.Preferences.InterestedCategories = req.InterestedCategories
	}
	if req.FollowedProfessions != nil {
		user.Preferences.FollowedProfessions = req.FollowedProfessions
	}
	if req.EmailNotifications != nil {
		user.Preferences.EmailNotifications = *req.EmailNotifications
	}
	if req.Theme != "" {
		user.Preferences.Theme = req.Theme
	}
	if req.Language != "" {
		user.Preferences.Language = req.Language
	}

	return s.userRepo.Update(ctx, user)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 验证原密码
	if !passwordutil.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return errors.New("incorrect password")
	}

	// 密码强度校验
	if len(req.NewPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// 加密新密码
	newPasswordHash, err := passwordutil.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = newPasswordHash
	return s.userRepo.Update(ctx, user)
}

// FavoriteService 收藏服务
type FavoriteService struct {
	favoriteRepo repository.FavoriteRepository
}

// NewFavoriteService 创建收藏服务实例
func NewFavoriteService(favoriteRepo repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{
		favoriteRepo: favoriteRepo,
	}
}

// CreateFavoriteRequest 创建收藏请求
type CreateFavoriteRequest struct {
	Type          string `json:"type"`
	TargetID      string `json:"target_id"`
	TargetTitle   string `json:"target_title"`
	TargetSummary string `json:"target_summary,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Note          string `json:"note,omitempty"`
}

// FavoriteItemResponse 收藏项响应
type FavoriteItemResponse struct {
	ID            string   `json:"id"`
	Type          string   `json:"type"`
	TargetID      string   `json:"target_id"`
	TargetTitle   string   `json:"target_title"`
	TargetSummary string   `json:"target_summary,omitempty"`
	Tags          []string `json:"tags"`
	Note          string   `json:"note,omitempty"`
	CreatedAt     string   `json:"created_at"`
}

// FavoriteListResponse 收藏列表响应
type FavoriteListResponse struct {
	Items      []*FavoriteItemResponse `json:"items"`
	Page       int64                   `json:"page"`
	PageSize   int64                   `json:"page_size"`
	Total      int64                   `json:"total"`
	TotalPages int64                   `json:"total_pages"`
	Stats      *repository.FavoriteStats `json:"stats"`
}

// CreateFavorite 创建收藏
func (s *FavoriteService) CreateFavorite(ctx context.Context, userID string, req *CreateFavoriteRequest) (*models.Favorite, error) {
	// 检查是否已收藏
	existing, err := s.favoriteRepo.FindByTarget(ctx, userID, req.Type, req.TargetID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing favorite: %w", err)
	}
	if existing != nil {
		return nil, errors.New("already favorited")
	}

	favorite := &models.Favorite{
		ID:            uuid.New().String(),
		UserID:        userID,
		Type:          req.Type,
		TargetID:      req.TargetID,
		TargetTitle:   req.TargetTitle,
		TargetSummary: req.TargetSummary,
		Tags:          req.Tags,
		Note:          req.Note,
	}

	err = s.favoriteRepo.Create(ctx, favorite)
	if err != nil {
		return nil, fmt.Errorf("failed to create favorite: %w", err)
	}

	return favorite, nil
}

// DeleteFavorite 删除收藏
func (s *FavoriteService) DeleteFavorite(ctx context.Context, userID, favoriteID string) error {
	favorite, err := s.favoriteRepo.GetByID(ctx, favoriteID)
	if err != nil {
		return fmt.Errorf("failed to get favorite: %w", err)
	}
	if favorite == nil {
		return errors.New("favorite not found")
	}

	// 验证所有权
	if favorite.UserID != userID {
		return errors.New("permission denied")
	}

	return s.favoriteRepo.Delete(ctx, favoriteID)
}

// BatchDeleteFavorites 批量删除收藏
func (s *FavoriteService) BatchDeleteFavorites(ctx context.Context, userID string, ids []string) (int64, error) {
	// TODO: 验证所有权
	err := s.favoriteRepo.BatchDelete(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to delete favorites: %w", err)
	}
	return int64(len(ids)), nil
}

// GetFavoriteList 获取收藏列表
func (s *FavoriteService) GetFavoriteList(ctx context.Context, userID string, page, pageSize int64, favType *string) (*FavoriteListResponse, error) {
	offset := (page - 1) * pageSize

	items, total, err := s.favoriteRepo.ListByUserID(ctx, userID, int(offset), int(pageSize), favType)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorites: %w", err)
	}

	// 获取统计
	stats, err := s.favoriteRepo.GetStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	// 转换响应
	itemResponses := make([]*FavoriteItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = &FavoriteItemResponse{
			ID:            item.ID,
			Type:          item.Type,
			TargetID:      item.TargetID,
			TargetTitle:   item.TargetTitle,
			TargetSummary: item.TargetSummary,
			Tags:          item.Tags,
			Note:          item.Note,
			CreatedAt:     item.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &FavoriteListResponse{
		Items:      itemResponses,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		Stats:      stats,
	}, nil
}

// GetUserCount 获取用户总数
func (s *UserService) GetUserCount() (int64, error) {
	return s.userRepo.GetUserCount()
}
