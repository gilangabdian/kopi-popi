package sales

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/gilangages/kopi-popi/internal/branch"
	"github.com/gilangages/kopi-popi/internal/catalog"
	"github.com/gilangages/kopi-popi/internal/inventory"
)

// Mock Repo
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateShift(shift *Shift) error {
	args := m.Called(shift)
	return args.Error(0)
}
func (m *MockRepository) UpdateShift(shift *Shift) error {
	args := m.Called(shift)
	return args.Error(0)
}
func (m *MockRepository) GetOpenShiftByCashier(cashierID string) (*Shift, error) {
	args := m.Called(cashierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Shift), args.Error(1)
}
func (m *MockRepository) CreateCart(cart *Cart) error {
	args := m.Called(cart)
	return args.Error(0)
}
func (m *MockRepository) GetCartByID(id string) (*Cart, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Cart), args.Error(1)
}
func (m *MockRepository) GetActiveCartByCustomer(customerID string) (*Cart, error) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Cart), args.Error(1)
}
func (m *MockRepository) GetActiveCartsByBranch(branchID int) ([]Cart, error) {
	args := m.Called(branchID)
	return args.Get(0).([]Cart), args.Error(1)
}
func (m *MockRepository) AddOrUpdateCartItem(item *CartItem) error {
	args := m.Called(item)
	return args.Error(0)
}
func (m *MockRepository) RemoveCartItem(cartID string, productID int) error {
	args := m.Called(cartID, productID)
	return args.Error(0)
}
func (m *MockRepository) ClearCartItems(cartID string) error {
	args := m.Called(cartID)
	return args.Error(0)
}
func (m *MockRepository) DeleteCart(cartID string) error {
	args := m.Called(cartID)
	return args.Error(0)
}
func (m *MockRepository) UpdateTransactionStatus(id string, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}
func (m *MockRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}
func (m *MockRepository) CreateTransaction(tx *gorm.DB, trx *Transaction) error {
	args := m.Called(tx, trx)
	return args.Error(0)
}
func (m *MockRepository) GetTransactionByID(id string) (*Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Transaction), args.Error(1)
}
func (m *MockRepository) GetTransactions(branchID *int, customerID *string, status *string, startDate *string, endDate *string) ([]Transaction, error) {
	args := m.Called(branchID, customerID, status, startDate, endDate)
	return args.Get(0).([]Transaction), args.Error(1)
}
func (m *MockRepository) GetShiftByID(id string) (*Shift, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Shift), args.Error(1)
}
func (m *MockRepository) CreateExpense(tx *gorm.DB, expense *Expense) error {
	args := m.Called(tx, expense)
	return args.Error(0)
}
func (m *MockRepository) GetExpensesByShiftID(shiftID string) ([]Expense, error) {
	args := m.Called(shiftID)
	return args.Get(0).([]Expense), args.Error(1)
}

// Mock Branch Service
type MockBranchService struct {
	mock.Mock
}

func (m *MockBranchService) CreateBranch(ctx context.Context, req branch.CreateBranchRequest) error { return nil }
func (m *MockBranchService) UpdateBranch(ctx context.Context, id int, req branch.UpdateBranchRequest) error { return nil }
func (m *MockBranchService) DeleteBranch(ctx context.Context, id int) error { return nil }
func (m *MockBranchService) GetAllBranches(ctx context.Context, role string, includeInactive bool) ([]branch.Branch, error) {
	args := m.Called(ctx, role, includeInactive)
	return args.Get(0).([]branch.Branch), args.Error(1)
}
func (m *MockBranchService) UpdateOperatingHours(ctx context.Context, id int, req branch.UpdateOperatingHoursRequest, role string) error { return nil }
func (m *MockBranchService) ToggleAcceptingOrders(ctx context.Context, id int, req branch.UpdateAcceptingOrdersRequest, role string, branchID *int) error { return nil }

// Mock Catalog Service
type MockCatalogService struct {
	mock.Mock
}
func (m *MockCatalogService) GetAllCategories(ctx context.Context) ([]catalog.Category, error) { return nil, nil }
func (m *MockCatalogService) CreateCategory(ctx context.Context, req catalog.CategoryRequest) error { return nil }
func (m *MockCatalogService) UpdateCategory(ctx context.Context, id int, req catalog.CategoryRequest) error { return nil }
func (m *MockCatalogService) DeleteCategory(ctx context.Context, id int) error { return nil }

func (m *MockCatalogService) GetAllMaterials(ctx context.Context) ([]catalog.Material, error) { return nil, nil }
func (m *MockCatalogService) CreateMaterial(ctx context.Context, req catalog.MaterialRequest) error { return nil }
func (m *MockCatalogService) UpdateMaterial(ctx context.Context, id int, req catalog.MaterialRequest) error { return nil }
func (m *MockCatalogService) DeleteMaterial(ctx context.Context, id int) error { return nil }

