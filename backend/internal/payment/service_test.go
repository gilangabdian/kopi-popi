package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/gilangages/kopi-popi/internal/notification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCK REPOSITORY ---
type mockPaymentRepo struct {
	mock.Mock
}

func (m *mockPaymentRepo) UpdateTransactionStatusToPaid(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func (m *mockPaymentRepo) UpdateTransactionStatusToFailed(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func (m *mockPaymentRepo) GetTransactionInfo(ctx context.Context, orderID string) (*TransactionInfo, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TransactionInfo), args.Error(1)
}

func (m *mockPaymentRepo) DeductStockForTransaction(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

// --- MOCK NOTIFICATION SERVICE ---
type mockNotificationService struct {
	mock.Mock
}

func (m *mockNotificationService) SendRegistrationOTPEmail(email string, userName string, otpCode string) {}
func (m *mockNotificationService) SendOTPResetEmail(email string, userName string, otpCode string)        {}
func (m *mockNotificationService) GetMyNotifications(ctx context.Context, userID string) ([]notification.Notification, error) {
	return nil, nil
}
func (m *mockNotificationService) MarkAsRead(ctx context.Context, userID string, id string) error {
	return nil
}
func (m *mockNotificationService) SendInvoiceEmail(customerEmail string, customerName string, transactionID string, amount float64) {
	m.Called(customerEmail, customerName, transactionID, amount)
}
func (m *mockNotificationService) SendOrderReadyEmail(customerEmail string, customerName string, transactionID string) {}
func (m *mockNotificationService) SendRestockRequestEmail(adminEmail string, branchName string, reqID string)        {}
func (m *mockNotificationService) SendRestockResultEmail(managerEmail string, branchName string, status string, notes string) {}
func (m *mockNotificationService) CreateInAppNotification(ctx context.Context, userID string, title string, message string, entityType string, entityID string) error {
	return nil
}
func (m *mockNotificationService) SendInAppNotification(roleName string, title string, message string, wsMessage string) error {
	args := m.Called(roleName, title, message, wsMessage)
	return args.Error(0)
}

// --- UNIT TESTS ---

func TestHandleWebhook_Success_Settlement(t *testing.T) {
	mockRepo := new(mockPaymentRepo)
	mockNotif := new(mockNotificationService)

	svc := NewService(mockRepo, mockNotif)

	orderID := "ORDER-123"
	payload := map[string]interface{}{
		"order_id":           orderID,
		"transaction_status": "settlement",
	}

	txInfo := &TransactionInfo{
		ID:            orderID,
		TotalAmount:   50000,
		CustomerName:  "Gilang",
		CustomerEmail: "gilang@example.com",
	}

	mockRepo.On("UpdateTransactionStatusToPaid", mock.Anything, orderID).Return(nil)
	mockRepo.On("GetTransactionInfo", mock.Anything, orderID).Return(txInfo, nil)
	mockRepo.On("DeductStockForTransaction", mock.Anything, orderID).Return(nil)

	mockNotif.On("SendInvoiceEmail", "gilang@example.com", "Gilang", orderID, float64(50000)).Return()

	err := svc.HandleWebhook(context.Background(), payload)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockNotif.AssertExpectations(t)
}

func TestHandleWebhook_Success_Cancel(t *testing.T) {
	mockRepo := new(mockPaymentRepo)
	svc := NewService(mockRepo, nil)

	orderID := "ORDER-456"
	payload := map[string]interface{}{
		"order_id":           orderID,
		"transaction_status": "cancel",
	}

	mockRepo.On("UpdateTransactionStatusToFailed", mock.Anything, orderID).Return(nil)

	err := svc.HandleWebhook(context.Background(), payload)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHandleWebhook_InvalidPayload(t *testing.T) {
	svc := NewService(nil, nil)

	// Missing order_id
	payload := map[string]interface{}{
		"transaction_status": "settlement",
	}

	err := svc.HandleWebhook(context.Background(), payload)
	assert.Error(t, err)
	assert.Equal(t, "invalid order_id in webhook", err.Error())
}

func TestHandleWebhook_UpdateStatusFailed(t *testing.T) {
	mockRepo := new(mockPaymentRepo)
	svc := NewService(mockRepo, nil)

	orderID := "ORDER-123"
	payload := map[string]interface{}{
		"order_id":           orderID,
		"transaction_status": "settlement",
	}

	expectedErr := errors.New("db error")
	mockRepo.On("UpdateTransactionStatusToPaid", mock.Anything, orderID).Return(expectedErr)

	err := svc.HandleWebhook(context.Background(), payload)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}
