package promo

import "time"

type Promo struct {
	ID                int        `gorm:"primaryKey" json:"id"`
	Code              string     `gorm:"size:50;not null;unique" json:"code"`
	Title             string     `gorm:"size:100;not null" json:"title"`
	Slug              string     `gorm:"size:255;unique" json:"slug"`
	DiscountType      string     `gorm:"type:enum('PERCENTAGE','FIXED');not null" json:"discount_type"`
	DiscountValue     float64    `gorm:"not null" json:"discount_value"`
	MaxDiscountAmount *float64   `json:"max_discount_amount,omitempty"`
	MinPurchaseAmount float64    `gorm:"not null;default:0" json:"min_purchase_amount"`
	ValidFrom         time.Time  `gorm:"not null" json:"valid_from"`
	ValidUntil        time.Time  `gorm:"not null" json:"valid_until"`
	IsActive          bool       `gorm:"not null;default:true" json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type PromoRequest struct {
	Code              string   `json:"code" binding:"required"`
	Title             string   `json:"title" binding:"required"`
	DiscountType      string   `json:"discount_type" binding:"required,oneof=PERCENTAGE FIXED"`
	DiscountValue     float64  `json:"discount_value" binding:"required,gt=0"`
	MaxDiscountAmount *float64 `json:"max_discount_amount,omitempty"`
	MinPurchaseAmount float64  `json:"min_purchase_amount" binding:"gte=0"`
	ValidFrom         string   `json:"valid_from" binding:"required"`
	ValidUntil        string   `json:"valid_until" binding:"required"`
}

type ValidatePromoRequest struct {
	Code        string  `json:"code" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required,gt=0"`
}

type ValidatePromoResponse struct {
	IsValid        bool    `json:"is_valid"`
	DiscountAmount float64 `json:"discount_amount"`
	FinalAmount    float64 `json:"final_amount"`
	Message        string  `json:"message"`
}
