CREATE TABLE participants (
                              id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                              room_id BIGINT UNSIGNED NOT NULL,
                              user_id BIGINT UNSIGNED NULL COMMENT 'NULL for anonymous participants',
                              display_name VARCHAR(100) NOT NULL,
                              xp_score INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Denormalized total XP',
                              is_anonymous BOOLEAN NOT NULL DEFAULT TRUE,
                              joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                              FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
                              FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,

                              UNIQUE KEY unique_user_room (room_id, user_id),

                              INDEX idx_participants_room (room_id),
                              INDEX idx_participants_xp_score (xp_score DESC),
                              INDEX idx_participants_joined_at (joined_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;