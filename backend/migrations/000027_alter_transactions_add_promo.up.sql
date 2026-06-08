ALTER TABLE transactions 
ADD COLUMN promo_code VARCHAR(50) NULL AFTER total_amount,
ADD COLUMN discount_amount DECIMAL(15, 2) NOT NULL DEFAULT 0 AFTER promo_code,
ADD COLUMN final_amount DECIMAL(15, 2) NOT NULL DEFAULT 0 AFTER discount_amount;

UPDATE transactions SET final_amount = total_amount;
