package service

import (
	"hot-ai-backend/internal/access"
	"hot-ai-backend/internal/models"
	"hot-ai-backend/internal/repository"
	"strconv"
	"time"
)

// ProfessionService 职业服务
type ProfessionService struct {
	professionRepo *repository.ProfessionRepository
}

// NewProfessionService 创建职业服务实例
func NewProfessionService(professionRepo *repository.ProfessionRepository) *ProfessionService {
	return &ProfessionService{
		professionRepo: professionRepo,
	}
}

// GetProfessionsRequest 获取职业列表请求参数
type GetProfessionsRequest struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	CategoryID uint   `json:"category_id"`
	RiskLevel  string `json:"risk_level"`
	Keyword    string `json:"keyword"`
}

// GetProfessionsResponse 获取职业列表响应
type GetProfessionsResponse struct {
	Professions []ProfessionView    `json:"professions"`
	Total       int64               `json:"total"`
	TotalPages  int                 `json:"total_pages"`
	Page        int                 `json:"page"`
	PageSize    int                 `json:"page_size"`
}

// GetProfessions 获取职业列表
func (s *ProfessionService) GetProfessions(page, pageSize int, categoryID uint, riskLevel, keyword string) (*GetProfessionsResponse, error) {
	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 12
	}

	professions, total, err := s.professionRepo.GetList(page, pageSize, categoryID, riskLevel, keyword)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	// 默认 level=0，handler 会用 ctx 里的 level 重新算并覆盖
	return &GetProfessionsResponse{
		Professions: ProfessionListView(professions, 0),
		Total:       total,
		TotalPages:  totalPages,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

// GetProfessionByID 根据 ID 获取职业详情
func (s *ProfessionService) GetProfessionByID(idStr string) (*models.Profession, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, err
	}
	profession, err := s.professionRepo.GetByID(uint(id))
	if err != nil {
		return nil, err
	}
	return profession, nil
}

// CreateProfessionRequest 创建职业请求
type CreateProfessionRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	CategoryID  uint   `json:"category_id"`
	RiskLevel   string `json:"risk_level"`
	RiskScore   int    `json:"risk_score"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status"`
}

// CreateProfession 创建职业
func (s *ProfessionService) CreateProfession(req *CreateProfessionRequest) (*models.Profession, error) {
	now := time.Now()
	profession := &models.Profession{
		Name:        req.Name,
		Slug:        req.Slug,
		Icon:        req.Icon,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		RiskLevel:   req.RiskLevel,
		RiskScore:   req.RiskScore,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		PublishedAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.professionRepo.Create(profession); err != nil {
		return nil, err
	}
	return profession, nil
}

// UpdateProfessionRequest 更新职业请求
type UpdateProfessionRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	CategoryID  uint   `json:"category_id"`
	RiskLevel   string `json:"risk_level"`
	RiskScore   int    `json:"risk_score"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status"`
}

// UpdateProfession 更新职业
func (s *ProfessionService) UpdateProfession(idStr string, req *UpdateProfessionRequest) (*models.Profession, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	// 先获取现有记录，保留有效的 PublishedAt
	existing, err := s.professionRepo.GetByID(uint(id))
	if err != nil {
		return nil, err
	}

	// 检查并修复无效的 PublishedAt（0000-00-00 或 year < 1000）
	validPublishedAt := existing.PublishedAt
	if validPublishedAt.IsZero() || validPublishedAt.Year() < 1000 {
		validPublishedAt = time.Now()
	}

	profession := &models.Profession{
		ID:          uint(id),
		Name:        req.Name,
		Slug:        req.Slug,
		Icon:        req.Icon,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		RiskLevel:   req.RiskLevel,
		RiskScore:   req.RiskScore,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		PublishedAt: validPublishedAt,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	if err := s.professionRepo.Update(profession); err != nil {
		return nil, err
	}
	return profession, nil
}

// DeleteProfession 删除职业
func (s *ProfessionService) DeleteProfession(idStr string) error {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return err
	}

	return s.professionRepo.Delete(uint(id))
}

