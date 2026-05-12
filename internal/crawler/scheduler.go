package crawler

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"hot-ai-backend/internal/models"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// StartScheduler 启动定时任务调度器
func StartScheduler(ctx context.Context, c CrawlerConf, db *gorm.DB, minifluxClient *MinifluxClient, aiClient *AIClient) {
	logx.Info("启动定时任务调度器...")

	ticker := time.NewTicker(time.Duration(c.Crawler.FetchInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logx.Info("定时任务调度器停止")
			return
		case <-ticker.C:
			logx.Info("执行定时抓取任务...")

			// 1. 处理 RSS 源（通过 Miniflux）
			if minifluxClient != nil && c.Miniflux.Enabled {
				processRSSSources(ctx, c, db, minifluxClient, aiClient)
			}

			// 2. 处理传统 HTML 源（根据配置开关决定）
			if c.Crawler.EnableTraditionalCrawling {
				ExecuteFetchTasks(ctx, c, db, aiClient)
			} else {
				logx.Debug("传统 HTML 爬虫已禁用，跳过执行")
			}
		}
	}
}

// ExecuteFetchTasks 执行抓取任务
func ExecuteFetchTasks(ctx context.Context, c CrawlerConf, db *gorm.DB, aiClient *AIClient) {
	if db == nil {
		logx.Error("数据库未初始化")
		return
	}

	startTime := time.Now()
	logx.Info("开始执行抓取任务...")

	// 1. 从数据库读取活跃的抓取源（按优先级排序）
	var sources []models.CrawlerSource
	now := time.Now()
	err := db.Where("status = ? AND (next_fetch_at IS NULL OR next_fetch_at <= ?)",
		models.CrawlerSourceStatusActive, now).
		Order("priority DESC, id").
		Limit(c.Crawler.ConcurrentFetch). // 限制并发数量
		Find(&sources).Error

	if err != nil {
		logx.Error("读取抓取源失败:", err)
		return
	}

	if len(sources) == 0 {
		logx.Info("本次没有需要抓取的源")
		return
	}

	logx.Infof("找到 %d 个待抓取源", len(sources))

	// 2. 并发执行抓取任务
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, c.Crawler.MaxConcurrency) // 信号量控制最大并发数

	for _, source := range sources {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量

		go func(src models.CrawlerSource) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			executeSingleFetchTask(ctx, c, db, src, aiClient)
		}(source)
	}

	// 等待所有任务完成
	wg.Wait()

	elapsed := time.Since(startTime)
	logx.Infof("抓取任务完成，耗时：%v", elapsed)
}

// executeSingleFetchTask 执行单个抓取任务
func executeSingleFetchTask(ctx context.Context, c CrawlerConf, db *gorm.DB, source models.CrawlerSource, aiClient *AIClient) {
	logx.Infof("开始抓取源：%s (%s)", source.Name, source.URL)

	fetchStartTime := time.Now()
	fetchLog := &models.CrawlerFetchLog{
		ID:             uuid.New().String(),
		SourceID:       source.ID,
		FetchStartedAt: fetchStartTime,
	}

	defer func() {
		// 更新抓取日志
		fetchCompletedAt := time.Now()
		fetchLog.FetchCompletedAt = &fetchCompletedAt
		fetchLog.DurationMs = int(time.Since(fetchStartTime).Milliseconds())
		db.Create(fetchLog)

		// 更新抓取源状态
		updateCrawlerSource(db, source.ID, fetchLog)
	}()

	// 执行 HTTP 请求
	resp, err := sendHTTPRequest(ctx, source)
	if err != nil {
		logx.Errorf("抓取源 %s 失败：%v", source.Name, err)
		fetchLog.Status = models.FetchLogStatusFailed
		fetchLog.ErrorMessage = err.Error()
		return
	}
	defer resp.Body.Close()

	fetchLog.StatusCode = resp.StatusCode

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("读取响应失败：%s: %v", source.Name, err)
		fetchLog.Status = models.FetchLogStatusFailed
		fetchLog.ErrorMessage = fmt.Sprintf("读取响应失败：%v", err)
		return
	}

	fetchLog.ResponseSize = len(body)

	// 检查 HTTP 状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logx.Errorf("抓取源 %s HTTP 错误：%d", source.Name, resp.StatusCode)
		fetchLog.Status = models.FetchLogStatusFailed
		fetchLog.ErrorMessage = fmt.Sprintf("HTTP 错误：%d", resp.StatusCode)
		return
	}

	// 3. 处理抓取结果（直接存储到数据库）
	itemsFetched, err := processFetchedContent(ctx, c, db, source, body, aiClient)
	if err != nil {
		logx.Errorf("处理抓取内容失败：%s: %v", source.Name, err)
		fetchLog.Status = models.FetchLogStatusFailed
		fetchLog.ErrorMessage = fmt.Sprintf("处理内容失败：%v", err)
		return
	}

	// 成功
	fetchLog.Status = models.FetchLogStatusSuccess
	fetchLog.ItemsFetched = itemsFetched
	logx.Infof("抓取源 %s 成功，获取到 %d 条内容", source.Name, itemsFetched)
}

