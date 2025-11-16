CREATE TABLE polls (
                       id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                       room_id BIGINT UNSIGNED NOT NULL,
                       question TEXT NOT NULL,
                       status ENUM('draft', 'active', 'closed') NOT NULL DEFAULT 'draft',
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       activated_at TIMESTAMP NULL,
                       closed_at TIMESTAMP NULL,

                       FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,

                       INDEX idx_polls_room (room_id),
                       INDEX idx_polls_status (status),
                       INDEX idx_polls_room_status (room_id, status),
                       INDEX idx_polls_activated_at (activated_at DESC),
                       INDEX idx_polls_created_at (created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;