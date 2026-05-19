CREATE TABLE IF NOT EXISTS `product_boms` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `product_id` INT NOT NULL,
  `material_id` INT NOT NULL,
  `quantity_needed` DECIMAL(10,2) NOT NULL,

  CONSTRAINT `fk_boms_products`
    FOREIGN KEY (`product_id`)
    REFERENCES `products` (`id`)
    ON DELETE CASCADE,

  CONSTRAINT `fk_boms_materials`
    FOREIGN KEY (`material_id`)
    REFERENCES `materials` (`id`)
);