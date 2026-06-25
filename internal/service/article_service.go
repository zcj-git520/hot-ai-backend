package service

import (
	"hot-ai-backend/internal/access"
	"hot-ai-backend/internal/models"
	"hot-ai-backend/internal/repository"
	"strconv"
)

// ArticleService 文章服务
type ArticleService struct {
	articleRepo *repository.ArticleRepository
}

// NewArticleService 创建文章服务实例
func NewArticleService(articleRepo *repository.ArticleRepository) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
	}
}

// GetArticlesRequest 获取文章列表请求参数
type GetArticlesRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Category string `json:"category"`
	Keyword  string `json:"keyword"`
}

// GetArticlesResponse 获取文章列表响应
type GetArticlesResponse struct {
	Articles    []models.Article `json:"articles"`
	Total       int64            `json:"total"`
	TotalPages  int              `json:"total_pages"`
	Page        int              `json:"page"`
	PageSize    int              `json:"page_size"`
}

// GetArticles 获取文章列表
func (s *ArticleService) GetArticles(req *GetArticlesRequest) (*GetArticlesResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	articles, total, err := s.articleRepo.GetList(req.Page, req.PageSize, req.Category, req.Keyword)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &GetArticlesResponse{
		Articles:   articles,
		Total:      total,
		TotalPages: totalPages,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

// GetArticleByID 根据ID获取文章详情
func (s *ArticleService) GetArticleByID(idStr string) (*models.Article, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	article, err := s.articleRepo.GetByID(uint(id))
	if err != nil {
		return nil, err
	}

	// 增加阅读量
	_ = s.articleRepo.IncrementViewCount(uint(id))

	return article, nil
}

// GetCategories 获取文章分类列表
func (s *ArticleService) GetCategories() ([]models.Category, error) {
	return s.articleRepo.GetCategories()
}

// GetArticleCount 获取文章总数
func (s *ArticleService) GetArticleCount() (int64, error) {
	return s.articleRepo.GetCount()
}

// ArticleView 文章详情响应 (含 access 决策结果)
type ArticleView struct {
	*models.Article
	IsLocked          bool                `json:"is_locked"`
	RequiredLevel     int                 `json:"required_level,omitempty"`
	RequiredLevelName string              `json:"required_level_name,omitempty"`
	Locked            *access.LockedContent `json:"locked,omitempty"`
}

// ToView 把 article 包成 view，根据 userLevel 算 access
func ToView(a *models.Article, userLevel int) *ArticleView {
	v := &ArticleView{Article: a}
	decision := access.Decide(userLevel, a.AccessLevel)
	v.IsLocked = !decision.Allow
	if !decision.Allow {
		v.RequiredLevel = a.AccessLevel
		v.RequiredLevelName = access.LevelName(a.AccessLevel)
		// 详情场景：不够级别时裁剪正文，给游客看 500 字预览
		preview, _ := access.TruncateContent(a.Content, access.GuestPreviewChars)
		a.Content = preview
		lp := access.LockedPlaceholder("文章", a.AccessLevel)
		v.Locked = &lp
	}
	return v
}

// ToViewList 给列表里每篇打 is_locked 标签
func ToViewList(articles []models.Article, userLevel int) []ArticleView {
	out := make([]ArticleView, 0, len(articles))
	for i := range articles {
		v := ArticleView{Article: &articles[i]}
		decision := access.Decide(userLevel, articles[i].AccessLevel)
		v.IsLocked = !decision.Allow
		if !decision.Allow {
			v.RequiredLevel = articles[i].AccessLevel
			v.RequiredLevelName = access.LevelName(articles[i].AccessLevel)
		}
		out = append(out, v)
	}
	return out
}
