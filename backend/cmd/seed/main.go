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

	// Seed data in correct order (respecting FK constraints)
	seedAdminUser()
	seedDemoUsers()
	seedOBCMasters() // Must be seeded before ProductionOrders
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

// seedOBCMasters membuat sample OBC Master records
func seedOBCMasters() {
	db := database.GetDB()

	// Check apakah OBC Master sudah ada
	var count int64
	db.Model(&models.OBCMaster{}).Count(&count)
	if count > 0 {
		log.Printf("ℹ️  OBC Masters sudah ada (%d records), skip seeding OBC Masters", count)
		return
	}

	now := time.Now()
	sampleOBCMasters := []models.OBCMaster{
		{
			OBCNumber:           "OBC000001",
			OBCDate:             &now,
			Material:            "HVS",
			Seri:                "A4-80",
			Warna:               "Putih",
			FactoryCode:         "FAC-001",
			QuantityOrdered:     50000,
			DueDate:             timePtr(time.Now().AddDate(0, 0, 2)),
			PlatNumber:          "PLAT-PROD-A123",
			MaterialDescription: "Kertas HVS A4 80gsm",
			BaseUnit:            "RIM",
		},
		{
			OBCNumber:           "OBC000002",
			OBCDate:             &now,
			Material:            "ART-PAPER",
			Seri:                "A3-120",
			Warna:               "Putih",
			FactoryCode:         "FAC-001",
			QuantityOrdered:     30000,
			DueDate:             timePtr(time.Now().AddDate(0, 0, 7)),
			PlatNumber:          "PLAT-PROD-B456",
			MaterialDescription: "Kertas Art Paper A3 120gsm",
			BaseUnit:            "RIM",
		},
		{
			OBCNumber:           "OBC000003",
			OBCDate:             &now,
			Material:            "DUPLEX",
			Seri:                "A4-250",
			Warna:               "Putih",
			FactoryCode:         "FAC-002",
			QuantityOrdered:     20000,
			DueDate:             timePtr(time.Now().AddDate(0, 0, 14)),
			PlatNumber:          "PLAT-PROD-C789",
			MaterialDescription: "Kertas Duplex A4 250gsm",
			BaseUnit:            "RIM",
		},
		{
			OBCNumber:           "OBC000004",
			OBCDate:             &now,
			Material:            "GLOSSY",
			Seri:                "A4-150",
			Warna:               "Putih",
			FactoryCode:         "FAC-001",
			QuantityOrdered:     40000,
			DueDate:             timePtr(time.Now().AddDate(0, 0, 3)),
			PlatNumber:          "PLAT-PROD-D012",
			MaterialDescription: "Kertas Glossy A4 150gsm",
			BaseUnit:            "RIM",
		},
		{
			OBCNumber:           "OBC000005",
			OBCDate:             &now,
			Material:            "MANILA",
			Seri:                "A4-100",
			Warna:               "Kuning",
			FactoryCode:         "FAC-002",
			QuantityOrdered:     25000,
			DueDate:             timePtr(time.Now().AddDate(0, 0, 10)),
			PlatNumber:          "PLAT-PROD-E345",
			MaterialDescription: "Kertas Manila A4 100gsm",
			BaseUnit:            "RIM",
		},
		{
			OBCNumber:           "OBC000006",
			OBCDate:             &now,
			Material:            "IVORY",
			Seri:                "A3-200",
			Warna:               "Putih",
			FactoryCode:         "FAC-001",
			QuantityOrdered:     15000,
			DueDate:             timePtr(time.Now().AddDate(0, 0, 5)),
			PlatNumber:          "PLAT-PROD-F678",
			MaterialDescription: "Kertas Ivory A3 200gsm",
			BaseUnit:            "RIM",
		},
		{
			OBCNumber:           "OBC000007",
			OBCDate:             &now,
			Material:            "MATTE",
			Seri:                "A4-180",
			Warna:               "Putih",
			FactoryCode:         "FAC-002",
			QuantityOrdered:     35000,
			DueDate:             timePtr(time.Now().AddDate(0, 0, 1)),
			PlatNumber:          "PLAT-PROD-G901",
			MaterialDescription: "Kertas Matte A4 180gsm",
			BaseUnit:            "RIM",
		},
		{
			OBCNumber:           "OBC000008",
			OBCDate:             &now,
			Material:            "CONCORDE",
			Seri:                "A4-220",
			Warna:               "Putih",
			FactoryCode:         "FAC-001",
			QuantityOrdered:     10000,
			DueDate:             timePtr(time.Now().AddDate(0, 0, 20)),
			PlatNumber:          "PLAT-PROD-H234",
			MaterialDescription: "Kertas Concorde A4 220gsm",
			BaseUnit:            "RIM",
		},
	}

	for i := range sampleOBCMasters {
		if err := db.Create(&sampleOBCMasters[i]).Error; err != nil {
			log.Printf("Warning: Failed to seed OBC Master %s: %v", sampleOBCMasters[i].OBCNumber, err)
			continue
		}

		log.Printf("✅ OBC Master seeded: %s - %s", 
			sampleOBCMasters[i].OBCNumber, 
			sampleOBCMasters[i].MaterialDescription)
	}

	log.Printf("\n📦 OBC Masters Summary: Total %d records seeded", len(sampleOBCMasters))
}

