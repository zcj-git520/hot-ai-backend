package recruit

import (
	"strings"
)

// AgentCaller 调用 agent 的抽象（mock + real）
type AgentCaller interface {
	Normalize(title, description string) (pid uint, confidence float64, err error)
}

// KeywordsLoader 加载所有 profession 关键词
type KeywordsLoader interface {
	AllKeywords() (map[uint][]string, error)
}

type Normalizer struct {
	agent    AgentCaller
	keywords KeywordsLoader
}

func NewNormalizer(agent AgentCaller, kw KeywordsLoader) *Normalizer {
	return &Normalizer{agent: agent, keywords: kw}
}

// Normalize 返回 (profession_id, method, error)
// method ∈ {"llm", "keyword", ""}；profession_id=0 表示未匹配
func (n *Normalizer) Normalize(job *RawJob) (uint, string, error) {
	if n.agent != nil {
		pid, _, err := n.agent.Normalize(job.Title, job.Description)
		if err == nil && pid > 0 {
			return pid, "llm", nil
		}
	}
	all, err := n.keywords.AllKeywords()
	if err != nil {
		return 0, "", err
	}
	best := matchByKeywords(job, all)
	if best > 0 {
		return best, "keyword", nil
	}
	return 0, "", nil
}

func matchByKeywords(job *RawJob, all map[uint][]string) uint {
	text := job.Title + " " + job.Description
	var best uint
	bestScore := 0
	for pid, kws := range all {
		score := 0
		for _, kw := range kws {
			if strings.Contains(text, kw) {
				score++
			}
		}
		if score > bestScore {
			bestScore = score
			best = pid
		}
	}
	return best
}
