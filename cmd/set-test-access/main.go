package main

import (
	"fmt"

	"hot-ai-backend/internal/database"
)

func main() {
	dsn := "root:shds@Admin123@tcp(192.168.50.109:13306)/hot_ai?charset=utf8mb4&parseTime=true&loc=Local"
	cfg, _ := database.ParseDSN(dsn)
	database.InitDB(*cfg)
	db := database.GetDB()

	// Tools: reset, then set tier 1 (logged-in) and tier 2 (member)
	db.Exec(`UPDATE tools SET access_level = 0`)
	db.Exec(`UPDATE tools SET access_level = 1 WHERE id IN (1, 3)`)    // ChatGPT, Notion AI
	db.Exec(`UPDATE tools SET access_level = 2 WHERE id IN (2, 4, 5)`) // Claude, Jasper, Copy.ai

	// Professions: id=1->1, id=2->2, id=3->2
	db.Exec(`UPDATE professions SET access_level = 0`)
	db.Exec(`UPDATE professions SET access_level = 1 WHERE id = 1`)
	db.Exec(`UPDATE professions SET access_level = 2 WHERE id IN (2, 3)`)

	// Learning paths: path 1 -> 0, path 2 -> 1, path 3+ -> 2
	db.Exec(`UPDATE learning_paths SET access_level = 0`)
	db.Exec(`UPDATE learning_paths SET access_level = 1 WHERE id = 2`)
	db.Exec(`UPDATE learning_paths SET access_level = 2 WHERE id IN (3, 4, 5)`)

	// Path chapters: path 1 chapters by order_index
	db.Exec(`UPDATE path_chapters SET access_level = 0`)
	db.Exec(`UPDATE path_chapters SET access_level = 1 WHERE path_id = 1 AND order_index = 2`)
	db.Exec(`UPDATE path_chapters SET access_level = 2 WHERE path_id = 1 AND order_index = 3`)

	fmt.Println("=== tools (by access_level DESC) ===")
	var tools []map[string]interface{}
	db.Raw(`SELECT id, name, access_level FROM tools ORDER BY access_level DESC, id LIMIT 10`).Scan(&tools)
	for _, t := range tools {
		fmt.Printf("  id=%v %-25s level=%v\n", t["id"], t["name"], t["access_level"])
	}

	fmt.Println("\n=== professions ===")
	var profs []map[string]interface{}
	db.Raw(`SELECT id, name, access_level FROM professions ORDER BY id LIMIT 5`).Scan(&profs)
	for _, t := range profs {
		fmt.Printf("  id=%v %-25s level=%v\n", t["id"], t["name"], t["access_level"])
	}

	fmt.Println("\n=== learning paths ===")
	var paths []map[string]interface{}
	db.Raw(`SELECT id, title, access_level FROM learning_paths ORDER BY id LIMIT 5`).Scan(&paths)
	for _, t := range paths {
		fmt.Printf("  id=%v %-25s level=%v\n", t["id"], t["title"], t["access_level"])
	}

	fmt.Println("\n=== path_chapters (path 1) ===")
	var chs []map[string]interface{}
	db.Raw(`SELECT id, order_index, title, access_level FROM path_chapters WHERE path_id=1 ORDER BY order_index`).Scan(&chs)
	for _, t := range chs {
		fmt.Printf("  id=%v ord=%v %-30s level=%v\n", t["id"], t["order_index"], t["title"], t["access_level"])
	}
}