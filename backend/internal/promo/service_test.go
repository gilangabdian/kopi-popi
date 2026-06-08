package promo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreatePromo(promo *Promo) error {
	args := m.Called(promo)
	return args.Error(0)
}
func (m *MockRepository) UpdatePromo(promo *Promo) error {
	args := m.Called(promo)
	return args.Error(0)
}
func (m *MockRepository) GetPromoByID(id int) (*Promo, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Promo), args.Error(1)
}
func (m *MockRepository) GetPromoByCode(code string) (*Promo, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Promo), args.Error(1)
}
func (m *MockRepository) GetActivePromos() ([]Promo, error) {
	args := m.Called()
	return args.Get(0).([]Promo), args.Error(1)
}
func (m *MockRepository) GetAllPromos() ([]Promo, error) {
	args := m.Called()
	return args.Get(0).([]Promo), args.Error(1)
}

func TestValidatePromo_SuccessFixed(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	now := time.Now()
	promo := &Promo{
		Code:              "FIXED10",
		DiscountType:      "FIXED",
		DiscountValue:     10000,
		MinPurchaseAmount: 50000,
		ValidFrom:         now.Add(-time.Hour),
		ValidUntil:        now.Add(time.Hour),
		IsActive:          true,
	}

	mockRepo.On("GetPromoByCode", "FIXED10").Return(promo, nil)

	resp, err := svc.ValidatePromo(context.Background(), "FIXED10", 100000)
	assert.NoError(t, err)
	assert.True(t, resp.IsValid)
	assert.Equal(t, 10000.0, resp.DiscountAmount)
	assert.Equal(t, 90000.0, resp.FinalAmount)
}

func TestValidatePromo_SuccessPercentageWithMax(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	now := time.Now()
	maxDiscount := 15000.0
	promo := &Promo{
		Code:              "DISC20",
		DiscountType:      "PERCENTAGE",
		DiscountValue:     20, // 20%
		MaxDiscountAmount: &maxDiscount,
		MinPurchaseAmount: 50000,
		ValidFrom:         now.Add(-time.Hour),
		ValidUntil:        now.Add(time.Hour),
		IsActive:          true,
	}

	mockRepo.On("GetPromoByCode", "DISC20").Return(promo, nil)

	// 20% of 100,000 = 20,000. Max is 15,000.
	resp, err := svc.ValidatePromo(context.Background(), "DISC20", 100000)
	assert.NoError(t, err)
	assert.True(t, resp.IsValid)
	assert.Equal(t, 15000.0, resp.DiscountAmount)
	assert.Equal(t, 85000.0, resp.FinalAmount)
}

func TestValidatePromo_BelowMinPurchase(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	now := time.Now()
	promo := &Promo{
		Code:              "FIXED10",
		DiscountType:      "FIXED",
		DiscountValue:     10000,
		MinPurchaseAmount: 50000,
		ValidFrom:         now.Add(-time.Hour),
		ValidUntil:        now.Add(time.Hour),
		IsActive:          true,
	}

	mockRepo.On("GetPromoByCode", "FIXED10").Return(promo, nil)

	resp, err := svc.ValidatePromo(context.Background(), "FIXED10", 40000)
	assert.NoError(t, err)
	assert.False(t, resp.IsValid)
	assert.Equal(t, "Total belanja belum memenuhi syarat minimum penggunaan promo", resp.Message)
}

func TestCalculateDiscount_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	now := time.Now()
	promo := &Promo{
		Code:              "DISC10",
		DiscountType:      "PERCENTAGE",
		DiscountValue:     10,
		MinPurchaseAmount: 0,
		ValidFrom:         now.Add(-time.Hour),
		ValidUntil:        now.Add(time.Hour),
		IsActive:          true,
	}

	mockRepo.On("GetPromoByCode", "DISC10").Return(promo, nil)

	discount, final, err := svc.CalculateDiscount("DISC10", 50000)
	assert.NoError(t, err)
	assert.Equal(t, 5000.0, discount)
	assert.Equal(t, 45000.0, final)
}
