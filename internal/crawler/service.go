package crawler

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// CrawlerConf 采集服务配置
type CrawlerConf struct {
	Name    string `json:"Name,optional,default=crawler-svc"`
	Consul  struct {
		Host string `json:",optional"`
		Port int    `json:",default=8500"`
		Key  string `json:",optional"`
	} `json:"Consul,optional"`
	// Redis 配置
	Redis struct {
		Host     string `json:",optional"`
		Port     int    `json:",default=6379"`
		Password string `json:",optional"`
		DB       int    `json:",default=0"`
		Type     string `json:",optional"`
	} `json:"Redis"`
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
		MaxConcurrency  int `json:",default=5"`
		Timeout         int `json:",default=30"`
		UserAgent       string `json:",optional"`
		ConcurrentFetch int `json:",default=10"`
		FetchInterval   int `json:",default=300"`
		MaxRetries      int `json:",default=3"`
	} `json:"Crawler,optional"`
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
	}

	// 2. 初始化 Redis 客户端
	redisClient := InitRedis(c)
	if redisClient == nil {
		logx.Error("Redis 初始化失败")
		return
	}

	// 3. 启动定时任务调度器
	go StartScheduler(ctx, c, db)

	// 4. 启动 Redis Stream 消费者
	go StartRedisStreamConsumer(ctx, c, redisClient, db)

	// 5. 初始化 AI 处理模块
	if c.AI.APIKey != "" && c.AI.APIKey != "your-api-key-here" {
		go StartAIProcessor(ctx, c, redisClient)
	}

	logx.Info("采集服务启动完成")
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
