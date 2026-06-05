package inventory_test

import (
	"testing"
	"gorm.io/gorm"

	"github.com/gilangages/kopi-popi/internal/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock for inventory.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetBranchStock(branchID int) ([]inventory.BranchInventory, error) {
	args := m.Called(branchID)
	return args.Get(0).([]inventory.BranchInventory), args.Error(1)
}
func (m *MockRepository) GetInventoryMovements(branchID int) ([]inventory.InventoryMovement, error) {
	args := m.Called(branchID)
	return args.Get(0).([]inventory.InventoryMovement), args.Error(1)
}
func (m *MockRepository) GetRestockRequests(branchID *int) ([]inventory.RestockRequest, error) {
	args := m.Called(branchID)
	return args.Get(0).([]inventory.RestockRequest), args.Error(1)
}
func (m *MockRepository) GetRestockRequestByID(id string) (*inventory.RestockRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*inventory.RestockRequest), args.Error(1)
}
func (m *MockRepository) CreateRestockRequest(req *inventory.RestockRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
func (m *MockRepository) UpdateRestockStatus(id string, status string, rejectionReason *string) error {
	args := m.Called(id, status, rejectionReason)
	return args.Error(0)
}
func (m *MockRepository) MarkAsDeliveredAndAddStock(requestID string) error {
	args := m.Called(requestID)
	return args.Error(0)
}
func (m *MockRepository) DeductStock(tx *gorm.DB, branchID int, materialID int, quantity float64, description string) error {
	args := m.Called(tx, branchID, materialID, quantity, description)
	return args.Error(0)
}

func TestCreateRestockRequest_Unauthorized(t *testing.T) {
	mockRepo := new(MockRepository)
	service := inventory.NewService(mockRepo, nil)

	req := &inventory.RestockRequest{
		Items: []inventory.RestockItem{{MaterialID: 1, QuantityRequested: 10}},
	}
	branchID := 1

	// Act: Cashier tries to create request
	err := service.CreateRestockRequest(req, "CASHIER", &branchID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "forbidden: only manager can create restock requests", err.Error())
	mockRepo.AssertNotCalled(t, "CreateRestockRequest")
}

func TestUpdateStatus_Approved_ByManager(t *testing.T) {
	mockRepo := new(MockRepository)
	service := inventory.NewService(mockRepo, nil)

	reqID := "req-1"
	mockRepo.On("GetRestockRequestByID", reqID).Return(&inventory.RestockRequest{ID: reqID, Status: "Pending"}, nil)

	branchID := 1

	// Act: Manager tries to approve their own request
	err := service.UpdateRestockStatus(reqID, "Approved", nil, "MANAGER", &branchID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "forbidden: only admin can approve or reject", err.Error())
	mockRepo.AssertNotCalled(t, "UpdateRestockStatus")
}

func TestMarkAsDelivered_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := inventory.NewService(mockRepo, nil)

	reqID := "req-1"
	branchID := 1
	mockRepo.On("GetRestockRequestByID", reqID).Return(&inventory.RestockRequest{ID: reqID, BranchID: branchID, Status: "Approved"}, nil)
	mockRepo.On("MarkAsDeliveredAndAddStock", reqID).Return(nil)

	// Act: Manager marks as delivered
	err := service.UpdateRestockStatus(reqID, "Delivered", nil, "MANAGER", &branchID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "MarkAsDeliveredAndAddStock", reqID)
}

func TestMarkAsDelivered_WrongBranch(t *testing.T) {
	mockRepo := new(MockRepository)
	service := inventory.NewService(mockRepo, nil)

	reqID := "req-1"
	actualBranchID := 2
	mockRepo.On("GetRestockRequestByID", reqID).Return(&inventory.RestockRequest{ID: reqID, BranchID: actualBranchID, Status: "Approved"}, nil)

	managerBranchID := 1

	// Act: Manager of Branch 1 tries to deliver a request for Branch 2
	err := service.UpdateRestockStatus(reqID, "Delivered", nil, "MANAGER", &managerBranchID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "forbidden: can only mark delivery for your own branch", err.Error())
	mockRepo.AssertNotCalled(t, "MarkAsDeliveredAndAddStock")
}

func TestUpdateStatus_Rejected_MissingReason(t *testing.T) {
	mockRepo := new(MockRepository)
	service := inventory.NewService(mockRepo, nil)

	reqID := "req-1"
	mockRepo.On("GetRestockRequestByID", reqID).Return(&inventory.RestockRequest{ID: reqID, Status: "Pending"}, nil)

	// Act: Admin tries to reject without reason
	err := service.UpdateRestockStatus(reqID, "Rejected", nil, "ADMIN", nil)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid: rejection_reason is mandatory when rejecting", err.Error())
	mockRepo.AssertNotCalled(t, "UpdateRestockStatus")
}

func TestUpdateStatus_Rejected_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := inventory.NewService(mockRepo, nil)

	reqID := "req-1"
	reason := "Out of stock"
	mockRepo.On("GetRestockRequestByID", reqID).Return(&inventory.RestockRequest{ID: reqID, Status: "Pending"}, nil)
	mockRepo.On("UpdateRestockStatus", reqID, "Rejected", &reason).Return(nil)

	// Act: Admin rejects with reason
	err := service.UpdateRestockStatus(reqID, "Rejected", &reason, "ADMIN", nil)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "UpdateRestockStatus", reqID, "Rejected", &reason)
}


