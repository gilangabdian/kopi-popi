CREATE TABLE IF NOT EXISTS `users` (
  `id` VARCHAR(36) PRIMARY KEY COMMENT 'UUID',
  `role_id` INT NOT NULL,
  `branch_id` INT NULL COMMENT 'Null for Admin & Customer',
  `name` VARCHAR(100) NOT NULL,
  `email` VARCHAR(150) UNIQUE NOT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `phone` VARCHAR(20),
  `profile_picture` VARCHAR(255) DEFAULT 'default-profile.png',
  `is_active` BOOLEAN NOT NULL DEFAULT TRUE,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT `fk_users_roles`
    FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),

  CONSTRAINT `fk_users_branches`
    FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`)
);