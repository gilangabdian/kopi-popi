package sales

import (
	"context"
	"errors"
	"time"

	"github.com/gilangages/kopi-popi/internal/branch"
	"github.com/gilangages/kopi-popi/internal/catalog"
	"github.com/gilangages/kopi-popi/internal/inventory"
	"github.com/gilangages/kopi-popi/internal/notification"
	"github.com/gilangages/kopi-popi/internal/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	// Shifts
	OpenShift(ctx context.Context, branchID int, cashierID string, req OpenShiftRequest) (*Shift, error)
	CloseShift(ctx context.Context, cashierID string, req CloseShiftRequest) error
	GetMyOpenShift(ctx context.Context, cashierID string) (*Shift, error)

	// Expenses
	RecordExpense(ctx context.Context, cashierID string, req RecordExpenseRequest) error
	GetMyExpenses(ctx context.Context, cashierID string) ([]Expense, error)

	// Carts (Online & Offline)
	AddCartItem(ctx context.Context, customerID *string, branchID int, req AddCartItemRequest) error
	InitOfflineCart(ctx context.Context, branchID int, req InitOfflineCartRequest) (*Cart, error)
	AddItemToOfflineCart(ctx context.Context, cartID string, branchID int, req AddCartItemRequest) error
	GetMyCart(ctx context.Context, customerID string) (*Cart, error)
	GetOfflineCarts(ctx context.Context, branchID int) ([]Cart, error)
	ClearCart(ctx context.Context, cartID string) error

	// Checkout
	Checkout(ctx context.Context, customerID *string, cashierID *string, req CheckoutRequest) (*Transaction, error)

	// Transactions
	GetTransactions(ctx context.Context, role string, reqBranchID *int, reqCustomerID *string, status *string, startDate *string, endDate *string) ([]Transaction, error)
	GetTransactionByID(ctx context.Context, id string, role string, reqBranchID *int, reqCustomerID *string) (*Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id string, status string, role string, reqBranchID *int) error
}

type service struct {
	repo       Repository
	branchSvc  branch.Service
	catalogSvc catalog.Service
	invSvc     inventory.Service
	notifSvc   notification.Service
	userSvc    user.Service
}

func NewService(repo Repository, branchSvc branch.Service, catalogSvc catalog.Service, invSvc inventory.Service, notifSvc notification.Service, userSvc user.Service) Service {
	return &service{
		repo:       repo,
		branchSvc:  branchSvc,
		catalogSvc: catalogSvc,
		invSvc:     invSvc,
		notifSvc:   notifSvc,
		userSvc:    userSvc,
	}
}

// --- SHIFTS ---

