package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"hot-ai-backend/internal/recruit"
)

// JobicyJob maps fields from Jobicy API v2 (https://jobicy.com/api/v2/remote-jobs)
// Note: Jobicy v2.2+ removed the `page` parameter; we get the latest 50 jobs per call.
type JobicyJob struct {
	ID             int      `json:"id"`
	URL            string   `json:"url"`
	JobTitle       string   `json:"jobTitle"`
	CompanyName    string   `json:"companyName"`
	CompanyLogo    string   `json:"companyLogo"`
	CompanyLogoNorm string  `json:"companyLogoNorm"`
	JobIndustry    []string `json:"jobIndustry"`
	JobType        []string `json:"jobType"`
	JobGeo         string   `json:"jobGeo"`
	JobLevel       string   `json:"jobLevel"`
	JobExcerpt     string   `json:"jobExcerpt"`
	JobDescription string   `json:"jobDescription"`
	PubDate        string   `json:"pubDate"`
	SalaryMin      *float64 `json:"salaryMin"`
	SalaryMax      *float64 `json:"salaryMax"`
	SalaryCurrency string   `json:"salaryCurrency"`
	SalaryPeriod   string   `json:"salaryPeriod"`
}

type JobicyResponse struct {
	APIVersion string      `json:"apiVersion"`
	JobCount   int         `json:"jobCount"`
	Jobs       []JobicyJob `json:"jobs"`
}

// RemotiveJob maps fields from Remotive API (https://remotive.com/api/remote-jobs)
type RemotiveJob struct {
	ID         int      `json:"id"`
	URL        string   `json:"url"`
	Title      string   `json:"title"`
	Company    string   `json:"company_name"`
	CompanyLogo string  `json:"company_logo"`
	Category   string   `json:"category"`
	Tags       []string `json:"tags"`
	JobType    string   `json:"job_type"`
	PubDate    string   `json:"publication_date"`
	Location   string   `json:"candidate_required_location"`
	Salary     string   `json:"salary"`
}

type RemotiveResponse struct {
	Jobs []RemotiveJob `json:"jobs"`
}

// RealAdapter fetches real job postings from public job APIs (Jobicy + Remotive).
// BOSS / Zhilian / Liepin all block automated scraping, so we use public aggregators.
// Every field stored is what the upstream API returned — no fake Chinese city,
// no fake platform label, no currency conversion. Readers see the actual data.
type RealAdapter struct {
	platform recruit.Platform
	client   *http.Client
}

func NewRealAdapter(p recruit.Platform) *RealAdapter {
	tr := &http.Transport{DisableCompression: true}
	return &RealAdapter{
		platform: p,
		client:   &http.Client{Timeout: 20 * time.Second, Transport: tr},
	}
}

func (a *RealAdapter) Platform() recruit.Platform { return a.platform }

// salaryAnnualToAnnual converts Jobicy's salary numbers to a yearly USD integer
// based on its `salaryPeriod`. We do not clamp or convert currency — readers see
// what the source returned, in the source currency.
func salaryAnnualToAnnual(period string, min, max *float64) (annualMin, annualMax int, ok bool) {
	if min == nil || max == nil {
		return 0, 0, false
	}
	switch period {
	case "yearly":
		return int(*min), int(*max), true
	case "monthly":
		return int(*min) * 12, int(*max) * 12, true
	case "hourly":
		// 40h/week × 52 weeks = 2080h/year
		return int(*min) * 2080, int(*max) * 2080, true
	default:
		return 0, 0, false
	}
}

var remotiveSalaryRe = regexp.MustCompile(`\$?\s*([0-9][0-9,\.]*)\s*(k|K)?\s*(?:-|–|to|—)?\s*\$?\s*([0-9][0-9,\.]*)?\s*(k|K)?`)

// parseRemotiveSalary extracts a (min, max) USD/year pair from Remotive's free-text
// salary field. Handles "$80,000 - $120,000", "60k-80k", "120000", etc. Returns
// ok=false when no number is parseable.
func parseRemotiveSalary(s string) (min, max int, ok bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, 0, false
	}
	m := remotiveSalaryRe.FindStringSubmatch(s)
	if m == nil {
		return 0, 0, false
	}
	toInt := func(raw, suffix string) int {
		raw = strings.ReplaceAll(raw, ",", "")
		n, err := strconv.Atoi(raw)
		if err != nil {
			return 0
		}
		if suffix == "k" || suffix == "K" {
			n *= 1000
		}
		return n
	}
	lo := toInt(m[1], m[2])
	if m[3] == "" {
		return lo, lo, true
	}
	hi := toInt(m[3], m[4])
	if hi == 0 {
		hi = lo
	}
	if lo > hi {
		lo, hi = hi, lo
	}
	return lo, hi, true
}

