package repository

import (
	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/models"
)

// ArticleRepository 文章仓储
type ArticleRepository struct{}

// NewArticleRepository 创建文章仓储实例
func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{}
}

// GetList 获取文章列表(支持分页和分类筛选)
func (r *ArticleRepository) GetList(page, pageSize int, categoryCode string, keyword string) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := database.GetDB().Model(&models.Article{})

	// 分类筛选 - 通过category_code关联查询
	if categoryCode != "" && categoryCode != "all" {
		query = query.Joins("JOIN categories ON articles.category_id = categories.id").
			Where("categories.code = ?", categoryCode)
	}

	// 关键词搜索 - 匹配标题和摘要
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("articles.title LIKE ? OR articles.summary LIKE ?", searchPattern, searchPattern)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("articles.published_at DESC").Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	// 填充关联数据
	for i := range articles {
		r.fillArticleData(&articles[i])
	}

	return articles, total, nil
}

// GetByID 根据ID获取文章详情
func (r *ArticleRepository) GetByID(id uint) (*models.Article, error) {
	var article models.Article
	if err := database.GetDB().Where("id = ?", id).First(&article).Error; err != nil {
		return nil, err
	}

	// 填充关联数据
	r.fillArticleData(&article)

	return &article, nil
}

// GetCategories 获取所有启用的文章分类
func (r *ArticleRepository) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := database.GetDB().Where("status = ?", 1).Order("sort_order ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// fillArticleData 填充文章的关联数据(来源名、分类名、标签、统计数据)
func (r *ArticleRepository) fillArticleData(article *models.Article) {
	db := database.GetDB()

	// 获取来源名称
	var source models.Source
	if err := db.Where("id = ?", article.SourceID).First(&source).Error; err == nil {
		article.SourceName = source.Name
	}

	// 获取分类名称
	var category models.Category
	if err := db.Where("id = ?", article.CategoryID).First(&category).Error; err == nil {
		article.CategoryName = category.Name
	}

	// 获取标签列表
	var relations []models.ArticleTagRelation
	if err := db.Where("article_id = ?", article.ID).Find(&relations).Error; err == nil {
		var tagIDs []uint
		for _, rel := range relations {
			tagIDs = append(tagIDs, rel.TagID)
		}

		if len(tagIDs) > 0 {
			var tags []models.Tag
			if err := db.Where("id IN ?", tagIDs).Find(&tags).Error; err == nil {
				for _, tag := range tags {
					article.Tags = append(article.Tags, tag.Name)
				}
			}
		}
	}

	// 获取统计数据
	var stats models.ArticleStats
	if err := db.Where("article_id = ?", article.ID).First(&stats).Error; err == nil {
		article.ViewCount = stats.ViewCount
		article.CommentCount = stats.CommentCount
		article.LikeCount = stats.LikeCount
	}
}

// IncrementViewCount 增加文章阅读量
func (r *ArticleRepository) IncrementViewCount(id uint) error {
	db := database.GetDB()

	// 先查找或创建统计记录
	var stats models.ArticleStats
	result := db.Where("article_id = ?", id).First(&stats)

	if result.Error != nil {
		// 不存在则创建
		stats = models.ArticleStats{
			ArticleID: id,
			ViewCount: 1,
		}
		return db.Create(&stats).Error
	}

	// 存在则更新
	return db.Model(&stats).Update("view_count", stats.ViewCount+1).Error
}

// GetCount 获取文章总数
func (r *ArticleRepository) GetCount() (int64, error) {
	var total int64
	if err := database.GetDB().Model(&models.Article{}).Where("status = ?", "published").Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
