package logic

import (
	"fmt"
	"strings"

	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/models"
)

// LearningPathLogic 学习路径管理逻辑
type LearningPathLogic struct{}

// NewLearningPathLogic 创建学习路径管理逻辑实例
func NewLearningPathLogic() *LearningPathLogic {
	return &LearningPathLogic{}
}

// CreateLearningPath 创建学习路径
func (l *LearningPathLogic) CreateLearningPath(req *CreateLearningPathRequest) (*models.LearningPath, error) {
	// 生成 slug
	slug := generateSlug(req.Title)

	path := &models.LearningPath{
		Title:          req.Title,
		Slug:           slug,
		Icon:           req.Icon,
		Description:    req.Description,
		Difficulty:     req.Difficulty,
		LevelLabel:     req.LevelLabel,
		LearningGoals:  toJSONArray(req.LearningGoals),
		TargetAudience: toJSONArray(req.TargetAudience),
		EstimatedDays:  req.EstimatedDays,
		EstimatedHours: req.EstimatedHours,
		CoverImage:     req.CoverImage,
		IsFeatured:     boolToInt(req.IsFeatured),
		IsActive:       1,
		Status:         0, // 待审核状态
	}

	if err := database.GetDB().Create(path).Error; err != nil {
		return nil, err
	}

	// 创建管理数据
	management := &models.LearningPathManagement{
		PathID: path.ID,
	}
	database.GetDB().Create(management)

	return path, nil
}

// UpdateLearningPath 更新学习路径
func (l *LearningPathLogic) UpdateLearningPath(id uint, req *UpdateLearningPathRequest) (*models.LearningPath, error) {
	path, err := l.getPathByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != "" {
		path.Title = req.Title
		path.Slug = generateSlug(req.Title)
	}
	if req.Icon != "" {
		path.Icon = req.Icon
	}
	if req.Description != "" {
		path.Description = req.Description
	}
	if req.Difficulty != "" {
		path.Difficulty = req.Difficulty
	}
	if req.LevelLabel != "" {
		path.LevelLabel = req.LevelLabel
	}
	if req.LearningGoals != nil {
		path.LearningGoals = toJSONArray(req.LearningGoals)
	}
	if req.TargetAudience != nil {
		path.TargetAudience = toJSONArray(req.TargetAudience)
	}
	if req.EstimatedDays > 0 {
		path.EstimatedDays = req.EstimatedDays
	}
	if req.EstimatedHours > 0 {
		path.EstimatedHours = req.EstimatedHours
	}
	if req.CoverImage != "" {
		path.CoverImage = req.CoverImage
	}
	if req.IsFeatured != nil {
		path.IsFeatured = boolToInt(*req.IsFeatured)
	}
	if req.SortOrder != nil {
		path.SortOrder = *req.SortOrder
	}

	if err := database.GetDB().Save(path).Error; err != nil {
		return nil, err
	}

	return path, nil
}

// DeleteLearningPath 删除学习路径
func (l *LearningPathLogic) DeleteLearningPath(id uint) error {
	// 软删除：更新状态
	result := database.GetDB().Model(&models.LearningPath{}).Where("id = ?", id).Update("status", 2)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("学习路径不存在")
	}
	return nil
}

