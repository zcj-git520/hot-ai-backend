// Package access 权限决策与内容裁剪
//
// 业务服务（content / learning-path / profession / tool）统一调用本包的 helper，
// 避免每个服务重复实现「是否放行 / 截断多少 / 占位文案」。
package access

import (
	"sync/atomic"
	"time"

	"hot-ai-backend/internal/models"
)

// 默认游客预览字数（按 rune 计）
const GuestPreviewChars = 500

// paywallEnabled 是否启用付费墙。true = 按 level 拦截；false = 所有人可读。
// 默认 false（全员可读），由各服务 main 启动时根据 YAML 中的 Paywall.Enabled 写入。
// 用 atomic.Bool 是因为读多写少（启动写一次，请求读），并发安全。
var paywallEnabled atomic.Bool

// SetPaywallEnabled 设置是否启用付费墙。建议在 main 启动时调用一次。
func SetPaywallEnabled(enabled bool) {
	paywallEnabled.Store(enabled)
}

// PaywallEnabled 返回当前付费墙是否启用。
func PaywallEnabled() bool {
	return paywallEnabled.Load()
}

// LevelName 级别名称（给前端 / 日志用）
func LevelName(level int) string {
	switch level {
	case 0:
		return "游客"
	case 1:
		return "普通用户"
	case 2:
		return "会员"
	default:
		return "未知"
	}
}

// Decision 决策结果
type Decision struct {
	Allow  bool
	Reason string
}

// Decide 是否允许访问完整内容
// userLevel: 请求方身份级别（0/1/2）
// contentLevel: 内容要求的级别（0/1/2）
//
// 付费墙关闭时直接放行（无论 userLevel/contentLevel），让所有用户可读完整内容。
// 付费墙开启时才走 level 比对。
func Decide(userLevel int, contentLevel int) Decision {
	if !paywallEnabled.Load() {
		return Decision{Allow: true, Reason: "paywall_disabled"}
	}
	if userLevel >= contentLevel {
		return Decision{Allow: true}
	}
	return Decision{Allow: false, Reason: "insufficient_level"}
}

// TruncateContent 按 rune 截断内容（避免半字汉字）
// 返回 (preview, truncated)
func TruncateContent(content string, maxChars int) (string, bool) {
	if maxChars <= 0 {
		return "", true
	}
	r := []rune(content)
	if len(r) <= maxChars {
		return content, false
	}
	return string(r[:maxChars]) + "…", true
}

// LockedContent 锁定时返回给前端的占位
type LockedContent struct {
	IsLocked          bool   `json:"is_locked"`
	RequiredLevel     int    `json:"required_level"`
	RequiredLevelName string `json:"required_level_name"`
	Title             string `json:"title"`
	Message           string `json:"message"`
	CTA               string `json:"cta"`
}

// LockedPlaceholder 生成锁定占位
// contentType: 内容类型描述（"文章"/"章节"/"职业"/"工具"）
// requiredLevel: 该内容要求的最低级别
func LockedPlaceholder(contentType string, requiredLevel int) LockedContent {
	cta := "登录后阅读"
	if requiredLevel >= 2 {
		cta = "升级会员"
	}
	return LockedContent{
		IsLocked:          true,
		RequiredLevel:     requiredLevel,
		RequiredLevelName: LevelName(requiredLevel),
		Title:             "此" + contentType + "为「" + LevelName(requiredLevel) + "专享」",
		Message:           "本" + contentType + "需要 " + LevelName(requiredLevel) + " 身份才能阅读全文。",
		CTA:               cta,
	}
}

// ComputeLevel 根据用户实时身份计算访问级别
// 规则：
//   - banned 用户视为 level=1 (不允许享受会员)
//   - 会员且未到期 → level=2
//   - 登录用户 → level=1
//   - 未登录 → level=0 (理论上不会传 user，调用方需先判断)
func ComputeLevel(u *models.User) int {
	if u == nil {
		return 0
	}
	if u.Status == models.UserStatusBanned {
		return 1
	}
	if u.IsMember && (u.MemberExpireAt == nil || u.MemberExpireAt.After(time.Now())) {
		return 2
	}
	return 1
}