// SearchProfessions 搜索职业
func (s *ProfessionService) SearchProfessions(keyword string, page, pageSize int) (*GetProfessionsResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 12
	}

	professions, total, err := s.professionRepo.SearchByKeyword(keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &GetProfessionsResponse{
		Professions: ProfessionListView(professions, 0),
		Total:       total,
		TotalPages:  totalPages,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

// GetCategories 获取职业分类列表
func (s *ProfessionService) GetCategories() ([]models.ProfessionCategory, error) {
	return s.professionRepo.GetCategories()
}

// GetRiskLevels 获取风险等级信息
func (s *ProfessionService) GetRiskLevels() []models.RiskLevelInfo {
	return s.professionRepo.GetRiskLevels()
}

// GetFeatured 获取精选职业
func (s *ProfessionService) GetFeatured(limit int) ([]models.Profession, error) {
	return s.professionRepo.GetFeatured(limit)
}

// GetProfessionCount 获取职业总数
func (s *ProfessionService) GetProfessionCount() (int64, error) {
	return s.professionRepo.GetProfessionCount()
}

// ProfessionFeaturedView 速览卡片视图（精简字段）
type ProfessionFeaturedView struct {
	ID                uint                 `json:"id"`
	Name              string               `json:"name"`
	CategoryID        uint                 `json:"category_id,omitempty"`
	RiskLevel         string               `json:"risk_level"`
	RiskScore         int                  `json:"risk_score"`
	AccessLevel       int                  `json:"access_level"`
	IsLocked          bool                 `json:"is_locked"`
	RequiredLevel     int                  `json:"required_level,omitempty"`
	RequiredLevelName string               `json:"required_level_name,omitempty"`
	ImpactSummary     string               `json:"impact_summary"`
	ImpactTimeline    models.JSONStringMap `json:"impact_timeline"`
}

// ToProfessionFeaturedView 转换 Profession 为速览卡片视图
func ToProfessionFeaturedView(p *models.Profession, userLevel int) *ProfessionFeaturedView {
	v := &ProfessionFeaturedView{
		ID:             p.ID,
		Name:           p.Name,
		CategoryID:     p.CategoryID,
		RiskLevel:      p.RiskLevel,
		RiskScore:      p.RiskScore,
		AccessLevel:    p.AccessLevel,
		ImpactTimeline: models.JSONStringMap{},
	}
	if p.ImpactAnalysis != nil {
		v.ImpactSummary = p.ImpactAnalysis.ImpactSummary
		v.ImpactTimeline = p.ImpactAnalysis.ImpactTimeline
	}
	decision := access.Decide(userLevel, p.AccessLevel)
	v.IsLocked = !decision.Allow
	if !decision.Allow {
		v.RequiredLevel = p.AccessLevel
		v.RequiredLevelName = access.LevelName(p.AccessLevel)
	}
	return v
}

// GetFeaturedAll 获取全部职业（速览用）
func (s *ProfessionService) GetFeaturedAll() ([]models.Profession, error) {
	return s.professionRepo.GetFeaturedAll()
}

// ProfessionView 职业详情响应 (含 access 决策)
type ProfessionView struct {
	*models.Profession
	IsLocked          bool                  `json:"is_locked"`
	RequiredLevel     int                   `json:"required_level,omitempty"`
	RequiredLevelName string                `json:"required_level_name,omitempty"`
	Locked            *access.LockedContent `json:"locked,omitempty"`
}

// ToProfessionView 把 Profession 包成 view，根据 userLevel 算 access
func ToProfessionView(p *models.Profession, userLevel int) *ProfessionView {
	v := &ProfessionView{Profession: p}
	decision := access.Decide(userLevel, p.AccessLevel)
	v.IsLocked = !decision.Allow
	if !decision.Allow {
		v.RequiredLevel = p.AccessLevel
		v.RequiredLevelName = access.LevelName(p.AccessLevel)
		preview, _ := access.TruncateContent(p.Description, access.GuestPreviewChars)
		p.Description = preview
		// 锁定时也清空子结构（影响分析、转型建议）
		p.ImpactAnalysis = nil
		p.TransitionAdvice = nil
		p.MarketData = nil
		lp := access.LockedPlaceholder("职业", p.AccessLevel)
		v.Locked = &lp
	}
	return v
}

// ProfessionListView 给列表里每条职业打 is_locked 标签
func ProfessionListView(professions []models.Profession, userLevel int) []ProfessionView {
	out := make([]ProfessionView, 0, len(professions))
	for i := range professions {
		v := ProfessionView{Profession: &professions[i]}
		decision := access.Decide(userLevel, professions[i].AccessLevel)
		v.IsLocked = !decision.Allow
		if !decision.Allow {
			v.RequiredLevel = professions[i].AccessLevel
			v.RequiredLevelName = access.LevelName(professions[i].AccessLevel)
		}
		out = append(out, v)
	}
	return out
}
