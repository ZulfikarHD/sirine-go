package main

import (
	"fmt"
	"log"
	"sirine-go/backend/config"
	"sirine-go/backend/database"
	"sirine-go/backend/models"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

/**
 * Script untuk migrate existing ProductionOrder data ke OBCMaster architecture
 * yang mencakup:
 * 1. Check existing production_orders data
 * 2. Create temporary OBCMaster records dari existing OBC data
 * 3. Update production_orders.obc_master_id dengan link ke OBCMaster
 * 4. Add foreign key constraint
 */

func main() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file tidak ditemukan, menggunakan default values")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := database.GetDB()

	log.Println("🔄 Starting OBC data migration...")

	// Step 1: Check if there's existing data
	var poCount int64
	if err := db.Model(&models.ProductionOrder{}).Count(&poCount).Error; err != nil {
		log.Fatal("Failed to count production orders:", err)
	}

	log.Printf("Found %d production orders in database", poCount)

	if poCount == 0 {
		log.Println("✅ No existing data to migrate. Adding foreign key constraint...")
		
		// Add foreign key constraint
		if err := addForeignKeyConstraint(db); err != nil {
			log.Fatal("Failed to add foreign key constraint:", err)
		}
		
		log.Println("✅ Migration completed successfully!")
		return
	}

	// Step 2: Check if migration is needed (obc_master_id column exists but is 0)
	var posWithoutOBC int64
	if err := db.Model(&models.ProductionOrder{}).Where("obc_master_id = 0").Count(&posWithoutOBC).Error; err != nil {
		log.Printf("Warning: Could not check obc_master_id column: %v", err)
		log.Println("This might mean the column doesn't exist yet or has no default value.")
		log.Println("Since we have existing POs without OBC links, we'll create a default OBC Master.")
	}

	log.Printf("Found %d production orders without OBC Master link", posWithoutOBC)

	// Step 3: Create a default OBCMaster untuk existing POs yang tidak punya OBC data
	log.Println("Creating default OBCMaster for existing production orders...")
	
	defaultOBC := &models.OBCMaster{
		OBCNumber:           "OBC-LEGACY-MIGRATION",
		Material:            "LEGACY",
		Seri:                "MIGRATION",
		Warna:               "N/A",
		FactoryCode:         "LEGACY",
		QuantityOrdered:     0,
		MaterialDescription: "Legacy data dari migration - belum ada OBC Master",
		Personalization:     "non Perso",
	}

	// Create atau get existing default OBC
	var existingDefault models.OBCMaster
	result := db.Where("obc_number = ?", defaultOBC.OBCNumber).First(&existingDefault)
	
	if result.Error != nil {
		// Create new default OBC
		if err := db.Create(defaultOBC).Error; err != nil {
			log.Fatal("Failed to create default OBC Master:", err)
		}
		log.Printf("✅ Created default OBC Master with ID: %d", defaultOBC.ID)
	} else {
		// Use existing default
		defaultOBC = &existingDefault
		log.Printf("✅ Using existing default OBC Master with ID: %d", defaultOBC.ID)
	}

	// Step 4: Update all production_orders dengan default OBC Master ID
	log.Println("Updating production orders with default OBC Master link...")
	
	// Update semua PO yang belum punya obc_master_id
	result = db.Model(&models.ProductionOrder{}).
		Where("obc_master_id = 0 OR obc_master_id IS NULL").
		Update("obc_master_id", defaultOBC.ID)
	
	if result.Error != nil {
		log.Fatal("Failed to update production orders:", result.Error)
	}
	
	log.Printf("✅ Updated %d production orders with default OBC Master", result.RowsAffected)

	// Step 5: Add foreign key constraint
	log.Println("Adding foreign key constraint...")
	
	if err := addForeignKeyConstraint(db); err != nil {
		log.Printf("Warning: Failed to add foreign key constraint: %v", err)
		log.Println("You may need to add it manually or check if it already exists.")
	} else {
		log.Println("✅ Foreign key constraint added successfully")
	}

	log.Println("✅ OBC data migration completed!")
	log.Println("")
	log.Println("📝 Next steps:")
	log.Println("1. Import real OBC data via Excel import endpoint")
	log.Println("2. Update existing POs to link to real OBC Masters")
	log.Println("3. Optional: Delete the legacy OBC Master (OBC-LEGACY-MIGRATION)")
}

/**
 * addForeignKeyConstraint menambahkan foreign key constraint antara
 * production_orders.obc_master_id dan obc_masters.id
 */
func addForeignKeyConstraint(db *gorm.DB) error {
	// Check if constraint already exists
	var count int64
	query := `
		SELECT COUNT(*) 
		FROM INFORMATION_SCHEMA.table_constraints 
		WHERE constraint_schema = DATABASE() 
		AND table_name = 'production_orders' 
		AND constraint_name = 'fk_obc_masters_production_orders'
	`
	
	if err := db.Raw(query).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check constraint: %w", err)
	}

	if count > 0 {
		log.Println("ℹ️ Foreign key constraint already exists")
		return nil
	}

	// Add constraint
	alterQuery := `
		ALTER TABLE production_orders 
		ADD CONSTRAINT fk_obc_masters_production_orders 
		FOREIGN KEY (obc_master_id) 
		REFERENCES obc_masters(id)
	`
	
	if err := db.Exec(alterQuery).Error; err != nil {
		return fmt.Errorf("failed to add constraint: %w", err)
	}

	return nil
}
