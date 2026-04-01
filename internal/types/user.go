package types

// UserProfile 用户资料
type UserProfile struct {
	// 用户 ID
	ID string `json:"id"`
	// 邮箱
	Email string `json:"email"`
	// 昵称
	Nickname string `json:"nickname"`
	// 头像 URL
	Avatar string `json:"avatar,omitempty"`
	// 个人简介
	Bio string `json:"bio,omitempty"`
	// 账号状态
	Status string `json:"status"`
	// 邮箱是否验证
	EmailVerified bool `json:"email_verified"`
	// 偏好设置
	Preferences UserPreferences `json:"preferences"`
	// 创建时间
	CreatedAt string `json:"created_at"`
	// 更新时间
	UpdatedAt string `json:"updated_at"`
}

// UserPreferences 用户偏好
type UserPreferences struct {
	// 感兴趣的分类
	InterestedCategories []string `json:"interested_categories,omitempty"`
	// 关注的职业
	FollowedProfessions []string `json:"followed_professions,omitempty"`
	// 邮件通知设置
	EmailNotifications EmailNotificationSettings `json:"email_notifications"`
	// 主题
	Theme string `json:"theme,omitempty"`
	// 语言
	Language string `json:"language,omitempty"`
}

// EmailNotificationSettings 邮件通知设置
type EmailNotificationSettings struct {
	// 系统更新通知
	SystemUpdate bool `json:"system_update"`
	// 新内容推送
	NewContent bool `json:"new_content"`
	// 周报
	WeeklyDigest bool `json:"weekly_digest"`
}

// UpdateUserProfileRequest 更新用户资料请求
type UpdateUserProfileRequest struct {
	// 昵称
	Nickname string `json:"nickname,omitempty"`
	// 头像 URL
	Avatar string `json:"avatar,omitempty"`
	// 个人简介
	Bio string `json:"bio,omitempty"`
}

// UpdateUserPreferencesRequest 更新用户偏好请求
type UpdateUserPreferencesRequest struct {
	// 感兴趣的分类
	InterestedCategories []string `json:"interested_categories,omitempty"`
	// 关注的职业
	FollowedProfessions []string `json:"followed_professions,omitempty"`
	// 邮件通知设置
	EmailNotifications *EmailNotificationSettings `json:"email_notifications,omitempty"`
	// 主题
	Theme string `json:"theme,omitempty"`
	// 语言
	Language string `json:"language,omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	// 旧密码
	OldPassword string `json:"old_password"`
	// 新密码
	NewPassword string `json:"new_password"`
}
