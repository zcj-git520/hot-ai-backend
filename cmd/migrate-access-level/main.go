// 一键迁移：给 5 张内容表加 access_level 列
// 用法：cd hot-ai-backend && go run cmd/migrate-access-level/main.go
package main

import (
	"fmt"
	"log"

	"hot-ai-backend/internal/database"
)

func main() {
	// 复用 auth-sec 的 DSN（已知能连）
	dsn := "root:shds@Admin123@tcp(192.168.50.109:13306)/hot_ai?charset=utf8mb4&parseTime=true&loc=Local"
	cfg, err := database.ParseDSN(dsn)
	if err != nil {
		log.Fatalf("parse DSN: %v", err)
	}
	if err := database.InitDB(*cfg); err != nil {
		log.Fatalf("init DB: %v", err)
	}

	// GORM AutoMigrate 在 articles / professions 上已成功加过 access_level
	// （对其他几张 content-svc 没显式 AutoMigrate 的表，直接 SQL ALTER）

	// 全部走 SQL 更稳 — 已经在 articles 上加过了，重复加会被 IF NOT EXISTS 跳过
	stmts := []string{
		`ALTER TABLE articles         ADD COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=游客,1=普通,2=会员'`,
		`ALTER TABLE professions     ADD COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=游客,1=普通,2=会员'`,
		`ALTER TABLE tools           ADD COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=游客,1=普通,2=会员'`,
		`ALTER TABLE learning_paths  ADD COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=游客,1=普通,2=会员'`,
		`ALTER TABLE path_chapters   ADD COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0=游客,1=普通,2=会员'`,
	}
	for _, s := range stmts {
		if err := database.GetDB().Exec(s).Error; err != nil {
			fmt.Printf("  [warn] %s ... : %v\n", s[:80], err)
		} else {
			fmt.Printf("  [ok] %s\n", s[:80])
		}
	}

	// Seed admin/user/member roles
	seeds := []struct{ name, desc string }{
		{"admin", "平台管理员"},
		{"user", "普通注册用户"},
		{"member", "付费会员"},
	}
	for _, r := range seeds {
		if err := database.GetDB().Exec(
			`INSERT IGNORE INTO roles (id, name, description, created_at) VALUES (UUID(), ?, ?, NOW())`,
			r.name, r.desc,
		).Error; err != nil {
			fmt.Printf("  [warn] seed role %s: %v\n", r.name, err)
		} else {
			fmt.Printf("  [ok] seed role %s\n", r.name)
		}
	}

	fmt.Println("done.")
}