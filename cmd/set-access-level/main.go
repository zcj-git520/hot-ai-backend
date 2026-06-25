// 设置指定 ID 的 access_level
// 用法：go run cmd/set-access-level/main.go <table> <id> <level>
package main

import (
	"fmt"
	"log"
	"os"

	"hot-ai-backend/internal/database"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("usage: go run cmd/set-access-level/main.go <articles|professions|tools|learning_paths|path_chapters> <id> <0|1|2>")
	}
	table := os.Args[1]
	id := os.Args[2]
	level := os.Args[3]

	allowed := map[string]bool{
		"articles": true, "professions": true, "tools": true,
		"learning_paths": true, "path_chapters": true,
	}
	if !allowed[table] {
		log.Fatalf("invalid table: %s", table)
	}

	dsn := "root:shds@Admin123@tcp(192.168.50.109:13306)/hot_ai?charset=utf8mb4&parseTime=true&loc=Local"
	cfg, err := database.ParseDSN(dsn)
	if err != nil {
		log.Fatalf("parse DSN: %v", err)
	}
	if err := database.InitDB(*cfg); err != nil {
		log.Fatalf("init DB: %v", err)
	}

	r := database.GetDB().Exec(fmt.Sprintf("UPDATE %s SET access_level = ? WHERE id = ?", table), level, id)
	if r.Error != nil {
		log.Fatalf("update: %v", r.Error)
	}
	fmt.Printf("[ok] %s id=%s → access_level=%s (%d row)\n", table, id, level, r.RowsAffected)
}
