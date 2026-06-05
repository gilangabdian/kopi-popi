ALTER TABLE `transactions` DROP COLUMN `customer_name`;

DROP TABLE IF EXISTS `email_verifications`;

ALTER TABLE `users` DROP COLUMN `is_verified`;
