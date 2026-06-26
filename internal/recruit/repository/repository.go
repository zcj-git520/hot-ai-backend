package repository

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	"hot-ai-backend/internal/recruit"
	"gorm.io/gorm"
)

// RawJobsRepo raw_jobs 仓储
type RawJobsRepo struct {
	db *gorm.DB
}

func NewRawJobsRepo(db *gorm.DB) *RawJobsRepo { return &RawJobsRepo{db: db} }

// Upsert 插入或返回已存在行
func (r *RawJobsRepo) Upsert(job *recruit.RawJob) (uint64, error) {
	var existingID uint64
	err := r.db.Raw(
		"SELECT id FROM recruit_raw_jobs WHERE platform = ? AND platform_job_id = ?",
		job.Platform, job.PlatformJobID,
	).Scan(&existingID).Error
	if err == nil && existingID > 0 {
		return existingID, nil
	}
	if err := r.db.Create(job).Error; err != nil {
		return 0, err
	}
	return job.ID, nil
}

// CountByProfession 统计某 profession 在指定时间窗口内的归一后岗位数
func (r *RawJobsRepo) CountByProfession(pid uint, sinceDays int) (int, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*) FROM recruit_normalized_jobs n
		WHERE n.profession_id = ? AND n.normalized_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
	`, pid, sinceDays).Scan(&count).Error
	return int(count), err
}

// NormalizedJobsRepo normalized_jobs 仓储
type NormalizedJobsRepo struct {
	db *gorm.DB
}

func NewNormalizedJobsRepo(db *gorm.DB) *NormalizedJobsRepo { return &NormalizedJobsRepo{db: db} }

func (r *NormalizedJobsRepo) Insert(n *recruit.NormalizedJob) (uint64, error) {
	hitsJSON, err := json.Marshal(n.AIKeywordHits)
	if err != nil {
		return 0, err
	}
	res := r.db.Exec(`
		INSERT INTO recruit_normalized_jobs
			(raw_job_id, profession_id, match_method, match_confidence,
			 ai_keywords_count, ai_keywords_total, ai_keyword_hits, normalized_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		n.RawJobID, n.ProfessionID, n.MatchMethod, n.MatchConfidence,
		n.AIKeywordsCount, n.AIKeywordsTotal, hitsJSON, n.NormalizedAt,
	)
	if res.Error != nil {
		return 0, res.Error
	}
	id := uint64(res.RowsAffected)
	return id, nil
}

func (r *NormalizedJobsRepo) ListSince(pid uint, sinceDays int) ([]recruit.NormalizedJob, error) {
	var out []recruit.NormalizedJob
	err := r.db.Raw(`
		SELECT id, raw_job_id, profession_id, match_method, match_confidence,
			ai_keywords_count, ai_keywords_total, normalized_at
		FROM recruit_normalized_jobs
		WHERE profession_id = ? AND normalized_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		ORDER BY normalized_at DESC
	`, pid, sinceDays).Scan(&out).Error
	return out, err
}

// ListToday 今日归一后的所有岗位
func (r *NormalizedJobsRepo) ListToday() ([]recruit.NormalizedJob, error) {
	var out []recruit.NormalizedJob
	err := r.db.Raw(`
		SELECT id, raw_job_id, profession_id, match_method, match_confidence,
			ai_keywords_count, ai_keywords_total, normalized_at
		FROM recruit_normalized_jobs
		WHERE normalized_at >= CURDATE()
	`).Scan(&out).Error
	return out, err
}

// DailyMetricsRepo daily_metrics 仓储
type DailyMetricsRepo struct {
	db *gorm.DB
}

func NewDailyMetricsRepo(db *gorm.DB) *DailyMetricsRepo { return &DailyMetricsRepo{db: db} }

func (r *DailyMetricsRepo) Upsert(m *recruit.DailyMetrics) error {
	return r.db.Save(m).Error
}

