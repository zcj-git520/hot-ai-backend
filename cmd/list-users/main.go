package main

import (
	"fmt"
	"log"

	"hot-ai-backend/internal/database"
)

func main() {
	dsn := "root:shds@Admin123@tcp(192.168.50.109:13306)/hot_ai?charset=utf8mb4&parseTime=true&loc=Local"
	cfg, err := database.ParseDSN(dsn)
	if err != nil { log.Fatalf("parse DSN: %v", err) }
	if err := database.InitDB(*cfg); err != nil { log.Fatalf("init DB: %v", err) }
	
	// Set article 689 access_level to 2 (member only)
	database.GetDB().Exec(`UPDATE articles SET access_level = 2 WHERE id = 689`)
	// Set profession 1 (if exists) access_level to 2
	database.GetDB().Exec(`UPDATE professions SET access_level = 2 WHERE id <= 2`)
	// Set tool 1 (if exists) access_level to 2
	database.GetDB().Exec(`UPDATE tools SET access_level = 2 WHERE id <= 3`)
	// Set first chapter of first learning path access_level to 2
	database.GetDB().Exec(`UPDATE path_chapters SET access_level = 2 WHERE path_id = 1 ORDER BY order_index LIMIT 2`)
	
	// Verify
	var arts []map[string]interface{}
	database.GetDB().Raw(`SELECT id, title, access_level FROM articles WHERE access_level > 0 LIMIT 5`).Scan(&arts)
	for _, a := range arts { fmt.Printf("article: %v\n", a) }
	var profs []map[string]interface{}
	database.GetDB().Raw(`SELECT id, name, access_level FROM professions WHERE access_level > 0 LIMIT 5`).Scan(&profs)
	for _, p := range profs { fmt.Printf("profession: %v\n", p) }
	var tools []map[string]interface{}
	database.GetDB().Raw(`SELECT id, name, access_level FROM tools WHERE access_level > 0 LIMIT 5`).Scan(&tools)
	for _, t := range tools { fmt.Printf("tool: %v\n", t) }
	var chs []map[string]interface{}
	database.GetDB().Raw(`SELECT id, path_id, title, access_level FROM path_chapters WHERE access_level > 0 LIMIT 5`).Scan(&chs)
	for _, c := range chs { fmt.Printf("chapter: %v\n", c) }
}
