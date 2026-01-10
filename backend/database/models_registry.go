package database

import (
	"sirine-go/backend/internal/counting"
	"sirine-go/backend/internal/cutting"
	"sirine-go/backend/models"
)

// ModelsRegistry merupakan central registry untuk semua database models
// yang akan di-migrate oleh GORM AutoMigrate
type ModelsRegistry struct {
	models []interface{}
	tables []string // Table names dalam urutan untuk drop (reverse order)
}

// NewModelsRegistry membuat instance baru dari ModelsRegistry
// dengan semua models yang diregister
func NewModelsRegistry() *ModelsRegistry {
	registry := &ModelsRegistry{
		models: make([]interface{}, 0),
		tables: make([]string, 0),
	}

	// Register semua models di sini
	// Format: registry.Register(&ModelStruct{}, "table_name")
	registry.Register(&models.User{}, "users")
	registry.Register(&models.UserSession{}, "user_sessions")
	registry.Register(&models.PasswordResetToken{}, "password_reset_tokens")
	registry.Register(&models.ActivityLog{}, "activity_logs")
	registry.Register(&models.Notification{}, "notifications")

	// OBC Master & Production Order models (OBCMaster HARUS sebelum ProductionOrder untuk foreign key)
	registry.Register(&models.OBCMaster{}, "obc_masters")
	registry.Register(&models.ProductionOrder{}, "production_orders")
	registry.Register(&models.POStageTracking{}, "po_stage_trackings")
	registry.Register(&models.KhazwalMaterialPreparation{}, "khazwal_material_preparations")

	// Khazwal Counting & Cutting models (Epic 2 & 3)
	registry.Register(&counting.KhazwalCountingResult{}, "khazwal_counting_results")
	registry.Register(&cutting.KhazwalCuttingResult{}, "khazwal_cutting_results")

	return registry
}

// Register menambahkan model dan table name ke registry
func (r *ModelsRegistry) Register(model interface{}, tableName string) {
	r.models = append(r.models, model)
	// Tables disimpan untuk rollback (akan di-reverse nanti)
	r.tables = append(r.tables, tableName)
}

// GetModels mengembalikan slice dari semua registered models
// untuk digunakan dalam AutoMigrate
func (r *ModelsRegistry) GetModels() []interface{} {
	return r.models
}

// GetTablesForRollback mengembalikan table names dalam reverse order
// untuk rollback yang aman (menghindari foreign key constraints)
func (r *ModelsRegistry) GetTablesForRollback() []string {
	reversed := make([]string, len(r.tables))
	for i, table := range r.tables {
		reversed[len(r.tables)-1-i] = table
	}
	return reversed
}

// GetTableCount mengembalikan jumlah tables yang registered
func (r *ModelsRegistry) GetTableCount() int {
	return len(r.tables)
}