// sendHTTPRequest 发送 HTTP 请求
func sendHTTPRequest(ctx context.Context, source models.CrawlerSource) (*http.Response, error) {
	timeout := time.Duration(source.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 60 * time.Second // 默认增加到 60 秒
	}

	// 创建带超时的 HTTP 客户端
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			// TLS 配置，跳过证书验证（用于测试）
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // 生产环境应设为 true
			},
		},
	}

	req, err := http.NewRequestWithContext(ctx, string(source.FetchMethod), source.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败：%w", err)
	}

	// 设置更真实的 User-Agent（模拟浏览器）
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	req.Header.Set("User-Agent", userAgent)

	// 设置 Accept 头
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// 如果有自定义请求头，合并进去
	if len(source.RequestHeaders) > 0 {
		var headers map[string]string
		if err := json.Unmarshal(source.RequestHeaders, &headers); err == nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
	}

	logx.Infof("发送请求到: %s (超时: %ds)", source.URL, int(timeout.Seconds()))

	// 执行请求，添加重试逻辑
	var resp *http.Response
	maxRetries := source.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = client.Do(req)
		if err == nil {
			return resp, nil
		}

		logx.Errorf("请求失败 (尝试 %d/%d): %v", attempt, maxRetries, err)

		if attempt < maxRetries {
			// 指数退避重试
			waitTime := time.Duration(attempt*2) * time.Second
			logx.Infof("等待 %v 后重试...", waitTime)
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("经过 %d 次重试后仍然失败: %w", maxRetries, err)
}

// processFetchedContent 处理抓取到的内容（支持两级爬取：列表页 + 详情页）
func processFetchedContent(ctx context.Context, c CrawlerConf, db *gorm.DB, source models.CrawlerSource, content []byte, aiClient *AIClient) (int, error) {
	logx.Infof("开始处理抓取内容: %s (类型: %s)", source.Name, source.SourceType)

	// 解析 parse_rules
	var rules map[string]interface{}
	if len(source.ParseRules) > 0 {
		if err := json.Unmarshal(source.ParseRules, &rules); err != nil {
			logx.Errorf("解析 parse_rules 失败: %v", err)
			rules = make(map[string]interface{})
		}
	}

	// 根据源类型选择不同的处理方式
	switch source.SourceType {
	case models.CrawlerSourceTypeHTML:
		// HTML 类型：先解析列表页获取文章链接，再逐个抓取详情
		return processHTMLSource(ctx, c, db, source, string(content), rules, aiClient)
	default:
		// 其他类型：直接处理内容（非 RSS，需要解析）
		articleData := models.Article{}

		if err := ProcessAndStoreArticle(ctx, articleData, db, false, aiClient); err != nil {
			return 0, fmt.Errorf("处理并存储文章失败：%w", err)
		}

		return 1, nil
	}
}

