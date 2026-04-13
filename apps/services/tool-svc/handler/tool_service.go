package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
	"hot-ai-backend/internal/service"
	"hot-ai-backend/internal/types"
)

type ToolServiceHandler struct {
	svc *service.ToolService
}

func NewToolServiceHandler(svc *service.ToolService) *ToolServiceHandler {
	return &ToolServiceHandler{svc: svc}
}

// ToolCategories 工具类别列表
func (h *ToolServiceHandler) ToolCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.svc.CategoryList(ctx)
	if err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(categories))
}

// ToolList 工具列表
func (h *ToolServiceHandler) ToolList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := make(map[string]interface{})

	// 解析查询参数
	query := r.URL.Query()
	
	// 页码和每页数量
	if pageStr := query.Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params["page"] = page
		}
	}
	if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			params["page_size"] = pageSize
		}
	}

	// 类别ID
	if categoryIDStr := query.Get("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil && categoryID > 0 {
			params["category_id"] = categoryID
		}
	}

	// 是否免费
	if isFreeStr := query.Get("is_free"); isFreeStr != "" {
		isFree := isFreeStr == "true"
		params["is_free"] = isFree
	}

	// 难度等级
	if difficulty := query.Get("difficulty"); difficulty != "" {
		params["difficulty"] = difficulty
	}

	// 最低评分
	if minRatingStr := query.Get("min_rating"); minRatingStr != "" {
		if minRating, err := strconv.ParseFloat(minRatingStr, 64); err == nil && minRating > 0 {
			params["min_rating"] = minRating
		}
	}

	// 搜索关键词
	if search := query.Get("search"); search != "" {
		params["search"] = search
	}

	// 排序字段
	if sortBy := query.Get("sort_by"); sortBy != "" {
		params["sort_by"] = sortBy
	} else {
		params["sort_by"] = "popularity"
	}

	// 排序方向
	if order := query.Get("order"); order != "" {
		params["order"] = order
	} else {
		params["order"] = "desc"
	}

	tools, total, err := h.svc.ToolList(ctx, params)
	if err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	// 返回数据
	result := map[string]interface{}{
		"list":      tools,
		"total":     total,
		"page":      params["page"],
		"page_size": params["page_size"],
	}

	httpx.OkJsonCtx(ctx, w, types.Success(result))
}

// ToolDetail 工具详情（通过slug）
func (h *ToolServiceHandler) ToolDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 从URL路径参数获取slug
	slug := r.PathValue("slug")

	// 如果PathValue不可用，尝试从URL path中手动提取
	if slug == "" {
		path := r.URL.Path
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '/' {
				slug = path[i+1:]
				break
			}
		}
	}

	if slug == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少工具slug"))
		return
	}

	tool, err := h.svc.ToolDetail(ctx, slug)
	if err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(tool))
}

// ToolDetailByID 工具详情（通过id）
func (h *ToolServiceHandler) ToolDetailByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 从URL路径参数获取ID
	idStr := r.PathValue("id")

	// 如果PathValue不可用，尝试从URL path中手动提取
	if idStr == "" {
		path := r.URL.Path
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '/' {
				idStr = path[i+1:]
				break
			}
		}
	}

	if idStr == "" {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("缺少工具ID"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		httpx.ErrorCtx(ctx, w, fmt.Errorf("无效的工具ID"))
		return
	}

	tool, err := h.svc.ToolDetailByID(ctx, uint(id))
	if err != nil {
		httpx.ErrorCtx(ctx, w, err)
		return
	}

	httpx.OkJsonCtx(ctx, w, types.Success(tool))
}
