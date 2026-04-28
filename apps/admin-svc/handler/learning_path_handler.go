package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"hot-ai-backend/apps/admin-svc/internal/logic"
	"hot-ai-backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// LearningPathHandler 学习路径管理处理器
type LearningPathHandler struct {
	learningPathLogic *logic.LearningPathLogic
}

// NewLearningPathHandler 创建学习路径处理器实例
func NewLearningPathHandler() *LearningPathHandler {
	return &LearningPathHandler{
		learningPathLogic: logic.NewLearningPathLogic(),
	}
}

// GetLearningPaths 获取学习路径列表
func (h *LearningPathHandler) GetLearningPaths(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 20

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
		pageSize = 20
	}

	difficulty := r.URL.Query().Get("difficulty")
	search := r.URL.Query().Get("search")

	var status *int
	if s := r.URL.Query().Get("status"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil {
			status = &parsed
		}
	}

	req := &logic.GetLearningPathsRequest{
		Page:       page,
		PageSize:   pageSize,
		Difficulty: difficulty,
		Search:     search,
		Status:     status,
	}

	paths, total, err := h.learningPathLogic.GetLearningPaths(req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"list":  paths,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	}))
}

// GetLearningPathByID 获取学习路径详情
func (h *LearningPathHandler) GetLearningPathByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	path, err := h.learningPathLogic.GetLearningPathByID(uint(id))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(path))
}

// CreateLearningPath 创建学习路径
func (h *LearningPathHandler) CreateLearningPath(w http.ResponseWriter, r *http.Request) {
	var req logic.CreateLearningPathRequest
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 验证必填字段
	if req.Title == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("标题不能为空"))
		return
	}
	if req.Difficulty == "" {
		req.Difficulty = "beginner"
	}
	if req.LevelLabel == "" {
		req.LevelLabel = "入门"
	}

	path, err := h.learningPathLogic.CreateLearningPath(&req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(path))
}

// UpdateLearningPath 更新学习路径
func (h *LearningPathHandler) UpdateLearningPath(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	var req logic.UpdateLearningPathRequest
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	path, err := h.learningPathLogic.UpdateLearningPath(uint(id), &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(path))
}

// DeleteLearningPath 删除学习路径
func (h *LearningPathHandler) DeleteLearningPath(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	if err := h.learningPathLogic.DeleteLearningPath(uint(id)); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "删除成功",
	}))
}

// GetChapters 获取章节列表
func (h *LearningPathHandler) GetChapters(w http.ResponseWriter, r *http.Request) {
	pathIDStr := r.PathValue("path_id")
	if pathIDStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	pathID, err := strconv.ParseUint(pathIDStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	chapters, err := h.learningPathLogic.GetChapters(uint(pathID))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(chapters))
}

// GetChapterByID 获取章节详情
func (h *LearningPathHandler) GetChapterByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少章节ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的章节ID"))
		return
	}

	chapter, err := h.learningPathLogic.GetChapterByID(uint(id))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(chapter))
}

// CreateChapter 创建章节
func (h *LearningPathHandler) CreateChapter(w http.ResponseWriter, r *http.Request) {
	pathIDStr := r.PathValue("path_id")
	if pathIDStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	pathID, err := strconv.ParseUint(pathIDStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	var req logic.CreateChapterRequest
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

	chapter, err := h.learningPathLogic.CreateChapter(uint(pathID), &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(chapter))
}

// UpdateChapter 更新章节
func (h *LearningPathHandler) UpdateChapter(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少章节ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的章节ID"))
		return
	}

	var req logic.UpdateChapterRequest
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	chapter, err := h.learningPathLogic.UpdateChapter(uint(id), &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(chapter))
}

// DeleteChapter 删除章节
func (h *LearningPathHandler) DeleteChapter(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少章节ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的章节ID"))
		return
	}

	if err := h.learningPathLogic.DeleteChapter(uint(id)); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "删除成功",
	}))
}

// PublishLearningPath 发布学习路径
func (h *LearningPathHandler) PublishLearningPath(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	if err := h.learningPathLogic.PublishLearningPath(uint(id)); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "发布成功",
	}))
}

// UnpublishLearningPath 下架学习路径
func (h *LearningPathHandler) UnpublishLearningPath(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	if err := h.learningPathLogic.UnpublishLearningPath(uint(id)); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "下架成功",
	}))
}

// SetFeatured 设置推荐状态
func (h *LearningPathHandler) SetFeatured(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	var req struct {
		Featured bool `json:"featured"`
	}
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	if err := h.learningPathLogic.SetFeatured(uint(id), req.Featured); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "设置成功",
	}))
}