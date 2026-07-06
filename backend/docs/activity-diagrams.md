# Activity Diagrams

Dokumen ini menjabarkan alur aktivitas (*activity diagram*) utama yang ada di sistem Kopi Popi POS & Inventory.

**🔗 Tautan Diagram Draw.io:** [Kopi Popi Activity Diagram](https://drive.google.com/file/d/1vKz73sIDmGUCSYnmQzVXrhzpqJ0mcA-c/view)

---

## 1. Manajemen Akun & Akses (General)
Alur ketika pengguna mengakses sistem. Memiliki jalur khusus berdasarkan niat awal pengguna dan sistem *Role*.
- **Aktor:** Pengguna (Semua Role), Sistem Backend
- **Alur Utama:**
  1. Pengguna membuka halaman Autentikasi.
  2. Pengguna memilih jalur (Registrasi, Lupa Password, atau Login).
  3. Saat login, sistem memverifikasi kredensial. Jika valid, *token session* dibuat.
  4. Sistem mengecek *Role*:
     - Jika **Customer**: Dialihkan ke beranda/katalog. Bisa lanjut mengedit profil.
     - Jika **Selain Customer**: Dialihkan ke *Dashboard* masing-masing (Kasir, Manager, Admin).

## 2. Transaksi Pelanggan "Click & Collect" & Notifikasi
Alur pemesanan online oleh pelanggan hingga pesanan diambil.
- **Aktor:** Customer, Sistem Backend, Cashier
- **Alur Utama:**
  1. Customer memilih Cabang, mengisi Keranjang, dan melakukan Checkout (Simulasi Dummy).
  2. Sistem memvalidasi pembayaran dan **secara paralel** melakukan tiga hal:
     - Memotong stok otomatis via BOM.
     - Menembak data ke layar "Monitor Antrean" kasir.
     - Mengirim notifikasi *order* berhasil ke Customer.
  3. Kasir menerima pesanan di layar, memprosesnya, dan menekan "Selesai" saat kopi sudah siap.
  4. Sistem memberi notifikasi "Pesanan Siap Diambil" ke Customer.
  5. Customer datang mengambil pesanan (opsional: unduh struk digital).

## 3. Operasional Shift Toko (POS & Audit)
Alur operasional harian kasir dari buka hingga tutup toko.
- **Aktor:** Branch Manager, Cashier, Sistem Backend
- **Alur Utama:**
  1. Manager membuatkan akun untuk Kasir.
  2. Kasir login ke Tablet POS dan menginput Modal Awal (*Buka Shift*).
  3. Kasir memproses transaksi *offline* (pelanggan langsung). Sistem memotong stok BOM otomatis setiap ada pesanan.
  4. Selesai *shift*, Kasir menghitung uang laci dan menginput totalnya (*Tutup Shift*).
  5. Sistem men-*generate* laporan penjualan harian secara otomatis.
  6. Manager memeriksa keakuratan *dashboard* hasil laporan tersebut.

## 4. Manajemen Inventaris & Supply Chain
Alur permintaan penyetokan ulang (*restock*) dari cabang ke kantor pusat.
- **Aktor:** Branch Manager, Sistem Backend, Super Admin
- **Alur Utama:**
  1. Manager memantau stok, (opsional) melakukan *Stock Opname* jika ada selisih.
  2. Manager mengajukan *Request Restock* bahan baku ke Pusat.
  3. Admin Pusat me-*review* permohonan:
     - Jika **Ditolak**: Admin memberi alasan, Manager menerima notifikasi.
     - Jika **Disetujui**: Admin menekan "Approve".
  4. Sistem mengeksekusi transfer stok otomatis (Gudang Pusat berkurang, Cabang bertambah).
  5. Sistem mengirim notifikasi "Barang Dikirim" ke Manager.

## 5. Master Data & Konfigurasi Global
Alur pengaturan utama sistem (*Command Center*) yang dilakukan oleh kantor pusat.
- **Aktor:** Super Admin, Sistem Backend
- **Alur Utama:**
  1. Admin login dan masuk ke *Command Center*.
  2. Melalui *Menu Utama*, Admin bisa berulang kali memilih berbagai cabang konfigurasi:
     - Kelola Cabang
     - Kelola Akun & Assign Manager
     - Kelola Master Bahan Baku
     - Kelola Master Produk Jadi
     - Kelola BOM (Resep)
     - Kelola Gudang Pusat
  3. Admin juga bisa melihat *Dashboard Analitik Global* yang menampilkan agregat performa seluruh cabang.
