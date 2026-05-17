package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"hot-ai-backend/internal/service"
	"hot-ai-backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// AdminChapterHandler 管理后台章节处理器
type AdminChapterHandler struct {
	svc *service.AdminService
}

// NewAdminChapterHandler 创建章节处理器实例
func NewAdminChapterHandler(svc *service.AdminService) *AdminChapterHandler {
	return &AdminChapterHandler{svc: svc}
}

// GetChapters 获取章节列表
func (h *AdminChapterHandler) GetChapters(w http.ResponseWriter, r *http.Request) {
	pathIDStr := getPathValue(r, "path_id")
	if pathIDStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	pathID, err := strconv.ParseUint(pathIDStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	chapters, err := h.svc.GetChapters(uint(pathID))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(chapters))
}

// GetChapterByID 获取章节详情
func (h *AdminChapterHandler) GetChapterByID(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少章节ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的章节ID"))
		return
	}

	chapter, err := h.svc.GetChapterByID(uint(id))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(chapter))
}

// CreateChapter 创建章节
func (h *AdminChapterHandler) CreateChapter(w http.ResponseWriter, r *http.Request) {
	pathIDStr := getPathValue(r, "path_id")
	if pathIDStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	pathID, err := strconv.ParseUint(pathIDStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	var req service.CreateChapterRequest
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	if req.Title == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("章节标题不能为空"))
		return
	}
	if req.ContentType == "" {
		req.ContentType = "article"
	}

	chapter, err := h.svc.CreateChapter(uint(pathID), &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(chapter))
}

// UpdateChapter 更新章节
func (h *AdminChapterHandler) UpdateChapter(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少章节ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的章节ID"))
		return
	}

	var req service.UpdateChapterRequest
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	chapter, err := h.svc.UpdateChapter(uint(id), &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(chapter))
}

// DeleteChapter 删除章节
func (h *AdminChapterHandler) DeleteChapter(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少章节ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的章节ID"))
		return
	}

	if err := h.svc.DeleteChapter(uint(id)); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "删除成功",
	}))
}
