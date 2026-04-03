package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"hot-ai-backend/internal/models"
)

// StartScheduler 启动定时任务调度器
func StartScheduler(ctx context.Context, c CrawlerConf, db *gorm.DB) {
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
			// TODO: 从数据库读取抓取源并执行抓取
			ExecuteFetchTasks(ctx, c, db)
		}
	}
}

// ExecuteFetchTasks 执行抓取任务
func ExecuteFetchTasks(ctx context.Context, c CrawlerConf, db *gorm.DB) {
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
			
			executeSingleFetchTask(ctx, c, db, src)
		}(source)
	}

	// 等待所有任务完成
	wg.Wait()
	
	elapsed := time.Since(startTime)
	logx.Infof("抓取任务完成，耗时：%v", elapsed)
}

// executeSingleFetchTask 执行单个抓取任务
func executeSingleFetchTask(ctx context.Context, c CrawlerConf, db *gorm.DB, source models.CrawlerSource) {
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
	
	// 3. 处理抓取结果
	itemsFetched, err := processFetchedContent(ctx, c, source, body)
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
		timeout = 30 * time.Second
	}
	
	client := &http.Client{
		Timeout: timeout,
	}
	
	req, err := http.NewRequestWithContext(ctx, string(source.FetchMethod), source.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败：%w", err)
	}
	
	// 设置默认 User-Agent
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
	
	return client.Do(req)
}

// processFetchedContent 处理抓取到的内容
func processFetchedContent(ctx context.Context, c CrawlerConf, source models.CrawlerSource, content []byte) (int, error) {
	// TODO: 根据 parse_rules 解析内容
	// 这里先简单处理，将原始数据推送到 Redis Stream
	
	redisClient := InitRedis(c)
	if redisClient == nil {
		return 0, fmt.Errorf("Redis 初始化失败")
	}
	defer redisClient.Close()
	
	// 推送到 Redis Stream
	streamKey := "crawler:articles"
	values := map[string]interface{}{
		"source_id":   source.ID,
		"source_name": source.Name,
		"source_url":  source.URL,
		"category":    source.Category,
		"content":     string(content),
		"fetched_at":  time.Now().Format(time.RFC3339),
	}
	
	_, err := redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		Values: values,
	}).Result()
	
	if err != nil {
		return 0, fmt.Errorf("推送到 Redis Stream 失败：%w", err)
	}
	
	logx.Infof("已将抓取内容推送到 Redis Stream: %s", source.Name)
	return 1, nil // 暂时返回 1 条
}

// updateCrawlerSource 更新抓取源状态
func updateCrawlerSource(db *gorm.DB, sourceID string, fetchLog *models.CrawlerFetchLog) {
	updates := map[string]interface{}{
		"last_fetch_at":      gorm.Expr("NOW()"),
		"next_fetch_at":      gorm.Expr("DATE_ADD(NOW(), INTERVAL (SELECT fetch_interval FROM crawler_sources WHERE id = ?) SECOND)", sourceID),
		"last_status_code":   fetchLog.StatusCode,
		"total_fetches":      gorm.Expr("total_fetches + 1"),
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