// processHTMLSource 处理 HTML 类型的抓取源（两级爬取）
func processHTMLSource(ctx context.Context, c CrawlerConf, db *gorm.DB, source models.CrawlerSource, htmlContent string, rules map[string]interface{}, aiClient *AIClient) (int, error) {
	logx.Infof("开始解析 HTML 列表页: %s", source.Name)

	// 第1步：从列表页提取文章链接
	articleLinks, err := ParseArticleLinks(htmlContent, source.URL, rules)
	if err != nil {
		return 0, fmt.Errorf("解析文章链接失败: %w", err)
	}

	if len(articleLinks) == 0 {
		logx.Info("未从列表页找到文章链接: " + source.Name)
		return 0, nil
	}

	logx.Infof("从列表页找到 %d 个文章链接", len(articleLinks))

	// 第2步：逐个抓取文章详情
	successCount := 0
	maxArticles := 10 // 每次最多抓取10篇文章，避免过度请求

	for i, link := range articleLinks {
		if i >= maxArticles {
			logx.Infof("已达到最大抓取数量 (%d)，停止抓取", maxArticles)
			break
		}

		select {
		case <-ctx.Done():
			logx.Info("上下文取消，停止抓取")
			return successCount, ctx.Err()
		default:
		}

		logx.Infof("[%d/%d] 抓取文章: %s - %s", i+1, len(articleLinks), link.Title, link.URL)

		// 添加延迟，避免请求过快
		time.Sleep(2 * time.Second)

		// 抓取文章详情页
		article, err := fetchArticleDetail(ctx, source, link, rules)
		if err != nil {
			logx.Errorf("抓取文章详情失败 [%s]: %v", link.URL, err)
			continue
		}

		// 补充元数据并转换为 models.Article
		articleData := map[string]interface{}{
			"source_id":   fmt.Sprintf("%d", source.ID),
			"source_name": source.Name,
			"category":    source.Category,
		}
		for k, v := range article {
			articleData[k] = v
		}

		// 转换为 models.Article
		articleModel := models.Article{
			Title:        getStringValue(articleData, "title"),
			Content:      getStringValue(articleData, "content"),
			Summary:      getStringValue(articleData, "summary"),
			Author:       getStringValue(articleData, "author"),
			OriginalURL:  getStringValue(articleData, "original_url"),
			SourceName:   getStringValue(articleData, "source_name"),
			CategoryName: getStringValue(articleData, "category"),
			Tags:         []string{},
		}

		// 解析发布时间
		if pubAtStr := getStringValue(articleData, "published_at"); pubAtStr != "" {
			if pubAt, err := time.Parse(time.RFC3339, pubAtStr); err == nil {
				articleModel.PublishedAt = pubAt
			} else {
				articleModel.PublishedAt = time.Now()
			}
		} else {
			articleModel.PublishedAt = time.Now()
		}

		// 存储到数据库（HTML 源需要解析）
		if err := ProcessAndStoreArticle(ctx, articleModel, db, false, aiClient); err != nil {
			logx.Errorf("存储文章失败 [%s]: %v", link.URL, err)
			continue
		}

		successCount++
		logx.Infof("成功抓取并存储文章: %s", articleModel.Title)
	}

	logx.Infof("HTML 源抓取完成: %s, 成功: %d/%d", source.Name, successCount, len(articleLinks))
	return successCount, nil
}

// fetchArticleDetail 抓取文章详情页
func fetchArticleDetail(ctx context.Context, source models.CrawlerSource, link ArticleLink, rules map[string]interface{}) (map[string]interface{}, error) {
	// 发送 HTTP 请求
	timeout := time.Duration(source.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", link.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置 User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; HotAI-Crawler/1.0)")

	// 如果有自定义请求头，合并进去
	if len(source.RequestHeaders) > 0 {
		var headers map[string]string
		if err := json.Unmarshal(source.RequestHeaders, &headers); err == nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
	}

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP 错误: %d", resp.StatusCode)
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析文章详情
	article, err := ParseArticleDetail(string(body), link.URL, rules)
	if err != nil {
		return nil, fmt.Errorf("解析文章详情失败: %w", err)
	}

	// 使用链接中的标题（如果解析结果为空）
	if article["title"] == "" || article["title"] == "未知文章" {
		article["title"] = link.Title
	}

	// 使用链接中的发布时间（如果解析结果为空）
	if article["published_at"] == "" && link.PublishedAt != nil {
		article["published_at"] = link.PublishedAt.Format(time.RFC3339)
	}

	return article, nil
}

