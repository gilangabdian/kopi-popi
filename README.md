# Kopi-Popi POS & Inventory System

Sistem Point of Sale (POS) dan Inventory Multi-Outlet untuk Kopi-Popi.

## Perancangan & Dokumentasi Sistem
Untuk memudahkan pemahaman alur sistem secara menyeluruh, dokumentasi rancangan sistem dipecah ke dalam beberapa file khusus yang berada di folder `backend/docs/`.
Silakan klik tautan di bawah ini untuk melihat detailnya:

1. 👤 [Use Cases & Aktor](backend/docs/use-cases.md) - Penjelasan fungsionalitas tiap pengguna (Customer, Cashier, Manager, Admin).
2. 🔄 [Activity Diagrams](backend/docs/activity-diagrams.md) - Alur kerja sistem (Autentikasi, Transaksi, Shift Kasir, Restock, dll).
3. 🗄️ [Skema Database](backend/docs/database.md) - Penjelasan ringkas mengenai tabel-tabel dan relasinya.
4. 🏗️ [Arsitektur & DDD](backend/docs/architecture.md) - Penjelasan struktur *Monorepo* dan implementasi *Domain-Driven Design*.

## Fitur Utama
1. **Users & Auth**: Role management (Admin, Manager, Cashier, Customer), JWT, verifikasi OTP.
2. **Branches**: Manajemen multi-cabang.
3. **Catalog**: Kategori, Material (Bahan Baku), dan Produk (beserta resep/BOM).
4. **Inventory**: Mutasi stok cabang, Restock Request (Manager -> Admin), Kedatangan stok gudang pusat, Alokasi pusat ke cabang.
5. **Sales & Transactions**:
   - Shift kasir (Buka/Tutup Kasir)
   - Manajemen Keranjang (Customer Online & Hold Bill Offline)
   - Checkout & Integrasi Pembayaran (Midtrans)
   - Pemrosesan Transaksi / Kitchen Display System (GET & UPDATE status pesanan)
6. **Dashboard & Analytics**: Laporan pendapatan penjualan (termasuk breakdown metode pembayaran), produk terlaris, dan histori shift kasir.
7. **Media & Uploads**: Endpoint terpusat (`POST /uploads`) untuk menyimpan gambar statis (misalnya cover artikel atau gambar produk).

## API Documentation
Dokumentasi API tersedia dalam format OpenAPI 3.0.
- Source files berada di folder `backend/api/paths/` dan `backend/api/schemas/`.
- File OpenAPI yang sudah dibundle (siap di-import ke Postman/Redocly) berada di: `backend/api/openapi-bundled.json`.

### Cara Update OpenAPI Bundled
Jika melakukan perubahan pada endpoint di folder `paths` atau `schemas`, lakukan bundling ulang menggunakan Redocly:
```bash
cd backend
npx @redocly/cli bundle api/openapi.json -o api/openapi-bundled.json
```

## Setup Lokal (Docker)
Jalankan perintah berikut di root folder untuk membangun database dan backend:
```bash
docker compose up -d --build
```
Aplikasi backend akan berjalan di `http://localhost:8080`.

## Frontend Web Application
Bagian frontend (klien) dibangun menggunakan **Next.js** (App Router) untuk memenuhi kebutuhan Landing Page (Customer) dan Dashboard/POS (Admin/Kasir).

### Teknologi yang Digunakan
- **Next.js 16** (React Framework)
- **Tailwind CSS** (Styling)
- **Embla Carousel** (Slider)
- **Iconify** (Icons)

### Cara Menjalankan Frontend
Pastikan Anda sudah menginstal Node.js versi terbaru, lalu jalankan perintah berikut dari root proyek:
```bash
cd frontend
npm install
npm run dev
```
Aplikasi web (Landing Page) akan berjalan di `http://localhost:3000`.
