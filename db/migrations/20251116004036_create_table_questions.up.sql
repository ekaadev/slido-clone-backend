CREATE TABLE questions (
                           id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                           room_id BIGINT UNSIGNED NOT NULL,
                           participant_id BIGINT UNSIGNED NOT NULL,
                           content TEXT NOT NULL,
                           upvote_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Denormalized upvote counter',
                           status ENUM('pending', 'answered', 'highlighted') NOT NULL DEFAULT 'pending',
                           is_validated_by_presenter BOOLEAN NOT NULL DEFAULT FALSE,
                           xp_awarded INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'Total XP earned from this question',
                           created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                           FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
                           FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE,

                           INDEX idx_questions_room (room_id),
                           INDEX idx_questions_participant (participant_id),
                           INDEX idx_questions_status (status),
                           INDEX idx_questions_upvote_count (upvote_count DESC),
                           INDEX idx_questions_validated (is_validated_by_presenter),
                           INDEX idx_questions_created_at (created_at DESC),

                           FULLTEXT INDEX ft_questions_content (content)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;