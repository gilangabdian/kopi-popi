CREATE TABLE `blogs` (
    `id` CHAR(36) PRIMARY KEY COMMENT 'UUID',
    `author_id` CHAR(36) NOT NULL COMMENT 'Relasi ke users (Admin)',
    `title` VARCHAR(255) NOT NULL,
    `content` TEXT NOT NULL,
    `estimated_read_time_mins` INT NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`author_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);
