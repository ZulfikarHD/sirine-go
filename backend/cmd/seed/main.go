package main

import (
	"log"
	"sirine-go/backend/config"
	"sirine-go/backend/database"
	"sirine-go/backend/models"
	"sirine-go/backend/services"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables dari file .env
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file tidak ditemukan, menggunakan default values")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("🌱 Starting database seeding...")

	// Seed data
	seedAdminUser()
	seedDemoUsers()

	log.Println("✅ Database seeding completed!")
}

// seedAdminUser membuat admin user default
func seedAdminUser() {
	db := database.GetDB()

	// Check apakah admin sudah ada
	var existingAdmin models.User
	if err := db.Where("nip = ?", "99999").First(&existingAdmin).Error; err == nil {
		log.Println("ℹ️  Admin user sudah ada, skip seeding admin")
		return
	}

	// Generate password hash secara dynamic
	passwordService := services.NewPasswordService()
	adminPassword := "Admin@123"
	adminPasswordHash, err := passwordService.HashPassword(adminPassword)
	if err != nil {
		log.Fatal("Failed to hash admin password:", err)
	}

	adminUser := models.User{
		NIP:                "99999",
		FullName:           "Administrator",
		Email:              "admin@sirine.local",
		Phone:              "081234567890",
		PasswordHash:       adminPasswordHash,
		Role:               models.RoleAdmin,
		Department:         models.DeptKhazwal,
		Shift:              models.ShiftPagi,
		Status:             models.StatusActive,
		MustChangePassword: false, // Admin tidak perlu ganti password
	}

	if err := db.Create(&adminUser).Error; err != nil {
		log.Fatal("Failed to seed admin user:", err)
	}

	log.Println("✅ Admin user seeded successfully")
	log.Println("   NIP: 99999")
	log.Println("   Email: admin@sirine.local")
	log.Println("   Password: Admin@123")
}

// seedDemoUsers membuat beberapa demo users untuk testing
func seedDemoUsers() {
	db := database.GetDB()

	// Generate password hash secara dynamic
	passwordService := services.NewPasswordService()
	demoPassword := "Demo@123"
	demoPasswordHash, err := passwordService.HashPassword(demoPassword)
	if err != nil {
		log.Fatal("Failed to hash demo password:", err)
	}

	demoUsers := []models.User{
		{
			NIP:          "10001",
			FullName:     "Manager Produksi",
			Email:        "manager@sirine.local",
			Phone:        "081234567891",
			PasswordHash: demoPasswordHash,
			Role:         models.RoleManager,
			Department:   models.DeptKhazwal,
			Shift:        models.ShiftPagi,
			Status:       models.StatusActive,
		},
		{
			NIP:          "20001",
			FullName:     "Staff Khazanah Awal",
			Email:        "khazwal@sirine.local",
			Phone:        "081234567892",
			PasswordHash: demoPasswordHash,
			Role:         models.RoleStaffKhazwal,
			Department:   models.DeptKhazwal,
			Shift:        models.ShiftPagi,
			Status:       models.StatusActive,
		},
		{
			NIP:          "30001",
			FullName:     "Operator Cetak",
			Email:        "cetak@sirine.local",
			Phone:        "081234567893",
			PasswordHash: demoPasswordHash,
			Role:         models.RoleOperatorCetak,
			Department:   models.DeptCetak,
			Shift:        models.ShiftSiang,
			Status:       models.StatusActive,
		},
		{
			NIP:          "40001",
			FullName:     "QC Inspector",
			Email:        "qc@sirine.local",
			Phone:        "081234567894",
			PasswordHash: demoPasswordHash,
			Role:         models.RoleQCInspector,
			Department:   models.DeptVerifikasi,
			Shift:        models.ShiftPagi,
			Status:       models.StatusActive,
		},
		{
			NIP:          "50001",
			FullName:     "Verifikator",
			Email:        "verifikator@sirine.local",
			Phone:        "081234567895",
			PasswordHash: demoPasswordHash,
			Role:         models.RoleVerifikator,
			Department:   models.DeptVerifikasi,
			Shift:        models.ShiftPagi,
			Status:       models.StatusActive,
		},
		{
			NIP:          "60001",
			FullName:     "Staff Khazanah Akhir",
			Email:        "khazkhir@sirine.local",
			Phone:        "081234567896",
			PasswordHash: demoPasswordHash,
			Role:         models.RoleStaffKhazkhir,
			Department:   models.DeptKhazkhir,
			Shift:        models.ShiftSiang,
			Status:       models.StatusActive,
		},
	}

	for _, user := range demoUsers {
		// Check apakah user sudah ada
		var existing models.User
		if err := db.Where("nip = ?", user.NIP).First(&existing).Error; err == nil {
			log.Printf("ℹ️  User %s (%s) sudah ada, skip", user.NIP, user.FullName)
			continue
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Warning: Failed to seed user %s: %v", user.NIP, err)
			continue
		}

		log.Printf("✅ Demo user seeded: %s - %s", user.NIP, user.FullName)
	}

	log.Println("\n📝 Demo users credentials:")
	log.Println("   Password untuk semua demo users: Demo@123")
}