func (s *service) OpenShift(ctx context.Context, branchID int, cashierID string, req OpenShiftRequest) (*Shift, error) {
	// Check if already open
	existing, err := s.repo.GetOpenShiftByCashier(cashierID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("conflict: you already have an open shift")
	}

	shift := &Shift{
		ID:           uuid.NewString(),
		BranchID:     branchID,
		CashierID:    cashierID,
		StartTime:    time.Now(),
		StartingCash: req.StartingCash,
		ExpectedCash: req.StartingCash,
		Status:       "Open",
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateShift(shift); err != nil {
		return nil, err
	}
	return shift, nil
}

func (s *service) CloseShift(ctx context.Context, cashierID string, req CloseShiftRequest) error {
	shift, err := s.repo.GetOpenShiftByCashier(cashierID)
	if err != nil {
		return err
	}
	if shift == nil {
		return errors.New("not found: no open shift found")
	}

	now := time.Now()
	shift.EndTime = &now
	shift.ActualCash = &req.ActualCash
	shift.Status = "Closed" // Wait for manager to verify
	shift.UpdatedAt = time.Now()

	err = s.repo.UpdateShift(shift)
	if err != nil {
		return err
	}

	// Notifikasi In-App ke Manager
	if s.notifSvc != nil && s.userSvc != nil {
		allUsers, err := s.userSvc.GetEmployees(ctx, "Admin", nil)
		if err == nil {
			for _, u := range allUsers {
				if u.BranchID != nil && *u.BranchID == shift.BranchID && u.RoleID == 2 { // Manager
					s.notifSvc.SendInAppNotification(u.ID, "Shift Ditutup", "Kasir baru saja menutup shift, harap verifikasi kas.", "INFO")
				}
			}
		}
	}

	return nil
}

func (s *service) GetMyOpenShift(ctx context.Context, cashierID string) (*Shift, error) {
	return s.repo.GetOpenShiftByCashier(cashierID)
}

// --- EXPENSES ---

func (s *service) RecordExpense(ctx context.Context, cashierID string, req RecordExpenseRequest) error {
	shift, err := s.repo.GetOpenShiftByCashier(cashierID)
	if err != nil {
		return errors.New("no open shift found")
	}

	return s.repo.WithTransaction(func(tx *gorm.DB) error {
		expense := &Expense{
			ShiftID:     shift.ID,
			Amount:      req.Amount,
			Description: req.Description,
		}
		if err := s.repo.CreateExpense(tx, expense); err != nil {
			return err
		}

		shift.ExpectedCash -= req.Amount
		if err := tx.Save(shift).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *service) GetMyExpenses(ctx context.Context, cashierID string) ([]Expense, error) {
	shift, err := s.repo.GetOpenShiftByCashier(cashierID)
	if err != nil {
		return nil, errors.New("no open shift found")
	}

	return s.repo.GetExpensesByShiftID(shift.ID)
}

// --- CARTS ---

func (s *service) AddCartItem(ctx context.Context, customerID *string, branchID int, req AddCartItemRequest) error {
	// 1. Cek apakah cabang menerima order
	branches, err := s.branchSvc.GetAllBranches(ctx, "Customer", true)
	if err != nil {
		return err
	}
	var targetBranch *branch.Branch
	for _, b := range branches {
		if b.ID == branchID {
			targetBranch = &b
			break
		}
	}
	if targetBranch == nil {
		return errors.New("branch not found")
	}
	if !targetBranch.IsAcceptingOrders || !targetBranch.IsActive {
		return errors.New("forbidden: branch is currently not accepting orders")
	}
	
	// Validasi Jam Operasional jika perlu (disini asumsinya ditangani)
	// if targetBranch.OpeningTime != nil && targetBranch.ClosingTime != nil {
	//   // Cek jam saat ini
	// }

	var cart *Cart

	if customerID != nil {
		// Scenario 1: Online Customer Cart
		cart, err = s.repo.GetActiveCartByCustomer(*customerID)
		if err != nil {
			return err
		}
		if cart == nil {
			// Bikin baru
			cart = &Cart{
				ID:         uuid.NewString(),
				CustomerID: customerID,
				BranchID:   branchID,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			if err := s.repo.CreateCart(cart); err != nil {
				return err
			}
		} else {
			// Pastikan branch nya sama, kalau beda, clear dulu
			if cart.BranchID != branchID {
				if err := s.repo.ClearCartItems(cart.ID); err != nil {
					return err
				}
				cart.BranchID = branchID
				cart.UpdatedAt = time.Now()
				// update cart branch? The repository doesn't have an update cart, so let's just clear for now or delete and recreate.
				s.repo.DeleteCart(cart.ID)
				cart = &Cart{
					ID:         uuid.NewString(),
					CustomerID: customerID,
					BranchID:   branchID,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				s.repo.CreateCart(cart)
			}
		}
	} else {
		// Kasir gak bisa pakai endpoint ini tanpa InitOfflineCart dulu
		return errors.New("unauthorized: offline cart requires cart initialization")
	}

	item := &CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Notes:     req.Notes,
	}

	return s.repo.AddOrUpdateCartItem(item)
}

func (s *service) InitOfflineCart(ctx context.Context, branchID int, req InitOfflineCartRequest) (*Cart, error) {
	cart := &Cart{
		ID:        uuid.NewString(),
		CartName:  &req.CartName,
		BranchID:  branchID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreateCart(cart); err != nil {
		return nil, err
	}
	return cart, nil
}

func (s *service) AddItemToOfflineCart(ctx context.Context, cartID string, branchID int, req AddCartItemRequest) error {
	cart, err := s.repo.GetCartByID(cartID)
	if err != nil {
		return err
	}
	if cart == nil {
		return errors.New("not found: cart not found")
	}
	if cart.BranchID != branchID {
		return errors.New("forbidden: cart does not belong to your branch")
	}
	if cart.CustomerID != nil {
		return errors.New("invalid: this cart belongs to a customer, not an offline hold bill")
	}

	item := &CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Notes:     req.Notes,
	}
	return s.repo.AddOrUpdateCartItem(item)
}

func (s *service) GetMyCart(ctx context.Context, customerID string) (*Cart, error) {
	return s.repo.GetActiveCartByCustomer(customerID)
}

func (s *service) GetOfflineCarts(ctx context.Context, branchID int) ([]Cart, error) {
	return s.repo.GetActiveCartsByBranch(branchID)
}

func (s *service) ClearCart(ctx context.Context, cartID string) error {
	return s.repo.ClearCartItems(cartID)
}

// --- CHECKOUT ---

func (s *service) Checkout(ctx context.Context, customerID *string, cashierID *string, req CheckoutRequest) (*Transaction, error) {
	cart, err := s.repo.GetCartByID(req.CartID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("not found: cart not found")
	}

	// Otorisasi: kalau customerID ada, harus match dengan cart.customerID
	if customerID != nil {
		if cart.CustomerID == nil || *cart.CustomerID != *customerID {
			return nil, errors.New("forbidden: not your cart")
		}
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("invalid: cart is empty")
	}

	var shiftID *string
	if cashierID != nil {
		// Kasir checkout
		shift, err := s.repo.GetOpenShiftByCashier(*cashierID)
		if err != nil {
			return nil, err
		}
		if shift == nil {
			return nil, errors.New("invalid: you must open a shift before checking out")
		}
		shiftID = &shift.ID
	}

	// 1. Hitung Subtotal & Kumpulkan Produk
	var totalAmount float64
	var details []TransactionDetail
	var productIDs []int

	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	boms, err := s.catalogSvc.GetProductsBOM(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	for _, item := range cart.Items {
		if item.Product == nil {
			return nil, errors.New("invalid: product not found in cart item")
		}
		subtotal := float64(item.Quantity) * item.Product.Price
		totalAmount += subtotal
		details = append(details, TransactionDetail{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
			Notes:     item.Notes,
		})
	}

	// Validasi Amount Tendered (Jika Cash)
	if req.PaymentMethod == "CASH" {
		if req.AmountTendered == nil || *req.AmountTendered < totalAmount {
			return nil, errors.New("invalid: amount tendered is less than total amount")
		}
	}

	status := "Paid"
	if req.OrderType == "Online_Pickup" || req.OrderType == "Online_Delivery" {
		status = "Waiting_Payment"
	}

	transaction := &Transaction{
		ID:            uuid.NewString(),
		BranchID:      cart.BranchID,
		CustomerID:    customerID,
		CustomerName:  req.CustomerName,
		CashierID:     cashierID,
		ShiftID:       shiftID,
		OrderType:     req.OrderType,
		PaymentMethod: req.PaymentMethod,
		TotalAmount:   totalAmount,
		Status:        status,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Details:       details,
	}

	// 2. Jalankan DB Transaction: Create Transaksi, Hapus Cart,
	// Jika Offline -> Kurangi Stok, Update Shift (Jika ada)
	err = s.repo.WithTransaction(func(tx *gorm.DB) error {
		// a. Create Transaction
		if err := s.repo.CreateTransaction(tx, transaction); err != nil {
			return err
		}

		if status == "Paid" {
			// b. Potong Stok (BOM) hanya jika langsung lunas (Offline)
			for _, item := range cart.Items {
				recipe := boms[item.ProductID]
				for _, r := range recipe {
					qtyToDeduct := r.QuantityNeeded * float64(item.Quantity)
					desc := "Sales Transaction: " + transaction.ID
					if err := s.invSvc.DeductStock(tx, cart.BranchID, r.MaterialID, qtyToDeduct, desc); err != nil {
						return err
					}
				}
			}

			// c. Tambah Expected Cash di Shift Kasir
			if shiftID != nil && req.PaymentMethod == "CASH" {
				var shift Shift
				if err := tx.First(&shift, "id = ?", *shiftID).Error; err != nil {
					return err
				}
				shift.ExpectedCash += totalAmount
				if err := tx.Save(&shift).Error; err != nil {
					return err
				}
			}
		}

		// d. Hapus Cart
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&CartItem{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&Cart{}, "id = ?", cart.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 3. Notifikasi
	// Online Payment invoice will be sent by Webhook, NOT here.
	// We only send notification here for Offline transactions (maybe E-Receipt if they provide email later)
	// But according to user: "beli offline atau bawa pulang maka bisa bayar cash atau lewat qris gitu... nota fisik". No email for offline.
	// So we don't send invoice email here anymore! It's fully handled by webhook for online.

	// Kasir di-notify saat transaksi lunas (Offline langsung masuk antrian, Online nunggu webhook)
	if s.notifSvc != nil {
		if status == "Paid" {
			// Kasir yg melayani gak butuh notif in-app karena dia yg input.
			// Tapi barista mungkin butuh? Sementara lewati saja.
		} else {
			// Online nunggu lunas, jadi ngga di-notify pas "Waiting_Payment"
			// Webhook yg akan trigger notif "Pesanan Lunas" ke kasir.
		}
	}

	return transaction, nil
}

func (s *service) GetTransactions(ctx context.Context, role string, reqBranchID *int, reqCustomerID *string, status *string, startDate *string, endDate *string) ([]Transaction, error) {
	var finalBranchID *int
	var finalCustomerID *string

	if role == "Customer" {
		if reqCustomerID == nil {
			return nil, errors.New("unauthorized: customer ID is required")
		}
		finalCustomerID = reqCustomerID
	} else if role == "Cashier" || role == "Manager" {
		if reqBranchID == nil {
			return nil, errors.New("forbidden: employee must have a branch assigned")
		}
		finalBranchID = reqBranchID
	} else if role == "Admin" || role == "ADMIN" {
		finalBranchID = reqBranchID // Admin can pass branch filter explicitly
		finalCustomerID = reqCustomerID
	} else {
		return nil, errors.New("unauthorized: invalid role")
	}

	return s.repo.GetTransactions(finalBranchID, finalCustomerID, status, startDate, endDate)
}

func (s *service) GetTransactionByID(ctx context.Context, id string, role string, reqBranchID *int, reqCustomerID *string) (*Transaction, error) {
	trx, err := s.repo.GetTransactionByID(id)
	if err != nil {
		return nil, err
	}
	if trx == nil {
		return nil, errors.New("not found: transaction not found")
	}

	// Otorisasi
	if role == "Customer" {
		if trx.CustomerID == nil || reqCustomerID == nil || *trx.CustomerID != *reqCustomerID {
			return nil, errors.New("forbidden: not your transaction")
		}
	} else if role == "Cashier" || role == "Manager" {
		if reqBranchID == nil || trx.BranchID != *reqBranchID {
			return nil, errors.New("forbidden: transaction does not belong to your branch")
		}
	}

	return trx, nil
}

func (s *service) UpdateTransactionStatus(ctx context.Context, id string, status string, role string, reqBranchID *int) error {
	if role == "Customer" {
		return errors.New("forbidden: customers cannot update transaction status")
	}

	trx, err := s.repo.GetTransactionByID(id)
	if err != nil {
		return err
	}
	if trx == nil {
		return errors.New("not found: transaction not found")
	}

	if (role == "Cashier" || role == "Manager") && (reqBranchID == nil || trx.BranchID != *reqBranchID) {
		return errors.New("forbidden: transaction does not belong to your branch")
	}

	// Validasi Flow Status (Sederhana)
	// Waiting_Payment -> Paid -> Preparing -> Ready -> Completed -> Cancelled
	if trx.Status == "Completed" || trx.Status == "Cancelled" {
		return errors.New("conflict: transaction is already finalized")
	}

	return s.repo.UpdateTransactionStatus(id, status)
}
