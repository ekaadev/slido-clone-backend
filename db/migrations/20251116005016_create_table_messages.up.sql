CREATE TABLE messages (
                          id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                          room_id BIGINT UNSIGNED NOT NULL,
                          participant_id BIGINT UNSIGNED NOT NULL,
                          content TEXT NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                          FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
                          FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE,

                          INDEX idx_messages_room (room_id),
                          INDEX idx_messages_participant (participant_id),
                          INDEX idx_messages_room_created (room_id, created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;