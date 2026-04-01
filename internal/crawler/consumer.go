package crawler

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// startRedisStreamConsumer 启动 Redis Stream 消费者
func StartRedisStreamConsumer(ctx context.Context, c CrawlerConf, client *redis.Client) {
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
				ProcessArticleMessage(ctx, c, msg)
				// 确认消息
				client.XAck(ctx, streamKey, groupKey, msg.ID)
			}
		}
	}
}

// ProcessArticleMessage 处理文章消息
func ProcessArticleMessage(ctx context.Context, c CrawlerConf, msg redis.XMessage) {
	logx.Infof("处理文章消息：%+v", msg.Values)
	// TODO: 实现文章处理逻辑
	// 1. 解析文章内容
	// 2. 清洗和去重
	// 3. 存储到 MongoDB
}
