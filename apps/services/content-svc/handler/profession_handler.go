package handler

import (
	"fmt"
	"net/http"

	"hot-ai-backend/internal/service"

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
	pageSize := 10
	riskLevel := r.URL.Query().Get("risk_level")

	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	resp, err := h.professionService.GetProfessions(&service.GetProfessionsRequest{
		Page:      page,
		PageSize:  pageSize,
		RiskLevel: riskLevel,
	})
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, resp)
}

// GetProfessionBySlug 根据slug获取职业详情
func (h *ProfessionHandler) GetProfessionBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少职业slug"))
		return
	}

	profession, err := h.professionService.GetProfessionBySlug(slug)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, profession)
}

// SearchProfessions 搜索职业
func (h *ProfessionHandler) SearchProfessions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少搜索关键词"))
		return
	}

	page := 1
	pageSize := 10

	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	resp, err := h.professionService.SearchProfessions(query, page, pageSize)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, resp)
}

// GetRiskLevels 获取风险等级信息
func (h *ProfessionHandler) GetRiskLevels(w http.ResponseWriter, r *http.Request) {
	riskLevels, err := h.professionService.GetRiskLevels()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, riskLevels)
}
