package auth

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	CreatePasswordReset(ctx context.Context, pwReset *PasswordReset) error
	FindPasswordResetByToken(ctx context.Context, token string) (*PasswordReset, error)
	UpdatePassword(ctx context.Context, email string, hashedPassword string) error
	DeletePasswordReset(ctx context.Context, email string) error
	CreateEmailVerification(ctx context.Context, ev *EmailVerification) error
	GetEmailVerification(ctx context.Context, email string) (*EmailVerification, error)
	DeleteEmailVerification(ctx context.Context, email string) error
	VerifyUserEmail(ctx context.Context, email string) error
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

func (r *authRepository) CreatePasswordReset(ctx context.Context, pwReset *PasswordReset) error {
	return r.db.WithContext(ctx).Create(pwReset).Error
}

func (r *authRepository) FindPasswordResetByToken(ctx context.Context, token string) (*PasswordReset, error) {
	var pwReset PasswordReset
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&pwReset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pwReset, nil
}

func (r *authRepository) UpdatePassword(ctx context.Context, email string, hashedPassword string) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Update("password_hash", hashedPassword).Error
}

func (r *authRepository) DeletePasswordReset(ctx context.Context, email string) error {
	return r.db.WithContext(ctx).Where("email = ?", email).Delete(&PasswordReset{}).Error
}

func (r *authRepository) CreateEmailVerification(ctx context.Context, ev *EmailVerification) error {
	return r.db.WithContext(ctx).Create(ev).Error
}

func (r *authRepository) GetEmailVerification(ctx context.Context, email string) (*EmailVerification, error) {
	var ev EmailVerification
	err := r.db.WithContext(ctx).Where("email = ?", email).Order("created_at desc").First(&ev).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ev, nil
}

func (r *authRepository) DeleteEmailVerification(ctx context.Context, email string) error {
	return r.db.WithContext(ctx).Where("email = ?", email).Delete(&EmailVerification{}).Error
}

func (r *authRepository) VerifyUserEmail(ctx context.Context, email string) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Update("is_verified", true).Error
}
