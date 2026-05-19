CREATE TABLE IF NOT EXISTS `carts` (
  `id` VARCHAR(36) PRIMARY KEY COMMENT 'UUID',
  `customer_id` VARCHAR(36) NOT NULL COMMENT 'UUID',
  `branch_id` INT NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT `fk_carts_users`
    FOREIGN KEY (`customer_id`)
    REFERENCES `users` (`id`)
    ON DELETE CASCADE,

  CONSTRAINT `fk_carts_branches`
    FOREIGN KEY (`branch_id`)
    REFERENCES `branches` (`id`)
);