ALTER TABLE restock_requests MODIFY COLUMN status ENUM('Pending', 'Approved', 'Rejected', 'Delivered') NOT NULL DEFAULT 'Pending';
