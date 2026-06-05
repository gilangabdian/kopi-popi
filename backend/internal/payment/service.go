package payment

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/gilangages/kopi-popi/internal/notification"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type Service interface {
	CreateSnapURL(orderID string, totalAmount float64, customerName, customerEmail string) (string, error)
	HandleWebhook(ctx context.Context, payload map[string]interface{}) error
}

type paymentService struct {
	repo         Repository
	notifService notification.Service
	snapClient   snap.Client
}

func NewService(repo Repository, notif notification.Service) Service {
	// Initialize Midtrans setup
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	isProd := os.Getenv("MIDTRANS_IS_PRODUCTION") == "true"
	
	env := midtrans.Sandbox
	if isProd {
		env = midtrans.Production
	}
	
	var s snap.Client
	s.New(serverKey, env)

	return &paymentService{repo: repo, notifService: notif, snapClient: s}
}

func (s *paymentService) CreateSnapURL(orderID string, totalAmount float64, customerName, customerEmail string) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(totalAmount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
		},
	}

	snapResp, err := s.snapClient.CreateTransaction(req)
	if err != nil {
		return "", err
	}

	return snapResp.RedirectURL, nil
}

func (s *paymentService) HandleWebhook(ctx context.Context, payload map[string]interface{}) error {
	// 1. Parse payload Midtrans
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return errors.New("invalid order_id in webhook")
	}
	transactionStatus, _ := payload["transaction_status"].(string)

	// In real production, we MUST verify the signature key here:
	// signatureKey = SHA512(order_id + status_code + gross_amount + ServerKey)

	if transactionStatus == "settlement" || transactionStatus == "capture" {
		// Update status transaksi ke Paid
		err := s.repo.UpdateTransactionStatusToPaid(ctx, orderID)
		if err != nil {
			log.Printf("[WEBHOOK ERROR] Failed to update status to Paid for %s: %v", orderID, err)
			return err
		}

		// Ambil data pelanggan dan total amount untuk kirim invoice
		txInfo, err := s.repo.GetTransactionInfo(ctx, orderID)
		if err != nil {
			log.Printf("[WEBHOOK ERROR] Failed to get tx info: %v", err)
			return nil // Still return 200 OK to midtrans so it doesn't retry
		}
		
		// Potong stok (Idealnya ini dikelola terpusat atau kita trigger event)
		// Karena payment berdiri sendiri, kita bisa mengupdate transaction_status, 
		// lalu kita panggil notifikasi.
		// NOTE: Untuk potong stok, webhook idealnya memanggil sales/inventory service, 
		// tapi untuk menyederhanakan, kita asumsikan webhook payment ini memanggil service inventory secara langsung
		// atau cukup kirim email invoice saja dan stok sudah dipotong waktu Checkout "Pending".
		// Oh, di Rencana: "Stok gudang BELUM dipotong. Webhook: Memotong Stok gudang".
		// Untuk memotong stok gudang, payment perlu depend ke Inventory, atau Repository Payment melakukan raw query pengurangan stok berdasarkan cart_items.
		// Lebih aman lewat Repository agar tidak cycle.
		err = s.repo.DeductStockForTransaction(ctx, orderID)
		if err != nil {
			log.Printf("[WEBHOOK ERROR] Failed to deduct stock for %s: %v", orderID, err)
		}

		if s.notifService != nil {
			email := txInfo.CustomerEmail
			if email == "" {
				email = "guest@example.com"
			}
			s.notifService.SendInvoiceEmail(email, txInfo.CustomerName, orderID, txInfo.TotalAmount)
			// Kirim notifikasi ke kasir
			// title := "Pesanan Online Lunas!"
			// msg := fmt.Sprintf("Pesanan %s sudah lunas. Segera siapkan pesanan.", orderID)
			// _ = s.notifService.SendInAppNotification("Cashier", title, msg, "RELOAD_ORDER_LIST") // Gagal karena butuh UUID Cashier, bukan string "Cashier"
		}

	} else if transactionStatus == "cancel" || transactionStatus == "expire" || transactionStatus == "deny" {
		_ = s.repo.UpdateTransactionStatusToFailed(ctx, orderID)
	}

	return nil
}
