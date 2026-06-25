package service

import (
	"hot-ai-backend/internal/access"
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
	Search     string `json:"search"`
}

// GetLearningPathsResponse 获取学习路径列表响应
type GetLearningPathsResponse struct {
	List  []models.LearningPath `json:"list"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
}

// GetLearningPaths 获取学习路径列表
func (s *LearningPathService) GetLearningPaths(req *GetLearningPathsRequest) (*GetLearningPathsResponse, error) {
	paths, total, err := s.learningPathRepo.GetList(req.Page, req.PageSize, req.Difficulty, req.Search)
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

// GetLearningPathCount 获取学习路径总数
func (s *LearningPathService) GetLearningPathCount() (int64, error) {
	return s.learningPathRepo.GetLearningPathCount()
}

// LearningPathView 学习路径详情响应 (含 access 决策)
type LearningPathView struct {
	*models.LearningPath
	IsLocked          bool                  `json:"is_locked"`
	RequiredLevel     int                   `json:"required_level,omitempty"`
	RequiredLevelName string                `json:"required_level_name,omitempty"`
	Locked            *access.LockedContent `json:"locked,omitempty"`
	// 内嵌章节也走 access 决策, 用 ChaptersView 替代原始 Chapters
	ChaptersView []PathChapterView `json:"chapters_view,omitempty"`
}

// ToLearningPathView 把 LearningPath 包成 view，根据 userLevel 算 access
// 同时给内嵌的 Chapters 打 access 标签 (effective = max(path, chapter))
func ToLearningPathView(p *models.LearningPath, userLevel int) *LearningPathView {
	v := &LearningPathView{LearningPath: p}
	decision := access.Decide(userLevel, p.AccessLevel)
	v.IsLocked = !decision.Allow
	if !decision.Allow {
		v.RequiredLevel = p.AccessLevel
		v.RequiredLevelName = access.LevelName(p.AccessLevel)
		preview, _ := access.TruncateContent(p.Description, access.GuestPreviewChars)
		p.Description = preview
		lp := access.LockedPlaceholder("学习路径", p.AccessLevel)
		v.Locked = &lp
	}
	// 内嵌章节走 PathChapterListView — 每章按 effective = max(父, 子) 判锁
	if len(p.Chapters) > 0 {
		v.ChaptersView = PathChapterListView(p.Chapters, userLevel, p.AccessLevel)
		// 同时把原始 chapters 清空避免泄漏无锁信息的 content (这里没有 content, 只为安全)
		p.Chapters = nil
	}
	return v
}

// LearningPathListView 给列表里每条路径打 is_locked 标签
func LearningPathListView(paths []models.LearningPath, userLevel int) []LearningPathView {
	out := make([]LearningPathView, 0, len(paths))
	for i := range paths {
		v := LearningPathView{LearningPath: &paths[i]}
		decision := access.Decide(userLevel, paths[i].AccessLevel)
		v.IsLocked = !decision.Allow
		if !decision.Allow {
			v.RequiredLevel = paths[i].AccessLevel
			v.RequiredLevelName = access.LevelName(paths[i].AccessLevel)
		}
		out = append(out, v)
	}
	return out
}

// PathChapterView 章节详情响应 (含 access 决策)
type PathChapterView struct {
	*models.PathChapter
	IsLocked          bool                  `json:"is_locked"`
	RequiredLevel     int                   `json:"required_level,omitempty"`
	RequiredLevelName string                `json:"required_level_name,omitempty"`
	Locked            *access.LockedContent `json:"locked,omitempty"`
}

// ToPathChapterView 把 PathChapter 包成 view，根据 userLevel 算 access
// 章节的 effective level = max(parent path level, chapter level)
func ToPathChapterView(c *models.PathChapter, userLevel int, pathAccessLevel ...int) *PathChapterView {
	effective := c.AccessLevel
	if len(pathAccessLevel) > 0 && pathAccessLevel[0] > effective {
		effective = pathAccessLevel[0]
	}
	v := &PathChapterView{PathChapter: c}
	decision := access.Decide(userLevel, effective)
	v.IsLocked = !decision.Allow
	if !decision.Allow {
		v.RequiredLevel = effective
		v.RequiredLevelName = access.LevelName(effective)
		preview, _ := access.TruncateContent(c.Content, access.GuestPreviewChars)
		c.Content = preview
		lp := access.LockedPlaceholder("章节", effective)
		v.Locked = &lp
	}
	return v
}

// PathChapterListView 给章节列表打 is_locked 标签
// pathAccessLevel 为父路径的 access level，章节 effective = max(父, 子)
func PathChapterListView(chapters []models.PathChapter, userLevel int, pathAccessLevel ...int) []PathChapterView {
	parentLevel := 0
	if len(pathAccessLevel) > 0 {
		parentLevel = pathAccessLevel[0]
	}
	out := make([]PathChapterView, 0, len(chapters))
	for i := range chapters {
		v := PathChapterView{PathChapter: &chapters[i]}
		effective := chapters[i].AccessLevel
		if parentLevel > effective {
			effective = parentLevel
		}
		decision := access.Decide(userLevel, effective)
		v.IsLocked = !decision.Allow
		if !decision.Allow {
			v.RequiredLevel = effective
			v.RequiredLevelName = access.LevelName(effective)
		}
		out = append(out, v)
	}
	return out
}
