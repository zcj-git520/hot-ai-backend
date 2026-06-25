package main

import (
	"fmt"

	"hot-ai-backend/internal/database"
)

func main() {
	dsn := "root:shds@Admin123@tcp(192.168.50.109:13306)/hot_ai?charset=utf8mb4&parseTime=true&loc=Local"
	cfg, _ := database.ParseDSN(dsn)
	database.InitDB(*cfg)

	// Show all tools with their slugs
	var tools []map[string]interface{}
	database.GetDB().Raw(`SELECT id, name, slug, access_level FROM tools ORDER BY id LIMIT 10`).Scan(&tools)
	for _, t := range tools { fmt.Printf("%v\n", t) }
}
