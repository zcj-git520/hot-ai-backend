package repository

import (
	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/models"
)

// LearningPathRepository 学习路径仓储
type LearningPathRepository struct{}

// NewLearningPathRepository 创建学习路径仓储实例
func NewLearningPathRepository() *LearningPathRepository {
	return &LearningPathRepository{}
}

// GetList 获取学习路径列表（支持分页、难度筛选）
func (r *LearningPathRepository) GetList(page, pageSize int, difficulty string) ([]models.LearningPath, int64, error) {
	var paths []models.LearningPath
	var total int64

	query := database.GetDB().Model(&models.LearningPath{}).Where("status = ? AND is_active = ?", 1, 1)

	// 难度筛选
	if difficulty != "" && difficulty != "all" {
		query = query.Where("difficulty = ?", difficulty)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(pageSize).Find(&paths).Error; err != nil {
		return nil, 0, err
	}

	// 填充章节数量
	for i := range paths {
		r.fillChapterCount(&paths[i])
	}

	return paths, total, nil
}

// GetByID 根据 ID 获取学习路径详情
func (r *LearningPathRepository) GetByID(id uint) (*models.LearningPath, error) {
	var path models.LearningPath
	if err := database.GetDB().Where("id = ? AND status = ? AND is_active = ?", id, 1, 1).First(&path).Error; err != nil {
		return nil, err
	}

	// 填充章节
	r.fillChapters(&path)

	// 增加访问量
	_ = r.IncrementViewCount(id)

	return &path, nil
}

// GetBySlug 根据 slug 获取学习路径详情
func (r *LearningPathRepository) GetBySlug(slug string) (*models.LearningPath, error) {
	var path models.LearningPath
	if err := database.GetDB().Where("slug = ? AND status = ? AND is_active = ?", slug, 1, 1).First(&path).Error; err != nil {
		return nil, err
	}

	// 填充章节
	r.fillChapters(&path)

	// 增加访问量
	_ = r.IncrementViewCount(path.ID)

	return &path, nil
}

// GetChapters 获取路径的所有章节
func (r *LearningPathRepository) GetChapters(pathID uint) ([]models.PathChapter, error) {
	var chapters []models.PathChapter
	if err := database.GetDB().
		Where("path_id = ? AND status = ?", pathID, 1).
		Order("order_index ASC").
		Find(&chapters).Error; err != nil {
		return nil, err
	}
	return chapters, nil
}

// GetChapterByID 根据章节 ID 获取详情
func (r *LearningPathRepository) GetChapterByID(chapterID uint) (*models.PathChapter, error) {
	var chapter models.PathChapter
	if err := database.GetDB().Where("id = ? AND status = ?", chapterID, 1).First(&chapter).Error; err != nil {
		return nil, err
	}
	return &chapter, nil
}

// GetChapterBySlug 根据路径 ID 和章节 slug 获取详情
func (r *LearningPathRepository) GetChapterBySlug(pathID uint, slug string) (*models.PathChapter, error) {
	var chapter models.PathChapter
	if err := database.GetDB().
		Where("path_id = ? AND slug = ? AND status = ?", pathID, slug, 1).
		First(&chapter).Error; err != nil {
		return nil, err
	}
	return &chapter, nil
}

// GetPrevNextChapter 获取前一章和下一章
func (r *LearningPathRepository) GetPrevNextChapter(pathID, currentChapterID uint) (*models.PathChapter, *models.PathChapter, error) {
	var prevChapter, nextChapter *models.PathChapter

	// 获取前一章
	var prev models.PathChapter
	if err := database.GetDB().
		Where("path_id = ? AND id < ? AND status = ?", pathID, currentChapterID, 1).
		Order("order_index DESC").
		First(&prev).Error; err == nil {
		prevChapter = &prev
	}

	// 获取下一章
	var next models.PathChapter
	if err := database.GetDB().
		Where("path_id = ? AND id > ? AND status = ?", pathID, currentChapterID, 1).
		Order("order_index ASC").
		First(&next).Error; err == nil {
		nextChapter = &next
	}

	return prevChapter, nextChapter, nil
}

// GetLevelInfo 获取难度等级信息（静态数据）
func (r *LearningPathRepository) GetLevelInfo() []models.LevelInfo {
	return []models.LevelInfo{
		{ID: "beginner", Level: "beginner", Name: "入门", Icon: "🌱", Description: "零基础友好，讲解基础概念和工具使用", Color: "#4ade80", MinHours: 20, MaxHours: 40},
		{ID: "intermediate", Level: "intermediate", Name: "进阶", Icon: "✍️", Description: "需要一定基础，讲解实用技巧和进阶应用", Color: "#60a5fa", MinHours: 30, MaxHours: 60},
		{ID: "advanced", Level: "advanced", Name: "高级", Icon: "🚀", Description: "需要扎实基础，讲解深度应用和开发技能", Color: "#a855f7", MinHours: 40, MaxHours: 100},
	}
}

// GetFeatured 获取推荐路径
func (r *LearningPathRepository) GetFeatured(limit int) ([]models.LearningPath, error) {
	var paths []models.LearningPath
	query := database.GetDB().Model(&models.LearningPath{}).
		Where("status = ? AND is_active = ? AND is_featured = ?", 1, 1, 1).
		Order("sort_order ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&paths).Error; err != nil {
		return nil, err
	}

	for i := range paths {
		r.fillChapterCount(&paths[i])
	}

	return paths, nil
}

// GetPathProgress 获取用户的学习进度
func (r *LearningPathRepository) GetPathProgress(userID string, pathID uint) (*models.LearningProgress, error) {
	var progress models.LearningProgress
	query := database.GetDB().Model(&models.LearningProgress{}).Where("path_id = ?", pathID)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	} else {
		// TODO: 未登录时使用 session_id
	}

	if err := query.First(&progress).Error; err != nil {
		return nil, err
	}
	return &progress, nil
}

