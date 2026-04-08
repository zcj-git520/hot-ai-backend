package service

import (
	"hot-ai-backend/internal/models"
	"hot-ai-backend/internal/repository"
)

// LearningPathService 学习路径服务
type LearningPathService struct {
	learningPathRepo *repository.LearningPathRepository
}

// NewLearningPathService 创建学习路径服务实例
func NewLearningPathService(learningPathRepo *repository.LearningPathRepository) *LearningPathService {
	return &LearningPathService{
		learningPathRepo: learningPathRepo,
	}
}

// GetLearningPathsRequest 获取学习路径列表请求
type GetLearningPathsRequest struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Difficulty string `json:"difficulty"`
}

// GetLearningPathsResponse 获取学习路径列表响应
type GetLearningPathsResponse struct {
	List  []models.LearningPath `json:"list"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
}

// GetLearningPaths 获取学习路径列表
func (s *LearningPathService) GetLearningPaths(req *GetLearningPathsRequest) (*GetLearningPathsResponse, error) {
	paths, total, err := s.learningPathRepo.GetList(req.Page, req.PageSize, req.Difficulty)
	if err != nil {
		return nil, err
	}

	return &GetLearningPathsResponse{
		List:  paths,
		Total: total,
		Page:  req.Page,
	}, nil
}

// GetLearningPathByID 根据ID获取学习路径详情
func (s *LearningPathService) GetLearningPathByID(id uint) (*models.LearningPath, error) {
	path, err := s.learningPathRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return path, nil
}

// GetLearningPathBySlug 根据slug获取学习路径详情
func (s *LearningPathService) GetLearningPathBySlug(slug string) (*models.LearningPath, error) {
	path, err := s.learningPathRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	return path, nil
}

// GetPathChapters 获取路径的所有章节
func (s *LearningPathService) GetPathChapters(pathID uint) ([]models.PathChapter, error) {
	chapters, err := s.learningPathRepo.GetChapters(pathID)
	if err != nil {
		return nil, err
	}
	return chapters, nil
}

// GetChapterByID 根据章节ID获取详情
func (s *LearningPathService) GetChapterByID(chapterID uint) (*models.PathChapter, error) {
	chapter, err := s.learningPathRepo.GetChapterByID(chapterID)
	if err != nil {
		return nil, err
	}
	return chapter, nil
}

// GetChapterBySlug 根据路径ID和章节slug获取详情
func (s *LearningPathService) GetChapterBySlug(pathID uint, slug string) (*models.PathChapter, error) {
	chapter, err := s.learningPathRepo.GetChapterBySlug(pathID, slug)
	if err != nil {
		return nil, err
	}
	return chapter, nil
}

// GetPrevNextChapter 获取前一章和下一章
func (s *LearningPathService) GetPrevNextChapter(pathID, currentChapterID uint) (prev, next *models.PathChapter, err error) {
	prev, next, err = s.learningPathRepo.GetPrevNextChapter(pathID, currentChapterID)
	if err != nil {
		return nil, nil, err
	}
	return prev, next, nil
}

// GetLevelInfo 获取难度等级信息
func (s *LearningPathService) GetLevelInfo() []models.LevelInfo {
	return s.learningPathRepo.GetLevelInfo()
}

// GetFeaturedPaths 获取推荐路径
func (s *LearningPathService) GetFeaturedPaths(limit int) ([]models.LearningPath, error) {
	paths, err := s.learningPathRepo.GetFeatured(limit)
	if err != nil {
		return nil, err
	}
	return paths, nil
}

// GetPathProgress 获取用户的学习进度
func (s *LearningPathService) GetPathProgress(userID string, pathID uint) (*models.LearningProgress, error) {
	progress, err := s.learningPathRepo.GetPathProgress(userID, pathID)
	if err != nil {
		return nil, err
	}
	return progress, nil
}

// GetCompletedChapters 获取用户已完成的章节列表
func (s *LearningPathService) GetCompletedChapters(userID string, pathID uint) ([]uint, error) {
	chapterIDs, err := s.learningPathRepo.GetCompletedChapters(userID, pathID)
	if err != nil {
		return nil, err
	}
	return chapterIDs, nil
}

// SaveProgressRequest 保存学习进度请求
type SaveProgressRequest struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	PathID    uint   `json:"path_id"`
	ChapterID uint   `json:"chapter_id"`
	Status    string `json:"status"`
	TimeSpent int    `json:"time_spent"`
	Notes     string `json:"notes"`
}

// SaveProgress 保存学习进度
func (s *LearningPathService) SaveProgress(req *SaveProgressRequest) error {
	progress := &models.LearningProgress{
		UserID:    req.UserID,
		SessionID: req.SessionID,
		PathID:    req.PathID,
		ChapterID: req.ChapterID,
		Status:    req.Status,
		TimeSpent: req.TimeSpent,
		Notes:     req.Notes,
	}

	if err := s.learningPathRepo.SaveProgress(progress); err != nil {
		return err
	}

	// 如果是完成状态，更新路径管理数据
	if req.Status == "completed" {
		_ = s.learningPathRepo.IncrementCompleteCount(req.PathID)
	}

	return nil
}

// IncrementStartCount 增加开始学习次数
func (s *LearningPathService) IncrementStartCount(pathID uint) error {
	return s.learningPathRepo.IncrementStartCount(pathID)
}
