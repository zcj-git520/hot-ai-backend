package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"hot-ai-backend/internal/service"
	"hot-ai-backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ProfessionHandler 职业处理器
type ProfessionHandler struct {
	professionService *service.ProfessionService
}

// NewProfessionHandler 创建职业处理器实例
func NewProfessionHandler(professionService *service.ProfessionService) *ProfessionHandler {
	return &ProfessionHandler{
		professionService: professionService,
	}
}

// GetProfessions 获取职业列表
func (h *ProfessionHandler) GetProfessions(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 12
	categoryID := uint(0)
	riskLevel := r.URL.Query().Get("risk_level")
	keyword := r.URL.Query().Get("keyword")

	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}
	if cid := r.URL.Query().Get("category_id"); cid != "" {
		if id, err := strconv.ParseUint(cid, 10, 32); err == nil {
			categoryID = uint(id)
		}
	}

	resp, err := h.professionService.GetProfessions(page, pageSize, categoryID, riskLevel, keyword)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(resp))
}

// GetProfessionByID 根据 ID 获取职业详情
func (h *ProfessionHandler) GetProfessionByID(w http.ResponseWriter, r *http.Request) {
	// 从URL路径参数获取ID
	idStr := r.PathValue("id")

	// 如果PathValue不可用(Go < 1.22),尝试从URL path中手动提取
	if idStr == "" {
		path := r.URL.Path
		// URL格式: /api/professions/123
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '/' {
				idStr = path[i+1:]
				break
			}
		}
	}

	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少职业ID"))
		return
	}

	_, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的职业ID"))
		return
	}

	profession, err := h.professionService.GetProfessionByID(idStr)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(profession))
}

// GetCategories 获取职业分类列表
func (h *ProfessionHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.professionService.GetCategories()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(categories))
}

// GetRiskLevels 获取风险等级信息
func (h *ProfessionHandler) GetRiskLevels(w http.ResponseWriter, r *http.Request) {
	riskLevels := h.professionService.GetRiskLevels()
	httpx.OkJsonCtx(r.Context(), w, types.Success(riskLevels))
}

// SearchProfessions 搜索职业
func (h *ProfessionHandler) SearchProfessions(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少搜索关键词"))
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

	resp, err := h.professionService.SearchProfessions(keyword, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(resp))
}