func (r *DailyMetricsRepo) GetByDate(professionID uint, date time.Time) (*recruit.DailyMetrics, error) {
	var m recruit.DailyMetrics
	err := r.db.Where("profession_id = ? AND metric_date = ?", professionID, date).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// KeywordsRepo 关键词仓储
type KeywordsRepo struct {
	db *gorm.DB
}

func NewKeywordsRepo(db *gorm.DB) *KeywordsRepo { return &KeywordsRepo{db: db} }

func (r *KeywordsRepo) ListByProfession(pid uint) ([]recruit.Keyword, error) {
	var kws []recruit.Keyword
	if err := r.db.Where("profession_id = ?", pid).Find(&kws).Error; err != nil {
		return nil, err
	}
	return kws, nil
}

func (r *KeywordsRepo) AllKeywords() (map[uint][]string, error) {
	var kws []recruit.Keyword
	if err := r.db.Find(&kws).Error; err != nil {
		return nil, err
	}
	out := make(map[uint][]string)
	for _, k := range kws {
		out[k.ProfessionID] = append(out[k.ProfessionID], k.Keyword)
	}
	return out, nil
}

// CrawlLogRepo 抓取日志仓储
type CrawlLogRepo struct {
	db *gorm.DB
}

func NewCrawlLogRepo(db *gorm.DB) *CrawlLogRepo { return &CrawlLogRepo{db: db} }

func (r *CrawlLogRepo) Insert(log *recruit.CrawlLog) error {
	return r.db.Create(log).Error
}

// ConfigRepo 配置仓储
type ConfigRepo struct {
	db *gorm.DB
}

func NewConfigRepo(db *gorm.DB) *ConfigRepo { return &ConfigRepo{db: db} }

func (r *ConfigRepo) Get(key string) (float64, error) {
	var v float64
	err := r.db.Raw("SELECT config_value FROM recruit_config WHERE config_key = ?", key).Scan(&v).Error
	return v, err
}

func (r *ConfigRepo) GetAll() (map[string]float64, error) {
	var rows []recruit.Config
	if err := r.db.Find(&rows).Error; err != nil {
		return nil, err
	}
	m := make(map[string]float64, len(rows))
	for _, c := range rows {
		m[c.ConfigKey] = c.ConfigValue
	}
	return m, nil
}

// JobListItem 单条岗位展示模型
type JobListItem struct {
	ID              uint64          `json:"id" gorm:"column:id"`
	Title           string          `json:"title" gorm:"column:title"`
	Company         string          `json:"company" gorm:"column:company"`
	City            string          `json:"city" gorm:"column:city"`
	SalaryMin       *int            `json:"salary_min" gorm:"column:salary_min"`
	SalaryMax       *int            `json:"salary_max" gorm:"column:salary_max"`
	Description     string          `json:"description" gorm:"column:description"`
	URL             string          `json:"url" gorm:"column:url"`
	Platform        string          `json:"platform" gorm:"column:platform"`
	CrawledAt       time.Time       `json:"crawled_at" gorm:"column:crawled_at"`
	ProfessionID    *uint           `json:"profession_id" gorm:"column:profession_id"`
	AIKeywordsCount int             `json:"ai_keyword_count" gorm:"column:ai_keywords_count"`
	AIKeywordsRaw   json.RawMessage `json:"-" gorm:"column:ai_keyword_hits"`
	AIKeywords      []string        `json:"ai_keywords" gorm:"-"`
}

// ListJobsParams 列表查询参数
type ListJobsParams struct {
	ProfessionID uint
	City         string
	MinSalary    int
	HasAI        bool
	Limit        int
	Offset       int
}

const maxJobListLimit = 100

// ListJobs 列出归一后的岗位（JOIN raw + normalized）。同时返回 total 用于分页。
func (r *RawJobsRepo) ListJobs(p ListJobsParams) ([]JobListItem, int, error) {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > maxJobListLimit {
		p.Limit = maxJobListLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}

	conds := []string{"n.normalized_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)"}
	args := []interface{}{}
	if p.ProfessionID > 0 {
		conds = append(conds, "n.profession_id = ?")
		args = append(args, p.ProfessionID)
	}
	if p.City != "" {
		conds = append(conds, "r.city = ?")
		args = append(args, p.City)
	}
	if p.MinSalary > 0 {
		conds = append(conds, "r.salary_min >= ?")
		args = append(args, p.MinSalary)
	}
	if p.HasAI {
		conds = append(conds, "n.ai_keywords_count > 0")
	}
	where := "WHERE " + strings.Join(conds, " AND ")

	// total
	var total int
	if err := r.db.Raw(
		"SELECT COUNT(*) FROM recruit_normalized_jobs n JOIN recruit_raw_jobs r ON r.id = n.raw_job_id "+where,
		args...,
	).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// page
	pageArgs := append([]interface{}{}, args...)
	pageArgs = append(pageArgs, p.Limit, p.Offset)
	var items []JobListItem
	err := r.db.Raw(`
		SELECT r.id, r.title, r.company, r.city, r.salary_min, r.salary_max,
		       r.description, r.url, r.platform, r.crawled_at,
		       n.profession_id, n.ai_keywords_count, n.ai_keyword_hits
		FROM recruit_normalized_jobs n
		JOIN recruit_raw_jobs r ON r.id = n.raw_job_id
	`+where+`
		ORDER BY r.crawled_at DESC
		LIMIT ? OFFSET ?
	`, pageArgs...).Scan(&items).Error
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		if len(items[i].AIKeywordsRaw) > 0 {
			_ = json.Unmarshal(items[i].AIKeywordsRaw, &items[i].AIKeywords)
		}
		if items[i].AIKeywords == nil {
			items[i].AIKeywords = []string{}
		}
	}
	return items, total, nil
}

// JobStats 聚合统计
type JobStats struct {
	TopCompanies       []CompanyCount      `json:"top_companies"`
	SalaryDistribution []SalaryBucketCount `json:"salary_distribution"`
	CityDistribution   map[string]float64  `json:"city_distribution"`
	AIKeywordsTop      []KeywordCount      `json:"ai_keywords_top"`
	SampleSize         int                 `json:"sample_size"`
}

