package sales

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	// Shifts
	CreateShift(shift *Shift) error
	UpdateShift(shift *Shift) error
	GetOpenShiftByCashier(cashierID string) (*Shift, error)
	GetShiftByID(id string) (*Shift, error)

	// Expenses
	CreateExpense(tx *gorm.DB, expense *Expense) error
	GetExpensesByShiftID(shiftID string) ([]Expense, error)

	// Carts
	CreateCart(cart *Cart) error
	GetCartByID(id string) (*Cart, error)
	GetActiveCartByCustomer(customerID string) (*Cart, error)
	GetActiveCartsByBranch(branchID int) ([]Cart, error) // Untuk Kasir lihat semua Hold Bill di cabangnya
	DeleteCart(id string) error

	// Cart Items
	AddOrUpdateCartItem(item *CartItem) error
	RemoveCartItem(cartID string, productID int) error
	ClearCartItems(cartID string) error

	// Transactions
	CreateTransaction(tx *gorm.DB, transaction *Transaction) error
	GetTransactionByID(id string) (*Transaction, error)
	GetTransactions(branchID *int, customerID *string, status *string, startDate *string, endDate *string) ([]Transaction, error)
	UpdateTransactionStatus(id string, status string) error

	// DB Transaction Helper
	WithTransaction(fn func(tx *gorm.DB) error) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// --- SHIFTS ---

func (r *repository) CreateShift(shift *Shift) error {
	return r.db.Create(shift).Error
}

func (r *repository) UpdateShift(shift *Shift) error {
	return r.db.Save(shift).Error
}

func (r *repository) GetOpenShiftByCashier(cashierID string) (*Shift, error) {
	var shift Shift
	err := r.db.Where("cashier_id = ? AND status = ?", cashierID, "Open").First(&shift).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &shift, nil
}

func (r *repository) GetShiftByID(id string) (*Shift, error) {
	var shift Shift
	err := r.db.First(&shift, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &shift, nil
}

// --- EXPENSES ---

func (r *repository) CreateExpense(tx *gorm.DB, expense *Expense) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(expense).Error
}

func (r *repository) GetExpensesByShiftID(shiftID string) ([]Expense, error) {
	var expenses []Expense
	err := r.db.Where("shift_id = ?", shiftID).Order("created_at desc").Find(&expenses).Error
	return expenses, err
}

// --- CARTS ---

func (r *repository) CreateCart(cart *Cart) error {
	return r.db.Create(cart).Error
}

func (r *repository) GetCartByID(id string) (*Cart, error) {
	var cart Cart
	err := r.db.Preload("Items").Preload("Items.Product").First(&cart, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *repository) GetActiveCartByCustomer(customerID string) (*Cart, error) {
	var cart Cart
	err := r.db.Preload("Items").Preload("Items.Product").
		Where("customer_id = ?", customerID).
		First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *repository) GetActiveCartsByBranch(branchID int) ([]Cart, error) {
	var carts []Cart
	// Ambil semua cart yang tidak punya customer_id (Offline Cart / Hold Bill) di cabang tersebut
	err := r.db.Preload("Items").Preload("Items.Product").
		Where("branch_id = ? AND customer_id IS NULL", branchID).
		Find(&carts).Error
	return carts, err
}

func (r *repository) DeleteCart(id string) error {
	return r.db.Delete(&Cart{}, "id = ?", id).Error
}

// --- CART ITEMS ---

func (r *repository) AddOrUpdateCartItem(item *CartItem) error {
	var existing CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", item.CartID, item.ProductID).First(&existing).Error
	
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(item).Error
	}

	// Update if exists
	existing.Quantity += item.Quantity
	if existing.Quantity <= 0 {
		return r.db.Delete(&existing).Error
	}
	
	if item.Notes != nil {
		existing.Notes = item.Notes
	}

	return r.db.Save(&existing).Error
}

func (r *repository) RemoveCartItem(cartID string, productID int) error {
	return r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).Delete(&CartItem{}).Error
}

func (r *repository) ClearCartItems(cartID string) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&CartItem{}).Error
}

// --- TRANSACTIONS ---

func (r *repository) CreateTransaction(tx *gorm.DB, transaction *Transaction) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(transaction).Error
}

func (r *repository) GetTransactionByID(id string) (*Transaction, error) {
	var transaction Transaction
	err := r.db.Preload("Details").Preload("Details.Product").First(&transaction, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *repository) GetTransactions(branchID *int, customerID *string, status *string, startDate *string, endDate *string) ([]Transaction, error) {
	var transactions []Transaction
	query := r.db.Preload("Details").Preload("Details.Product")

	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}
	if customerID != nil {
		query = query.Where("customer_id = ?", *customerID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *startDate+" 00:00:00", *endDate+" 23:59:59")
	}

	err := query.Order("created_at desc").Find(&transactions).Error
	return transactions, err
}

func (r *repository) UpdateTransactionStatus(id string, status string) error {
	return r.db.Model(&Transaction{}).Where("id = ?", id).Update("status", status).Error
}

// --- DB TX HELPER ---

func (r *repository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
