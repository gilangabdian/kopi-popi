ALTER TABLE `carts` MODIFY `customer_id` VARCHAR(36) NULL COMMENT 'UUID (NULL for Offline Cashier Carts)';
ALTER TABLE `carts` ADD COLUMN `cart_name` VARCHAR(100) NULL COMMENT 'E.g., Meja 4 / Mas Baju Merah';

ALTER TABLE `branches` ADD COLUMN `opening_time` TIME NULL COMMENT 'Format HH:MM:SS';
ALTER TABLE `branches` ADD COLUMN `closing_time` TIME NULL COMMENT 'Format HH:MM:SS';
ALTER TABLE `branches` ADD COLUMN `is_accepting_orders` BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'Store Live Status';
