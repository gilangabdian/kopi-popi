package branch

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context, includeInactive bool) ([]Branch, error)
	FindByID(ctx context.Context, id int) (*Branch, error)
	FindByIDOrSlug(ctx context.Context, idOrSlug string) (*Branch, error)
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
	Create(ctx context.Context, branch *Branch) error
	Update(ctx context.Context, branch *Branch) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAll(ctx context.Context, includeInactive bool) ([]Branch, error) {
	var branches []Branch
	query := r.db.WithContext(ctx)
	
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	err := query.Find(&branches).Error
	return branches, err
}

func (r *repository) FindByID(ctx context.Context, id int) (*Branch, error) {
	var branch Branch
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&branch).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &branch, nil
}

func (r *repository) FindByIDOrSlug(ctx context.Context, idOrSlug string) (*Branch, error) {
	var branch Branch
	err := r.db.WithContext(ctx).Where("id = ? OR slug = ?", idOrSlug, idOrSlug).First(&branch).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &branch, nil
}

func (r *repository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Branch{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

func (r *repository) Create(ctx context.Context, branch *Branch) error {
	return r.db.WithContext(ctx).Create(branch).Error
}

func (r *repository) Update(ctx context.Context, branch *Branch) error {
	return r.db.WithContext(ctx).Save(branch).Error
}
