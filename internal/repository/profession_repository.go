package repository

import (
	"strings"

	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/models"
)

// ProfessionRepository 职业仓储
type ProfessionRepository struct{}

// NewProfessionRepository 创建职业仓储实例
func NewProfessionRepository() *ProfessionRepository {
	return &ProfessionRepository{}
}

// GetList 获取职业列表（支持分页、分类、风险等级筛选、关键词搜索）
func (r *ProfessionRepository) GetList(page, pageSize int, categoryID uint, riskLevel, keyword string) ([]models.Profession, int64, error) {
	var professions []models.Profession
	var total int64

	query := database.GetDB().Model(&models.Profession{}).Where("professions.status = ?", 1)

	// 分类筛选
	if categoryID > 0 {
		query = query.Where("professions.category_id = ?", categoryID)
	}

	// 风险等级筛选
	if riskLevel != "" && riskLevel != "all" {
		query = query.Where("professions.risk_level = ?", riskLevel)
	}

	// 关键词搜索
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("professions.name LIKE ? OR professions.description LIKE ?", searchPattern, searchPattern)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("professions.sort_order ASC, professions.risk_score DESC").Offset(offset).Limit(pageSize).Find(&professions).Error; err != nil {
		return nil, 0, err
	}

	// 填充关联数据
	for i := range professions {
		r.fillProfessionData(&professions[i])
	}

	return professions, total, nil
}

// GetByID 根据 ID 获取职业详情
func (r *ProfessionRepository) GetByID(id uint) (*models.Profession, error) {
	var profession models.Profession
	if err := database.GetDB().Where("id = ? AND status = ?", id, 1).First(&profession).Error; err != nil {
		return nil, err
	}

	// 填充关联数据
	r.fillProfessionData(&profession)

	// 增加访问量
	_ = r.IncrementViewCount(id)

	return &profession, nil
}

