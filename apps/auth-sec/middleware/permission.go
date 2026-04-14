package middleware

import (
	"errors"
	"net/http"
	"strings"

	"hot-ai-backend/internal/repository"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// PermissionMiddleware 权限中间件
type PermissionMiddleware struct {
	userRepo repository.UserRepository
}

// NewPermissionMiddleware 创建权限中间件实例
func NewPermissionMiddleware(userRepo repository.UserRepository) *PermissionMiddleware {
	return &PermissionMiddleware{
		userRepo: userRepo,
	}
}

// RequirePermission 需要指定权限
func (m *PermissionMiddleware) RequirePermission(permissionName string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserIDFromContext(r.Context())
			if userID == "" {
				httpx.ErrorCtx(r.Context(), w, errors.New("user not found in context"))
				return
			}

			// 检查用户是否有权限
			hasPermission, err := m.userRepo.HasPermission(r.Context(), userID, permissionName)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, errors.New("failed to check permission"))
				return
			}

			if !hasPermission {
				httpx.ErrorCtx(r.Context(), w, errors.New("permission denied"))
				return
			}

			next(w, r)
		}
	}
}

// RequireRole 需要指定角色
func (m *PermissionMiddleware) RequireRole(roleName string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			roles := GetRolesFromContext(r.Context())
			
			hasRole := false
			for _, role := range roles {
				if role == roleName {
					hasRole = true
					break
				}
			}

			if !hasRole {
				httpx.ErrorCtx(r.Context(), w, errors.New("role required: "+roleName))
				return
			}

			next(w, r)
		}
	}
}

// RequireAnyRole 需要任意一个指定角色
func (m *PermissionMiddleware) RequireAnyRole(roleNames ...string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			roles := GetRolesFromContext(r.Context())
			
			hasRole := false
			for _, role := range roles {
				for _, requiredRole := range roleNames {
					if role == requiredRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				httpx.ErrorCtx(r.Context(), w, errors.New("one of the following roles required: "+strings.Join(roleNames, ", ")))
				return
			}

			next(w, r)
		}
	}
}
