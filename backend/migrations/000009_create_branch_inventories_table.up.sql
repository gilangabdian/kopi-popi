CREATE TABLE IF NOT EXISTS `branch_inventories` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `branch_id` INT NOT NULL,
  `material_id` INT NOT NULL,
  `quantity` DECIMAL(12,2) NOT NULL DEFAULT 0,
  `last_updated` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT `fk_inv_branches`
    FOREIGN KEY (`branch_id`)
    REFERENCES `branches` (`id`),

  CONSTRAINT `fk_inv_materials`
    FOREIGN KEY (`material_id`)
    REFERENCES `materials` (`id`)
);