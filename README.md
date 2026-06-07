# Kopi-Popi POS & Inventory System

Sistem Point of Sale (POS) dan Inventory Multi-Outlet untuk Kopi-Popi.

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
