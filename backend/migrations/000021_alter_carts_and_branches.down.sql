ALTER TABLE `branches` DROP COLUMN `is_accepting_orders`;
ALTER TABLE `branches` DROP COLUMN `closing_time`;
ALTER TABLE `branches` DROP COLUMN `opening_time`;

ALTER TABLE `carts` DROP COLUMN `cart_name`;
ALTER TABLE `carts` MODIFY `customer_id` VARCHAR(36) NOT NULL COMMENT 'UUID';
