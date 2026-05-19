CREATE TABLE IF NOT EXISTS `password_resets` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `email` VARCHAR(150) NOT NULL,
  `token` VARCHAR(255) NOT NULL,
  `expires_at` TIMESTAMP NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT `fk_pw_resets_users`
    FOREIGN KEY (`email`)
    REFERENCES `users` (`email`)
    ON DELETE CASCADE
);