CREATE TABLE promos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE COMMENT 'Contoh: MANTAP20',
    title VARCHAR(100) NOT NULL COMMENT 'Judul Promo, misal: Diskon Gajian 20%',
    discount_type ENUM('PERCENTAGE', 'FIXED') NOT NULL,
    discount_value DECIMAL(15, 2) NOT NULL COMMENT 'Jika PERCENTAGE maka nilainya 1-100, jika FIXED maka nilai Rupiah',
    max_discount_amount DECIMAL(15, 2) COMMENT 'Maksimal potongan (khusus untuk tipe PERCENTAGE)',
    min_purchase_amount DECIMAL(15, 2) NOT NULL DEFAULT 0 COMMENT 'Minimal belanja untuk bisa pakai promo ini',
    valid_from DATETIME NOT NULL,
    valid_until DATETIME NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
