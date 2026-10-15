CREATE TABLE
    `refresh_tokens` (
        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
        `user_id` BIGINT UNSIGNED NOT NULL,
        `token_hash` VARCHAR(255) NOT NULL,
        `user_agent` VARCHAR(255) NULL,
        `ip_address` VARCHAR(45) NULL,
        `is_revoked` BOOLEAN DEFAULT FALSE NOT NULL,
        `expires_at` TIMESTAMP NOT NULL,
        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        -- Foreign Key
        CONSTRAINT `fk_refresh_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
    );

CREATE UNIQUE INDEX `idx_refresh_tokens_token_hash` ON `refresh_tokens` (`token_hash`);

CREATE INDEX `idx_refresh_tokens_user_id` ON `refresh_tokens` (`user_id`);