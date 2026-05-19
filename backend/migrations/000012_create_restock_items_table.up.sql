CREATE TABLE IF NOT EXISTS `restock_items` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `request_id` VARCHAR(36) NOT NULL COMMENT 'UUID',
  `material_id` INT NOT NULL,
  `quantity_requested` DECIMAL(12,2) NOT NULL,

  CONSTRAINT `fk_restockitems_req`
    FOREIGN KEY (`request_id`)
    REFERENCES `restock_requests` (`id`)
    ON DELETE CASCADE,

  CONSTRAINT `fk_restockitems_mat`
    FOREIGN KEY (`material_id`)
    REFERENCES `materials` (`id`),

  UNIQUE (`request_id`, `material_id`)
);