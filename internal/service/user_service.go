package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

// GetUsersResponse 获取用户列表响应
type GetUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"total_pages"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`      // "admin" | "user"
	Status    string `json:"status"`    // "active" | "disabled"
	CreatedAt string `json:"created_at"`
	LastLogin string `json:"last_login,omitempty"`
}

// GetAdminLogsResponse 获取操作日志响应
type GetAdminLogsResponse struct {
	Logs      []AdminLogResponse `json:"logs"`
	Total     int64              `json:"total"`
	TotalPages int               `json:"total_pages"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
}

// AdminLogResponse 操作日志响应
type AdminLogResponse struct {
	ID        uint   `json:"id"`
	Action    string `json:"action"`
	Detail    string `json:"detail,omitempty"`
	AdminUser string `json:"admin_user,omitempty"`
	IP        string `json:"ip,omitempty"`
	CreatedAt string `json:"created_at"`
}

// GetUsers 获取用户列表（分页、筛选、搜索）
func (s *UserService) GetUsers(page, pageSize int, role, status, search string) (*GetUsersResponse, error) {
	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 12
	}

	users, total, err := s.userRepo.GetUsers(page, pageSize, role, status)
	if err != nil {
		return nil, err
	}

	// 获取用户角色信息
	userResponses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		roles, err := s.userRepo.GetUserRoles(context.Background(), user.ID)
		if err != nil {
			return nil, err
		}

		// 取第一个角色作为用户的角色
		roleName := "user"
		if len(roles) > 0 {
			roleName = roles[0].Name
		}

		// 状态映射：inactive/banned -> disabled
		statusStr := string(user.Status)
		if user.Status == models.UserStatusInactive || user.Status == models.UserStatusBanned {
			statusStr = "disabled"
		}

		// 搜索过滤（nickname/email）
		if search != "" {
			found := false
			lowercaseSearch := strings.ToLower(search)
			if strings.Contains(strings.ToLower(user.Nickname), lowercaseSearch) {
				found = true
			}
			if strings.Contains(strings.ToLower(user.Email), lowercaseSearch) {
				found = true
			}
			if !found {
				continue
			}
		}

		var lastLogin string
		if user.LastLoginAt != nil {
			lastLogin = user.LastLoginAt.Format("2006-01-02T15:04:05Z")
		}

		userResponses = append(userResponses, UserResponse{
			ID:        user.ID,
			Username:  user.Nickname,
			Email:     user.Email,
			Role:      roleName,
			Status:    statusStr,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			LastLogin: lastLogin,
		})
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &GetUsersResponse{
		Users:      userResponses,
		Total:      total,
		TotalPages: totalPages,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetUserByID 获取用户详情
func (s *UserService) GetUserByID(idStr string) (*UserResponse, error) {
	user, err := s.userRepo.GetUserByID(idStr)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 获取用户角色
	roles, err := s.userRepo.GetUserRoles(context.Background(), user.ID)
	if err != nil {
		return nil, err
	}

	roleName := "user"
	if len(roles) > 0 {
		roleName = roles[0].Name
	}

	// 状态映射：inactive/banned -> disabled
	statusStr := string(user.Status)
	if user.Status == models.UserStatusInactive || user.Status == models.UserStatusBanned {
		statusStr = "disabled"
	}

	var lastLogin string
	if user.LastLoginAt != nil {
		lastLogin = user.LastLoginAt.Format("2006-01-02T15:04:05Z")
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.Nickname,
		Email:     user.Email,
		Role:      roleName,
		Status:    statusStr,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		LastLogin: lastLogin,
	}, nil
}

// UpdateUserRole 更新用户角色
func (s *UserService) UpdateUserRole(idStr string, role string) error {
	// 验证角色值
	if role != "admin" && role != "user" && role != "" {
		return errors.New("invalid role: must be 'admin' or 'user'")
	}

	err := s.userRepo.UpdateUserRole(idStr, role)
	if err != nil {
		return err
	}

	// 记录操作日志
	adminLog := &models.AdminOperationLog{
		AdminUserID:  idStr, // 实际操作时需要从上下文获取当前管理员ID
		TargetUserID: idStr,
		Action:       "update_role",
		Detail:       "Changed role to " + role,
	}
	return s.userRepo.CreateAdminLog(adminLog)
}

// DisableUser 禁用用户 (active -> banned)
func (s *UserService) DisableUser(idStr string) error {
	err := s.userRepo.UpdateUserStatus(idStr, models.UserStatusBanned)
	if err != nil {
		return err
	}

	// 记录操作日志
	adminLog := &models.AdminOperationLog{
		AdminUserID:  idStr,
		TargetUserID: idStr,
		Action:       "disable_user",
		Detail:       "Disabled user",
	}
	return s.userRepo.CreateAdminLog(adminLog)
}

// EnableUser 启用用户 (banned -> active)
func (s *UserService) EnableUser(idStr string) error {
	err := s.userRepo.UpdateUserStatus(idStr, models.UserStatusActive)
	if err != nil {
		return err
	}

	// 记录操作日志
	adminLog := &models.AdminOperationLog{
		AdminUserID:  idStr,
		TargetUserID: idStr,
		Action:       "enable_user",
		Detail:       "Enabled user",
	}
	return s.userRepo.CreateAdminLog(adminLog)
}

// GetAdminLogs 获取操作日志
func (s *UserService) GetAdminLogs(userIDStr string, page, pageSize int) (*GetAdminLogsResponse, error) {
	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 12
	}

	logs, total, err := s.userRepo.GetAdminLogs(userIDStr, page, pageSize)
	if err != nil {
		return nil, err
	}

	logResponses := make([]AdminLogResponse, 0, len(logs))
	for _, log := range logs {
		// 获取管理员用户名
		adminUser := ""
		if log.AdminUserID != "" {
			adminUser = log.AdminUserID
		}

		logResponses = append(logResponses, AdminLogResponse{
			ID:        log.ID,
			Action:    log.Action,
			Detail:    log.Detail,
			AdminUser: adminUser,
			IP:        log.IP,
			CreatedAt: log.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &GetAdminLogsResponse{
		Logs:      logResponses,
		Total:     total,
		TotalPages: totalPages,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(email, nickname, password, role string) (*UserResponse, error) {
	// 验证必填字段
	if email == "" {
		return nil, errors.New("邮箱不能为空")
	}
	if nickname == "" {
		return nil, errors.New("昵称不能为空")
	}
	if password == "" {
		return nil, errors.New("密码不能为空")
	}

	// 验证邮箱格式
	if !strings.Contains(email, "@") {
		return nil, errors.New("邮箱格式不正确")
	}

	// 验证角色
	if role != "admin" && role != "user" {
		role = "user"
	}

	// 验证密码强度 (至少6位)
	if len(password) < 6 {
		return nil, errors.New("密码至少6位")
	}

	// 生成UUID
	userID := uuid.New().String()

	// 密码哈希
	passwordHash, err := passwordutil.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建用户
	user := &models.User{
		ID:           userID,
		Email:        email,
		PasswordHash: passwordHash,
		Nickname:     nickname,
		Status:       models.UserStatusActive,
	}

	// 创建用户到数据库
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	// 更新用户角色
	if err := s.userRepo.UpdateUserRole(userID, role); err != nil {
		return nil, fmt.Errorf("设置用户角色失败: %v", err)
	}

	// 返回用户信息
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Nickname,
		Email:     user.Email,
		Role:      role,
		Status:    "active",
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// UpdatePassword 修改密码
func (s *UserService) UpdatePassword(userID, newPassword string) error {
	// 验证必填字段
	if userID == "" {
		return errors.New("用户ID不能为空")
	}
	if newPassword == "" {
		return errors.New("新密码不能为空")
	}

	// 验证密码强度 (至少6位)
	if len(newPassword) < 6 {
		return errors.New("密码至少6位")
	}

	// 密码哈希
	passwordHash, err := passwordutil.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(userID, passwordHash); err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}

	return nil
}
