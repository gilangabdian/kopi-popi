CREATE TABLE IF NOT EXISTS `transactions` (
  `id` VARCHAR(36) PRIMARY KEY COMMENT 'UUID',
  `branch_id` INT NOT NULL,
  `customer_id` VARCHAR(36) NULL COMMENT 'UUID',
  `cashier_id` VARCHAR(36) NULL COMMENT 'UUID',
  `shift_id` VARCHAR(36) NULL COMMENT 'UUID',

  `order_type`
    ENUM(
      'Online_Pickup',
      'Offline_DineIn',
      'Offline_Takeaway'
    ) NOT NULL,

  `payment_method` VARCHAR(50) NOT NULL,

  `total_amount` DECIMAL(12,2) NOT NULL,

  `status`
    ENUM(
      'Waiting_Payment',
      'Paid',
      'Preparing',
      'Ready',
      'Completed',
      'Cancelled'
    ) NOT NULL DEFAULT 'Waiting_Payment',

  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT `fk_trans_branches`
    FOREIGN KEY (`branch_id`)
    REFERENCES `branches` (`id`),

  CONSTRAINT `fk_trans_users`
    FOREIGN KEY (`customer_id`)
    REFERENCES `users` (`id`),

  CONSTRAINT `fk_trans_cashier`
    FOREIGN KEY (`cashier_id`)
    REFERENCES `users` (`id`),

  CONSTRAINT `fk_trans_shifts`
    FOREIGN KEY (`shift_id`)
    REFERENCES `shifts` (`id`)
);