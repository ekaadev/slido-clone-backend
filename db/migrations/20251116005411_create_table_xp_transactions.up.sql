CREATE TABLE xp_transactions (
                                 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                 participant_id BIGINT UNSIGNED NOT NULL,
                                 room_id BIGINT UNSIGNED NOT NULL,
                                 points INT NOT NULL COMMENT 'Can be positive or negative',
                                 source_type ENUM(
        'poll',
        'question_created',
        'upvote_received',
        'presenter_validated'
    ) NOT NULL,
                                 source_id BIGINT UNSIGNED NOT NULL COMMENT 'Polymorphic: ID from polls/questions/votes',
                                 created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                 FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE,
                                 FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,

                                 INDEX idx_xp_transactions_participant (participant_id),
                                 INDEX idx_xp_transactions_room (room_id),
                                 INDEX idx_xp_transactions_source (source_type, source_id),
                                 INDEX idx_xp_transactions_created_at (created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;