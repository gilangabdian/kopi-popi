CREATE TABLE IF NOT EXISTS `restock_requests` (
  `id` VARCHAR(36) PRIMARY KEY COMMENT 'UUID',
  `branch_id` INT NOT NULL,
  `requested_by` VARCHAR(36) NOT NULL COMMENT 'UUID',
  `status` ENUM('Pending', 'Approved', 'Rejected') NOT NULL DEFAULT 'Pending',
  `reason` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT `fk_restockreq_branches`
    FOREIGN KEY (`branch_id`)
    REFERENCES `branches` (`id`),

  CONSTRAINT `fk_restockreq_users`
    FOREIGN KEY (`requested_by`)
    REFERENCES `users` (`id`)
);