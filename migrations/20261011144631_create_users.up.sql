CREATE TABLE
    `users` (
        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(100) NOT NULL,
        `email` VARCHAR(150) NOT NULL,
        `is_email_verified` BOOLEAN NOT NULL DEFAULT FALSE,
        `email_verified_at` TIMESTAMP NULL DEFAULT NULL,
        `password_hash` VARCHAR(255) NOT NULL,
        `phone` VARCHAR(20) DEFAULT NULL,
        `address` TEXT DEFAULT NULL,
        `role` VARCHAR(30) NOT NULL DEFAULT 'user',
        `is_active` TINYINT (1) NOT NULL DEFAULT 1,
        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` TIMESTAMP DEFAULT NULL,
        PRIMARY KEY (`id`),
        UNIQUE KEY `users_email_unique` (`email`),
        INDEX `idx_users_role` (`role`),
        INDEX `idx_users_is_active` (`is_active`),
        INDEX `idx_users_deleted_at` (`deleted_at`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;