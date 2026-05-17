package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"hot-ai-backend/internal/service"
	"hot-ai-backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// AdminLearningPathHandler 管理后台学习路径处理器
type AdminLearningPathHandler struct {
	svc *service.AdminService
}

// NewAdminLearningPathHandler 创建学习路径处理器实例
func NewAdminLearningPathHandler(svc *service.AdminService) *AdminLearningPathHandler {
	return &AdminLearningPathHandler{svc: svc}
}

// GetLearningPaths 获取学习路径列表
func (h *AdminLearningPathHandler) GetLearningPaths(w http.ResponseWriter, r *http.Request) {
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

	req := &service.AdminGetLearningPathsRequest{
		Page:       page,
		PageSize:   pageSize,
		Difficulty: difficulty,
		Search:     search,
		Status:     status,
	}

	paths, total, err := h.svc.GetLearningPaths(req)
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
func (h *AdminLearningPathHandler) GetLearningPathByID(w http.ResponseWriter, r *http.Request) {
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

	path, err := h.svc.GetLearningPathByID(uint(id))
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(path))
}

// CreateLearningPath 创建学习路径
func (h *AdminLearningPathHandler) CreateLearningPath(w http.ResponseWriter, r *http.Request) {
	var req service.CreateLearningPathRequest
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

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

	path, err := h.svc.CreateLearningPath(&req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(path))
}

// UpdateLearningPath 更新学习路径
func (h *AdminLearningPathHandler) UpdateLearningPath(w http.ResponseWriter, r *http.Request) {
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

	var req service.UpdateLearningPathRequest
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	path, err := h.svc.UpdateLearningPath(uint(id), &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(path))
}

// DeleteLearningPath 删除学习路径
func (h *AdminLearningPathHandler) DeleteLearningPath(w http.ResponseWriter, r *http.Request) {
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

	if err := h.svc.DeleteLearningPath(uint(id)); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "删除成功",
	}))
}

// SubmitReview 提交审核
func (h *AdminLearningPathHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
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

	if err := h.svc.SubmitReview(uint(id)); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "提交审核成功",
	}))
}

// Approve 审核通过
func (h *AdminLearningPathHandler) Approve(w http.ResponseWriter, r *http.Request) {
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

	if err := h.svc.Approve(uint(id)); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "审核通过",
	}))
}

// Reject 审核拒绝
func (h *AdminLearningPathHandler) Reject(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Reason string `json:"reason"`
	}
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	if err := h.svc.Reject(uint(id), req.Reason); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "已拒绝",
	}))
}

// Publish 发布学习路径
func (h *AdminLearningPathHandler) Publish(w http.ResponseWriter, r *http.Request) {
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

	if err := h.svc.PublishLearningPath(uint(id)); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "发布成功",
	}))
}

// Unpublish 下架学习路径
func (h *AdminLearningPathHandler) Unpublish(w http.ResponseWriter, r *http.Request) {
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

	if err := h.svc.UnpublishLearningPath(uint(id)); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "下架成功",
	}))
}

// SetFeatured 设置推荐状态
func (h *AdminLearningPathHandler) SetFeatured(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Featured bool `json:"featured"`
	}
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	if err := h.svc.SetFeatured(uint(id), req.Featured); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"success": true,
		"message": "设置成功",
	}))
}
