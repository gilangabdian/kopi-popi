package inventory

import (
	"time"

	"github.com/gilangages/kopi-popi/internal/catalog"
)

// BranchInventory merepresentasikan stok material di suatu cabang
type BranchInventory struct {
	ID         int              `json:"id" gorm:"primaryKey;autoIncrement"`
	BranchID   int              `json:"branch_id"`
	MaterialID int              `json:"material_id"`
	Material   catalog.Material `json:"material" gorm:"foreignKey:MaterialID"`
	Quantity   float64          `json:"quantity" gorm:"type:decimal(12,2)"`
}

// InventoryMovement merepresentasikan log perpindahan stok (Kartu Stok)
type InventoryMovement struct {
	ID           string           `json:"id" gorm:"primaryKey"`
	BranchID     int              `json:"branch_id"`
	MaterialID   int              `json:"material_id"`
	Material     catalog.Material `json:"material" gorm:"foreignKey:MaterialID"`
	MovementType string           `json:"movement_type"` // 'IN', 'OUT', 'ADJUSTMENT'
	Quantity     float64          `json:"quantity" gorm:"type:decimal(12,2)"`
	Description  string           `json:"description"`
	CreatedAt    time.Time        `json:"created_at"`
}

// RestockRequest merepresentasikan surat permintaan restock
type RestockRequest struct {
	ID          string        `json:"id" gorm:"primaryKey"` // UUID
	BranchID    int           `json:"branch_id"`
	RequestedBy string        `json:"requested_by"` // UUID User
	Status          string        `json:"status"`       // 'Pending', 'Approved', 'Rejected', 'Delivered'
	Reason          string        `json:"reason"`
	RejectionReason *string       `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Items       []RestockItem `json:"items" gorm:"foreignKey:RequestID"`
}

// RestockItem merepresentasikan detail barang dari sebuah permintaan restock
type RestockItem struct {
	ID                int              `json:"id" gorm:"primaryKey;autoIncrement"`
	RequestID         string           `json:"request_id"`
	MaterialID        int              `json:"material_id"`
	Material          catalog.Material `json:"material" gorm:"foreignKey:MaterialID"`
	QuantityRequested float64          `json:"quantity_requested" gorm:"type:decimal(12,2)"`
}
