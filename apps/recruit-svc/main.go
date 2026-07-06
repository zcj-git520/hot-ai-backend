package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"gorm.io/gorm"

	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/recruit"
	"hot-ai-backend/internal/recruit/adapters"
	"hot-ai-backend/internal/recruit/aikeywords"
	"hot-ai-backend/internal/recruit/repository"
	"hot-ai-backend/internal/recruit/score"
)

var configFile = flag.String("f", "apps/recruit-svc/etc/recruit-svc.yaml", "the config file")

type RecruitSvcConf struct {
	rest.RestConf
	DataSource struct {
		MySQL struct {
			DSN string `json:",optional"`
		}
	}
	Agent  string `json:",default=http://localhost:8889"`
	UseLLM bool   `json:",default=false"`
}

type handler struct {
	db      *gorm.DB
	useLLM  bool
	agentURL string
}

func main() {
	flag.Parse()
	var c RecruitSvcConf
	if err := conf.Load(*configFile, &c); err != nil {
		fmt.Fprintf(os.Stderr, "load config error: %v\n", err)
		os.Exit(1)
	}
	if c.DataSource.MySQL.DSN == "" {
		fmt.Fprintln(os.Stderr, "no DSN")
		os.Exit(1)
	}
	dbConfig, err := database.ParseDSN(c.DataSource.MySQL.DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse DSN: %v\n", err)
		os.Exit(1)
	}
	if err := database.InitDB(*dbConfig); err != nil {
		fmt.Fprintf(os.Stderr, "init DB: %v\n", err)
		os.Exit(1)
	}

	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	h := &handler{db: database.DB, useLLM: c.UseLLM, agentURL: c.Agent}
	server.AddRoute(rest.Route{Method: "GET", Path: "/api/recruit/health", Handler: h.health})
	server.AddRoute(rest.Route{Method: "POST", Path: "/api/recruit/run-now", Handler: h.runNow})
	server.AddRoute(rest.Route{Method: "GET", Path: "/api/recruit/pulse/:id", Handler: h.pulse})
	server.AddRoute(rest.Route{Method: "GET", Path: "/api/recruit/jobs", Handler: h.listJobs})
	server.AddRoute(rest.Route{Method: "GET", Path: "/api/recruit/jobs/stats", Handler: h.jobsStats})
	server.AddRoute(rest.Route{Method: "GET", Path: "/api/recruit/jobs/:id", Handler: h.jobDetail})

	fmt.Printf("Starting recruit-svc at %s:%d (useLLM=%v)...\n", c.Host, c.Port, c.UseLLM)
	server.Start()
}

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","service":"recruit-svc"}`))
}

type runResult struct {
	Platforms  map[string]int `json:"platforms"`
	Raw        int            `json:"raw"`
	Normalized int            `json:"normalized"`
	Professions int           `json:"professions"`
	Scores     map[uint]int   `json:"scores"`
}

