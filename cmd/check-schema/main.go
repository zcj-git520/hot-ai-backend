package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"hot-ai-backend/internal/database"
)

func main() {
	dsn := "root:shds@Admin123@tcp(192.168.50.109:13306)/hot_ai?charset=utf8mb4&parseTime=true&loc=Local"
	cfg, _ := database.ParseDSN(dsn)
	database.InitDB(*cfg)
	db := database.GetDB()

	pw, _ := bcrypt.GenerateFromPassword([]byte("Test1234!"), bcrypt.DefaultCost)

	// Reset password for perm-test@example.com (member) and 12345@qq.com (normal)
	res := db.Exec(`UPDATE users SET password_hash = ? WHERE email IN (?, ?)`,
		string(pw), "perm-test@example.com", "12345@qq.com")
	fmt.Printf("Updated %d users\n", res.RowsAffected)

	var users []map[string]interface{}
	db.Raw(`SELECT email, is_member FROM users WHERE email IN (?, ?)`,
		"perm-test@example.com", "12345@qq.com").Scan(&users)
	for _, u := range users {
		fmt.Printf("  %v\n", u)
	}
}