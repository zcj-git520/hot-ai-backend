package crawler

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	
	// 1. 解析文章内容
	sourceID, _ := msg.Values["source_id"].(string)
	sourceName, _ := msg.Values["source_name"].(string)
	sourceURL, _ := msg.Values["source_url"].(string)
	category, _ := msg.Values["category"].(string)
	content, _ := msg.Values["content"].(string)
	fetchedAt, _ := msg.Values["fetched_at"].(string)
	
	if content == "" {
		logx.Error("内容为空，跳过处理")
		return
	}
	
	// 2. 清洗和去重（这里简化处理）
	// TODO: 实现更复杂的内容清洗逻辑
	cleanedContent := cleanContent(content)
	
	// 3. 构建文章数据
	articleData := map[string]interface{}{
		"id":          uuid.New().String(),
		"source_id":   sourceID,
		"source_name": sourceName,
		"source_url":  sourceURL,
		"category":    category,
		"title":       extractTitle(cleanedContent), // TODO: 实现标题提取
		"content":     cleanedContent,
		"url":         sourceURL,
		"fetched_at":  fetchedAt,
		"status":      "pending_process", // 待处理
	}
	
	// 4. 推送到 AI 处理队列（如果配置了 AI）
	if c.AI.APIKey != "" && c.AI.APIKey != "your-api-key-here" {
		pushToAIQueue(ctx, c, articleData)
	} else {
		// 5. 直接存储到 MongoDB（如果没有 AI 处理）
		storeToMongoDB(ctx, c, articleData)
	}
	
	logx.Infof("文章消息处理完成：%s", sourceName)
}

// cleanContent 清洗内容
func cleanContent(content string) string {
	// TODO: 实现 HTML 标签清理、空白字符处理等
	return content
}

// extractTitle 提取标题
func extractTitle(content string) string {
	// TODO: 根据内容提取标题
	return "未命名文章"
}

// pushToAIQueue 推送到 AI 处理队列
func pushToAIQueue(ctx context.Context, c CrawlerConf, articleData map[string]interface{}) {
	redisClient := InitRedis(c)
	if redisClient == nil {
		logx.Error("Redis 初始化失败，无法推送 AI 队列")
		return
	}
	defer redisClient.Close()
	
	aiStreamKey := "ai:articles_to_process"
	_, err := redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: aiStreamKey,
		Values: articleData,
	}).Result()
	
	if err != nil {
		logx.Error("推送 AI 队列失败:", err)
	} else {
		logx.Info("已推送至 AI 处理队列")
	}
}

// storeToMongoDB 存储到 MongoDB
func storeToMongoDB(ctx context.Context, c CrawlerConf, articleData map[string]interface{}) {
	// TODO: 实现 MongoDB 存储逻辑
	logx.Info("存储到 MongoDB（待实现）")
}
