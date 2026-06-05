package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gilangages/kopi-popi/internal/notification"
	"github.com/gilangages/kopi-popi/pkg/hash"
	"github.com/gilangages/kopi-popi/pkg/jwt"
	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*User, error)
	VerifyEmail(ctx context.Context, req VerifyEmailRequest) error
	Login(ctx context.Context, req LoginRequest) (string, *User, error)
	ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req ResetPasswordRequest) error
}

type authService struct {
	repo         Repository
	notifService notification.Service
}

func NewService(repo Repository, notif notification.Service) Service {
	return &authService{repo: repo, notifService: notif}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	// 1. Cek ketersediaan email
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email is already registered")
	}

	// 2. Hash password
	hashedPassword, err := hash.MakeHash(req.Password)
	if err != nil {
		return nil, err
	}

	// 3. Buat User ID baru
	newID := uuid.New().String()

	// 4. Siapkan object user
	var phonePtr *string
	if req.Phone != "" {
		phonePtr = &req.Phone
	}

	user := &User{
		ID:           newID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        phonePtr,
	}

	// 5. Simpan ke database
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// 6. Generate OTP (6 digit)
	otpStr := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

	ev := &EmailVerification{
		Email:     user.Email,
		OTPCode:   otpStr,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	_ = s.repo.DeleteEmailVerification(ctx, user.Email) // Clean up old
	err = s.repo.CreateEmailVerification(ctx, ev)
	if err != nil {
		return nil, err
	}

	// 7. Kirim Registration OTP Email (Async)
	if s.notifService != nil {
		s.notifService.SendRegistrationOTPEmail(user.Email, user.Name, otpStr)
	}

	return user, nil
}

func (s *authService) VerifyEmail(ctx context.Context, req VerifyEmailRequest) error {
	// 1. Dapatkan OTP dari database
	ev, err := s.repo.GetEmailVerification(ctx, req.Email)
	if err != nil {
		return err
	}
	if ev == nil {
		return errors.New("tidak ada permintaan verifikasi untuk email ini")
	}

	// 2. Cek kecocokan
	if ev.OTPCode != req.OTPCode {
		return errors.New("kode OTP salah")
	}

	// 3. Cek expired
	if time.Now().After(ev.ExpiresAt) {
		return errors.New("kode OTP telah kedaluwarsa")
	}

	// 4. Update user
	err = s.repo.VerifyUserEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	// 5. Hapus data verifikasi
	_ = s.repo.DeleteEmailVerification(ctx, req.Email)

	return nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (string, *User, error) {
	// 1. Cari user berdasarkan email
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid email or password")
	}

	// 2. Cocokkan hash password
	isValid := hash.CheckHash(req.Password, user.PasswordHash)
	if !isValid {
		return "", nil, errors.New("invalid email or password")
	}

	// 2.5 Cek status verifikasi email
	if !user.IsVerified {
		return "", nil, errors.New("Mohon verifikasi email Anda terlebih dahulu")
	}

	// 3. Generate JWT
	// Asumsi default role name
	roleName := "Customer"
	if user.RoleID == 1 { // Asumsi ID 1 adalah Admin
		roleName = "Admin"
	} else if user.RoleID == 2 {
		roleName = "Manager"
	} else if user.RoleID == 3 {
		roleName = "Cashier"
	}

	token, err := jwt.GenerateToken(user.ID, user.Name, roleName, user.BranchID, req.RememberMe)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	// 1. Cek apakah email terdaftar
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		// Security best practice: Jangan beri tahu email tidak ada,
		// pura-pura sukses agar tidak bisa ditebak (enumeration)
		return nil
	}

	// 2. Generate Reset Token (UUID biasa sudah cukup aman)
	resetToken := uuid.New().String()

	// 3. Simpan ke Database
	pwReset := &PasswordReset{
		Email:     req.Email,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(time.Hour * 1), // Berlaku 1 jam
	}

	// Hapus token lama jika ada agar tidak menumpuk
	_ = s.repo.DeletePasswordReset(ctx, req.Email)

	err = s.repo.CreatePasswordReset(ctx, pwReset)
	if err != nil {
		return err
	}

	// 4. Kirim email reset password
	if s.notifService != nil {
		s.notifService.SendOTPResetEmail(user.Email, user.Name, resetToken)
	}
	
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	// 1. Cari token di database
	pwReset, err := s.repo.FindPasswordResetByToken(ctx, req.Token)
	if err != nil {
		return err
	}
	if pwReset == nil {
		return errors.New("invalid or expired token")
	}

	// 2. Cek kedaluwarsa
	if time.Now().After(pwReset.ExpiresAt) {
		_ = s.repo.DeletePasswordReset(ctx, pwReset.Email) // Bersihkan token
		return errors.New("invalid or expired token")
	}

	// 3. Hash password baru
	hashedPassword, err := hash.MakeHash(req.NewPassword)
	if err != nil {
		return err
	}

	// 4. Update password di tabel users
	err = s.repo.UpdatePassword(ctx, pwReset.Email, hashedPassword)
	if err != nil {
		return err
	}

	// 5. Hapus token agar tidak bisa dipakai 2x
	_ = s.repo.DeletePasswordReset(ctx, pwReset.Email)

	return nil
}
