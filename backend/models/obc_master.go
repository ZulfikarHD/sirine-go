package models

import (
	"time"

	"gorm.io/gorm"
)

// OBCMaster merupakan model untuk master data OBC (Order Batch Confirmation)
// yang mencakup semua spesifikasi produk dari SAP dan akan di-reference oleh ProductionOrder
type OBCMaster struct {
	ID                   uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	OBCNumber            string         `gorm:"uniqueIndex;type:varchar(20);not null" json:"obc_number" binding:"required"`
	OBCDate              *time.Time     `gorm:"type:date" json:"obc_date"`
	Material             string         `gorm:"index;type:varchar(50)" json:"material"`
	Seri                 string         `gorm:"index;type:varchar(50)" json:"seri"`
	Warna                string         `gorm:"index;type:varchar(50)" json:"warna"`
	FactoryCode          string         `gorm:"index;type:varchar(50)" json:"factory_code"`
	QuantityOrdered      int            `json:"quantity_ordered"`
	JHT                  string         `gorm:"type:varchar(100)" json:"jht"`
	RPB                  float64        `gorm:"type:decimal(15,2)" json:"rpb"`
	HJE                  float64        `gorm:"type:decimal(15,2)" json:"hje"`
	BPB                  int            `json:"bpb"`
	Rencet               int            `json:"rencet"`
	DueDate              *time.Time     `gorm:"type:date" json:"due_date"`
	Personalization      string         `gorm:"type:varchar(20)" json:"personalization"`
	AdhesiveType         string         `gorm:"type:varchar(50)" json:"adhesive_type"`
	GR                   string         `gorm:"type:varchar(50)" json:"gr"`
	PlatNumber           string         `gorm:"type:varchar(50)" json:"plat_number"`
	Type                 string         `gorm:"type:varchar(50)" json:"type"`
	CreatedOn            *time.Time     `gorm:"type:date" json:"created_on"`
	SalesDocument        string         `gorm:"type:varchar(50)" json:"sales_document"`
	ItemCode             string         `gorm:"type:varchar(50)" json:"item_code"`
	MaterialDescription  string         `gorm:"type:varchar(255)" json:"material_description"`
	BaseUnit             string         `gorm:"type:varchar(20)" json:"base_unit"`
	PCACategory          string         `gorm:"type:varchar(50)" json:"pca_category"`
	AlcoholPercentage    float64        `gorm:"type:decimal(5,2)" json:"alcohol_percentage"`
	HPTLContent          float64        `gorm:"type:decimal(5,2)" json:"hptl_content"`
	RegionCode           string         `gorm:"type:varchar(20)" json:"region_code"`
	OBCInitial           string         `gorm:"type:varchar(50)" json:"obc_initial"`
	Allocation           string         `gorm:"type:varchar(255)" json:"allocation"`
	TotalOrderOBC        int            `json:"total_order_obc"`
	PlantCode            string         `gorm:"type:varchar(10)" json:"plant_code"`
	Unit                 string         `gorm:"type:varchar(20)" json:"unit"`
	ProductionYear       int            `json:"production_year"`
	ExciseRatePerLiter   float64        `gorm:"type:decimal(15,2)" json:"excise_rate_per_liter"`
	PCAVolume            float64        `gorm:"type:decimal(15,2)" json:"pca_volume"`
	MMEAColorCode        string         `gorm:"type:varchar(50)" json:"mmea_color_code"`
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	ProductionOrders []ProductionOrder `gorm:"foreignKey:OBCMasterID" json:"production_orders,omitempty"`
}

// TableName menentukan nama tabel di database
func (OBCMaster) TableName() string {
	return "obc_masters"
}

// GetDisplayName mengembalikan display name untuk OBC
// yang berisi kombinasi OBC Number dan Material
func (o *OBCMaster) GetDisplayName() string {
	if o.Material != "" {
		return o.OBCNumber + " - " + o.Material
	}
	return o.OBCNumber
}

// HasProductionOrders memeriksa apakah OBC memiliki Production Orders
func (o *OBCMaster) HasProductionOrders() bool {
	return len(o.ProductionOrders) > 0
}

// GetTotalPOQuantity menghitung total quantity dari semua Production Orders
// yang terkait dengan OBC ini
func (o *OBCMaster) GetTotalPOQuantity() int {
	total := 0
	for _, po := range o.ProductionOrders {
		total += po.QuantityOrdered
	}
	return total
}

// CalculateTotalWithBuffer menghitung total quantity dengan buffer 6%
// sesuai dengan rumus: Total = QTY + (QTY * 6%)
func (o *OBCMaster) CalculateTotalWithBuffer() int {
	buffer := float64(o.QuantityOrdered) * 0.06
	return o.QuantityOrdered + int(buffer)
}

// IsPersonalized memeriksa apakah OBC merupakan produk personalisasi
func (o *OBCMaster) IsPersonalized() bool {
	return o.Personalization == "Perso" || o.Personalization == "PERSO"
}