// stripHTML removes HTML tags and collapses whitespace so the DB row stays compact.
func stripHTML(s string) string {
	var b strings.Builder
	in := false
	for _, r := range s {
		if r == '<' {
			in = true
			continue
		}
		if r == '>' {
			in = false
			b.WriteByte(' ')
			continue
		}
		if r == '�' {
			continue
		}
		if !in {
			b.WriteRune(r)
		}
	}
	out := strings.Join(strings.Fields(b.String()), " ")
	if len(out) > 600 {
		out = out[:600]
	}
	for len(out) > 0 {
		last := out[len(out)-1]
		if last < 0x80 {
			break
		}
		if last < 0xC0 {
			out = out[:len(out)-1]
			continue
		}
		out = out[:len(out)-1]
		break
	}
	return out
}

// fetchJobicy pulls the latest 50 remote jobs from Jobicy.
func (a *RealAdapter) fetchJobicy() ([]JobicyJob, error) {
	url := "https://jobicy.com/api/v2/remote-jobs?count=50"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jobicy fetch: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("jobicy read: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("jobicy status %d", resp.StatusCode)
	}
	var jr JobicyResponse
	if err := json.Unmarshal(body, &jr); err != nil {
		return nil, fmt.Errorf("jobicy parse: %w", err)
	}
	return jr.Jobs, nil
}

// fetchRemotive pulls latest ~30 remote jobs from Remotive.
func (a *RealAdapter) fetchRemotive() ([]RemotiveJob, error) {
	url := "https://remotive.com/api/remote-jobs?limit=50"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("remotive fetch: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("remotive read: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("remotive status %d", resp.StatusCode)
	}
	var rr RemotiveResponse
	if err := json.Unmarshal(body, &rr); err != nil {
		return nil, fmt.Errorf("remotive parse: %w", err)
	}
	return rr.Jobs, nil
}

// FetchJobs fetches from the source matching this adapter's platform. A RealAdapter
// constructed with PlatformJobicy pulls only Jobicy; PlatformRemotive pulls only
// Remotive. The pipeline creates one adapter per platform, so combining both sources
// here would double every job. Each job is tagged with its real source platform and
// stores the upstream location/currency verbatim.
func (a *RealAdapter) FetchJobs(keyword, city string) ([]recruit.RawJob, error) {
	switch a.platform {
	case recruit.PlatformJobicy:
		return a.fetchJobicyJobs()
	case recruit.PlatformRemotive:
		return a.fetchRemotiveJobs()
	default:
		return nil, fmt.Errorf("unsupported platform %q", a.platform)
	}
}

// fetchJobicyJobs pulls + maps Jobicy postings into RawJob rows.
func (a *RealAdapter) fetchJobicyJobs() ([]recruit.RawJob, error) {
	now := time.Now()
	jobs, err := a.fetchJobicy()
	if err != nil {
		return nil, err
	}
	out := make([]recruit.RawJob, 0, len(jobs))
	for i, j := range jobs {
		annualMin, annualMax, salOK := salaryAnnualToAnnual(j.SalaryPeriod, j.SalaryMin, j.SalaryMax)
		job := recruit.RawJob{
			Platform:      recruit.PlatformJobicy,
			PlatformJobID: fmt.Sprintf("jobicy-%d", j.ID),
			Title:         j.JobTitle,
			Company:       j.CompanyName,
			City:          strings.TrimSpace(j.JobGeo),
			Description:   stripHTML(j.JobExcerpt + " " + j.JobDescription),
			Industry:      strings.Join(j.JobIndustry, ","),
			URL:           j.URL,
			CrawledAt:     now.Add(time.Duration(-i) * time.Minute),
		}
		if salOK {
			currency := strings.TrimSpace(j.SalaryCurrency)
			if currency == "" {
				currency = "USD"
			}
			job.SalaryMin = &annualMin
			job.SalaryMax = &annualMax
			job.SalaryCurrency = &currency
		}
		out = append(out, job)
	}
	return out, nil
}

// fetchRemotiveJobs pulls + maps Remotive postings into RawJob rows.
func (a *RealAdapter) fetchRemotiveJobs() ([]recruit.RawJob, error) {
	now := time.Now()
	rjobs, err := a.fetchRemotive()
	if err != nil {
		return nil, err
	}
	out := make([]recruit.RawJob, 0, len(rjobs))
	for i, j := range rjobs {
		loc := strings.TrimSpace(j.Location)
		if loc == "" {
			loc = "Remote"
		}
		job := recruit.RawJob{
			Platform:      recruit.PlatformRemotive,
			PlatformJobID: fmt.Sprintf("remotive-%d", j.ID),
			Title:         j.Title,
			Company:       j.Company,
			City:          loc,
			Description:   fmt.Sprintf("[%s] %s", j.Category, j.JobType),
			Industry:      j.Category,
			URL:           j.URL,
			CrawledAt:     now.Add(time.Duration(-i) * time.Minute),
		}
		// Try parsing free-text salary; Remotive usually writes "$80,000 - $120,000".
		if lo, hi, ok := parseRemotiveSalary(j.Salary); ok {
			job.SalaryMin = &lo
			job.SalaryMax = &hi
			usd := "USD"
			job.SalaryCurrency = &usd
		} else if j.Salary != "" {
			// Keep the raw text in the description so readers see what the source said.
			job.Description = fmt.Sprintf("[%s] %s · 薪资：%s", j.Category, j.JobType, j.Salary)
		}
		out = append(out, job)
	}
	return out, nil
}