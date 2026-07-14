package middleware

import (
	"context"
	"net/http"
	"strings"

	"hot-ai-backend/internal/utils/jwtutil"
)

// ContextKey 上下文键类型
type ContextKey string

const (
	CtxUserID   ContextKey = "userID"
	CtxUserRoles ContextKey = "roles"
	CtxUserLevel ContextKey = "userLevel"
)

// LevelResolver 根据用户 ID 实时查询访问级别。
// 返回值：
//   - >0：使用此值覆盖 JWT 里的旧 level（用于处理「DB 里已是会员但 JWT 没刷新」的情况）
//   - 0：忽略，沿用 JWT 的 level
type LevelResolver func(ctx context.Context, userID string) int

// OptionalAuth 有 token 就解析 + 塞 ctx；没有或非法就当匿名（level=0）
// 不像 AuthMiddleware 那样 reject —— 业务接口要的就是「匿名也能进」
//
// resolve 可为 nil；非 nil 时，当 JWT 里 level < 2 且能拿到 userID，
// 会回调一次 resolver 拿真实 level（处理用户升级会员后没重新登录导致 JWT level 过期的场景）。
func OptionalAuth(jwtSecret string, resolve LevelResolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			level := 0
			userID := ""
			roles := []string{}

			if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
				tokenString := strings.TrimPrefix(h, "Bearer ")
				if claims, err := jwtutil.ParseAccessToken(jwtutil.Config{Secret: jwtSecret}, tokenString); err == nil {
					userID = claims.UserID
					roles = claims.Roles
					level = claims.Level
				}
				// 解析失败不报错 —— 当匿名继续
			}

			// JWT 里 level 不是 2 但能拿到 userID 时，回查一次 DB 拿最新 level
			// （典型场景：用户被升级会员 / 会员过期，但 token 还是旧 level）
			if resolve != nil && level < 2 && userID != "" {
				if fresh := resolve(r.Context(), userID); fresh > level {
					level = fresh
				}
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, CtxUserID, userID)
			ctx = context.WithValue(ctx, CtxUserRoles, roles)
			ctx = context.WithValue(ctx, CtxUserLevel, level)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext 从上下文获取用户 ID（可能为空字符串，表示匿名）
func GetUserIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(CtxUserID).(string); ok {
		return v
	}
	return ""
}

// GetUserRolesFromContext 从上下文获取用户角色
func GetUserRolesFromContext(ctx context.Context) []string {
	if v, ok := ctx.Value(CtxUserRoles).([]string); ok {
		return v
	}
	return []string{}
}

// GetUserLevelFromContext 从上下文获取用户访问级别（0=匿名, 1=普通, 2=会员/管理员）
func GetUserLevelFromContext(ctx context.Context) int {
	if v, ok := ctx.Value(CtxUserLevel).(int); ok {
		return v
	}
	return 0
}