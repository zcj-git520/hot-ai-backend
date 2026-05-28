package handler

import (
	"fmt"
	"net/http"

	"hot-ai-backend/internal/service"
	"hot-ai-backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// UpdateUserRoleRequest 更新角色请求
type UpdateUserRoleRequest struct {
	Role string `json:"role"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	Password string `json:"password"`
}

// UserHandler 用户管理处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers 获取用户列表
// GET /api/admin/users
// Query: page, page_size, search, role, status
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 12
	search := ""
	role := ""
	status := ""

	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}
	search = r.URL.Query().Get("search")
	role = r.URL.Query().Get("role")
	status = r.URL.Query().Get("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}

	resp, err := h.userService.GetUsers(page, pageSize, role, status, search)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(resp))
}

// GetUserByID 获取用户详情
// GET /api/admin/users/:id
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, ErrMissingID)
		return
	}

	user, err := h.userService.GetUserByID(idStr)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(user))
}

// UpdateUserRole 更新用户角色
// PUT /api/admin/users/:id/role
// Body: {"role": "admin" | "user"}
func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, ErrMissingID)
		return
	}

	var req UpdateUserRoleRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	err := h.userService.UpdateUserRole(idStr, req.Role)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(nil))
}

// DisableUser 禁用用户
// POST /api/admin/users/:id/disable
func (h *UserHandler) DisableUser(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, ErrMissingID)
		return
	}

	err := h.userService.DisableUser(idStr)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(nil))
}

// EnableUser 启用用户
// POST /api/admin/users/:id/enable
func (h *UserHandler) EnableUser(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, ErrMissingID)
		return
	}

	err := h.userService.EnableUser(idStr)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(nil))
}

// GetAdminLogs 获取操作日志
// GET /api/admin/users/:id/logs
// Query: page, page_size
func (h *UserHandler) GetAdminLogs(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, ErrMissingID)
		return
	}

	page := 1
	pageSize := 12

	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}

	resp, err := h.userService.GetAdminLogs(idStr, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(resp))
}

// CreateUser 创建用户
// POST /api/admin/users
// Body: {"email": "", "nickname": "", "password": "", "role": "user"}
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	user, err := h.userService.CreateUser(req.Email, req.Nickname, req.Password, req.Role)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(user))
}

// UpdatePassword 修改密码
// PUT /api/admin/users/:id/password
// Body: {"password": "newpassword"}
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, ErrMissingID)
		return
	}

	var req UpdatePasswordRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	err := h.userService.UpdatePassword(idStr, req.Password)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(nil))
}