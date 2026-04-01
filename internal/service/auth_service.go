package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hot-ai-backend/internal/cache"
	"hot-ai-backend/internal/models"
	"hot-ai-backend/internal/repository"
	"hot-ai-backend/internal/utils/jwtutil"
	"hot-ai-backend/internal/utils/mailutil"
	"hot-ai-backend/internal/utils/passwordutil"

	"github.com/google/uuid"
)

// AuthService 认证服务
type AuthService struct {
	userRepo       repository.UserRepository
	favoriteRepo   repository.FavoriteRepository
	jwtConfig      jwtutil.Config
	mailSecret     string // 邮件验证 token 密钥
	emailCache     *cache.EmailCache // 邮箱验证码缓存
	smtpConfig     *mailutil.SMTPConfig // SMTP 配置
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository, favoriteRepo repository.FavoriteRepository, secret string) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		favoriteRepo: favoriteRepo,
		jwtConfig: jwtutil.Config{
			Secret:        secret,
			AccessExpire:  24 * time.Hour,
			RefreshExpire: 7 * 24 * time.Hour,
			Issuer:        "hot-ai-platform",
		},
		mailSecret: secret, // 使用相同的密钥，也可以单独配置
	}
}

// SetEmailCache 设置邮箱缓存（依赖注入）
func (s *AuthService) SetEmailCache(emailCache *cache.EmailCache) {
	s.emailCache = emailCache
}

// SetSMTPConfig 设置 SMTP 配置
func (s *AuthService) SetSMTPConfig(config *mailutil.SMTPConfig) {
	s.smtpConfig = config
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    int64          `json:"expires_in"`
	User         *UserInfo      `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Nickname  string   `json:"nickname"`
	Avatar    string   `json:"avatar,omitempty"`
	Roles     []string `json:"roles"`
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse 刷新 Token 响应
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest, verificationCode string) (*RegisterResponse, error) {
	// 1. 验证邮箱验证码
	if s.emailCache != nil {
		err := s.emailCache.VerifyCode(ctx, req.Email, verificationCode)
		if err != nil {
			return nil, fmt.Errorf("验证码错误：%w", err)
		}
	}

	// 2. 检查邮箱是否已存在
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// 密码强度校验
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// 加密密码
	passwordHash, err := passwordutil.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		Nickname:     req.Nickname,
		Status:       models.UserStatusActive,
		Preferences: models.UserPreferences{
			EmailNotifications: models.EmailNotificationSettings{
				SystemUpdate: true,
				NewContent:   true,
				WeeklyDigest: true,
			},
			Theme:    "light",
			Language: "zh-CN",
		},
	}

	// 分配默认角色（user）
	defaultRoleID := "role_user"

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 分配默认角色
	err = s.userRepo.AddRole(ctx, user.ID, defaultRoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	return &RegisterResponse{
		UserID:  user.ID,
		Message: "注册成功，请登录",
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest, ip string) (*LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// 验证账号状态
	if user.Status != models.UserStatusActive {
		return nil, errors.New("account is inactive or banned")
	}

	// 验证密码
	if !passwordutil.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// 获取用户角色
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// 生成 Token
	tokenPair, err := jwtutil.GenerateTokenPair(s.jwtConfig, user.ID, roleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 更新最后登录信息
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID, ip)

	return &LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: &UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Roles:    roleNames,
		},
	}, nil
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// 解析 Refresh Token
	claims, err := jwtutil.ParseRefreshToken(s.jwtConfig, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 获取用户
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 验证账号状态
	if user.Status != models.UserStatusActive {
		return nil, errors.New("account is inactive or banned")
	}

	// 获取用户角色
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	// 生成新的 Access Token
	newTokenPair, err := jwtutil.GenerateTokenPair(s.jwtConfig, user.ID, roleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &RefreshTokenResponse{
		AccessToken: newTokenPair.AccessToken,
		ExpiresIn:   newTokenPair.ExpiresIn,
	}, nil
}

// Logout 登出
func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	// TODO: 将 Token 加入黑名单（Redis 实现）
	// 这里暂时只记录日志
	return nil
}

// GetUserProfile 获取用户信息
func (s *AuthService) GetUserProfile(ctx context.Context, userID string) (*UserInfo, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 获取用户角色
	roles, err := s.userRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return &UserInfo{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Roles:    roleNames,
	}, nil
}

// SendVerificationEmail 发送验证邮件
func (s *AuthService) SendVerificationEmail(ctx context.Context, email string) error {
	// 1. 获取用户
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 2. 检查邮箱是否已验证
	if user.EmailVerified {
		return errors.New("邮箱已验证")
	}

	// 3. 生成验证 token
	token, err := mailutil.GenerateVerificationToken(user.ID, user.Email, s.mailSecret)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// 4. 构建验证链接
	verifyURL := fmt.Sprintf(
		"http://localhost:8080/api/auth/verify-email?token=%s",
		token,
	)

	// 5. TODO: 发送邮件
	// 这里暂时只记录日志，实际项目中需要集成 SMTP 或第三方邮件服务
	fmt.Printf("[邮件发送] 收件人：%s\n", user.Email)
	fmt.Printf("[邮件发送] 主题：验证您的 Hot AI 账号\n")
	fmt.Printf("[邮件发送] 验证链接：%s\n", verifyURL)
	fmt.Printf("[邮件发送] 过期时间：%s\n", time.Now().Add(24*time.Hour).Format("2006-01-02 15:04:05"))

	// 实际项目中可以使用以下方式发送邮件：
	// - 使用 net/smtp 包直接通过 SMTP 发送
	// - 使用第三方服务（SendGrid、Mailgun、阿里云邮件推送等）
	// - 使用消息队列异步发送

	return nil
}

// VerifyEmail 验证邮箱
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	// 1. 解析并验证 token
	tokenData, err := mailutil.ParseVerificationToken(token, s.mailSecret)
	if err != nil {
		return fmt.Errorf("无效的验证链接：%w", err)
	}

	// 2. 获取用户
	user, err := s.userRepo.GetByID(ctx, tokenData.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 3. 验证邮箱是否匹配
	if user.Email != tokenData.Email {
		return errors.New("邮箱地址不匹配")
	}

	// 4. 检查是否已验证
	if user.EmailVerified {
		return errors.New("邮箱已验证")
	}

	// 5. 更新邮箱验证状态
	user.EmailVerified = true
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// SendRegistrationCode 发送注册验证码
func (s *AuthService) SendRegistrationCode(ctx context.Context, email string) error {
	// 1. 检查邮箱是否已注册
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return errors.New("该邮箱已注册")
	}

	// 2. 检查发送频率
	if s.emailCache != nil {
		err = s.emailCache.CheckSendFrequency(ctx, email)
		if err != nil {
			return err
		}
	}

	// 3. 生成 6 位数字验证码
	code := mailutil.GenerateVerificationCode()

	// 4. 存储验证码到 Redis（5 分钟有效期）
	if s.emailCache != nil {
		err = s.emailCache.StoreVerificationCode(ctx, email, code, 5*time.Minute)
		if err != nil {
			return fmt.Errorf("failed to store code: %w", err)
		}

		// 记录发送次数
		err = s.emailCache.RecordSendAttempt(ctx, email)
		if err != nil {
			fmt.Printf("[警告] 记录发送次数失败：%v\n", err)
		}
	}

	// 5. 发送邮件验证码
	err = mailutil.SendVerificationEmail(s.smtpConfig, email, code)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
