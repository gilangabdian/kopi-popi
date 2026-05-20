package auth

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &authRepository{db}
}

func (r *authRepository) CreateUser(ctx context.Context, user *User) error {
	// Jika role_id tidak di-set, kita set secara otomatis ke Customer
	if user.RoleID == 0 {
		var roleID int
		err := r.db.WithContext(ctx).Table("roles").Select("id").Where("name = ?", "Customer").Scan(&roleID).Error
		if err != nil {
			return err
		}
		if roleID == 0 {
			return errors.New("customer role not found in database")
		}
		user.RoleID = roleID
	}

	return r.db.WithContext(ctx).Create(user).Error
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &user, nil
}
