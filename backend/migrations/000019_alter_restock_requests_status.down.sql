ALTER TABLE restock_requests MODIFY COLUMN status ENUM('Pending', 'Approved', 'Rejected') NOT NULL DEFAULT 'Pending';
