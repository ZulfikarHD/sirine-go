package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sirine-go/backend/config"
	"sirine-go/backend/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables dari file .env
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file tidak ditemukan, menggunakan default values")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Check command argument
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		migrateUp(cfg)
	case "down":
		migrateDown(cfg)
	case "fresh":
		migrateFresh(cfg)
	case "create":
		createDatabase(cfg)
	case "drop":
		dropDatabase(cfg)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run cmd/migrate/main.go [command]")
	fmt.Println("\nCommands:")
	fmt.Println("  create  - Buat database")
	fmt.Println("  drop    - Hapus database")
	fmt.Println("  up      - Jalankan migrations")
	fmt.Println("  down    - Rollback migrations")
	fmt.Println("  fresh   - Drop database, create ulang, dan migrate")
}

// createDatabase membuat database jika belum ada
func createDatabase(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}
	defer db.Close()

	// Buat database dengan character set dan collation yang benar
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}

	log.Printf("✅ Database '%s' created successfully!", cfg.DBName)
}

// dropDatabase menghapus database
func dropDatabase(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}
	defer db.Close()

	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.DBName)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Failed to drop database:", err)
	}

	log.Printf("✅ Database '%s' dropped successfully!", cfg.DBName)
}

// migrateUp menjalankan migrations dengan GORM AutoMigrate
func migrateUp(cfg *config.Config) {
	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Running migrations...")

	// Get models registry
	registry := database.NewModelsRegistry()
	
	// Auto migrate all registered models
	if err := database.AutoMigrate(registry.GetModels()...); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Printf("✅ Migrations completed! (%d tables migrated)", registry.GetTableCount())
}

// migrateDown rollback migrations (drop all tables)
func migrateDown(cfg *config.Config) {
	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Rolling back migrations...")

	db := database.GetDB()

	// Get models registry
	registry := database.NewModelsRegistry()
	
	// Drop tables dalam urutan terbalik untuk menghindari foreign key constraints
	tables := registry.GetTablesForRollback()

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Printf("Warning: Failed to drop table %s: %v", table, err)
		} else {
			log.Printf("Dropped table: %s", table)
		}
	}

	log.Printf("✅ Rollback completed! (%d tables dropped)", registry.GetTableCount())
}

// migrateFresh drop database, create ulang, dan migrate
func migrateFresh(cfg *config.Config) {
	log.Println("🔄 Starting fresh migration...")
	
	dropDatabase(cfg)
	createDatabase(cfg)
	migrateUp(cfg)
	
	log.Println("✅ Fresh migration completed!")
}
