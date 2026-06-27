package recruit

import (
	"encoding/json"
	"time"
)

// Platform 招聘平台
type Platform string

const (
	PlatformBoss    Platform = "boss"
	PlatformZhilian Platform = "zhilian"
	PlatformLiepin  Platform = "liepin"
)

// RawJob 原始岗位
type RawJob struct {
	ID            uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Platform      Platform  `gorm:"column:platform;type:enum('boss','zhilian','liepin')"`
	PlatformJobID string    `gorm:"column:platform_job_id"`
	Title         string    `gorm:"column:title"`
	Company       string    `gorm:"column:company"`
	City          string    `gorm:"column:city"`
	SalaryMin     *int      `gorm:"column:salary_min"`
	SalaryMax     *int      `gorm:"column:salary_max"`
	Description   string    `gorm:"column:description;type:text"`
	Skills        []string  `gorm:"column:skills;type:json"`
	Industry      string    `gorm:"column:industry"`
	URL           string    `gorm:"column:url"`
	CrawledAt     time.Time `gorm:"column:crawled_at"`
}

func (RawJob) TableName() string { return "recruit_raw_jobs" }

// NormalizedJob 归一后岗位
type NormalizedJob struct {
	ID              uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	RawJobID        uint64    `gorm:"column:raw_job_id"`
	ProfessionID    *uint     `gorm:"column:profession_id"`
	MatchMethod     string    `gorm:"column:match_method;type:enum('llm','keyword','manual')"`
	MatchConfidence *float64  `gorm:"column:match_confidence"`
	AIKeywordsCount int       `gorm:"column:ai_keywords_count"`
	AIKeywordsTotal int       `gorm:"column:ai_keywords_total"`
	AIKeywordHits   []string  `gorm:"column:ai_keyword_hits;type:json"`
	NormalizedAt    time.Time `gorm:"column:normalized_at"`
}

func (NormalizedJob) TableName() string { return "recruit_normalized_jobs" }

// DailyMetrics 每日 4 维指标
type DailyMetrics struct {
	ID                   uint64          `gorm:"column:id;primaryKey;autoIncrement"`
	ProfessionID         uint            `gorm:"column:profession_id"`
	MetricDate           time.Time       `gorm:"column:metric_date"`
	JobCount             int             `gorm:"column:job_count"`
	JobCountPrev30d      int             `gorm:"column:job_count_prev_30d"`
	SalaryMedian         *float64        `gorm:"column:salary_median"`
	SalaryMedianPrev90d  *float64        `gorm:"column:salary_median_prev_90d"`
	AIPenetrationRate    *float64        `gorm:"column:ai_penetration_rate"`
	GeoDistribution      json.RawMessage `gorm:"column:geo_distribution;type:json"`
	IndustryDistribution json.RawMessage `gorm:"column:industry_distribution;type:json"`
	SampleSize           int             `gorm:"column:sample_size"`
}

func (DailyMetrics) TableName() string { return "recruit_daily_metrics" }

// Keyword 关键词
type Keyword struct {
	ID           uint   `gorm:"column:id;primaryKey;autoIncrement"`
	ProfessionID uint   `gorm:"column:profession_id"`
	Keyword      string `gorm:"column:keyword"`
	Weight       int    `gorm:"column:weight"`
}

func (Keyword) TableName() string { return "recruit_keywords" }

// CrawlLog 抓取日志
type CrawlLog struct {
	ID           uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	Platform     Platform   `gorm:"column:platform;type:enum('boss','zhilian','liepin')"`
	StartedAt    time.Time  `gorm:"column:started_at"`
	FinishedAt   *time.Time `gorm:"column:finished_at"`
	Status       string     `gorm:"column:status;type:enum('success','partial','failed')"`
	JobsFetched  int        `gorm:"column:jobs_fetched"`
	JobsNew      int        `gorm:"column:jobs_new"`
	ErrorMessage string     `gorm:"column:error_message;type:text"`
}

func (CrawlLog) TableName() string { return "recruit_crawl_log" }

// Config 配置
type Config struct {
	ID          uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	ConfigKey   string    `gorm:"column:config_key"`
	ConfigValue float64   `gorm:"column:config_value"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Config) TableName() string { return "recruit_config" }
