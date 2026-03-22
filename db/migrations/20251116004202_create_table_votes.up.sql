CREATE TABLE votes (
    id BIGSERIAL PRIMARY KEY,
    question_id BIGINT NOT NULL,
    participant_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_votes_question FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    CONSTRAINT fk_votes_participant FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE,
    CONSTRAINT unique_vote_per_question UNIQUE (question_id, participant_id)
);

CREATE INDEX idx_votes_question ON votes (question_id);
CREATE INDEX idx_votes_participant ON votes (participant_id);
CREATE INDEX idx_votes_created_at ON votes (created_at DESC);
