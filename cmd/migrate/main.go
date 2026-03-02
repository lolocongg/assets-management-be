package main

import (
	"fmt"
	"os"

	"github.com/davidcm146/assets-management-be.git/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	if len(os.Args) < 2 {
		panic("please provide 'up' or 'down' as argument")
	}
	cmd := os.Args[1]
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	migrationDir := "./internal/database/migrations"

	switch cmd {
	case "up":
		if err := database.MigrateUp(conn, migrationDir); err != nil {
			panic(fmt.Sprintf("migration up failed: %v", err))
		}
		fmt.Println("Migration up completed.")
	case "down":
		if err := database.MigrateDown(conn, migrationDir); err != nil {
			panic(fmt.Sprintf("migration down failed: %v", err))
		}
		fmt.Println("Migration down completed.")
	default:
		panic("unknown command, use 'up' or 'down'")
	}
}
