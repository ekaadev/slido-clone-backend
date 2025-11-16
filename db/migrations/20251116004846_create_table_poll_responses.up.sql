CREATE TABLE poll_responses (
                                id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                poll_id BIGINT UNSIGNED NOT NULL,
                                participant_id BIGINT UNSIGNED NOT NULL,
                                poll_option_id BIGINT UNSIGNED NOT NULL,
                                created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE,
                                FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE,
                                FOREIGN KEY (poll_option_id) REFERENCES poll_options(id) ON DELETE CASCADE,

                                UNIQUE KEY unique_poll_response (poll_id, participant_id),

                                INDEX idx_poll_responses_poll (poll_id),
                                INDEX idx_poll_responses_participant (participant_id),
                                INDEX idx_poll_responses_option (poll_option_id),
                                INDEX idx_poll_responses_created_at (created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