func (h *handler) runNow(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	result, err := h.runPipeline(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *handler) pulse(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		var params struct {
			ID uint64 `path:"id"`
		}
		if err := httpx.ParsePath(r, &params); err != nil || params.ID == 0 {
			http.Error(w, "invalid id", 400)
			return
		}
		idStr = strconv.FormatUint(params.ID, 10)
	}
	pid, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	var md struct {
		MarketConfidenceScore *int       `json:"market_confidence_score"`
		LastMetricsDate       *time.Time `json:"last_metrics_date"`
		DataFreshness         string     `json:"data_freshness"`
		MetricWindow          string     `json:"metric_window"`
	}
	if err := h.db.Raw(`
		SELECT market_confidence_score, last_metrics_date, data_freshness, metric_window
		FROM profession_market_data
		WHERE profession_id = ?
	`, pid).Scan(&md).Error; err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 拉今日 daily_metrics
	var dm recruit.DailyMetrics
	today := time.Now().Format("2006-01-02")
	h.db.Raw(`SELECT * FROM recruit_daily_metrics WHERE profession_id = ? AND metric_date = ?`,
		pid, today).Scan(&dm)

	out := map[string]interface{}{
		"profession_id":           pid,
		"market_confidence_score":  md.MarketConfidenceScore,
		"data_freshness":          md.DataFreshness,
		"last_metrics_date":       md.LastMetricsDate,
		"metric_window":           md.MetricWindow,
		"job_count":               dm.JobCount,
		"salary_median":           dm.SalaryMedian,
		"ai_penetration_rate":     dm.AIPenetrationRate,
		"geo_distribution":        json.RawMessage(dm.GeoDistribution),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *handler) runPipeline(ctx context.Context) (*runResult, error) {
	rawRepo := repository.NewRawJobsRepo(h.db)
	normRepo := repository.NewNormalizedJobsRepo(h.db)
	dmRepo := repository.NewDailyMetricsRepo(h.db)
	kwRepo := repository.NewKeywordsRepo(h.db)
	cfgRepo := repository.NewConfigRepo(h.db)

	// 1) 加载职业关键词
	allKw, err := kwRepo.AllKeywords()
	if err != nil {
		return nil, fmt.Errorf("load keywords: %w", err)
	}

	// 2) 拉每个 profession 的代表性关键词（每个 profession 取 weight 最高的前 3）
	professions := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9}
	profKeywords := map[uint][]string{}
	for _, pid := range professions {
		kw, _ := kwRepo.ListByProfession(pid)
		sort.Slice(kw, func(i, j int) bool { return kw[i].Weight > kw[j].Weight })
		max := 3
		if len(kw) < max {
			max = len(kw)
		}
		for i := 0; i < max; i++ {
			profKeywords[pid] = append(profKeywords[pid], kw[i].Keyword)
		}
	}

	result := &runResult{
		Platforms:  map[string]int{},
		Scores:     map[uint]int{},
	}

	// 4) 抓取（真实数据：直接走 Jobicy + Remotive 公共 API，不再贴 boss/智联/猎聘 假牌）
	platforms := []struct {
		name  string
		p     recruit.Platform
		fetch func() []recruit.RawJob
	}{
		{"jobicy", recruit.PlatformJobicy, func() []recruit.RawJob {
			a := adapters.NewRealAdapter(recruit.PlatformJobicy)
			jobs, err := a.FetchJobs("", "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "real adapter jobicy: %v\n", err)
				return nil
			}
			return jobs
		}},
		{"remotive", recruit.PlatformRemotive, func() []recruit.RawJob {
			a := adapters.NewRealAdapter(recruit.PlatformRemotive)
			jobs, err := a.FetchJobs("", "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "real adapter remotive: %v\n", err)
				return nil
			}
			return jobs
		}},
	}

	allRaw := []recruit.RawJob{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, p := range platforms {
		wg.Add(1)
		go func(p struct {
			name  string
			p     recruit.Platform
			fetch func() []recruit.RawJob
		}) {
			defer wg.Done()
			jobs := p.fetch()
			mu.Lock()
			result.Platforms[p.name] = len(jobs)
			allRaw = append(allRaw, jobs...)
			mu.Unlock()
		}(p)
	}
	wg.Wait()

	// 5) Upsert raw
	for i := range allRaw {
		if _, err := rawRepo.Upsert(&allRaw[i]); err != nil {
			return nil, fmt.Errorf("upsert raw: %w", err)
		}
	}
	result.Raw = len(allRaw)

	// 归一前先清掉这些 raw_job_id 今天已生成的归一行，保证重复执行不累积
	rawIDs := make([]uint64, 0, len(allRaw))
	for i := range allRaw {
		if allRaw[i].ID > 0 {
			rawIDs = append(rawIDs, allRaw[i].ID)
		}
	}
	if err := normRepo.DeleteTodayByRawJobIDs(rawIDs); err != nil {
		return nil, fmt.Errorf("clear today norm: %w", err)
	}

	// 6) 归一（关键词 fallback，无 LLM）
	for i := range allRaw {
		job := allRaw[i]
		text := job.Title + " " + job.Description
		bestPID, bestScore := uint(0), 0
		for pid, kws := range allKw {
			score := 0
			for _, kw := range kws {
				if contains(text, kw) {
					score++
				}
			}
			if score > bestScore {
				bestScore = score
				bestPID = pid
			}
		}
		hits, total := aikeywords.CountInText(text)
		aiHits := aikeywords.HitKeywords(text)
		var pid *uint
		var matchMethod string
		var conf *float64
		if bestPID > 0 {
			pid = &bestPID
			matchMethod = "keyword"
			c := 0.6
			conf = &c
		} else {
			matchMethod = "keyword"
		}
		norm := recruit.NormalizedJob{
			RawJobID:        job.ID,
			ProfessionID:    pid,
			MatchMethod:     matchMethod,
			MatchConfidence: conf,
			AIKeywordsCount: hits,
			AIKeywordsTotal: total,
			AIKeywordHits:   aiHits,
			NormalizedAt:    time.Now(),
		}
		if _, err := normRepo.Insert(&norm); err != nil {
			return nil, fmt.Errorf("insert norm: %w", err)
		}
	}

	// 7) 聚合
	today := time.Now()
	todayNorm, _ := normRepo.ListToday()
	result.Normalized = len(todayNorm)
	groupedByProf := map[uint][]recruit.NormalizedJob{}
	aiHitByProf := map[uint]int{}
	cityByProf := map[uint]map[string]int{}
	salaryByProf := map[uint][]float64{}
	for _, n := range todayNorm {
		if n.ProfessionID == nil {
			continue
		}
		pid := *n.ProfessionID
		groupedByProf[pid] = append(groupedByProf[pid], n)
		aiHitByProf[pid] += n.AIKeywordsCount
	}
	// 拉 raw 的薪资/城市
	for _, raw := range allRaw {
		// 反查 norm 拿到 pid
		var n recruit.NormalizedJob
		h.db.Raw(`SELECT id, raw_job_id, profession_id, match_method, match_confidence,
			ai_keywords_count, ai_keywords_total, normalized_at
			FROM recruit_normalized_jobs WHERE raw_job_id = ? ORDER BY id DESC LIMIT 1`, raw.ID).Scan(&n)
		if n.ProfessionID == nil {
			continue
		}
		pid := *n.ProfessionID
		if raw.SalaryMin != nil && raw.SalaryMax != nil {
			mid := float64(*raw.SalaryMin+*raw.SalaryMax) / 2
			salaryByProf[pid] = append(salaryByProf[pid], mid)
		}
		if raw.City != "" {
			if cityByProf[pid] == nil {
				cityByProf[pid] = map[string]int{}
			}
			cityByProf[pid][raw.City]++
		}
	}

	// 8) 算分
	_ = cfgRepo
	cfg, _ := cfgRepo.GetAll()
	w := score.Weights{
		DemandDecay:      cfg["weight.demand_decay"],
		SalaryDrop:       cfg["weight.salary_drop"],
		AIPenetration:    cfg["weight.ai_penetration"],
		DistributionConc: cfg["weight.distribution_concentration"],
	}

	for pid, items := range groupedByProf {
		n := len(items)
		aiRate := 0.0
		if n > 0 {
			aiRate = float64(aiHitByProf[pid]) / float64(n)
		}
		var median *float64
		if sals := salaryByProf[pid]; len(sals) > 0 {
			sort.Float64s(sals)
			m := sals[len(sals)/2]
			median = &m
		}
		geoBytes, _ := json.Marshal(normalizeDist(cityByProf[pid]))

		// 拿 30 天前数据
		var prevCount int
		h.db.Raw(`
			SELECT IFNULL(SUM(job_count),0) FROM recruit_daily_metrics
			WHERE profession_id = ? AND metric_date BETWEEN DATE_SUB(?, INTERVAL 60 DAY) AND DATE_SUB(?, INTERVAL 30 DAY)
		`, pid, today, today).Scan(&prevCount)
		var prevSalary *float64
		var prevRows []float64
		h.db.Raw(`
			SELECT salary_median FROM recruit_daily_metrics
			WHERE profession_id = ? AND metric_date BETWEEN DATE_SUB(?, INTERVAL 120 DAY) AND DATE_SUB(?, INTERVAL 90 DAY)
		`, pid, today, today).Scan(&prevRows)
		if len(prevRows) > 0 {
			sort.Float64s(prevRows)
			m := prevRows[len(prevRows)/2]
			prevSalary = &m
		}

		dm := recruit.DailyMetrics{
			ProfessionID:        pid,
			MetricDate:          today,
			JobCount:            n,
			JobCountPrev30d:     prevCount,
			SalaryMedian:        median,
			SalaryMedianPrev90d: prevSalary,
			AIPenetrationRate:   &aiRate,
			GeoDistribution:     geoBytes,
			SampleSize:          n,
		}
		_ = dmRepo.Upsert(&dm)

		cs := score.Compute(dm, w)
		scoreVal := int(math.Round(cs))

		// 写回 profession_market_data
		freshness := "fresh"
		h.db.Exec(`
			UPDATE profession_market_data
			SET market_confidence_score = ?, last_metrics_date = ?, data_freshness = ?, metric_window = '30d'
			WHERE profession_id = ?
		`, scoreVal, today.Format("2006-01-02"), freshness, pid)

		result.Professions++
		result.Scores[pid] = scoreVal
	}

	return result, nil
}

func runMockAdapter(p recruit.Platform, profKeywords map[uint][]string, cities []string) []recruit.RawJob {
	// 每个 profession 的代表性关键词，每个城市，5 条
	out := []recruit.RawJob{}
	now := time.Now()
	idx := 0
	for _, kws := range profKeywords {
		for _, city := range cities {
			for _, kw := range kws {
				for v := 0; v < 2; v++ {
					_ = v
					min, max := 15000+idx*100, 25000+idx*100
					out = append(out, recruit.RawJob{
						Platform:      p,
						PlatformJobID: fmt.Sprintf("%s-mock-%d", p, idx),
						Title:         fmt.Sprintf("%s 工程师", kw),
						Company:       "某公司",
						City:          city,
						SalaryMin:     &min,
						SalaryMax:     &max,
						Description:   fmt.Sprintf("招 %s，熟悉 AI/GPT/LLM 工具优先", kw),
						CrawledAt:     now,
					})
					idx++
				}
			}
		}
	}
	return out
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (stringIndex(haystack, needle) >= 0)
}

func stringIndex(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func normalizeDist(counts map[string]int) map[string]float64 {
	out := map[string]float64{}
	total := 0
	for _, c := range counts {
		total += c
	}
	if total == 0 {
		return out
	}
	for k, c := range counts {
		out[k] = float64(c) / float64(total)
	}
	return out
}

func (h *handler) listJobs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	p := repository.ListJobsParams{}

	if v := q.Get("profession_id"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			p.ProfessionID = uint(n)
		}
	}
	if v := q.Get("city"); v != "" {
		p.City = v
	}
	if v := q.Get("min_salary"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			p.MinSalary = n
		}
	}
	if v := q.Get("has_ai"); v == "true" || v == "1" {
		p.HasAI = true
	}
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			p.Limit = n
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			p.Offset = n
		}
	}

	repo := repository.NewRawJobsRepo(h.db)
	jobs, total, err := repo.ListJobs(p)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if jobs == nil {
		jobs = []repository.JobListItem{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jobs":   jobs,
		"total":  total,
		"limit":  p.Limit,
		"offset": p.Offset,
	})
}

