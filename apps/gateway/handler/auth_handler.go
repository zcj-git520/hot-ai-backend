package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"hot-ai-backend/internal/service"
	"hot-ai-backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// SetAuthService 设置认证服务（依赖注入）
func (h *AuthHandler) SetAuthService(authService *service.AuthService) {
	h.authService = authService
}

// Login 用户登录
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 1. 前端校验（后端再次校验）
	if err := validateLoginRequest(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 2. 获取客户端 IP 地址
	ip := getClientIP(r)

	// 3. 转换为 service 层的请求类型
	serviceReq := &service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// 4. 调用服务层处理登录逻辑
	resp, err := h.authService.Login(r.Context(), serviceReq, ip)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 5. 返回成功响应
	httpx.OkJsonCtx(r.Context(), w, types.Success(resp))
}

// Register 用户注册
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 1. 前端校验（后端再次校验）
	if err := validateRegisterRequest(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 2. 转换为 service 层的请求类型
	serviceReq := &service.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	}

	// 3. 调用服务层处理注册逻辑（需要验证码）
	resp, err := h.authService.Register(r.Context(), serviceReq, req.VerificationCode)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 4. 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.SuccessWithMessage("注册成功", resp))
}

// validateRegisterRequest 校验注册请求
// validateLoginRequest 校验登录请求
func validateLoginRequest(req *types.LoginRequest) error {
	// 邮箱格式校验
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return types.NewBadRequestError("邮箱格式不正确")
	}

	// 密码非空校验
	if len(req.Password) < 6 {
		return types.NewBadRequestError("密码长度不能少于 6 个字符")
	}

	return nil
}

// getClientIP 获取客户端真实 IP 地址
func getClientIP(r *http.Request) string {
	// 尝试从 X-Forwarded-For 获取（适用于代理/负载均衡场景）
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For 可能包含多个 IP，取第一个
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}

	// 尝试从 X-Real-IP 获取
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// 从 RemoteAddr 获取（直连场景）
	ip := r.RemoteAddr
	if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
		ip = ip[:colonIndex]
	}

	// IPv6 本地回环地址处理
	if ip == "::1" || ip == "[::1]" {
		return "127.0.0.1"
	}

	return ip
}

func validateRegisterRequest(req *types.RegisterRequest) error {
	// 邮箱校验：标准邮箱格式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return types.NewBadRequestError("邮箱格式不正确")
	}

	// 密码确认校验：两次输入的密码必须一致
	if req.Password != req.PasswordConfirm {
		return types.NewBadRequestError("两次输入的密码不一致")
	}

	// 密码强度校验：8-20 位，包含大小写字母 + 数字 + 特殊字符
	if len(req.Password) < 8 || len(req.Password) > 20 {
		return types.NewBadRequestError("密码长度必须在 8-20 位之间")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	for _, char := range req.Password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return types.NewBadRequestError("密码必须包含大写字母")
	}
	if !hasLower {
		return types.NewBadRequestError("密码必须包含小写字母")
	}
	if !hasDigit {
		return types.NewBadRequestError("密码必须包含数字")
	}
	if !hasSpecial {
		return types.NewBadRequestError("密码必须包含特殊字符")
	}

	// 昵称校验：2-20 个字符，允许中文、字母、数字
	if len(req.Nickname) < 2 || len(req.Nickname) > 20 {
		return types.NewBadRequestError("昵称长度必须在 2-20 个字符之间")
	}

	return nil
}

// SendVerificationEmail 发送验证邮件
func (h *AuthHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 邮箱格式校验
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError("邮箱格式不正确"))
		return
	}

	// 调用服务层发送验证邮件（用于已注册用户验证邮箱）
	err := h.authService.SendVerificationEmail(r.Context(), req.Email)
	if err != nil {
		// 根据错误类型返回不同的状态码
		if err.Error() == "用户不存在" {
			httpx.ErrorCtx(r.Context(), w, types.NewNotFoundError(err.Error()))
		} else if err.Error() == "邮箱已验证" {
			httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError(err.Error()))
		} else {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		return
	}

	// 返回成功
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.SuccessWithMessage("验证邮件已发送，请检查邮箱", map[string]interface{}{
		"email": req.Email,
	}))
}

// SendRegistrationCode 发送注册验证码
func (h *AuthHandler) SendRegistrationCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 邮箱格式校验
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError("邮箱格式不正确"))
		return
	}

	// 调用服务层发送注册验证码
	err := h.authService.SendRegistrationCode(r.Context(), req.Email)
	if err != nil {
		// 根据错误类型返回不同的状态码
		if err.Error() == "该邮箱已注册" {
			httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError(err.Error()))
		} else if err.Error() == "操作过于频繁" || err.Error() == "操作过于频繁，请 1 小时后再试" {
			httpx.ErrorCtx(r.Context(), w, &types.APIError{
				Code:    429,
				Message: err.Error(),
			})
		} else {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		return
	}

	// 返回成功（不返回验证码，提示用户查看邮箱）
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.SuccessWithMessage("验证码已发送，请检查邮箱", map[string]interface{}{
		"email":      req.Email,
		"expires_in": 300, // 5 分钟
	}))
}

// VerifyEmail 验证邮箱
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError("缺少验证 token"))
		return
	}

	// 调用服务层验证邮箱
	err := h.authService.VerifyEmail(r.Context(), token)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 返回成功
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.Success("邮箱验证成功"))
}

// RefreshToken 刷新 Token
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req types.RefreshTokenRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 1. 校验 Refresh Token 非空
	if req.RefreshToken == "" {
		httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError("Refresh Token 不能为空"))
		return
	}

	// 2. 转换为 service 层的请求类型
	serviceReq := &service.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	// 3. 调用服务层处理刷新 Token 逻辑
	resp, err := h.authService.RefreshToken(r.Context(), serviceReq)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 4. 返回成功响应
	httpx.OkJsonCtx(r.Context(), w, types.Success(resp))
}

// Logout 用户登出
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. 从请求头中获取 Access Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError("缺少认证信息"))
		return
	}

	// 2. 提取 Token（格式：Bearer <token>）
	var accessToken string
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		accessToken = parts[1]
	} else {
		httpx.ErrorCtx(r.Context(), w, types.NewBadRequestError("Token 格式不正确"))
		return
	}

	// 3. 调用服务层处理登出逻辑
	err := h.authService.Logout(r.Context(), accessToken)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 4. 返回成功响应
	httpx.OkJsonCtx(r.Context(), w, types.SuccessWithMessage("登出成功", nil))
}
