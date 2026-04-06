package repository

import (
	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/models"
	"strings"
)

// ProfessionRepository 职业仓储
type ProfessionRepository struct{}

// NewProfessionRepository 创建职业仓储实例
func NewProfessionRepository() *ProfessionRepository {
	return &ProfessionRepository{}
}

// GetList 获取职业列表(支持分页和风险等级筛选)
func (r *ProfessionRepository) GetList(page, pageSize int, riskLevel string) ([]models.Profession, int64, error) {
	var professions []models.Profession
	var total int64

	query := database.GetDB().Model(&models.Profession{}).Where("is_active = ?", true)

	// 风险等级筛选
	if riskLevel != "" && riskLevel != "all" {
		query = query.Where("risk_level = ?", riskLevel)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("automation_rate DESC").Offset(offset).Limit(pageSize).Find(&professions).Error; err != nil {
		return nil, 0, err
	}

	return professions, total, nil
}

// GetBySlug 根据slug获取职业详情
func (r *ProfessionRepository) GetBySlug(slug string) (*models.Profession, error) {
	var profession models.Profession
	if err := database.GetDB().Where("slug = ? AND is_active = ?", slug, true).First(&profession).Error; err != nil {
		return nil, err
	}
	return &profession, nil
}

// Search 搜索职业
func (r *ProfessionRepository) Search(query string, page, pageSize int) ([]models.Profession, int64, error) {
	var professions []models.Profession
	var total int64

	searchQuery := "%" + strings.ToLower(query) + "%"
	db := database.GetDB().Model(&models.Profession{}).Where("is_active = ? AND (LOWER(name) LIKE ? OR LOWER(description) LIKE ?)", true, searchQuery, searchQuery)

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Order("automation_rate DESC").Offset(offset).Limit(pageSize).Find(&professions).Error; err != nil {
		return nil, 0, err
	}

	return professions, total, nil
}

// GetRiskLevels 获取所有风险等级信息
func (r *ProfessionRepository) GetRiskLevels() ([]models.RiskLevelInfo, error) {
	var riskLevels []models.RiskLevelInfo
	if err := database.GetDB().Order("sort_order ASC").Find(&riskLevels).Error; err != nil {
		return nil, err
	}
	return riskLevels, nil
}

// IncrementViewCount 增加职业页面访问量
func (r *ProfessionRepository) IncrementViewCount(id string) error {
	return database.GetDB().Model(&models.Profession{}).Where("id = ?", id).UpdateColumn("view_count", database.GetDB().Raw("view_count + 1")).Error
}
