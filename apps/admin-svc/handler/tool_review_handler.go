package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/models"
	"hot-ai-backend/internal/repository"
	"hot-ai-backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// AdminToolReviewHandler 管理后台工具审核处理器
type AdminToolReviewHandler struct {
	reviewRepo *repository.ToolReviewRepository
}

// NewAdminToolReviewHandler 创建工具审核处理器实例
func NewAdminToolReviewHandler() *AdminToolReviewHandler {
	return &AdminToolReviewHandler{
		reviewRepo: repository.NewToolReviewRepository(),
	}
}

// GetPendingTools 获取待审核工具列表
func (h *AdminToolReviewHandler) GetPendingTools(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	search := r.URL.Query().Get("search")
	reviewStatus := r.URL.Query().Get("review_status")

	// 查询工具列表，默认只查待审核
	var tools []models.Tool
	var total int64

	query := database.GetDB().Table("tools")
	if reviewStatus != "" {
		query = query.Where("review_status = ?", reviewStatus)
	}
	// 不传 review_status 参数时，默认显示所有工具
	_ = reviewStatus // 消除未使用警告

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Order("submitted_at DESC").Offset(offset).Limit(pageSize).Find(&tools).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(map[string]interface{}{
		"list":      tools,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

// GetToolDetailForReview 获取工具详情（含审核历史）
func (h *AdminToolReviewHandler) GetToolDetailForReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少工具ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("无效的工具ID"))
		return
	}

	// 查询工具详情
	var tool models.Tool
	if err := database.GetDB().Table("tools").Where("id = ?", uint(id)).First(&tool).Error; err != nil {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("工具不存在"))
		return
	}

	// 查询审核历史
	reviews, err := h.reviewRepo.GetByToolID(ctx, uint(id))
	if err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	result := map[string]interface{}{
		"tool":   tool,
		"reviews": reviews,
	}

	httpx.OkJsonCtx(ctx, w, types.Success(result))
}

// ApproveTool 审核通过
func (h *AdminToolReviewHandler) ApproveTool(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少工具ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("无效的工具ID"))
		return
	}

	// 获取当前管理员ID（从请求头或上下文获取，这里假设从header获取）
	adminID := r.Header.Get("X-Admin-ID")
	if adminID == "" {
		adminID = "system"
	}

	// 更新工具状态为已审核通过
	result := database.GetDB().Table("tools").Where("id = ?", id).Updates(map[string]interface{}{
		"review_status": "approved",
	})
	if result.Error != nil {
		httpx.ErrorCtx(ctx, w, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("工具不存在"))
		return
	}

	// 创建审核记录 (临时注释，等表创建后再启用)
	// record := &models.ToolReviewRecord{
	// 	ToolID:  uint(id),
	// 	AdminID: adminID,
	// 	Action:  "approve",
	// 	Reason:  "",
	// }
	// if err := h.reviewRepo.Create(ctx, record); err != nil {
	// 	httpx.ErrorCtx(ctx, w, err)
	// 	return
	// }

	httpx.OkJsonCtx(ctx, w, types.Success(map[string]interface{}{
		"success": true,
		"message": "审核通过",
	}))
}

// RejectTool 审核拒绝
func (h *AdminToolReviewHandler) RejectTool(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少工具ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("无效的工具ID"))
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	if req.Reason == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("拒绝原因不能为空"))
		return
	}

	// 获取当前管理员ID
	adminID := r.Header.Get("X-Admin-ID")
	if adminID == "" {
		adminID = "system"
	}

	// 更新工具状态为已拒绝
	result := database.GetDB().Table("tools").Where("id = ?", id).Updates(map[string]interface{}{
		"review_status":   "rejected",
		"revision_reason": req.Reason,
	})
	if result.Error != nil {
		httpx.ErrorCtx(ctx, w, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("工具不存在"))
		return
	}

	// 创建审核记录
	record := &models.ToolReviewRecord{
		ToolID:  uint(id),
		AdminID: adminID,
		Action:  "reject",
		Reason:  req.Reason,
	}
	if err := h.reviewRepo.Create(ctx, record); err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(map[string]interface{}{
		"success": true,
		"message": "已拒绝",
	}))
}

// RequestRevision 退回修改
func (h *AdminToolReviewHandler) RequestRevision(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少工具ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("无效的工具ID"))
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := httpx.ParseJsonBody(r, &req); err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	if req.Reason == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("退回原因不能为空"))
		return
	}

	// 获取当前管理员ID
	adminID := r.Header.Get("X-Admin-ID")
	if adminID == "" {
		adminID = "system"
	}

	// 更新工具状态为需修改
	result := database.GetDB().Table("tools").Where("id = ?", id).Updates(map[string]interface{}{
		"review_status":   "revision_requested",
		"revision_reason": req.Reason,
	})
	if result.Error != nil {
		httpx.ErrorCtx(ctx, w, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("工具不存在"))
		return
	}

	// 创建审核记录
	record := &models.ToolReviewRecord{
		ToolID:  uint(id),
		AdminID: adminID,
		Action:  "request_revision",
		Reason:  req.Reason,
	}
	if err := h.reviewRepo.Create(ctx, record); err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(map[string]interface{}{
		"success": true,
		"message": "已退回修改",
	}))
}

// SetToolOnline 上线工具
func (h *AdminToolReviewHandler) SetToolOnline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少工具ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("无效的工具ID"))
		return
	}

	result := database.GetDB().Table("tools").Where("id = ?", id).Updates(map[string]interface{}{
		"is_online": true,
	})
	if result.Error != nil {
		httpx.ErrorCtx(ctx, w, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("工具不存在"))
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(map[string]interface{}{
		"success": true,
		"message": "工具已上线",
	}))
}

// SetToolOffline 下线工具
func (h *AdminToolReviewHandler) SetToolOffline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少工具ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("无效的工具ID"))
		return
	}

	result := database.GetDB().Table("tools").Where("id = ?", id).Updates(map[string]interface{}{
		"is_online": false,
	})
	if result.Error != nil {
		httpx.ErrorCtx(ctx, w, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("工具不存在"))
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(map[string]interface{}{
		"success": true,
		"message": "工具已下线",
	}))
}