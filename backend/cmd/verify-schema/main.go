package main

import (
	"fmt"
	"log"
	"sirine-go/backend/config"
	"sirine-go/backend/database"

	"github.com/joho/godotenv"
)

/**
 * Script untuk verify database schema setelah migration
 * yang mencakup:
 * 1. Check obc_masters table structure
 * 2. Check production_orders table structure
 * 3. Verify foreign key constraints
 * 4. Check indexes
 */

type ColumnInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default *string
	Extra   string
}

type IndexInfo struct {
	Table      string
	NonUnique  int
	KeyName    string
	SeqInIndex int
	ColumnName string
}

type ForeignKeyInfo struct {
	ConstraintName            string
	TableName                 string
	ColumnName                string
	ReferencedTableName       string
	ReferencedColumnName      string
	UpdateRule                string
	DeleteRule                string
}

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

	fmt.Println("================================================================================")
	fmt.Println("DATABASE SCHEMA VERIFICATION")
	fmt.Println("================================================================================")
	fmt.Println()

	// 1. Verify obc_masters table
	fmt.Println("📋 OBC MASTERS TABLE")
	fmt.Println("--------------------------------------------------------------------------------")
	
	var obcColumns []ColumnInfo
	if err := db.Raw("DESCRIBE obc_masters").Scan(&obcColumns).Error; err != nil {
		log.Fatal("Failed to describe obc_masters:", err)
	}

	fmt.Printf("Columns: %d\n\n", len(obcColumns))
	fmt.Printf("%-30s %-20s %-6s %-6s %-10s\n", "FIELD", "TYPE", "NULL", "KEY", "EXTRA")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, col := range obcColumns {
		defaultVal := "NULL"
		if col.Default != nil {
			defaultVal = *col.Default
		}
		fmt.Printf("%-30s %-20s %-6s %-6s %-10s\n", 
			col.Field, col.Type, col.Null, col.Key, col.Extra)
		_ = defaultVal // avoid unused variable
	}
	fmt.Println()

	// 2. Verify production_orders table
	fmt.Println("📋 PRODUCTION ORDERS TABLE")
	fmt.Println("--------------------------------------------------------------------------------")
	
	var poColumns []ColumnInfo
	if err := db.Raw("DESCRIBE production_orders").Scan(&poColumns).Error; err != nil {
		log.Fatal("Failed to describe production_orders:", err)
	}

	fmt.Printf("Columns: %d\n\n", len(poColumns))
	fmt.Printf("%-30s %-20s %-6s %-6s %-10s\n", "FIELD", "TYPE", "NULL", "KEY", "EXTRA")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, col := range poColumns {
		fmt.Printf("%-30s %-20s %-6s %-6s %-10s\n", 
			col.Field, col.Type, col.Null, col.Key, col.Extra)
	}
	fmt.Println()

	// 3. Verify indexes on obc_masters
	fmt.Println("🔍 OBC MASTERS INDEXES")
	fmt.Println("--------------------------------------------------------------------------------")
	
	var obcIndexes []IndexInfo
	query := `
		SELECT TABLE_NAME as 'Table', NON_UNIQUE as NonUnique, INDEX_NAME as KeyName, 
		       SEQ_IN_INDEX as SeqInIndex, COLUMN_NAME as ColumnName
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'obc_masters'
		ORDER BY INDEX_NAME, SEQ_IN_INDEX
	`
	if err := db.Raw(query).Scan(&obcIndexes).Error; err != nil {
		log.Fatal("Failed to get obc_masters indexes:", err)
	}

	currentIndex := ""
	for _, idx := range obcIndexes {
		if idx.KeyName != currentIndex {
			unique := "NON-UNIQUE"
			if idx.NonUnique == 0 {
				unique = "UNIQUE"
			}
			fmt.Printf("\n%-30s [%s]\n", idx.KeyName, unique)
			currentIndex = idx.KeyName
		}
		fmt.Printf("  └─ %s\n", idx.ColumnName)
	}
	fmt.Println()

	// 4. Verify indexes on production_orders
	fmt.Println("🔍 PRODUCTION ORDERS INDEXES")
	fmt.Println("--------------------------------------------------------------------------------")
	
	var poIndexes []IndexInfo
	query = `
		SELECT TABLE_NAME as 'Table', NON_UNIQUE as NonUnique, INDEX_NAME as KeyName, 
		       SEQ_IN_INDEX as SeqInIndex, COLUMN_NAME as ColumnName
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'production_orders'
		ORDER BY INDEX_NAME, SEQ_IN_INDEX
	`
	if err := db.Raw(query).Scan(&poIndexes).Error; err != nil {
		log.Fatal("Failed to get production_orders indexes:", err)
	}

	currentIndex = ""
	for _, idx := range poIndexes {
		if idx.KeyName != currentIndex {
			unique := "NON-UNIQUE"
			if idx.NonUnique == 0 {
				unique = "UNIQUE"
			}
			fmt.Printf("\n%-30s [%s]\n", idx.KeyName, unique)
			currentIndex = idx.KeyName
		}
		fmt.Printf("  └─ %s\n", idx.ColumnName)
	}
	fmt.Println()

	// 5. Verify foreign keys
	fmt.Println("🔗 FOREIGN KEY CONSTRAINTS")
	fmt.Println("--------------------------------------------------------------------------------")
	
	var foreignKeys []ForeignKeyInfo
	query = `
		SELECT 
			kcu.CONSTRAINT_NAME as ConstraintName,
			kcu.TABLE_NAME as TableName,
			kcu.COLUMN_NAME as ColumnName,
			kcu.REFERENCED_TABLE_NAME as ReferencedTableName,
			kcu.REFERENCED_COLUMN_NAME as ReferencedColumnName,
			rc.UPDATE_RULE as UpdateRule,
			rc.DELETE_RULE as DeleteRule
		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
		JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc 
			ON kcu.CONSTRAINT_NAME = rc.CONSTRAINT_NAME
			AND kcu.CONSTRAINT_SCHEMA = rc.CONSTRAINT_SCHEMA
		WHERE kcu.TABLE_SCHEMA = DATABASE()
		AND kcu.TABLE_NAME = 'production_orders'
		AND kcu.REFERENCED_TABLE_NAME IS NOT NULL
		ORDER BY kcu.CONSTRAINT_NAME
	`
	if err := db.Raw(query).Scan(&foreignKeys).Error; err != nil {
		log.Fatal("Failed to get foreign keys:", err)
	}

	if len(foreignKeys) == 0 {
		fmt.Println("⚠️  No foreign keys found on production_orders table")
	} else {
		for _, fk := range foreignKeys {
			fmt.Printf("\nConstraint: %s\n", fk.ConstraintName)
			fmt.Printf("  Table:      %s (%s)\n", fk.TableName, fk.ColumnName)
			fmt.Printf("  References: %s (%s)\n", fk.ReferencedTableName, fk.ReferencedColumnName)
			fmt.Printf("  On Update:  %s\n", fk.UpdateRule)
			fmt.Printf("  On Delete:  %s\n", fk.DeleteRule)
		}
	}
	fmt.Println()

	// 6. Verify data
	fmt.Println("📊 DATA SUMMARY")
	fmt.Println("-" * 80)
	
	var obcCount int64
	var poCount int64
	
	if err := db.Raw("SELECT COUNT(*) FROM obc_masters WHERE deleted_at IS NULL").Scan(&obcCount).Error; err != nil {
		log.Fatal("Failed to count obc_masters:", err)
	}
	
	if err := db.Raw("SELECT COUNT(*) FROM production_orders WHERE deleted_at IS NULL").Scan(&poCount).Error; err != nil {
		log.Fatal("Failed to count production_orders:", err)
	}

	fmt.Printf("OBC Masters:        %d records\n", obcCount)
	fmt.Printf("Production Orders:  %d records\n", poCount)
	fmt.Println()

	// Check orphaned POs
	var orphanedCount int64
	query = `
		SELECT COUNT(*) 
		FROM production_orders po
		LEFT JOIN obc_masters obc ON po.obc_master_id = obc.id
		WHERE po.deleted_at IS NULL AND obc.id IS NULL
	`
	if err := db.Raw(query).Scan(&orphanedCount).Error; err != nil {
		log.Printf("Warning: Failed to check orphaned POs: %v", err)
	} else {
		if orphanedCount > 0 {
			fmt.Printf("⚠️  Found %d orphaned Production Orders (no linked OBC Master)\n", orphanedCount)
		} else {
			fmt.Printf("✅ All Production Orders have valid OBC Master links\n")
		}
	}

	fmt.Println()
	fmt.Println("=" * 80)
	fmt.Println("✅ SCHEMA VERIFICATION COMPLETED")
	fmt.Println("=" * 80)
}
