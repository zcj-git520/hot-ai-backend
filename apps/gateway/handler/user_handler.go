package handler

import (
	"errors"
	"net/http"

	"hot-ai-backend/apps/gateway/middleware"
	"hot-ai-backend/internal/service"

	"github.com/gorilla/mux"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService     *service.UserService
	favoriteService *service.FavoriteService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService *service.UserService, favoriteService *service.FavoriteService) *UserHandler {
	return &UserHandler{
		userService:     userService,
		favoriteService: favoriteService,
	}
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		httpx.ErrorCtx(r.Context(), w, errors.New("unauthorized"))
		return
	}

	resp, err := h.userService.GetUserProfile(r.Context(), userID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, resp)
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		httpx.ErrorCtx(r.Context(), w, errors.New("unauthorized"))
		return
	}

	var req service.UserProfileRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	err := h.userService.UpdateUserProfile(r.Context(), userID, &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
		"success": true,
	})
}

// UpdatePreferences 更新偏好设置
func (h *UserHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		httpx.ErrorCtx(r.Context(), w, errors.New("unauthorized"))
		return
	}

	var req service.UserPreferencesRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	err := h.userService.UpdateUserPreferences(r.Context(), userID, &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
		"success": true,
	})
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		httpx.ErrorCtx(r.Context(), w, errors.New("unauthorized"))
		return
	}

	var req service.ChangePasswordRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	err := h.userService.ChangePassword(r.Context(), userID, &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
		"success": true,
	})
}

// GetFavorites 获取收藏列表
func (h *UserHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		httpx.ErrorCtx(r.Context(), w, errors.New("unauthorized"))
		return
	}

	// 解析查询参数
	page := int64(1)
	pageSize := int64(20)
	favType := r.URL.Query().Get("type")
	var favTypePtr *string
	if favType != "" {
		favTypePtr = &favType
	}

	resp, err := h.favoriteService.GetFavoriteList(r.Context(), userID, page, pageSize, favTypePtr)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, resp)
}

// AddFavorite 添加收藏
func (h *UserHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		httpx.ErrorCtx(r.Context(), w, errors.New("unauthorized"))
		return
	}

	var req service.CreateFavoriteRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	favorite, err := h.favoriteService.CreateFavorite(r.Context(), userID, &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
		"success": true,
		"favorite": favorite,
	})
}

// DeleteFavorite 删除收藏
func (h *UserHandler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == "" {
		httpx.ErrorCtx(r.Context(), w, errors.New("unauthorized"))
		return
	}

	// 从 URL 路径获取收藏 ID - 使用 gorilla/mux 解析
	vars := mux.Vars(r)
	favoriteID := vars["id"]
	if favoriteID == "" {
		httpx.ErrorCtx(r.Context(), w, errors.New("missing favorite id"))
		return
	}

	err := h.favoriteService.DeleteFavorite(r.Context(), userID, favoriteID)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
		"success": true,
	})
}
