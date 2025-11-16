CREATE TABLE rooms (
                       id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                       room_code VARCHAR(20) NOT NULL UNIQUE,
                       title VARCHAR(255) NOT NULL,
                       presenter_id BIGINT UNSIGNED NOT NULL,
                       status ENUM('active', 'closed') NOT NULL DEFAULT 'active',
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       closed_at TIMESTAMP NULL,

                       FOREIGN KEY (presenter_id) REFERENCES users(id) ON DELETE CASCADE,

                       INDEX idx_rooms_room_code (room_code),
                       INDEX idx_rooms_status (status),
                       INDEX idx_rooms_presenter (presenter_id),
                       INDEX idx_rooms_created_at (created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;