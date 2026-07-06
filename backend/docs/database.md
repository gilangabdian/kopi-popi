# Skema Database & Relasi

Dokumen ini menjabarkan struktur pangkalan data (*database*) dari sistem Kopi Popi, yang mencakup manajemen hak akses, inventaris, hingga point-of-sale (POS).

**🔗 Tautan Skema Lengkap (DBML):** [Kopi Popi dbdiagram.io](https://dbdiagram.io/d/kopi-popi-6a0c3cdf9f1f8ec47b4f0cb6)

---

## Ringkasan Eksekutif
Secara keseluruhan, database ini menggunakan **17 tabel relasional** yang bisa dikelompokkan ke dalam 7 kategori fungsional:
1. **Users & Access Management** (3 tabel)
2. **Branches & Operational** (3 tabel)
3. **Master Data & Promos** (5 tabel)
4. **Inventory & Logistics** (6 tabel)
5. **Cart & Transactions** (4 tabel)
6. **Notifications** (1 tabel)
7. **Blogs** (1 tabel)

---

## Penjelasan Kategori dan Tabel

### 1. Users & Access Management
Berfungsi mengatur siapa saja yang bisa masuk ke sistem beserta data profilnya.
- **`roles`**: Berisi definisi peran (Admin, Manager, Cashier, Customer).
- **`users`**: Tabel sentral tempat semua data profil pengguna disimpan (UUID, Role ID, Branch ID untuk pekerja toko, Email, Password, Poin Loyalitas).
- **`email_verifications` & `password_resets`**: Tabel penunjang untuk keamanan (kode OTP dan token *reset* sandi).

### 2. Branches & Operational
Berfungsi menaungi data cabang fisik toko dan operasional pergantian *shift* kasir.
- **`branches`**: Menyimpan data identitas tiap cabang, termasuk Gudang Pusat (HQ) yang umumnya memiliki ID=1.
- **`shifts`**: Mencatat sesi kerja kasir, mulai dari modal kas (*starting cash*) saat membuka toko, hingga pendapatan bersih aktual saat menutup toko.
- **`expenses`**: Mencatat jika ada pengeluaran insidental selama *shift* (misalnya: beli galon air).

### 3. Master Data & Promos
Berfungsi sebagai pusat katalog menu dan promo perusahaan.
- **`categories`**: Kategori pengelompokan menu/bahan (Kopi, Non-Kopi, Susu, dll).
- **`products`**: Tabel menu barang jadi yang siap dijual ke *customer* (misalnya: Es Kopi Susu).
- **`materials`**: Tabel bahan mentah penyusun produk (misalnya: Susu UHT, Kopi Robusta).
- **`product_boms` (Bill of Materials)**: Tabel relasional yang menghubungkan 1 Produk (misal Kopi Susu) dengan banyak Material (misal 50ml espresso, 100ml susu). Di sinilah kunci otomatisasi pengurangan stok terjadi.
- **`promos`**: Tabel penyimpan kode promo, tipe diskon (persen/nominal), dan tanggal berlaku.

### 4. Inventory & Logistics
Berfungsi melacak persediaan fisik di tiap cabang dan mengatur aliran distribusi dari Pusat.
- **`branch_inventories`**: Menyimpan data aktual jumlah/stok bahan mentah (Material) di suatu Cabang tertentu.
- **`inventory_movements`**: Mencatat riwayat historis mutasi barang (Masuk, Keluar, Penyesuaian/Opname).
- **`restock_requests` & `restock_items`**: Menangani permintaan suplai barang dari Manager Cabang ke Admin Pusat.
- **`incoming_stocks` & `incoming_stock_items`**: Menangani rekap penerimaan pasokan barang baru dari pihak ketiga/Supplier ke Gudang Pusat.

### 5. Cart & Transactions
Berfungsi menangani siklus pembelanjaan dari keranjang hingga pembayaran lunas.
- **`carts` & `cart_items`**: Tempat penampungan sementara sebelum *Customer* melakukan *Checkout*.
- **`transactions`**: Tabel krusial yang merekam nota pembelian. Menyimpan ID Cabang, ID Customer (jika online), ID Kasir & Shift (jika offline), tipe pesanan (Takeaway/Dine-in), total bayar, serta status transaksi (Menunggu Pembayaran -> Sedang Diproses -> Selesai).
- **`transaction_details`**: Merinci Produk Jadi (menu) apa saja yang dibeli dalam satu struk Transaksi beserta *subtotal*-nya.

### 6. Notifications
- **`notifications`**: Tabel khusus untuk mencatat pesan/pemberitahuan sistem ke UUID *User* spesifik (baik Customer maupun Admin).

### 7. Blogs
- **`blogs`**: Menampung konten/artikel berita yang dibuat oleh Admin/Manager, menghubungkan penulisnya (`author_id`) ke tabel `users`.

---

## Hubungan Antar Data (Relasi Kunci)
* **User -> Role & Branch**: Setiap `user` memiliki satu `role`. Apabila dia seorang Kasir atau Manager, dia akan dipetakan secara statis ke satu `branch`.
* **Branch -> Shift & Transaction**: Sebuah `branch` menaungi banyak `shift` dan mencetak banyak `transaction`.
* **Product -> BOM -> Material**: Untuk membuat 1 `product`, sistem membaca tabel `product_boms` untuk melihat `material` apa saja yang perlu ditarik dari `branch_inventories` secara atomik saat `transaction` berstatus lunas.
* **Request -> Material**: `restock_requests` terikat pada Cabang pemohon, berisi banyak `restock_items` yang masing-masing merujuk ke ID `material`. Apabila disetujui, nilainya dimutasikan ke `branch_inventories`.
