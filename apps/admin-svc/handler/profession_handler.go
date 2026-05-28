package handler

import (
	"errors"
	"fmt"
	"net/http"

	"hot-ai-backend/internal/service"
	"hot-ai-backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

var ErrMissingID = errors.New("缺少ID参数")

// ProfessionHandler 职业管理处理器
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
	riskLevel := ""
	keyword := ""

	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}
	if cid := r.URL.Query().Get("category_id"); cid != "" {
		fmt.Sscanf(cid, "%d", &categoryID)
	}
	riskLevel = r.URL.Query().Get("risk_level")
	keyword = r.URL.Query().Get("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 12
	}

	resp, err := h.professionService.GetProfessions(page, pageSize, categoryID, riskLevel, keyword)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(resp))
}

// GetProfessionByID 获取职业详情
func (h *ProfessionHandler) GetProfessionByID(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, ErrMissingID)
		return
	}

	profession, err := h.professionService.GetProfessionByID(idStr)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(profession))
}

// CreateProfession 创建职业
func (h *ProfessionHandler) CreateProfession(w http.ResponseWriter, r *http.Request) {
	var req service.CreateProfessionRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	profession, err := h.professionService.CreateProfession(&req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(profession))
}

// UpdateProfession 更新职业
func (h *ProfessionHandler) UpdateProfession(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, ErrMissingID)
		return
	}

	var req service.UpdateProfessionRequest
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	profession, err := h.professionService.UpdateProfession(idStr, &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(profession))
}

// DeleteProfession 删除职业
func (h *ProfessionHandler) DeleteProfession(w http.ResponseWriter, r *http.Request) {
	idStr := getPathValue(r, "id")
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, ErrMissingID)
		return
	}

	err := h.professionService.DeleteProfession(idStr)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(nil))
}

// GetCategories 获取职业分类
func (h *ProfessionHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.professionService.GetCategories()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(categories))
}

// GetRiskLevels 获取风险等级
func (h *ProfessionHandler) GetRiskLevels(w http.ResponseWriter, r *http.Request) {
	riskLevels := h.professionService.GetRiskLevels()
	httpx.OkJsonCtx(r.Context(), w, types.Success(riskLevels))
}