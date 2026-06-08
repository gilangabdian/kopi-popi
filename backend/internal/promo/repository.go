package promo

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	CreatePromo(promo *Promo) error
	UpdatePromo(promo *Promo) error
	GetPromoByID(id int) (*Promo, error)
	GetPromoByCode(code string) (*Promo, error)
	GetActivePromos() ([]Promo, error)
	GetAllPromos() ([]Promo, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreatePromo(promo *Promo) error {
	return r.db.Create(promo).Error
}

func (r *repository) UpdatePromo(promo *Promo) error {
	return r.db.Save(promo).Error
}

func (r *repository) GetPromoByID(id int) (*Promo, error) {
	var promo Promo
	err := r.db.First(&promo, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("promo not found")
		}
		return nil, err
	}
	return &promo, nil
}

func (r *repository) GetPromoByCode(code string) (*Promo, error) {
	var promo Promo
	err := r.db.Where("code = ?", code).First(&promo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("promo code not found")
		}
		return nil, err
	}
	return &promo, nil
}

func (r *repository) GetActivePromos() ([]Promo, error) {
	var promos []Promo
	// Active promo: is_active = true, and current time is between valid_from and valid_until
	err := r.db.Where("is_active = ? AND valid_from <= NOW() AND valid_until >= NOW()", true).
		Order("created_at desc").Find(&promos).Error
	return promos, err
}

func (r *repository) GetAllPromos() ([]Promo, error) {
	var promos []Promo
	err := r.db.Order("created_at desc").Find(&promos).Error
	return promos, err
}
