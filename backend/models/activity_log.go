package models

import (
	"encoding/json"
	"time"
)

// ActivityAction merupakan enum untuk jenis aksi
type ActivityAction string

const (
	ActionCreate         ActivityAction = "CREATE"
	ActionUpdate         ActivityAction = "UPDATE"
	ActionDelete         ActivityAction = "DELETE"
	ActionLogin          ActivityAction = "LOGIN"
	ActionLogout         ActivityAction = "LOGOUT"
	ActionPasswordChange ActivityAction = "PASSWORD_CHANGE"
)

// ActivityLog merupakan model untuk audit trail
// yang mencakup tracking semua critical actions dalam sistem
type ActivityLog struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64         `gorm:"not null;index" json:"user_id"`
	Action     ActivityAction `gorm:"type:enum('CREATE','UPDATE','DELETE','LOGIN','LOGOUT','PASSWORD_CHANGE');not null;index" json:"action"`
	EntityType string         `gorm:"type:varchar(50);not null;index" json:"entity_type"` // Table name atau entity type
	EntityID   *uint64        `gorm:"type:bigint unsigned" json:"entity_id"`
	Changes    json.RawMessage `gorm:"type:json" json:"changes"` // Before/after values dalam JSON format
	IPAddress  string         `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent  string         `gorm:"type:text" json:"user_agent"`
	CreatedAt  time.Time      `gorm:"autoCreateTime;index" json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName menentukan nama tabel di database
func (ActivityLog) TableName() string {
	return "activity_logs"
}

// ChangeData merupakan struktur untuk before/after values
type ChangeData struct {
	Before interface{} `json:"before"`
	After  interface{} `json:"after"`
}

// SetChanges mengset changes dalam format JSON
func (a *ActivityLog) SetChanges(before, after interface{}) error {
	changeData := ChangeData{
		Before: before,
		After:  after,
	}
	data, err := json.Marshal(changeData)
	if err != nil {
		return err
	}
	a.Changes = data
	return nil
}

// GetChanges mengambil changes dari JSON
func (a *ActivityLog) GetChanges() (*ChangeData, error) {
	var changeData ChangeData
	if len(a.Changes) == 0 {
		return nil, nil
	}
	err := json.Unmarshal(a.Changes, &changeData)
	if err != nil {
		return nil, err
	}
	return &changeData, nil
}
