package repository

import (
	"context"
	"errors"
	"time"

	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/models"

	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 基础 CRUD
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error

	// 查询
	FindUsers(ctx context.Context, offset, limit int, status *models.UserStatus) ([]*models.User, int64, error)
	UpdateLastLogin(ctx context.Context, id, ip string) error
	GetUserCount() (int64, error)

	// 角色管理
	AddRole(ctx context.Context, userID, roleID string) error
	RemoveRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*models.Role, error)

	// 权限查询
	GetUserPermissions(ctx context.Context, userID string) ([]*models.Permission, error)
	HasPermission(ctx context.Context, userID, permissionName string) (bool, error)

	// 管理后台方法
	GetUsers(page, pageSize int, role, status string) ([]models.User, int64, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUserStatus(id string, status models.UserStatus) error
	UpdateUserRole(userID string, role string) error
	CreateAdminLog(log *models.AdminOperationLog) error
	GetAdminLogs(targetUserID string, page, pageSize int) ([]models.AdminOperationLog, int64, error)
	CreateUser(user *models.User) error
	UpdatePassword(userID string, passwordHash string) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.GetDB(),
	}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据 ID 获取用户
func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// FindUsers 查询用户列表
func (r *userRepository) FindUsers(ctx context.Context, offset, limit int, status *models.UserStatus) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return users, 0, nil
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

// UpdateLastLogin 更新最后登录信息
func (r *userRepository) UpdateLastLogin(ctx context.Context, id, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at": now,
			"last_login_ip": ip,
		}).Error
}

// AddRole 为用户添加角色
func (r *userRepository) AddRole(ctx context.Context, userID, roleID string) error {
	userRole := &models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.WithContext(ctx).Create(userRole).Error
}

// RemoveRole 移除用户角色
func (r *userRepository) RemoveRole(ctx context.Context, userID, roleID string) error {
	return r.db.WithContext(ctx).Delete(&models.UserRole{}, "user_id = ? AND role_id = ?", userID, roleID).Error
}

// GetUserRoles 获取用户角色列表
func (r *userRepository) GetUserRoles(ctx context.Context, userID string) ([]*models.Role, error) {
	var roles []*models.Role
	err := r.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

// GetUserPermissions 获取用户权限列表
func (r *userRepository) GetUserPermissions(ctx context.Context, userID string) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Group("permissions.id").
		Find(&permissions).Error
	return permissions, err
}

// HasPermission 检查用户是否有指定权限
func (r *userRepository) HasPermission(ctx context.Context, userID, permissionName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.name = ?", userID, permissionName).
		Count(&count).Error
	return count > 0, err
}

// GetUserCount 获取用户总数
func (r *userRepository) GetUserCount() (int64, error) {
	var total int64
	if err := r.db.Model(&models.User{}).Where("status = ?", models.UserStatusActive).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// GetUsers 获取用户列表（分页、筛选）
func (r *userRepository) GetUsers(page, pageSize int, role, status string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	// 角色筛选 - 需要 JOIN user_roles 和 roles 表
	if role != "" && role != "all" {
		query = query.
			Joins("JOIN user_roles ON users.id = user_roles.user_id").
			Joins("JOIN roles ON user_roles.role_id = roles.id").
			Where("roles.name = ?", role)
	}

	// 状态筛选
	if status != "" && status != "all" {
		query = query.Where("users.status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return users, 0, nil
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("users.created_at DESC").Find(&users).Error
	return users, total, err
}

// GetUserByID 根据ID获取用户
func (r *userRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserStatus 更新用户状态 (active -> inactive/banned, or vice versa)
func (r *userRepository) UpdateUserStatus(id string, status models.UserStatus) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateUserRole 更新用户角色 (需要操作 user_roles 表)
// role: "admin" 或 "user"
func (r *userRepository) UpdateUserRole(userID string, role string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除用户的所有角色
		if err := tx.Delete(&models.UserRole{}, "user_id = ?", userID).Error; err != nil {
			return err
		}

		// 如果不是空角色，则添加新角色
		if role != "" {
			// 查询角色ID
			var roleRecord models.Role
			if err := tx.Where("name = ?", role).First(&roleRecord).Error; err != nil {
				return err
			}

			// 添加新角色
			userRole := &models.UserRole{
				UserID: userID,
				RoleID: roleRecord.ID,
			}
			if err := tx.Create(userRole).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// CreateAdminLog 创建管理操作日志
func (r *userRepository) CreateAdminLog(log *models.AdminOperationLog) error {
	return r.db.Create(log).Error
}

// GetAdminLogs 获取操作日志
func (r *userRepository) GetAdminLogs(targetUserID string, page, pageSize int) ([]models.AdminOperationLog, int64, error) {
	var logs []models.AdminOperationLog
	var total int64

	query := r.db.Model(&models.AdminOperationLog{})

	if targetUserID != "" {
		query = query.Where("target_user_id = ?", targetUserID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return logs, 0, nil
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

// CreateUser 创建用户
func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// UpdatePassword 更新密码
func (r *userRepository) UpdatePassword(userID string, passwordHash string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("password_hash", passwordHash).Error
}
