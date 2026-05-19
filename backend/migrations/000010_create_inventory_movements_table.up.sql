CREATE TABLE IF NOT EXISTS `inventory_movements` (
  `id` VARCHAR(36) PRIMARY KEY,
  `branch_id` INT NOT NULL,
  `material_id` INT NOT NULL,
  `movement_type` ENUM('IN', 'OUT', 'ADJUSTMENT') NOT NULL,
  `quantity` DECIMAL(12,2) NOT NULL,
  `description` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT `fk_invmov_branches`
    FOREIGN KEY (`branch_id`)
    REFERENCES `branches` (`id`),

  CONSTRAINT `fk_invmov_materials`
    FOREIGN KEY (`material_id`)
    REFERENCES `materials` (`id`)
);