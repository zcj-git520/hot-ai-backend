package adapters

import (
	"fmt"
	"net/http"
	"time"

	"hot-ai-backend/internal/recruit"
)

// BossAdapter Boss直聘（带 Mock 模式）
type BossAdapter struct {
	client  *http.Client
	baseURL string
	mock    bool
}

func NewBossAdapter(client *http.Client, baseURL string) *BossAdapter {
	return &BossAdapter{client: client, baseURL: baseURL}
}

func NewBossMockAdapter() *BossAdapter {
	return &BossAdapter{mock: true, baseURL: "mock://boss"}
}

func (a *BossAdapter) Platform() recruit.Platform { return recruit.PlatformBoss }

func (a *BossAdapter) FetchJobs(keyword, city string) ([]recruit.RawJob, error) {
	if a.mock {
		return mockJobs(recruit.PlatformBoss, keyword, city), nil
	}
	url := fmt.Sprintf("%s/web/geek/job?query=%s&city=%s&page=1", a.baseURL, keyword, city)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("boss: status %d", resp.StatusCode)
	}
	html := readAll(resp.Body)
	cards, err := ParseJobCards(html, ".job-card")
	if err != nil {
		return nil, err
	}
	return cardsToRaw(cards, recruit.PlatformBoss), nil
}

// ZhilianAdapter 智联招聘
type ZhilianAdapter struct {
	client  *http.Client
	baseURL string
	mock    bool
}

func NewZhilianAdapter(client *http.Client, baseURL string) *ZhilianAdapter {
	return &ZhilianAdapter{client: client, baseURL: baseURL}
}

func NewZhilianMockAdapter() *ZhilianAdapter {
	return &ZhilianAdapter{mock: true, baseURL: "mock://zhilian"}
}

func (a *ZhilianAdapter) Platform() recruit.Platform { return recruit.PlatformZhilian }

func (a *ZhilianAdapter) FetchJobs(keyword, city string) ([]recruit.RawJob, error) {
	if a.mock {
		return mockJobs(recruit.PlatformZhilian, keyword, city), nil
	}
	url := fmt.Sprintf("%s/jobs/searchresult.ashx?kw=%s&jl=%s&p=1", a.baseURL, keyword, city)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("zhilian: status %d", resp.StatusCode)
	}
	html := readAll(resp.Body)
	cards, err := ParseJobCards(html, ".joblist_item")
	if err != nil {
		return nil, err
	}
	return cardsToRaw(cards, recruit.PlatformZhilian), nil
}

// LiepinAdapter 猎聘
type LiepinAdapter struct {
	client  *http.Client
	baseURL string
	mock    bool
}

func NewLiepinAdapter(client *http.Client, baseURL string) *LiepinAdapter {
	return &LiepinAdapter{client: client, baseURL: baseURL}
}

func NewLiepinMockAdapter() *LiepinAdapter {
	return &LiepinAdapter{mock: true, baseURL: "mock://liepin"}
}

func (a *LiepinAdapter) Platform() recruit.Platform { return recruit.PlatformLiepin }

func (a *LiepinAdapter) FetchJobs(keyword, city string) ([]recruit.RawJob, error) {
	if a.mock {
		return mockJobs(recruit.PlatformLiepin, keyword, city), nil
	}
	url := fmt.Sprintf("%s/zhaopin/?key=%s&city=%s&currentPage=1", a.baseURL, keyword, city)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("liepin: status %d", resp.StatusCode)
	}
	html := readAll(resp.Body)
	cards, err := ParseJobCards(html, ".so-job-item")
	if err != nil {
		return nil, err
	}
	return cardsToRaw(cards, recruit.PlatformLiepin), nil
}

func cardsToRaw(cards []JobCard, p recruit.Platform) []recruit.RawJob {
	now := time.Now()
	out := make([]recruit.RawJob, 0, len(cards))
	for _, c := range cards {
		min, max, _ := ParseSalaryRange(c.Salary)
		var pmin, pmax *int
		if min > 0 {
			pmin = &min
		}
		if max > 0 {
			pmax = &max
		}
		out = append(out, recruit.RawJob{
			Platform:      p,
			PlatformJobID: c.JobID,
			Title:         c.Title,
			Company:       c.Company,
			City:          c.City,
			SalaryMin:     pmin,
			SalaryMax:     pmax,
			URL:           c.Href,
			CrawledAt:     now,
		})
	}
	return out
}

// mockJobs 模拟器，按 city 和 keyword 生成确定的岗位数据
func mockJobs(p recruit.Platform, keyword, city string) []recruit.RawJob {
	now := time.Now()
	out := []recruit.RawJob{}
	cities := []string{"北京", "上海", "深圳", "广州", "杭州", "成都"}
	salaries := []string{"15-25K", "20-35K", "10-18K", "25-40K", "12-20K"}
	companies := []string{"某互联网公司", "某 AI 公司", "某数据公司", "某设计公司", "某创业公司"}
	for i := 0; i < 8; i++ {
		sal := salaries[i%len(salaries)]
		min, max, _ := ParseSalaryRange(sal)
		city2 := cities[(i+len(city))%len(cities)]
		if city != "" {
			city2 = city
		}
		out = append(out, recruit.RawJob{
			Platform:      p,
			PlatformJobID: fmt.Sprintf("%s-%s-%d", p, keyword, i),
			Title:         fmt.Sprintf("%s %s", keyword, []string{"工程师", "高级", "资深", "助理"}[i%4]),
			Company:       companies[i%len(companies)],
			City:          city2,
			SalaryMin:     &min,
			SalaryMax:     &max,
			Description:   fmt.Sprintf("我们需要一个%s。要求熟悉 %s。AI 工具 GPT/Claude 经验优先。", keyword, keyword),
			CrawledAt:     now,
		})
	}
	return out
}
