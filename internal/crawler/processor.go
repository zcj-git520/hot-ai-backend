package crawler

import (
	"context"
	"hot-ai-backend/internal/models"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// ProcessAndStoreArticle 处理并存储文章到数据库
// isRSS: 是否为 RSS 源，RSS 源跳过内容清洗和标题提取
func ProcessAndStoreArticle(ctx context.Context, article models.Article, db *gorm.DB, isRSS bool, translateClient *TranslateClient) error {
	logx.Infof("处理文章内容: %v (RSS: %v)", article.Title, isRSS)

	// 1. 检查必要字段
	if article.Content == "" {
		logx.Error("内容为空，跳过处理")
		return nil
	}

	// 2. 处理内容（RSS 源跳过清洗）
	if !isRSS {
		// HTML 源需要清洗内容
		article.Content = cleanContent(article.Content)
		if article.Content == "" {
			logx.Error("清洗后内容为空，跳过处理")
			return nil
		}
	}

	// 3. 翻译文章（如果启用翻译服务）
	if translateClient != nil {
		err := translateClient.TranslateArticle(ctx, &article)
		if err != nil {
			logx.Errorf("翻译文章失败: %v", err)
		}
		logx.Infof("文章翻译成功")

	}

	// 4. 提取标签（从 article.Tags 切片）
	var tags []string
	if len(article.Tags) > 0 {
		tags = make([]string, len(article.Tags))
		copy(tags, article.Tags)
	}

	// 5. 检查数据库连接
	if db == nil {
		logx.Error("数据库连接未初始化")
		return nil
	}

	// 6. 获取或创建分类
	//categoryName := "默认分类"
	//if article.CategoryName != "" {
	//	categoryName = article.CategoryName
	//}
	//categoryModel, err := createOrUpdateCategory(db, categoryName)
	//if err != nil {
	//	logx.Errorf("创建/更新分类失败: %v", err)
	//	return err
	//}

	// 7. 获取或创建来源
	sourceName := article.SourceName
	if sourceName == "" {
		sourceName = "未知来源"
	}
	sourceURL := article.OriginalURL
	if sourceURL == "" {
		sourceURL = "http://unknown.com"
	}
	source, err := createOrUpdateSource(db, sourceName, sourceURL)
	if err != nil {
		logx.Errorf("创建/更新来源失败: %v", err)
		return err
	}

	// 8. 准备文章数据
	now := time.Now()
	publishedAt := now
	if !article.PublishedAt.IsZero() {
		publishedAt = article.PublishedAt
	}

	// 设置文章字段（RSS 源跳过标题提取）
	if article.Title == "" && !isRSS {
		// HTML 源需要从内容中提取标题
		article.Title = extractTitle(article.Content)
		if article.Title == "" {
			article.Title = "未命名文章"
		}
	}
	article.Summary = truncateString(article.Content, 200)
	article.SummaryEn = truncateString(article.ContentEn, 200)
	article.SourceID = source.ID
	article.CategoryID = 1
	if article.Author == "" {
		article.Author = "未知作者"
	}
	article.PublishedAt = publishedAt
	if article.Status == 0 {
		article.Status = 1 // 1-已发布
	}

	// 9. 创建文章
	if err := db.Create(&article).Error; err != nil {
		logx.Errorf("创建文章失败: %v", err)
		return err
	}

	// 10. 创建文章统计
	stats := &models.ArticleStats{
		ArticleID:    article.ID,
		ViewCount:    0,
		CommentCount: 0,
		LikeCount:    0,
	}
	db.Create(stats)

	// 11. 处理标签关联
	if len(tags) > 0 {
		processArticleTags(db, article.ID, tags)
	}

	logx.Infof("文章存储成功: %s (ID: %d)", article.Title, article.ID)
	return nil
}

// cleanContent 清洗 HTML 内容
func cleanContent(content string) string {
	if content == "" {
		return ""
	}

	// 1. 解码 HTML 实体
	content = html.UnescapeString(content)

	// 2. 去除 HTML 标签
	reg := regexp.MustCompile(`<[^>]*>`)
	content = reg.ReplaceAllString(content, " ")

	// 3. 去除多余空白字符
	content = strings.Join(strings.Fields(content), " ")

	// 4. 去除首尾空白
	content = strings.TrimSpace(content)

	return content
}

// extractTitle 从内容中提取标题
func extractTitle(content string) string {
	// 方法1: 查找 <title> 标签
	titleRegex := regexp.MustCompile(`<title>(.*?)</title>`)
	matches := titleRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// 方法2: 查找 H1 标签
	h1Regex := regexp.MustCompile(`<h1[^>]*>(.*?)</h1>`)
	matches = h1Regex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(html.UnescapeString(matches[1]))
	}

	// 方法3: 查找 meta description
	metaRegex := regexp.MustCompile(`<meta[^>]*name=["']description["'][^>]*content=["'](.*?)["']`)
	matches = metaRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(html.UnescapeString(matches[1]))
	}

	// 默认返回空
	return ""
}

