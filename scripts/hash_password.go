package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run hash_password.go <password>")
		fmt.Println("Example: go run hash_password.go admin123")
		os.Exit(1)
	}

	password := os.Args[1]

	// Hash password với bcrypt (cost 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=================================")
	fmt.Println("Password Hash Generator")
	fmt.Println("=================================")
	fmt.Printf("Original password: %s\n", password)
	fmt.Printf("Hashed password: %s\n", string(hashedPassword))
	fmt.Println("=================================")
	fmt.Println("\nSQL để tạo user:")
	fmt.Printf("INSERT INTO users (username, password, role, created_at, updated_at)\n")
	fmt.Printf("VALUES ('admin', '%s', 'admin', NOW(), NOW());\n", string(hashedPassword))
}
