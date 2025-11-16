CREATE TABLE poll_options (
                              id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                              poll_id BIGINT UNSIGNED NOT NULL,
                              option_text VARCHAR(255) NOT NULL,
                              vote_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Denormalized vote counter',
                              `order` TINYINT UNSIGNED NOT NULL,

                              FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE,

                              INDEX idx_poll_options_poll (poll_id),
                              INDEX idx_poll_options_order (poll_id, `order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;