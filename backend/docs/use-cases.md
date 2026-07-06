# Use Cases

Dokumen ini menjabarkan kasus penggunaan (*use case*) dari setiap aktor yang berinteraksi dengan sistem Kopi Popi POS & Inventory.

**🔗 Tautan Diagram Figma:** [Kopi Popi Use Case Diagram](https://www.figma.com/board/cTVzMwDQakaILje5H1yMef/Untitled?node-id=1-46&t=IEWfZwEMce8ccC1m-1)

---

## 1. Customer (Pelanggan)
Aktor pelanggan yang berinteraksi langsung melalui aplikasi web *Customer-facing* (Landing Page / Web Store).

- **Melakukan Registrasi Akun**
  - *Include*: Verifikasi Kode OTP.
- **Melakukan Login**
  - *Include*: Verifikasi Kredensial.
  - *Extend*: Lupa Password / Reset Kata Sandi.
- **Kelola Profil Akun**
  - *Extend*: Edit Data Diri / Password.
- **Melihat Katalog Menu**
  - *Include*: Memilih / Mengubah Cabang Toko.
- **Kelola Keranjang Belanja**
- **Checkout Pesanan Online**
  - *Include*: Simulasi Pembayaran Dummy.
  - *Include*: Potong Stok Bahan Baku Otomatis (Sistem melihat resep pesanan via BOM dan otomatis mengurangi bahan baku di cabang tersebut).
- **Melacak Status Pesanan**
- **Melihat Pusat Notifikasi** (Melihat riwayat notifikasi pesanan/sistem).
- **Melihat Riwayat Transaksi**
  - *Extend*: Unduh Struk Digital.

---

## 2. Cashier (Kasir Toko)
Aktor kasir yang bertugas di masing-masing cabang, beroperasi menggunakan layar Tablet/POS.

- **Melakukan Login**
  - *Include*: Verifikasi Kredensial.
- **Kelola Shift Kerja**
  - *Include*: Input Modal Awal (Buka Shift).
  - *Include*: Input Uang Fisik di Laci (Tutup Shift). *(Sistem otomatis men-generate laporan)*.
- **Proses Pesanan Offline (Layar POS)**
  - *Include*: Potong Stok Bahan Baku Otomatis (via BOM).
- **Monitor Antrean Pesanan Online**
  - *Include*: Update Status Pesanan (Otomatis memicu notifikasi Email/Web ke Customer).

---

## 3. Branch Manager (Kepala Cabang)
Aktor yang bertanggung jawab penuh terhadap operasional satu lokasi cabang toko tertentu.

- **Melakukan Login**
- **Kelola Akun Kasir Cabang** (Create, Read, Update, Suspend/Nonaktifkan akun staf kasir).
- **Melihat Dashboard Analitik Cabang** (Melihat hasil laporan *generate* otomatis dari setiap Tutup Shift Kasir).
- **Kelola Stok Fisik Bahan Baku Cabang** (Melihat sisa Susu, Kopi, Gelas, dll secara aktual).
  - *Extend*: Lakukan Stock Opname (Menyesuaikan jumlah stok fisik di toko apabila ada bahan yang tumpah, rusak, atau kedaluwarsa).
- **Request Restock Bahan Baku ke Pusat** (Meminta pasokan ulang dari gudang pusat).

---

## 4. Super Admin (Head Office / Pusat)
Aktor dengan tingkat akses tertinggi yang mengatur jalannya bisnis secara global dari kantor pusat (HQ).

- **Melakukan Login**
- **Kelola Data Cabang** (Membuat entitas/lokasi cabang baru).
- **Kelola Data Akun Global**
  - *Include*: Assign Akun Branch Manager ke Cabang tertentu (via dropdown opsi cabang yang sudah ada).
- **Kelola Master Data Bahan Baku** (Membuat entitas barang mentah seperti Susu UHT, Biji Kopi, Gelas Cup).
- **Kelola Berita / Artikel Blog**
- **Kelola Master Data Produk Jadi** (Membuat entitas menu jualan seperti Es Kopi Susu).
- **Kelola Resep Produk / BOM** (Menautkan 1 Produk Jadi dengan takaran Bahan Baku spesifik).
- **Kelola Stok Gudang Pusat** (Menginput data bahan baku yang masuk dari Supplier ke gudang utama perusahaan).
- **Review Request Restock dari Cabang**
  - *Include*: Approve Request & Transfer Stok Bahan Baku Otomatis (Jika disetujui, stok Gudang Pusat berkurang dan stok Cabang bertambah).
  - *Extend*: Tolak Request & Input Alasan (Misalnya stok pusat habis, atau jumlah permintaan dinilai berlebihan).
- **Melihat Dashboard Analitik Global** (Melakukan pantauan bisnis pada semua cabang).
