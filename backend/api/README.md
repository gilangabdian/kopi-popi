# API Documentation - Multi Outlet POS Kopi-Popi & Inventory

Dokumentasi API ini ditulis menggunakan standar **OpenAPI 3.0.3**. Untuk menjaga agar dokumen tetap mudah dibaca, dikelola (_maintainable_), dan *scalable*, dokumentasi ini dipisah berdasarkan konsep **Domain-Driven Design**.

---

## 📂 Struktur Direktori Spesifikasi

Seluruh file spesifikasi JSON dibagi menjadi komponen-komponen berikut:

| Folder / File | Deskripsi | Contoh Isi / Domain |
| :--- | :--- | :--- |
| **`openapi.json`** | File **Root / Utama**. Tempat mendeklarasikan Info API, Server, Security Global, dan mapping Path. | `"title": "Multi Outlet POS"`, mapping `$ref` ke folder `paths/`. |
| **`paths/`** | Berisi definisi semua endpoint (Method, Params, Request Body). Dipisah per fitur/domain. | `auth.json`, `users.json`, `branches.json`, `catalogues.json`, `media.json` |
| **`schemas/`** | Berisi definisi **Model/Struct** JSON agar *reusable* dan tidak perlu ditulis berulang-ulang. | `User`, `Product`, `Branch`, `ErrorResponse` |
| **`responses/`** | Berisi definisi **Format Response Standard** aplikasi. | `GenericSuccess`, `GenericError` |

---

## ✅ Modul / Domain yang Telah Diimplementasikan

Hingga saat ini, sistem *backend* telah memiliki implementasi penuh untuk domain berikut:

1. **Auth (`internal/auth`)**: Menangani otentikasi (Register, Login, Forgot Password).
2. **User (`internal/user`)**: Manajemen pengguna, staf, hak akses, manajemen profil pribadi, serta **Sistem Membership & Loyalty Points**.
3. **Branch (`internal/branch`)**: CRUD untuk manajemen cabang kedai kopi.
4. **Catalog (`internal/catalog`)**: Manajemen Kategori, Material (Bahan Baku), dan Produk (termasuk Resep / Bill of Materials). Dilengkapi proteksi rahasia resep.
5. **Media (`internal/media`)**: Sentralisasi unggah gambar/file dengan proteksi ekstensi (JPG, PNG) dan batas ukuran (5MB). Mendukung folder dinamis.
6. **Inventory (`internal/inventory`)**: Manajemen stok fisik cabang, buku riwayat mutasi (kartu stok), dan alur logistik permintaan barang antar cabang (Restock Requests) yang sudah dilengkapi dengan fitur persetujuan dan pemberian alasan penolakan (*Rejection Reason*) oleh Admin.
7. **Sales (`internal/sales`)**: Inti dari sistem POS (Point of Sales). Mengurus *Shift* Kasir (Buka/Tutup Kasir dengan rekonsiliasi kas beserta pencatatan *Petty Cash* / Pengeluaran Mendadak), *Cart* Online & Offline, proses *Checkout*, dan **Kitchen Display System (KDS)** untuk pemantauan antrean dan pembaruan status pesanan (Paid -> Preparing -> Ready -> Completed), terintegrasi penuh dengan pemotongan & penambahan **Loyalty Points**.
8. **Notifications (`internal/notification`)**: Menangani *In-App Notifications* (Notifikasi web) dan *Email Notifications* Asynchronous via SMTP. Modul ini menjadi fondasi utama untuk pengiriman e-Receipt (Invoice), peringatan stok minimum, dan persetujuan Restock.
9. **Payment (`internal/payment`)**: Integrasi Gateway Pembayaran menggunakan **Midtrans**. Modul ini memisahkan logika khusus pembayaran seperti Request Snap Token (URL Pembayaran) dan Webhook Handler untuk memproses callback dari Midtrans yang secara otomatis memotong stok dan mengubah status Transaksi.
10. **Analytics (`internal/analytics`)**: Modul pelaporan dan *dashboard*. Berfungsi melakukan agregasi data untuk menghasilkan laporan pendapatan/penjualan (beserta rincian per metode pembayaran), produk paling laris (*top products*), serta rekap selisih uang kasir (*shifts report*).
11. **Promo (`internal/promo`)**: Modul pemasaran (*marketing*) untuk pengelolaan kupon dan kode diskon tingkat transaksi (*cart-level*). Mendukung potongan harga berbentuk persentase (dilengkapi batas maksimum diskon) dan nominal tetap (*fixed*).
12. **Blogs (`internal/blogs`)**: Modul konten manajemen untuk membuat dan mengelola artikel/blog, lengkap dengan kalkulasi estimasi waktu baca otomatis.