func (h *handler) jobsStats(w http.ResponseWriter, r *http.Request) {
	pidStr := r.URL.Query().Get("profession_id")
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil || pid == 0 {
		http.Error(w, "profession_id required", 400)
		return
	}
	repo := repository.NewRawJobsRepo(h.db)
	stats, err := repo.GetJobsStats(uint(pid))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"profession_id":       pid,
		"top_companies":       stats.TopCompanies,
		"salary_distribution": stats.SalaryDistribution,
		"city_distribution":   stats.CityDistribution,
		"ai_keywords_top":     stats.AIKeywordsTop,
		"sample_size":         stats.SampleSize,
	})
}

func (h *handler) jobDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		var params struct {
			ID uint64 `path:"id"`
		}
		if err := httpx.ParsePath(r, &params); err != nil || params.ID == 0 {
			http.Error(w, "invalid id", 400)
			return
		}
		idStr = strconv.FormatUint(params.ID, 10)
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	var item repository.JobListItem
	err = h.db.Raw(`
		SELECT r.id, r.title, r.company, r.city, r.salary_min, r.salary_max,
		       r.description, r.url, r.platform, r.crawled_at,
		       n.profession_id, n.ai_keywords_count, n.ai_keyword_hits
		FROM recruit_normalized_jobs n
		JOIN recruit_raw_jobs r ON r.id = n.raw_job_id
		WHERE r.id = ?
		LIMIT 1
	`, id).Scan(&item).Error
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if item.ID == 0 {
		http.Error(w, "not found", 404)
		return
	}
	if len(item.AIKeywordsRaw) > 0 {
		_ = json.Unmarshal(item.AIKeywordsRaw, &item.AIKeywords)
	}
	if item.AIKeywords == nil {
		item.AIKeywords = []string{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
