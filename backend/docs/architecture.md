# Arsitektur & Struktur Direktori

Sistem Kopi Popi dibangun dengan menggunakan pendekatan **Monorepo** yang memisahkan ranah Frontend dan Backend, serta menerapkan filosofi **Domain-Driven Design (DDD)** pada layer kapabilitas layanannya.

---

## 1. Konsep Monorepo
Seluruh basis kode proyek ini ditempatkan dalam satu repositori (*repository*) tunggal (Monorepo), namun secara logis terpisah menjadi dua layanan utama (Microservices/Modules):

- **`/frontend`**: Aplikasi *client-facing* yang dibangun menggunakan **Next.js** (React) dan Tailwind CSS. Bagian ini menangani seluruh *User Interface* (UI) baik untuk Customer (Landing Page) maupun Panel POS/Admin.
- **`/backend`**: Aplikasi API yang dibangun menggunakan **Golang** (*Asumsi berdasarkan praktik umum DDD di backend modern*). Bagian ini beroperasi tanpa tampilan antarmuka dan bertugas khusus mengolah *business logic*, memproses keamanan otentikasi, serta berinteraksi dengan *Database*.

Pemisahan ini memungkinkan skalabilitas mandiri (misal backend *down* tidak mematikan *landing page* statis, atau frontend bisa di-*deploy* ke Vercel sementara backend ke AWS/GCP).

---

## 2. Implementasi Domain-Driven Design (DDD) di Backend
**Domain-Driven Design (DDD)** adalah pola perancangan arsitektur perangkat lunak yang berfokus pada "Domain Bisnis Utama" dari aplikasi tersebut. Daripada mengelompokkan *file* berdasarkan sifat teknisnya (misal: memisahkan semua *controller* di satu folder dan semua *model* di folder lain), DDD menuntut kita mengelompokkan *file* berdasarkan **Konteks Fungsionalitasnya (Bounded Context)**.

Berdasarkan *Database Schema* yang ada, struktur folder `backend` Kopi Popi disarankan untuk ditata dengan modul domain (berada di `backend/internal/domains/`) sebagai berikut:

### 📁 Domain `users`
**Fungsi:** Menangani seluruh siklus hidup entitas manusia/aktor.
- **Tanggung jawab:** Autentikasi (Login/Register), pengaturan Role (Kasir, Manager), profil pengguna, poin loyalitas, serta verifikasi *password* dan OTP.
- **Cakupan Tabel:** `users`, `roles`, `email_verifications`, `password_resets`.

### 📁 Domain `branches`
**Fungsi:** Menangani logika pembukaan toko dan manajemen *shift* pekerja.
- **Tanggung jawab:** Menyimpan data jam buka/tutup cabang, serta merekam pembukaan dan penutupan *shift* kasir beserta audit uang kas di laci.
- **Cakupan Tabel:** `branches`, `shifts`, `expenses`.

### 📁 Domain `catalog`
**Fungsi:** Menangani master data apa yang dijual dan apa pembentuknya.
- **Tanggung jawab:** Mengelola jenis Produk Jadi (Menu), Kategori, Bahan Baku (Material), dan Promosi. Domain ini sangat penting karena menyimpan resep atau *Bill of Materials* (BOM).
- **Cakupan Tabel:** `categories`, `products`, `materials`, `product_boms`, `promos`.

### 📁 Domain `inventory`
**Fungsi:** Menangani aspek persediaan, logistik, dan aliran barang.
- **Tanggung jawab:** Mencatat stok fisik (*stock on-hand*) di setiap cabang, mencatat arus keluar/masuk (*movements*), serta menangani permohonan suplai barang dari Cabang ke Gudang Pusat (*Restock Requests*).
- **Cakupan Tabel:** `branch_inventories`, `inventory_movements`, `restock_requests`, `restock_items`, `incoming_stocks`, `incoming_stock_items`.

### 📁 Domain `transactions`
**Fungsi:** Menangani siklus *checkout* hingga kasir.
- **Tanggung jawab:** Menyimpan status sementara di Keranjang Belanja (*Cart*), merekam bukti nota pembayaran (*Transactions*), menghitung subtotal diskon, serta menginisiasi pemicu (*trigger*) otomatis ke domain `inventory` untuk memotong stok saat pesanan sukses (menggunakan data resep dari domain `catalog`).
- **Cakupan Tabel:** `carts`, `cart_items`, `transactions`, `transaction_details`.

### 📁 Domain `notifications`
**Fungsi:** Layanan notifikasi terpusat.
- **Tanggung jawab:** Mengirimkan *alert* ke *user* tertentu (seperti pesan "Pesanan Siap Diambil" atau "Request Restock Ditolak").
- **Cakupan Tabel:** `notifications`.

### 📁 Domain `blogs`
**Fungsi:** Layanan *Content Management*.
- **Tanggung jawab:** Mengatur konten artikel berita perusahaan dan estimasi waktu bacanya.
- **Cakupan Tabel:** `blogs`.

---

## Kesimpulan Arsitektur
Dengan struktur **DDD** ini, jika tim pengembang (*developer*) mendapat tugas "Tolong perbaiki fitur Checkout Kasir", ia tidak perlu mencari *file* yang berserakan. Ia cukup masuk ke direktori `internal/domains/transactions/` karena seluruh *Controller, Service/Usecase, dan Repository* terkait transaksi pembelian terkandung sepenuhnya di dalam satu modul tersebut. Ini melahirkan kode yang sangat mudah dikelola (*maintainable*) di masa depan.
