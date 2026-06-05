ALTER TABLE `users` ADD COLUMN `is_verified` BOOLEAN DEFAULT FALSE;
UPDATE `users` SET `is_verified` = TRUE;

CREATE TABLE IF NOT EXISTS `email_verifications` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `email` VARCHAR(100) NOT NULL,
  `otp_code` VARCHAR(10) NOT NULL,
  `expires_at` DATETIME NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE `transactions` ADD COLUMN `customer_name` VARCHAR(100) NULL;
