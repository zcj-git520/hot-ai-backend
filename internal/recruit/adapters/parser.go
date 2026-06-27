package adapters

import (
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type JobCard struct {
	Title   string
	Company string
	City    string
	Salary  string
	Href    string
	JobID   string
}

func ParseJobCards(html []byte, selector string) ([]JobCard, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}
	var cards []JobCard
	doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".job-title").First().Text())
		href, _ := s.Find(".job-title").First().Attr("href")
		jobID := extractJobID(href)
		cards = append(cards, JobCard{
			Title:   title,
			Company: strings.TrimSpace(s.Find(".company").First().Text()),
			City:    strings.TrimSpace(s.Find(".city").First().Text()),
			Salary:  strings.TrimSpace(s.Find(".salary").First().Text()),
			Href:    href,
			JobID:   jobID,
		})
	})
	return cards, nil
}

var salaryRe = regexp.MustCompile(`(\d+)-?(\d+)?K`)

func ParseSalaryRange(s string) (min, max int, ok bool) {
	s = strings.TrimSpace(s)
	if s == "" || strings.Contains(s, "面议") {
		return 0, 0, false
	}
	m := salaryRe.FindStringSubmatch(s)
	if len(m) < 2 {
		return 0, 0, false
	}
	mn, _ := strconv.Atoi(m[1])
	mx := mn
	if len(m) >= 3 && m[2] != "" {
		mx, _ = strconv.Atoi(m[2])
	}
	return mn * 1000, mx * 1000, true
}

func extractJobID(href string) string {
	parts := strings.Split(href, "/")
	if len(parts) > 0 {
		last := parts[len(parts)-1]
		if last != "" {
			return last
		}
	}
	if idx := strings.Index(href, "id="); idx >= 0 {
		return href[idx+3:]
	}
	return href
}

func readAll(r io.Reader) []byte {
	b, _ := io.ReadAll(r)
	return b
}