// updateCrawlerSource 更新抓取源状态
func updateCrawlerSource(db *gorm.DB, sourceID string, fetchLog *models.CrawlerFetchLog) {
	updates := map[string]interface{}{
		"last_fetch_at":    gorm.Expr("NOW()"),
		"next_fetch_at":    gorm.Expr("DATE_ADD(NOW(), INTERVAL @fetch_interval SECOND)", gorm.Expr("@fetch_interval = (SELECT fetch_interval FROM (SELECT fetch_interval FROM crawler_sources WHERE id = ?) AS tmp)"), sourceID),
		"last_status_code": fetchLog.StatusCode,
		"total_fetches":    gorm.Expr("total_fetches + 1"),
	}

	if fetchLog.Status == models.FetchLogStatusSuccess {
		updates["successful_fetches"] = gorm.Expr("successful_fetches + 1")
		updates["consecutive_failures"] = 0
		updates["last_error_message"] = ""
		if fetchLog.ItemsFetched > 0 {
			updates["total_items"] = gorm.Expr("total_items + ?", fetchLog.ItemsFetched)
		}
	} else {
		updates["failed_fetches"] = gorm.Expr("failed_fetches + 1")
		updates["consecutive_failures"] = gorm.Expr("consecutive_failures + 1")
		updates["last_error_message"] = fetchLog.ErrorMessage

		// 如果连续失败超过一定次数，暂停抓取源
		var source models.CrawlerSource
		if db.First(&source, sourceID).Error == nil {
			if source.ConsecutiveFailures >= source.MaxRetries {
				updates["status"] = models.CrawlerSourceStatusError
			}
		}
	}

	db.Model(&models.CrawlerSource{}).Where("id = ?", sourceID).Updates(updates)
}

// processRSSSources 处理 RSS 源（直接从 Miniflux 获取所有 feeds）
func processRSSSources(ctx context.Context, c CrawlerConf, db *gorm.DB, minifluxClient *MinifluxClient, aiClient *AIClient) {
	logx.Info("开始处理 RSS 源（从 Miniflux 直接获取）...")

	// 1. 直接从 Miniflux 获取所有订阅源
	feeds, err := minifluxClient.GetFeeds()
	if err != nil {
		logx.Errorf("从 Miniflux 获取订阅源失败: %v", err)
		return
	}

	if len(feeds) == 0 {
		logx.Info("Miniflux 中没有订阅源")
		return
	}

	logx.Infof("从 Miniflux 获取到 %d 个订阅源", len(feeds))

	// 2. 逐个处理每个 feed
	for _, feed := range feeds {
		select {
		case <-ctx.Done():
			logx.Info("上下文取消，停止 RSS 抓取")
			return
		default:
		}

		logx.Infof("处理 Feed: %s (%s)", feed.Title, feed.FeedURL)

		// 创建虚拟的 CrawlerSource 用于记录和处理
		virtualSource := models.CrawlerSource{
			Name:       feed.Title,
			URL:        feed.FeedURL,
			SourceType: models.CrawlerSourceTypeRSS,
			Category:   "rss",
			Status:     models.CrawlerSourceStatusActive,
		}

		// 从 Miniflux 获取该 feed 的最新文章
		articles, err := fetchArticlesFromMinifluxFeed(ctx, minifluxClient, feed, virtualSource, aiClient)
		if err != nil {
			logx.Errorf("从 Miniflux 获取文章失败 [%s]: %v", feed.Title, err)
			continue
		}

		logx.Infof("Feed %s 处理完成，获取 %d 篇文章", feed.Title, len(articles))

		// 添加延迟，避免请求过快
		time.Sleep(2 * time.Second)
	}

	logx.Info("RSS 源处理完成")
}

