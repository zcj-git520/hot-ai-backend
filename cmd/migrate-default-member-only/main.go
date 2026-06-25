// 一键迁移：把 access_level 默认值改成 2，给 users 加 is_member/member_expire_at
// 用法：cd hot-ai-backend && go run cmd/migrate-default-member-only/main.go
package main

import (
	"fmt"
	"log"

	"hot-ai-backend/internal/database"
)

func main() {
	dsn := "root:shds@Admin123@tcp(192.168.50.109:13306)/hot_ai?charset=utf8mb4&parseTime=true&loc=Local"
	cfg, err := database.ParseDSN(dsn)
	if err != nil {
		log.Fatalf("parse DSN: %v", err)
	}
	if err := database.InitDB(*cfg); err != nil {
		log.Fatalf("init DB: %v", err)
	}

	// 1) 改默认值为 2 (会员专享)
	alterDefaults := []string{
		`ALTER TABLE articles        MODIFY COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '0=游客, 1=普通用户, 2=会员'`,
		`ALTER TABLE professions    MODIFY COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '0=游客, 1=普通用户, 2=会员'`,
		`ALTER TABLE tools          MODIFY COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '0=游客, 1=普通用户, 2=会员'`,
		`ALTER TABLE learning_paths MODIFY COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '0=游客, 1=普通用户, 2=会员'`,
		`ALTER TABLE path_chapters  MODIFY COLUMN access_level TINYINT UNSIGNED NOT NULL DEFAULT 2 COMMENT '0=游客, 1=普通用户, 2=会员'`,
	}
	for _, s := range alterDefaults {
		if err := database.GetDB().Exec(s).Error; err != nil {
			fmt.Printf("  [warn] alter default: %v\n", err)
		} else {
			fmt.Println("  [ok] alter default for one table")
		}
	}

	// 2) 把所有现有数据设为 2 (会员专享)
	updates := []string{
		`UPDATE articles        SET access_level = 2 WHERE access_level = 0`,
		`UPDATE professions    SET access_level = 2 WHERE access_level = 0`,
		`UPDATE tools          SET access_level = 2 WHERE access_level = 0`,
		`UPDATE learning_paths SET access_level = 2 WHERE access_level = 0`,
		`UPDATE path_chapters  SET access_level = 2 WHERE access_level = 0`,
	}
	for _, s := range updates {
		r := database.GetDB().Exec(s)
		if r.Error != nil {
			fmt.Printf("  [warn] update: %v\n", r.Error)
		} else {
			fmt.Printf("  [ok] updated %d rows\n", r.RowsAffected)
		}
	}

	// 3) 给 users 加会员字段 (MySQL 不支持 IF NOT EXISTS，先查再改)
	type colInfo struct {
		Field string
		Type  string
	}
	checkCol := func(table, field string) bool {
		var c colInfo
		err := database.GetDB().Raw(
			"SELECT COLUMN_NAME as field, DATA_TYPE as type FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?",
			table, field,
		).Scan(&c).Error
		return err == nil && c.Field != ""
	}
	if !checkCol("users", "is_member") {
		if err := database.GetDB().Exec(`ALTER TABLE users ADD COLUMN is_member TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否为付费会员'`).Error; err != nil {
			fmt.Printf("  [warn] add is_member: %v\n", err)
		} else {
			fmt.Println("  [ok] add is_member")
		}
	} else {
		fmt.Println("  [skip] is_member 已存在")
	}
	if !checkCol("users", "member_expire_at") {
		if err := database.GetDB().Exec(`ALTER TABLE users ADD COLUMN member_expire_at DATETIME NULL COMMENT '会员到期时间（NULL = 永不过期）'`).Error; err != nil {
			fmt.Printf("  [warn] add member_expire_at: %v\n", err)
		} else {
			fmt.Println("  [ok] add member_expire_at")
		}
	} else {
		fmt.Println("  [skip] member_expire_at 已存在")
	}
	// idx_is_member
	var idxCount int
	database.GetDB().Raw(
		"SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND INDEX_NAME = 'idx_is_member'",
	).Scan(&idxCount)
	if idxCount == 0 {
		if err := database.GetDB().Exec(`ALTER TABLE users ADD INDEX idx_is_member (is_member)`).Error; err != nil {
			fmt.Printf("  [warn] add idx: %v\n", err)
		} else {
			fmt.Println("  [ok] add idx_is_member")
		}
	} else {
		fmt.Println("  [skip] idx_is_member 已存在")
	}

	fmt.Println("done.")
}
