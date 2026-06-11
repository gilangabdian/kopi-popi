package blogs

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, blog *Blog) error
	FindAll(ctx context.Context, limit int, offset int) ([]Blog, int64, error)
	FindByIDOrSlug(ctx context.Context, idOrSlug string) (*Blog, error)
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
	Update(ctx context.Context, blog *Blog) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, blog *Blog) error {
	return r.db.WithContext(ctx).Create(blog).Error
}

func (r *repository) FindAll(ctx context.Context, limit int, offset int) ([]Blog, int64, error) {
	var blogs []Blog
	var total int64

	err := r.db.WithContext(ctx).Model(&Blog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Order("created_at desc").Limit(limit).Offset(offset).Find(&blogs).Error
	return blogs, total, err
}

func (r *repository) FindByIDOrSlug(ctx context.Context, idOrSlug string) (*Blog, error) {
	var blog Blog
	err := r.db.WithContext(ctx).Where("id = ? OR slug = ?", idOrSlug, idOrSlug).First(&blog).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("blog not found")
		}
		return nil, err
	}
	return &blog, nil
}

func (r *repository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Blog{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) Update(ctx context.Context, blog *Blog) error {
	return r.db.WithContext(ctx).Save(blog).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Blog{}, "id = ?", id).Error
}
