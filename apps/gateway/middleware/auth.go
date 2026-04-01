package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"hot-ai-backend/internal/utils/jwtutil"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	jwtConfig jwtutil.Config
}

// NewAuthMiddleware 创建认证中间件实例
func NewAuthMiddleware(secret string, expireSeconds int64) *AuthMiddleware {
	return &AuthMiddleware{
		jwtConfig: jwtutil.Config{
			Secret:       secret,
			AccessExpire: time.Duration(expireSeconds) * time.Second,
			Issuer:       "hot-ai-platform",
		},
	}
}

// Handle 认证处理
func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 Authorization header 获取 Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("missing authorization header"))
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			httpx.ErrorCtx(r.Context(), w, errors.New("invalid authorization format"))
			return
		}

		tokenString := parts[1]

		// 验证 Token
		claims, err := jwtutil.ParseAccessToken(m.jwtConfig, tokenString)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.New("invalid or expired token"))
			return
		}

		// 将用户信息存入上下文
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "roles", claims.Roles)

		next(w, r.WithContext(ctx))
	}
}

// GetUserIDFromContext 从上下文获取用户 ID
func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return ""
	}
	return userID
}

// GetRolesFromContext 从上下文获取用户角色
func GetRolesFromContext(ctx context.Context) []string {
	roles, ok := ctx.Value("roles").([]string)
	if !ok {
		return []string{}
	}
	return roles
}
