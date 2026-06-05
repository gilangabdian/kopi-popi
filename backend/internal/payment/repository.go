package payment

import (
	"context"
	"gorm.io/gorm"
)

type TransactionInfo struct {
	ID            string
	TotalAmount   float64
	CustomerName  string
	CustomerEmail string
}

type Repository interface {
	UpdateTransactionStatusToPaid(ctx context.Context, orderID string) error
	UpdateTransactionStatusToFailed(ctx context.Context, orderID string) error
	GetTransactionInfo(ctx context.Context, orderID string) (*TransactionInfo, error)
	DeductStockForTransaction(ctx context.Context, orderID string) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &paymentRepository{db}
}

func (r *paymentRepository) UpdateTransactionStatusToPaid(ctx context.Context, orderID string) error {
	return r.db.WithContext(ctx).Table("transactions").
		Where("id = ?", orderID).
		Update("status", "Paid").Error
}

func (r *paymentRepository) UpdateTransactionStatusToFailed(ctx context.Context, orderID string) error {
	return r.db.WithContext(ctx).Table("transactions").
		Where("id = ?", orderID).
		Update("status", "Failed").Error
}

func (r *paymentRepository) GetTransactionInfo(ctx context.Context, orderID string) (*TransactionInfo, error) {
	var info TransactionInfo
	// Join with users to get email
	err := r.db.WithContext(ctx).Table("transactions").
		Select("transactions.id, transactions.total_amount, COALESCE(transactions.customer_name, users.name) as customer_name, users.email as customer_email").
		Joins("LEFT JOIN users ON users.id = transactions.customer_id").
		Where("transactions.id = ?", orderID).
		Scan(&info).Error

	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *paymentRepository) DeductStockForTransaction(ctx context.Context, orderID string) error {
	// Raw query to deduct stock based on BOM for each cart_item
	// This is a complex query but prevents cyclic dependencies between payment and inventory domains
	query := `
		UPDATE branch_inventories bi
		JOIN (
			SELECT
				t.branch_id,
				mi.material_id,
				SUM(td.quantity * mi.quantity_needed) as total_qty_needed
			FROM transactions t
			JOIN transaction_details td ON td.transaction_id = t.id
			JOIN product_boms mi ON mi.product_id = td.product_id
			WHERE t.id = ?
			GROUP BY t.branch_id, mi.material_id
		) needed ON bi.branch_id = needed.branch_id AND bi.material_id = needed.material_id
		SET bi.quantity = bi.quantity - needed.total_qty_needed
	`
	return r.db.WithContext(ctx).Exec(query, orderID).Error
}