// timePtr helper function untuk create *time.Time
func timePtr(t time.Time) *time.Time {
	return &t
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

	// Get OBC Masters untuk reference
	var obcMasters []models.OBCMaster
	if err := db.Order("id ASC").Find(&obcMasters).Error; err != nil {
		log.Fatal("Failed to fetch OBC Masters for PO seeding:", err)
	}

	if len(obcMasters) == 0 {
		log.Fatal("No OBC Masters found. Please run seedOBCMasters first.")
	}

	// Sample PO data dengan berbagai status dan priority
	// Setiap PO di-link ke OBC Master yang sesuai
	samplePOs := []models.ProductionOrder{
		{
			PONumber:                  1001,
			OBCMasterID:               obcMasters[0].ID, // Link to OBC000001
			OBCNumber:                 obcMasters[0].OBCNumber,
			SAPCustomerCode:           "CUST-001",
			SAPProductCode:            "PROD-A123",
			ProductName:               obcMasters[0].MaterialDescription,
			QuantityOrdered:           50000,
			QuantityTargetLembarBesar: 1000,
			EstimatedRims:             100,
			OrderDate:                 time.Now().AddDate(0, 0, -5),
			DueDate:                   time.Now().AddDate(0, 0, 2),
			Priority:                  models.PriorityUrgent,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
			Notes:                     "PO urgent untuk customer priority",
		},
		{
			PONumber:                  1002,
			OBCMasterID:               obcMasters[1].ID, // Link to OBC000002
			OBCNumber:                 obcMasters[1].OBCNumber,
			SAPCustomerCode:           "CUST-002",
			SAPProductCode:            "PROD-B456",
			ProductName:               obcMasters[1].MaterialDescription,
			QuantityOrdered:           30000,
			QuantityTargetLembarBesar: 600,
			EstimatedRims:             60,
			OrderDate:                 time.Now().AddDate(0, 0, -3),
			DueDate:                   time.Now().AddDate(0, 0, 7),
			Priority:                  models.PriorityNormal,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
		},
		{
			PONumber:                  1003,
			OBCMasterID:               obcMasters[2].ID, // Link to OBC000003
			OBCNumber:                 obcMasters[2].OBCNumber,
			SAPCustomerCode:           "CUST-003",
			SAPProductCode:            "PROD-C789",
			ProductName:               obcMasters[2].MaterialDescription,
			QuantityOrdered:           20000,
			QuantityTargetLembarBesar: 400,
			EstimatedRims:             40,
			OrderDate:                 time.Now().AddDate(0, 0, -7),
			DueDate:                   time.Now().AddDate(0, 0, 14),
			Priority:                  models.PriorityLow,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
		},
		{
			PONumber:                  1004,
			OBCMasterID:               obcMasters[3].ID, // Link to OBC000004
			OBCNumber:                 obcMasters[3].OBCNumber,
			SAPCustomerCode:           "CUST-001",
			SAPProductCode:            "PROD-D012",
			ProductName:               obcMasters[3].MaterialDescription,
			QuantityOrdered:           40000,
			QuantityTargetLembarBesar: 800,
			EstimatedRims:             80,
			OrderDate:                 time.Now().AddDate(0, 0, -2),
			DueDate:                   time.Now().AddDate(0, 0, 3),
			Priority:                  models.PriorityUrgent,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
			Notes:                     "Rush order - priority tinggi",
		},
		{
			PONumber:                  1005,
			OBCMasterID:               obcMasters[4].ID, // Link to OBC000005
			OBCNumber:                 obcMasters[4].OBCNumber,
			SAPCustomerCode:           "CUST-004",
			SAPProductCode:            "PROD-E345",
			ProductName:               obcMasters[4].MaterialDescription,
			QuantityOrdered:           25000,
			QuantityTargetLembarBesar: 500,
			EstimatedRims:             50,
			OrderDate:                 time.Now().AddDate(0, 0, -4),
			DueDate:                   time.Now().AddDate(0, 0, 10),
			Priority:                  models.PriorityNormal,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
		},
		{
			PONumber:                  1006,
			OBCMasterID:               obcMasters[5].ID, // Link to OBC000006
			OBCNumber:                 obcMasters[5].OBCNumber,
			SAPCustomerCode:           "CUST-005",
			SAPProductCode:            "PROD-F678",
			ProductName:               obcMasters[5].MaterialDescription,
			QuantityOrdered:           15000,
			QuantityTargetLembarBesar: 300,
			EstimatedRims:             30,
			OrderDate:                 time.Now().AddDate(0, 0, -1),
			DueDate:                   time.Now().AddDate(0, 0, 5),
			Priority:                  models.PriorityNormal,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
		},
		{
			PONumber:                  1007,
			OBCMasterID:               obcMasters[6].ID, // Link to OBC000007
			OBCNumber:                 obcMasters[6].OBCNumber,
			SAPCustomerCode:           "CUST-002",
			SAPProductCode:            "PROD-G901",
			ProductName:               obcMasters[6].MaterialDescription,
			QuantityOrdered:           35000,
			QuantityTargetLembarBesar: 700,
			EstimatedRims:             70,
			OrderDate:                 time.Now().AddDate(0, 0, -6),
			DueDate:                   time.Now().AddDate(0, 0, 1),
			Priority:                  models.PriorityUrgent,
			PriorityScore:             0,
			CurrentStage:              models.StageKhazwalMaterialPrep,
			CurrentStatus:             models.StatusWaitingMaterialPrep,
			Notes:                     "Deadline ketat - perlu perhatian segera",
		},
		{
			PONumber:                  1008,
			OBCMasterID:               obcMasters[7].ID, // Link to OBC000008
			OBCNumber:                 obcMasters[7].OBCNumber,
			SAPCustomerCode:           "CUST-006",
			SAPProductCode:            "PROD-H234",
			ProductName:               obcMasters[7].MaterialDescription,
			QuantityOrdered:           10000,
			QuantityTargetLembarBesar: 200,
			EstimatedRims:             20,
			OrderDate:                 time.Now().AddDate(0, 0, -8),
			DueDate:                   time.Now().AddDate(0, 0, 20),
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

		// Get the corresponding OBC Master untuk plat code
		var obcMaster models.OBCMaster
		db.First(&obcMaster, samplePOs[i].OBCMasterID)

		// Create corresponding Khazwal Material Preparation record
		khazwalPrep := models.KhazwalMaterialPreparation{
			ProductionOrderID:    samplePOs[i].ID,
			SAPPlatCode:          obcMaster.PlatNumber, // Use actual plat number from OBC Master
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