type CompanyCount struct {
	Company string `json:"company"`
	Count   int    `json:"count"`
}

type SalaryBucketCount struct {
	Bucket string `json:"bucket"`
	Count  int    `json:"count"`
}

type KeywordCount struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
}

var salaryBuckets = []struct {
	Min, Max int
	Label    string
}{
	{0, 10000, "<10k"},
	{10000, 20000, "10-20k"},
	{20000, 30000, "20-30k"},
	{30000, 50000, "30-50k"},
	{50000, 80000, "50-80k"},
	{80000, 999999, "80k+"},
}

// GetJobsStats 返回指定 profession 的 4 维聚合
func (r *RawJobsRepo) GetJobsStats(professionID uint) (*JobStats, error) {
	stats := &JobStats{
		CityDistribution: map[string]float64{},
	}

	// 1) Top companies
	if err := r.db.Raw(`
		SELECT r.company, COUNT(*) AS count
		FROM recruit_normalized_jobs n
		JOIN recruit_raw_jobs r ON r.id = n.raw_job_id
		WHERE n.profession_id = ? AND n.normalized_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
		GROUP BY r.company
		ORDER BY count DESC
		LIMIT 5
	`, professionID).Scan(&stats.TopCompanies).Error; err != nil {
		return nil, err
	}

	// 2) Salary distribution
	type sdRow struct {
		Min   *int `gorm:"column:salary_min"`
		Max   *int `gorm:"column:salary_max"`
		Count int  `gorm:"column:count"`
	}
	var sdRows []sdRow
	if err := r.db.Raw(`
		SELECT r.salary_min, r.salary_max, COUNT(*) AS count
		FROM recruit_normalized_jobs n
		JOIN recruit_raw_jobs r ON r.id = n.raw_job_id
		WHERE n.profession_id = ? AND n.normalized_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
		  AND r.salary_min IS NOT NULL AND r.salary_max IS NOT NULL
		GROUP BY r.salary_min, r.salary_max
	`, professionID).Scan(&sdRows).Error; err != nil {
		return nil, err
	}
	bucketCount := map[string]int{}
	for _, row := range sdRows {
		if row.Min == nil || row.Max == nil {
			continue
		}
		mid := (*row.Min + *row.Max) / 2
		for _, b := range salaryBuckets {
			if mid >= b.Min && mid < b.Max {
				bucketCount[b.Label] += row.Count
				break
			}
		}
	}
	for _, b := range salaryBuckets {
		if c := bucketCount[b.Label]; c > 0 {
			stats.SalaryDistribution = append(stats.SalaryDistribution, SalaryBucketCount{
				Bucket: b.Label, Count: c,
			})
		}
	}

	// 3) City distribution
	var cityRows []struct {
		City  string `gorm:"column:city"`
		Count int    `gorm:"column:count"`
	}
	if err := r.db.Raw(`
		SELECT r.city, COUNT(*) AS count
		FROM recruit_normalized_jobs n
		JOIN recruit_raw_jobs r ON r.id = n.raw_job_id
		WHERE n.profession_id = ? AND n.normalized_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
		  AND r.city != ''
		GROUP BY r.city
	`, professionID).Scan(&cityRows).Error; err != nil {
		return nil, err
	}
	totalCity := 0
	for _, c := range cityRows {
		totalCity += c.Count
	}
	if totalCity > 0 {
		for _, c := range cityRows {
			stats.CityDistribution[c.City] = float64(c.Count) / float64(totalCity)
		}
	}

	// 4) AI keywords top
	var kwRows []struct {
		Hits  string `gorm:"column:ai_keyword_hits"`
		Count int    `gorm:"column:c"`
	}
	if err := r.db.Raw(`
		SELECT n.ai_keyword_hits, COUNT(*) AS c
		FROM recruit_normalized_jobs n
		WHERE n.profession_id = ? AND n.normalized_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
		  AND n.ai_keywords_count > 0
		GROUP BY n.ai_keyword_hits
	`, professionID).Scan(&kwRows).Error; err != nil {
		return nil, err
	}
	kwAgg := map[string]int{}
	for _, row := range kwRows {
		var hits []string
		if err := json.Unmarshal([]byte(row.Hits), &hits); err != nil {
			continue
		}
		for _, h := range hits {
			kwAgg[h] += row.Count
		}
	}
	type kv struct {
		K string
		V int
	}
	var sorted []kv
	for k, v := range kwAgg {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].V > sorted[j].V })
	for i, item := range sorted {
		if i >= 5 {
			break
		}
		stats.AIKeywordsTop = append(stats.AIKeywordsTop, KeywordCount{
			Keyword: item.K, Count: item.V,
		})
	}

	// 5) Sample size
	if err := r.db.Raw(`
		SELECT COUNT(*) FROM recruit_normalized_jobs n
		WHERE n.profession_id = ? AND n.normalized_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
	`, professionID).Scan(&stats.SampleSize).Error; err != nil {
		return nil, err
	}

	return stats, nil
}
