package analytics

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetSalesReport(ctx context.Context, branchID *int, startDate, endDate string) ([]SalesReportResponse, error)
	GetTopProducts(ctx context.Context, branchID *int, startDate, endDate string, limit int) ([]TopProductResponse, error)
	GetShiftReports(ctx context.Context, branchID *int, startDate, endDate string) ([]ShiftReportResponse, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetSalesReport(ctx context.Context, branchID *int, startDate, endDate string) ([]SalesReportResponse, error) {
	validStatuses := []string{"Paid", "Preparing", "Ready", "Completed"}

	// 1. Dapatkan Total Revenue & Count per hari
	type DailySales struct {
		Date              string
		TotalRevenue      float64
		TotalTransactions int
	}
	var dailySales []DailySales

	query := r.db.Table("transactions").
		Select("DATE(created_at) as date, SUM(total_amount) as total_revenue, COUNT(id) as total_transactions").
		Where("status IN ?", validStatuses)

	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	err := query.Group("DATE(created_at)").Order("date ASC").Scan(&dailySales).Error
	if err != nil {
		return nil, err
	}

	// 2. Dapatkan Payment Breakdown per hari
	type PaymentAgg struct {
		Date   string
		Method string
		Amount float64
	}
	var paymentAggs []PaymentAgg

	queryPay := r.db.Table("transactions").
		Select("DATE(created_at) as date, payment_method as method, SUM(total_amount) as amount").
		Where("status IN ?", validStatuses)

	if branchID != nil {
		queryPay = queryPay.Where("branch_id = ?", *branchID)
	}
	if startDate != "" && endDate != "" {
		queryPay = queryPay.Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	err = queryPay.Group("DATE(created_at), payment_method").Scan(&paymentAggs).Error
	if err != nil {
		return nil, err
	}

	// Gabungkan hasil
	var results []SalesReportResponse
	for _, ds := range dailySales {
		var breakdown []PaymentMethodBreakdown
		for _, pa := range paymentAggs {
			if pa.Date == ds.Date {
				breakdown = append(breakdown, PaymentMethodBreakdown{
					Method: pa.Method,
					Amount: pa.Amount,
				})
			}
		}

		results = append(results, SalesReportResponse{
			Date:              ds.Date,
			TotalRevenue:      ds.TotalRevenue,
			TotalTransactions: ds.TotalTransactions,
			PaymentBreakdown:  breakdown,
		})
	}

	return results, nil
}

func (r *repository) GetTopProducts(ctx context.Context, branchID *int, startDate, endDate string, limit int) ([]TopProductResponse, error) {
	validStatuses := []string{"Paid", "Preparing", "Ready", "Completed"}
	var topProducts []TopProductResponse

	query := r.db.Table("transaction_details").
		Select("transaction_details.product_id, products.name as product_name, SUM(transaction_details.quantity) as total_quantity_sold, SUM(transaction_details.subtotal) as total_revenue").
		Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
		Joins("JOIN products ON products.id = transaction_details.product_id").
		Where("transactions.status IN ?", validStatuses)

	if branchID != nil {
		query = query.Where("transactions.branch_id = ?", *branchID)
	}
	if startDate != "" && endDate != "" {
		query = query.Where("transactions.created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	err := query.Group("transaction_details.product_id, products.name").
		Order("total_quantity_sold DESC").
		Limit(limit).
		Scan(&topProducts).Error

	return topProducts, err
}

func (r *repository) GetShiftReports(ctx context.Context, branchID *int, startDate, endDate string) ([]ShiftReportResponse, error) {
	var shiftReports []ShiftReportResponse

	query := r.db.Table("shifts").
		Select("shifts.id as shift_id, users.name as cashier_name, shifts.branch_id, shifts.start_time, shifts.end_time, shifts.expected_cash, shifts.actual_cash, (shifts.actual_cash - shifts.expected_cash) as difference").
		Joins("JOIN users ON users.id = shifts.cashier_id").
		Where("shifts.status = ?", "Closed")

	if branchID != nil {
		query = query.Where("shifts.branch_id = ?", *branchID)
	}
	if startDate != "" && endDate != "" {
		// Menggunakan waktu tutup shift (end_time) sebagai patokan filter tanggal
		query = query.Where("shifts.end_time BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	err := query.Order("shifts.end_time DESC").Scan(&shiftReports).Error
	return shiftReports, err
}
