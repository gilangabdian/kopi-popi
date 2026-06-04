package branch

import "time"

// Branch merepresentasikan struktur data tabel branches di MySQL
type Branch struct {
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name              string    `json:"name"`
	Address           string    `json:"address"`
	IsActive          bool      `json:"is_active"`
	OpeningTime       *string   `json:"opening_time"`
	ClosingTime       *string   `json:"closing_time"`
	IsAcceptingOrders bool      `json:"is_accepting_orders"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CreateBranchRequest merepresentasikan payload untuk membuat cabang baru
type CreateBranchRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

// UpdateBranchRequest merepresentasikan payload untuk memperbarui cabang
type UpdateBranchRequest struct {
	Name     *string `json:"name,omitempty"`
	Address  *string `json:"address,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// UpdateOperatingHoursRequest merepresentasikan payload untuk mengubah jam operasional
type UpdateOperatingHoursRequest struct {
	OpeningTime *string `json:"opening_time" binding:"required"`
	ClosingTime *string `json:"closing_time" binding:"required"`
}

// UpdateAcceptingOrdersRequest merepresentasikan payload untuk tombol darurat buka/tutup toko
type UpdateAcceptingOrdersRequest struct {
	IsAcceptingOrders *bool `json:"is_accepting_orders" binding:"required"`
}
