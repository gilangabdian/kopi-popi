package analytics

import "time"

// PaymentMethodBreakdown merepresentasikan rincian pendapatan per metode pembayaran.
type PaymentMethodBreakdown struct {
	Method string  `json:"method"`
	Amount float64 `json:"amount"`
}

// SalesReportResponse merepresentasikan laporan pendapatan.
type SalesReportResponse struct {
	Date              string                   `json:"date"`               // YYYY-MM-DD
	TotalRevenue      float64                  `json:"total_revenue"`
	TotalTransactions int                      `json:"total_transactions"`
	PaymentBreakdown  []PaymentMethodBreakdown `json:"payment_breakdown"`
}

// TopProductResponse merepresentasikan laporan produk terlaris.
type TopProductResponse struct {
	ProductID         int     `json:"product_id"`
	ProductName       string  `json:"product_name"`
	TotalQuantitySold int     `json:"total_quantity_sold"`
	TotalRevenue      float64 `json:"total_revenue"` // Opsional, subtotal dari produk ini
}

// ShiftReportResponse merepresentasikan laporan tutup kasir.
type ShiftReportResponse struct {
	ShiftID      string    `json:"shift_id"`
	CashierName  string    `json:"cashier_name"`
	BranchID     int       `json:"branch_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	ExpectedCash float64   `json:"expected_cash"`
	ActualCash   float64   `json:"actual_cash"`
	Difference   float64   `json:"difference"`
}
