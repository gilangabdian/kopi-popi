CREATE TABLE IF NOT EXISTS `shifts` (
  `id` VARCHAR(36) PRIMARY KEY COMMENT 'UUID',
  `branch_id` INT NOT NULL,
  `cashier_id` VARCHAR(36) NOT NULL COMMENT 'UUID',
  `start_time` TIMESTAMP NOT NULL,
  `end_time` TIMESTAMP NULL,
  `starting_cash` DECIMAL(12,2) NOT NULL,
  `expected_cash` DECIMAL(12,2) NOT NULL,
  `actual_cash` DECIMAL(12,2) NULL,
  `status` ENUM('Open', 'Closed', 'Verified') NOT NULL DEFAULT 'Open',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT `fk_shifts_branches`
    FOREIGN KEY (`branch_id`)
    REFERENCES `branches` (`id`),

  CONSTRAINT `fk_shifts_users`
    FOREIGN KEY (`cashier_id`)
    REFERENCES `users` (`id`)
);