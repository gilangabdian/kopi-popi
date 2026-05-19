CREATE TABLE IF NOT EXISTS `transaction_details` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `transaction_id` VARCHAR(36) NOT NULL COMMENT 'UUID',
  `product_id` INT NOT NULL,
  `quantity` INT NOT NULL,
  `subtotal` DECIMAL(12,2) NOT NULL,
  `notes` VARCHAR(255),

  CONSTRAINT `fk_tdet_trans`
    FOREIGN KEY (`transaction_id`)
    REFERENCES `transactions` (`id`)
    ON DELETE CASCADE,

  CONSTRAINT `fk_tdet_products`
    FOREIGN KEY (`product_id`)
    REFERENCES `products` (`id`)
);