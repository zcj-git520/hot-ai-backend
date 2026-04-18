package crawler

import (
	"context"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"hot-ai-backend/internal/models"
)

// startRedisStreamConsumer 启动 Redis Stream 消费者
func StartRedisStreamConsumer(ctx context.Context, c CrawlerConf, client *redis.Client, db *gorm.DB) {
	logx.Info("启动 Redis Stream 消费者...")

	streamKey := "crawler:articles"
	groupKey := "crawler-group"
	consumerKey := "consumer-1"

	// 创建消费者组
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	err := client.XGroupCreateMkStream(ctxTimeout, streamKey, groupKey, "0").Err()
	cancel()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		logx.Error("创建消费者组失败:", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			logx.Info("Redis Stream 消费者停止")
			return
		default:
			// 读取消息
			args := &redis.XReadGroupArgs{
				Group:    groupKey,
				Consumer: consumerKey,
				Streams:  []string{streamKey, ">"},
				Count:    10,
				Block:    5 * time.Second,
			}

			streams, err := client.XReadGroup(ctx, args).Result()
			if err != nil {
				if err != redis.Nil {
					logx.Error("读取 Stream 失败:", err)
				}
				continue
			}

			if len(streams) == 0 || len(streams[0].Messages) == 0 {
				continue
			}

			// 处理消息
			for _, msg := range streams[0].Messages {
				ProcessArticleMessage(ctx, c, msg, db)
				// 确认消息
				client.XAck(ctx, streamKey, groupKey, msg.ID)
			}
		}
	}
}

// ProcessArticleMessage 处理文章消息
func ProcessArticleMessage(ctx context.Context, c CrawlerConf, msg redis.XMessage, db *gorm.DB) {
	logx.Infof("处理文章消息：%+v", msg.Values)

	// 1. 解析文章内容
	sourceID, _ := msg.Values["source_id"].(string)
	sourceName, _ := msg.Values["source_name"].(string)
	sourceURL, _ := msg.Values["source_url"].(string)
	category, _ := msg.Values["category"].(string)
	content, _ := msg.Values["content"].(string)
	fetchedAt, _ := msg.Values["fetched_at"].(string)

	if content == "" {
		logx.Error("内容为空，跳过处理")
		return
	}

	// 2. 清洗内容并提取标题
	cleanedContent := cleanContent(content)
	title := extractTitle(cleanedContent)
	if title == "" {
		title = "未命名文章"
	}

	// 3. 提取标签
	var tags []string
	if tagStr, ok := msg.Values["tags"].(string); ok && tagStr != "" {
		tagStr = strings.TrimSpace(tagStr)
		if tagStr != "" {
			tags = strings.Split(tagStr, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}
	}

	// 4. 构建文章数据
	articleData := map[string]interface{}{
		"id":          uuid.New().String(),
		"source_id":   sourceID,
		"source_name": sourceName,
		"source_url":  sourceURL,
		"category":    category,
		"title":       title,
		"content":     cleanedContent,
		"author":      msg.Values["author"],
		"fetched_at":  fetchedAt,
		"tags":        strings.Join(tags, ","),
	}

	// 5. 存储到 MySQL
	storeToMySQL(ctx, db, articleData)

	logx.Infof("文章消息处理完成：%s", sourceName)
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

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// storeToMySQL 存储文章到 MySQL 数据库
func storeToMySQL(_ context.Context, db *gorm.DB, articleData map[string]interface{}) {
	if db == nil {
		logx.Error("数据库连接未初始化")
		return
	}

	// 1. 获取或创建分类
	categoryName := articleData["category"].(string)
	category, err := createOrUpdateCategory(db, categoryName)
	if err != nil {
		logx.Errorf("创建/更新分类失败: %v", err)
		return
	}

	// 2. 获取或创建来源
	sourceName := articleData["source_name"].(string)
	sourceURL := articleData["source_url"].(string)
	source, err := createOrUpdateSource(db, sourceName, sourceURL)
	if err != nil {
		logx.Errorf("创建/更新来源失败: %v", err)
		return
	}

	// 3. 清洗内容
	cleanedContent := cleanContent(articleData["content"].(string))
	title := articleData["title"].(string)

	// 4. 提取标签
	var tags []string
	if tagStr, ok := articleData["tags"].(string); ok && tagStr != "" {
		tagStr = strings.TrimSpace(tagStr)
		if tagStr != "" {
			tags = strings.Split(tagStr, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}
	}

	// 5. 创建文章
	now := time.Now()
	article := &models.Article{
		Title:           truncateString(title, 200),
		Summary:         truncateString(cleanedContent, 500),
		Content:         cleanedContent,
		OriginalURL:     sourceURL,
		SourceID:        uint(source.ID),
		CategoryID:      uint(category.ID),
		Author:          articleData["author"].(string),
		PublishedAt:     now,
		Status:          "published",
	}

	if err := db.Create(article).Error; err != nil {
		logx.Errorf("创建文章失败: %v", err)
		return
	}

	// 6. 创建文章统计
	stats := &models.ArticleStats{
		ArticleID:    article.ID,
		ViewCount:    0,
		CommentCount: 0,
		LikeCount:    0,
	}
	db.Create(stats)

	// 7. 处理标签关联
	if len(tags) > 0 {
		processArticleTags(db, article.ID, tags)
	}

	// 8. 记录抓取日志
	if sourceID, ok := articleData["source_id"].(string); ok {
		fetchLog := &models.CrawlerFetchLog{
			ID:             uuid.New().String(),
			SourceID:       sourceID,
			FetchStartedAt: now,
			FetchCompletedAt: &now,
			DurationMs:     0,
			StatusCode:     200,
			ResponseSize:   0,
			ItemsFetched:   1,
			Status:         "success",
			CreatedAt:      now,
		}
		db.Create(fetchLog)
	}

	logx.Infof("文章存储成功: %s (ID: %d)", article.Title, article.ID)
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
			Name:     categoryName,
			Code:     strings.ToLower(strings.ReplaceAll(categoryName, " ", "_")),
			Color:    colors[categoryCounter%len(colors)],
			SortOrder: 0,
			Status:   1,
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
				Name: tagName,
				Type: 0,
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
