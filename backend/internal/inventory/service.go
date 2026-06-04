package inventory

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Service interface {
	GetBranchStock(branchID int, requestingRole string, requestingBranchID *int) ([]BranchInventory, error)
	GetInventoryMovements(branchID int, requestingRole string, requestingBranchID *int) ([]InventoryMovement, error)
	GetRestockRequests(requestingRole string, requestingBranchID *int) ([]RestockRequest, error)
	CreateRestockRequest(req *RestockRequest, requestingRole string, requestingBranchID *int) error
	UpdateRestockStatus(id string, newStatus string, rejectionReason *string, requestingRole string, requestingBranchID *int) error
	DeductStock(tx interface{}, branchID int, materialID int, quantity float64, description string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetBranchStock(branchID int, requestingRole string, requestingBranchID *int) ([]BranchInventory, error) {
	requestingRole = strings.ToUpper(requestingRole)
	if requestingRole != "ADMIN" && (requestingBranchID == nil || *requestingBranchID != branchID) {
		return nil, errors.New("forbidden: can only access your own branch stock")
	}
	return s.repo.GetBranchStock(branchID)
}

func (s *service) GetInventoryMovements(branchID int, requestingRole string, requestingBranchID *int) ([]InventoryMovement, error) {
	requestingRole = strings.ToUpper(requestingRole)
	if requestingRole != "ADMIN" && (requestingBranchID == nil || *requestingBranchID != branchID) {
		return nil, errors.New("forbidden: can only access your own branch movements")
	}
	return s.repo.GetInventoryMovements(branchID)
}

func (s *service) GetRestockRequests(requestingRole string, requestingBranchID *int) ([]RestockRequest, error) {
	requestingRole = strings.ToUpper(requestingRole)
	if requestingRole == "ADMIN" {
		return s.repo.GetRestockRequests(nil) // Get all
	}
	if requestingBranchID == nil {
		return nil, errors.New("forbidden: branch ID missing for manager")
	}
	return s.repo.GetRestockRequests(requestingBranchID)
}

func (s *service) CreateRestockRequest(req *RestockRequest, requestingRole string, requestingBranchID *int) error {
	requestingRole = strings.ToUpper(requestingRole)
	if requestingRole != "MANAGER" {
		return errors.New("forbidden: only manager can create restock requests")
	}
	if requestingBranchID == nil {
		return errors.New("forbidden: branch ID missing")
	}
	
	// Force the branch ID to be the manager's branch
	req.BranchID = *requestingBranchID
	req.Status = "Pending"
	
	if len(req.Items) == 0 {
		return errors.New("invalid: request must have at least one item")
	}

	return s.repo.CreateRestockRequest(req)
}

func (s *service) UpdateRestockStatus(id string, newStatus string, rejectionReason *string, requestingRole string, requestingBranchID *int) error {
	requestingRole = strings.ToUpper(requestingRole)
	req, err := s.repo.GetRestockRequestByID(id)
	if err != nil {
		return err
	}
	if req == nil {
		return errors.New("not found: restock request not found")
	}

	// Admin actions: Approve or Reject
	if newStatus == "Approved" || newStatus == "Rejected" {
		if requestingRole != "ADMIN" {
			return errors.New("forbidden: only admin can approve or reject")
		}
		if req.Status != "Pending" {
			return errors.New("conflict: can only approve/reject pending requests")
		}
		if newStatus == "Rejected" && (rejectionReason == nil || *rejectionReason == "") {
			return errors.New("invalid: rejection_reason is mandatory when rejecting")
		}
		if newStatus == "Approved" {
			rejectionReason = nil
		}
		return s.repo.UpdateRestockStatus(id, newStatus, rejectionReason)
	}

	// Manager action: Mark as Delivered
	if newStatus == "Delivered" {
		if requestingRole != "MANAGER" {
			return errors.New("forbidden: only manager can mark as delivered")
		}
		if requestingBranchID == nil || *requestingBranchID != req.BranchID {
			return errors.New("forbidden: can only mark delivery for your own branch")
		}
		if req.Status != "Approved" {
			return errors.New("conflict: can only deliver approved requests")
		}
		
		// This handles both changing status and updating inventory via DB Transaction
		return s.repo.MarkAsDeliveredAndAddStock(id)
	}

	return errors.New("invalid: unknown status")
}

func (s *service) DeductStock(tx interface{}, branchID int, materialID int, quantity float64, description string) error {
	importGormDB, ok := tx.(*gorm.DB)
	if !ok && tx != nil {
		return errors.New("invalid transaction type")
	}
	return s.repo.DeductStock(importGormDB, branchID, materialID, quantity, description)
}