// truncateString 截断字符串（按字符而非字节，避免UTF-8多字节字符被截断）
func truncateString(s string, maxLen int) string {
	if len([]rune(s)) <= maxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxLen]) + "..."
}

// createOrUpdateCategory 创建或更新分类
func createOrUpdateCategory(db *gorm.DB, categoryName string) (*models.Category, error) {
	var category models.Category
	err := db.Where("name = ?", categoryName).First(&category).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新分类
		colors := []string{"#FF5733", "#33FF57", "#3357FF", "#F333FF", "#33FFF5", "#FFFF33", "#FF8C33", "#8C33FF"}
		categoryCounter := getAndUpdateCategoryCounter()
		category = models.Category{
			Name:      categoryName,
			Code:      strings.ToLower(strings.ReplaceAll(categoryName, " ", "_")),
			Color:     colors[categoryCounter%len(colors)],
			SortOrder: 0,
			Status:    1,
		}
		if err := db.Create(&category).Error; err != nil {
			return nil, err
		}
		logx.Infof("创建分类: %s", categoryName)
		return &category, nil
	} else if err != nil {
		return nil, err
	}

	return &category, nil
}

// categoryCounter 分类计数器
var categoryCounter int

// getAndUpdateCategoryCounter 获取并更新分类计数器
func getAndUpdateCategoryCounter() int {
	categoryCounter++
	return categoryCounter
}

// createOrUpdateSource 创建或更新来源
func createOrUpdateSource(db *gorm.DB, name, domain string) (*models.Source, error) {
	var source models.Source
	err := db.Where("domain = ?", domain).First(&source).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新来源
		source = models.Source{
			Name:             name,
			Domain:           domain,
			LogoURL:          "",
			Description:      "",
			ReliabilityScore: 5,
			Status:           1,
		}
		if err := db.Create(&source).Error; err != nil {
			return nil, err
		}
		logx.Infof("创建来源: %s (%s)", name, domain)
		return &source, nil
	} else if err != nil {
		return nil, err
	}

	return &source, nil
}

// processArticleTags 处理文章标签关联
func processArticleTags(db *gorm.DB, articleID uint, tags []string) {
	for _, tagName := range tags {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

		// 获取或创建标签
		var tag models.Tag
		err := db.Where("name = ?", tagName).First(&tag).Error

		if err == gorm.ErrRecordNotFound {
			tag = models.Tag{
				Name:   tagName,
				Type:   0,
				Status: 1,
			}
			if err := db.Create(&tag).Error; err != nil {
				logx.Errorf("创建标签失败: %v", err)
				continue
			}
			logx.Infof("创建标签: %s", tagName)
		} else if err != nil {
			logx.Errorf("查询标签失败: %v", err)
			continue
		}

		// 创建文章标签关联
		relation := models.ArticleTagRelation{
			ArticleID: articleID,
			TagID:     tag.ID,
		}

		if err := db.Create(&relation).Error; err != nil {
			logx.Errorf("创建文章标签关联失败: %v", err)
		}
	}
}
