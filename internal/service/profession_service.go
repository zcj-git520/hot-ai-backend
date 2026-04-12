package service

import (
	"hot-ai-backend/internal/models"
	"hot-ai-backend/internal/repository"
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
	Professions []models.Profession `json:"professions"`
	Total       int64               `json:"total"`
	TotalPages  int                 `json:"total_pages"`
	Page        int                 `json:"page"`
	PageSize    int                 `json:"page_size"`
}

// GetProfessions 获取职业列表
func (s *ProfessionService) GetProfessions(req *GetProfessionsRequest) (*GetProfessionsResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 12
	}

	professions, total, err := s.professionRepo.GetList(req.Page, req.PageSize, req.CategoryID, req.RiskLevel, req.Keyword)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &GetProfessionsResponse{
		Professions: professions,
		Total:       total,
		TotalPages:  totalPages,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}, nil
}

// GetProfessionByID 根据 ID 获取职业详情
func (s *ProfessionService) GetProfessionByID(id uint) (*models.Profession, error) {
	profession, err := s.professionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return profession, nil
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
		Professions: professions,
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
