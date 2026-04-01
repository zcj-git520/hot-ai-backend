package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeNotFound    = errors.New("验证码不存在或已过期")
	ErrCodeInvalid     = errors.New("验证码错误")
	ErrCodeExpired     = errors.New("验证码已过期")
	ErrSendTooFrequent = errors.New("发送过于频繁")
)

// EmailCache 邮箱验证码缓存
type EmailCache struct {
	client *redis.Client
	prefix string
}

// NewEmailCache 创建邮箱缓存实例
func NewEmailCache(client *redis.Client) *EmailCache {
	return &EmailCache{
		client: client,
		prefix: "email:verify:",
	}
}

// StoreVerificationCode 存储验证码
func (c *EmailCache) StoreVerificationCode(ctx context.Context, email, code string, expire time.Duration) error {
	key := c.getVerifyCodeKey(email)
	
	data := map[string]interface{}{
		"code":      code,
		"created_at": time.Now().Unix(),
	}
	
	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal code data: %w", err)
	}
	
	return c.client.Set(ctx, key, value, expire).Err()
}

// VerifyCode 验证邮箱验证码
func (c *EmailCache) VerifyCode(ctx context.Context, email, inputCode string) error {
	key := c.getVerifyCodeKey(email)
	
	// 获取存储的验证码
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCodeNotFound
		}
		return fmt.Errorf("failed to get code: %w", err)
	}
	
	var storedData map[string]interface{}
	if err := json.Unmarshal(val, &storedData); err != nil {
		return fmt.Errorf("failed to unmarshal code data: %w", err)
	}
	
	storedCode, ok := storedData["code"].(string)
	if !ok {
		return ErrCodeInvalid
	}
	
	// 比较验证码
	if storedCode != inputCode {
		return ErrCodeInvalid
	}
	
	// 验证成功后删除验证码（防止重复使用）
	_ = c.client.Del(ctx, key)
	
	return nil
}

// CheckSendFrequency 检查发送频率（1 分钟内最多发送 1 次，1 小时内最多 5 次）
func (c *EmailCache) CheckSendFrequency(ctx context.Context, email string) error {
	// 检查 1 分钟内的发送次数
	minuteKey := c.getMinuteLimitKey(email)
	minuteCount, err := c.client.Get(ctx, minuteKey).Int64()
	if err == nil && minuteCount >= 1 {
		return ErrSendTooFrequent
	}
	
	// 检查 1 小时内的发送次数
	hourKey := c.getHourLimitKey(email)
	hourCount, err := c.client.Get(ctx, hourKey).Int64()
	if err == nil && hourCount >= 5 {
		return errors.New("操作过于频繁，请 1 小时后再试")
	}
	
	return nil
}

// RecordSendAttempt 记录发送尝试
func (c *EmailCache) RecordSendAttempt(ctx context.Context, email string) error {
	// 增加 1 分钟计数
	minuteKey := c.getMinuteLimitKey(email)
	pipe := c.client.Pipeline()
	pipe.Incr(ctx, minuteKey)
	pipe.Expire(ctx, minuteKey, time.Minute)
	
	// 增加 1 小时计数
	hourKey := c.getHourLimitKey(email)
	pipe.Incr(ctx, hourKey)
	pipe.Expire(ctx, hourKey, time.Hour)
	
	_, err := pipe.Exec(ctx)
	return err
}

// DeleteVerificationCode 删除验证码
func (c *EmailCache) DeleteVerificationCode(ctx context.Context, email string) error {
	key := c.getVerifyCodeKey(email)
	return c.client.Del(ctx, key).Err()
}

func (c *EmailCache) getVerifyCodeKey(email string) string {
	return fmt.Sprintf("%scode:%s", c.prefix, email)
}

func (c *EmailCache) getMinuteLimitKey(email string) string {
	return fmt.Sprintf("%slimit:minute:%s", c.prefix, email)
}

func (c *EmailCache) getHourLimitKey(email string) string {
	return fmt.Sprintf("%slimit:hour:%s", c.prefix, email)
}