// PublishLearningPath 发布学习路径
func (l *LearningPathLogic) PublishLearningPath(id uint) error {
	result := database.GetDB().Model(&models.LearningPath{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        1,
		"is_active":     1,
		"published_at":  database.GetDB().NowFunc(),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("学习路径不存在")
	}
	return nil
}

// UnpublishLearningPath 下架学习路径
func (l *LearningPathLogic) UnpublishLearningPath(id uint) error {
	result := database.GetDB().Model(&models.LearningPath{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    0,
		"is_active": 0,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("学习路径不存在")
	}
	return nil
}

// SetFeatured 设置推荐状态
func (l *LearningPathLogic) SetFeatured(id uint, featured bool) error {
	result := database.GetDB().Model(&models.LearningPath{}).Where("id = ?", id).Update("is_featured", boolToInt(featured))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("学习路径不存在")
	}
	return nil
}

// CreateChapter 创建章节
func (l *LearningPathLogic) CreateChapter(pathID uint, req *CreateChapterRequest) (*models.PathChapter, error) {
	// 检查路径是否存在
	path, err := l.getPathByID(pathID)
	if err != nil {
		return nil, fmt.Errorf("学习路径不存在")
	}

	slug := generateSlug(req.Title)

	// 获取最大排序索引
	var maxOrder int
	database.GetDB().Model(&models.PathChapter{}).Where("path_id = ?", pathID).Select("COALESCE(MAX(order_index), 0)").Scan(&maxOrder)

	chapter := &models.PathChapter{
		PathID:        path.ID,
		Title:         req.Title,
		Slug:          slug,
		Description:   req.Description,
		ContentType:   req.ContentType,
		Content:       req.Content,
		VideoURL:      req.VideoURL,
		ExternalLinks: toJSONArray(req.ExternalLinks),
		EstimatedHours: req.EstimatedHours,
		OrderIndex:    maxOrder + 1,
		IsFree:       boolToInt(req.IsFree),
		Status:        1,
	}

	if err := database.GetDB().Create(chapter).Error; err != nil {
		return nil, err
	}

	// 更新路径的章节数量
	l.updateChapterCount(pathID)

	return chapter, nil
}

// UpdateChapter 更新章节
func (l *LearningPathLogic) UpdateChapter(id uint, req *UpdateChapterRequest) (*models.PathChapter, error) {
	chapter, err := l.getChapterByID(id)
	if err != nil {
		return nil, fmt.Errorf("章节不存在")
	}

	// 更新字段
	if req.Title != "" {
		chapter.Title = req.Title
		chapter.Slug = generateSlug(req.Title)
	}
	if req.Description != "" {
		chapter.Description = req.Description
	}
	if req.ContentType != "" {
		chapter.ContentType = req.ContentType
	}
	if req.Content != "" {
		chapter.Content = req.Content
	}
	if req.VideoURL != "" {
		chapter.VideoURL = req.VideoURL
	}
	if req.ExternalLinks != nil {
		chapter.ExternalLinks = toJSONArray(req.ExternalLinks)
	}
	if req.EstimatedHours > 0 {
		chapter.EstimatedHours = req.EstimatedHours
	}
	if req.OrderIndex != nil {
		chapter.OrderIndex = *req.OrderIndex
	}
	if req.IsFree != nil {
		chapter.IsFree = boolToInt(*req.IsFree)
	}

	if err := database.GetDB().Save(chapter).Error; err != nil {
		return nil, err
	}

	return chapter, nil
}

// DeleteChapter 删除章节
func (l *LearningPathLogic) DeleteChapter(id uint) error {
	chapter, err := l.getChapterByID(id)
	if err != nil {
		return fmt.Errorf("章节不存在")
	}

	pathID := chapter.PathID

	// 软删除
	result := database.GetDB().Model(&models.PathChapter{}).Where("id = ?", id).Update("status", 2)
	if result.Error != nil {
		return result.Error
	}

	// 更新路径的章节数量
	l.updateChapterCount(pathID)

	return nil
}

// GetLearningPaths 获取学习路径列表（管理后台，不限制 status）
func (l *LearningPathLogic) GetLearningPaths(req *GetLearningPathsRequest) ([]models.LearningPath, int64, error) {
	var paths []models.LearningPath
	var total int64

	query := database.GetDB().Model(&models.LearningPath{}).Where("status != ?", 2) // 排除已删除

	if req.Difficulty != "" {
		query = query.Where("difficulty = ?", req.Difficulty)
	}
	if req.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&paths).Error; err != nil {
		return nil, 0, err
	}

	// 填充章节数量
	for i := range paths {
		l.fillChapterCount(&paths[i])
	}

	return paths, total, nil
}

// GetLearningPathByID 获取学习路径详情
func (l *LearningPathLogic) GetLearningPathByID(id uint) (*models.LearningPath, error) {
	var path models.LearningPath
	if err := database.GetDB().Where("id = ? AND status != ?", id, 2).First(&path).Error; err != nil {
		return nil, err
	}

	l.fillChapters(&path)
	l.fillManagementData(&path)

	return &path, nil
}

