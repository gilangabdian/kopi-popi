CREATE TABLE IF NOT EXISTS `materials` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `category_id` INT NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `unit` VARCHAR(20) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT `fk_materials_categories`
    FOREIGN KEY (`category_id`)
    REFERENCES `categories` (`id`)
);