// fetchArticlesFromMinifluxFeed 从 Miniflux 获取指定 Feed 的文章
func fetchArticlesFromMinifluxFeed(ctx context.Context, client *MinifluxClient, feed Feed, source models.CrawlerSource, aiClient *AIClient) ([]map[string]interface{}, error) {
	logx.Infof("从 Miniflux 获取 Feed 文章: %s (ID: %d)", feed.Title, feed.ID)

	// 获取该 Feed 的最新条目（最多 10 篇）
	limit := 10
	entries, err := client.GetFeedEntries(feed.ID, limit)
	if err != nil {
		return nil, fmt.Errorf("获取条目失败: %w", err)
	}

	if len(entries) == 0 {
		logx.Info("没有找到新文章")
		return nil, nil
	}

	logx.Infof("找到 %d 篇新文章", len(entries))

	var articles []map[string]interface{}

	// 处理每篇文章
	for i, entry := range entries {
		select {
		case <-ctx.Done():
			return articles, ctx.Err()
		default:
		}

		logx.Infof("[%d/%d] 处理文章: %s", i+1, len(entries), entry.Title)

		// 方法1: 使用 Miniflux 的 fetch-content API
		fullEntry, err := client.FetchEntryContent(entry.ID)
		if err != nil {
			logx.Info("获取文章内容失败 [%s]: %v", entry.Title, err)
			continue
		}
		if len(fullEntry) < 100 {
			logx.Info("文章内容过短 [%s]: %s", entry.Title, fullEntry)
			continue
		}

		// 解析发布时间
		publishedAt := entry.Date
		if publishedAt.IsZero() {
			publishedAt = time.Now()
		}

		// 从 entry.Tags 提取标签
		tags := entry.Tags
		if tags == nil {
			tags = []string{}
		}

		// 构建 Article 模型
		articleModel := models.Article{
			Title:        entry.Title,
			Content:      fullEntry,
			Summary:      truncateString(entry.Content, 200),
			Author:       entry.Author,
			PublishedAt:  publishedAt,
			OriginalURL:  entry.URL,
			SourceName:   feed.Title,
			CategoryName: source.Category,
			Tags:         tags,
			Status:       1, // 1-已发布
			// 注意: SourceID 和 CategoryID 会在 ProcessAndStoreArticle 中自动处理
			// 因为这里使用的是虚拟的 CrawlerSource，没有真实的 ID
		}

		// 存储到数据库（RSS 源不需要解析文章内容）
		if err := ProcessAndStoreArticle(ctx, articleModel, getDatabaseConnection(), true, aiClient); err != nil {
			logx.Errorf("存储文章失败 [%s]: %v", entry.Title, err)
		} else {
			// 将 Article 模型转换为 map 用于返回
			articleMap := map[string]interface{}{
				"id":           articleModel.ID,
				"title":        articleModel.Title,
				"content":      articleModel.Content,
				"summary":      articleModel.Summary,
				"author":       articleModel.Author,
				"published_at": articleModel.PublishedAt.Format(time.RFC3339),
				"original_url": articleModel.OriginalURL,
				"source_name":  articleModel.SourceName,
				"category":     articleModel.CategoryName,
			}
			articles = append(articles, articleMap)
			logx.Infof("成功存储文章: %s (ID: %d)", entry.Title, articleModel.ID)
		}

		// 添加延迟
		time.Sleep(1 * time.Second)
	}

	// 处理完成后，标记该 Feed 的所有条目为已读
	if len(articles) > 0 {
		logx.Infof("准备标记 Feed %d 为已读", feed.ID)
		if err := client.MarkFeedAsRead(feed.ID); err != nil {
			logx.Errorf("标记 Feed %d 为已读失败: %v", feed.ID, err)
		} else {
			logx.Infof("成功标记 Feed %d 为已读，共 %d 篇文章", feed.ID, len(articles))
		}
	}

	return articles, nil
}

// getDatabaseConnection 获取数据库连接（需要在 service.go 中暴露）
var globalDB *gorm.DB

// SetGlobalDB 设置全局数据库连接（在 service.go 中调用）
func SetGlobalDB(db *gorm.DB) {
	globalDB = db
}

func getDatabaseConnection() *gorm.DB {
	return globalDB
}

// getStringValue 从 map 中安全地获取字符串值
func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
