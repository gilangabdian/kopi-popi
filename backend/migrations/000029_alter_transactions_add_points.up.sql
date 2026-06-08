ALTER TABLE `transactions` 
ADD COLUMN `points_redeemed` INT NOT NULL DEFAULT 0,
ADD COLUMN `points_earned` INT NOT NULL DEFAULT 0;
