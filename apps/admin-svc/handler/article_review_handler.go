package handler

import (
	"fmt"
	"net/http"

	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/models"
	"hot-ai-backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// AdminArticleReviewHandler 文章审核处理器
type AdminArticleReviewHandler struct{}

// NewAdminArticleReviewHandler 创建文章审核处理器实例
func NewAdminArticleReviewHandler() *AdminArticleReviewHandler {
	return &AdminArticleReviewHandler{}
}

// GetPendingArticles 获取待审核文章列表
func (h *AdminArticleReviewHandler) GetPendingArticles(w http.ResponseWriter, r *http.Request) {
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

	// 查询文章列表
	var articles []models.Article
	var total int64

	query := database.GetDB().Table("articles")

	// 按审核状态筛选
	switch reviewStatus {
	case "draft":
		query = query.Where("status = ?", 0)
	case "published":
		query = query.Where("status = ?", 1)
	case "rejected":
		query = query.Where("status = ?", 3) // 3 表示被拒绝
	default:
		// 默认显示待审核和已发布的
		if reviewStatus == "" {
			query = query.Where("status IN (0, 1, 3)")
		} else {
			query = query.Where("status = ?", reviewStatus)
		}
	}

	if search != "" {
		query = query.Where("title LIKE ? OR summary LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(map[string]interface{}{
		"list":      articles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

// GetArticleDetailForReview 获取文章详情
func (h *AdminArticleReviewHandler) GetArticleDetailForReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// 尝试从路径获取
		idStr = r.PathValue("id")
	}

	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少文章ID"))
		return
	}

	var article models.Article
	if err := database.GetDB().Where("id = ?", idStr).First(&article).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(article))
}

// ApproveArticle 审核通过
func (h *AdminArticleReviewHandler) ApproveArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少文章ID"))
		return
	}

	var article models.Article
	if err := database.GetDB().Where("id = ?", idStr).First(&article).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	// 更新状态为已发布
	if err := database.GetDB().Model(&article).Updates(map[string]interface{}{
		"status": 1,
	}).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(map[string]interface{}{
		"message": "审核通过",
		"success": true,
	}))
}

// RejectArticle 拒绝文章
func (h *AdminArticleReviewHandler) RejectArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少文章ID"))
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	var article models.Article
	if err := database.GetDB().Where("id = ?", idStr).First(&article).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	// 更新状态为已拒绝
	if err := database.GetDB().Model(&article).Updates(map[string]interface{}{
		"status":           3, // 3 = 已拒绝
		"rejection_reason": req.Reason,
	}).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(map[string]interface{}{
		"message": "已拒绝",
		"success": true,
	}))
}

// PublishArticle 直接发布文章
func (h *AdminArticleReviewHandler) PublishArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少文章ID"))
		return
	}

	var article models.Article
	if err := database.GetDB().Where("id = ?", idStr).First(&article).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	// 更新状态为已发布
	if err := database.GetDB().Model(&article).Updates(map[string]interface{}{
		"status":        1,
		"published_at":  database.GetDB().NowFunc(),
	}).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(map[string]interface{}{
		"message": "已发布",
		"success": true,
	}))
}

// UnpublishArticle 下架文章
func (h *AdminArticleReviewHandler) UnpublishArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少文章ID"))
		return
	}

	var article models.Article
	if err := database.GetDB().Where("id = ?", idStr).First(&article).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	// 更新状态为待审核
	if err := database.GetDB().Model(&article).Updates(map[string]interface{}{
		"status": 0,
	}).Error; err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(map[string]interface{}{
		"message": "已下架",
		"success": true,
	}))
}
