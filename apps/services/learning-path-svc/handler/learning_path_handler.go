package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"hot-ai-backend/internal/service"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// LearningPathHandler 学习路径处理器
type LearningPathHandler struct {
	learningPathService *service.LearningPathService
}

// NewLearningPathHandler 创建学习路径处理器实例
func NewLearningPathHandler(learningPathService *service.LearningPathService) *LearningPathHandler {
	return &LearningPathHandler{
		learningPathService: learningPathService,
	}
}

// getPathValue 获取路径参数（兼容 Go 1.21）
func getPathValue(r *http.Request, name string) string {
	// 优先使用 Go 1.22+ PathValue
	val := r.PathValue(name)
	if val != "" {
		return val
	}
	// 回退：手动从 URL path 中提取最后一个/后的内容
	path := r.URL.Path
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			val = path[i+1:]
			break
		}
	}
	return val
}

// GetLearningPaths 获取学习路径列表
func (h *LearningPathHandler) GetLearningPaths(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 12

	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	difficulty := r.URL.Query().Get("difficulty")

	resp, err := h.learningPathService.GetLearningPaths(&service.GetLearningPathsRequest{
		Page:       page,
		PageSize:   pageSize,
		Difficulty: difficulty,
	})
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, resp)
}

// GetLearningPathByID 根据ID获取学习路径详情
func (h *LearningPathHandler) GetLearningPathByID(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	path, err := h.learningPathService.GetLearningPathByID(uint(id))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, path)
}

// GetLearningPathBySlug 根据slug获取学习路径详情
func (h *LearningPathHandler) GetLearningPathBySlug(w http.ResponseWriter, r *http.Request) {
	slug := getPathValue(r, "slug")
	if slug == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径slug"))
		return
	}

	path, err := h.learningPathService.GetLearningPathBySlug(slug)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, path)
}

// GetPathChapters 获取路径的所有章节
func (h *LearningPathHandler) GetPathChapters(w http.ResponseWriter, r *http.Request) {
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

	chapters, err := h.learningPathService.GetPathChapters(uint(pathID))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, chapters)
}

// GetChapterByID 根据章节ID获取详情
func (h *LearningPathHandler) GetChapterByID(w http.ResponseWriter, r *http.Request) {
	chapterIDStr := getPathValue(r, "chapter_id")
	if chapterIDStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少章节ID"))
		return
	}

	chapterID, err := strconv.ParseUint(chapterIDStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的章节ID"))
		return
	}

	chapter, err := h.learningPathService.GetChapterByID(uint(chapterID))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, chapter)
}

// GetChapterBySlug 根据路径slug和章节slug获取章节详情
func (h *LearningPathHandler) GetChapterBySlug(w http.ResponseWriter, r *http.Request) {
	pathSlug := getPathValue(r, "path_slug")
	chapterSlug := getPathValue(r, "chapter_slug")
	pathIDStr := r.URL.Query().Get("path_id")

	if pathSlug == "" || chapterSlug == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径slug或章节slug"))
		return
	}

	// 获取路径信息
	path, err := h.learningPathService.GetLearningPathBySlug(pathSlug)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("路径不存在"))
		return
	}

	// 优先使用URL参数中的path_id
	pathID := path.ID
	if pathIDStr != "" {
		if id, err := strconv.ParseUint(pathIDStr, 10, 32); err == nil {
			pathID = uint(id)
		}
	}

	chapter, err := h.learningPathService.GetChapterBySlug(pathID, chapterSlug)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 获取前一章和下一章
	prev, next, _ := h.learningPathService.GetPrevNextChapter(pathID, chapter.ID)

	result := map[string]interface{}{
		"chapter": chapter,
		"prev":    prev,
		"next":    next,
	}

	httpx.OkJsonCtx(r.Context(), w, result)
}

// GetFeaturedPaths 获取推荐路径
func (h *LearningPathHandler) GetFeaturedPaths(w http.ResponseWriter, r *http.Request) {
	limit := 4
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	paths, err := h.learningPathService.GetFeaturedPaths(limit)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, paths)
}

// GetLevelInfo 获取难度等级信息
func (h *LearningPathHandler) GetLevelInfo(w http.ResponseWriter, r *http.Request) {
	levelInfo := h.learningPathService.GetLevelInfo()
	httpx.OkJsonCtx(r.Context(), w, levelInfo)
}

// GetPathProgress 获取用户的学习进度
func (h *LearningPathHandler) GetPathProgress(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	pathIDStr := r.URL.Query().Get("path_id")

	if pathIDStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	pathID, err := strconv.ParseUint(pathIDStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	progress, err := h.learningPathService.GetPathProgress(userID, uint(pathID))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, progress)
}

// GetCompletedChapters 获取用户已完成的章节列表
func (h *LearningPathHandler) GetCompletedChapters(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	pathIDStr := r.URL.Query().Get("path_id")

	if pathIDStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	pathID, err := strconv.ParseUint(pathIDStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	chapterIDs, err := h.learningPathService.GetCompletedChapters(userID, uint(pathID))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
		"completed_chapters": chapterIDs,
	})
}

// SaveProgress 保存学习进度
func (h *LearningPathHandler) SaveProgress(w http.ResponseWriter, r *http.Request) {
	var req service.SaveProgressRequest
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	if err := h.learningPathService.SaveProgress(&req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
		"success": true,
		"message": "进度已保存",
	})
}

// GetPathDashboard 获取路径学习仪表盘（综合统计）
func (h *LearningPathHandler) GetPathDashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	pathIDStr := r.URL.Query().Get("path_id")

	if pathIDStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少路径ID"))
		return
	}

	pathID, err := strconv.ParseUint(pathIDStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的路径ID"))
		return
	}

	// 获取路径信息
	path, err := h.learningPathService.GetLearningPathByID(uint(pathID))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 获取章节列表
	chapters, err := h.learningPathService.GetPathChapters(uint(pathID))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	// 获取已完成的章节
	completedChapters, err := h.learningPathService.GetCompletedChapters(userID, uint(pathID))
	if err != nil {
		completedChapters = []uint{}
	}

	httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
		"path":               path,
		"chapters":           chapters,
		"completed_chapters": completedChapters,
		"progress": map[string]interface{}{
			"total_chapters":      len(chapters),
			"completed_count":     len(completedChapters),
			"progress_percentage": float64(len(completedChapters)) / float64(len(chapters)) * 100,
		},
	})
}
