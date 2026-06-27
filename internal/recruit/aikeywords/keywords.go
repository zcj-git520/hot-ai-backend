package aikeywords

import "strings"

var keywords = []string{
	"AI", "人工智能", "GPT", "LLM", "大模型", "AIGC", "提示词", "Prompt", "RAG",
	"智能体", "Agent", "Copilot", "文心一言", "通义千问", "Kimi", "DeepSeek",
	"Sora", "Midjourney", "Stable Diffusion", "数字人", "向量检索", "知识库",
	"prompt", "agent", "copilot", "embedding", "vector", "fine-tuning",
	"multimodal", "generative",
}

// Contains 判断 text 是否含 AI 关键词（大小写不敏感）
func Contains(text string) bool {
	lower := strings.ToLower(text)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

// CountInText 统计命中关键词次数与总 token 数
func CountInText(text string) (hits, total int) {
	lower := strings.ToLower(text)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			hits++
		}
	}
	tokens := strings.Fields(text)
	total = len(tokens)
	if total == 0 {
		total = 1
	}
	return
}

// HitKeywords 返回所有命中的关键词
func HitKeywords(text string) []string {
	lower := strings.ToLower(text)
	hits := []string{}
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			hits = append(hits, kw)
		}
	}
	return hits
}
