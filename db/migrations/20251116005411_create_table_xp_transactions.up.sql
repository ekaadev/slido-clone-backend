CREATE TABLE xp_transactions (
    id BIGSERIAL PRIMARY KEY,
    participant_id BIGINT NOT NULL,
    room_id BIGINT NOT NULL,
    points INT NOT NULL,
    source_type VARCHAR(30) NOT NULL CHECK (source_type IN ('poll', 'question_created', 'upvote_received', 'presenter_validated')),
    source_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_xp_transactions_participant FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE,
    CONSTRAINT fk_xp_transactions_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE INDEX idx_xp_transactions_participant ON xp_transactions (participant_id);
CREATE INDEX idx_xp_transactions_room ON xp_transactions (room_id);
CREATE INDEX idx_xp_transactions_source ON xp_transactions (source_type, source_id);
CREATE INDEX idx_xp_transactions_created_at ON xp_transactions (created_at DESC);
