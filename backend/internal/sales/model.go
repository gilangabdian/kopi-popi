package sales

import (
	"time"

	"github.com/gilangages/kopi-popi/internal/catalog"
)

// Shift merepresentasikan tabel shifts
type Shift struct {
	ID           string     `json:"id" gorm:"primaryKey"`
	BranchID     int        `json:"branch_id"`
	CashierID    string     `json:"cashier_id"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	StartingCash float64    `json:"starting_cash"`
	ExpectedCash float64    `json:"expected_cash"`
	ActualCash   *float64   `json:"actual_cash"`
	Status       string     `json:"status"` // Open, Closed, Verified
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Cart merepresentasikan tabel carts (bisa untuk Pelanggan Online / Kasir Offline Hold Bill)
type Cart struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	CustomerID *string    `json:"customer_id,omitempty"` // Null jika cart milik kasir (Offline)
	CartName   *string    `json:"cart_name,omitempty"`   // Contoh: "Meja 4", "Mas Baju Merah" (Offline)
	BranchID   int        `json:"branch_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Items      []CartItem `json:"items" gorm:"foreignKey:CartID"`
}

// CartItem merepresentasikan tabel cart_items
type CartItem struct {
	ID        int              `json:"id" gorm:"primaryKey;autoIncrement"`
	CartID    string           `json:"cart_id"`
	ProductID int              `json:"product_id"`
	Quantity  int              `json:"quantity"`
	Notes     *string          `json:"notes,omitempty"`
	Product   *catalog.Product `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"` // Untuk response
}

// Transaction merepresentasikan tabel transactions
type Transaction struct {
	ID            string              `json:"id" gorm:"primaryKey"`
	BranchID      int                 `json:"branch_id"`
	CustomerID    *string             `json:"customer_id,omitempty"`
	CustomerName  *string             `json:"customer_name,omitempty"`
	CashierID     *string             `json:"cashier_id,omitempty"`
	ShiftID       *string             `json:"shift_id,omitempty"`
	OrderType     string              `json:"order_type"` // Online_Pickup, Offline_DineIn, Offline_Takeaway
	PaymentMethod string              `json:"payment_method"` // "CASH", "QRIS", "DEBIT", dll
	TotalAmount   float64             `json:"total_amount"`
	Status        string              `json:"status"` // Waiting_Payment, Paid, Preparing, Ready, Completed, Cancelled
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	Details       []TransactionDetail `json:"details" gorm:"foreignKey:TransactionID"`
}

// TransactionDetail merepresentasikan tabel transaction_details
type TransactionDetail struct {
	ID            int              `json:"id" gorm:"primaryKey;autoIncrement"`
	TransactionID string           `json:"transaction_id"`
	ProductID     int              `json:"product_id"`
	Quantity      int              `json:"quantity"`
	Subtotal      float64          `json:"subtotal"`
	Notes         *string          `json:"notes,omitempty"`
	Product       *catalog.Product `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"`
}

// --- REQUEST PAYLOADS ---

type OpenShiftRequest struct {
	StartingCash float64 `json:"starting_cash" binding:"required"`
}

type CloseShiftRequest struct {
	ActualCash float64 `json:"actual_cash" binding:"required"`
}

type AddCartItemRequest struct {
	ProductID int     `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Notes     *string `json:"notes,omitempty"`
}

type InitOfflineCartRequest struct {
	CartName string `json:"cart_name" binding:"required"`
}

type CheckoutRequest struct {
	CartID        string  `json:"cart_id" binding:"required"`
	OrderType     string  `json:"order_type" binding:"required"` // Online_Pickup, Offline_DineIn, Offline_Takeaway
	PaymentMethod string  `json:"payment_method" binding:"required"` // CASH, QRIS, dll
	CustomerName  *string `json:"customer_name,omitempty"` // Wajib diisi kasir (offline)
	AmountTendered *float64 `json:"amount_tendered,omitempty"` // Khusus Cash, uang yang diberikan pelanggan (untuk kembalian - opsional)
}

type UpdateTransactionStatusPayload struct {
	Status string `json:"status" binding:"required,oneof=Paid Preparing Ready Completed Cancelled"`
}
