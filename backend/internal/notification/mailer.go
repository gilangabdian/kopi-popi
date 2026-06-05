package notification

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// mailTemplate adalah struktur dasar HTML email yang seragam (dengan logo/styling)
const baseHTMLTemplate = `
<!DOCTYPE html>
<html>
<head>
<style>
	body { font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; }
	.container { background-color: #ffffff; padding: 30px; border-radius: 8px; max-width: 600px; margin: 0 auto; box-shadow: 0 4px 8px rgba(0,0,0,0.1); }
	.header { text-align: center; border-bottom: 1px solid #ddd; padding-bottom: 20px; margin-bottom: 20px; }
	.header h1 { color: #d35400; margin: 0; }
	.content { color: #333; line-height: 1.6; }
	.footer { text-align: center; font-size: 12px; color: #888; margin-top: 30px; border-top: 1px solid #ddd; padding-top: 10px; }
	.highlight { background-color: #fcf3cf; padding: 10px; border-left: 4px solid #f1c40f; margin: 20px 0; }
	.btn { display: inline-block; padding: 10px 20px; background-color: #d35400; color: #fff; text-decoration: none; border-radius: 5px; margin-top: 20px; }
</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>☕ Kopi-Popi</h1>
		</div>
		<div class="content">
			{{.Body}}
		</div>
		<div class="footer">
			<p>&copy; {{.Year}} Kopi-Popi System. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`

// Fungsi untuk membungkus konten HTML ke dalam Base Template
func wrapWithBaseTemplate(bodyHTML string) (string, error) {
	tmpl, err := template.New("email").Parse(baseHTMLTemplate)
	if err != nil {
		return "", err
	}

	data := struct {
		Body template.HTML
		Year int
	}{
		Body: template.HTML(bodyHTML),
		Year: time.Now().Year(),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// === Kumpulan Template Spesifik ===

func getInvoiceTemplate(customerName string, transactionID string, amount float64) (string, error) {
	body := `
		<h2>Terima Kasih Atas Pesanan Anda, ` + customerName + `!</h2>
		<p>Pesanan Anda telah berhasil diproses. Berikut adalah rincian pembayaran Anda:</p>
		<div class="highlight">
			<strong>Nomor Transaksi:</strong> ` + transactionID + `<br>
			<strong>Total Pembayaran:</strong> Rp ` + formatCurrency(amount) + `
		</div>
		<p>Tim barista kami sedang menyiapkan pesanan Anda dengan penuh cinta. Kami akan mengabari Anda kembali jika pesanan sudah siap diambil!</p>
	`
	return wrapWithBaseTemplate(body)
}

func getOrderReadyTemplate(customerName string, transactionID string) (string, error) {
	body := `
		<h2>Pesanan Anda Sudah Siap! 🎉</h2>
		<p>Halo ` + customerName + `! Kopi favorit Anda untuk transaksi <strong>#` + transactionID + `</strong> sudah selesai diracik dan siap untuk diambil atau dinikmati.</p>
		<p>Segera kunjungi konter pengambilan. Sampai jumpa!</p>
	`
	return wrapWithBaseTemplate(body)
}

func getRestockRequestTemplate(branchName string, reqID string) (string, error) {
	body := `
		<h2>Permintaan Restock Baru</h2>
		<p>Admin, terdapat permintaan pengisian ulang stok (Restock) dari cabang <strong>` + branchName + `</strong>.</p>
		<div class="highlight">
			<strong>ID Request:</strong> ` + reqID + `
		</div>
		<p>Silakan login ke Dashboard Kopi-Popi untuk menyetujui atau menolak permintaan ini.</p>
	`
	return wrapWithBaseTemplate(body)
}

func getRestockResultTemplate(branchName string, reqID string, status string) (string, error) {
	color := "green"
	if status == "Rejected" {
		color = "red"
	}
	body := `
		<h2>Update Permintaan Restock</h2>
		<p>Permintaan restock untuk cabang <strong>` + branchName + `</strong> (ID: ` + reqID + `) telah diperbarui.</p>
		<div class="highlight" style="border-left-color: ` + color + `">
			Status Saat Ini: <strong style="color:` + color + `">` + status + `</strong>
		</div>
	`
	return wrapWithBaseTemplate(body)
}

func getWelcomeTemplate(userName string) (string, error) {
	body := `
		<h2>Selamat Datang di Kopi-Popi, ` + userName + `! ☕</h2>
		<p>Terima kasih telah bergabung dengan sistem Kopi-Popi. Akun Anda telah berhasil didaftarkan.</p>
		<p>Nikmati pengalaman mengelola dan memesan kopi terbaik dengan sistem terpadu kami.</p>
		<a href="#" class="btn">Mulai Jelajahi</a>
	`
	return wrapWithBaseTemplate(body)
}

func getRegistrationOTPTemplate(userName, email, otpCode string) string {
	verifyLink := fmt.Sprintf("http://localhost:3000/verify?email=%s", email)
	return fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px; background-color: #f4f4f4;">
			<div style="max-width: 600px; margin: auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 4px 8px rgba(0,0,0,0.1);">
				<h2 style="color: #FF5722; text-align: center;">Verifikasi Akun Kopi-Popi Anda ☕</h2>
				<p>Halo <strong>%s</strong>,</p>
				<p>Terima kasih telah mendaftar di Kopi-Popi. Untuk mengaktifkan akun Anda, masukkan 6 digit kode OTP di bawah ini:</p>
				<div style="text-align: center; margin: 30px 0;">
					<span style="font-size: 32px; font-weight: bold; letter-spacing: 5px; color: #333; background: #eee; padding: 10px 20px; border-radius: 5px;">%s</span>
				</div>
				<p style="text-align: center; color: #e53935; font-size: 14px;">Kode OTP ini hanya berlaku selama 5 menit.</p>
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="background-color: #FF5722; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; font-weight: bold;">Verifikasi Akun Saya Sekarang</a>
				</div>
				<p>Jika Anda tidak merasa mendaftar di Kopi-Popi, abaikan saja email ini.</p>
				<hr style="border: 1px solid #ddd; margin-top: 30px;">
				<p style="font-size: 12px; color: #888; text-align: center;">© 2026 Kopi-Popi. All rights reserved.</p>
			</div>
		</div>
	`, userName, otpCode, verifyLink)
}

func getOTPResetTemplate(userName string, otpCode string) (string, error) {
	body := `
		<h2>Permintaan Reset Password</h2>
		<p>Halo ` + userName + `! Kami menerima permintaan untuk mengatur ulang kata sandi akun Anda.</p>
		<p>Berikut adalah kode OTP Anda. Kode ini hanya berlaku selama 15 menit.</p>
		<div class="highlight" style="text-align: center; font-size: 24px; letter-spacing: 5px; font-weight: bold;">
			` + otpCode + `
		</div>
		<p>Jika Anda tidak merasa melakukan permintaan ini, abaikan email ini dan pastikan akun Anda aman.</p>
	`
	return wrapWithBaseTemplate(body)
}

// Helper (simple formatter, in production you'd use a real formatting lib)
func formatCurrency(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
