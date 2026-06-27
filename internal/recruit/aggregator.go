package recruit

import (
	"encoding/json"
	"sort"
	"time"
)

// RawDataSource 给 Aggregator 提供 raw job 的薪资/城市聚合
type RawDataSource interface {
	SalariesByProfession(pid uint) []float64
	CityCountsByProfession(pid uint) map[string]int
}

type Aggregator struct {
	source RawDataSource
}

func NewAggregator(src RawDataSource) *Aggregator {
	return &Aggregator{source: src}
}

// ComputeDailyMetrics 从今日归一后的 jobs 算出每个 profession 的 4 维指标
func (a *Aggregator) ComputeDailyMetrics(jobs []NormalizedJob, today time.Time) ([]DailyMetrics, error) {
	grouped := groupByProfession(jobs)
	results := make([]DailyMetrics, 0, len(grouped))
	for pid, items := range grouped {
		aiHit, _, n := aiStats(items)
		aiRate := 0.0
		if n > 0 {
			aiRate = float64(aiHit) / float64(n)
		}
		salaries := a.source.SalariesByProfession(pid)
		var median *float64
		if len(salaries) > 0 {
			sort.Float64s(salaries)
			m := salaries[len(salaries)/2]
			median = &m
		}
		cities := a.source.CityCountsByProfession(pid)
		geoBytes, _ := json.Marshal(normalizeDist(cities))
		dm := DailyMetrics{
			ProfessionID:      pid,
			MetricDate:        today,
			JobCount:          n,
			AIPenetrationRate: &aiRate,
			SalaryMedian:      median,
			GeoDistribution:   geoBytes,
			SampleSize:        n,
		}
		results = append(results, dm)
	}
	return results, nil
}

func groupByProfession(jobs []NormalizedJob) map[uint][]NormalizedJob {
	m := map[uint][]NormalizedJob{}
	for _, j := range jobs {
		if j.ProfessionID == nil {
			continue
		}
		pid := *j.ProfessionID
		m[pid] = append(m[pid], j)
	}
	return m
}

func aiStats(items []NormalizedJob) (hits, total, n int) {
	for _, j := range items {
		n++
		hits += j.AIKeywordsCount
		total += j.AIKeywordsTotal
	}
	return
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
