package notification

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/google/uuid"
)

type Service interface {
	// In-App
	GetMyNotifications(ctx context.Context, userID string) ([]Notification, error)
	MarkAsRead(ctx context.Context, userID string, id string) error

	// Auth Triggers
	SendRegistrationOTPEmail(email string, userName string, otpCode string)
	SendOTPResetEmail(email string, userName string, otpCode string)
	SendInvoiceEmail(customerEmail string, customerName string, transactionID string, amount float64)
	SendOrderReadyEmail(customerEmail string, customerName string, transactionID string)
	SendRestockRequestEmail(adminEmail string, branchName string, reqID string)
	SendRestockResultEmail(managerEmail string, branchName string, reqID string, status string)

	// In-App (Sync / Internal Use)
	SendInAppNotification(userID string, title string, message string, notifType string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// --- IN-APP NOTIFICATIONS ---

func (s *service) GetMyNotifications(ctx context.Context, userID string) ([]Notification, error) {
	return s.repo.GetNotificationsByUser(userID)
}

func (s *service) MarkAsRead(ctx context.Context, userID string, id string) error {
	return s.repo.MarkAsRead(id, userID)
}

func (s *service) SendInAppNotification(userID string, title string, message string, notifType string) error {
	notif := &Notification{
		ID:      uuid.NewString(),
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    notifType,
		IsRead:  false,
	}
	return s.repo.SaveNotification(notif)
}

// --- EMAIL NOTIFICATIONS (ASYNC) ---

func (s *service) sendEmailAsync(to string, subject string, bodyHTML string) {
	go func() {
		smtpHost := os.Getenv("SMTP_HOST")
		smtpPort := os.Getenv("SMTP_PORT")
		smtpUser := os.Getenv("SMTP_USER")
		smtpPass := os.Getenv("SMTP_PASSWORD")
		senderName := os.Getenv("SMTP_SENDER_NAME")

		// Jika SMTP belum di-setup, log ke terminal saja (Mock Email)
		if smtpUser == "" || smtpPass == "" {
			log.Printf("[MOCK EMAIL] To: %s | Subject: %s\n", to, subject)
			log.Println("--- Email Body Start ---")
			log.Println(bodyHTML)
			log.Println("--- Email Body End ---")
			return
		}

		// Kirim email beneran
		from := senderName + " <" + smtpUser + ">"
		
		// Setup headers
		headers := make(map[string]string)
		headers["From"] = from
		headers["To"] = to
		headers["Subject"] = subject
		headers["MIME-Version"] = "1.0"
		headers["Content-Type"] = `text/html; charset="utf-8"`

		message := ""
		for k, v := range headers {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + bodyHTML

		auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, []byte(message))
		if err != nil {
			log.Printf("[EMAIL ERROR] Failed to send email to %s: %v\n", to, err)
		} else {
			log.Printf("[EMAIL SUCCESS] Sent to %s | Subject: %s\n", to, subject)
		}
	}()
}

// --- SPECIFIC EMAIL HANDLERS ---

func (s *service) SendRegistrationOTPEmail(email string, name string, otpCode string) {
	subject := "Verifikasi Akun Kopi-Popi Anda"
	body := getRegistrationOTPTemplate(name, email, otpCode)
	s.sendEmailAsync(email, subject, body)
}

func (s *service) SendOTPResetEmail(email string, name string, otpCode string) {
	html, err := getOTPResetTemplate(name, otpCode)
	if err != nil {
		log.Printf("[EMAIL TEMPLATE ERROR] %v\n", err)
		return
	}
	s.sendEmailAsync(email, "Kode OTP Reset Password", html)
}

func (s *service) SendInvoiceEmail(customerEmail string, customerName string, transactionID string, amount float64) {
	html, err := getInvoiceTemplate(customerName, transactionID, amount)
	if err != nil {
		log.Printf("[EMAIL TEMPLATE ERROR] %v\n", err)
		return
	}
	s.sendEmailAsync(customerEmail, "Invoice Pembelian Anda", html)
}

func (s *service) SendOrderReadyEmail(customerEmail string, customerName string, transactionID string) {
	html, err := getOrderReadyTemplate(customerName, transactionID)
	if err != nil {
		log.Printf("[EMAIL TEMPLATE ERROR] %v\n", err)
		return
	}
	s.sendEmailAsync(customerEmail, "Pesanan Anda Sudah Siap!", html)
}

func (s *service) SendRestockRequestEmail(adminEmail string, branchName string, reqID string) {
	html, err := getRestockRequestTemplate(branchName, reqID)
	if err != nil {
		log.Printf("[EMAIL TEMPLATE ERROR] %v\n", err)
		return
	}
	s.sendEmailAsync(adminEmail, "Permintaan Restock Baru - "+branchName, html)
}

func (s *service) SendRestockResultEmail(managerEmail string, branchName string, reqID string, status string) {
	html, err := getRestockResultTemplate(branchName, reqID, status)
	if err != nil {
		log.Printf("[EMAIL TEMPLATE ERROR] %v\n", err)
		return
	}
	s.sendEmailAsync(managerEmail, "Update Restock: "+status, html)
}
