package crawler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"hot-ai-backend/internal/models"
)

// CreateCrawlerSource 创建抓取源配置
func CreateCrawlerSource(ctx context.Context, db *gorm.DB, source *models.CrawlerSource) error {
	result := db.Create(source)
	if result.Error != nil {
		logx.Errorf("创建抓取源失败: %v", result.Error)
		return result.Error
	}
	logx.Infof("创建抓取源成功: %s", source.Name)
	return nil
}

// GetAllCrawlerSources 获取所有抓取源
func GetAllCrawlerSources(db *gorm.DB) ([]models.CrawlerSource, error) {
	var sources []models.CrawlerSource
	err := db.Find(&sources).Error
	if err != nil {
		logx.Errorf("获取抓取源列表失败: %v", err)
		return nil, err
	}
	return sources, nil
}

// GetActiveCrawlerSources 获取活跃的抓取源
func GetActiveCrawlerSources(db *gorm.DB, limit int) ([]models.CrawlerSource, error) {
	var sources []models.CrawlerSource
	err := db.Where("status = ?", models.CrawlerSourceStatusActive).
		Order("priority DESC, id").
		Limit(limit).
		Find(&sources).Error
	if err != nil {
		logx.Errorf("获取活跃抓取源失败: %v", err)
		return nil, err
	}
	return sources, nil
}

// GetCrawlerSourceByID 根据ID获取抓取源
func GetCrawlerSourceByID(db *gorm.DB, id string) (*models.CrawlerSource, error) {
	var source models.CrawlerSource
	err := db.Where("id = ?", id).First(&source).Error
	if err != nil {
		logx.Errorf("获取抓取源失败 (ID: %s): %v", id, err)
		return nil, err
	}
	return &source, nil
}

// UpdateCrawlerSourceStatus 更新抓取源状态
func UpdateCrawlerSourceStatus(db *gorm.DB, sourceID string, status models.CrawlerSourceStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	err := db.Model(&models.CrawlerSource{}).Where("id = ?", sourceID).Updates(updates).Error
	if err != nil {
		logx.Errorf("更新抓取源状态失败 (ID: %s): %v", sourceID, err)
		return err
	}
	return nil
}

// UpdateCrawlerSourceLastFetch 更新抓取源的最后一次抓取时间
func UpdateCrawlerSourceLastFetch(db *gorm.DB, sourceID string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"last_fetch_at": now,
	}
	err := db.Model(&models.CrawlerSource{}).Where("id = ?", sourceID).Updates(updates).Error
	if err != nil {
		logx.Errorf("更新抓取源抓取时间失败 (ID: %s): %v", sourceID, err)
		return err
	}
	return nil
}

// UpdateCrawlerSourceStats 更新抓取源统计信息
func UpdateCrawlerSourceStats(db *gorm.DB, sourceID string, fetchLog *models.CrawlerFetchLog) error {
	updates := map[string]interface{}{
		"last_status_code":   fetchLog.StatusCode,
		"last_fetch_at":      gorm.Expr("NOW()"),
		"total_fetches":      gorm.Expr("total_fetches + 1"),
		"consecutive_failures": gorm.Expr("consecutive_failures + 1"),
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
		updates["last_error_message"] = fetchLog.ErrorMessage

		// 如果连续失败超过最大重试次数，暂停抓取源
		var source models.CrawlerSource
		if db.First(&source, sourceID).Error == nil {
			if source.ConsecutiveFailures >= source.MaxRetries {
				updates["status"] = models.CrawlerSourceStatusError
				logx.Errorf("抓取源 %s 连续失败 %d 次，暂停抓取", source.Name, source.ConsecutiveFailures)
			}
		}
	}

	err := db.Model(&models.CrawlerSource{}).Where("id = ?", sourceID).Updates(updates).Error
	if err != nil {
		logx.Errorf("更新抓取源统计失败 (ID: %s): %v", sourceID, err)
		return err
	}
	return nil
}

// CreateCrawlerFetchLog 创建抓取日志
func CreateCrawlerFetchLog(db *gorm.DB, fetchLog *models.CrawlerFetchLog) error {
	result := db.Create(fetchLog)
	if result.Error != nil {
		logx.Errorf("创建抓取日志失败: %v", result.Error)
		return result.Error
	}
	return nil
}

// CreateCrawlerFetchLogWithContent 创建抓取日志并包含原始内容
func CreateCrawlerFetchLogWithContent(db *gorm.DB, fetchLog *models.CrawlerFetchLog) error {
	result := db.Create(fetchLog)
	if result.Error != nil {
		logx.Errorf("创建抓取日志失败: %v", result.Error)
		return result.Error
	}
	return nil
}

