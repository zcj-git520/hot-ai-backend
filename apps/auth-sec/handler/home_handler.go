package handler

import (
	"net/http"

	"hot-ai-backend/internal/service"
	"hot-ai-backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// HomeHandler 首页处理器
type HomeHandler struct {
	articleService *service.ArticleService
	profService    *service.ProfessionService
	learningPathService *service.LearningPathService
	userService    *service.UserService
	toolService    *service.ToolService
}

// NewHomeHandler 创建首页处理器实例
func NewHomeHandler(
	articleService *service.ArticleService,
	profService *service.ProfessionService,
	learningPathService *service.LearningPathService,
	userService *service.UserService,
	toolService *service.ToolService,
) *HomeHandler {
	return &HomeHandler{
		articleService: articleService,
		profService:    profService,
		learningPathService: learningPathService,
		userService:    userService,
		toolService:    toolService,
	}
}

// GetHomeStats 首页统计信息
func (h *HomeHandler) GetHomeStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.getHomeStats()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(stats))
}

// GetHomeTopics 热门话题
func (h *HomeHandler) GetHomeTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := h.getHomeTopics()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"topics": topics,
	}))
}

// GetHomeRankings 热门排行
func (h *HomeHandler) GetHomeRankings(w http.ResponseWriter, r *http.Request) {
	rankings, err := h.getHomeRankings()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, types.Success(map[string]interface{}{
		"items": rankings,
	}))
}

// getHomeStats 获取首页统计数据
func (h *HomeHandler) getHomeStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取文章总数
	articleCount, err := h.articleService.GetArticleCount()
	if err != nil {
		return nil, err
	}
	stats["articles"] = articleCount
	stats["articlesGrowth"] = "12%" // 示例增长数据

	// 获取职业总数
	professionCount, err := h.profService.GetProfessionCount()
	if err != nil {
		return nil, err
	}
	stats["professions"] = professionCount
	stats["professionsGrowth"] = "8%"

	// 获取学习路径总数
	learningPathCount, err := h.learningPathService.GetLearningPathCount()
	if err != nil {
		return nil, err
	}
	stats["learningPaths"] = learningPathCount
	stats["learningPathsGrowth"] = "5%"

	// 获取用户总数
	userCount, err := h.userService.GetUserCount()
	if err != nil {
		return nil, err
	}
	stats["users"] = userCount
	stats["usersGrowth"] = "23%"

	return stats, nil
}

// getHomeTopics 获取热门话题
func (h *HomeHandler) getHomeTopics() ([]map[string]interface{}, error) {
	// 返回模拟数据，后续可以替换为真实数据
	topics := []map[string]interface{}{
		{
			"id":    1,
			"title": "AI 编程助手会取代程序员吗？",
			"summary": "随着 Cursor、GitHub Copilot 等工具的普及，程序员的工作方式正在发生深刻变化。是辅助工具还是替代者？",
			"hot":   "15.2k",
			"trend": "上升 25%",
			"rank":  1,
		},
		{
			"id":    2,
			"title": "2026 年最值得学习的 AI 技能",
			"summary": "从提示词工程到 AI 应用开发，盘点今年最具价值的 AI 相关技能和职业发展路径。",
			"hot":   "12.8k",
			"trend": "上升 18%",
			"rank":  2,
		},
	}

	return topics, nil
}

// getHomeRankings 获取热门排行
func (h *HomeHandler) getHomeRankings() ([]map[string]interface{}, error) {
	// 返回模拟数据，后续可以替换为真实数据
	rankings := []map[string]interface{}{
		{
			"id":    1,
			"title": "GPT-5 发布：AI 能力再次飞跃",
			"hot":   "25.3k",
		},
		{
			"id":    2,
			"title": "AI 绘画工具对比：Midjourney vs Stable Diffusion",
			"hot":   "18.7k",
		},
		{
			"id":    3,
			"title": "提示词工程入门指南",
			"hot":   "15.2k",
		},
		{
			"id":    4,
			"title": "AI 如何改变软件开发行业",
			"hot":   "12.5k",
		},
		{
			"id":    5,
			"title": "2026 年 AI 投资趋势分析",
			"hot":   "10.8k",
		},
	}

	return rankings, nil
}