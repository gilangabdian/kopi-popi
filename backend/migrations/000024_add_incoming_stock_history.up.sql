CREATE TABLE incoming_stocks (
    id CHAR(36) PRIMARY KEY,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE incoming_stock_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    incoming_stock_id CHAR(36) NOT NULL,
    material_id INT NOT NULL,
    quantity DECIMAL(12,2) NOT NULL,
    supplier_name VARCHAR(255),
    supplier_phone VARCHAR(50),
    FOREIGN KEY (incoming_stock_id) REFERENCES incoming_stocks(id) ON DELETE CASCADE,
    FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE
);
