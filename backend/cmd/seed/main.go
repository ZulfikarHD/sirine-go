package main

import (
	"log"
	"time"
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
	seedProductionOrders()

	log.Println("✅ Database seeding completed!")
}

// seedAdminUser membuat admin user default
func seedAdminUser() {
	db := database.GetDB()

	// Check apakah admin sudah ada (check by NIP or email)
	var existingAdmin models.User
	if err := db.Where("nip = ? OR email = ?", "99999", "admin@sirine.local").First(&existingAdmin).Error; err == nil {
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
		log.Printf("Warning: Failed to seed admin user (might already exist): %v", err)
		return
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

// seedProductionOrders membuat sample Production Orders untuk testing
func seedProductionOrders() {
	db := database.GetDB()

	// Check apakah PO sudah ada
	var count int64
	db.Model(&models.ProductionOrder{}).Count(&count)
	if count > 0 {
		log.Printf("ℹ️  Production Orders sudah ada (%d records), skip seeding PO", count)
		return
	}

	// Sample PO data dengan berbagai status dan priority
	samplePOs := []models.ProductionOrder{
		{
			PONumber:                  1001,
			OBCNumber:                 "OBC000001",
			SAPCustomerCode:           "CUST-001",
			SAPProductCode:            "PROD-A123",
			ProductName:               "Kertas HVS A4 80gsm",
			QuantityOrdered:           50000,
			QuantityTargetLembarBesar: 1000,
			EstimatedRims:             100,
			OrderDate:                 time.Now().AddDate(0, 0, -5),
			DueDate:                   time.Now().AddDate(0, 0, 2), // 2 hari dari sekarang - URGENT
			Priority:                  models.PriorityUrgent,
			PriorityScore:             0, // akan di-calculate
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
			Notes:                     "PO urgent untuk customer priority",
		},
		{
			PONumber:                  1002,
			OBCNumber:                 "OBC000002",
			SAPCustomerCode:           "CUST-002",
			SAPProductCode:            "PROD-B456",
			ProductName:               "Kertas Art Paper A3 120gsm",
			QuantityOrdered:           30000,
			QuantityTargetLembarBesar: 600,
			EstimatedRims:             60,
			OrderDate:                 time.Now().AddDate(0, 0, -3),
			DueDate:                   time.Now().AddDate(0, 0, 7), // 7 hari dari sekarang
			Priority:                  models.PriorityNormal,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
		},
		{
			PONumber:                  1003,
			OBCNumber:                 "OBC000003",
			SAPCustomerCode:           "CUST-003",
			SAPProductCode:            "PROD-C789",
			ProductName:               "Kertas Duplex A4 250gsm",
			QuantityOrdered:           20000,
			QuantityTargetLembarBesar: 400,
			EstimatedRims:             40,
			OrderDate:                 time.Now().AddDate(0, 0, -7),
			DueDate:                   time.Now().AddDate(0, 0, 14), // 14 hari dari sekarang
			Priority:                  models.PriorityLow,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
		},
		{
			PONumber:                  1004,
			OBCNumber:                 "OBC000004",
			SAPCustomerCode:           "CUST-001",
			SAPProductCode:            "PROD-D012",
			ProductName:               "Kertas Glossy A4 150gsm",
			QuantityOrdered:           40000,
			QuantityTargetLembarBesar: 800,
			EstimatedRims:             80,
			OrderDate:                 time.Now().AddDate(0, 0, -2),
			DueDate:                   time.Now().AddDate(0, 0, 3), // 3 hari dari sekarang
			Priority:                  models.PriorityUrgent,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
			Notes:                     "Rush order - priority tinggi",
		},
		{
			PONumber:                  1005,
			OBCNumber:                 "OBC000005",
			SAPCustomerCode:           "CUST-004",
			SAPProductCode:            "PROD-E345",
			ProductName:               "Kertas Manila A4 100gsm",
			QuantityOrdered:           25000,
			QuantityTargetLembarBesar: 500,
			EstimatedRims:             50,
			OrderDate:                 time.Now().AddDate(0, 0, -4),
			DueDate:                   time.Now().AddDate(0, 0, 10), // 10 hari dari sekarang
			Priority:                  models.PriorityNormal,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
		},
		{
			PONumber:                  1006,
			OBCNumber:                 "OBC000006",
			SAPCustomerCode:           "CUST-005",
			SAPProductCode:            "PROD-F678",
			ProductName:               "Kertas Ivory A3 200gsm",
			QuantityOrdered:           15000,
			QuantityTargetLembarBesar: 300,
			EstimatedRims:             30,
			OrderDate:                 time.Now().AddDate(0, 0, -1),
			DueDate:                   time.Now().AddDate(0, 0, 5), // 5 hari dari sekarang
			Priority:                  models.PriorityNormal,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
		},
		{
			PONumber:                  1007,
			OBCNumber:                 "OBC000007",
			SAPCustomerCode:           "CUST-002",
			SAPProductCode:            "PROD-G901",
			ProductName:               "Kertas Matte A4 180gsm",
			QuantityOrdered:           35000,
			QuantityTargetLembarBesar: 700,
			EstimatedRims:             70,
			OrderDate:                 time.Now().AddDate(0, 0, -6),
			DueDate:                   time.Now().AddDate(0, 0, 1), // 1 hari dari sekarang - SANGAT URGENT
			Priority:                  models.PriorityUrgent,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
			Notes:                     "Deadline ketat - perlu perhatian segera",
		},
		{
			PONumber:                  1008,
			OBCNumber:                 "OBC000008",
			SAPCustomerCode:           "CUST-006",
			SAPProductCode:            "PROD-H234",
			ProductName:               "Kertas Concorde A4 220gsm",
			QuantityOrdered:           10000,
			QuantityTargetLembarBesar: 200,
			EstimatedRims:             20,
			OrderDate:                 time.Now().AddDate(0, 0, -8),
			DueDate:                   time.Now().AddDate(0, 0, 20), // 20 hari dari sekarang
			Priority:                  models.PriorityLow,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
		},
	}

	// Insert PO data dan calculate priority score
	for i := range samplePOs {
		samplePOs[i].UpdatePriorityScore()

		if err := db.Create(&samplePOs[i]).Error; err != nil {
			log.Printf("Warning: Failed to seed PO %d: %v", samplePOs[i].PONumber, err)
			continue
		}

		// Create corresponding Khazwal Material Preparation record
		khazwalPrep := models.KhazwalMaterialPreparation{
			ProductionOrderID:    samplePOs[i].ID,
			SAPPlatCode:          "PLAT-" + samplePOs[i].SAPProductCode,
			KertasBlankoQuantity: samplePOs[i].QuantityTargetLembarBesar,
			TintaRequirements:    []byte(`{"cyan": 2, "magenta": 2, "yellow": 2, "black": 3}`),
			Status:               models.MaterialPrepPending,
		}

		if err := db.Create(&khazwalPrep).Error; err != nil {
			log.Printf("Warning: Failed to seed Khazwal Prep for PO %d: %v", samplePOs[i].PONumber, err)
		}

		log.Printf("✅ PO seeded: %d - %s (Priority: %s, Score: %d)", 
			samplePOs[i].PONumber, 
			samplePOs[i].ProductName, 
			samplePOs[i].Priority,
			samplePOs[i].PriorityScore)
	}

	log.Println("\n📦 Production Orders Summary:")
	log.Printf("   Total POs: %d", len(samplePOs))
	log.Println("   Status: WAITING_MATERIAL_PREP")
	log.Println("   Stage: KHAZWAL_MATERIAL_PREP")
}
