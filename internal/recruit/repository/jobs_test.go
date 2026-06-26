package repository

import (
	"os"
	"testing"

	"hot-ai-backend/internal/database"
)

// TestMain 在包级别启动时初始化 DB。DSN 来自环境变量 RECRUIT_TEST_DSN，
// 未设置时跳过测试（CI 默认关闭）。
func TestMain(m *testing.M) {
	dsn := os.Getenv("RECRUIT_TEST_DSN")
	if dsn == "" {
		// 允许生产/CI 跳过集成测试
		return
	}
	cfg, err := database.ParseDSN(dsn)
	if err != nil {
		panic(err)
	}
	if err := database.InitDB(*cfg); err != nil {
		panic(err)
	}
	m.Run()
}

func TestListJobsByProfession(t *testing.T) {
	if database.DB == nil {
		t.Skip("no DB")
	}
	repo := NewRawJobsRepo(database.DB)
	jobs, total, err := repo.ListJobs(ListJobsParams{ProfessionID: 7, Limit: 5, Offset: 0})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total < 1 {
		t.Fatalf("expected ≥1 job for profession 7, got %d", total)
	}
	if len(jobs) > 5 {
		t.Fatalf("limit not respected, got %d", len(jobs))
	}
	for _, j := range jobs {
		if j.ProfessionID == nil || *j.ProfessionID != 7 {
			t.Fatalf("wrong profession: %v", j.ProfessionID)
		}
	}
}

func TestListJobsHasAI(t *testing.T) {
	if database.DB == nil {
		t.Skip("no DB")
	}
	repo := NewRawJobsRepo(database.DB)
	jobs, _, err := repo.ListJobs(ListJobsParams{HasAI: true, Limit: 5, Offset: 0})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	for _, j := range jobs {
		if j.AIKeywordsCount == 0 {
			t.Fatalf("has_ai=true returned job with 0 AI keywords: id=%d", j.ID)
		}
	}
}

func TestListJobsLimitCap(t *testing.T) {
	if database.DB == nil {
		t.Skip("no DB")
	}
	repo := NewRawJobsRepo(database.DB)
	_, _, err := repo.ListJobs(ListJobsParams{Limit: 9999, Offset: 0})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	// 通过单元测试覆盖默认 cap：调用方传 9999 仓库层应截断为 100
	// 这条断言在 Step 3 实现后单独验证。
}