// SaveProgress 保存学习进度
func (r *LearningPathRepository) SaveProgress(progress *models.LearningProgress) error {
	db := database.GetDB()

	// 检查是否已存在
	var existing models.LearningProgress
	existQuery := db.Model(&models.LearningProgress{}).
		Where("path_id = ? AND chapter_id = ?", progress.PathID, progress.ChapterID)

	if progress.UserID != "" {
		existQuery = existQuery.Where("user_id = ?", progress.UserID)
	} else if progress.SessionID != "" {
		existQuery = existQuery.Where("session_id = ?", progress.SessionID)
	}

	if err := existQuery.First(&existing).Error; err == nil {
		// 更新
		progress.ID = existing.ID
		return db.Save(progress).Error
	}

	// 新建
	return db.Create(progress).Error
}

// GetCompletedChapters 获取已完成的章节列表
func (r *LearningPathRepository) GetCompletedChapters(userID string, pathID uint) ([]uint, error) {
	var progressList []models.LearningProgress
	query := database.GetDB().Model(&models.LearningProgress{}).
		Where("path_id = ? AND status = ?", pathID, "completed")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Pluck("chapter_id", &progressList).Error; err != nil {
		return nil, err
	}

	chapterIDs := make([]uint, len(progressList))
	for i, p := range progressList {
		chapterIDs[i] = p.ChapterID
	}
	return chapterIDs, nil
}

// IncrementViewCount 增加路径访问量
func (r *LearningPathRepository) IncrementViewCount(pathID uint) error {
	return database.GetDB().Model(&models.LearningPath{}).
		Where("id = ?", pathID).
		UpdateColumn("student_count", database.GetDB().Raw("student_count + 1")).Error
}

// IncrementStartCount 增加开始学习次数
func (r *LearningPathRepository) IncrementStartCount(pathID uint) error {
	return database.GetDB().Model(&models.LearningPathManagement{}).
		Where("path_id = ?", pathID).
		UpdateColumn("start_count", database.GetDB().Raw("start_count + 1")).Error
}

// IncrementCompleteCount 增加完成次数
func (r *LearningPathRepository) IncrementCompleteCount(pathID uint) error {
	return database.GetDB().Model(&models.LearningPathManagement{}).
		Where("path_id = ?", pathID).
		UpdateColumn("complete_count", database.GetDB().Raw("complete_count + 1")).Error
}

// fillChapterCount 填充路径的章节数量
func (r *LearningPathRepository) fillChapterCount(path *models.LearningPath) {
	var count int64
	database.GetDB().Model(&models.PathChapter{}).
		Where("path_id = ? AND status = ?", path.ID, 1).
		Count(&count)
	path.ChapterCount = int(count)
}

// fillChapters 填充路径的章节列表
func (r *LearningPathRepository) fillChapters(path *models.LearningPath) {
	var chapters []models.PathChapter
	if err := database.GetDB().
		Where("path_id = ? AND status = ?", path.ID, 1).
		Order("order_index ASC").
		Find(&chapters).Error; err != nil {
		path.Chapters = nil
		return
	}
	path.Chapters = chapters
}

// parseJSONList 解析 JSON 字符串为字符串数组
