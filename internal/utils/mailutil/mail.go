package mailutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Config 邮件服务配置
type Config struct {
	SMTPServer string
	SMTPPort   int
	Username   string
	Password   string
	FromName   string
	FromEmail  string
}

// EmailVerificationToken 邮箱验证 Token（用于验证链接）
type EmailVerificationToken struct {
	UserID    string
	Email     string
	ExpiresAt time.Time
}

// EmailVerificationCode 邮箱验证码（用于注册）
type EmailVerificationCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
}

// GenerateVerificationToken 生成邮箱验证 Token
func GenerateVerificationToken(userID, email, secret string) (string, error) {
	if userID == "" || email == "" || secret == "" {
		return "", errors.New("invalid parameters")
	}

	// 设置过期时间（24 小时）
	expiresAt := time.Now().Add(24 * time.Hour)

	// 构建 payload
	payload := fmt.Sprintf("%s|%s|%d", userID, email, expiresAt.Unix())

	// 生成签名
	signature := generateHMAC(payload, secret)

	// 组合 token: payload.signature
	token := fmt.Sprintf("%s.%s", payload, signature)

	// Base64 编码
	return base64.URLEncoding.EncodeToString([]byte(token)), nil
}

// ParseVerificationToken 解析并验证邮箱验证 Token
func ParseVerificationToken(tokenString, secret string) (*EmailVerificationToken, error) {
	if tokenString == "" || secret == "" {
		return nil, errors.New("invalid parameters")
	}

	// Base64 解码
	decoded, err := base64.URLEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, errors.New("invalid token format")
	}

	// 分离 payload 和 signature
	parts := splitToken(string(decoded))
	if len(parts) != 2 {
		return nil, errors.New("invalid token structure")
	}

	payload := parts[0]
	signature := parts[1]

	// 验证签名
	expectedSignature := generateHMAC(payload, secret)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, errors.New("invalid token signature")
	}

	// 解析 payload
	payloadParts := splitStringByDelimiter(payload, '|')
	if len(payloadParts) != 3 {
		return nil, errors.New("invalid token payload")
	}

	userID := payloadParts[0]
	email := payloadParts[1]
	expiresAtUnix := payloadParts[2]

	// 解析过期时间
	var expiresAt int64
	_, err = fmt.Sscanf(expiresAtUnix, "%d", &expiresAt)
	if err != nil {
		return nil, errors.New("invalid token expiration")
	}

	// 检查是否过期
	if time.Now().Unix() > expiresAt {
		return nil, errors.New("token has expired")
	}

	return &EmailVerificationToken{
		UserID:    userID,
		Email:     email,
		ExpiresAt: time.Unix(expiresAt, 0),
	}, nil
}

// generateHMAC 生成 HMAC 签名
func generateHMAC(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// splitToken 分割 token
func splitToken(token string) []string {
	for i := len(token) - 1; i >= 0; i-- {
		if token[i] == '.' {
			return []string{token[:i], token[i+1:]}
		}
	}
	return []string{token}
}

// splitStringByDelimiter 按分隔符分割字符串
func splitStringByDelimiter(s string, delimiter rune) []string {
	var parts []string
	var current string

	for _, r := range s {
		if r == delimiter {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	parts = append(parts, current)
	return parts
}

// GenerateVerificationCode 生成 6 位数字验证码
func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
