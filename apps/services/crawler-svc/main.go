package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/crawler-svc.yaml", "the config file")

func main() {
	flag.Parse()

	var c conf.Config
	if err := conf.Load(*configFile, &c); err != nil {
		fmt.Fprintf(os.Stderr, "load config error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Starting crawler-svc...")

	// TODO: 初始化采集服务
	// 1. 启动定时任务调度器
	// 2. 启动 Redis Stream 消费者
	// 3. 初始化 AI 处理模块

	startCrawlerService(c)
}

func startCrawlerService(c conf.Config) {
	// 待实现
	// 1. 从数据库读取抓取源
	// 2. 定时执行抓取任务
	// 3. 内容清洗和去重
	// 4. 发布 ArticleCreated 事件到 Redis Stream
}