// GetCrawlerFetchLogs 获取抓取日志列表
func GetCrawlerFetchLogs(db *gorm.DB, limit int) ([]models.CrawlerFetchLog, error) {
	var logs []models.CrawlerFetchLog
	err := db.Order("created_at DESC").Limit(limit).Find(&logs).Error
	if err != nil {
		logx.Errorf("获取抓取日志列表失败: %v", err)
		return nil, err
	}
	return logs, nil
}

// GetCrawlerFetchLogByID 根据ID获取抓取日志
func GetCrawlerFetchLogByID(db *gorm.DB, id string) (*models.CrawlerFetchLog, error) {
	var log models.CrawlerFetchLog
	err := db.Where("id = ?", id).First(&log).Error
	if err != nil {
		logx.Errorf("获取抓取日志失败 (ID: %s): %v", id, err)
		return nil, err
	}
	return &log, nil
}

// GetFetchLogsBySourceID 获取指定抓取源的日志
func GetFetchLogsBySourceID(db *gorm.DB, sourceID string, limit int) ([]models.CrawlerFetchLog, error) {
	var logs []models.CrawlerFetchLog
	err := db.Where("source_id = ?", sourceID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	if err != nil {
		logx.Errorf("获取抓取源日志失败 (SourceID: %s): %v", sourceID, err)
		return nil, err
	}
	return logs, nil
}

// GetFetchLogsByStatus 获取指定状态的抓取日志
func GetFetchLogsByStatus(db *gorm.DB, status string, limit int) ([]models.CrawlerFetchLog, error) {
	var logs []models.CrawlerFetchLog
	err := db.Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	if err != nil {
		logx.Errorf("获取抓取日志失败 (Status: %s): %v", status, err)
		return nil, err
	}
	return logs, nil
}

// CleanOldFetchLogs 清理旧的抓取日志
func CleanOldFetchLogs(db *gorm.DB, days int) error {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	err := db.Where("created_at < ?", cutoffTime).Delete(&models.CrawlerFetchLog{}).Error
	if err != nil {
		logx.Errorf("清理旧抓取日志失败: %v", err)
		return err
	}
	logx.Infof("已清理 %d 天前的抓取日志", days)
	return nil
}

// UpdateCrawlerSourceNextFetchTime 更新抓取源的下次抓取时间
func UpdateCrawlerSourceNextFetchTime(db *gorm.DB, sourceID string) error {
	var source models.CrawlerSource
	err := db.Where("id = ?", sourceID).First(&source).Error
	if err != nil {
		logx.Errorf("获取抓取源失败 (ID: %s): %v", sourceID, err)
		return err
	}

	// 设置下次抓取时间
	nextFetchAt := time.Now().Add(time.Duration(source.FetchInterval) * time.Second)
	updates := map[string]interface{}{
		"next_fetch_at": nextFetchAt,
	}

	err = db.Model(&models.CrawlerSource{}).Where("id = ?", sourceID).Updates(updates).Error
	if err != nil {
		logx.Errorf("更新抓取源下次抓取时间失败 (ID: %s): %v", sourceID, err)
		return err
	}
	logx.Infof("更新抓取源 %s 下次抓取时间: %s", source.Name, nextFetchAt.Format(time.RFC3339))
	return nil
}

// SaveCrawlerSourceToDatabase 保存抓取源配置到数据库
func SaveCrawlerSourceToDatabase(db *gorm.DB, source models.CrawlerSource) error {
	// 检查是否已存在
	var existing models.CrawlerSource
	err := db.Where("id = ?", source.ID).First(&existing).Error

	if err == nil {
		// 已存在，更新
		source.UpdatedAt = time.Now()
		err = db.Save(&source).Error
		if err != nil {
			logx.Errorf("更新抓取源失败 (ID: %s): %v", source.ID, err)
			return err
		}
		logx.Infof("更新抓取源成功: %s", source.Name)
	} else if err == gorm.ErrRecordNotFound {
		// 不存在，创建
		source.ID = uuid.New().String()
		source.CreatedAt = time.Now()
		source.UpdatedAt = time.Now()
		err = db.Create(&source).Error
		if err != nil {
			logx.Errorf("创建抓取源失败: %v", err)
			return err
		}
		logx.Infof("创建抓取源成功: %s", source.Name)
	} else {
		logx.Errorf("查询抓取源失败 (ID: %s): %v", source.ID, err)
		return err
	}
	return nil
}

// BatchCreateCrawlerSources 批量创建抓取源
func BatchCreateCrawlerSources(db *gorm.DB, sources []models.CrawlerSource) error {
	now := time.Now()
	for i := range sources {
		sources[i].ID = uuid.New().String()
		sources[i].CreatedAt = now
		sources[i].UpdatedAt = now
	}

	result := db.Create(&sources)
	if result.Error != nil {
		logx.Errorf("批量创建抓取源失败: %v", result.Error)
		return result.Error
	}
	logx.Infof("批量创建抓取源成功，数量: %d", len(sources))
	return nil
}
