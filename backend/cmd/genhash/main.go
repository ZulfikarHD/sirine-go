package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "make":
		makeHash()
	case "check":
		checkHash()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Hash Utility for Password Management")
	fmt.Println("\nUsage:")
	fmt.Println("  go run cmd/hash/main.go <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  make <password>              Generate bcrypt hash untuk password")
	fmt.Println("  check <password> <hash>      Verify apakah password match dengan hash")
	fmt.Println("\nExamples:")
	fmt.Println("  go run cmd/hash/main.go make \"MyPassword123\"")
	fmt.Println("  go run cmd/hash/main.go check \"MyPassword123\" \"$2a$12$...\"")
}

// makeHash generates bcrypt hash untuk password
func makeHash() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Password tidak diberikan")
		fmt.Println("\nUsage: go run cmd/hash/main.go make <password>")
		os.Exit(1)
	}

	password := os.Args[2]

	if password == "" {
		fmt.Println("Error: Password tidak boleh kosong")
		os.Exit(1)
	}

	// Generate bcrypt hash dengan cost 12
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Fatal("Failed to generate hash:", err)
	}

	fmt.Println("\n✅ Hash generated successfully!")
	fmt.Printf("\nPassword: %s\n", password)
	fmt.Printf("Hash:     %s\n", string(hash))
	fmt.Println("\nCopy hash di atas untuk digunakan di seeder atau database.")
}

// checkHash verifies apakah password match dengan hash
func checkHash() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Password atau hash tidak lengkap")
		fmt.Println("\nUsage: go run cmd/hash/main.go check <password> <hash>")
		os.Exit(1)
	}

	password := os.Args[2]
	hash := os.Args[3]

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	
	if err == nil {
		fmt.Println("\n✅ Password MATCH!")
		fmt.Printf("\nPassword: %s\n", password)
		fmt.Printf("Hash:     %s\n", hash)
		fmt.Println("\nVerifikasi berhasil - password sesuai dengan hash.")
	} else {
		fmt.Println("\n❌ Password TIDAK MATCH!")
		fmt.Printf("\nPassword: %s\n", password)
		fmt.Printf("Hash:     %s\n", hash)
		fmt.Println("\nVerifikasi gagal - password tidak sesuai dengan hash.")
		os.Exit(1)
	}
}
