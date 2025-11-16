CREATE TABLE votes (
                       id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                       question_id BIGINT UNSIGNED NOT NULL,
                       participant_id BIGINT UNSIGNED NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                       FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
                       FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE,

                       UNIQUE KEY unique_vote_per_question (question_id, participant_id),

                       INDEX idx_votes_question (question_id),
                       INDEX idx_votes_participant (participant_id),
                       INDEX idx_votes_created_at (created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;