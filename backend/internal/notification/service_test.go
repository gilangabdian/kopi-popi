package notification

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository adalah mock untuk antarmuka Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveNotification(notification *Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockRepository) GetNotificationsByUser(userID string) ([]Notification, error) {
	args := m.Called(userID)
	return args.Get(0).([]Notification), args.Error(1)
}

func (m *MockRepository) MarkAsRead(id string, userID string) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func TestGetMyNotifications(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	userID := "user-123"

	expectedNotifs := []Notification{
		{ID: "notif-1", UserID: userID, Title: "Test", Message: "Msg", Type: "INFO", IsRead: false, CreatedAt: time.Now()},
	}

	mockRepo.On("GetNotificationsByUser", userID).Return(expectedNotifs, nil)

	notifs, err := service.GetMyNotifications(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, notifs, 1)
	assert.Equal(t, "notif-1", notifs[0].ID)
	mockRepo.AssertExpectations(t)
}

func TestMarkAsRead(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	ctx := context.Background()
	userID := "user-123"
	notifID := "notif-1"

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("MarkAsRead", notifID, userID).Return(nil).Once()

		err := service.MarkAsRead(ctx, userID, notifID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		mockRepo.On("MarkAsRead", notifID, userID).Return(errors.New("not found")).Once()

		err := service.MarkAsRead(ctx, userID, notifID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSendInAppNotification(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	userID := "user-123"

	mockRepo.On("SaveNotification", mock.AnythingOfType("*notification.Notification")).Return(nil).Once()

	err := service.SendInAppNotification(userID, "Low Stock", "Stock is low", "WARNING")
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Untuk unit test Email (Async), karena menggunakan goroutine dan tidak mengembalikan error/value,
// kita cukup jalankan untuk memastikan tidak ada panic/crash saat memanggil templatenya.
// Di dunia nyata mock interface Mailer bisa diterapkan, namun Mock Email log sudah cukup aman di sini.
func TestSendEmails_NoPanic(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	assert.NotPanics(t, func() {
		service.SendRegistrationOTPEmail("test@example.com", "Budi", "123456")
		service.SendOTPResetEmail("test@example.com", "Budi", "123456")
		service.SendInvoiceEmail("test@example.com", "Budi", "TRX-123", 50000)
		service.SendOrderReadyEmail("test@example.com", "Budi", "TRX-123")
		service.SendRestockRequestEmail("admin@example.com", "Cabang 1", "REQ-123")
		service.SendRestockResultEmail("manager@example.com", "Cabang 1", "REQ-123", "Approved")
	})
}
