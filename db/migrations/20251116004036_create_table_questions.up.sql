CREATE TABLE questions (
    id BIGSERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL,
    participant_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    upvote_count INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'answered', 'highlighted')),
    is_validated_by_presenter BOOLEAN NOT NULL DEFAULT FALSE,
    xp_awarded INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_questions_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    CONSTRAINT fk_questions_participant FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE
);

CREATE INDEX idx_questions_room ON questions (room_id);
CREATE INDEX idx_questions_participant ON questions (participant_id);
CREATE INDEX idx_questions_status ON questions (status);
CREATE INDEX idx_questions_upvote_count ON questions (upvote_count DESC);
CREATE INDEX idx_questions_validated ON questions (is_validated_by_presenter);
CREATE INDEX idx_questions_created_at ON questions (created_at DESC);
