// 把指定 email 的用户设为会员 (is_member=1)
// 用法：go run cmd/promote-member/main.go testuser1@example.com
package main

import (
	"fmt"
	"log"
	"os"

	"hot-ai-backend/internal/database"
)

func main() {
	email := ""
	if len(os.Args) > 1 {
		email = os.Args[1]
	} else {
		log.Fatal("usage: go run cmd/promote-member/main.go <email>")
	}

	dsn := "root:shds@Admin123@tcp(192.168.50.109:13306)/hot_ai?charset=utf8mb4&parseTime=true&loc=Local"
	cfg, err := database.ParseDSN(dsn)
	if err != nil {
		log.Fatalf("parse DSN: %v", err)
	}
	if err := database.InitDB(*cfg); err != nil {
		log.Fatalf("init DB: %v", err)
	}

	r := database.GetDB().Exec(`UPDATE users SET is_member = 1 WHERE email = ?`, email)
	if r.Error != nil {
		log.Fatalf("update: %v", r.Error)
	}
	if r.RowsAffected == 0 {
		log.Fatalf("no user with email=%s", email)
	}
	fmt.Printf("[ok] promoted %s to member (%d row)\n", email, r.RowsAffected)
}
