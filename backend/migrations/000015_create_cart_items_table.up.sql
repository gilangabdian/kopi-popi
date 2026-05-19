CREATE TABLE IF NOT EXISTS `cart_items` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `cart_id` VARCHAR(36) NOT NULL COMMENT 'UUID',
  `product_id` INT NOT NULL,
  `quantity` INT NOT NULL DEFAULT 1,
  `notes` VARCHAR(255),

  CONSTRAINT `fk_cartitems_carts`
    FOREIGN KEY (`cart_id`)
    REFERENCES `carts` (`id`)
    ON DELETE CASCADE,

  CONSTRAINT `fk_cartitems_products`
    FOREIGN KEY (`product_id`)
    REFERENCES `products` (`id`),

  UNIQUE (`cart_id`, `product_id`)
);