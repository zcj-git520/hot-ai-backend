package service

import (
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
}

// GetArticlesResponse 获取文章列表响应
type GetArticlesResponse struct {
	Articles []models.Article `json:"articles"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
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

	articles, total, err := s.articleRepo.GetList(req.Page, req.PageSize, req.Category)
	if err != nil {
		return nil, err
	}

	return &GetArticlesResponse{
		Articles: articles,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
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
