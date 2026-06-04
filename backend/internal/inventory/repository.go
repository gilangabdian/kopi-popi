package inventory

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetBranchStock(branchID int) ([]BranchInventory, error)
	GetInventoryMovements(branchID int) ([]InventoryMovement, error)
	GetRestockRequests(branchID *int) ([]RestockRequest, error)
	GetRestockRequestByID(id string) (*RestockRequest, error)
	CreateRestockRequest(req *RestockRequest) error
	UpdateRestockStatus(id string, status string, rejectionReason *string) error
	MarkAsDeliveredAndAddStock(requestID string) error
	DeductStock(tx *gorm.DB, branchID int, materialID int, quantity float64, description string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetBranchStock(branchID int) ([]BranchInventory, error) {
	var stocks []BranchInventory
	err := r.db.Preload("Material").Where("branch_id = ?", branchID).Find(&stocks).Error
	return stocks, err
}

func (r *repository) GetInventoryMovements(branchID int) ([]InventoryMovement, error) {
	var movements []InventoryMovement
	err := r.db.Preload("Material").Where("branch_id = ?", branchID).Order("created_at desc").Find(&movements).Error
	return movements, err
}

func (r *repository) GetRestockRequests(branchID *int) ([]RestockRequest, error) {
	var requests []RestockRequest
	query := r.db.Preload("Items").Preload("Items.Material")
	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}
	err := query.Order("created_at desc").Find(&requests).Error
	return requests, err
}

func (r *repository) GetRestockRequestByID(id string) (*RestockRequest, error) {
	var req RestockRequest
	err := r.db.Preload("Items").Preload("Items.Material").First(&req, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

func (r *repository) CreateRestockRequest(req *RestockRequest) error {
	req.ID = uuid.NewString()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	return r.db.Create(req).Error
}

func (r *repository) UpdateRestockStatus(id string, newStatus string, rejectionReason *string) error {
	updates := map[string]interface{}{
		"status": newStatus,
	}
	if rejectionReason != nil {
		updates["rejection_reason"] = *rejectionReason
	}
	return r.db.Model(&RestockRequest{}).Where("id = ?", id).Updates(updates).Error
}

// MarkAsDeliveredAndAddStock uses a database transaction to ensure atomicity
func (r *repository) MarkAsDeliveredAndAddStock(requestID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Dapatkan request dan items-nya
		var req RestockRequest
		if err := tx.Preload("Items").First(&req, "id = ?", requestID).Error; err != nil {
			return err
		}

		if req.Status == "Delivered" {
			return errors.New("request has already been delivered")
		}

		// 2. Ubah status jadi Delivered
		if err := tx.Model(&req).Update("status", "Delivered").Error; err != nil {
			return err
		}

		// 3. Proses tiap item
		for _, item := range req.Items {
			// a. Tambahkan Movement
			movement := InventoryMovement{
				ID:           uuid.NewString(),
				BranchID:     req.BranchID,
				MaterialID:   item.MaterialID,
				MovementType: "IN",
				Quantity:     item.QuantityRequested,
				Description:  "Restock Delivered - Request ID: " + req.ID,
				CreatedAt:    time.Now(),
			}
			if err := tx.Create(&movement).Error; err != nil {
				return err
			}

			// b. Update Branch Inventory (Upsert)
			var inv BranchInventory
			err := tx.Where("branch_id = ? AND material_id = ?", req.BranchID, item.MaterialID).First(&inv).Error
			
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err // Error DB lainnya
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Insert baru jika material belum pernah ada di cabang
				inv = BranchInventory{
					BranchID:   req.BranchID,
					MaterialID: item.MaterialID,
					Quantity:   item.QuantityRequested,
				}
				if err := tx.Create(&inv).Error; err != nil {
					return err
				}
			} else {
				// Update quantity jika sudah ada
				if err := tx.Model(&inv).Update("quantity", inv.Quantity+item.QuantityRequested).Error; err != nil {
					return err
				}
			}
		}

		// Semua sukses, commit otomatis
		return nil
	})
}

// DeductStock dipanggil saat transaksi sales berhasil. Kita biarkan stok negatif untuk kebebasan aplikasi POS
func (r *repository) DeductStock(tx *gorm.DB, branchID int, materialID int, quantity float64, description string) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	// 1. Catat Movement
	movement := InventoryMovement{
		ID:           uuid.NewString(),
		BranchID:     branchID,
		MaterialID:   materialID,
		MovementType: "OUT",
		Quantity:     quantity,
		Description:  description,
		CreatedAt:    time.Now(),
	}
	if err := db.Create(&movement).Error; err != nil {
		return err
	}

	// 2. Potong Stok (Upsert agar bisa langsung minus jika barang baru di menu tapi stok di DB belum diisi)
	var inv BranchInventory
	err := db.Where("branch_id = ? AND material_id = ?", branchID, materialID).First(&inv).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		inv = BranchInventory{
			BranchID:   branchID,
			MaterialID: materialID,
			Quantity:   -quantity, // Stok jadi minus karena transaksi, tapi data stok asli belum pernah dimasukkan.
		}
		return db.Create(&inv).Error
	} else {
		return db.Model(&inv).Update("quantity", inv.Quantity-quantity).Error
	}
}
