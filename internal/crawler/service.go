package crawler

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// CrawlerConf 采集服务配置
type CrawlerConf struct {
	Name   string `json:"Name,optional,default=crawler-svc"`
	Consul struct {
		Host string `json:",optional"`
		Port int    `json:",default=8500"`
		Key  string `json:",optional"`
	} `json:"Consul,optional"`
	// 数据库配置
	DataSource struct {
		MySQL struct {
			DSN string `json:",optional"`
		}
		MongoDB struct {
			URI      string `json:",optional"`
			Database string `json:",optional"`
		} `json:"MongoDB,optional"`
	} `json:"DataSource,optional"`
	// AI 配置
	AI struct {
		Provider string `json:",optional"`
		APIKey   string `json:",optional"`
		Model    string `json:",optional"`
	} `json:"AI,optional"`
	// 爬虫配置
	Crawler struct {
		MaxConcurrency            int    `json:",default=5"`
		Timeout                   int    `json:",default=30"`
		UserAgent                 string `json:",optional"`
		ConcurrentFetch           int    `json:",default=10"`
		FetchInterval             int    `json:",default=300"`
		MaxRetries                int    `json:",default=3"`
		EnableTraditionalCrawling bool   `json:",default=true"` // 是否启用传统HTML网页爬取
	} `json:"Crawler,optional"`
	// Miniflux 配置（RSS 源）
	Miniflux struct {
		BaseURL string `json:",optional"`
		APIKey  string `json:",optional"`
		Enabled bool   `json:",default=false"`
	} `json:"Miniflux,optional"`
	// AI Server 服务配置
	AIServer struct {
		BaseURL string `json:",optional"`
		Enabled bool   `json:",default=false"`
		Model   string `json:",optional"`
	} `json:"AIServer,optional"`
	// 日志配置
	Log struct {
		Mode     string `json:",optional"`
		Path     string `json:",optional"`
		Level    string `json:",optional"`
		Encoding string `json:",optional"`
	} `json:"Log,optional"`
}

// StartCrawlerService 启动采集服务
func StartCrawlerService(ctx context.Context, c CrawlerConf) {
	logx.Info("初始化采集服务...")

	// 1. 初始化数据库连接
	var db *gorm.DB
	if c.DataSource.MySQL.DSN != "" {
		var err error
		db, err = gorm.Open(mysql.Open(c.DataSource.MySQL.DSN), &gorm.Config{})
		if err != nil {
			logx.Error("连接数据库失败:", err)
			return
		}
		logx.Info("数据库连接成功")

		// 设置全局数据库连接（供 crawler 模块使用）
		SetGlobalDB(db)
	}

	// 2. 初始化 Miniflux 客户端（如果启用）
	var minifluxClient *MinifluxClient
	if c.Miniflux.Enabled && c.Miniflux.BaseURL != "" && c.Miniflux.APIKey != "" {
		minifluxClient = NewMinifluxClient(c.Miniflux.BaseURL, c.Miniflux.APIKey)
		logx.Infof("Miniflux 客户端初始化成功: %s", c.Miniflux.BaseURL)

		// 测试连接
		//feeds, err := minifluxClient.GetFeeds()
		//if err != nil {
		//	logx.Errorf("Miniflux 连接测试失败: %v", err)
		//} else {
		//	logx.Infof("Miniflux 中有 %d 个订阅源", len(feeds))
		//}
	} else {
		logx.Info("Miniflux 未启用，将使用传统爬虫模式")
	}

	// 3. 初始化 AI Server 客户端（用于翻译等功能）
	var translateClient *TranslateClient
	if c.AIServer.Enabled && c.AIServer.BaseURL != "" {
		translateClient = NewTranslateClient(c.AIServer.BaseURL)
		if c.AIServer.Model != "" {
			translateClient.model = c.AIServer.Model
		}
		logx.Infof("AI Server 客户端初始化成功: %s (模型: %s)", c.AIServer.BaseURL, translateClient.model)
	} else {
		logx.Info("AI Server 服务未启用，文章将不会被翻译")
	}

	// 4. 启动定时任务调度器
	go StartScheduler(ctx, c, db, minifluxClient, translateClient)

	logx.Info("采集服务启动完成（支持 RSS + 传统爬虫）")
}

// InitLog 初始化日志配置
func InitLog(c CrawlerConf) {
	if c.Log.Mode == "file" {
		logx.MustSetup(logx.LogConf{
			Mode:     c.Log.Mode,
			Path:     c.Log.Path,
			Level:    c.Log.Level,
			Encoding: c.Log.Encoding,
		})
	} else {
		logx.MustSetup(logx.LogConf{
			Mode:  "console",
			Level: "info",
		})
	}
}
