package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"hot-ai-backend/internal/crawler"
)

var configFile = flag.String("f", "apps\\crawler-svc\\etc\\crawler-svc.yaml", "the config file")

func main() {
	flag.Parse()

	var c crawler.CrawlerConf
	if err := conf.Load(*configFile, &c); err != nil {
		fmt.Fprintf(os.Stderr, "load config error: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	crawler.InitLog(c)

	fmt.Println("Starting crawler-svc...")

	// 创建服务上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动采集服务
	crawler.StartCrawlerService(ctx, c)

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down crawler-svc...")
	cancel()
	time.Sleep(time.Second * 2) // 等待 goroutine 清理
}
