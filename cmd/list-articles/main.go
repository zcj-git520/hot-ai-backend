package main

import (
	"fmt"
	"log"

	"hot-ai-backend/internal/database"
)

func main() {
	dsn := "root:shds@Admin123@tcp(192.168.50.109:13306)/hot_ai?charset=utf8mb4&parseTime=true&loc=Local"
	cfg, _ := database.ParseDSN(dsn)
	if err := database.InitDB(*cfg); err != nil { log.Fatal(err) }
	type row struct{ ID uint; Title string; AccessLevel int }
	var rows []row
	database.GetDB().Raw("SELECT id, title, access_level FROM articles ORDER BY id LIMIT 10").Scan(&rows)
	for _, r := range rows {
		t := r.Title
		if len(t) > 40 { t = t[:40] + "..." }
		fmt.Printf("id=%d  level=%d  %s\n", r.ID, r.AccessLevel, t)
	}
}
