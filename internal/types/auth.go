package types

// LoginRequest 登录请求
type LoginRequest struct {
	// 邮箱
	Email string `json:"email" example:"user@example.com"`
	// 密码
	Password string `json:"password" example:"SecurePass123!"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	// 邮箱
	Email string `json:"email" example:"user@example.com"`
	// 密码
	Password string `json:"password" example:"SecurePass123!"`
	// 确认密码
	PasswordConfirm string `json:"password_confirm" example:"SecurePass123!"`
	// 昵称
	Nickname string `json:"nickname" example:"用户名"`
	// 邮箱验证码
	VerificationCode string `json:"verification_code" example:"123456"`
}

// UserInfo 用户信息
type UserInfo struct {
	// 用户 ID
	ID string `json:"id"`
	// 邮箱
	Email string `json:"email"`
	// 昵称
	Nickname string `json:"nickname"`
	// 头像 URL
	Avatar string `json:"avatar,omitempty"`
	// 角色列表
	Roles []string `json:"roles"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	// Access Token
	AccessToken string `json:"access_token"`
	// Refresh Token
	RefreshToken string `json:"refresh_token"`
	// 过期时间 (秒)
	ExpiresIn int64 `json:"expires_in"`
	// 用户信息
	User *UserInfo `json:"user"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	// 用户 ID
	UserID string `json:"user_id"`
	// 消息
	Message string `json:"message"`
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	// Refresh Token
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse 刷新 Token 响应
type RefreshTokenResponse struct {
	// Access Token
	AccessToken string `json:"access_token"`
	// 过期时间 (秒)
	ExpiresIn int64 `json:"expires_in"`
}