---

## 🔐 Authentication Strategy

Sistem menggunakan global authentication berbasis **JWT (JSON Web Token)**. Secara *default*, semua endpoint membutuhkan autentikasi kecuali jika di-override dengan `"security": []`.

**Format Header HTTP yang harus dikirim:**
```http
Authorization: Bearer <your_jwt_token_here>
```

---

## 📸 Panduan Upload Gambar (Media Module)

Karena Golang didesain bebas (*unopinionated*), Kopi-Popi menggunakan modul **Media** terpusat untuk segala kebutuhan upload gambar (seperti Profil atau Foto Produk).

**Alur Upload Gambar:**
1. *Frontend* mengirim file fisik gambar (`multipart/form-data`) ke endpoint `POST /uploads`. Parameter `folder` bisa diisi `products`, `profiles`, atau `blogs`.
2. *Backend* akan memvalidasi file (maksimal 5MB, format JPG/PNG/WEBP).
3. Jika lolos, file disimpan secara lokal di dalam container Docker di folder `/app/uploads` (di-mapping ke Host `./uploads` lewat Docker Volume agar tidak hilang).
4. *Backend* mengembalikan URL publik (misal: `http://localhost:8080/uploads/blogs/1234.jpg`).
5. *Frontend* menyimpan URL tersebut, lalu memanggil `POST /products`, `PATCH /users/me`, atau `POST /blogs` dengan mengirim data JSON berisi `image_url`, `profile_picture`, atau menyisipkannya di `content` blog menggunakan URL tersebut.

---


Endpoint berikut **WAJIB** menyertakan token JWT yang valid.

| Fitur / Domain | Contoh Endpoint Utama | Akses Ideal |
| :--- | :--- | :--- |
| **User Profile** | `/users/me`, `/users/me/password` | *Semua Role* |
| **Users Mgmt** | `/users`, `/users/managers` | *Admin & Manager* |
| **Branches** | `/branches` (POST, PUT, DELETE) | *Admin* |
| **Catalog (Material)**| `/materials` (CRUD) | *Admin & Manager* |
| **Catalog (Product)**| `/products` (POST, PUT, DELETE) | *Admin* |
| **Media / Uploads** | `/uploads` (`POST`) | *Semua Pengguna Terdaftar* |
| **Inventory**| `/inventories/branches/:id`| *Admin & Manager* |
| **Inventory**| `/inventories/restocks` | *Admin & Manager* |
| **Analytics (Reports)**| `/reports/sales`, `/reports/top-products` | *Admin & Manager* |

---

## 🛠️ Cara Menampilkan Secara Visual (Swagger UI)

Karena spesifikasi JSON dipisah ke beberapa file, cara terbaik untuk membacanya selama tahap _development_:

1. **Menggunakan VS Code (Lokal)**
   - Install ekstensi **OpenAPI (Swagger) Editor** (oleh 42Crunch).
   - Buka file `openapi.json`.
   - Tekan `Shift + Alt + O` untuk menampilkan UI interaktif.

2. **Bundle Menjadi 1 File Utuh (Untuk di Deploy / Redocly)**
   Gunakan Redocly CLI untuk membundel spesifikasi modular menjadi satu file:
   ```bash
   npx @redocly/cli bundle openapi.json -o openapi-bundled.json
   ```
   *File `openapi-bundled.json` ini sudah tersedia di dalam folder ini dan bisa diimpor ke Postman atau Swagger UI.*