// GetCategories 获取所有启用的职业分类
func (r *ProfessionRepository) GetCategories() ([]models.ProfessionCategory, error) {
	var categories []models.ProfessionCategory
	if err := database.GetDB().Where("status = ?", 1).Order("sort_order ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetImpactAnalysis 获取职业影响分析
func (r *ProfessionRepository) GetImpactAnalysis(professionID uint) (*models.ProfessionImpactAnalysis, error) {
	var analysis models.ProfessionImpactAnalysis
	if err := database.GetDB().Where("profession_id = ?", professionID).First(&analysis).Error; err != nil {
		return nil, err
	}
	return &analysis, nil
}

// GetTransitionAdvice 获取职业转型建议
func (r *ProfessionRepository) GetTransitionAdvice(professionID uint) (*models.ProfessionTransitionAdvice, error) {
	var advice models.ProfessionTransitionAdvice
	if err := database.GetDB().Where("profession_id = ?", professionID).First(&advice).Error; err != nil {
		return nil, err
	}
	return &advice, nil
}

// GetMarketData 获取职业市场数据
func (r *ProfessionRepository) GetMarketData(professionID uint) (*models.ProfessionMarketData, error) {
	var data models.ProfessionMarketData
	if err := database.GetDB().Where("profession_id = ?", professionID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

// IncrementViewCount 增加职业页面访问量
func (r *ProfessionRepository) IncrementViewCount(id uint) error {
	return database.GetDB().Model(&models.Profession{}).
		Where("id = ?", id).
		UpdateColumn("view_count", database.GetDB().Raw("view_count + 1")).Error
}

// GetRiskLevels 获取风险等级展示信息（静态数据）
func (r *ProfessionRepository) GetRiskLevels() []models.RiskLevelInfo {
	return []models.RiskLevelInfo{
		{ID: "extreme", Level: "extreme", Name: "极高风险", Icon: "🔴", Description: "高度重复、规则明确的工作内容", Color: "#ef4444", MinScore: 80, MaxScore: 100},
		{ID: "high", Level: "high", Name: "高风险", Icon: "🟠", Description: "大部分工作内容可被 AI 替代", Color: "#f97316", MinScore: 60, MaxScore: 79},
		{ID: "medium", Level: "medium", Name: "中风险", Icon: "🟡", Description: "部分工作内容可被 AI 辅助", Color: "#eab308", MinScore: 40, MaxScore: 59},
		{ID: "low", Level: "low", Name: "低风险", Icon: "🟢", Description: "需要创造力和复杂决策", Color: "#22c55e", MinScore: 0, MaxScore: 39},
	}
}

// fillProfessionData 填充职业的关联数据 (JSON 字段由自定义类型自动解析)
func (r *ProfessionRepository) fillProfessionData(profession *models.Profession) {
	db := database.GetDB()

	// 获取分类名称
	if profession.CategoryID > 0 {
		var category models.ProfessionCategory
		if err := db.Where("id = ?", profession.CategoryID).First(&category).Error; err == nil {
			profession.CategoryName = category.Name
		}
	}

	// 获取影响分析 (JSON 字段由自定义类型自动解析)
	analysis, err := r.GetImpactAnalysis(profession.ID)
	if err == nil {
		profession.ImpactAnalysis = analysis
	}

	// 获取转型建议 (JSON 字段由自定义类型自动解析)
	advice, err := r.GetTransitionAdvice(profession.ID)
	if err == nil {
		profession.TransitionAdvice = advice
	}

	// 获取市场数据
	market, err := r.GetMarketData(profession.ID)
	if err == nil {
		profession.MarketData = market
	}
}

// Search 搜索职业
func (r *ProfessionRepository) Search(keyword string, page, pageSize int) ([]models.Profession, int64, error) {
	return r.GetList(page, pageSize, 0, "", keyword)
}

// GetFeatured 获取精选职业
func (r *ProfessionRepository) GetFeatured(limit int) ([]models.Profession, error) {
	var professions []models.Profession
	query := database.GetDB().Model(&models.Profession{}).
		Where("professions.status = ?", 1).
		Order("professions.sort_order ASC, professions.view_count DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&professions).Error; err != nil {
		return nil, err
	}

	for i := range professions {
		r.fillProfessionData(&professions[i])
	}

	return professions, nil
}

// SearchByKeyword 根据关键词搜索职业（支持 FULLTEXT）
func (r *ProfessionRepository) SearchByKeyword(keyword string, page, pageSize int) ([]models.Profession, int64, error) {
	var professions []models.Profession
	var total int64

	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return r.GetList(page, pageSize, 0, "", "")
	}

	db := database.GetDB().Model(&models.Profession{}).Where("professions.status = ?", 1)

	// 尝试 FULLTEXT 搜索
	searchQuery := "MATCH(name, description) AGAINST(? IN BOOLEAN MODE)"
	var countResult int64
	if err := db.Where(searchQuery, keyword).Count(&countResult).Error; err == nil && countResult > 0 {
		db = db.Where(searchQuery, keyword)
	} else {
		// 降级为 LIKE 搜索
		searchPattern := "%" + keyword + "%"
		db = db.Where("professions.name LIKE ? OR professions.description LIKE ?", searchPattern, searchPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("professions.risk_score DESC").Offset(offset).Limit(pageSize).Find(&professions).Error; err != nil {
		return nil, 0, err
	}

	for i := range professions {
		r.fillProfessionData(&professions[i])
	}

	return professions, total, nil
}

// GetProfessionCount 获取职业总数
func (r *ProfessionRepository) GetProfessionCount() (int64, error) {
	var total int64
	if err := database.GetDB().Model(&models.Profession{}).Where("status = ?", 1).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// Create 创建职业
func (r *ProfessionRepository) Create(profession *models.Profession) error {
	return database.GetDB().Create(profession).Error
}

// Update 更新职业
func (r *ProfessionRepository) Update(profession *models.Profession) error {
	return database.GetDB().Save(profession).Error
}

// Delete 删除职业（软删除）
func (r *ProfessionRepository) Delete(id uint) error {
	return database.GetDB().Model(&models.Profession{}).Where("id = ?", id).Update("status", 0).Error
}
