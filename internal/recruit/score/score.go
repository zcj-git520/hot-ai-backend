package score

import (
	"encoding/json"
	"math"

	"hot-ai-backend/internal/recruit"
)

type Weights struct {
	DemandDecay      float64
	SalaryDrop       float64
	AIPenetration    float64
	DistributionConc float64
}

// Compute 计算 market_confidence_score (0-100)
func Compute(m recruit.DailyMetrics, w Weights) float64 {
	decay := 0.0
	if m.JobCountPrev30d > 0 {
		decay = clamp((float64(m.JobCountPrev30d)-float64(m.JobCount))/float64(m.JobCountPrev30d)*100, 0, 100)
	}

	sal := 0.0
	if m.SalaryMedian != nil && m.SalaryMedianPrev90d != nil && *m.SalaryMedianPrev90d > 0 {
		sal = clamp((*m.SalaryMedianPrev90d-*m.SalaryMedian)/(*m.SalaryMedianPrev90d)*100, 0, 100)
	}

	ai := 0.0
	if m.AIPenetrationRate != nil {
		ai = clamp(*m.AIPenetrationRate*100, 0, 100)
	}

	conc := 0.0
	if len(m.GeoDistribution) > 0 {
		conc = distributionConcentration(m.GeoDistribution)
	}

	rawWeights := []float64{w.DemandDecay, w.SalaryDrop, w.AIPenetration, w.DistributionConc}
	values := []float64{decay, sal, ai, conc}

	// 哪一维有有效数据：decay 用 JobCountPrev30d>0 也算有效；其他非 nil/非 0 也算
	valid := []bool{m.JobCountPrev30d > 0, m.SalaryMedianPrev90d != nil, m.AIPenetrationRate != nil, len(m.GeoDistribution) > 0}

	var sumW, dot float64
	for i := range rawWeights {
		if valid[i] {
			sumW += rawWeights[i]
			dot += rawWeights[i] * values[i]
		}
	}
	if sumW == 0 {
		return 100.0
	}
	return clamp(100-dot/sumW, 0, 100)
}

func distributionConcentration(raw json.RawMessage) float64 {
	var dist map[string]float64
	if err := json.Unmarshal(raw, &dist); err != nil {
		return 0
	}
	if len(dist) == 0 {
		return 0
	}
	var h float64
	for _, p := range dist {
		if p > 0 {
			h -= p * math.Log2(p)
		}
	}
	hMax := math.Log2(float64(len(dist)))
	if hMax == 0 {
		return 0
	}
	return clamp(100*(1-h/hMax), 0, 100)
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
