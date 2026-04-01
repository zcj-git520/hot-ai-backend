package crawler

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// InitRedis 初始化 Redis 客户端
func InitRedis(c CrawlerConf) *redis.Client {
	addr := c.Redis.Host
	if c.Redis.Port != 0 && addr != "" {
		// 如果 Host 已经包含端口，则不使用 Port
		if !ContainsPort(addr) {
			addr = fmt.Sprintf("%s:%d", addr, c.Redis.Port)
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logx.Error("Redis 连接失败:", err)
		return nil
	}

	logx.Info("Redis 连接成功:", addr)
	return client
}

// ContainsPort 检查地址是否已包含端口
func ContainsPort(addr string) bool {
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return true
		}
	}
	return false
}
