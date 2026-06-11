package promo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gosimple/slug"
)

type Service interface {
	CreatePromo(ctx context.Context, req PromoRequest) error
	UpdatePromo(ctx context.Context, id int, req PromoRequest) error
	GetPromoByIDOrSlug(ctx context.Context, idOrSlug string) (*Promo, error)
	GetPromos(ctx context.Context, role string) ([]Promo, error)
	ValidatePromo(ctx context.Context, code string, totalAmount float64) (*ValidatePromoResponse, error)
	CalculateDiscount(code string, totalAmount float64) (discountAmount float64, finalAmount float64, err error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) generateUniqueSlug(title string, currentID int) (string, error) {
	baseSlug := slug.Make(title)
	finalSlug := baseSlug
	counter := 1

	for {
		exists, err := s.repo.CheckSlugExists(finalSlug)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
		existingPromo, _ := s.repo.GetPromoByIDOrSlug(finalSlug)
		if existingPromo != nil && existingPromo.ID == currentID {
			break
		}
		finalSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
	return finalSlug, nil
}

func (s *service) CreatePromo(ctx context.Context, req PromoRequest) error {
	validFrom, err := time.Parse(time.RFC3339, req.ValidFrom)
	if err != nil {
		return errors.New("invalid valid_from format, use RFC3339")
	}

	validUntil, err := time.Parse(time.RFC3339, req.ValidUntil)
	if err != nil {
		return errors.New("invalid valid_until format, use RFC3339")
	}

	if validFrom.After(validUntil) {
		return errors.New("valid_from cannot be after valid_until")
	}

	finalSlug, err := s.generateUniqueSlug(req.Title, 0)
	if err != nil {
		return err
	}

	promo := &Promo{
		Code:              req.Code,
		Title:             req.Title,
		Slug:              finalSlug,
		DiscountType:      req.DiscountType,
		DiscountValue:     req.DiscountValue,
		MaxDiscountAmount: req.MaxDiscountAmount,
		MinPurchaseAmount: req.MinPurchaseAmount,
		ValidFrom:         validFrom,
		ValidUntil:        validUntil,
		IsActive:          true,
	}

	return s.repo.CreatePromo(promo)
}

func (s *service) UpdatePromo(ctx context.Context, id int, req PromoRequest) error {
	promo, err := s.repo.GetPromoByID(id)
	if err != nil {
		return err
	}

	validFrom, err := time.Parse(time.RFC3339, req.ValidFrom)
	if err != nil {
		return errors.New("invalid valid_from format, use RFC3339")
	}

	validUntil, err := time.Parse(time.RFC3339, req.ValidUntil)
	if err != nil {
		return errors.New("invalid valid_until format, use RFC3339")
	}

	if validFrom.After(validUntil) {
		return errors.New("valid_from cannot be after valid_until")
	}

	finalSlug, err := s.generateUniqueSlug(req.Title, promo.ID)
	if err != nil {
		return err
	}

	promo.Code = req.Code
	promo.Title = req.Title
	promo.Slug = finalSlug
	promo.DiscountType = req.DiscountType
	promo.DiscountValue = req.DiscountValue
	promo.MaxDiscountAmount = req.MaxDiscountAmount
	promo.MinPurchaseAmount = req.MinPurchaseAmount
	promo.ValidFrom = validFrom
	promo.ValidUntil = validUntil

	return s.repo.UpdatePromo(promo)
}

func (s *service) GetPromoByIDOrSlug(ctx context.Context, idOrSlug string) (*Promo, error) {
	return s.repo.GetPromoByIDOrSlug(idOrSlug)
}

func (s *service) GetPromos(ctx context.Context, role string) ([]Promo, error) {
	if role == "Admin" {
		return s.repo.GetAllPromos()
	}
	return s.repo.GetActivePromos()
}

func (s *service) ValidatePromo(ctx context.Context, code string, totalAmount float64) (*ValidatePromoResponse, error) {
	promo, err := s.repo.GetPromoByCode(code)
	if err != nil {
		return &ValidatePromoResponse{
			IsValid: false,
			Message: "Kode promo tidak ditemukan",
		}, nil
	}

	if !promo.IsActive {
		return &ValidatePromoResponse{
			IsValid: false,
			Message: "Kode promo sudah tidak aktif",
		}, nil
	}

	now := time.Now()
	if now.Before(promo.ValidFrom) || now.After(promo.ValidUntil) {
		return &ValidatePromoResponse{
			IsValid: false,
			Message: "Kode promo sudah kadaluarsa atau belum dimulai",
		}, nil
	}

	if totalAmount < promo.MinPurchaseAmount {
		return &ValidatePromoResponse{
			IsValid: false,
			Message: "Total belanja belum memenuhi syarat minimum penggunaan promo",
		}, nil
	}

	discountAmount := 0.0
	if promo.DiscountType == "FIXED" {
		discountAmount = promo.DiscountValue
	} else if promo.DiscountType == "PERCENTAGE" {
		discountAmount = totalAmount * (promo.DiscountValue / 100.0)
		if promo.MaxDiscountAmount != nil && discountAmount > *promo.MaxDiscountAmount {
			discountAmount = *promo.MaxDiscountAmount
		}
	}

	// Discount cannot be larger than the total amount
	if discountAmount > totalAmount {
		discountAmount = totalAmount
	}

	finalAmount := totalAmount - discountAmount

	return &ValidatePromoResponse{
		IsValid:        true,
		DiscountAmount: discountAmount,
		FinalAmount:    finalAmount,
		Message:        "Promo berhasil digunakan!",
	}, nil
}

func (s *service) CalculateDiscount(code string, totalAmount float64) (float64, float64, error) {
	resp, err := s.ValidatePromo(context.Background(), code, totalAmount)
	if err != nil {
		return 0, totalAmount, err
	}
	if !resp.IsValid {
		return 0, totalAmount, errors.New(resp.Message)
	}
	return resp.DiscountAmount, resp.FinalAmount, nil
}