func (m *MockCatalogService) GetAllProducts(ctx context.Context, categoryID *int, search string) ([]catalog.Product, error) { return nil, nil }
func (m *MockCatalogService) GetProductDetail(ctx context.Context, id int, role string, includeRecipe bool) (*catalog.Product, error) { return nil, nil }
func (m *MockCatalogService) GetProductsBOM(ctx context.Context, productIDs []int) (map[int][]catalog.ProductBOM, error) {
	args := m.Called(ctx, productIDs)
	return args.Get(0).(map[int][]catalog.ProductBOM), args.Error(1)
}
func (m *MockCatalogService) CreateProduct(ctx context.Context, req catalog.ProductRequest) error { return nil }
func (m *MockCatalogService) UpdateProduct(ctx context.Context, id int, req catalog.ProductRequest) error { return nil }
func (m *MockCatalogService) DeleteProduct(ctx context.Context, id int) error { return nil }


// Mock Inventory Service
type MockInventoryService struct {
	mock.Mock
}

func (m *MockInventoryService) GetBranchStock(branchID int, requestingRole string, requestingBranchID *int) ([]inventory.BranchInventory, error) { return nil, nil }
func (m *MockInventoryService) GetInventoryMovements(branchID int, requestingRole string, requestingBranchID *int) ([]inventory.InventoryMovement, error) { return nil, nil }
func (m *MockInventoryService) GetRestockRequests(requestingRole string, requestingBranchID *int) ([]inventory.RestockRequest, error) { return nil, nil }
func (m *MockInventoryService) CreateRestockRequest(req *inventory.RestockRequest, requestingRole string, requestingBranchID *int) error { return nil }
func (m *MockInventoryService) UpdateRestockStatus(id string, newStatus string, rejectionReason *string, requestingRole string, requestingBranchID *int) error { return nil }
func (m *MockInventoryService) DeductStock(tx interface{}, branchID int, materialID int, quantity float64, description string) error {
	args := m.Called(tx, branchID, materialID, quantity, description)
	return args.Error(0)
}

func TestOpenShift_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, nil, nil, nil, nil, nil)
	
	cashierID := uuid.NewString()
	mockRepo.On("GetOpenShiftByCashier", cashierID).Return(nil, nil)
	mockRepo.On("CreateShift", mock.AnythingOfType("*sales.Shift")).Return(nil)
	
	req := OpenShiftRequest{StartingCash: 50000}
	shift, err := svc.OpenShift(context.Background(), 1, cashierID, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, shift)
	assert.Equal(t, 50000.0, shift.StartingCash)
	assert.Equal(t, "Open", shift.Status)
}

func TestOpenShift_AlreadyOpen(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, nil, nil, nil, nil, nil)
	
	cashierID := uuid.NewString()
	mockRepo.On("GetOpenShiftByCashier", cashierID).Return(&Shift{Status: "Open"}, nil)
	
	req := OpenShiftRequest{StartingCash: 50000}
	shift, err := svc.OpenShift(context.Background(), 1, cashierID, req)
	
	assert.Error(t, err)
	assert.Nil(t, shift)
	assert.Equal(t, "conflict: you already have an open shift", err.Error())
}

func TestInitOfflineCart(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, nil, nil, nil, nil, nil)
	
	mockRepo.On("CreateCart", mock.AnythingOfType("*sales.Cart")).Return(nil)
	
	cartName := "Meja 1"
	req := InitOfflineCartRequest{CartName: cartName}
	cart, err := svc.InitOfflineCart(context.Background(), 1, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, cart)
	assert.Equal(t, "Meja 1", *cart.CartName)
	assert.Equal(t, 1, cart.BranchID)
}

func TestAddItemToOfflineCart(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, nil, nil, nil, nil, nil)
	
	cartID := uuid.NewString()
	cartName := "Meja 2"
	mockRepo.On("GetCartByID", cartID).Return(&Cart{
		ID: cartID,
		BranchID: 1,
		CartName: &cartName,
	}, nil)
	mockRepo.On("AddOrUpdateCartItem", mock.AnythingOfType("*sales.CartItem")).Return(nil)
	
	req := AddCartItemRequest{ProductID: 1, Quantity: 2}
	err := svc.AddItemToOfflineCart(context.Background(), cartID, 1, req)
	
	assert.NoError(t, err)
}

func TestCheckout_EmptyCart(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, nil, nil, nil, nil, nil)
	
	cartID := uuid.NewString()
	customerID := uuid.NewString()
	
	mockRepo.On("GetCartByID", cartID).Return(&Cart{
		ID: cartID,
		CustomerID: &customerID,
		BranchID: 1,
		Items: []CartItem{}, // empty
	}, nil)
	
	req := CheckoutRequest{CartID: cartID}
	trx, err := svc.Checkout(context.Background(), &customerID, nil, req)
	
	assert.Error(t, err)
	assert.Nil(t, trx)
	assert.Equal(t, "invalid: cart is empty", err.Error())
}
