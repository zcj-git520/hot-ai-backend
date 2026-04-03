package models

import (
	"encoding/json"
	"time"
)

// CrawlerSourceType 抓取源类型
type CrawlerSourceType string

const (
	CrawlerSourceTypeRSS     CrawlerSourceType = "rss"
	CrawlerSourceTypeSitemap CrawlerSourceType = "sitemap"
	CrawlerSourceTypeAPI     CrawlerSourceType = "api"
	CrawlerSourceTypeHTML    CrawlerSourceType = "html"
)

// FetchMethod 请求方法
type FetchMethod string

const (
	FetchMethodGET  FetchMethod = "get"
	FetchMethodPOST FetchMethod = "post"
)

// ParserType 解析器类型
type ParserType string

const (
	ParserTypeJSON  ParserType = "json"
	ParserTypeXPATH ParserType = "xpath"
	ParserTypeCSS   ParserType = "css"
	ParserTypeRegex ParserType = "regex"
	ParserTypeHTML  ParserType = "html"
)

// CrawlerSourceStatus 抓取源状态
type CrawlerSourceStatus string

const (
	CrawlerSourceStatusActive  CrawlerSourceStatus = "active"
	CrawlerSourceStatusInactive CrawlerSourceStatus = "inactive"
	CrawlerSourceStatusError   CrawlerSourceStatus = "error"
	CrawlerSourceStatusPaused  CrawlerSourceStatus = "paused"
)

// FetchLogStatus 抓取日志状态
type FetchLogStatus string

const (
	FetchLogStatusSuccess FetchLogStatus = "success"
	FetchLogStatusFailed  FetchLogStatus = "failed"
	FetchLogStatusTimeout FetchLogStatus = "timeout"
	FetchLogStatusSkipped FetchLogStatus = "skipped"
)

// CrawlerSource 抓取源配置
type CrawlerSource struct {
	ID                string              `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	Name              string              `gorm:"column:name;type:varchar(100);not null" json:"name"`
	URL               string              `gorm:"column:url;type:varchar(500);not null" json:"url"`
	SourceType        CrawlerSourceType   `gorm:"column:source_type;type:enum('rss','sitemap','api','html');not null;default:'html'" json:"source_type"`
	Category          string              `gorm:"column:category;type:varchar(50)" json:"category"`
	FetchMethod       FetchMethod         `gorm:"column:fetch_method;type:enum('get','post');not null;default:'get'" json:"fetch_method"`
	RequestHeaders    json.RawMessage     `gorm:"column:request_headers;type:json" json:"request_headers,omitempty"`
	RequestBody       string              `gorm:"column:request_body;type:text" json:"request_body,omitempty"`
	
	// 解析配置
	ParserType        ParserType          `gorm:"column:parser_type;type:enum('json','xpath','css','regex','html');not null;default:'html'" json:"parser_type"`
	ParseRules        json.RawMessage     `gorm:"column:parse_rules;type:json" json:"parse_rules,omitempty"`
	
	// 调度配置
	FetchInterval     int                 `gorm:"column:fetch_interval;default:300" json:"fetch_interval"`
	Timeout           int                 `gorm:"column:timeout;default:30" json:"timeout"`
	MaxRetries        int                 `gorm:"column:max_retries;default:3" json:"max_retries"`
	Priority          int                 `gorm:"column:priority;default:1" json:"priority"`
	
	// 状态管理
	Status            CrawlerSourceStatus `gorm:"column:status;type:enum('active','inactive','error','paused');not null;default:'active'" json:"status"`
	LastFetchAt       *time.Time          `gorm:"column:last_fetch_at" json:"last_fetch_at"`
	NextFetchAt       *time.Time          `gorm:"column:next_fetch_at" json:"next_fetch_at"`
	LastStatusCode    int                 `gorm:"column:last_status_code" json:"last_status_code"`
	LastErrorMessage  string              `gorm:"column:last_error_message;type:text" json:"last_error_message,omitempty"`
	ConsecutiveFailures int               `gorm:"column:consecutive_failures;default:0" json:"consecutive_failures"`
	
	// 统计信息
	TotalFetches      int                 `gorm:"column:total_fetches;default:0" json:"total_fetches"`
	SuccessfulFetches int                 `gorm:"column:successful_fetches;default:0" json:"successful_fetches"`
	FailedFetches     int                 `gorm:"column:failed_fetches;default:0" json:"failed_fetches"`
	TotalItems        int                 `gorm:"column:total_items;default:0" json:"total_items"`
	
	// 元数据
	Description       string              `gorm:"column:description;type:text" json:"description,omitempty"`
	Tags              json.RawMessage     `gorm:"column:tags;type:json" json:"tags,omitempty"`
	CreatedBy         string              `gorm:"column:created_by;type:varchar(36)" json:"created_by,omitempty"`
	CreatedAt         time.Time           `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time           `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (CrawlerSource) TableName() string {
	return "crawler_sources"
}

// IsReadyToFetch 检查是否到了抓取时间
func (s *CrawlerSource) IsReadyToFetch() bool {
	if s.Status != CrawlerSourceStatusActive {
		return false
	}
	
	if s.NextFetchAt == nil {
		return true
	}
	
	return time.Now().After(*s.NextFetchAt)
}

// CrawlerFetchLog 抓取日志
type CrawlerFetchLog struct {
	ID                 string         `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	SourceID           string         `gorm:"column:source_id;type:varchar(36);not null;index:idx_source_id" json:"source_id"`
	FetchStartedAt     time.Time      `gorm:"column:fetch_started_at;not null;index:idx_fetch_time" json:"fetch_started_at"`
	FetchCompletedAt   *time.Time     `gorm:"column:fetch_completed_at" json:"fetch_completed_at"`
	DurationMs         int            `gorm:"column:duration_ms" json:"duration_ms"`
	StatusCode         int            `gorm:"column:status_code" json:"status_code"`
	ResponseSize       int            `gorm:"column:response_size" json:"response_size"`
	ItemsFetched       int            `gorm:"column:items_fetched;default:0" json:"items_fetched"`
	Status             FetchLogStatus `gorm:"column:status;type:enum('success','failed','timeout','skipped');not null" json:"status"`
	ErrorMessage       string         `gorm:"column:error_message;type:text" json:"error_message,omitempty"`
	RawResponse        string         `gorm:"column:raw_response;type:longtext" json:"raw_response,omitempty"`
	CreatedAt          time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	
	// 关联
	Source CrawlerSource `gorm:"foreignKey:SourceID" json:"-"`
}

// TableName 指定表名
func (CrawlerFetchLog) TableName() string {
	return "crawler_fetch_logs"
}