// GetChapters 获取章节列表
func (l *LearningPathLogic) GetChapters(pathID uint) ([]models.PathChapter, error) {
	var chapters []models.PathChapter
	if err := database.GetDB().Where("path_id = ? AND status != ?", pathID, 2).Order("order_index ASC").Find(&chapters).Error; err != nil {
		return nil, err
	}
	return chapters, nil
}

// GetChapterByID 获取章节详情
func (l *LearningPathLogic) GetChapterByID(id uint) (*models.PathChapter, error) {
	return l.getChapterByID(id)
}

// ========== 私有方法 ==========

func (l *LearningPathLogic) getPathByID(id uint) (*models.LearningPath, error) {
	var path models.LearningPath
	if err := database.GetDB().Where("id = ? AND status != ?", id, 2).First(&path).Error; err != nil {
		return nil, err
	}
	return &path, nil
}

func (l *LearningPathLogic) getChapterByID(id uint) (*models.PathChapter, error) {
	var chapter models.PathChapter
	if err := database.GetDB().Where("id = ? AND status != ?", id, 2).First(&chapter).Error; err != nil {
		return nil, err
	}
	return &chapter, nil
}

func (l *LearningPathLogic) fillChapterCount(path *models.LearningPath) {
	var count int64
	database.GetDB().Model(&models.PathChapter{}).Where("path_id = ? AND status = ?", path.ID, 1).Count(&count)
	path.ChapterCount = int(count)
}

func (l *LearningPathLogic) fillChapters(path *models.LearningPath) {
	var chapters []models.PathChapter
	if err := database.GetDB().Where("path_id = ? AND status = ?", path.ID, 1).Order("order_index ASC").Find(&chapters).Error; err != nil {
		path.Chapters = nil
		return
	}
	path.Chapters = chapters
}

func (l *LearningPathLogic) fillManagementData(path *models.LearningPath) {
	var management models.LearningPathManagement
	if err := database.GetDB().Where("path_id = ?", path.ID).First(&management).Error; err == nil {
		path.ManagementData = &management
	}
}

func (l *LearningPathLogic) updateChapterCount(pathID uint) {
	var count int64
	database.GetDB().Model(&models.PathChapter{}).Where("path_id = ? AND status = ?", pathID, 1).Count(&count)
	database.GetDB().Model(&models.LearningPath{}).Where("id = ?", pathID).Update("chapter_count", count)
}

// ========== 辅助函数 ==========

func generateSlug(title string) string {
	// 简单slug生成：中文转拼音或转小写英文
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// 移除非字母数字和连字符的字符
	var result []rune
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result = append(result, r)
		}
	}
	return string(result)
}

func toJSONArray(arr []string) string {
	if arr == nil || len(arr) == 0 {
		return "[]"
	}
	result := "["
	for i, s := range arr {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("\"%s\"", s)
	}
	result += "]"
	return result
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ========== 请求/响应结构 ==========

type GetLearningPathsRequest struct {
	Page       int
	PageSize   int
	Difficulty string
	Search     string
	Status     *int
}

type CreateLearningPathRequest struct {
	Title          string
	Icon           string
	Description    string
	Difficulty     string
	LevelLabel     string
	LearningGoals  []string
	TargetAudience []string
	EstimatedDays  int
	EstimatedHours int
	CoverImage     string
	IsFeatured     bool
}

type UpdateLearningPathRequest struct {
	Title          *string
	Icon           *string
	Description    *string
	Difficulty     *string
	LevelLabel     *string
	LearningGoals  []string
	TargetAudience []string
	EstimatedDays  *int
	EstimatedHours *int
	CoverImage     *string
	IsFeatured     *bool
	SortOrder      *int
}

type CreateChapterRequest struct {
	Title         string
	Description   string
	ContentType   string
	Content       string
	VideoURL      string
	ExternalLinks []string
	EstimatedHours int
	IsFree        bool
}

type UpdateChapterRequest struct {
	Title         *string
	Description   *string
	ContentType   *string
	Content       *string
	VideoURL      *string
	ExternalLinks []string
	EstimatedHours *int
	OrderIndex    *int
	IsFree        *bool
}