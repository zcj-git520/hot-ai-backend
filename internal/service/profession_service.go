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
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	RiskLevel string `json:"risk_level"`
}

// GetProfessionsResponse 获取职业列表响应
type GetProfessionsResponse struct {
	Professions []models.Profession `json:"professions"`
	Total       int64               `json:"total"`
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
		req.PageSize = 10
	}

	professions, total, err := s.professionRepo.GetList(req.Page, req.PageSize, req.RiskLevel)
	if err != nil {
		return nil, err
	}

	return &GetProfessionsResponse{
		Professions: professions,
		Total:       total,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}, nil
}

// GetProfessionBySlug 根据slug获取职业详情
func (s *ProfessionService) GetProfessionBySlug(slug string) (*models.Profession, error) {
	profession, err := s.professionRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	// 增加访问量
	_ = s.professionRepo.IncrementViewCount(profession.ID)

	return profession, nil
}

// SearchProfessions 搜索职业
func (s *ProfessionService) SearchProfessions(query string, page, pageSize int) (*GetProfessionsResponse, error) {
	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	professions, total, err := s.professionRepo.Search(query, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &GetProfessionsResponse{
		Professions: professions,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}, nil
}

// GetRiskLevels 获取风险等级信息
func (s *ProfessionService) GetRiskLevels() ([]models.RiskLevelInfo, error) {
	return s.professionRepo.GetRiskLevels()
}
